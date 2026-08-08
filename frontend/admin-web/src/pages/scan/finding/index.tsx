import { searchFormLayout } from '@/constants/formLayout';
import {
  FINDING_SEVERITY_BADGE,
  findingSeverityText,
} from '@/constants/scanStatus';
import { getFindingList, type Finding } from '@/services/antd-gin-api/finding';
import type { ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Badge, Empty } from 'antd';
import React from 'react';

const FindingList: React.FC = () => {
  const columns: ProColumns<Finding>[] = [
    {
      title: '严重级别',
      dataIndex: 'severity',
      key: 'severity',
      width: 100,
      valueType: 'select',
      valueEnum: {
        critical: { text: '严重' },
        high: { text: '高危' },
        medium: { text: '中危' },
        low: { text: '低危' },
        info: { text: '信息' },
      },
      render: (_, record) => (
        <Badge
          status={FINDING_SEVERITY_BADGE[record.severity] ?? 'default'}
          text={findingSeverityText(record.severity)}
        />
      ),
    },
    {
      title: '规则',
      dataIndex: 'rule_display_name',
      key: 'rule_display_name',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '所属任务',
      dataIndex: 'job_code',
      key: 'job_code',
      width: 220,
      ellipsis: true,
      copyable: true,
    },
    {
      title: '检测规则',
      dataIndex: 'rule_code',
      key: 'rule_code',
      hideInTable: true,
      fieldProps: { placeholder: '规则代码' },
    },
    {
      title: '位置',
      dataIndex: 'location',
      key: 'location',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '发现时间',
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
      fieldProps: { placeholder: '标题或描述' },
    },
  ];

  return (
    <PageContainer title="发现项">
      <ProTable<Finding>
        rowKey="finding_code"
        search={{ ...searchFormLayout }}
        locale={{ emptyText: <Empty description="暂无发现项" /> }}
        request={async (params) => {
          const res = await getFindingList({
            page: params.current || 1,
            pageSize: params.pageSize || 10,
            job_code: params.job_code,
            rule_code: params.rule_code,
            severity: params.severity,
            keyword: params.keyword,
          });
          if (res.code === 200 && res.data) {
            return {
              data: res.data.items,
              success: true,
              total: res.data.total,
            };
          }
          return { data: [], success: false, total: 0 };
        }}
        columns={columns}
        onRow={(record) => ({
          onClick: () =>
            history.push(`/scan/finding/detail/${record.finding_code}`),
          style: { cursor: 'pointer' },
        })}
      />
    </PageContainer>
  );
};

export default FindingList;
