export default [
  {
    path: '/user',
    layout: false,
    routes: [{ path: '/user/login', component: './user/login' }],
  },
  {
    path: '/account',
    hideInMenu: true,
    access: 'authenticated',
    routes: [
      { path: '/account/center', component: './account/center' },
      { path: '/account/password', component: './account/password' },
    ],
  },
  {
    path: '/welcome',
    icon: 'smile',
    component: './Welcome',
    access: 'canViewWelcome',
  },
  {
    path: '/dashboard',
    icon: 'dashboard',
    component: './Dashboard',
    access: 'canViewDashboard',
  },
  {
    path: '/monitor',
    icon: 'monitor',
    access: 'canViewMonitorSection',
    routes: [
      {
        path: '/monitor/online-users',
        component: './monitor/online-users',
        access: 'canViewMonitor',
      },
      {
        path: '/monitor/service',
        component: './monitor/service',
        access: 'canViewMonitor',
      },
      {
        path: '/monitor/cache',
        component: './monitor/cache',
        access: 'canViewCache',
      },
      {
        path: '/monitor/data',
        component: './monitor/data',
        access: 'canViewData',
      },
    ],
  },
  { path: '/cockpit', redirect: '/dashboard' },
  {
    path: '/scan',
    icon: 'BugOutlined',
    access: 'canViewScanSection',
    routes: [
      { path: '/scan', redirect: '/scan/asset' },
      {
        path: '/scan/asset',
        component: './scan/asset',
        access: 'canViewScanAsset',
      },
      {
        path: '/scan/asset/create',
        component: './scan/asset/form',
        hideInMenu: true,
        access: 'canViewScanAsset',
      },
      {
        path: '/scan/asset/edit/:code',
        component: './scan/asset/form',
        hideInMenu: true,
        access: 'canViewScanAsset',
      },
      {
        path: '/scan/job',
        component: './scan/job',
        access: 'canViewScanJob',
      },
      {
        path: '/scan/job/create',
        component: './scan/job/create',
        hideInMenu: true,
        access: 'canViewScanJob',
      },
      {
        path: '/scan/job/detail/:code',
        component: './scan/job/detail',
        hideInMenu: true,
        access: 'canViewScanJob',
      },
      {
        path: '/scan/finding',
        component: './scan/finding',
        access: 'canViewScanFinding',
      },
      {
        path: '/scan/finding/detail/:code',
        component: './scan/finding/detail',
        hideInMenu: true,
        access: 'canViewScanFinding',
      },
      {
        path: '/scan/detection-rule',
        component: './scan/detection-rule',
        access: 'canViewScanDetectionRule',
      },
      {
        path: '/scan/report',
        component: './scan/report',
        access: 'canViewScanReport',
      },
    ],
  },
  {
    path: '/system',
    icon: 'setting',
    access: 'canViewSystem',
    routes: [
      {
        path: '/system/user',
        component: './system/user',
        access: 'canViewSystemUser',
      },
      {
        path: '/system/role',
        component: './system/role',
        access: 'canViewSystemRole',
      },
      {
        path: '/system/menu',
        component: './system/menu',
        access: 'canViewSystemMenu',
      },
      {
        path: '/system/dept',
        component: './system/dept',
        access: 'canViewSystemDept',
      },
      {
        path: '/system/operation-log',
        component: './system/operation-log',
        access: 'canViewOperationLog',
      },
    ],
  },
  { path: '/', redirect: '/welcome' },
  { path: '*', layout: false, component: './404' },
];
