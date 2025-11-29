---
title: "Markdown 语法完整指南"
date: 2024-01-11T09:00:00+08:00
categories: ["教程"]
tags: ["Markdown", "写作", "语法"]
toc: true
---

完整的 Markdown 语法参考指南。

<!--more-->

## 标题

```markdown
# 一级标题
## 二级标题
### 三级标题
#### 四级标题
##### 五级标题
###### 六级标题
```

## 文本样式

**粗体文本**

*斜体文本*

***粗斜体文本***

~~删除线~~

## 列表

### 无序列表

- 项目 1
- 项目 2
  - 子项目 2.1
  - 子项目 2.2
- 项目 3

### 有序列表

1. 第一项
2. 第二项
3. 第三项

### 任务列表

- [x] 已完成任务
- [ ] 未完成任务
- [ ] 待办事项

## 链接和图片

### 链接

[Hugo 官网](https://gohugo.io)

[Resound 主题](https://github.com/ifeitao/hugo-theme-resound)

### 图片

```markdown
![替代文本](/images/example.png)
```

## 引用

> 这是一个引用。
> 
> 可以跨多行。

嵌套引用：

> 第一层引用
>> 第二层引用
>>> 第三层引用

## 代码

### 行内代码

使用 `code` 表示行内代码。

### 代码块

```python
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

print(fibonacci(10))
```

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Resound!")
}
```

```rust
fn main() {
    println!("Hello, Resound!");
}
```

## 表格

| 左对齐 | 居中对齐 | 右对齐 |
|:-------|:--------:|-------:|
| 内容1  | 内容2    | 内容3  |
| A      | B        | C      |
| 较长的内容 | 中等 | 短 |

复杂表格：

| 功能 | Resound | 其他主题 | 说明 |
|------|:-------:|:--------:|------|
| 响应式 | ✅ | ⚠️ | 完美支持 |
| 搜索 | ✅ | ❌ | 前端搜索 |
| TOC | ✅ | ⚠️ | 自动高亮 |
| 幻灯片 | ✅ | ❌ | Reveal.js |

## 分割线

---

***

___

## HTML 标签

可以在 Markdown 中使用 HTML：

<div style="color: red;">
这是红色文本
</div>

<details>
<summary>点击展开</summary>

这是隐藏的内容。

可以包含任何 Markdown 语法。

</details>

## 脚注

这是一个带脚注的文本[^1]。

这是另一个脚注[^note]。

[^1]: 这是脚注的内容。
[^note]: 这是另一个脚注的详细说明。

## 定义列表

术语1
: 定义1

术语2
: 定义2a
: 定义2b

## 缩写

HTML 是超文本标记语言。

*[HTML]: Hyper Text Markup Language

## Emoji 表情

使用 emoji: :smile: :heart: :thumbsup: :rocket:

GitHub 风格: :octocat: :shipit:

## 快捷键

<kbd>Ctrl</kbd> + <kbd>C</kbd> 复制

<kbd>Ctrl</kbd> + <kbd>V</kbd> 粘贴

<kbd>Ctrl</kbd> + <kbd>Z</kbd> 撤销

## 上标和下标

H~2~O (下标)

X^2^ (上标)

## 高亮

==高亮文本==

## 提示框

> **提示**
> 
> 这是一个提示信息。

> **警告**
> 
> 这是一个警告信息。

> **注意**
> 
> 这是一个注意事项。

## 总结

Markdown 是一种轻量级标记语言，易于学习和使用。Resound 主题完美支持所有 Markdown 语法！

---

**Happy Writing!** ✍️
