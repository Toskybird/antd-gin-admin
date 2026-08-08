import { SCAN_CHILD_ROUTE_PATHS } from '@/constants/scanMenu';
import {
  assetStatusText,
  findingSeverityText,
  isScanJobActive,
  scanJobStatusText,
  scanPolicyText,
} from '@/constants/scanStatus';
import routes from '../../../config/routes';

describe('scan pages', () => {
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

  it('keeps SCAN_CHILD_ROUTE_PATHS aligned with menu list routes', () => {
    const scan = findRoute('/scan');
    expect(scan).toBeTruthy();

    const menuPaths = (scan.routes || [])
      .filter(
        (r: any) =>
          typeof r.path === 'string' &&
          r.component &&
          !r.hideInMenu &&
          !r.redirect,
      )
      .map((r: any) => r.path);

    expect(menuPaths).toEqual([...SCAN_CHILD_ROUTE_PATHS]);
  });

  it('registers hideInMenu form and detail routes', () => {
    const scan = findRoute('/scan');
    const hiddenPaths = (scan.routes || [])
      .filter((r: any) => r.hideInMenu)
      .map((r: any) => r.path);

    expect(hiddenPaths).toEqual(
      expect.arrayContaining([
        '/scan/asset/create',
        '/scan/asset/edit/:code',
        '/scan/job/create',
        '/scan/job/detail/:code',
        '/scan/finding/detail/:code',
      ]),
    );
  });
});

describe('scan status helpers', () => {
  it('maps asset status labels', () => {
    expect(assetStatusText('active')).toBe('启用');
    expect(assetStatusText('disabled')).toBe('停用');
    expect(assetStatusText('unknown')).toBe('unknown');
  });

  it('maps scan job status labels', () => {
    expect(scanJobStatusText('queued')).toBe('排队中');
    expect(scanJobStatusText('running')).toBe('执行中');
    expect(scanJobStatusText('succeeded')).toBe('已完成');
  });

  it('detects active scan jobs', () => {
    expect(isScanJobActive('queued')).toBe(true);
    expect(isScanJobActive('running')).toBe(true);
    expect(isScanJobActive('succeeded')).toBe(false);
  });

  it('maps policy and severity labels', () => {
    expect(scanPolicyText('quick')).toBe('快速');
    expect(scanPolicyText('standard')).toBe('标准');
    expect(findingSeverityText('critical')).toBe('严重');
    expect(findingSeverityText('info')).toBe('信息');
  });
});
