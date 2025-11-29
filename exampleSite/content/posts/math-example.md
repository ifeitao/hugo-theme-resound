---
title: "数学公式示例"
date: 2024-01-14T14:00:00+08:00
categories: ["示例"]
tags: ["MathJax", "数学", "公式"]
mathjax: true
toc: true
---

本文展示 Resound 主题中 MathJax 数学公式的渲染效果。

<!--more-->

## 行内公式

在文本中可以使用行内公式，例如：质能方程 $E = mc^2$ 是物理学中最著名的公式之一。

勾股定理：$a^2 + b^2 = c^2$

欧拉公式：$e^{i\pi} + 1 = 0$

## 块级公式

### 积分

高斯积分：

$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$

分部积分：

$$
\int u dv = uv - \int v du
$$

### 求和

$$
\sum_{i=1}^{n} i = \frac{n(n+1)}{2}
$$

无穷级数：

$$
\sum_{n=0}^{\infty} \frac{1}{2^n} = 2
$$

### 矩阵

$$
\begin{pmatrix}
a & b \\
c & d
\end{pmatrix}
\begin{pmatrix}
x \\
y
\end{pmatrix}
=
\begin{pmatrix}
ax + by \\
cx + dy
\end{pmatrix}
$$

单位矩阵：

$$
I = \begin{pmatrix}
1 & 0 & 0 \\
0 & 1 & 0 \\
0 & 0 & 1
\end{pmatrix}
$$

### 方程组

线性方程组：

$$
\begin{cases}
x + y = 5 \\
2x - y = 1
\end{cases}
$$

### 微分方程

$$
\frac{dy}{dx} + P(x)y = Q(x)
$$

### 极限

$$
\lim_{x \to \infty} \frac{1}{x} = 0
$$

$$
\lim_{n \to \infty} \left(1 + \frac{1}{n}\right)^n = e
$$

## 复杂公式

### 傅里叶变换

$$
F(\omega) = \int_{-\infty}^{\infty} f(t) e^{-i\omega t} dt
$$

### 薛定谔方程

$$
i\hbar\frac{\partial}{\partial t}\Psi(\mathbf{r},t) = \hat{H}\Psi(\mathbf{r},t)
$$

### 麦克斯韦方程组

$$
\begin{align}
\nabla \cdot \mathbf{E} &= \frac{\rho}{\epsilon_0} \\
\nabla \cdot \mathbf{B} &= 0 \\
\nabla \times \mathbf{E} &= -\frac{\partial \mathbf{B}}{\partial t} \\
\nabla \times \mathbf{B} &= \mu_0\mathbf{J} + \mu_0\epsilon_0\frac{\partial \mathbf{E}}{\partial t}
\end{align}
$$

### 贝叶斯公式

$$
P(A|B) = \frac{P(B|A)P(A)}{P(B)}
$$

## 特殊符号

### 希腊字母

$\alpha, \beta, \gamma, \delta, \epsilon, \zeta, \eta, \theta$

$\Alpha, \Beta, \Gamma, \Delta, \Epsilon, \Zeta, \Eta, \Theta$

$\pi, \phi, \psi, \omega, \Omega, \Sigma, \Lambda$

### 运算符

$\times, \div, \pm, \mp, \cdot, \ast, \circ$

$\leq, \geq, \neq, \approx, \equiv, \sim$

$\in, \notin, \subset, \supset, \subseteq, \supseteq$

$\cup, \cap, \emptyset, \infty, \partial, \nabla$

### 箭头

$\rightarrow, \leftarrow, \Rightarrow, \Leftarrow, \leftrightarrow, \Leftrightarrow$

$\uparrow, \downarrow, \Uparrow, \Downarrow$

## 使用方法

在文章的 front matter 中启用 MathJax：

```yaml
---
mathjax: true
---
```

然后使用 `$...$` 表示行内公式，使用 `$$...$$` 表示块级公式。

---

**MathJax 让数学公式变得优雅！** ✨
