# Changelog

本文件记录 `helloagents/` 知识库（SSOT）相关的所有重要变更。
格式基于 Keep a Changelog，版本号遵循语义化版本（SemVer）。

## [Unreleased]

### 新增
- 新增 Windows 一键重建部署脚本 `docker-redeploy.bat`（基于 Docker Compose，支持 `dev/stable/prod` 与 `--dry-run`）。

### 变更
- 同步上游 `NoFxAiOS/nofx` 的 `dev` 分支最新提交到本地 `my-custom`，完成冲突解决与兼容性修复，并将结果推送到 fork。

## [0.1.0] - 2026-01-11

### 新增
- 初始化 `helloagents/` 知识库目录结构与核心文档（overview/arch/api/data/modules）。
