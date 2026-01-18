# docker

## 目的
提供容器化构建与运行配置，支撑本地开发、测试与部署环境的一致性。

## 模块概述
- **目录:** `docker/`
- **职责:** 镜像构建、运行时配置、环境依赖编排（以代码/配置为准）
- **状态:** ✅可用（持续补齐文档）
- **最后更新:** 2026-01-17

## 使用方式

### Windows 一键重建部署
- 入口脚本: 项目根目录 `docker-redeploy.bat`
- 双击运行: 默认按 `docker-compose.yml`（本地源码构建）重建并 `up -d`
- 可选参数:
  - `dev/stable/prod`: 选择对应 compose 文件
  - `--no-cache`: dev 模式全量构建
  - `--pull`: dev 模式拉取基础镜像；stable/prod 拉取服务镜像
  - `--down`: 先 `down`（不删除 volumes）
  - `--dry-run`: 仅打印将要执行的命令，不实际执行

### Compose 文件说明
- `docker-compose.yml`: 本地构建（适合开发/自定义改动后重建部署）
- `docker-compose.stable.yml`: 稳定版镜像（`stable` tag）
- `docker-compose.prod.yml`: 最新版镜像（`latest` tag）

### 构建与数据目录（重要）
- 构建镜像时已通过 `.dockerignore` 忽略 `data/`，避免构建上下文过大导致构建缓慢或失败。
- 运行时数据通过 compose volume 挂载 `./data:/app/data` 持久化，不依赖镜像内置数据目录。

## 规范
- 容器化相关变更应在方案包中记录“如何验证”（例如启动检查点、健康检查、关键端口）。
- 与 `docker-compose*.yml`、`nginx/` 的联动关系应在变更中逐步补齐。

## API接口
- 不适用。

## 数据模型
- 不适用。

## 依赖
- 与部署拓扑相关，参考 `helloagents/wiki/arch.md` 并以配置为准。

## 变更历史
- [202601172143_docker_redeploy_bat](../../history/2026-01/202601172143_docker_redeploy_bat/) - 新增 Windows 一键重建部署脚本入口（bat）。
