# Changelog

## [2.1.0] - 2026-08-27

### Added
- 新增 Wiki 探测状态统计与统一接口调用指标，供 Starcat Admin Console 展示业务数据和调用趋势。

### Changed
- GitHub Actions 不再部署独立 Fly App；官方生产环境统一由 `starcat-api` 聚合服务发布。

## [2.0.0] - 2026-08-07

### Changed
- 导出可装配 `server` 包；依赖 `starcat-api-kit`。
- `/api/v1/ping` 改用 kit `httputil.HandlePingV1`（envelope 契约不变）。
- `server.FromEnv` 改用 kit `env`。

## [Unreleased]

### Added
- **R-03 (2026-06-11)**：新增 `GET /api/v1/ping` 端点，专给 Starcat 客户端「测试连接」按钮用。
  - 走 BearerAuth 中间件，鉴权通过返回 200 + envelope `{data: {service: "wiki", ok: true}}`；
    无效 / 缺失 Key → 401；服务故障 → 5xx。
  - 实现：`internal/handler/ping.go` + `internal/handler/ping_test.go`（7 case）。
  - 设计意图：之前客户端 auth probe 用 `GET /api/v1/wikis` 缺 `owner` / `repo` 参数会触发 400，
    现在统一走 ping，语义更清晰。
  - 跨项目约定：本 `ping.go` 与 trending / weekly / sharing 三个项目「除 import path 外 byte-level 一致」。

## [1.0.0] - 2026-06-10
- 初始化 `starcat-wiki-api` 项目
- 支持 DeepWiki, Zread, Google Code Wiki 探测
- 支持 SWR 分级缓存
