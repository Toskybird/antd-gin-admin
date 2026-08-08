import {
  FINDING_SEVERITY_BADGE,
  findingSeverityText,
} from '@/constants/scanStatus';
import { getFinding, type Finding } from '@/services/antd-gin-api/finding';
import { PageContainer, ProCard } from '@ant-design/pro-components';
import { history, Link, useParams } from '@umijs/max';
import { Badge, Button, Descriptions, Empty, Spin } from 'antd';
import React, { useCallback, useEffect, useState } from 'react';

const FindingDetail: React.FC = () => {
  const { code } = useParams<{ code: string }>();
  const [finding, setFinding] = useState<Finding | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadFinding = useCallback(async () => {
    if (!code) {
      return;
    }
    setError(false);
    setLoading(true);
    try {
      const res = await getFinding(code);
      if (res.code === 200 && res.data) {
        setFinding(res.data);
      } else {
        setError(true);
      }
    } catch {
      setError(true);
    } finally {
      setLoading(false);
    }
  }, [code]);

  useEffect(() => {
    loadFinding();
  }, [loadFinding]);

  if (error && !finding) {
    return (
      <PageContainer
        title="发现项详情"
        onBack={() => history.push('/scan/finding')}
      >
        <ProCard>
          <Empty description="加载失败">
            <Button type="primary" onClick={loadFinding}>
              重试
            </Button>
          </Empty>
        </ProCard>
      </PageContainer>
    );
  }

  return (
    <PageContainer
      title="发现项详情"
      onBack={() => history.push('/scan/finding')}
    >
      <Spin spinning={loading}>
        {finding && (
          <ProCard>
            <Descriptions column={1} bordered>
              <Descriptions.Item label="严重级别">
                <Badge
                  status={FINDING_SEVERITY_BADGE[finding.severity] ?? 'default'}
                  text={findingSeverityText(finding.severity)}
                />
              </Descriptions.Item>
              <Descriptions.Item label="规则代码">
                {finding.rule_code}
              </Descriptions.Item>
              <Descriptions.Item label="规则名称">
                {finding.rule_display_name}
              </Descriptions.Item>
              <Descriptions.Item label="标题">{finding.title}</Descriptions.Item>
              <Descriptions.Item label="描述">
                {finding.description || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="证据">
                <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>
                  {finding.evidence || '-'}
                </pre>
              </Descriptions.Item>
              <Descriptions.Item label="位置">
                {finding.location || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="关联任务">
                <Link to={`/scan/job/detail/${finding.job_code}`}>
                  {finding.job_code}
                </Link>
              </Descriptions.Item>
              <Descriptions.Item label="发现时间">
                {finding.created_at}
              </Descriptions.Item>
            </Descriptions>
          </ProCard>
        )}
      </Spin>
    </PageContainer>
  );
};

export default FindingDetail;
