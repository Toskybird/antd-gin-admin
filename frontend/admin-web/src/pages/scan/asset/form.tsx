import {
  modalFormLayout,
  modalFormSubmitterLayout,
} from '@/constants/formLayout';
import {
  createAsset,
  getAsset,
  updateAsset,
  type CreateAssetRequest,
  type UpdateAssetRequest,
} from '@/services/antd-gin-api/asset';
import { getDeptTree, type Dept } from '@/services/antd-gin-api/dept';
import { hasPermission } from '@/utils/permission';
import {
  PageContainer,
  ProCard,
  ProForm,
  ProFormText,
  ProFormTextArea,
  ProFormTreeSelect,
} from '@ant-design/pro-components';
import { history, useModel, useParams } from '@umijs/max';
import { Button, Form, message, Spin } from 'antd';
import React, { useEffect, useState } from 'react';

type TreeOption = {
  title: string;
  value: string;
  children?: TreeOption[];
};

const buildDeptTreeOptions = (depts: Dept[]): TreeOption[] =>
  depts.map((dept) => ({
    title: dept.dept_name,
    value: dept.dept_code,
    children: dept.children?.length
      ? buildDeptTreeOptions(dept.children)
      : undefined,
  }));

const AssetForm: React.FC = () => {
  const { code } = useParams<{ code: string }>();
  const isEdit = !!code;
  const [loading, setLoading] = useState(isEdit);
  const [initialValues, setInitialValues] = useState<
    CreateAssetRequest | undefined
  >();
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;
  const canCreate = hasPermission(currentUser, 'scan:asset:create');
  const canUpdate = hasPermission(currentUser, 'scan:asset:update');

  useEffect(() => {
    if (!isEdit || !code) {
      return;
    }
    setLoading(true);
    getAsset(code)
      .then((res) => {
        if (res.code === 200 && res.data) {
          setInitialValues({
            name: res.data.name,
            root_url: res.data.root_url,
            dept_code: res.data.dept_code,
            remark: res.data.remark,
          });
        }
      })
      .finally(() => setLoading(false));
  }, [code, isEdit]);

  const handleFinish = async (values: CreateAssetRequest) => {
    if (isEdit && code) {
      if (!canUpdate) {
        message.error('无编辑权限');
        return false;
      }
      const payload: UpdateAssetRequest = {
        name: values.name,
        root_url: values.root_url,
        dept_code: values.dept_code,
        remark: values.remark,
      };
      const res = await updateAsset(code, payload);
      if (res.code === 200) {
        message.success('保存成功');
        history.push('/scan/asset');
        return true;
      }
      return false;
    }
    if (!canCreate) {
      message.error('无创建权限');
      return false;
    }
    const res = await createAsset(values);
    if (res.code === 200) {
      message.success('创建成功');
      history.push('/scan/asset');
      return true;
    }
    return false;
  };

  return (
    <PageContainer
      title={isEdit ? '编辑资产' : '新建资产'}
      onBack={() => history.push('/scan/asset')}
    >
      <ProCard>
        <Spin spinning={loading}>
          <ProForm<CreateAssetRequest>
            {...modalFormLayout}
            style={{ maxWidth: 640 }}
            initialValues={initialValues}
            key={initialValues ? 'loaded' : isEdit ? 'loading' : 'create'}
            onFinish={handleFinish}
            submitter={{
              render: (_, dom) => (
                <Form.Item {...modalFormSubmitterLayout}>
                  {dom}
                  <Button
                    style={{ marginLeft: 8 }}
                    onClick={() => history.push('/scan/asset')}
                  >
                    取消
                  </Button>
                </Form.Item>
              ),
            }}
          >
            <ProFormText
              name="name"
              label="名称"
              rules={[{ required: true, message: '请输入资产名称' }]}
            />
            <ProFormText
              name="root_url"
              label="根 URL"
              rules={[
                { required: true, message: '请输入根 URL' },
                { type: 'url', message: '请输入有效的 URL' },
              ]}
            />
            <ProFormTreeSelect
              name="dept_code"
              label="所属部门"
              rules={[{ required: true, message: '请选择所属部门' }]}
              request={async () => {
                const res = await getDeptTree();
                if (res.code === 200 && res.data) {
                  return buildDeptTreeOptions(res.data);
                }
                return [];
              }}
              fieldProps={{
                showSearch: true,
                treeDefaultExpandAll: true,
                placeholder: '请选择部门',
              }}
            />
            <ProFormTextArea name="remark" label="备注" fieldProps={{ rows: 3 }} />
          </ProForm>
        </Spin>
      </ProCard>
    </PageContainer>
  );
};

export default AssetForm;
