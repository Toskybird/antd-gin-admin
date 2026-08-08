package scanengine

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
)

// NucleiCLI runs ProjectDiscovery nuclei as a subprocess.
type NucleiCLI struct {
	BinPath        string
	TemplatesDir   string
	Timeout        time.Duration
	ExtraArgs      []string
	LookPath       func(file string) (string, error) // injectable for tests; default exec.LookPath
	CommandContext func(ctx context.Context, name string, arg ...string) *exec.Cmd
}

// Available reports whether nuclei binary can be resolved.
func (n *NucleiCLI) Available() bool {
	_, err := n.resolveBin()
	return err == nil
}

func (n *NucleiCLI) resolveBin() (string, error) {
	look := n.LookPath
	if look == nil {
		look = exec.LookPath
	}
	bin := strings.TrimSpace(n.BinPath)
	if bin == "" {
		bin = "nuclei"
	}
	if strings.Contains(bin, "/") || strings.Contains(bin, `\`) {
		if _, err := os.Stat(bin); err != nil {
			return "", err
		}
		return bin, nil
	}
	if p, err := look(bin); err == nil {
		return p, nil
	}
	candidates := []string{}
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		candidates = append(candidates, filepath.Join(gopath, "bin", bin))
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, "go", "bin", bin))
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("nuclei binary %q not found in PATH or GOPATH/bin", bin)
}

// Scan executes nuclei against the target and maps JSONL findings.
func (n *NucleiCLI) Scan(ctx context.Context, req repointerfaces.ScanRequest) ([]repointerfaces.EngineFinding, error) {
	if strings.TrimSpace(req.EntryURL) == "" {
		return nil, fmt.Errorf("entry URL is required")
	}
	if len(req.EnabledRules) == 0 {
		return []repointerfaces.EngineFinding{}, nil
	}

	bin, err := n.resolveBin()
	if err != nil {
		return nil, fmt.Errorf("nuclei not found: %w", err)
	}

	timeout := n.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	args := []string{
		"-u", req.EntryURL,
		"-jsonl",
		"-silent",
		"-nc",
		"-ot", // omit template blob
		"-or", // omit raw request/response
		"-id", strings.Join(req.EnabledRules, ","),
	}
	if sev := severityForPolicy(req.Policy); sev != "" {
		args = append(args, "-severity", sev)
	}
	if dir := strings.TrimSpace(n.TemplatesDir); dir != "" {
		args = append(args, "-t", resolveTemplatesDir(dir))
	}
	args = append(args, n.ExtraArgs...)

	newCmd := n.CommandContext
	if newCmd == nil {
		newCmd = exec.CommandContext
	}
	cmd := newCmd(runCtx, bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	// nuclei exits non-zero in some versions when findings exist or templates missing;
	// still try to parse stdout. Fail hard only if no usable output and process error.
	findings, parseErr := ParseNucleiJSONL(stdout.Bytes())
	if parseErr != nil {
		return nil, parseErr
	}
	if err != nil && len(findings) == 0 && runCtx.Err() != nil {
		return nil, fmt.Errorf("nuclei canceled or timed out: %w", runCtx.Err())
	}
	if err != nil && len(findings) == 0 && stdout.Len() == 0 {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("nuclei failed: %s", msg)
	}
	return findings, nil
}

func severityForPolicy(policy string) string {
	switch strings.TrimSpace(strings.ToLower(policy)) {
	case "quick":
		return "critical,high"
	case "standard":
		return "critical,high,medium"
	case "deep":
		return "critical,high,medium,low,info"
	default:
		return "critical,high,medium"
	}
}

type nucleiJSONLEvent struct {
	TemplateID        string   `json:"template-id"`
	Info              struct {
		Name        string `json:"name"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
	} `json:"info"`
	MatchedAt         string   `json:"matched-at"`
	Host              string   `json:"host"`
	Type              string   `json:"type"`
	Matcher           string   `json:"matcher-name"`
	ExtractedResults  []string `json:"extracted-results"`
}

// ParseNucleiJSONL converts nuclei -jsonl stdout into EngineFinding rows.
func ParseNucleiJSONL(raw []byte) ([]repointerfaces.EngineFinding, error) {
	out := make([]repointerfaces.EngineFinding, 0)
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	// Increase buffer for large lines (though we omit raw/template).
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var ev nucleiJSONLEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			return nil, fmt.Errorf("nuclei jsonl line %d: %w", lineNo, err)
		}
		ruleCode := strings.TrimSpace(ev.TemplateID)
		if ruleCode == "" {
			continue
		}
		severity := strings.ToLower(strings.TrimSpace(ev.Info.Severity))
		if severity == "" {
			severity = "info"
		}
		title := strings.TrimSpace(ev.Info.Name)
		if title == "" {
			title = ruleCode
		}
		location := strings.TrimSpace(ev.MatchedAt)
		if location == "" {
			location = strings.TrimSpace(ev.Host)
		}
		evidence := strings.TrimSpace(ev.Matcher)
		if evidence == "" && len(ev.ExtractedResults) > 0 {
			evidence = strings.Join(ev.ExtractedResults, ", ")
		}
		if evidence == "" {
			evidence = strings.TrimSpace(ev.Type)
		}
		out = append(out, repointerfaces.EngineFinding{
			RuleCode:    ruleCode,
			Severity:    severity,
			Title:       title,
			Description: strings.TrimSpace(ev.Info.Description),
			Evidence:    evidence,
			Location:    location,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// ResolveEngine picks Nuclei when available, otherwise the provided fallback.
func ResolveEngine(nuclei *NucleiCLI, fallback repointerfaces.ScanEngineAdapter) (repointerfaces.ScanEngineAdapter, string) {
	if nuclei != nil && nuclei.Available() {
		return nuclei, "nuclei"
	}
	if fallback == nil {
		fallback = &noopEngine{}
	}
	return fallback, "fake"
}

type noopEngine struct{}

func (noopEngine) Scan(_ context.Context, _ repointerfaces.ScanRequest) ([]repointerfaces.EngineFinding, error) {
	return []repointerfaces.EngineFinding{}, nil
}

func resolveTemplatesDir(dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}
	if st, err := os.Stat(dir); err == nil && st.IsDir() {
		if abs, err := filepath.Abs(dir); err == nil {
			return abs
		}
		return dir
	}
	// When server is started from repo root: backend/configs/nuclei-templates
	alt := filepath.Join("backend", dir)
	if st, err := os.Stat(alt); err == nil && st.IsDir() {
		if abs, err := filepath.Abs(alt); err == nil {
			return abs
		}
		return alt
	}
	return dir
}
