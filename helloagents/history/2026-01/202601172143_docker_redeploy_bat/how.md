# 技术设计: Windows 一键重建部署脚本

## 技术方案

### 核心技术
- Docker Desktop + Docker Compose v2（`docker compose`），兼容旧版 `docker-compose` 作为兜底。

### 实现要点
- 使用 `pushd "%~dp0"` 强制脚本在仓库根目录执行。
- 参数解析：
  - 模式：`dev/stable/prod` → 选择 compose 文件
  - 选项：`--no-cache` / `--pull` / `--down` / `--dry-run` / `--help`
- 运行顺序：
  - 可选 `down --remove-orphans`（不删卷）
  - `dev` 模式优先 `build`（可带 `--no-cache` / `--pull`）
  - `stable/prod` 模式默认 `pull`
  - `up -d --remove-orphans`
  - `ps` 输出当前状态（便于排查）

## 架构设计
无。

## 架构决策 ADR
无。

## API设计
无。

## 数据模型
无。

## 安全与性能
- **安全:** 不触碰密钥/PII；不执行危险清理命令（如 `docker system prune`、`down -v`）。
- **性能:** 默认使用 build cache；可通过 `--no-cache` 强制全量构建。

## 测试与部署
- **测试:** `cmd /c docker-redeploy.bat --help` 与 `cmd /c docker-redeploy.bat dev --dry-run`
- **部署:** 通过脚本执行 `docker compose` 完成重建部署。

