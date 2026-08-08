import { searchFormLayout } from '@/constants/formLayout';
import {
  isScanJobActive,
  SCAN_JOB_STATUS_BADGE,
  SCAN_POLICY_TEXT,
  scanJobStatusText,
} from '@/constants/scanStatus';
import {
  cancelScanJob,
  getScanJobList,
  type ScanJob,
} from '@/services/antd-gin-api/scan_job';
import { hasPermission } from '@/utils/permission';
import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { history, useModel } from '@umijs/max';
import { Badge, Button, Empty, message, Modal } from 'antd';
import React, { useRef, useState } from 'react';

const ScanJobList: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const [polling, setPolling] = useState<number | undefined>(undefined);
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const canCreate = hasPermission(currentUser, 'scan:job:create');
  const canCancel = hasPermission(currentUser, 'scan:job:cancel');

  const handleCancel = (record: ScanJob) => {
    Modal.confirm({
      title: '确认取消',
      content: `确定要取消扫描任务 "${record.job_code}" 吗？`,
      onOk: async () => {
        const res = await cancelScanJob(record.job_code);
        if (res.code === 200) {
          message.success('已取消');
          actionRef.current?.reload();
        }
      },
    });
  };

  const columns: ProColumns<ScanJob>[] = [
    {
      title: '任务标识',
      dataIndex: 'job_code',
      key: 'job_code',
      width: 220,
      ellipsis: true,
      copyable: true,
      hideInSearch: true,
    },
    {
      title: '入口 URL',
      dataIndex: 'entry_url',
      key: 'entry_url',
      ellipsis: true,
      copyable: true,
      hideInSearch: true,
    },
    {
      title: '策略',
      dataIndex: 'policy',
      key: 'policy',
      width: 90,
      valueType: 'select',
      valueEnum: {
        quick: { text: '快速' },
        standard: { text: '标准' },
        deep: { text: '深度' },
      },
      render: (_, record) => SCAN_POLICY_TEXT[record.policy] ?? record.policy,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: {
        queued: { text: '排队中' },
        running: { text: '执行中' },
        succeeded: { text: '已完成' },
        failed: { text: '失败' },
        cancelled: { text: '已取消' },
      },
      render: (_, record) => (
        <Badge
          status={SCAN_JOB_STATUS_BADGE[record.status] ?? 'default'}
          text={scanJobStatusText(record.status)}
        />
      ),
    },
    {
      title: '发现项数',
      dataIndex: 'finding_count',
      key: 'finding_count',
      width: 90,
      hideInSearch: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 170,
      hideInSearch: true,
    },
    {
      title: '关键字',
      dataIndex: 'keyword',
      key: 'keyword',
      hideInTable: true,
      fieldProps: { placeholder: '入口 URL 或任务标识' },
    },
    {
      title: '操作',
      key: 'action',
      width: 140,
      fixed: 'right',
      hideInSearch: true,
      render: (_, record) =>
        [
          <Button
            key="detail"
            type="link"
            size="small"
            onClick={() =>
              history.push(`/scan/job/detail/${record.job_code}`)
            }
          >
            详情
          </Button>,
          canCancel &&
            isScanJobActive(record.status) && (
              <Button
                key="cancel"
                type="link"
                size="small"
                danger
                onClick={() => handleCancel(record)}
              >
                取消
              </Button>
            ),
        ].filter(Boolean),
    },
  ];

  return (
    <PageContainer title="扫描任务">
      <ProTable<ScanJob>
        actionRef={actionRef}
        rowKey="job_code"
        search={{ ...searchFormLayout }}
        polling={polling}
        toolBarRender={() =>
          canCreate
            ? [
                <Button
                  type="primary"
                  key="create"
                  onClick={() => history.push('/scan/job/create')}
                >
                  <PlusOutlined /> 创建扫描
                </Button>,
              ]
            : []
        }
        locale={{
          emptyText: (
            <Empty description="暂无扫描任务">
              {canCreate && (
                <Button
                  type="primary"
                  onClick={() => history.push('/scan/job/create')}
                >
                  创建扫描
                </Button>
              )}
            </Empty>
          ),
        }}
        request={async (params) => {
          const res = await getScanJobList({
            page: params.current || 1,
            pageSize: params.pageSize || 10,
            status: params.status,
            policy: params.policy,
            keyword: params.keyword,
          });
          if (res.code === 200 && res.data) {
            const hasActive = res.data.items.some((job) =>
              isScanJobActive(job.status),
            );
            setPolling(hasActive ? 3000 : undefined);
            return {
              data: res.data.items,
              success: true,
              total: res.data.total,
            };
          }
          setPolling(undefined);
          return { data: [], success: false, total: 0 };
        }}
        columns={columns}
      />
    </PageContainer>
  );
};

export default ScanJobList;
