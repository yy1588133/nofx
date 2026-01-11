# 架构设计

本文件用于沉淀系统的总体架构、边界与关键决策索引（ADR）。
如与代码实现不一致，以代码为准并同步更新。

---

## 总体架构（概念图）

> 说明：该图为基于目录结构的概念级视图，用于快速理解系统边界；实际调用链与依赖以代码为准。

```mermaid
flowchart TD
    UI[web/ 前端或交互层] --> API[api/ 接口接入层]
    API --> AUTH[auth/ 认证授权]
    API --> MARKET[market/ 行情与市场]
    API --> TRADER[trader/ 交易执行]
    API --> BACKTEST[backtest/ 回测]
    API --> STORE[store/ 存储]
    API --> DATA[data/ 数据处理]
    API --> LOG[logger/ 日志]

    CFG[config/ 配置] --> API
    SEC[security/ 安全] --> API
    PROVIDER[provider/ 外部适配] --> MARKET
    PROVIDER --> TRADER

    API --> GW[nginx/ 网关与反向代理（如启用）]
    GW --> UI
```

---

## 技术栈（基于仓库结构扫描）

- **核心/后端:** Go
- **前端/控制台:** TypeScript/JavaScript
- **脚本与运维:** Shell/PowerShell、Docker Compose、Nginx 配置
- **辅助:** Python（用于脚本/工具的可能性，具体以实际使用为准）

---

## 核心流程（示例：请求到业务处理）

> 说明：该时序为通用示意，便于后续在变更中补齐“真实路径”。

```mermaid
sequenceDiagram
    participant Client as Client(web/CLI)
    participant API as API(api/)
    participant Auth as Auth(auth/)
    participant Svc as Service(market/trader/...)
    participant Store as Store(store/)

    Client->>API: 发起请求
    API->>Auth: 认证/鉴权（如需要）
    Auth-->>API: 通过/拒绝
    API->>Svc: 执行业务逻辑
    Svc->>Store: 读写数据（如需要）
    Store-->>Svc: 返回结果
    Svc-->>API: 返回业务结果
    API-->>Client: 返回响应
```

---

## 重大架构决策（ADR 索引）

> 完整 ADR 存储在各变更的 `how.md` 中，本章节仅提供索引。

| adr_id | title | date | status | affected_modules | details |
|--------|-------|------|--------|------------------|---------|

