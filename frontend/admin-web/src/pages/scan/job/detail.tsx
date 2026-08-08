import {
  FINDING_SEVERITY_BADGE,
  findingSeverityText,
  isScanJobActive,
  SCAN_JOB_STATUS_BADGE,
  SCAN_POLICY_TEXT,
  scanJobStatusText,
} from '@/constants/scanStatus';
import { getFindingList, type Finding } from '@/services/antd-gin-api/finding';
import { createScanReport } from '@/services/antd-gin-api/scan_report';
import {
  cancelScanJob,
  getScanJob,
  type ScanJob,
} from '@/services/antd-gin-api/scan_job';
import { hasPermission } from '@/utils/permission';
import type { ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProCard, ProTable } from '@ant-design/pro-components';
import { history, useModel, useParams } from '@umijs/max';
import {
  Badge,
  Button,
  Descriptions,
  Empty,
  message,
  Modal,
  Space,
  Spin,
} from 'antd';
import React, { useCallback, useEffect, useState } from 'react';

const ScanJobDetail: React.FC = () => {
  const { code } = useParams<{ code: string }>();
  const [job, setJob] = useState<ScanJob | null>(null);
  const [findings, setFindings] = useState<Finding[]>([]);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [error, setError] = useState(false);
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const canCancel = hasPermission(currentUser, 'scan:job:cancel');
  const canCreateReport = hasPermission(currentUser, 'scan:report:create');

  const loadData = useCallback(async () => {
    if (!code) {
      return;
    }
    setError(false);
    try {
      const [jobRes, findingRes] = await Promise.all([
        getScanJob(code),
        getFindingList({ job_code: code, page: 1, pageSize: 20 }),
      ]);
      if (jobRes.code === 200 && jobRes.data) {
        setJob(jobRes.data);
      } else {
        setError(true);
      }
      if (findingRes.code === 200 && findingRes.data) {
        setFindings(findingRes.data.items);
      }
    } catch {
      setError(true);
    } finally {
      setLoading(false);
    }
  }, [code]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  useEffect(() => {
    if (!job || !isScanJobActive(job.status)) {
      return;
    }
    const timer = window.setInterval(() => {
      loadData();
    }, 3000);
    return () => window.clearInterval(timer);
  }, [job?.status, loadData]);

  const handleCancel = () => {
    if (!job) {
      return;
    }
    Modal.confirm({
      title: '确认取消',
      content: `确定要取消扫描任务 "${job.job_code}" 吗？`,
      onOk: async () => {
        const res = await cancelScanJob(job.job_code);
        if (res.code === 200) {
          message.success('已取消');
          loadData();
        }
      },
    });
  };

  const handleGenerateReport = async () => {
    if (!job) {
      return;
    }
    setGenerating(true);
    try {
      const res = await createScanReport({ job_code: job.job_code });
      if (res.code === 200) {
        message.success('报告已生成');
      }
    } finally {
      setGenerating(false);
    }
  };

  const canGenerateReport =
    job &&
    (job.status === 'succeeded' || job.status === 'cancelled') &&
    canCreateReport;

  const findingColumns: ProColumns<Finding>[] = [
    {
      title: '严重级别',
      dataIndex: 'severity',
      key: 'severity',
      width: 100,
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
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
  ];

  if (error && !job) {
    return (
      <PageContainer title="任务详情" onBack={() => history.push('/scan/job')}>
        <ProCard>
          <Empty description="加载失败">
            <Button type="primary" onClick={() => { setLoading(true); loadData(); }}>
              重试
            </Button>
          </Empty>
        </ProCard>
      </PageContainer>
    );
  }

  return (
    <PageContainer
      title="任务详情"
      onBack={() => history.push('/scan/job')}
      extra={
        job && (
          <Space>
            {canCancel && isScanJobActive(job.status) && (
              <Button danger onClick={handleCancel}>
                取消
              </Button>
            )}
            {canGenerateReport && (
              <>
                <Button loading={generating} onClick={handleGenerateReport}>
                  生成报告
                </Button>
                <Button
                  onClick={() =>
                    history.push(`/scan/report?job_code=${job.job_code}`)
                  }
                >
                  查看报告
                </Button>
              </>
            )}
          </Space>
        )
      }
    >
      <Spin spinning={loading}>
        {job && (
          <>
            <ProCard title="任务信息" style={{ marginBottom: 16 }}>
              <Descriptions column={2}>
                <Descriptions.Item label="任务标识">
                  {job.job_code}
                </Descriptions.Item>
                <Descriptions.Item label="状态">
                  <Badge
                    status={SCAN_JOB_STATUS_BADGE[job.status] ?? 'default'}
                    text={scanJobStatusText(job.status)}
                  />
                </Descriptions.Item>
                <Descriptions.Item label="入口 URL">
                  {job.entry_url}
                </Descriptions.Item>
                <Descriptions.Item label="策略">
                  {SCAN_POLICY_TEXT[job.policy] ?? job.policy}
                </Descriptions.Item>
                <Descriptions.Item label="最大深度">
                  {job.max_depth}
                </Descriptions.Item>
                <Descriptions.Item label="最大页面数">
                  {job.max_pages}
                </Descriptions.Item>
                <Descriptions.Item label="发现项数">
                  {job.finding_count}
                </Descriptions.Item>
                <Descriptions.Item label="创建时间">
                  {job.created_at}
                </Descriptions.Item>
                {job.finished_at && (
                  <Descriptions.Item label="结束时间">
                    {job.finished_at}
                  </Descriptions.Item>
                )}
              </Descriptions>
            </ProCard>

            <ProCard title="发现项摘要">
              {findings.length === 0 ? (
                <Empty description="未发现检测结果" />
              ) : (
                <ProTable<Finding>
                  rowKey="finding_code"
                  search={false}
                  options={false}
                  pagination={false}
                  dataSource={findings}
                  columns={findingColumns}
                  onRow={(record) => ({
                    onClick: () =>
                      history.push(
                        `/scan/finding/detail/${record.finding_code}`,
                      ),
                    style: { cursor: 'pointer' },
                  })}
                />
              )}
            </ProCard>
          </>
        )}
      </Spin>
    </PageContainer>
  );
};

export default ScanJobDetail;
