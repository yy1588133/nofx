# 任务清单: Windows 一键重建部署脚本

目录: `helloagents/plan/202601172143_docker_redeploy_bat/`

---

## 1. Docker 脚本
- [√] 1.1 新增 `docker-redeploy.bat`，实现 compose 文件选择、参数解析与错误处理，验证 why.md#核心场景-需求-一键重建并部署（Windows）-场景-默认本地构建（dev）
- [√] 1.2 增加 `--help`/`--dry-run` 行为并验证输出，验证 why.md#核心场景-需求-一键重建并部署（Windows）-场景-构建调试（dry-run）

## 2. 文档更新
- [√] 2.1 更新 `helloagents/wiki/modules/docker.md`，补充脚本用法与注意事项，验证 why.md#变更内容
- [√] 2.2 更新 `helloagents/CHANGELOG.md`，记录新增脚本与用途，验证 why.md#变更内容
- [√] 2.3 更新 `helloagents/history/index.md`，增加归档索引，验证 why.md#变更内容

## 3. 安全检查
- [√] 3.1 检查脚本不包含破坏性数据删除操作（`down -v` / `docker system prune` 等），验证 why.md#风险评估

## 4. 测试
- [√] 4.1 运行 `cmd /c docker-redeploy.bat --help` 与 `cmd /c docker-redeploy.bat dev --dry-run`，验证输出与退出码

## 5. 归档
- [√] 5.1 迁移方案包至 `helloagents/history/2026-01/` 并更新索引，完成后在本文件标记任务状态
