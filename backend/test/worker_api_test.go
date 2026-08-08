package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"antd-gin-admin-backend/internal/model/dto"
	"antd-gin-admin-backend/internal/model/entity"
	"antd-gin-admin-backend/internal/repository/interfaces"
	repoMock "antd-gin-admin-backend/internal/repository/mock"
	detectionrulesvc "antd-gin-admin-backend/internal/service/detectionrule"
	scanjobsvc "antd-gin-admin-backend/internal/service/scanjob"
	"antd-gin-admin-backend/internal/service/worker"
)

func TestWorkerFakeEngineWritesFindings(t *testing.T) {
	jobRepo := repoMock.NewScanJobMockRepository()
	assetRepo := repoMock.NewAssetMockRepository()
	queue := repoMock.NewMemoryScanJobQueue()
	findingRepo := repoMock.NewFindingMockRepository()
	ruleRepo := repoMock.NewDetectionRuleMockRepository()
	ruleSvc := detectionrulesvc.New(ruleRepo, &repoMock.StaticDetectionRuleSource{Rules: []entity.DetectionRule{
		{RuleCode: "exposed-panels", DisplayName: "Exposed Admin Panels"},
		{RuleCode: "tech-detect", DisplayName: "Technology Detection"},
	}})
	_, err := ruleSvc.Sync(context.Background())
	require.NoError(t, err)
	_, err = ruleSvc.SetEnabled(context.Background(), "tech-detect", false)
	require.NoError(t, err)

	require.NoError(t, assetRepo.Create(context.Background(), &entity.Asset{
		AssetCode: "AST-1", Name: "门户", RootURL: "https://portal.example.com", DeptCode: "D1", Status: "active",
	}))
	jobSvc := scanjobsvc.New(jobRepo, assetRepo, queue)
	job, err := jobSvc.Create(context.Background(), &dto.CreateScanJobRequest{
		Source: "asset", AssetCode: "AST-1", Policy: "standard",
	}, "U1", "D1")
	require.NoError(t, err)

	engine := &repoMock.FakeScanEngine{Findings: []interfaces.EngineFinding{
		{RuleCode: "exposed-panels", Severity: "high", Title: "Open admin", Evidence: "200", Location: "/admin"},
		{RuleCode: "tech-detect", Severity: "info", Title: "Nginx", Evidence: "header", Location: "/"},
	}}
	runner := worker.NewRunner(jobRepo, findingRepo, queue, engine, ruleSvc)
	ok, err := runner.ProcessNext(context.Background())
	require.NoError(t, err)
	require.True(t, ok)

	got, err := jobRepo.GetByCode(context.Background(), job.JobCode)
	require.NoError(t, err)
	assert.Equal(t, "succeeded", got.Status)
	assert.Equal(t, 1, got.FindingCount)

	findings, err := findingRepo.ListByJob(context.Background(), job.JobCode)
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Equal(t, "exposed-panels", findings[0].RuleCode)
}

func TestWorkerSkipsCancelledJob(t *testing.T) {
	jobRepo := repoMock.NewScanJobMockRepository()
	findingRepo := repoMock.NewFindingMockRepository()
	queue := repoMock.NewMemoryScanJobQueue()
	ruleRepo := repoMock.NewDetectionRuleMockRepository()
	ruleSvc := detectionrulesvc.New(ruleRepo, &repoMock.StaticDetectionRuleSource{Rules: []entity.DetectionRule{
		{RuleCode: "exposed-panels", DisplayName: "Exposed"},
	}})
	_, _ = ruleSvc.Sync(context.Background())
	engine := &repoMock.FakeScanEngine{Findings: []interfaces.EngineFinding{
		{RuleCode: "exposed-panels", Severity: "high", Title: "x"},
	}}

	require.NoError(t, jobRepo.Create(context.Background(), &entity.ScanJob{
		JobCode: "JOB-CANCEL", EntryURL: "https://x.example.com", Status: "cancelled", Policy: "quick",
	}))
	require.NoError(t, queue.Enqueue(context.Background(), "JOB-CANCEL"))

	runner := worker.NewRunner(jobRepo, findingRepo, queue, engine, ruleSvc)
	ok, err := runner.ProcessNext(context.Background())
	require.NoError(t, err)
	require.True(t, ok)

	got, _ := jobRepo.GetByCode(context.Background(), "JOB-CANCEL")
	assert.Equal(t, "cancelled", got.Status)
	count, _ := findingRepo.CountByJob(context.Background(), "JOB-CANCEL")
	assert.Equal(t, int64(0), count)
}
