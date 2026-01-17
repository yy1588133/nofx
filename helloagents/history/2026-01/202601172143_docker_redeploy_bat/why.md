# 变更提案: Windows 一键重建部署脚本

## 需求背景
当前项目主要通过 `docker compose` 启动（`docker-compose.yml` 本地构建；`docker-compose.*.yml` 镜像部署）。在 Windows 下，常见需求是通过双击执行脚本完成“切到项目目录 → 重新构建镜像 → 重建容器并后台运行”，减少手工输入命令与路径错误。

## 变更内容
1. 新增 `docker-redeploy.bat`：双击即可基于 Docker Compose 重建并部署本项目。
2. 支持可选参数（`dev/stable/prod`、`--no-cache`、`--pull`、`--down`、`--dry-run`），兼顾一键与可控。
3. 同步更新知识库中 Docker 模块文档与变更记录。

## 影响范围
- **模块:** docker / scripts（运行脚本）
- **文件:**
  - `docker-redeploy.bat`
  - `helloagents/wiki/modules/docker.md`
  - `helloagents/CHANGELOG.md`
  - `helloagents/history/index.md`
- **API:** 无
- **数据:** 无（默认不删除 volumes，仅重建容器）

## 核心场景

### 需求: 一键重建并部署（Windows）
**模块:** docker
在 Windows 文件管理器中双击脚本即可完成重建部署，且脚本能自动切换到仓库根目录，避免在错误目录执行导致构建失败。

#### 场景: 默认本地构建（dev）
开发者在仓库根目录双击 `docker-redeploy.bat`。
- 预期结果：使用 `docker-compose.yml` 重新构建后端/前端镜像并 `up -d`，服务启动后可访问端口（前端默认 3000）。

#### 场景: 指定镜像部署（stable/prod）
用户在命令行执行 `docker-redeploy.bat stable` 或 `docker-redeploy.bat prod`。
- 预期结果：使用对应 compose 文件拉取镜像并启动容器。

#### 场景: 构建调试（dry-run）
用户执行 `docker-redeploy.bat dev --dry-run`。
- 预期结果：仅打印将要执行的 compose 命令，不实际执行。

## 风险评估
- **风险:** 误删除数据（容器重建影响 volumes）。
- **缓解:** 脚本默认不执行 `docker compose down -v`；即便使用 `--down` 也仅执行 `down`（不删除 volumes）。

