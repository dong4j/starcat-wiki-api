# AGENTS.md — starcat-wiki-api

> **唯一协作规范源**：本仓库根目录 `AGENTS.md` 是本项目协作规范的唯一正文维护源。
> 开工前还必须阅读并遵守上级 [`../AGENTS.md`](../AGENTS.md) 的跨仓规则。

## 项目概述

外部文档站索引探测服务：检测 DeepWiki / Zread / Google Code Wiki 是否已索引某 GitHub 仓库，返回跳转链接。v2 采用 SWR 缓存、probing 预占位、按源限速并行探测、错误分类重试与宕机恢复。生产经 `starcat-api` 聚合部署。

## 技术栈

- Go 1.25.0 · `net/http`
- `modernc.org/sqlite` · `github.com/robfig/cron/v3`
- `github.com/starcat-app/starcat-api-kit` v0.3.0
- `github.com/joho/godotenv`

## 关键目录

```
cmd/server/           # 入口
server/               # 可导出装配
internal/probe/       # deepwiki / zread / codewiki 探测器
internal/store/       # SQLite + migrations
internal/scheduler/   # cron 清理与错误重试
internal/handler/     # REST + admin；启动期 probing 恢复见 recover.go
Makefile              # 构建产物名为 bin/wiki（非 server）
```

## 开发与测试命令

```bash
cp .env.example .env          # 填入 API_KEYS
make run                      # go run ./cmd/server，PORT=5004
make build                    # 输出 bin/wiki
make check                    # fmt-check + vet + test（PR 前）
make docker-run               # 映射 5004:5004
```

CI（`.github/workflows/go.yml`）：与 sharing 系列相同（gofmt · vet · docker build · test -race）。

环境变量见 `.env.example`：`STORE_FILE`、`METRICS_STORE_FILE`、`API_KEYS`、`PROBE_USER_AGENT`、`ENABLE_CODEWIKI_BATCHEXECUTE`、`CACHE_*`、`PROBE_*_INTERVAL_MS`、`RETRY_*`。

## 代码与架构约束

- **鉴权**：所有 `/api/v1/*` 必须 Bearer；`/healthz` 不鉴权。
- **三源探测**：DeepWiki（json_api）、Zread（json_api）、CodeWiki（batchexecute RPC 或 URL fallback）。
- **v2 异步**：batch 先写 probing → 秒返 → 独立 source worker 并行探测；启动期卡死 probing 由 `internal/handler/recover.go` 恢复。
- **按源限速**：各 wiki 站独立间隔（默认 1000ms），互不影响。
- **错误策略**：network/timeout 重试；429 长间隔；403 放弃。
- **背景任务**：`server.Service.StartBackground()` 只启动 **cron** 与 **启动期 probing 恢复**（`RecoverPendingProbes`），进程生命周期内**只能调用一次**。独立入口 `cmd/server/main.go` 已调用；`starcat-api` 聚合入口当前**未**调用——聚合模式下 cron 与启动期恢复未启用，但**请求触发的异步 batch 探测仍正常运行**；勿误判背景能力已全部失效或已全部启用。
- 本服务 **不需要** `GITHUB_TOKENS`；trending/weekly 可通过 `WIKI_API_URL`+`WIKI_API_KEY` 触发预热。

## 安全与数据边界

- 禁止入库：`.env`、`wiki.db`、`wiki-metrics.db`、`bin/`、`logs/`。
- `PROBE_USER_AGENT` 使用浏览器 UA 以降低反爬拦截；勿在日志输出完整 API Key。

## 部署与发布禁令

未经 dong4j 明确授权，禁止：`make release`、`scripts/deploy.sh`、`fly deploy`、`git push`/`git tag`。生产 Fly 仅经 `starcat-api`。
