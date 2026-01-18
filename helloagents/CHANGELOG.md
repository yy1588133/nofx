# Changelog

本文件记录 `helloagents/` 知识库（SSOT）相关的所有重要变更。
格式基于 Keep a Changelog，版本号遵循语义化版本（SemVer）。

## [Unreleased]

### 新增
- 新增 Windows 一键重建部署脚本 `docker-redeploy.bat`（基于 Docker Compose，支持 `dev/stable/prod` 与 `--dry-run`）。

### 修复
- 修复前端 Docker 构建因 TypeScript 类型缺失导致的编译失败（补齐 `RiskControlConfig` / `Exchange` 字段）。
- 修复 Docker 环境下 SQLite 旧库启动时 `AutoMigrate` 触发 `*_temp` 重建表失败（`NOT NULL constraint failed`）导致后端容器反复重启的问题：SQLite 改为“仅建表（create-only）”迁移策略，并对关键表做兼容修复。

### 变更
- 优化 Docker 构建上下文：忽略 `data/` 目录，避免传输超大上下文导致构建缓慢或取消。
- 同步上游 `NoFxAiOS/nofx` 的 `dev` 分支最新提交到本地 `my-custom`，完成冲突解决与兼容性修复，并将结果推送到 fork。

## [0.1.0] - 2026-01-11

### 新增
- 初始化 `helloagents/` 知识库目录结构与核心文档（overview/arch/api/data/modules）。
