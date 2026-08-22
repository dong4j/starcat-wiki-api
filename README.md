# Starcat Wiki API

<!-- starcat-promo:start -->
<div align="center">
<a href="https://starcat.ink"><img src="https://raw.githubusercontent.com/starcat-app/starcat-pro/main/banner.webp" width="100%" alt="Starcat" /></a>

<p><strong>Self-hostable support API for Starcat external documentation index checks.</strong></p>
<p>Starcat is a native macOS app that turns GitHub Stars into a searchable, organized and AI-assisted local knowledge base. Version 1.4.0 includes README rendering, knowledge-base RAG, GitHub notifications, My Projects, library and repository insights, macOS desktop widgets, tags and private notes, release tracking, repository health signals, AI summaries, semantic search, browser plugins, Alfred / uTools / Raycast search integrations, and self-hostable support APIs.</p>

<a href="https://github.com/starcat-app/homebrew-starcat"><img src="https://img.shields.io/badge/Install%20with-Homebrew-FBBF24?style=for-the-badge&logo=homebrew&logoColor=white" width="220" alt="Install with Homebrew"/></a>
<br/>
<sub><a href="./README-ZH.md">中文说明</a></sub>
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

**Preferred install method:**

```bash
brew tap starcat-app/starcat
brew trust starcat-app/starcat
brew install --cask starcat
```

**Useful links:**

