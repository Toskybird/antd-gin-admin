import { searchFormLayout } from '@/constants/formLayout';
import {
  getDetectionRuleList,
  syncDetectionRules,
  updateDetectionRuleStatus,
  type DetectionRule,
} from '@/services/antd-gin-api/detection_rule';
import { hasPermission } from '@/utils/permission';
import { SyncOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { useModel } from '@umijs/max';
import { Button, Empty, message, Switch, Typography } from 'antd';
import React, { useRef, useState } from 'react';

const DetectionRuleList: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const [syncing, setSyncing] = useState(false);
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const canSync = hasPermission(currentUser, 'scan:detection-rule:sync');
  const canUpdate = hasPermission(currentUser, 'scan:detection-rule:update');

  const handleSync = async () => {
    setSyncing(true);
    try {
      const res = await syncDetectionRules();
      if (res.code === 200) {
        message.success('同步完成');
        actionRef.current?.reload();
      }
    } finally {
      setSyncing(false);
    }
  };

  const handleToggle = async (record: DetectionRule, checked: boolean) => {
    const res = await updateDetectionRuleStatus(record.rule_code, checked);
    if (res.code === 200) {
      message.success(checked ? '已启用' : '已停用');
      actionRef.current?.reload();
    } else {
      message.error('更新失败');
      actionRef.current?.reload();
    }
  };

  const columns: ProColumns<DetectionRule>[] = [
    {
      title: '规则代码',
      dataIndex: 'rule_code',
      key: 'rule_code',
      width: 220,
      ellipsis: true,
      copyable: true,
    },
    {
      title: '展示名',
      dataIndex: 'display_name',
      key: 'display_name',
      ellipsis: true,
    },
    {
      title: '启用',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 100,
      hideInSearch: true,
      render: (_, record) => (
        <Switch
          checked={record.enabled}
          disabled={!canUpdate}
          onChange={(checked) => handleToggle(record, checked)}
        />
      ),
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      key: 'updated_at',
      width: 170,
      hideInSearch: true,
    },
  ];

  return (
    <PageContainer
      title="检测规则"
      content={
        <Typography.Text type="secondary">
          规则由系统自模板同步，不可自定义正文。
        </Typography.Text>
      }
    >
      <ProTable<DetectionRule>
        actionRef={actionRef}
        rowKey="rule_code"
        search={{ ...searchFormLayout, filterType: 'light' }}
        toolBarRender={() =>
          canSync
            ? [
                <Button
                  key="sync"
                  icon={<SyncOutlined />}
                  loading={syncing}
                  onClick={handleSync}
                >
                  同步规则
                </Button>,
              ]
            : []
        }
        locale={{
          emptyText: (
            <Empty description="暂无检测规则，请先执行同步" />
          ),
        }}
        request={async (params) => {
          const res = await getDetectionRuleList({
            page: params.current || 1,
            pageSize: params.pageSize || 10,
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
      />
    </PageContainer>
  );
};

export default DetectionRuleList;
