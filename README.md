# LuMusic

LuMusic 是基于 Go Music DL 的 Web 音乐搜索、播放、下载和本地曲库管理服务。本仓库只维护 Docker 版本，不包含桌面、Android、iOS 或其他平台的打包产物。

当前版本已同步上游 `v1.1.2`。累计包含 WebDAV 同步、播放速度控制、批量复制下载链接、基础路径配置，以及排行榜（多平台榜单聚合）、咪咕音频流下载与播放、下载去重与文件存在性校验、设置弹窗标签页重构、汽水音乐会员整曲取流等更新。

## 快速启动

```bash
mkdir -p data
docker compose up -d
```

启动后直接访问 <http://localhost:8080/>。首页不会跳转，旧地址 <http://localhost:8080/music> 仍然兼容。

默认 Compose 配置使用 `ghcr.io/twamix/lumusic:latest`，并把当前目录的 `data/` 挂载到容器 `/home/appuser/data`，用于保存下载文件、曲库索引、设置和登录信息。

## 从源码构建

```bash
docker compose -f docker-compose.dev.yml up -d --build
```

也可以直接构建镜像：

```bash
docker build -t lumusic:local .
docker run -d --name LuMusic \
  -p 8080:8080 \
  -v "${PWD}/data:/home/appuser/data" \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  lumusic:local
```

镜像基于 Alpine，内置 `ffmpeg` 和 `ffprobe`，支持 `linux/amd64` 与 `linux/arm64`。

## 基础路径

默认内部 Web 路径为 `/music`，但根路径 `/` 会直接显示可用首页。反向代理需要自定义二级路径时，可以覆盖容器启动命令：

```bash
./music-dl web --port 8080 --no-browser --base-path /dl
```

## WebDAV 同步

Web 设置中可以开启“下载时同步到 WebDAV”，并配置 WebDAV 地址、用户名、密码和远端目录。密码保存在服务端 `data/settings.db`，设置接口不会回显密码。

## 数据与配置

请定期备份 `data/`。如果通过反向代理公开服务，建议配置管理员账号，并使用 HTTPS 和访问控制。普通搜索、播放、下载和歌单浏览保持公开，系统设置和 Cookie 管理需要管理员登录。

## 自动发布

推送到 `main` 分支或创建 `v*` 标签后，GitHub Actions 会构建多架构 Docker 镜像并发布到：

```text
ghcr.io/twamix/lumusic
```

## 许可证

本项目遵循 GNU Affero General Public License v3.0，详情见 [LICENSE](LICENSE)。