- Home and downloads: https://starcat.ink
- Mac App Store: search for Starcat for GitHub
- Current Direct build: https://starcat.ink/downloads/Starcat-1.4.0-arm64.dmg
- Public support and release notes: https://github.com/starcat-app/starcat-pro
- Starcat App Homebrew tap: https://github.com/starcat-app/homebrew-starcat
- CLI / MCP: [starcat-cli](https://github.com/starcat-app/starcat-cli) / [Homebrew tap](https://github.com/starcat-app/homebrew-starcat-cli)
- AI Agent Skill: https://github.com/starcat-app/starcat-skill
- Browser plugins: [Chrome](https://github.com/starcat-app/starcat-chrome-plugin) / [Safari](https://github.com/starcat-app/starcat-safari-plugin)
- Launcher integrations: [Alfred](https://github.com/starcat-app/starcat-alfred-workflow) / [uTools](https://github.com/starcat-app/starcat-utools-plugin) / [Raycast](https://github.com/starcat-app/starcat-raycast-extension)
- Documentation: https://github.com/starcat-app/starcat-docs
- Website source: https://github.com/starcat-app/starcat-site
- Localization: https://github.com/starcat-app/starcat-localization

**Self-hostable support APIs:**

- [starcat-sharing-api](https://github.com/starcat-app/starcat-sharing-api)
- [starcat-trending-api](https://github.com/starcat-app/starcat-trending-api)
- [starcat-weekly-api](https://github.com/starcat-app/starcat-weekly-api)
- [starcat-wiki-api](https://github.com/starcat-app/starcat-wiki-api)
- [starcat-recommend-api](https://github.com/starcat-app/starcat-recommend-api)
- [starcat-discovery-api](https://github.com/starcat-app/starcat-discovery-api)

> Starcat provides hosted defaults for normal users. This API is open source so advanced users can inspect it, run it locally, or deploy their own instance.
<!-- starcat-promo:end -->

An external documentation index probe that checks whether DeepWiki, Zread, or Google Code Wiki has indexed a GitHub repository and returns the corresponding links.

> **v2** (2026-06-11): simplified statuses + probing placeholders + parallel probes with per-source rate limits + error-aware retries + crash recovery.

## Features

- **Three-source probing**: DeepWiki (`json_api`), Zread (`json_api`), and Google Code Wiki (`batchexecute` RPC)
- **SWR cache**: Stale-While-Revalidate with synchronous cold-start probes and asynchronous refreshes for stale data
- **v2 asynchronous batch processing**: store `probing` placeholders first → return immediately → probe in parallel with three independent source workers
- **v2 per-source rate limits**: each wiki source has an independent request interval (1s by default)
- **v2 error retries**: retry network errors and timeouts → retry 429 responses after a longer interval → stop immediately on 403 responses
- **v2 crash recovery**: automatically resume `probing` tasks at startup
- **Bearer Token authentication**: authentication is required for every `/api/v1/*` endpoint

## Quick Start

### Requirements

- Go 1.25+

### Run Locally

```bash
cp .env.example .env
# Edit .env and set API_KEYS
cd starcat-wiki-api
go run ./cmd/server/
```

The default port is `5004`.

### .env Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `5004` |
| `STORE_FILE` | SQLite database path | `./wiki.db` |
| `API_KEYS` | Comma-separated Bearer Token allowlist | Required |
| `PROBE_USER_AGENT` | HTTP request User-Agent | Chrome 126 |
| `ENABLE_CODEWIKI_BATCHEXECUTE` | Enable precise probing through the CodeWiki RPC | `false` |

#### Cache TTL

| Variable | Description | Default |
|----------|-------------|---------|
| `CACHE_INDEXED_TTL_HOURS` | Cache lifetime for `indexed` results | `168` (7 days) |
| `CACHE_NOT_INDEXED_TTL_HOURS` | Cache lifetime for `not_indexed` results | `24` |
| `CACHE_PROBING_TTL_MINUTES` | `probing` timeout for stalled-task recovery | `10` |
| `CACHE_ERROR_TTL_MINUTES` | Cache lifetime for `error` results | `30` |

#### Per-Source Rate Limits

| Variable | Description | Default |
|----------|-------------|---------|
| `PROBE_ZREAD_INTERVAL_MS` | Zread request interval (ms) | `1000` |
| `PROBE_DEEPWIKI_INTERVAL_MS` | DeepWiki request interval (ms) | `1000` |
| `PROBE_CODEWIKI_INTERVAL_MS` | CodeWiki request interval (ms) | `1000` |

#### Retry Policy

| Variable | Description | Default |
|----------|-------------|---------|
| `RETRY_MAX_ATTEMPTS` | Maximum retry attempts before marking a result `not_indexed` | `3` |
| `RETRY_INTERVAL_MINUTES` | Retry interval in minutes (4× longer for 429 responses) | `30` |

## API

All data endpoints require an `Authorization: Bearer <api-key>` header.

### `GET /api/v1/ping` (Authentication Required)

Returns the service identity and the build version injected from the release tag:

```json
{"schema_version":1,"data":{"service":"wiki","version":"1.2.3","ok":true}}
```

### `GET /api/v1/wikis?owner=X&repo=Y` (Authentication Required)

Probe a single GitHub repository.

| Parameter | Type | Description |
|-----------|------|-------------|
| `owner` | string | Repository owner |
| `repo` | string | Repository name |

Probe behavior (SWR):

| Scenario | Behavior | `cache_status` |
|----------|----------|----------------|
| No cache | Probe all three sources synchronously and write the results to the database | `cold` |
| Fresh cache | Return the cached results immediately | `fresh` |
| Stale cache | Return the stale results and refresh them asynchronously in the background | `stale` |

**200 response**:

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

Status reference:

| status | Meaning |
|--------|---------|
| `indexed` | The wiki has indexed the repository |
| `not_indexed` | The wiki has not indexed the repository, or all retries have been exhausted |
| `probing` | **Internal status** exposed as `not_indexed`; updated when the probe completes |
| `error` | **Internal status** exposed as `not_indexed` while waiting for a scheduled retry |

### `POST /api/v1/wikis/batch` (Authentication Required)

Probe repositories in a batch (asynchronous in v2).

**Request body**:

```json
{
  "repos": ["facebook/react", "torvalds/linux", "golang/go"]
}
```

- Up to 50 repositories
- The endpoint returns immediately (about 0.02s), while three source workers probe in parallel in the background
- Fresh cached results are returned directly; repositories with no cache or stale entries are stored as `probing` before asynchronous probing begins

**200 response**:

```json
{
  "schema_version": 1,
  "data": {
    "results": {
      "facebook/react": [{ "...": "cached result" }]
    },
    "new_probes": 2
  },
  "meta": {
    "generated_at": "2026-06-11T10:00:00+08:00"
  }
}
```

### `POST /internal/sync/probe` (Authentication Required, Admin)

Manually trigger probe synchronization (reserved).

### `POST /internal/refresh/owner` (Authentication Required, Admin)

Manually refresh the probe cache for every repository owned by the specified owner (reserved).

### `GET /healthz` (Public)

Health check that returns `ok`.

## Authentication

All `/api/v1/*` and `/internal/*` endpoints require an `Authorization: Bearer <api-key>` header.

Generate a new key:

```bash
bash ../scripts/gen-api-key.sh
```

## Project Structure

```
starcat-wiki-api/
├── cmd/server/main.go              # Entry point: wires the store, probes, and scheduler
├── internal/
│   ├── probe/                      # Three-source probe implementations
│   │   ├── types.go                #   Statuses, interfaces, and error classification
│   │   ├── base.go                 #   HTTP client and User-Agent
│   │   ├── zread.go                #   Zread probe
│   │   ├── deepwiki.go             #   DeepWiki probe
│   │   └── codewiki.go             #   Google Code Wiki RPC probe
│   ├── store/                      # SQLite persistence
│   │   ├── store.go                #   Interface definitions
│   │   ├── sqlite.go              #   Implementation and methods added in v2
│   │   └── migrations.go          #   Schema and compatibility migrations
│   ├── handler/                    # HTTP handler
│   │   ├── handler.go             #   JSON response helpers
│   │   ├── probe.go               #   Probe endpoints (GET/POST batch) and asynchronous workers
│   │   ├── recover.go             #   Recovery of probing tasks at startup
│   │   ├── retry.go               #   Scheduled error retries
│   │   └── admin.go               #   Admin endpoint
│   ├── scheduler/                  # cron scheduling
│   │   └── cron.go                #   Stale refreshes and error retries
│   ├── middleware/                 # Bearer authentication middleware
│   └── model/                      # Data models (Envelope)
├── .env.example                    # Configuration template
├── Makefile
└── Dockerfile
```

## Probe Flow

```
POST /api/v1/wikis/batch  {"repos": ["a/b", "c/d"]}
│
├─ Phase 1 (synchronous, < 50ms)
│   ├─ Query DB: fresh cache → add to response
│   └─ No cache/stale → INSERT OR IGNORE three probing rows
│
├─ Phase 2 (synchronous)
│   └─ Return { results: {...}, new_probes: N }
│
└─ Phase 3 (asynchronous, background)
    ├─ [worker.zread]    Probe one row → sleep 1s → next row
    ├─ [worker.deepwiki] Probe one row → sleep 1s → next row
    └─ [worker.codewiki] Probe one row → sleep 1s → next row
        │
        └─ For each result → update DB: probing → indexed / not_indexed / error
```

## Deployment (Fly.io)

```bash
fly secrets set \
  API_KEYS="sk-starcat-prodKey1,sk-starcat-prodKey2" \
  ENABLE_CODEWIKI_BATCHEXECUTE="true" \
  STORE_FILE="/data/wiki.db" \
  -a starcat-wiki-api

fly deploy -a starcat-wiki-api
```

## Technology

- **net/http**: Go standard library with no framework dependency
- **modernc.org/sqlite**: pure Go SQLite with no CGO dependency and cross-compilation support
- **robfig/cron/v3**: scheduled stale-data cleanup and error retries
- **godotenv**: `.env` file loading
