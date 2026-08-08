package worker

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"antd-gin-admin-backend/internal/model/entity"
	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
	scanjobsvc "antd-gin-admin-backend/internal/service/scanjob"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

// Runner consumes queued scan jobs and writes findings via an engine adapter.
type Runner struct {
	jobs      repointerfaces.ScanJobRepository
	findings  repointerfaces.FindingRepository
	queue     repointerfaces.ScanJobQueue
	engine    repointerfaces.ScanEngineAdapter
	rules     serviceinterfaces.DetectionRuleService
}

func NewRunner(
	jobs repointerfaces.ScanJobRepository,
	findings repointerfaces.FindingRepository,
	queue repointerfaces.ScanJobQueue,
	engine repointerfaces.ScanEngineAdapter,
	rules serviceinterfaces.DetectionRuleService,
) *Runner {
	return &Runner{jobs: jobs, findings: findings, queue: queue, engine: engine, rules: rules}
}

// ProcessNext dequeues one job and processes it. Returns false if queue empty.
func (r *Runner) ProcessNext(ctx context.Context) (bool, error) {
	jobCode, ok, err := r.queue.Dequeue(ctx)
	if err != nil || !ok {
		return false, err
	}
	return true, r.ProcessJob(ctx, jobCode)
}

// ProcessJob runs one scan job by code.
func (r *Runner) ProcessJob(ctx context.Context, jobCode string) error {
	job, err := r.jobs.GetByCode(ctx, jobCode)
	if err != nil {
		return err
	}
	if job == nil {
		return nil
	}
	if job.Status == scanjobsvc.StatusCancelled {
		return nil
	}
	job.Status = scanjobsvc.StatusRunning
	if err := r.jobs.Update(ctx, job); err != nil {
		return err
	}

	// Re-check cancel before engine work.
	job, err = r.jobs.GetByCode(ctx, jobCode)
	if err != nil {
		return err
	}
	if job.Status == scanjobsvc.StatusCancelled {
		return nil
	}

	enabledRules := []string{}
	if r.rules != nil {
		enabledRules, err = r.rules.ListEnabledCodes(ctx)
		if err != nil {
			return err
		}
	}

	findings, err := r.engine.Scan(ctx, repointerfaces.ScanRequest{
		JobCode:      job.JobCode,
		EntryURL:     job.EntryURL,
		Policy:       job.Policy,
		MaxDepth:     job.MaxDepth,
		MaxPages:     job.MaxPages,
		EnabledRules: enabledRules,
	})
	if err != nil {
		job.Status = scanjobsvc.StatusFailed
		now := time.Now()
		job.FinishedAt = &now
		_ = r.jobs.Update(ctx, job)
		return err
	}

	// Cancel may have happened during scan.
	fresh, err := r.jobs.GetByCode(ctx, jobCode)
	if err != nil {
		return err
	}
	if fresh != nil && fresh.Status == scanjobsvc.StatusCancelled {
		return nil
	}

	for _, f := range findings {
		code, err := generateFindingCode()
		if err != nil {
			return err
		}
		if err := r.findings.Create(ctx, &entity.Finding{
			FindingCode: code,
			JobCode:     jobCode,
			RuleCode:    f.RuleCode,
			Severity:    f.Severity,
			Title:       f.Title,
			Description: f.Description,
			Evidence:    f.Evidence,
			Location:    f.Location,
		}); err != nil {
			return err
		}
	}
	count, err := r.findings.CountByJob(ctx, jobCode)
	if err != nil {
		return err
	}
	now := time.Now()
	job.Status = scanjobsvc.StatusSucceeded
	job.FindingCount = int(count)
	job.FinishedAt = &now
	return r.jobs.Update(ctx, job)
}

func generateFindingCode() (string, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("FND-%s-%s", time.Now().Format("150405"), hex.EncodeToString(b[:])), nil
}
