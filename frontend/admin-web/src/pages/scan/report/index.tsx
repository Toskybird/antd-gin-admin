import { searchFormLayout } from '@/constants/formLayout';
import {
  downloadScanReport,
  getScanReportList,
  type ScanReport,
} from '@/services/antd-gin-api/scan_report';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { useSearchParams } from '@umijs/max';
import { Button, Empty, message } from 'antd';
import React, { useRef, useState } from 'react';

const triggerBlobDownload = (blob: Blob, filename: string) => {
  const url = window.URL.createObjectURL(blob);
  const anchor = document.createElement('a');
  anchor.href = url;
  anchor.download = filename;
  anchor.click();
  window.URL.revokeObjectURL(url);
};

const ScanReportList: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const [searchParams] = useSearchParams();
  const initialJobCode = searchParams.get('job_code') || undefined;
  const [downloadingCode, setDownloadingCode] = useState<string | null>(null);

  const handleDownload = async (record: ScanReport) => {
    setDownloadingCode(record.report_code);
    try {
      const blob = await downloadScanReport(record.report_code);
      const filename = `scan-report-${record.job_code}-v${record.version}.html`;
      triggerBlobDownload(blob, filename);
      message.success('下载已开始');
    } catch {
      message.error('下载失败');
    } finally {
      setDownloadingCode(null);
    }
  };

  const columns: ProColumns<ScanReport>[] = [
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
      width: 80,
      hideInSearch: true,
    },
    {
      title: '所属任务',
      dataIndex: 'job_code',
      key: 'job_code',
      width: 220,
      ellipsis: true,
      copyable: true,
      initialValue: initialJobCode,
    },
    {
      title: '生成时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 170,
      hideInSearch: true,
    },
    {
      title: '格式',
      dataIndex: 'format',
      key: 'format',
      width: 80,
      hideInSearch: true,
      render: () => 'HTML',
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      hideInSearch: true,
      render: (_, record) => (
        <Button
          type="link"
          size="small"
          loading={downloadingCode === record.report_code}
          onClick={() => handleDownload(record)}
        >
          下载
        </Button>
      ),
    },
  ];

  return (
    <PageContainer title="扫描报告">
      <ProTable<ScanReport>
        actionRef={actionRef}
        rowKey="report_code"
        search={{ ...searchFormLayout }}
        params={{ job_code: initialJobCode }}
        form={{ initialValues: { job_code: initialJobCode } }}
        locale={{ emptyText: <Empty description="暂无报告" /> }}
        request={async (params) => {
          const res = await getScanReportList({
            page: params.current || 1,
            pageSize: params.pageSize || 10,
            job_code: params.job_code,
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

export default ScanReportList;
