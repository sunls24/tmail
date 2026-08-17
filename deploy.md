# 🧰 自建部署教程 (v2.0.0+)

## 邮件接收方式

支持以下两种方式，均会将原始邮件转发到 tmail 的 `/api/report` 接口：

- Cloudflare Email Routing + Worker
- 服务器直接运行 `tmail-smtpd`

两种方式可以单独使用，也可以同时使用。

## 使用 Cloudflare Worker

- 首先开启邮件转发，按照官方流程来就行

- 创建一个 Workers，模板随便选都可以

![workers-create](doc/workers-create.webp)

创建好之后点击`Code editor`编辑代码，将[此处](doc/workers.js)的代码粘贴进去，需要将其中的域名`mail.sunls.de`替换为自己的，然后别忘记点击`Save and deploy`部署:

![workers-edit](doc/workers-edit.webp)

- 然后需要添加一条`Catch-All`的规则，注意要选择`Send to a Worker`，如图:

![email-routing.png](doc/email-routing.webp)

## 直接通过 SMTP 接收邮件

运行镜像内的 `tmail-smtpd` 即可直接接收 SMTP 邮件。它需要与 tmail 主服务配置相同的 `REPORT_HMAC_SECRET`：

```shell
docker run --name tmail-smtpd -d --restart unless-stopped \
  --network host \
  -e 'TMAIL_REPORT_URL=http://127.0.0.1:3000/api/report' \
  -e 'REPORT_HMAC_SECRET=your-report-hmac-secret' \
  -e 'DOMAIN_LIST=example.com' \
  sunls24/tmail:latest /app/tmail-smtpd
```

- `SMTP_ADDR`：SMTP 监听地址，默认 `:25`
- `TMAIL_REPORT_URL`：必填，tmail 的完整 `/api/report` 地址
- `DOMAIN_LIST`：可选；配置后只接收列表内的域名，不配置则不限制收件域名
- `REPORT_MAX_BODY_SIZE`：单封邮件大小上限，默认 `268435456`（256 MiB）

域名需要将 MX 记录指向运行 SMTP 服务的服务器，并确保公网 TCP 25 端口可以访问。非 root 用户直接运行二进制时，需要授予监听低位端口的权限，或通过端口映射将公网 25 转发到其他监听端口。

SMTP 服务每个事务只接受一个收件人；标准邮件服务器会自动为其他收件人创建后续投递事务。

## 环境变量配置

### 数据库配置

**目前仅支持 PostgreSQL**

- `DB_HOST`: 数据库地址
- `DB_PORT`: 数据库地址端口
- `DB_PASS`: 数据库密码
- `DB_NAME`: 数据库名称，默认`tmail`

### 必须

- `DOMAIN_LIST`: 支持的域名列表，使用`,`分割，例如: `isco.eu.org,chato.eu.org`

### 非必须

- `ADMIN_ADDRESS`: 管理员邮箱地址，可以查看所有邮件
- `HOST`: 服务监听地址，默认为`127.0.0.1`
- `PORT`: 服务监听端口，默认为`3000`
- `API_KEY`: API 调用密钥；启用人机验证时，可通过 `X-API-Key` 请求头跳过验证
- `REPORT_HMAC_SECRET`: Worker 或 SMTP 服务调用 `/api/report` 时使用的共享密钥；配置后 `/api/report` 要求签名
- `REPORT_MAX_BODY_SIZE`: 单封邮件请求体上限，单位字节，默认 `268435456`（256 MiB）
- `TURNSTILE_SITE_KEY`: Cloudflare Turnstile Site Key
- `TURNSTILE_SECRET_KEY`: Cloudflare Turnstile Secret Key
- `TURNSTILE_COOKIE_TTL`: 启用人机验证时的 Cookie 有效时间，默认为`6h`
- `DEBUG`: 本地 HTTP 调试时设置为`true`，生产环境不要开启

`TURNSTILE_SITE_KEY` 和 `TURNSTILE_SECRET_KEY` 必须同时设置；两者都不设置时关闭人机验证。

启用 `REPORT_HMAC_SECRET` 后，需要在 Worker 或 SMTP 服务中配置相同的密钥。

本地开发可以使用 Cloudflare 官方测试密钥：

```text
TURNSTILE_SITE_KEY=1x00000000000000000000AA
TURNSTILE_SECRET_KEY=1x0000000000000000000000000000000AA
DEBUG=true
```

## 部署

_请修改其中的环境变量配置_

### Docker

```shell
docker run --name tmail -d --restart unless-stopped -e 'DB_HOST=127.0.0.1' -e 'DB_PASS=postgres' -e 'HOST=0.0.0.0' -e 'DOMAIN_LIST=isco.eu.org,chato.eu.org' -e 'REPORT_HMAC_SECRET=your-report-hmac-secret' -e 'TURNSTILE_SITE_KEY=your-site-key' -e 'TURNSTILE_SECRET_KEY=your-secret-key' -p 3000:3000 sunls24/tmail:latest
```

### Docker Compose & Caddy (推荐)

_如果不需要反向代理，需要设置`HOST=0.0.0.0`环境变量_

**docker-compose.yaml**

```yaml
version: "3.0"

services:
  tmail:
    container_name: tmail
    image: sunls24/tmail:latest
    network_mode: host
    restart: unless-stopped
    environment:
      - "DB_HOST=127.0.0.1"
      - "DB_PASS=postgres"
      - "DOMAIN_LIST=isco.eu.org,chato.eu.org"
      - "REPORT_HMAC_SECRET=your-report-hmac-secret"
      - "TURNSTILE_SITE_KEY=your-site-key"
      - "TURNSTILE_SECRET_KEY=your-secret-key"
    volumes:
      - ./tmail:/app/fs
```

**Caddyfile**

```text
mail.example.com {
	encode zstd gzip
	@cache path /_astro/* /*.webp /favicon.svg
	header @cache Cache-Control "public, max-age=31536000, immutable"
	reverse_proxy 127.0.0.1:3000
}
```
