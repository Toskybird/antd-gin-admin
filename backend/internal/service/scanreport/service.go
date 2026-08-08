package scanreport

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"antd-gin-admin-backend/internal/apperrors"
	"antd-gin-admin-backend/internal/model/entity"
	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
	serviceinterfaces "antd-gin-admin-backend/internal/service/interfaces"
)

const (
	jobSucceeded = "succeeded"
	jobCancelled = "cancelled"
	formatHTML   = "html"
)

type Service struct {
	reports  repointerfaces.ScanReportRepository
	jobs     repointerfaces.ScanJobRepository
	findings repointerfaces.FindingRepository
	storage  repointerfaces.ObjectStorage
}

func New(
	reports repointerfaces.ScanReportRepository,
	jobs repointerfaces.ScanJobRepository,
	findings repointerfaces.FindingRepository,
	storage repointerfaces.ObjectStorage,
) serviceinterfaces.ScanReportService {
	return &Service{reports: reports, jobs: jobs, findings: findings, storage: storage}
}

func (s *Service) Generate(ctx context.Context, jobCode, creatorUserCode string) (*entity.ScanReport, error) {
	jobCode = strings.TrimSpace(jobCode)
	if jobCode == "" {
		return nil, apperrors.New(http.StatusBadRequest, "job_code_required", "必须提供任务编码")
	}
	job, err := s.jobs.GetByCode(ctx, jobCode)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, apperrors.New(http.StatusNotFound, "job_not_found", "扫描任务不存在")
	}
	if job.Status != jobSucceeded && job.Status != jobCancelled {
		return nil, apperrors.New(http.StatusBadRequest, "job_not_reportable", "仅成功或已取消的任务可生成报告")
	}

	version, err := s.reports.NextVersion(ctx, jobCode)
	if err != nil {
		return nil, err
	}
	findings, err := s.findings.ListByJob(ctx, jobCode)
	if err != nil {
		return nil, err
	}
	body := renderHTML(job, findings, version)
	reportCode, err := generateReportCode()
	if err != nil {
		return nil, err
	}
	objectKey := fmt.Sprintf("scan-reports/%s/v%d-%s.html", jobCode, version, reportCode)
	if err := s.storage.Put(ctx, objectKey, []byte(body), "text/html; charset=utf-8"); err != nil {
		return nil, err
	}
	report := &entity.ScanReport{
		ReportCode: reportCode,
		JobCode:    jobCode,
		Version:    version,
		ObjectKey:  objectKey,
		Format:     formatHTML,
		CreatedBy:  creatorUserCode,
	}
	if err := s.reports.Create(ctx, report); err != nil {
		return nil, err
	}
	return s.reports.GetByCode(ctx, reportCode)
}

func (s *Service) ListWithScope(ctx context.Context, page, pageSize int, jobCode string, scope *serviceinterfaces.UserDataScope) ([]*entity.ScanReport, int64, error) {
	filter := repointerfaces.ScanReportListFilter{
		JobCode: strings.TrimSpace(jobCode),
		All:     scope == nil || scope.All,
	}
	if !filter.All && scope != nil {
		codes, err := s.visibleJobCodes(ctx, scope)
		if err != nil {
			return nil, 0, err
		}
		filter.JobCodes = codes
	}
	return s.reports.List(ctx, page, pageSize, filter)
}

func (s *Service) Download(ctx context.Context, reportCode string, scope *serviceinterfaces.UserDataScope) ([]byte, string, string, error) {
	report, err := s.reports.GetByCode(ctx, reportCode)
	if err != nil {
		return nil, "", "", err
	}
	if report == nil {
		return nil, "", "", nil
	}
	if scope != nil && !scope.All {
		job, err := s.jobs.GetByCode(ctx, report.JobCode)
		if err != nil {
			return nil, "", "", err
		}
		if job == nil || !jobVisible(job, scope) {
			return nil, "", "", nil
		}
	}
	content, err := s.storage.Get(ctx, report.ObjectKey)
	if err != nil {
		return nil, "", "", err
	}
	filename := fmt.Sprintf("%s-v%d.html", report.JobCode, report.Version)
	return content, "text/html; charset=utf-8", filename, nil
}

func (s *Service) visibleJobCodes(ctx context.Context, scope *serviceinterfaces.UserDataScope) ([]string, error) {
	jobs, _, err := s.jobs.List(ctx, 1, 10000, repointerfaces.ScanJobListFilter{
		All:       false,
		DeptCodes: append([]string{}, scope.DeptCodes...),
		CreatedBy: scope.UserCode,
	})
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0, len(jobs))
	for _, job := range jobs {
		codes = append(codes, job.JobCode)
	}
	return codes, nil
}

func jobVisible(job *entity.ScanJob, scope *serviceinterfaces.UserDataScope) bool {
	if scope == nil || scope.All {
		return true
	}
	for _, d := range scope.DeptCodes {
		if d == job.DeptCode {
			return true
		}
	}
	if scope.UserCode != "" && job.CreatedBy == scope.UserCode {
		return true
	}
	return false
}

func generateReportCode() (string, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return fmt.Sprintf("RPT-%s-%s", time.Now().Format("150405"), hex.EncodeToString(b[:])), nil
}

func renderHTML(job *entity.ScanJob, findings []*entity.Finding, version int) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><head><meta charset=\"utf-8\"><title>Scan Report ")
	b.WriteString(html.EscapeString(job.JobCode))
	b.WriteString("</title></head><body>")
	b.WriteString("<h1>扫描报告</h1>")
	b.WriteString("<p>任务：")
	b.WriteString(html.EscapeString(job.JobCode))
	b.WriteString("</p><p>入口：")
	b.WriteString(html.EscapeString(job.EntryURL))
	b.WriteString("</p><p>状态：")
	b.WriteString(html.EscapeString(job.Status))
	b.WriteString("</p><p>版本：")
	b.WriteString(fmt.Sprintf("%d", version))
	b.WriteString("</p><h2>发现项</h2><ul>")
	for _, f := range findings {
		b.WriteString("<li>[")
		b.WriteString(html.EscapeString(f.Severity))
		b.WriteString("] ")
		b.WriteString(html.EscapeString(f.Title))
		b.WriteString(" (")
		b.WriteString(html.EscapeString(f.RuleCode))
		b.WriteString(")</li>")
	}
	b.WriteString("</ul></body></html>")
	return b.String()
}
