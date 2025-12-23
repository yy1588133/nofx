# NOFX 代码审查报告（2025-12-21）

> 范围：后端 Go + 前端 React/Vite（重点阅读：`web/src/components/BacktestPage.tsx`），并对关键安全/架构/性能路径做抽样检查。
>
> 目标：主动发现潜在问题，识别技术债、架构隐患、性能瓶颈、安全漏洞与可维护性风险；给出可执行改进建议与优先级。

---

## 0. 结论摘要（Executive Summary）

### 亮点

- **密码哈希使用 bcrypt**：`auth/auth.go` 使用 `bcrypt.GenerateFromPassword`/`CompareHashAndPassword`，方向正确。
- **存储加密（AES-GCM）+ 可选传输加密（RSA-OAEP + AES-GCM）**：`crypto/crypto.go` 提供了数据库密钥加密与浏览器 → 服务端加密能力。
- **SSRF 防护**：外部数据源请求使用 `security.SafeGet()`，并对私网/保留地址段做了拦截（`security/url_validator.go`）。

### 主要风险（需要优先处理）

- **P0：敏感信息可能被写入日志**（API key / secret / private key 等）。例如：
  - `api/server.go` 在更新模型与交易所配置后记录 `req.Models` / `req.Exchanges`（包含密钥字段），这是典型的“加密做了但日志泄露”的高危问题。
- **P0：公开的解密接口可能成为“解密预言机”**：`/api/crypto/decrypt` 在未认证情况下返回明文（`api/crypto_handler.go` + `api/server.go` 路由注册）。一旦攻击者拿到任何加密载荷（例如通过日志/浏览器扩展/代理等），即可直接调用该接口获得明文。
- **P0：CORS 过宽**：`api/server.go` 的 `corsMiddleware()` 返回 `Access-Control-Allow-Origin: *`。配合 Bearer Token、以及前端 token 存储在 `localStorage`，一旦出现 XSS，将更容易被跨域滥用。
- **P1：SSE（Debate Stream）鉴权实现疑似有缺陷且存在 token 泄露风险**：
  - 前端 `web/src/lib/api.ts` 使用 query string 传 token：`.../stream?token=...`。
  - 后端 `authMiddleware()` 仅从 `Authorization` header 读取 token（未见读取 query token 的实现）。这会导致：
    - 功能上：SSE 可能无法在浏览器原生 `EventSource` 下工作。
    - 安全上：token 放在 URL 会写入日志/历史记录/Referer，属于不推荐实践。

---

## 1. 代码库结构与技术栈（盘点）

- 后端：Go（Gin），SQLite（`modernc.org/sqlite`），WebSocket（`gorilla/websocket`）。入口：`main.go`。
- 前端：React 18 + TypeScript + Vite；状态与数据：Zustand + SWR；图表：Recharts、`lightweight-charts`；动画：framer-motion。
- 部署：Docker + Compose + Nginx。

---

## 2. 高优先级问题清单（按优先级）

> 说明：
>
> - **P0**：可直接导致密钥泄露/账户资金风险/远程攻击面扩大，或严重破坏核心安全假设。
> - **P1**：高概率导致生产事故、明显安全弱点或关键功能不可用。
> - **P2**：技术债/体验/性能/维护成本问题，短期可接受但应排期。

### P0（必须尽快修复）

#### P0-1 日志泄露敏感信息（AI key / Exchange key / Private key）

- 证据：
  - `api/server.go:1495`：`logger.Infof("✓ AI model config updated: %+v", req.Models)`
  - `api/server.go:1616`：`logger.Infof("✓ Exchange config updated: %+v", req.Exchanges)`
  - 这些结构体通常包含 `api_key/secret_key/passphrase/private_key` 等敏感字段。
- 风险：
  - 日志常被收集到集中式平台或磁盘长期留存，一旦泄露等价于“明文存储密钥”。
  - 即便开启传输加密/数据库加密也无法缓解。
- 建议：
  - **立即**：删除或替换上述日志，改为只记录计数、ID、是否启用等不敏感字段。
  - 增加统一的“敏感字段脱敏器”（masker），例如对 `sk-...`、`0x...`、长字符串只保留前后 3~4 位。

