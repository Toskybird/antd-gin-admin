# 前端主题自参考工程一次对齐移植并自管

`admin-web` 需要全局风格与多主题能力，且与 `frontend/web_react_code` 技术栈一致。我们决定从该参考工程**一次拷贝**主题体系（含 `theme.ts`、ProLayout Token、`global.less` 主题段、顶栏切换与 `localStorage` 持久化），迁入后由 `admin-web` **自管**，不抽共享包、不持续双向同步。

主题键名为 `light` / `dark` / `default` / `night`（默认壳原参考键 `webray` 已重命名为 `default`；读到旧值 `webray` 时迁移写回 `default`）。默认主题（`default`）的 `colorPrimary` 与 `colorLink` 固定为 `#4366B0`；`dark` 可保留 Ant 蓝 `#1890ff`。主色经 `rootContainer` 的 `ConfigProvider`（含 `cssVar`）注入。

登录页保持现状；除 `/dashboard` 外的壳内页面按 Design Token 与 `createStyles` 消除硬编码色。未选择「仅改主色」或「抽 monorepo 共享主题包」，是因为前者达不到对齐目标，后者超出本轮「只换皮」范围且参考工程并非本产品运行时依赖。
