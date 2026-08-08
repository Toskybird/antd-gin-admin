package scanengine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
)

func TestParseNucleiJSONL(t *testing.T) {
	raw := []byte(`{"template-id":"http-missing-security-headers","info":{"name":"Missing Security Headers","severity":"info","description":"headers missing"},"matched-at":"https://example.com","type":"http","matcher-name":"x-frame-options"}
{"template-id":"exposed-panels","info":{"name":"Exposed Admin","severity":"high","description":"panel"},"matched-at":"https://example.com/admin","type":"http"}
`)
	findings, err := ParseNucleiJSONL(raw)
	require.NoError(t, err)
	require.Len(t, findings, 2)
	assert.Equal(t, "http-missing-security-headers", findings[0].RuleCode)
	assert.Equal(t, "info", findings[0].Severity)
	assert.Equal(t, "Missing Security Headers", findings[0].Title)
	assert.Equal(t, "x-frame-options", findings[0].Evidence)
	assert.Equal(t, "https://example.com", findings[0].Location)
	assert.Equal(t, "high", findings[1].Severity)
	assert.Equal(t, "https://example.com/admin", findings[1].Location)
}

func TestParseNucleiJSONLEmpty(t *testing.T) {
	findings, err := ParseNucleiJSONL(nil)
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func TestNucleiCLISkipsWhenNoEnabledRules(t *testing.T) {
	n := &NucleiCLI{BinPath: "/nonexistent/nuclei"}
	out, err := n.Scan(context.Background(), repointerfaces.ScanRequest{
		EntryURL:     "https://example.com",
		EnabledRules: nil,
	})
	require.NoError(t, err)
	assert.Empty(t, out)
}

// local stub to avoid importing mock cycle
type stubEngine struct {
	name string
}

func (s *stubEngine) Scan(context.Context, repointerfaces.ScanRequest) ([]repointerfaces.EngineFinding, error) {
	return []repointerfaces.EngineFinding{{RuleCode: s.name}}, nil
}

func TestResolveEngine(t *testing.T) {
	fallback := &stubEngine{name: "fake"}
	n := &NucleiCLI{
		BinPath: "/nonexistent/nuclei-binary",
	}
	eng, kind := ResolveEngine(n, fallback)
	assert.Equal(t, "fake", kind)
	out, err := eng.Scan(context.Background(), repointerfaces.ScanRequest{})
	require.NoError(t, err)
	assert.Equal(t, "fake", out[0].RuleCode)

	n.BinPath = "nuclei"
	n.LookPath = func(string) (string, error) { return "/usr/local/bin/nuclei", nil }
	eng, kind = ResolveEngine(n, fallback)
	assert.Equal(t, "nuclei", kind)
	_, ok := eng.(*NucleiCLI)
	assert.True(t, ok)
}

func TestSeverityForPolicy(t *testing.T) {
	assert.Equal(t, "critical,high,medium,info", severityForPolicy("quick"))
	assert.Equal(t, "critical,high,medium,low,info", severityForPolicy("standard"))
	assert.Equal(t, "critical,high,medium,low,info", severityForPolicy("deep"))
}
