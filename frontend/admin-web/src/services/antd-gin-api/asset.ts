import { request } from '@umijs/max';

export interface Asset {
  asset_code: string;
  name: string;
  root_url: string;
  dept_code: string;
  remark?: string;
  status: 'active' | 'disabled';
  updated_at: string;
}

export interface PageResult<T> {
  items: T[];
  page: number;
  size: number;
  total: number;
}

export interface CreateAssetRequest {
  name: string;
  root_url: string;
  dept_code: string;
  remark?: string;
}

export interface UpdateAssetRequest {
  name?: string;
  root_url?: string;
  dept_code?: string;
  remark?: string;
  status?: 'active' | 'disabled';
}

export interface ListAssetParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  status?: string;
}

/** 获取资产列表 */
export async function getAssetList(params?: ListAssetParams) {
  return request<{
    code: number;
    data: PageResult<Asset>;
    message?: string;
  }>('/api/v1/scan/assets', {
    method: 'GET',
    params,
  });
}

/** 获取资产详情 */
export async function getAsset(assetCode: string) {
  return request<{
    code: number;
    data: Asset;
    message?: string;
  }>(`/api/v1/scan/assets/${assetCode}`, {
    method: 'GET',
  });
}

/** 创建资产 */
export async function createAsset(data: CreateAssetRequest) {
  return request<{
    code: number;
    data: Asset;
    message?: string;
  }>('/api/v1/scan/assets', {
    method: 'POST',
    data,
  });
}

/** 更新资产 */
export async function updateAsset(assetCode: string, data: UpdateAssetRequest) {
  return request<{
    code: number;
    data: Asset;
    message?: string;
  }>(`/api/v1/scan/assets/${assetCode}`, {
    method: 'PUT',
    data,
  });
}

/** 删除资产 */
export async function deleteAsset(assetCode: string) {
  return request<{
    code: number;
    data: { deleted: boolean };
    message?: string;
  }>(`/api/v1/scan/assets/${assetCode}`, {
    method: 'DELETE',
  });
}