#### P0-2 `/api/crypto/decrypt` 未认证的解密接口

- 证据：
  - `api/server.go:107-110`：公开路由 `api.POST("/crypto/decrypt", ...)`。
  - `api/crypto_handler.go:59-77`：直接返回 `plaintext`。
- 风险：
  - 该接口可被用作解密预言机：攻击者只要能获取任何合法加密 payload（例如通过客户端日志、代理、浏览器扩展、错误上报等），即可在无需认证的情况下拿到明文。
  - 目前 payload AAD 中包含 `userId/sessionId/ts/purpose`（见 `web/src/lib/crypto.ts` 与 `crypto/crypto.go`），但后端并未做强校验。
- 建议（任选其一，推荐 A）：
  - A. **移除该对外 API**（最佳）：前端不应需要服务端回传明文。服务端内部解密后直接处理即可。
  - B. 将 `/crypto/decrypt` 挂到 `protected` 组，要求 JWT；并校验 AAD 中的 `userId` 与当前登录用户一致，同时校验 `ts` 与重放窗口。
  - C. 如果仅用于调试，默认禁用且只允许 localhost 或显式 DEBUG 开关启用。

#### P0-3 CORS 过宽导致攻击面扩大

- 证据：`api/server.go:74-87` 的 `corsMiddleware()`：
  - `Access-Control-Allow-Origin: *`
  - `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
  - `Access-Control-Allow-Headers: Content-Type, Authorization`
- 风险：
  - 生产环境一般应限制来源域；过宽 CORS 会在出现 XSS/供应链脚本注入时加速横向利用。
- 建议：
  - 生产环境：改为 allowlist（配置项），例如仅允许 Web UI 域名。
  - 若前后端同域（Nginx 反代），可完全不需要开放跨域。

#### P0-4 默认策略/前端内置第三方 API URL + `auth=cm_...` 固化

- 证据：
  - `web/src/components/strategy/CoinSourceEditor.tsx:6-7`、`web/src/components/strategy/IndicatorEditor.tsx:5-7`
  - `store/strategy.go:228-264` 默认策略配置中也硬编码了 `http://nofxaios.com:30006` 与 `auth=cm_...`。
- 风险：
  - 内置“看似 token/凭证”的参数，可能引发滥用与不可控依赖。
  - HTTP 明文依赖（后端拉取不受浏览器 mixed content 影响，但仍存在链路窃听与内容篡改风险）。
- 建议：
  - 将默认 URL 移到 `.env`/配置文件，并在 UI 中以“示例值”呈现而不是“真实默认”。
  - 对外部数据源做可观测性与降级：失败时明确提示，不要静默使用不可信数据。

---

### P1（尽快修复/明确方案）

#### P1-1 SSE（Debate Stream）鉴权与 token 泄露

- 证据：
  - 前端：`web/src/lib/api.ts:774-777` 使用 `new EventSource(...?token=${token})`。
  - 后端：`api/server.go:2167-2207` 的 `authMiddleware()` 仅从 `Authorization` header 取 token。
- 影响：
  - 功能：浏览器原生 `EventSource` 无法自定义 header，若后端不支持 query token，则 SSE 可能不可用。
  - 安全：token 放 URL 会进日志、Referer、历史记录。
- 建议：
  - A. 改用 **cookie（HttpOnly）** 承载 session，并设置 `SameSite` 与 CSRF 防护；SSE 将自动携带 cookie。
  - B. 或改为 WebSocket，并在握手阶段传 token（仍不推荐 URL）。
  - C. 若必须 query token：后端显式读取并验证，同时关闭 Referer 泄露（加安全头），并对该 token 设置短时效。

#### P1-2 JWT Secret 存在危险默认值

- 证据：`config/config.go:46-48`：若未配置 `JWT_SECRET`，使用固定字符串 `default-jwt-secret-change-in-production`。
- 风险：
  - 用户若误部署缺少 `.env`，将直接导致可伪造 JWT。
