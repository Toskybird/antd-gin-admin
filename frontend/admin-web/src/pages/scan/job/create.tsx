import {
  modalFormLayout,
  modalFormSubmitterLayout,
} from '@/constants/formLayout';
import { getAssetList, type Asset } from '@/services/antd-gin-api/asset';
import {
  createScanJob,
  type CreateScanJobRequest,
} from '@/services/antd-gin-api/scan_job';
import { hasPermission } from '@/utils/permission';
import {
  PageContainer,
  ProCard,
  ProForm,
  ProFormDigit,
  ProFormRadio,
  ProFormSelect,
  ProFormText,
} from '@ant-design/pro-components';
import { history, useModel } from '@umijs/max';
import { Button, Form, message } from 'antd';
import React, { useState } from 'react';

const ScanJobCreate: React.FC = () => {
  const [source, setSource] = useState<'asset' | 'temp_url'>('asset');
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const canCreate = hasPermission(currentUser, 'scan:job:create');

  const handleFinish = async (values: CreateScanJobRequest & { source: 'asset' | 'temp_url' }) => {
    if (!canCreate) {
      message.error('无创建权限');
      return false;
    }
    const payload: CreateScanJobRequest = {
      source: values.source,
      policy: values.policy,
      max_depth: values.max_depth,
      max_pages: values.max_pages,
    };
    if (values.source === 'asset') {
      payload.asset_code = values.asset_code;
    } else {
      payload.entry_url = values.entry_url;
    }
    const res = await createScanJob(payload);
    if (res.code === 200 && res.data) {
      message.success('扫描任务已创建');
      history.push(`/scan/job/detail/${res.data.job_code}`);
      return true;
    }
    return false;
  };

  return (
    <PageContainer title="创建扫描" onBack={() => history.push('/scan/job')}>
      <ProCard>
        <ProForm<CreateScanJobRequest & { source: 'asset' | 'temp_url' }>
          {...modalFormLayout}
          style={{ maxWidth: 640 }}
          initialValues={{
            source: 'asset',
            policy: 'standard',
            max_depth: 3,
            max_pages: 100,
          }}
          onFinish={handleFinish}
          submitter={{
            searchConfig: { submitText: '开始扫描' },
            render: (_, dom) => (
              <Form.Item {...modalFormSubmitterLayout}>
                {dom}
                <Button
                  style={{ marginLeft: 8 }}
                  onClick={() => history.push('/scan/job')}
                >
                  取消
                </Button>
              </Form.Item>
            ),
          }}
        >
          <ProFormRadio.Group
            name="source"
            label="目标来源"
            options={[
              { label: '选择资产', value: 'asset' },
              { label: '临时 URL', value: 'temp_url' },
            ]}
            fieldProps={{
              onChange: (e) => setSource(e.target.value),
            }}
          />
          {source === 'asset' ? (
            <ProFormSelect
              name="asset_code"
              label="资产"
              rules={[{ required: true, message: '请选择资产' }]}
              request={async () => {
                const res = await getAssetList({
                  page: 1,
                  pageSize: 500,
                  status: 'active',
                });
                if (res.code === 200 && res.data) {
                  return res.data.items.map((asset: Asset) => ({
                    label: `${asset.name} (${asset.root_url})`,
                    value: asset.asset_code,
                  }));
                }
                return [];
              }}
              fieldProps={{ showSearch: true, placeholder: '请选择启用中的资产' }}
            />
          ) : (
            <ProFormText
              name="entry_url"
              label="入口 URL"
              rules={[
                { required: true, message: '请输入入口 URL' },
                { type: 'url', message: '请输入有效的 URL' },
              ]}
            />
          )}
          <ProFormDigit
            name="max_depth"
            label="最大爬取深度"
            min={1}
            max={10}
            rules={[{ required: true, message: '请输入最大爬取深度' }]}
          />
          <ProFormDigit
            name="max_pages"
            label="最大页面数"
            min={1}
            max={10000}
            rules={[{ required: true, message: '请输入最大页面数' }]}
          />
          <ProFormRadio.Group
            name="policy"
            label="扫描策略"
            options={[
              { label: '快速', value: 'quick' },
              { label: '标准', value: 'standard' },
              { label: '深度', value: 'deep' },
            ]}
          />
        </ProForm>
      </ProCard>
    </PageContainer>
  );
};

export default ScanJobCreate;
