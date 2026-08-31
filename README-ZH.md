# Starcat Wiki API

<!-- starcat-promo:start -->
<div align="center">
<a href="https://starcat.ink"><img src="https://raw.githubusercontent.com/starcat-app/starcat-pro/main/banner.webp" width="100%" alt="Starcat" /></a>

<p><strong>这是 Starcat 外部文档站索引探测的可自部署支撑服务。</strong></p>
<p>Starcat 是一款原生 macOS 应用，可以把 GitHub Stars 变成可搜索、可整理、可用 AI 追问的本地知识库，并通过桌面客户端、插件、CLI 与可自部署服务组成完整生态。</p>

<a href="https://github.com/starcat-app/homebrew-starcat"><img src="https://img.shields.io/badge/Install%20with-Homebrew-FBBF24?style=for-the-badge&logo=homebrew&logoColor=white" width="220" alt="Install with Homebrew"/></a>
<br/>
<sub><a href="./README.md">English</a></sub>
</div>

<div align="center">
<a href="https://starcat.ink"><img src="https://img.shields.io/badge/website-starcat.ink-38BDF8?style=flat&color=blue" alt="website"/></a>
<a href="https://github.com/starcat-app/starcat-pro"><img src="https://img.shields.io/badge/support-starcat--pro-lightgrey.svg?style=flat&color=blue" alt="support"/></a>
<a href="https://github.com/starcat-app/homebrew-starcat"><img src="https://img.shields.io/badge/install-homebrew-lightgrey.svg?style=flat&color=blue" alt="homebrew"/></a>
<a href="https://github.com/starcat-app/starcat-localization"><img src="https://img.shields.io/badge/localization-open-lightgrey.svg?style=flat&color=blue" alt="localization"/></a>
</div>

<div align="center">
<img width="900" src="https://raw.githubusercontent.com/starcat-app/starcat-pro/main/main.webp" alt="Starcat main window"/>
</div>

**首选 Homebrew 安装：**

```bash
brew tap starcat-app/starcat
brew trust starcat-app/starcat
brew install --cask starcat
```

**相关链接：**

