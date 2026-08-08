/** Asset status display */
export const ASSET_STATUS_TEXT: Record<string, string> = {
  active: '启用',
  disabled: '停用',
};

export const ASSET_STATUS_BADGE: Record<
  string,
  'success' | 'default' | 'error' | 'processing' | 'warning'
> = {
  active: 'success',
  disabled: 'default',
};

export function assetStatusText(status: string): string {
  return ASSET_STATUS_TEXT[status] ?? status;
}

/** Scan job status display */
export const SCAN_JOB_STATUS_TEXT: Record<string, string> = {
  queued: '排队中',
  running: '执行中',
  succeeded: '已完成',
  failed: '失败',
  cancelled: '已取消',
};

export const SCAN_JOB_STATUS_BADGE: Record<
  string,
  'success' | 'default' | 'error' | 'processing' | 'warning'
> = {
  queued: 'processing',
  running: 'processing',
  succeeded: 'success',
  failed: 'error',
  cancelled: 'default',
};

export function scanJobStatusText(status: string): string {
  return SCAN_JOB_STATUS_TEXT[status] ?? status;
}

export function isScanJobActive(status: string): boolean {
  return status === 'queued' || status === 'running';
}

/** Scan policy display */
export const SCAN_POLICY_TEXT: Record<string, string> = {
  quick: '快速',
  standard: '标准',
  deep: '深度',
};

export function scanPolicyText(policy: string): string {
  return SCAN_POLICY_TEXT[policy] ?? policy;
}

/** Finding severity display */
export const FINDING_SEVERITY_TEXT: Record<string, string> = {
  critical: '严重',
  high: '高危',
  medium: '中危',
  low: '低危',
  info: '信息',
};

export const FINDING_SEVERITY_BADGE: Record<
  string,
  'success' | 'default' | 'error' | 'processing' | 'warning'
> = {
  critical: 'error',
  high: 'error',
  medium: 'warning',
  low: 'default',
  info: 'processing',
};

export function findingSeverityText(severity: string): string {
  return FINDING_SEVERITY_TEXT[severity] ?? severity;
}
