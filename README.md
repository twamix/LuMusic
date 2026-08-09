# LuMusic

LuMusic 是基于 Go Music DL 的 Web 音乐搜索、播放、下载和本地曲库管理服务。本仓库只维护 Docker 版本，不包含桌面、Android、iOS 或其他平台的打包产物。

## 快速启动

```bash
mkdir -p data
docker compose up -d
```

启动后访问 <http://localhost:8080/>。旧地址 <http://localhost:8080/music> 仍然兼容。

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

镜像基于 Alpine，内置 `ffmpeg` 和 `ffprobe`，支持 `linux/amd64` 与 `linux/arm64`。容器启动命令为：

```text
./music-dl web --port 8080 --no-browser
```

## 数据与配置

请定期备份 `data/`。如果通过反向代理公开服务，建议在 LuMusic 的 Web 设置中配置管理员账号，并使用 HTTPS 和访问控制。

## 自动发布

推送到 `main` 分支或创建 `v*` 标签后，GitHub Actions 会构建多架构 Docker 镜像并发布到：

```text
ghcr.io/twamix/lumusic
```

## 许可证

本项目遵循仓库中的 [LICENSE](LICENSE) 文件。
