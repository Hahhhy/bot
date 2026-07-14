# 技术笔记：NapCat 部署方式选型

- **日期**: 2026-07-13
- **决策问题**: NapCat（QQ 连接器）用 Docker 还是 npm 部署？
- **结论**: npm 直接安装

---

## 背景

QQ bot 项目需要 NapCat 作为 QQ 和 Go 程序之间的 WebSocket 桥梁。NapCat 支持两种部署方式：
- Docker（官方镜像 `mlikiowa/napcat-docker`）
- npm 全局安装（经验证 NapCat 不在 npm 上，此方案不可行）

当前环境：Arch Linux，Node.js v26.4.0 已安装，Docker 未安装。

## 决策分析

### 方案对比

| 维度 | npm | Docker |
|------|-----|--------|
| 需安装软件 | 0（Node.js 已有） | Docker ~200MB + daemon + 用户组配置 |
| sudo 次数 | 0 | 3（安装、启动 daemon、usermod） |
| 故障面 | 1 层（NapCat 本身） | 4 层（daemon → 网络 → 容器 → NapCat） |
| 调试 | 直接输出到终端 | `docker logs napcat` |
| 进程管理 | Ctrl+C 即停 | `docker start/stop/restart` |
| 实际可行性 | ❌ NapCat 不在 npm 上 | ✅ |

### 适用场景分析

Docker 解决的核心问题是**环境隔离和跨平台可复现部署**。当前场景：

- NapCat 和 Go bot 同机 localhost 通信
- 无多实例需求
- 无跨平台部署需求

为一个单机 localhost 通信引入 Docker，属于过度设计（违反 YAGNI）。

### YAGNI 原则

> `.clinerules` Rule 5: 不添加参数、选项或抽象 "just in case"。

NapCat 实际分发方式：GitHub Release 下载 zip 包（~400MB），Arch 用户通过 AUR 包 `napcat-qq` 安装。

## 结论修正

**选择 AUR 安装**：`yay -S napcat-qq`

NapCat 不在 npm 上，npm 方案不可行。AUR 的 `napcat-qq` (v4.18.9) 是最新维护的包，由 Arch 社区维护。

下载约 400MB，由于 GitHub 在国内网络不稳定，建议在终端手动执行。

## 参考资料

- NapCat 仓库: https://github.com/NapNeko/NapCatQQ
- 当前项目规划: `QQ bot.md`