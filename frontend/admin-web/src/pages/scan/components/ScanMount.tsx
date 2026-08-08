import { PageContainer } from '@ant-design/pro-components';
import { Card, Typography } from 'antd';
import React from 'react';

type ScanMountProps = {
  title: string;
};

/** Temporary mount point until each scan page Issue lands. */
const ScanMount: React.FC<ScanMountProps> = ({ title }) => (
  <PageContainer title={title}>
    <Card>
      <Typography.Title level={4} style={{ marginTop: 0 }}>
        功能页挂载点
      </Typography.Title>
      <Typography.Paragraph type="secondary" style={{ marginBottom: 0 }}>
        扫描业务页面在此区域渲染。当前为「{title}」占位，完整列表/表单由对应 UI Issue 交付。
      </Typography.Paragraph>
    </Card>
  </PageContainer>
);

export default ScanMount;
