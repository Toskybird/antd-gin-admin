import { request } from '@umijs/max';
import type { PageResult } from './asset';

export interface Finding {
  finding_code: string;
  job_code: string;
  rule_code: string;
  rule_display_name: string;
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info';
  title: string;
  description: string;
  evidence: string;
  location: string;
  created_at: string;
}

export interface ListFindingParams {
  page?: number;
  pageSize?: number;
  job_code?: string;
  rule_code?: string;
  severity?: string;
  keyword?: string;
}

/** 获取发现项列表 */
export async function getFindingList(params?: ListFindingParams) {
  return request<{
    code: number;
    data: PageResult<Finding>;
    message?: string;
  }>('/api/v1/scan/findings', {
    method: 'GET',
    params,
  });
}

/** 获取发现项详情 */
export async function getFinding(findingCode: string) {
  return request<{
    code: number;
    data: Finding;
    message?: string;
  }>(`/api/v1/scan/findings/${findingCode}`, {
    method: 'GET',
  });
}
