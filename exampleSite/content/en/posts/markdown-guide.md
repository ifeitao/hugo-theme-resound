---
title: "Complete Markdown Syntax Guide"
date: 2024-01-11T09:00:00+08:00
categories: ["Tutorial"]
tags: ["Markdown", "Writing", "Syntax"]
toc: true
---

Complete Markdown syntax reference guide.

<!--more-->

## Headers

```markdown
# H1 Header
## H2 Header
### H3 Header
#### H4 Header
##### H5 Header
###### H6 Header
```

## Text Styles

**Bold Text**

*Italic Text*

***Bold and Italic Text***

~~Strikethrough~~

## Lists

### Unordered Lists

- Item 1
- Item 2
  - Sub-item 2.1
  - Sub-item 2.2
- Item 3

### Ordered Lists

1. First item
2. Second item
3. Third item

### Task Lists

- [x] Completed task
- [ ] Incomplete task
- [ ] To-do item

## Links and Images

### Links

[Hugo Official Site](https://gohugo.io)

[Resound Theme](https://github.com/ifeitao/hugo-theme-resound)

### Images

```markdown
![Alt Text](/images/example.png)
```

## Blockquotes

> This is a quote.
> 
> It can span multiple lines.

Nested quotes:

> First level quote
>> Second level quote
>>> Third level quote

## Code

### Inline Code

Use `code` for inline code.

### Code Blocks

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

## Tables

| Left Align | Center Align | Right Align |
|:-----------|:------------:|------------:|
| Content 1  | Content 2    | Content 3   |
| A          | B            | C           |
| Long content | Medium     | Short       |

Complex table:

| Feature | Resound | Other Themes | Description |
|---------|:-------:|:------------:|-------------|
| Responsive | ✅ | ⚠️ | Perfect support |
| Search | ✅ | ❌ | Frontend search |
| TOC | ✅ | ⚠️ | Auto highlight |
| Slides | ✅ | ❌ | Reveal.js |

## Horizontal Rules

---

***

___

## HTML Tags

You can use HTML in Markdown:

<div style="color: red;">
This is red text
</div>

<details>
<summary>Click to expand</summary>

This is hidden content.

Can contain any Markdown syntax.

</details>

## Footnotes

This is text with a footnote[^1].

This is another footnote[^note].

[^1]: This is the content of the footnote.
[^note]: This is another detailed explanation of the footnote.

## Definition Lists

Term 1
: Definition 1

Term 2
: Definition 2a
: Definition 2b

## Abbreviations

HTML is Hyper Text Markup Language.

*[HTML]: Hyper Text Markup Language

## Emoji

Using emoji: :smile: :heart: :thumbsup: :rocket:

GitHub style: :octocat: :shipit:

## Keyboard Keys

<kbd>Ctrl</kbd> + <kbd>C</kbd> Copy

<kbd>Ctrl</kbd> + <kbd>V</kbd> Paste

<kbd>Ctrl</kbd> + <kbd>Z</kbd> Undo

## Superscript and Subscript

H~2~O (subscript)

X^2^ (superscript)

## Highlight

==Highlighted text==

## Callouts

> **Tip**
> 
> This is a tip message.

> **Warning**
> 
> This is a warning message.

> **Note**
> 
> This is a note message.

## Summary

Markdown is a lightweight markup language that is easy to learn and use. Resound theme perfectly supports all Markdown syntax!

---

**Happy Writing!** ✍️