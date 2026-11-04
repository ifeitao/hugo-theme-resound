# brand/ — 主题默认图标的源文件

主题自带的 `static/` 图标出自次生进化（SecEvo Labs）的标记「此身 · 年轮 · 火种」：
三圈等距外扩，三个缺口同指 315°，火种红圆点留在圈内。200 单位画布，线宽 9 / 12.15。

用主题搭自己的站，直接替换 `static/` 下的同名文件即可，这一目录可以整个删掉。

| 文件 | 用途 |
|---|---|
| `icon-tile.svg` | 进化绿圆角方 + 三圈，生成 `favicon.svg`、`favicon-32x32.png`、`favicon-96x96.png` |
| `icon-tile-mini.svg` | 同一枚方标，只留一环，用于 ≤16px 档（三圈在这个尺寸会糊成一团灰） |
| `icon-disc.svg` | 满版绿圆章 + 纸白细描边，生成 `logo.png`（页眉 64px 圆形）与安装图标 |

改完源文件后重新生成 `static/` 下的全部位图：

```sh
python3 brand/make-icons.py
```

需要本机装有 Google Chrome（无头渲染）和 Pillow（拼装 `.ico`）。
