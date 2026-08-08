import { request } from '@umijs/max';
import type { PageResult } from './asset';

export interface ScanJob {
  job_code: string;
  asset_code?: string;
  entry_url: string;
  dept_code: string;
  policy: 'quick' | 'standard' | 'deep';
  max_depth: number;
  max_pages: number;
  status: 'queued' | 'running' | 'succeeded' | 'failed' | 'cancelled';
  finding_count: number;
  temp_url: boolean;
  created_at: string;
  finished_at?: string;
}

export interface CreateScanJobRequest {
  source: 'asset' | 'temp_url';
  asset_code?: string;
  entry_url?: string;
  policy: 'quick' | 'standard' | 'deep';
  max_depth?: number;
  max_pages?: number;
}

export interface ListScanJobParams {
  page?: number;
  pageSize?: number;
  status?: string;
  policy?: string;
  keyword?: string;
}

/** 获取扫描任务列表 */
export async function getScanJobList(params?: ListScanJobParams) {
  return request<{
    code: number;
    data: PageResult<ScanJob>;
    message?: string;
  }>('/api/v1/scan/scan-jobs', {
    method: 'GET',
    params,
  });
}

/** 获取扫描任务详情 */
export async function getScanJob(jobCode: string) {
  return request<{
    code: number;
    data: ScanJob;
    message?: string;
  }>(`/api/v1/scan/scan-jobs/${jobCode}`, {
    method: 'GET',
  });
}

/** 创建扫描任务 */
export async function createScanJob(data: CreateScanJobRequest) {
  return request<{
    code: number;
    data: ScanJob;
    message?: string;
  }>('/api/v1/scan/scan-jobs', {
    method: 'POST',
    data,
  });
}

/** 取消扫描任务 */
export async function cancelScanJob(jobCode: string) {
  return request<{
    code: number;
    data: ScanJob;
    message?: string;
  }>(`/api/v1/scan/scan-jobs/${jobCode}/cancel`, {
    method: 'POST',
  });
}
