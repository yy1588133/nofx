# 任务清单: 修复 Docker 构建失败（前端类型 + 构建上下文）

目录: `helloagents/plan/202601172228_fix_docker_build_types/`

---

## 1. 前端类型修复（阻断构建）
- [√] 1.1 在 `web/src/types.ts` 补齐 `RiskControlConfig` 缺失字段（`min_hold_minutes`、`cooldown_after_close_minutes`、`min_tp_cost_multiplier`、`min_stop_loss_distance_pct`），确保 `npm run build` 通过
- [√] 1.2 在 `web/src/types.ts` 为 `Exchange` 增加 `initial_balance?: number`（paper 类型兼容），修复 TS2339

## 2. Docker 构建加速（减少上下文体积）
- [√] 2.1 更新 `.dockerignore` 忽略 `data/`，避免传输超大构建上下文导致构建卡顿/取消

## 3. 知识库同步
- [√] 3.1 更新 `helloagents/wiki/modules/web.md`，记录本次构建修复点
- [√] 3.2 更新 `helloagents/wiki/modules/docker.md`，补充 `.dockerignore` 与数据目录说明

## 4. 验证
- [√] 4.1 执行 `docker compose -f docker-compose.yml build nofx-frontend`
- [√] 4.2 执行 `docker compose -f docker-compose.yml build nofx`（验证构建上下文与后端镜像构建流程）

## 5. 归档
- [√] 5.1 迁移方案包至 `helloagents/history/2026-01/` 并更新 `helloagents/history/index.md`
