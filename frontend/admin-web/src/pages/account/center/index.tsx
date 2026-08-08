import { updateCurrentProfile } from "@/services/antd-gin-api/user";
import { modalFormLayout } from "@/constants/formLayout";
import {
  PageContainer,
  ProCard,
  ProForm,
  type ProFormInstance,
  ProFormText,
} from "@ant-design/pro-components";
import { useModel } from "@umijs/max";
import { Descriptions, message, Tag } from "antd";
import React, { useEffect, useRef } from "react";

const AccountCenter: React.FC = () => {
  const { initialState, setInitialState } = useModel("@@initialState");
  const currentUser = initialState?.currentUser;
  const formRef = useRef<ProFormInstance>(null);

  useEffect(() => {
    formRef.current?.setFieldsValue({
      nickname: currentUser?.nickname,
      email: currentUser?.email,
      phone: currentUser?.phone,
    });
  }, [currentUser]);

  return (
    <PageContainer>
      <ProCard>
        <ProForm
          {...modalFormLayout}
          formRef={formRef}
          style={{ maxWidth: 640 }}
          initialValues={{
            nickname: currentUser?.nickname,
            email: currentUser?.email,
            phone: currentUser?.phone,
          }}
          onFinish={async (values) => {
            const res = await updateCurrentProfile(values);
            if (res.code === 200 && res.data) {
              setInitialState((state) => ({
                ...state,
                currentUser: {
                  ...state?.currentUser,
                  nickname: res.data.nickname,
                  email: res.data.email,
                  phone: res.data.phone,
                },
              }));
              message.success("保存成功");
              return true;
            }
            return false;
          }}
        >
          <ProFormText name="nickname" label="昵称" />
          <ProFormText
            name="email"
            label="邮箱"
            rules={[{ type: "email", message: "邮箱格式不正确" }]}
          />
          <ProFormText name="phone" label="手机号" />
        </ProForm>
        <Descriptions
          bordered
          column={{ xs: 1, sm: 1, md: 2 }}
          size="middle"
          style={{ marginTop: 24 }}
        >
          <Descriptions.Item label="用户名">
            {currentUser?.name || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="用户编码">
            {currentUser?.userid || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="部门">
            {currentUser?.deptName || currentUser?.deptCode || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            {currentUser?.status === undefined ? (
              "-"
            ) : (
              <Tag color={currentUser.status === 1 ? "green" : "red"}>
                {currentUser.status === 1 ? "启用" : "停用"}
              </Tag>
            )}
          </Descriptions.Item>
          <Descriptions.Item label="账号类型">
            <Tag color={currentUser?.isSuperAdmin ? "blue" : "default"}>
              {currentUser?.isSuperAdmin ? "超级管理员" : "普通用户"}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="角色" span={2}>
            {currentUser?.roles && currentUser.roles.length > 0
              ? currentUser.roles.map((role) => (
                  <Tag key={role.role_code || role.role_key}>
                    {role.role_name || role.role_key || role.role_code}
                  </Tag>
                ))
              : "-"}
          </Descriptions.Item>
        </Descriptions>
      </ProCard>
    </PageContainer>
  );
};

export default AccountCenter;
