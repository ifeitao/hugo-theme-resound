# 部署 Resound 主题示例站点到 Netlify

本指南说明如何将 Resound 主题的 exampleSite 部署到 Netlify。

## 🚀 快速部署

### 方法 1: 通过 Netlify 网站部署

1. **登录 Netlify**
   - 访问 [netlify.com](https://netlify.com)
   - 使用 GitHub 账号登录

2. **导入仓库**
   - 点击 "Add new site" → "Import an existing project"
   - 选择 "Deploy with GitHub"
   - 授权并选择 `ifeitao/hugo-theme-resound` 仓库

3. **配置构建设置**
   
   Netlify 会自动检测 `netlify.toml` 配置文件，无需手动设置。以下是自动应用的配置：
   
   - **Build command**: `cd exampleSite && hugo --gc --minify -b $DEPLOY_PRIME_URL`
   - **Publish directory**: `exampleSite/public`
   - **Hugo version**: `0.152.2`

4. **部署站点**
   - 点击 "Deploy site"
   - 等待构建完成（通常需要 1-2 分钟）

5. **访问站点**
   - 构建完成后，Netlify 会提供一个随机域名（如 `random-name-123456.netlify.app`）
   - 可以在 "Site settings" → "Domain management" 中自定义域名

### 方法 2: 使用 Netlify CLI

如果你已安装 Netlify CLI：

```bash
# 安装 Netlify CLI (如果尚未安装)
npm install -g netlify-cli

# 登录 Netlify
netlify login

# 在主题根目录运行
netlify init

# 按照提示选择:
# - Create & configure a new site
# - 选择你的团队
# - 输入站点名称（可选）

# 部署
netlify deploy --prod
```

## 📝 netlify.toml 配置说明

项目已包含 `netlify.toml` 配置文件：

```toml
[build]
  publish = "exampleSite/public"
  command = "cd exampleSite && hugo --gc --minify -b $DEPLOY_PRIME_URL"

[build.environment]
  HUGO_VERSION = "0.152.2"
  HUGO_ENABLEGITINFO = "true"

[context.production.environment]
  HUGO_ENV = "production"

[context.deploy-preview]
  command = "cd exampleSite && hugo --gc --minify --buildFuture -b $DEPLOY_PRIME_URL"

[context.branch-deploy]
  command = "cd exampleSite && hugo --gc --minify -b $DEPLOY_PRIME_URL"

[[redirects]]
  from = "/*"
  to = "/404.html"
  status = 404
```

### 配置说明

- **publish**: 指定发布目录为 `exampleSite/public`
- **command**: 构建命令，会先进入 exampleSite 目录再执行 hugo
- **HUGO_VERSION**: 指定 Hugo 版本为 0.152.2
- **HUGO_ENABLEGITINFO**: 启用 Git 信息（用于显示最后修改时间）
- **deploy-preview**: 为 PR 预览配置单独的构建命令
- **redirects**: 配置 404 页面重定向

## 🔧 自定义配置

### 更新 baseURL

部署后，需要更新 `exampleSite/hugo.toml` 中的 `baseURL`：

```toml
baseURL = 'https://your-site-name.netlify.app/'
```

或者使用自定义域名：

```toml
baseURL = 'https://your-domain.com/'
```

### 设置环境变量

如果需要额外的环境变量（如分析、评论系统等），在 Netlify 中设置：

1. 进入 Site settings → Build & deploy → Environment
2. 点击 "Edit variables"
3. 添加所需的环境变量

## 🌍 多语言站点

exampleSite 支持中英文双语：

- **中文**: `https://your-site.netlify.app/`
- **英文**: `https://your-site.netlify.app/en/`

## ✨ 功能验证

部署完成后，验证以下功能：

- ✅ 主页显示最新文章
- ✅ 搜索功能正常（需要 `index.json` 生成）
- ✅ 归档页面按时间分组
- ✅ 分类和标签页面
- ✅ 暗黑模式切换
- ✅ 响应式布局
- ✅ 数学公式渲染（MathJax）
- ✅ Mermaid 图表
- ✅ 幻灯片模式

## 🐛 常见问题

### 搜索功能不工作

确保 `exampleSite/hugo.toml` 中包含输出格式配置：

```toml
[outputs]
  home = ["HTML", "RSS", "JSON"]
```

### 样式或脚本未加载

检查 `baseURL` 是否正确设置为你的 Netlify 域名。

### 构建失败

1. 检查 Hugo 版本是否兼容（建议使用 0.152.2+）
2. 查看 Netlify 构建日志获取详细错误信息
3. 本地测试构建命令：
   ```bash
   cd exampleSite && hugo --gc --minify
   ```

## 🔄 持续部署

配置完成后，每次推送到 main 分支都会自动触发部署：

```bash
git add .
git commit -m "Update site content"
git push origin main
```

Netlify 会自动检测更新并重新构建站点。

## 📊 部署状态徽章

在 README.md 中添加部署状态徽章：

```markdown
[![Netlify Status](https://api.netlify.com/api/v1/badges/YOUR-SITE-ID/deploy-status)](https://app.netlify.com/sites/YOUR-SITE-NAME/deploys)
```

在 Netlify 的 Site settings → Status badges 中获取你的徽章代码。

## 🔗 相关链接

- [Netlify 文档](https://docs.netlify.com/)
- [Hugo 部署文档](https://gohugo.io/hosting-and-deployment/hosting-on-netlify/)
- [Resound 主题文档](https://github.com/ifeitao/hugo-theme-resound)

---

如有问题，请提交 [Issue](https://github.com/ifeitao/hugo-theme-resound/issues)。