- 官网与下载: https://starcat.ink
- Mac App Store: 搜索 Starcat for GitHub
- 公开支持与发布说明: https://github.com/starcat-app/starcat-pro
- Starcat App Homebrew tap: https://github.com/starcat-app/homebrew-starcat
- CLI / MCP: [starcat-cli](https://github.com/starcat-app/starcat-cli) / [Homebrew tap](https://github.com/starcat-app/homebrew-starcat-cli)
- AI Agent Skill: https://github.com/starcat-app/starcat-skill
- 浏览器插件: [Chrome](https://github.com/starcat-app/starcat-chrome-plugin) / [Safari](https://github.com/starcat-app/starcat-safari-plugin)
- 启动器集成: [Alfred](https://github.com/starcat-app/starcat-alfred-workflow) / [uTools](https://github.com/starcat-app/starcat-utools-plugin) / [Raycast](https://github.com/starcat-app/starcat-raycast-extension)
- 官方文档: https://github.com/starcat-app/starcat-docs
- 官网源码: https://github.com/starcat-app/starcat-site
- 本地化: https://github.com/starcat-app/starcat-localization

**可自部署支撑 API：**

- [starcat-sharing-api](https://github.com/starcat-app/starcat-sharing-api)
- [starcat-trending-api](https://github.com/starcat-app/starcat-trending-api)
- [starcat-weekly-api](https://github.com/starcat-app/starcat-weekly-api)
- [starcat-wiki-api](https://github.com/starcat-app/starcat-wiki-api)
- [starcat-recommend-api](https://github.com/starcat-app/starcat-recommend-api)
- [starcat-discovery-api](https://github.com/starcat-app/starcat-discovery-api)

> Starcat 为普通用户提供默认托管服务。这个 API 开源出来，是为了让进阶用户可以审查实现、本地运行，或部署自己的实例。
<!-- starcat-promo:end -->

外部文档站索引探测服务，探测 DeepWiki / Zread / Google Code Wiki 是否已索引某个 GitHub 仓库，返回跳转链接。

> **v2**（2026-06-11）：状态简化 + probing 预占位 + 按源限速并行探测 + 错误分类重试 + 宕机恢复。

## 特性

- **三源探测**：DeepWiki（json_api）、Zread（json_api）、Google Code Wiki（batchexecute RPC）
- **SWR 缓存**：Stale-While-Revalidate，冷启动同步探测，过期数据异步刷新
- **v2 Batch 异步**：先入库 probing → 秒返响应 → 3 个独立 source worker 并行探测
- **v2 按源限速**：每个 wiki 站独立间隔控制（默认 1s），互不影响
- **v2 错误重试**：network/timeout → 重试，429 → 长间隔重试，403 → 直接放弃
- **v2 宕机恢复**：启动时自动恢复 probing 中的任务
- **Bearer Token 鉴权**：所有 `/api/v1/*` 端点强制鉴权

## 快速开始

### 环境要求

- Go 1.25+

### 本地运行

```bash
cp .env.example .env
# 编辑 .env，填入 API_KEYS
cd starcat-wiki-api
go run ./cmd/server/
```

默认端口 `5004`。

### .env 配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `PORT` | 服务端口 | `5004` |
| `STORE_FILE` | SQLite 数据库路径 | `./wiki.db` |
| `METRICS_STORE_FILE` | 独立请求指标 SQLite 路径 | `./wiki-metrics.db` |
| `API_KEYS` | Bearer Token 白名单（逗号分隔） | 必填 |
| `PROBE_USER_AGENT` | HTTP 请求 UA | Chrome 126 |
| `ENABLE_CODEWIKI_BATCHEXECUTE` | 启用 CodeWiki RPC 精确探测 | `false` |

#### 缓存 TTL

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `CACHE_INDEXED_TTL_HOURS` | indexed 缓存时长 | `168`（7 天） |
| `CACHE_NOT_INDEXED_TTL_HOURS` | not_indexed 缓存时长 | `24` |
| `CACHE_PROBING_TTL_MINUTES` | probing 超时（卡死恢复） | `10` |
| `CACHE_ERROR_TTL_MINUTES` | error 缓存时长 | `30` |

#### 按源限速

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `PROBE_ZREAD_INTERVAL_MS` | zread 请求间隔 (ms) | `1000` |
| `PROBE_DEEPWIKI_INTERVAL_MS` | deepwiki 请求间隔 (ms) | `1000` |
| `PROBE_CODEWIKI_INTERVAL_MS` | codewiki 请求间隔 (ms) | `1000` |

#### 重试策略

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `RETRY_MAX_ATTEMPTS` | 最大重试次数（超过→not_indexed） | `3` |
| `RETRY_INTERVAL_MINUTES` | 重试间隔（分钟，429 用 4 倍） | `30` |

## API 接口

所有数据接口需要 `Authorization: Bearer <api-key>` 头。

### `GET /api/v1/ping`（需鉴权）

返回服务标识，以及由发布 tag 注入的构建版本：

```json
{"schema_version":1,"data":{"service":"wiki","version":"1.2.3","ok":true}}
```

### `GET /api/v1/wikis?owner=X&repo=Y`（需鉴权）

探测单个 GitHub 仓库。

| 参数 | 类型 | 说明 |
|------|------|------|
| `owner` | string | 仓库所有者 |
| `repo` | string | 仓库名 |

探测逻辑（SWR）：

| 场景 | 行为 | `cache_status` |
|------|------|---------------|
| 无缓存 | 同步探测 3 源，写入 DB | `cold` |
| 缓存新鲜 | 直接返回缓存 | `fresh` |
| 缓存过期 | 返回过期缓存 + 后台异步刷新 | `stale` |

**响应 200**：

```json
{
  "schema_version": 1,
  "data": [
    {
      "source": "zread",
      "status": "indexed",
      "url": "https://zread.ai/facebook/react",
      "probeMethod": "json_api",
      "matchedSignals": ["api_status_success"]
    },
    {
      "source": "deepwiki",
      "status": "indexed",
      "url": "https://deepwiki.com/facebook/react",
      "probeMethod": "json_api",
      "matchedSignals": ["api_status_completed"]
    },
    {
      "source": "codewiki",
      "status": "indexed",
      "url": "https://codewiki.google/github.com/facebook/react",
      "probeMethod": "batchexecute_fetch",
      "matchedSignals": ["rpc_ok", "marker_repo_id_matched", "marker_overview_found", "marker_large_payload"]
    }
  ],
  "meta": {
    "generated_at": "2026-06-11T10:00:00+08:00",
    "cache_status": "fresh"
  }
}
```

状态说明：

| status | 含义 |
|--------|------|
| `indexed` | 该 wiki 已收录此仓库 |
| `not_indexed` | 该 wiki 未收录此仓库（或放弃重试后的终态） |
| `probing` | **内部状态**，对外映射为 `not_indexed`，探测完成后更新 |
| `error` | **内部状态**，对外映射为 `not_indexed`，等待定时重试 |

### `POST /api/v1/wikis/batch`（需鉴权）

批量探测（v2 异步）。

**请求体**：

```json
{
  "repos": ["facebook/react", "torvalds/linux", "golang/go"]
}
```

- 最多 50 个 repo
- 接口秒返（→0.02s），后台 3 个 source worker 并行探测
- 已有新鲜缓存的直接返回，无缓存/过期的预入库 `probing` 后异步探测

**响应 200**：

```json
{
  "schema_version": 1,
  "data": {
    "results": {
      "facebook/react": [{ "...": "已有缓存的结果" }]
    },
    "new_probes": 2
  },
  "meta": {
    "generated_at": "2026-06-11T10:00:00+08:00"
  }
}
```

### `POST /internal/sync/probe`（需鉴权，Admin）

手动触发探测同步（预留）。

### `POST /internal/refresh/owner`（需鉴权，Admin）

手动刷新指定 owner 下所有 repo 的探测缓存（预留）。

### `GET /healthz`（公开）

健康检查，返回 `ok`。

## 运营与调用指标

- `GET /internal/stats`：探测覆盖、来源/状态分布、重试积压和新鲜度。
- `GET /internal/probe-errors?limit=1..100`：受限近期探测错误。
- `GET /internal/metrics/{summary,timeseries,routes,status-codes}`：路由调用量、错误与延迟聚合。

接口统一使用 Service API Key；指标不保留凭据或原始请求。

## 鉴权

所有 `/api/v1/*` 和 `/internal/*` 端点需要 `Authorization: Bearer <api-key>` 头。

生成新 key：

```bash
bash ../scripts/gen-api-key.sh
```

## 项目结构

```
starcat-wiki-api/
├── cmd/server/main.go              # 入口：装配 store + probes + scheduler
├── internal/
│   ├── probe/                      # 三源探测实现
│   │   ├── types.go                #   状态、接口、错误分类
│   │   ├── base.go                 #   HTTP 客户端 + UA
│   │   ├── zread.go                #   Zread 探测
│   │   ├── deepwiki.go             #   DeepWiki 探测
│   │   └── codewiki.go             #   Google Code Wiki RPC 探测
│   ├── store/                      # SQLite 持久化
│   │   ├── store.go                #   接口定义
│   │   ├── sqlite.go              #   实现 + v2 新增方法
│   │   └── migrations.go          #   Schema + 兼容迁移
│   ├── handler/                    # HTTP handler
│   │   ├── handler.go             #   工具函数（JSON 响应）
│   │   ├── probe.go               #   探测接口（GET/POST batch）+ 异步 worker
│   │   ├── recover.go             #   启动恢复 probing 任务
│   │   ├── retry.go               #   错误重试定时任务
│   │   └── admin.go               #   Admin endpoint
│   ├── scheduler/                  # cron 定时调度
│   │   └── cron.go                #   过期刷新 + 错误重试
│   ├── middleware/                 # Bearer 鉴权中间件
│   └── model/                      # 数据模型（Envelope）
├── .env.example                    # 配置模板
├── Makefile
└── Dockerfile
```

## 探测流程

```
POST /api/v1/wikis/batch  {"repos": ["a/b", "c/d"]}
│
├─ Phase 1（同步，< 50ms）
│   ├─ 查 DB：有新鲜缓存 → 加入响应
│   └─ 无缓存/过期 → INSERT OR IGNORE 3 条 probing
│
├─ Phase 2（同步）
│   └─ 返回响应 { results: {...}, new_probes: N }
│
└─ Phase 3（异步，后台）
    ├─ [worker.zread]    逐条探测 → sleep 1s → 下一条
    ├─ [worker.deepwiki] 逐条探测 → sleep 1s → 下一条
    └─ [worker.codewiki] 逐条探测 → sleep 1s → 下一条
        │
        └─ 每条结果 → 更新 DB: probing → indexed / not_indexed / error
```

## 部署（Fly.io）

```bash
fly secrets set \
  API_KEYS="sk-starcat-prodKey1,sk-starcat-prodKey2" \
  ENABLE_CODEWIKI_BATCHEXECUTE="true" \
  STORE_FILE="/data/wiki.db" \
  METRICS_STORE_FILE="/data/wiki-metrics.db" \
  -a starcat-wiki-api

fly deploy -a starcat-wiki-api
```

## 技术选型

- **net/http**：Go 标准库，无框架依赖
- **modernc.org/sqlite**：纯 Go SQLite（无 CGO，可交叉编译）
- **robfig/cron/v3**：定时调度（过期清理 + 错误重试）
- **godotenv**：.env 文件加载
