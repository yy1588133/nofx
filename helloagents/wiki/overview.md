# nofx

> 本文件包含项目级别的核心信息。详细模块文档见 `modules/` 目录。

---

## 1. 项目概述

### 目标与背景
基于当前仓库结构，项目包含后端（Go）与前端（TypeScript/JavaScript）实现，并按多个业务/基础能力目录拆分（如 `api/`、`auth/`、`market/`、`trader/` 等）。
本概述用于沉淀项目级信息，避免分散在零散文档中。

### 范围
- **范围内:** 系统核心能力（配置、接口接入、鉴权、市场/交易相关能力、存储与工具链等）及其工程化支撑（脚本、部署配置）。
- **范围外:** 与代码无关的运营决策与业务承诺；任何无法通过代码/测试验证的“保证性结论”。

### 干系人
- **负责人:** 未在知识库中登记（建议在首次引入正式流程时补充）。

---

## 2. 模块索引（基于目录命名的初始划分）

> 说明：下表为 `~init` 阶段基于目录名称的初始模块划分，用于建立知识库索引；实际依赖与边界以代码为准，并在后续变更中逐步校准。

| 模块名称 | 职责（初始推断） | 状态（文档） | 文档 |
|---------|------------------|-------------|------|
| api | 对外接口/接入层 | 📝规划中 | [modules/api.md](modules/api.md) |
| auth | 认证与授权 | 📝规划中 | [modules/auth.md](modules/auth.md) |
| backtest | 回测相关能力 | 📝规划中 | [modules/backtest.md](modules/backtest.md) |
| config | 配置与约定 | 📝规划中 | [modules/config.md](modules/config.md) |
| crypto | 加密/密钥相关能力 | 📝规划中 | [modules/crypto.md](modules/crypto.md) |
| data | 数据与数据处理 | 📝规划中 | [modules/data.md](modules/data.md) |
| debate | 讨论/策略评估相关能力（按目录名） | 📝规划中 | [modules/debate.md](modules/debate.md) |
| decision | 决策相关能力（按目录名） | 📝规划中 | [modules/decision.md](modules/decision.md) |
| docker | 容器化与镜像构建 | 📝规划中 | [modules/docker.md](modules/docker.md) |
| docs | 仓库附带文档与说明 | 📝规划中 | [modules/docs.md](modules/docs.md) |
| experience | 经验/实验相关能力（按目录名） | 📝规划中 | [modules/experience.md](modules/experience.md) |
| hook | 钩子/扩展点（按目录名） | 📝规划中 | [modules/hook.md](modules/hook.md) |
| logger | 日志能力 | 📝规划中 | [modules/logger.md](modules/logger.md) |
| manager | 管理/调度能力（按目录名） | 📝规划中 | [modules/manager.md](modules/manager.md) |
| market | 市场数据与行情相关能力 | 📝规划中 | [modules/market.md](modules/market.md) |
| mcp | MCP 相关能力（按目录名） | 📝规划中 | [modules/mcp.md](modules/mcp.md) |
| nginx | Nginx 配置与网关相关 | 📝规划中 | [modules/nginx.md](modules/nginx.md) |
| provider | 外部提供方/适配层（按目录名） | 📝规划中 | [modules/provider.md](modules/provider.md) |
| scripts | 脚本与自动化 | 📝规划中 | [modules/scripts.md](modules/scripts.md) |
| security | 安全与合规模块 | 📝规划中 | [modules/security.md](modules/security.md) |
| store | 存储/持久化相关能力 | 📝规划中 | [modules/store.md](modules/store.md) |
| trader | 交易执行相关能力 | 📝规划中 | [modules/trader.md](modules/trader.md) |
| web | 前端与交互界面 | 📝规划中 | [modules/web.md](modules/web.md) |

---

## 3. 快速链接

- [技术约定](../project.md)
- [架构设计](arch.md)
- [API 手册](api.md)
- [数据模型](data.md)
- [变更历史](../history/index.md)

