---
title: "Hugo 静态站点生成器入门"
date: 2024-01-10T11:00:00+08:00
categories: ["技术"]
tags: ["Hugo", "静态站点", "JAMStack"]
toc: true
---

Hugo 是世界上最快的静态站点生成器，使用 Go 语言编写。

<!--more-->

## 什么是 Hugo？

Hugo 是一个用 Go 语言编写的静态站点生成器，具有以下特点：

- ⚡ **极速构建** - 毫秒级构建速度
- 🎯 **简单易用** - 学习曲线平缓
- 🔧 **功能强大** - 内置丰富功能
- 📦 **单一二进制** - 无需依赖
- 🌐 **多语言支持** - 内置国际化

## 安装 Hugo

### macOS

```bash
brew install hugo
```

### Linux

```bash
# Debian/Ubuntu
sudo apt-get install hugo

# Arch Linux
sudo pacman -S hugo
```

### Windows

```powershell
choco install hugo-extended
```

## 创建站点

### 初始化项目

```bash
# 创建新站点
hugo new site my-blog

# 进入目录
cd my-blog

# 初始化 Git
git init
```

### 添加主题

```bash
# 使用 Resound 主题
git submodule add \
  https://github.com/ifeitao/hugo-theme-resound.git \
  themes/resound
```

### 配置站点

编辑 `hugo.toml`:

```toml
baseURL = 'https://example.com'
languageCode = 'zh-CN'
title = '我的博客'
theme = 'resound'

[params]
  slogan = '记录生活，分享知识'
  author = '你的名字'
```

## 创建内容

### 创建文章

```bash
hugo new posts/my-first-post.md
```

### 编辑文章

```markdown
---
title: "我的第一篇文章"
date: 2024-01-10T10:00:00+08:00
categories: ["技术"]
tags: ["Hugo"]
draft: false
---

这是文章内容。
```

### Front Matter

Hugo 支持三种 Front Matter 格式：

**YAML:**
```yaml
---
title: "文章标题"
date: 2024-01-10
---
```

**TOML:**
```toml
+++
title = "文章标题"
date = 2024-01-10
+++
```

**JSON:**
```json
{
  "title": "文章标题",
  "date": "2024-01-10"
}
```

## 本地开发

### 启动开发服务器

```bash
hugo server -D
```

参数说明：
- `-D` - 显示草稿
- `--navigateToChanged` - 自动导航到修改的页面
- `--disableFastRender` - 禁用快速渲染

### 热重载

Hugo 支持热重载，修改文件后浏览器会自动刷新。

## 构建站点

### 生产构建

```bash
hugo --minify
```

生成的文件在 `public/` 目录。

### 构建选项

```bash
# 指定环境
hugo --environment production

# 生成草稿
hugo -D

# 清理旧文件
hugo --cleanDestinationDir

# 详细输出
hugo -v
```

## 目录结构

```
my-blog/
├── archetypes/       # 内容原型
├── assets/           # 待处理资源
├── content/          # 内容文件
│   └── posts/        # 文章目录
├── data/             # 数据文件
├── layouts/          # 自定义布局
├── public/           # 生成的站点
├── static/           # 静态文件
├── themes/           # 主题目录
│   └── resound/      # Resound 主题
└── hugo.toml         # 配置文件
```

## 内容组织

### 章节（Sections）

```
content/
├── posts/            # 文章章节
│   ├── _index.md
│   └── post-1.md
├── docs/             # 文档章节
│   └── guide.md
└── about.md          # 单页
```

### 分类法（Taxonomies）

默认的分类法：
- categories（分类）
- tags（标签）

自定义分类法：

```toml
[taxonomies]
  category = 'categories'
  tag = 'tags'
  series = 'series'
```

## 短代码（Shortcodes）

### 内置短代码

**YouTube:**
```
{{</* youtube w7Ft2ymGmfc */>}}
```

**Twitter:**
```
{{</* tweet user="hugo" id="123456" */>}}
```

**图片:**
```
{{</* figure src="/images/example.jpg" title="示例图片" */>}}
```

### 自定义短代码

创建 `layouts/shortcodes/note.html`:

```html
<div class="note">
  {{ .Inner }}
</div>
```

使用：
```
{{</* note */>}}
这是一个提示框
{{</* /note */>}}
```

## 数据文件

在 `data/` 目录创建 YAML/JSON/TOML 文件：

```yaml
# data/features.yaml
- name: 响应式
  icon: 📱
- name: 快速
  icon: ⚡
```

在模板中使用：

```go-html-template
{{ range .Site.Data.features }}
  <div>{{ .icon }} {{ .name }}</div>
{{ end }}
```

## 部署

### GitHub Pages

```yaml
# .github/workflows/hugo.yml
name: Deploy Hugo site to Pages

on:
  push:
    branches: ["main"]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: 'latest'
          extended: true
      - run: hugo --minify
      - uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

### Netlify

创建 `netlify.toml`:

```toml
[build]
  publish = "public"
  command = "hugo --minify"

[context.production.environment]
  HUGO_VERSION = "0.120.0"
  HUGO_ENV = "production"
```

### Vercel

在项目设置中配置：
- Build Command: `hugo --minify`
- Output Directory: `public`

## 性能优化

### 图片处理

```go-html-template
{{ $image := resources.Get "images/photo.jpg" }}
{{ $resized := $image.Resize "800x" }}
<img src="{{ $resized.RelPermalink }}">
```

### 资源打包

```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $css := $css | minify | fingerprint }}
<link rel="stylesheet" href="{{ $css.RelPermalink }}">
```

## 最佳实践

1. **使用版本控制** - Git 管理内容
2. **模块化内容** - 使用 partials 和 shortcodes
3. **图片优化** - 使用 Hugo 的图片处理功能
4. **缓存优化** - 启用资源指纹
5. **SEO优化** - 设置完整的 meta 标签

## 学习资源

- [Hugo 官方文档](https://gohugo.io/documentation/)
- [Hugo Discourse](https://discourse.gohugo.io/)
- [Hugo GitHub](https://github.com/gohugoio/hugo)
- [Resound 主题](https://github.com/ifeitao/hugo-theme-resound)

---

**开始使用 Hugo 构建你的站点吧！** 🚀
