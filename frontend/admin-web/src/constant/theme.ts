export const PRIMARY_COLOR = '#4366B0';
export const DEFAULT_THEME = 'default';

/** 旧键兼容：localStorage 曾使用 webray */
export const LEGACY_THEME_KEY = 'webray';

export const THEME = {
  light: {
    navTheme: 'light',
    colors: 'customize',
    colorPrimary: PRIMARY_COLOR,
    colorLink: PRIMARY_COLOR,
    token: {
      pageContainer: {
        paddingBlockPageContainerContent: 16,
        paddingInlinePageContainerContent: 24,
      },
      header: {
        colorHeaderTitle: PRIMARY_COLOR,
      },
    },
  },
  dark: {
    navTheme: 'light',
    colorPrimary: '#1890ff',
    colorLink: '#1890ff',
    colors: 'default',
    token: {
      // 参见ts声明，demo 见文档，通过token 修改样式
      //https://procomponents.ant.design/components/layout#%E9%80%9A%E8%BF%87-token-%E4%BF%AE%E6%94%B9%E6%A0%B7%E5%BC%8F
      colorText: 'rgba(255, 255, 255, 0.85)', // 主文本色
      colorTextSecondary: 'rgba(255, 255, 255, 0.65)', // 次要文本色
      colorTextTertiary: 'rgba(255, 255, 255, 0.45)', // 第三级文本色
      colorTextQuaternary: 'rgba(255, 255, 255, 0.25)', // 第四级文本色

      // 特定组件文本颜色
      colorTextBase: 'rgba(255, 255, 255, 0.85)', // 基础文本色
      colorTextDescription: 'rgba(255, 255, 255, 0.65)', // 描述文本色
      colorTextDisabled: 'rgba(255, 255, 255, 0.25)', // 禁用状态文本色
      header: {
        // 头部背景
        colorBgHeader: '#001529',
        colorHeaderTitle: '#fff',
        colorTextMenu: 'rgb(255, 255, 255,0.65)',
        colorTextRightActionsItem: 'rgb(255, 255, 255,0.65)',
        colorBgRightActionsItemHover: 'rgba(199, 221, 255, 0.08)',
      },
      pageContainer: {
        paddingBlockPageContainerContent: 16,
        paddingInlinePageContainerContent: 24,
      },
      sider: {
        colorMenuBackground: '#001529', //背景
        colorTextMenu: 'rgb(255, 255, 255,0.65)', //menuItem 的字体颜
        colorBgMenuItemCollapsedElevated: '#001529 !import', //收起 menuItem 的弹出菜单背景颜色
        colorTextMenuSecondary: 'rgb(255, 255, 255,0.65)',
        colorTextMenuSelected: '#fff',
        colorTextMenuItemHover: '#fff',
        colorBgMenuItemSelected: '#1890ff',
      },
    },
  },
  default: {
    navTheme: 'light',
    colors: 'customize',
    colorPrimary: PRIMARY_COLOR,
    colorLink: PRIMARY_COLOR,
    token: {
      // 参见ts声明，demo 见文档，通过token 修改样式
      //https://procomponents.ant.design/components/layout#%E9%80%9A%E8%BF%87-token-%E4%BF%AE%E6%94%B9%E6%A0%B7%E5%BC%8F
      colorText: 'rgba(255, 255, 255, 0.85)', // 主文本色
      colorTextSecondary: 'rgba(255, 255, 255, 0.65)', // 次要文本色
      colorTextTertiary: 'rgba(255, 255, 255, 0.45)', // 第三级文本色
      colorTextQuaternary: 'rgba(255, 255, 255, 0.25)', // 第四级文本色

      // 特定组件文本颜色
      colorTextBase: 'rgba(255, 255, 255, 0.85)', // 基础文本色
      colorTextDescription: 'rgba(255, 255, 255, 0.65)', // 描述文本色
      colorTextDisabled: 'rgba(255, 255, 255, 0.25)', // 禁用状态文本色
      header: {
        // 头部背景
        colorBgHeader: '#001529',
        colorHeaderTitle: '#fff',
        colorTextMenu: 'rgb(255, 255, 255,0.65)',
        colorTextMenuActive: '#fff',
        colorTextMenuItemHover: '#fff',
        colorTextMenuSelected: PRIMARY_COLOR,
        colorTextRightActionsItem: 'rgb(255, 255, 255,0.65)',
      },
      pageContainer: {
        paddingBlockPageContainerContent: 16,
        paddingInlinePageContainerContent: 24,
      },
      sider: {
        colorMenuBackground: '#001529', //背景
        colorTextMenu: 'rgb(255, 255, 255,0.65)', //menuItem 的字体颜
        colorBgMenuItemCollapsedElevated: '#001529', //收起 menuItem 的弹出菜单背景颜色
        colorTextMenuSecondary: 'rgb(255, 255, 255,0.65)',
        colorTextMenuActive: '#fff',
        colorBgRightActionsItemHover: 'rgba(199, 221, 255, 0.08)',
        colorTextMenuItemHover: '#fff',
        colorTextMenuSelected: '#fff',
        colorBgMenuItemSelected: PRIMARY_COLOR,
      },
    },
  },
  // mix: {
  //   navTheme: 'light',
  //   colorPrimary: '#1890ff',
  //   colors: 'default',
  //   token: {
  //     colorText: 'rgba(255, 255, 255, 0.85)', // 主文本色
  //     colorTextSecondary: 'rgba(255, 255, 255, 0.65)', // 次要文本色
  //     colorTextTertiary: 'rgba(255, 255, 255, 0.45)', // 第三级文本色
  //     colorTextQuaternary: 'rgba(255, 255, 255, 0.25)', // 第四级文本色

  //     // 特定组件文本颜色
  //     colorTextBase: 'rgba(255, 255, 255, 0.85)', // 基础文本色
  //     colorTextDescription: 'rgba(255, 255, 255, 0.65)', // 描述文本色
  //     colorTextDisabled: 'rgba(255, 255, 255, 0.25)', // 禁用状态文本色
  //     header: {
  //       // 头部背景
  //       colorBgHeader: 'transparent',
  //     },
  //     sider: {
  //       colorMenuBackground: '#001529', //背景
  //       colorTextMenu: 'rgb(255, 255, 255,0.65)', //menuItem 的字体颜
  //       colorBgMenuItemCollapsedElevated: '#001529 !import', //收起 menuItem 的弹出菜单背景颜色
  //       colorTextMenuSecondary: 'rgb(255, 255, 255,0.65)',
  //       colorTextMenuSelected: '#fff',
  //       colorBgMenuItemSelected: '#1890ff',
  //     },
  //   },
  // },
  night: {
    navTheme: 'dark',
    colorPrimary: PRIMARY_COLOR,
    colorLink: PRIMARY_COLOR,
    colors: 'night',
    token: {
      bgLayout: 'rgb(5, 6, 8)',
      colorText: 'rgba(255, 255, 255, 0.85)', // 主文本色
      colorTextSecondary: 'rgba(255, 255, 255, 0.65)', // 次要文本色
      colorTextTertiary: 'rgba(255, 255, 255, 0.45)', // 第三级文本色
      colorTextQuaternary: 'rgba(255, 255, 255, 0.25)', // 第四级文本色

      header: {
        // 头部背景
        colorBgHeader: 'rgb(5, 6, 8)',
        colorHeaderTitle: 'rgba(199, 221, 255, 0.85)',
        colorTextMenu: 'rgba(199, 221, 255, 0.85)',
        colorTextRightActionsItem: 'rgba(199, 221, 255, 0.85)',
      },
      pageContainer: {
        colorBgPageContainer: 'rgb(5, 6, 8)',
      },
      sider: {
        colorMenuBackground: 'rgb(5, 6, 8)', //背景
        colorTextMenu: 'rgba(199, 221, 255, 0.85)', //menuItem 的字体颜
        colorBgMenuItemCollapsedElevated: 'rgb(5, 6, 8) !import', //收起 menuItem 的弹出菜单背景颜色
        colorTextMenuSecondary: 'rgba(199, 221, 255, 0.85)',
        colorTextMenuSelected: 'rgba(199, 221, 255)',
        colorBgMenuItemSelected: PRIMARY_COLOR,
        colorTextMenuItemHover: 'rgba(199, 221, 255)',
      },
    },
  },
};
