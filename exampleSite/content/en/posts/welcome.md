---
title: "Welcome to Resound Theme"
date: 2024-01-15T10:00:00+08:00
categories: ["Tutorial"]
tags: ["Hugo", "Resound", "Getting Started"]
toc: true
---

Welcome to Resound - a Hugo theme that makes your content resound!

<!--more-->

## Theme Features

Resound is a feature-rich and elegantly designed Hugo blog theme with the following features:

### Core Features

1. **Responsive Design** - Perfectly adapted for mobile, tablet, and desktop
2. **Sidebar Navigation** - Drawer-style design with multi-level menus
3. **Table of Contents** - Automatically generated with scroll highlighting
4. **Site Search** - Fast and accurate frontend search
5. **Archive and Categories** - Grouped by year and month for easy browsing

### Advanced Features

- 🎤 **Slide Mode** - Based on Reveal.js
- 📊 **Math Formulas** - MathJax support
- 📈 **Chart Rendering** - Mermaid support
- 💬 **Comment System** - Giscus integration
- 📊 **Analytics** - Google Analytics

## Quick Start

### Install Theme

```bash
git submodule add https://github.com/ifeitao/hugo-theme-resound.git themes/resound
```

### Configure Site

Set in `hugo.toml`:

```toml
theme = 'resound'

[params]
  slogan = 'Your site slogan'
  author = 'Your name'
```

### Create Posts

```bash
hugo new posts/my-first-post.md
```

## Usage Tips

### Enable Table of Contents

Add to the front matter of your post:

```yaml
toc: true
```

### Code Highlighting

Syntax highlighting support for multiple programming languages:

```python
def hello_world():
    print("Hello, Resound!")
    return True
```

```javascript
function greet(name) {
  console.log(`Hello, ${name}!`);
  return true;
}
```

### Table Support

| Feature | Support | Description |
|---------|---------|-------------|
| Responsive | ✅ | Fully supported |
| Search | ✅ | Frontend search |
| Comments | ✅ | Giscus |
| Slides | ✅ | Reveal.js |

### Lists and Quotes

**Unordered List:**

- First item
- Second item
  - Sub-item 2.1
  - Sub-item 2.2
- Third item

**Ordered List:**

1. First
2. Second
3. Last

**Quote:**

> Resound - Make Your Content Resound
>
> Let your content resound

## Advanced Features

Check other sample posts to learn more about advanced features:

- [Math Formula Example](/en/posts/math-example/)
- [Mermaid Chart Example](/en/posts/mermaid-example/)
- [Slide Demo](/en/posts/slide-demo/)

## Get Help

- 📖 [Full Documentation](https://github.com/ifeitao/hugo-theme-resound)
- 🐛 [Report Issues](https://github.com/ifeitao/hugo-theme-resound/issues)
- 💡 [Feature Suggestions](https://github.com/ifeitao/hugo-theme-resound/discussions)

---

**Start your creative journey!** 🚀