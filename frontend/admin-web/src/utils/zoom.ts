export const ZOOM_STORAGE_KEY = 'app-zoom';

export const ZOOM_PRESETS = [80, 90, 100, 110, 125] as const;

export type ZoomPercent = (typeof ZOOM_PRESETS)[number];

export const DEFAULT_ZOOM: ZoomPercent = 100;

const supportsCssZoom = (): boolean => {
  if (typeof document === 'undefined') {
    return false;
  }
  return 'zoom' in document.documentElement.style;
};

export const isZoomPercent = (value: number): value is ZoomPercent =>
  (ZOOM_PRESETS as readonly number[]).includes(value);

export const resolveZoom = (): ZoomPercent => {
  if (typeof localStorage === 'undefined') {
    return DEFAULT_ZOOM;
  }
  const raw = Number(localStorage.getItem(ZOOM_STORAGE_KEY));
  if (Number.isFinite(raw) && isZoomPercent(raw)) {
    return raw;
  }
  return DEFAULT_ZOOM;
};

/** 优先 CSS zoom；旧 Firefox 等用 transform: scale 兜底 */
export const applyZoom = (percent: ZoomPercent): void => {
  if (typeof document === 'undefined') {
    return;
  }
  const root = document.documentElement;
  const ratio = percent / 100;

  if (supportsCssZoom()) {
    root.style.zoom = `${percent}%`;
    root.style.removeProperty('transform');
    root.style.removeProperty('transform-origin');
    root.style.removeProperty('width');
    root.style.removeProperty('height');
    return;
  }

  root.style.removeProperty('zoom');
  root.style.transformOrigin = '0 0';
  root.style.transform = `scale(${ratio})`;
  root.style.width = `${100 / ratio}%`;
  root.style.height = `${100 / ratio}%`;
};

export const setZoom = (percent: ZoomPercent): void => {
  localStorage.setItem(ZOOM_STORAGE_KEY, String(percent));
  applyZoom(percent);
};
