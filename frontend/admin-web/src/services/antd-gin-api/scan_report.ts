import { request } from '@umijs/max';
import type { PageResult } from './asset';

export interface ScanReport {
  report_code: string;
  job_code: string;
  version: number;
  object_key: string;
  format: string;
  created_at: string;
}

export interface ListScanReportParams {
  page?: number;
  pageSize?: number;
  job_code?: string;
}

export interface CreateScanReportRequest {
  job_code: string;
}

/** 获取扫描报告列表 */
export async function getScanReportList(params?: ListScanReportParams) {
  return request<{
    code: number;
    data: PageResult<ScanReport>;
    message?: string;
  }>('/api/v1/scan/scan-reports', {
    method: 'GET',
    params,
  });
}

/** 生成扫描报告 */
export async function createScanReport(data: CreateScanReportRequest) {
  return request<{
    code: number;
    data: ScanReport;
    message?: string;
  }>('/api/v1/scan/scan-reports', {
    method: 'POST',
    data,
  });
}

/** 下载扫描报告（HTML blob） */
export async function downloadScanReport(reportCode: string) {
  return request<Blob>(`/api/v1/scan/scan-reports/${reportCode}/download`, {
    method: 'GET',
    responseType: 'blob',
  });
}