- 建议：
  - 启动时强制校验：生产模式下缺失 JWT secret 直接拒绝启动；或至少打印高亮致命警告并退出。

#### P1-3 缺少登录/OTP 等接口的限流与防爆破

- 证据：`api/server.go` 登录/注册/OTP 验证路径未见限流中间件。
- 风险：
  - 暴力破解与撞库；OTP 口令可被频繁尝试。
- 建议：
  - 按 IP + 账号维度的 rate limit（滑动窗口/令牌桶），并对失败计数做退避。

#### P1-4 前端 token 存储在 `localStorage`

- 证据：`web/src/contexts/AuthContext.tsx`、`web/src/lib/httpClient.ts`、`web/src/lib/api.ts` 中多处读取 `localStorage.getItem('auth_token')`。
- 风险：
  - 任何 XSS 都可以直接读取 token 并完全接管账户。
- 建议：
  - 推荐切换到 HttpOnly Cookie；或至少引入 CSP + 依赖审计 + 输出编码，降低 XSS 发生概率。

#### P1-5 注册页面依赖第三方 QR 生成服务（隐私/可用性）

- 证据：`web/src/components/RegisterPage.tsx:382` 使用 `https://api.qrserver.com/...`。
- 风险：
  - 将 OTP QR 内容（otpauth URI）发送到第三方；泄露即等价于 2FA 种子泄露。
- 建议：
  - 客户端本地生成二维码（引入轻量 QR 库），或服务端生成 PNG 并受鉴权保护。

---

### P2（技术债/性能/可维护性）

#### P2-1 `api/server.go` 与 `BacktestPage.tsx` 文件体积过大，维护成本高

- 证据：`api/server.go` 约 2880 行；`web/src/components/BacktestPage.tsx` 约 2549 行。
- 影响：
  - 修改风险大（容易引入回归），测试覆盖难。
- 建议：
  - 后端：按领域拆分 handler（auth、exchange、model、backtest、debate…），保留路由注册聚合。
  - 前端：将 Backtest 拆成容器组件 + 子模块（表单、Runs 列表、Metrics、Charts、Trades、Decisions、Compare 等）。

#### P2-2 Backtest 页面轮询频率较高，可能造成后端压力

- 证据：`web/src/components/BacktestPage.tsx` 多个 `useSWR` 带固定 `refreshInterval`：status 2s；equity/trades/decisions/runs 5s；metrics 10s；aiModels/strategies 30s。
- 备注：
  - 当 `selectedRunId` 存在时会持续轮询，当前代码未根据 `status.state` 动态降频或在完成后停止。
- 建议：
  - 对“运行中”与“已完成”状态使用不同刷新策略；完成后停止轮询。
  - 使用 SSE/WebSocket 推送状态变更。

#### P2-3 图表/marker 对齐已使用二分查找（仍可进一步简化/下沉）

- 证据：
  - `BacktestChart` 与 `CandlestickChartComponent` 使用 `findClosestIndexSorted`（二分查找）在已排序时间数组中定位最近点（见 `web/src/components/BacktestPage.tsx:99`、`:265`、`:516`）。
- 建议：
  - 若数据量继续增长：可考虑后端直接返回对齐后的 marker；或在前端构建更直接的时间索引以减少重复计算。

#### P2-4 `CandlestickChartComponent` 已实现取消/竞态处理（可标记已修复）

- 证据：
  - klines 拉取使用 `AbortController` + requestId（latest wins），并在 cleanup 中 abort（`web/src/components/BacktestPage.tsx:468`）。
  - chart 创建与数据应用拆分为多个 effect，且仅在卸载时 `chart.remove()`（`web/src/components/BacktestPage.tsx:405`、`:457`）。
- 建议：
  - 若后续继续重构：保持“创建 chart / 拉取数据 / 应用数据 / 更新 markers”分离的结构，避免回归。

#### P2-5 Provider 侧排序已使用 `sort.Slice`（可标记已修复）

- 证据：`provider/data_provider.go:176-179` 使用 `sort.Slice` 按 Score 降序排序。
- 建议：
  - 若该条仅用于追踪技术债：可从清单中移除或标注 Done。

