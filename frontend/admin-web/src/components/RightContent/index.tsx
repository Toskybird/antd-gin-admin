import { DEFAULT_THEME, LEGACY_THEME_KEY, THEME } from '@/constant/theme';
import {
  MenuOutlined,
  QuestionCircleOutlined,
  SunOutlined,
} from '@ant-design/icons';
import { useModel } from '@umijs/max';
import { Dropdown } from 'antd';
import type { FC } from 'react';
import { useEffect } from 'react';

export type SiderTheme = 'light' | 'dark';

export type LayoutType = 'mix' | 'top';

export const DEFAULT_LAYOUT: LayoutType = 'mix';

export const resolveLayoutType = (): LayoutType => {
  const saved = localStorage.getItem('layout-type');
  if (saved === 'mix' || saved === 'top') {
    return saved;
  }
  return DEFAULT_LAYOUT;
};

export const Question: FC = () => {
  return (
    <a
      href="https://pro.ant.design/docs/getting-started"
      target="_blank"
      rel="noreferrer"
      style={{
        display: 'inline-flex',
        padding: '4px',
        fontSize: '18px',
        color: 'inherit',
      }}
    >
      <QuestionCircleOutlined />
    </a>
  );
};

export const SelectLayout = () => {
  const items = [
    { label: '左侧菜单', key: 'mix' },
    { label: '顶部菜单', key: 'top' },
  ];

  return (
    <Dropdown
      menu={{
        items,
        selectable: true,
        selectedKeys: [resolveLayoutType()],
        onClick: (e) => {
          const key = e.key as LayoutType;
          localStorage.setItem('layout-type', key);
          window.location.reload();
        },
      }}
    >
      <MenuOutlined style={{ marginRight: 8, fontSize: 18, cursor: 'pointer' }} />
    </Dropdown>
  );
};

type ThemeKey = keyof typeof THEME;

export const SelectTheme = ({
  onThemeChange = (_theme?: string) => {},
}: {
  onThemeChange?: (theme?: string) => void | Promise<void>;
}) => {
  const items = [
    { label: '默认主题', key: 'default' },
    { label: '深色主题', key: 'dark' },
    { label: '浅色主题', key: 'light' },
    { label: '夜间主题', key: 'night' },
  ];

  const { initialState } = useModel('@@initialState');

  useEffect(() => {
    let savedTheme = localStorage.getItem('app-theme');
    if (savedTheme === LEGACY_THEME_KEY) {
      savedTheme = DEFAULT_THEME;
      localStorage.setItem('app-theme', DEFAULT_THEME);
    }
    if (savedTheme) {
      document.documentElement.setAttribute('data-theme', savedTheme);
    }
  }, []);

  const selectedTheme =
    (() => {
      const saved = localStorage.getItem('app-theme');
      if (saved === LEGACY_THEME_KEY) {
        return DEFAULT_THEME;
      }
      return saved;
    })() ||
    (initialState?.settings as { currentTheme?: string } | undefined)
      ?.currentTheme ||
    DEFAULT_THEME;

  return (
    <Dropdown
      menu={{
        items,
        selectable: true,
        selectedKeys: [selectedTheme],
        onClick: (e) => {
          const key = e.key as ThemeKey;
          localStorage.setItem('app-theme', key);
          document.documentElement.setAttribute('data-theme', key);
          void onThemeChange(key);
        },
      }}
    >
      <SunOutlined style={{ marginRight: 8, fontSize: 18, cursor: 'pointer' }} />
    </Dropdown>
  );
};
