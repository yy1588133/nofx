# 变更历史索引

本文件记录所有已完成变更（方案包）的索引，便于追溯与查询。

---

## 索引

| 时间戳 | 功能名称 | 类型 | 状态 | 方案包路径 |
|--------|----------|------|------|------------|
| 202601172143 | docker_redeploy_bat | 功能 | ✅已完成 | [2026-01/202601172143_docker_redeploy_bat](2026-01/202601172143_docker_redeploy_bat/) |
| 202601172228 | fix_docker_build_types | 修复 | ✅已完成 | [2026-01/202601172228_fix_docker_build_types](2026-01/202601172228_fix_docker_build_types/) |
| 202601172337 | fix_sqlite_automigrate | 修复 | ✅已完成 | [2026-01/202601172337_fix_sqlite_automigrate](2026-01/202601172337_fix_sqlite_automigrate/) |

---

## 按月归档

> 归档目录结构：`helloagents/history/YYYY-MM/YYYYMMDDHHMM_<feature>/`

### 2026-01

- [202601172143_docker_redeploy_bat](2026-01/202601172143_docker_redeploy_bat/) - 新增 Windows 一键重建部署 bat 脚本
- [202601172228_fix_docker_build_types](2026-01/202601172228_fix_docker_build_types/) - 修复前端类型导致的 Docker 构建失败，并优化构建上下文
- [202601172337_fix_sqlite_automigrate](2026-01/202601172337_fix_sqlite_automigrate/) - 修复 SQLite 旧库在 Docker 启动时 AutoMigrate 触发 *_temp 重建表失败