---

## 3. 重点文件审查：`web/src/components/BacktestPage.tsx`

### 3.1 架构与可维护性

- 单文件聚合了：表单、轮询数据、图表（Recharts + lightweight-charts）、交易列表、决策日志、对比视图、toast、动画等。
- 建议拆分：
  - `BacktestPage` 只负责 orchestrate（runId、tab、filters）。
  - `BacktestWizard`（步骤表单）
  - `BacktestRunsList`（历史运行列表 + compare）
  - `BacktestOverview`（metrics + positions）
  - `BacktestEquityChart`（Recharts）
  - `BacktestCandlestickChart`（lightweight-charts）
  - `BacktestTradesTable` / `BacktestDecisions`。

### 3.2 性能与体验

- 轮询较多：建议在 `status.state !== running` 时降低频率或停止。
- `CandlestickChartComponent` 每次切换/数据变化会重建 chart；建议复用 chart instance。
- marker 计算建议后端对齐，或前端用二分/索引。

### 3.3 安全与隐私

- Backtest 页本身主要是展示与触发回测，不直接处理密钥，但调用了依赖同一鉴权体系的 API；因此需依赖全局安全治理（CORS、token 存储、CSP、日志脱敏等）。

---

## 4. 后端安全与架构审查（抽样）

### 4.1 认证与会话

- JWT：`auth/auth.go`，过期 24h。
- 黑名单：内存保存，重启丢失（权衡可以接受，但要在文档中明确）。
- 建议：
  - 若要支持强制登出与多实例部署：黑名单应迁移到共享存储（Redis）或改短 TTL + refresh token。

### 4.2 安全头缺失

- Nginx 配置（`nginx/nginx.conf`）未添加 CSP/HSTS/X-Frame-Options 等。
- 建议：
  - 至少：`X-Frame-Options: DENY`、`X-Content-Type-Options: nosniff`、`Referrer-Policy: no-referrer`。
  - 如果继续使用 localStorage token：CSP 必做（并尽量禁用 `unsafe-inline`）。

### 4.3 外部请求与 SSRF

- Provider 使用 `security.SafeGet`，并拦截私网 IP，方向正确。
- 注意：`security.ValidateURL` 对 DNS 解析失败会放行（随后连接阶段可能报错）。这通常 OK，但应记录为设计选择。

---

## 5. 建议落地路线图（可执行）

### 5.1 24 小时内（P0）

1. 移除/脱敏日志：`api/server.go` 中所有可能打印密钥的日志。
2. 移除或加固 `/api/crypto/decrypt`：至少需要认证 + AAD 校验。
3. 生产环境限制 CORS：允许列表化。
4. 移除注册页第三方 QR 服务依赖，改本地生成。

### 5.2 1-2 周内（P1）

1. SSE 鉴权方案定稿（推荐 cookie session），修复 token query 方案。
2. 加入登录/OTP/注册限流。
3. 启动时强制要求 `JWT_SECRET`（以及其它关键安全配置）。

### 5.3 1-2 月（P2）

1. 拆分 `api/server.go` 与 `BacktestPage.tsx`。
2. Backtest 数据推送化（WebSocket/SSE）替换频繁轮询。
3. Provider 排序/HTTP 等实现做性能与可读性优化。

---

## 6. 附录：发现点索引（便于定位）

- CORS：`api/server.go:74-87`
- 公开解密：`api/server.go:107-110` + `api/crypto_handler.go:59-77`
- 可能泄露密钥日志：`api/server.go:1495`、`api/server.go:1616`
- token in URL（SSE）：`web/src/lib/api.ts:774-777`
- 第三方 QR：`web/src/components/RegisterPage.tsx:382`
- 默认第三方数据源 URL：
  - `web/src/components/strategy/CoinSourceEditor.tsx:6-7`
  - `web/src/components/strategy/IndicatorEditor.tsx:5-7`
  - `store/strategy.go:228-264`

---

> 备注：本报告避免记录任何真实密钥内容；请确保 `.env`、数据库文件与日志文件不会进入版本控制或被公开分发。
