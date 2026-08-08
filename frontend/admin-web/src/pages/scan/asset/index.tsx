import { searchFormLayout } from '@/constants/formLayout';
import {
  ASSET_STATUS_BADGE,
  assetStatusText,
} from '@/constants/scanStatus';
import {
  deleteAsset,
  getAssetList,
  updateAsset,
  type Asset,
} from '@/services/antd-gin-api/asset';
import { hasPermission } from '@/utils/permission';
import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { history, useModel } from '@umijs/max';
import { Badge, Button, Empty, message, Modal } from 'antd';
import React, { useRef } from 'react';

const AssetList: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const canCreate = hasPermission(currentUser, 'scan:asset:create');
  const canUpdate = hasPermission(currentUser, 'scan:asset:update');
  const canDelete = hasPermission(currentUser, 'scan:asset:delete');

  const handleToggleStatus = (record: Asset) => {
    const nextStatus = record.status === 'active' ? 'disabled' : 'active';
    const actionText = nextStatus === 'active' ? '启用' : '停用';
    Modal.confirm({
      title: `确认${actionText}`,
      content: `确定要${actionText}资产 "${record.name}" 吗？`,
      onOk: async () => {
        const res = await updateAsset(record.asset_code, { status: nextStatus });
        if (res.code === 200) {
          message.success(`${actionText}成功`);
          actionRef.current?.reload();
        }
      },
    });
  };

  const handleDelete = (record: Asset) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除资产 "${record.name}" 吗？`,
      onOk: async () => {
        const res = await deleteAsset(record.asset_code);
        if (res.code === 200) {
          message.success('删除成功');
          actionRef.current?.reload();
        }
      },
    });
  };

  const columns: ProColumns<Asset>[] = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '根 URL',
      dataIndex: 'root_url',
      key: 'root_url',
      ellipsis: true,
      copyable: true,
      hideInSearch: true,
    },
    {
      title: '所属部门',
      dataIndex: 'dept_code',
      key: 'dept_code',
      width: 120,
      hideInSearch: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: {
        active: { text: '启用' },
        disabled: { text: '停用' },
      },
      render: (_, record) => (
        <Badge
          status={ASSET_STATUS_BADGE[record.status] ?? 'default'}
          text={assetStatusText(record.status)}
        />
      ),
    },
    {
      title: '关键字',
      dataIndex: 'keyword',
      key: 'keyword',
      hideInTable: true,
      fieldProps: { placeholder: '名称或 URL' },
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      key: 'updated_at',
      width: 170,
      hideInSearch: true,
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      fixed: 'right',
      hideInSearch: true,
      render: (_, record) =>
        [
          canUpdate && (
            <Button
              key="edit"
              type="link"
              size="small"
              onClick={() =>
                history.push(`/scan/asset/edit/${record.asset_code}`)
              }
            >
              编辑
            </Button>
          ),
          canUpdate && (
            <Button
              key="toggle"
              type="link"
              size="small"
              onClick={() => handleToggleStatus(record)}
            >
              {record.status === 'active' ? '停用' : '启用'}
            </Button>
          ),
          canDelete && (
            <Button
              key="delete"
              type="link"
              size="small"
              danger
              onClick={() => handleDelete(record)}
            >
              删除
            </Button>
          ),
        ].filter(Boolean),
    },
  ];

  return (
    <PageContainer title="资产">
      <ProTable<Asset>
        actionRef={actionRef}
        rowKey="asset_code"
        search={{ ...searchFormLayout }}
        toolBarRender={() =>
          canCreate
            ? [
                <Button
                  type="primary"
                  key="create"
                  onClick={() => history.push('/scan/asset/create')}
                >
                  <PlusOutlined /> 新建资产
                </Button>,
              ]
            : []
        }
        locale={{
          emptyText: (
            <Empty description="暂无资产">
              {canCreate && (
                <Button
                  type="primary"
                  onClick={() => history.push('/scan/asset/create')}
                >
                  新建资产
                </Button>
              )}
            </Empty>
          ),
        }}
        request={async (params) => {
          const res = await getAssetList({
            page: params.current || 1,
            pageSize: params.pageSize || 10,
            keyword: params.keyword,
            status: params.status,
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

export default AssetList;
