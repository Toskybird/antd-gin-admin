# 扫描报告落对象存储（S3 兼容），私有化默认 MinIO，v1 导出 HTML

**Status:** superseded by ADR-0004

扫描报告为可下载的多版本文件，且 API 与 Worker 可能多实例部署。我们决定报告**正文存对象存储（S3 兼容）**，数据库仅存元数据与对象键；私有化默认在 **docker-compose 中附带 MinIO**，并允许改配到外部 S3。v1 报告格式为 **HTML**（可多次生成多版本），暂不做 PDF。不采用仅本地磁盘路径或把大文件塞进 PostgreSQL。
