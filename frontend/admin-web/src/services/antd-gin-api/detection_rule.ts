import { request } from '@umijs/max';
import type { PageResult } from './asset';

export interface DetectionRule {
  rule_code: string;
  display_name: string;
  enabled: boolean;
  updated_at: string;
}

export interface ListDetectionRuleParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  enabled?: boolean;
}

/** 获取检测规则列表 */
export async function getDetectionRuleList(params?: ListDetectionRuleParams) {
  return request<{
    code: number;
    data: PageResult<DetectionRule>;
    message?: string;
  }>('/api/v1/scan/detection-rules', {
    method: 'GET',
    params,
  });
}

/** 同步检测规则 */
export async function syncDetectionRules() {
  return request<{
    code: number;
    data: { synced: number };
    message?: string;
  }>('/api/v1/scan/detection-rules/sync', {
    method: 'POST',
  });
}

/** 更新检测规则启用状态 */
export async function updateDetectionRuleStatus(
  ruleCode: string,
  enabled: boolean,
) {
  return request<{
    code: number;
    data: DetectionRule;
    message?: string;
  }>(`/api/v1/scan/detection-rules/${ruleCode}/status`, {
    method: 'PUT',
    data: { enabled },
  });
}
