---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
layout: "gallery"
draft: false
description: ""
---

在这里写引导文字。图片请在正文用 Markdown 语法插入：`![可选描述](路径或外链)`，支持相对路径与外链。若不写描述，会用文件名作为说明。正文里的图片只用来生成网格和幻灯片，不会单独展示。

Write intro text here. Add images in the body with standard Markdown: `![optional caption](path-or-url)`. Both relative paths and external URLs are supported. If no caption is provided, the filename is used. Inline images are consumed for the gallery grid and lightbox only.
