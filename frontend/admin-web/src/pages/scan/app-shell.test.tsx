import access from '@/access';
import { SCAN_CHILD_ROUTE_PATHS } from '@/constants/scanMenu';
import Loading from '@/loading';
import routes from '../../../config/routes';

describe('web-admin app-shell scan section (default)', () => {
  const findRoute = (path: string, list: any[] = routes as any[]): any | undefined => {
    for (const route of list) {
      if (route.path === path) {
        return route;
      }
      if (route.routes) {
        const found = findRoute(path, route.routes);
        if (found) {
          return found;
        }
      }
    }
    return undefined;
  };

  it('exposes 漏洞扫描 child routes in PRD order', () => {
    const scan = findRoute('/scan');
    expect(scan).toBeTruthy();
    expect(scan.access).toBe('canViewScanSection');

    const childPaths = (scan.routes || [])
      .filter((r: any) => typeof r.path === 'string' && r.component)
      .map((r: any) => r.path);

    expect(childPaths).toEqual([...SCAN_CHILD_ROUTE_PATHS]);
  });

  it('grants scan section access from scan list permissions', () => {
    const result = access({
      currentUser: {
        name: '管理员',
        isSuperAdmin: false,
        permissions: [
          'scan:asset:list',
          'scan:job:list',
          'scan:finding:list',
          'scan:detection-rule:list',
          'scan:report:list',
        ],
      },
    });

    expect(result.canViewScanSection).toBe(true);
    expect(result.canViewScanAsset).toBe(true);
    expect(result.canViewScanJob).toBe(true);
    expect(result.canViewScanFinding).toBe(true);
    expect(result.canViewScanDetectionRule).toBe(true);
    expect(result.canViewScanReport).toBe(true);
  });
});

describe('web-admin app-shell scan section (loading)', () => {
  it('route loading uses Skeleton', () => {
    expect(Loading).toBeTruthy();
    // Observable contract: page-level loading is antd Skeleton (DESIGN §5).
    expect(Loading.toString()).toContain('Skeleton');
  });
});
