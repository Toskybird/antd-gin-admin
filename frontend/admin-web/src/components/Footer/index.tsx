import { useLocation } from '@umijs/max';
import { theme } from 'antd';
import React from 'react';

/** 仪表盘全屏沉浸式展示时隐藏全局页脚，避免底边出一条浅色横带 */
const Footer: React.FC = () => {
  const { pathname } = useLocation();
  const { token } = theme.useToken();
  const isDashboard =
    pathname === '/dashboard' || pathname.endsWith('/dashboard');
  if (isDashboard) {
    return null;
  }

  return (
    <footer
      style={{
        padding: '12px 0 20px',
        textAlign: 'center',
        fontSize: 12,
        color: token.colorTextTertiary,
      }}
    >
      项目来源：
      <a
        href="https://github.com/Toskybird/antd-gin-admin"
        target="_blank"
        rel="noreferrer"
        style={{
          color: token.colorLink,
        }}
      >
        GitHub
      </a>
    </footer>
  );
};

export default Footer;
