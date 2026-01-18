# 任务清单: fix_sqlite_automigrate

目录: `helloagents/plan/202601172337_fix_sqlite_automigrate/`

> 说明: 本次为“轻量迭代”方案包，仅包含 task.md。

---

## 1. 问题定位与复现
- [√] 1.1 复现 Docker 容器启动时 SQLite 迁移失败（`*_temp` 表 `NOT NULL constraint failed`）。
- [√] 1.2 核对容器实际使用的数据库文件路径与挂载（`./data` → `/app/data`，`DB_PATH` 默认 `data/data.db`）。

## 2. SQLite 迁移兼容修复（Store 层）
- [√] 2.1 在 `store/` 增加 SQLite “create-only” 迁移辅助，避免对已存在表执行 `AutoMigrate` 改表导致失败。
- [√] 2.2 针对 `users` 表增加兼容修复（补列/修复空值/去重/补索引），避免遗留数据触发迁移失败。
- [√] 2.3 针对 `ai_models` 表增加兼容修复（补列/修复空值/补索引），避免遗留数据触发迁移失败。
- [√] 2.4 将核心表初始化迁移切换为 “create-only” 行为（仅缺表时建表，不改已有表）。

## 3. 验证
- [√] 3.1 执行 `docker compose -f docker-compose.yml build nofx` 构建后端镜像。
- [√] 3.2 执行 `docker compose -f docker-compose.yml up -d --force-recreate nofx`，验证后端容器健康并成功启动 API。

## 4. 知识库同步
- [√] 4.1 更新 `helloagents/wiki/modules/store.md`，记录 SQLite 迁移兼容策略与注意事项。
- [√] 4.2 更新 `helloagents/CHANGELOG.md`、`helloagents/history/index.md` 并归档方案包至 `helloagents/history/`。

