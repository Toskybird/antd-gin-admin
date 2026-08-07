import {
  AvatarDropdown,
  AvatarName,
  Footer,
  SelectLayout,
  SelectTheme,
  SelectZoom,
} from '@/components';
import { resolveLayoutType } from '@/components/RightContent';
import { DEFAULT_THEME, LEGACY_THEME_KEY, PRIMARY_COLOR, THEME } from '@/constant/theme';
import { currentUser as queryCurrentUser } from '@/services/ant-design-pro/api';
import { getMenuTree, type Menu as APIMenu } from '@/services/antd-gin-api/menu';
import { applyZoom, resolveZoom } from '@/utils/zoom';
import type {
  Settings as LayoutSettings,
  MenuDataItem,
} from '@ant-design/pro-components';
import '@ant-design/v5-patch-for-react-19';
import type { RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { history } from '@umijs/max';
import { ConfigProvider, theme as antdTheme } from 'antd';
import React from 'react';
import defaultSettings from '../config/defaultSettings';
import { errorConfig } from './requestErrorConfig';

const loginPath = '/user/login';

type ThemeKey = keyof typeof THEME;

type AppSettings = Partial<LayoutSettings> & {
  currentTheme?: string;
  colors?: string;
};

const resolveThemeKey = (): ThemeKey => {
  let saved = localStorage.getItem('app-theme');
  if (saved === LEGACY_THEME_KEY) {
    saved = DEFAULT_THEME;
    localStorage.setItem('app-theme', DEFAULT_THEME);
  }
  if (saved && saved in THEME) {
    return saved as ThemeKey;
  }
  return DEFAULT_THEME as ThemeKey;
};

const buildSettings = (themeKey: ThemeKey): AppSettings =>
  ({
    ...defaultSettings,
    ...THEME[themeKey],
    currentTheme: themeKey,
  }) as AppSettings;

if (typeof document !== 'undefined') {
  document.documentElement.setAttribute('data-theme', resolveThemeKey());
  applyZoom(resolveZoom());
}

/**
 * @see https://umijs.org/docs/api/runtime-config#getinitialstate
 * */
export async function getInitialState(): Promise<{
  settings?: AppSettings;
  currentUser?: API.CurrentUser;
  menuData?: MenuDataItem[];
  loading?: boolean;
  fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
}> {
  const buildMenuData = (menus: APIMenu[]): MenuDataItem[] => {
    return (menus || [])
      .filter((m) => m.status === 1 && m.menu_type !== 'F')
      .map((m) => ({
        name: m.menu_name,
        path: m.path || `/${m.menu_code}`,
        children: m.children ? buildMenuData(m.children) : undefined,
      }));
  };
  const fetchMenuData = async () => {
    try {
      const res = await getMenuTree();
      if (res.code === 200 && Array.isArray(res.data)) {
        return buildMenuData(res.data);
      }
    } catch (_e) {
      // ignore and fallback to static menu
    }
    return [];
  };
  const fetchUserInfo = async () => {
    try {
      const msg = await queryCurrentUser({
        skipErrorHandler: true,
      });
      return msg.data;
    } catch (_error) {
      history.push(loginPath);
    }
    return undefined;
  };

  const currentTheme = resolveThemeKey();
  document.documentElement.setAttribute('data-theme', currentTheme);
  const settings = buildSettings(currentTheme);

  // 如果不是登录页面，执行
  const { location } = history;
  if (
    ![loginPath, '/user/register', '/user/register-result'].includes(
      location.pathname,
    )
  ) {
    const currentUser = await fetchUserInfo();
    const menuData = await fetchMenuData();
    return {
      fetchUserInfo,
      currentUser,
      menuData,
      settings,
    };
  }
  return {
    fetchUserInfo,
    settings,
  };
}

/** 根容器注入 antd 主色/链接色；夜间主题启用 darkAlgorithm */
export const rootContainer = (container: React.ReactNode) => {
  const currentTheme = resolveThemeKey();
  const isNight = currentTheme === 'night';
  const themeConfig = THEME[currentTheme];
  const colorPrimary = themeConfig.colorPrimary || PRIMARY_COLOR;
  const colorLink = themeConfig.colorLink || colorPrimary;
  return (
    <ConfigProvider
      theme={{
        cssVar: true,
        algorithm: isNight
          ? antdTheme.darkAlgorithm
          : antdTheme.defaultAlgorithm,
        token: {
          colorPrimary,
          colorLink,
        },
      }}
    >
      {container}
    </ConfigProvider>
  );
};

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({
  initialState,
}) => {
  const currentTheme =
    (initialState?.settings as AppSettings | undefined)?.currentTheme ||
    DEFAULT_THEME;
  const currentLayout = resolveLayoutType();

  return {
    actionsRender: () => [
      <SelectLayout key="SelectLayout" />,
      <SelectTheme
        key="SelectTheme"
        onThemeChange={async () => {
          // rootContainer 仅在启动时求值，刷新以同步 algorithm / Token
          window.location.reload();
        }}
      />,
      <SelectZoom key="SelectZoom" />,
    ],
    avatarProps: {
      src: initialState?.currentUser?.avatar,
      title: <AvatarName />,
      style: {
        color: 'inherit',
      },
      render: (_, avatarChildren) => {
        return <AvatarDropdown menu>{avatarChildren}</AvatarDropdown>;
      },
    },
    footerRender: () => <Footer />,
    breadcrumbRender: false,
    onPageChange: () => {
      const { location } = history;
      document.documentElement.setAttribute('data-theme', currentTheme);
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.push(loginPath);
      }
    },
    links: [],
    menuHeaderRender: undefined,
    menuDataRender: (_menuData) => {
      if (initialState?.menuData && initialState.menuData.length > 0) {
        return initialState.menuData;
      }
      return [];
    },
    ...initialState?.settings,
    // 覆盖 defaultSettings / 主题合并结果，以用户布局偏好为准
    layout: currentLayout,
  };
};

/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
export const request: RequestConfig = {
  // 开发环境使用代理，生产环境需要配置实际的后端地址
  baseURL:
    process.env.NODE_ENV === 'production'
      ? process.env.REACT_APP_API_BASE_URL || ''
      : '',
  ...errorConfig,
};
