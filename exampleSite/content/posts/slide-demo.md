---
title: "Resound 主题特性演示"
date: 2024-01-12T10:00:00+08:00
categories: ["演示"]
tags: ["幻灯片", "Reveal.js", "演示"]
layout: "slide"
theme: "black"
transition: "slide"
titlepage: true
mathjax: true
---

# Resound 主题

一个让你的内容产生回响的 Hugo 主题

..

## 主要特性

- 响应式设计
- 侧边栏导航
- 文章目录 (TOC)
- 站内搜索
- 幻灯片支持

..

## 为什么选择 Resound?

- ⚡ **快速** - Hugo 极速构建
- 🎨 **优雅** - 简洁的视觉设计
- 🔍 **强大** - 丰富的功能特性
- 📱 **响应式** - 完美适配各种设备

...

# 核心功能

..

## 响应式设计

完美适配

- 📱 手机 (< 570px)
- 📱 平板 (570px - 840px)
- 💻 桌面 (840px - 1024px)
- 🖥️ 大屏 (> 1024px)

..

## 搜索功能

- 前端搜索
- 快速响应
- 支持中文
- 多关键词

..

## 文章目录

- 自动生成
- 滚动高亮
- 固定位置
- 响应式显示

...

# 高级特性

..

## 幻灯片模式

基于 Reveal.js

- 多种主题
- 转场效果
- 键盘导航
- 触摸支持

..

## 数学公式

MathJax 支持

质能方程：

$$E = mc^2$$

欧拉公式：

$$e^{i\pi} + 1 = 0$$

..

## 代码高亮

```python
def hello_world():
    print("Hello, Resound!")
    return True
```

```javascript
const greet = (name) => {
  console.log(`Hello, ${name}!`);
};
```

...

# 如何使用

..

## 安装

```bash
git submodule add \
  https://github.com/ifeitao/hugo-theme-resound.git \
  themes/resound
```

..

## 配置

```toml
theme = 'resound'

[params]
  slogan = '让你的内容产生回响'
  author = '你的名字'
```

..

## 创建内容

```bash
# 创建文章
hugo new posts/my-post.md

# 创建幻灯片
hugo new posts/my-slide.md
```

在 front matter 中设置：

```yaml
layout: "slide"
theme: "black"
```

...

# 幻灯片操作

..

## 键盘快捷键

- ← → 左右切换
- ↑ ↓ 上下切换
- Esc 概览模式
- F 全屏模式
- S 演讲者视图

..

## 主题选择

可选主题：

- black (默认)
- white
- league
- sky
- beige
- simple
- serif

..

## 转场效果

可选效果：

- none
- fade
- slide
- convex
- concave
- zoom

...

# 获取帮助

..

## 资源链接

- 📖 [文档](https://github.com/ifeitao/hugo-theme-resound)
- 🐛 [Issues](https://github.com/ifeitao/hugo-theme-resound/issues)
- 💡 [Discussions](https://github.com/ifeitao/hugo-theme-resound/discussions)

..

## 社区

- GitHub: @ifeitao
- 主题: Resound
- 协议: MIT

...

# Thanks!

## 开始你的创作之旅

**Resound - Make Your Content Resound** 🎵

---

按 Esc 查看概览
