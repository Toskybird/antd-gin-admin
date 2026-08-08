import { useLocation } from '@umijs/max';
import React from 'react';

/** 仪表盘全屏沉浸式展示时隐藏全局页脚，避免底边出一条浅色横带 */
const Footer: React.FC = () => {
  const { pathname } = useLocation();
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
      }}
    />
  );
};

export default Footer;
