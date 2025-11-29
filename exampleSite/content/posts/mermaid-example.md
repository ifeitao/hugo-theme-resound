---
title: "Mermaid 图表示例"
date: 2024-01-13T16:00:00+08:00
categories: ["示例"]
tags: ["Mermaid", "图表", "可视化"]
mermaid: true
toc: true
---

本文展示 Resound 主题中 Mermaid 图表的渲染效果。

<!--more-->

## 流程图

### 简单流程图

```mermaid
graph TD
    A[开始] --> B{是否登录?}
    B -->|是| C[显示主页]
    B -->|否| D[跳转登录页]
    C --> E[结束]
    D --> E
```

### 复杂流程图

```mermaid
graph LR
    A[用户请求] --> B{验证Token}
    B -->|有效| C[查询数据库]
    B -->|无效| D[返回401]
    C --> E{数据存在?}
    E -->|是| F[返回数据]
    E -->|否| G[返回404]
    F --> H[记录日志]
    G --> H
    D --> H
    H --> I[结束]
```

## 时序图

### API调用时序图

```mermaid
sequenceDiagram
    participant C as 客户端
    participant S as 服务器
    participant D as 数据库
    
    C->>S: 发送请求
    activate S
    S->>D: 查询数据
    activate D
    D-->>S: 返回结果
    deactivate D
    S-->>C: 返回响应
    deactivate S
```

### 用户认证流程

```mermaid
sequenceDiagram
    participant U as 用户
    participant F as 前端
    participant B as 后端
    participant DB as 数据库
    
    U->>F: 输入用户名密码
    F->>B: POST /login
    B->>DB: 验证凭据
    DB-->>B: 返回用户信息
    B->>B: 生成Token
    B-->>F: 返回Token
    F-->>U: 登录成功
```

## 类图

```mermaid
classDiagram
    class Animal {
        +String name
        +int age
        +eat()
        +sleep()
    }
    class Dog {
        +String breed
        +bark()
    }
    class Cat {
        +String color
        +meow()
    }
    Animal <|-- Dog
    Animal <|-- Cat
```

## 状态图

```mermaid
stateDiagram-v2
    [*] --> 草稿
    草稿 --> 审核中: 提交审核
    审核中 --> 已发布: 审核通过
    审核中 --> 草稿: 审核拒绝
    已发布 --> 已归档: 归档
    已归档 --> [*]
```

## 甘特图

```mermaid
gantt
    title 项目开发计划
    dateFormat  YYYY-MM-DD
    section 设计阶段
    需求分析           :a1, 2024-01-01, 7d
    UI设计             :a2, after a1, 5d
    section 开发阶段
    前端开发           :b1, after a2, 10d
    后端开发           :b2, after a2, 12d
    section 测试阶段
    集成测试           :c1, after b1 b2, 5d
    上线部署           :c2, after c1, 2d
```

## 饼图

```mermaid
pie title 编程语言使用占比
    "Python" : 35
    "JavaScript" : 30
    "Go" : 20
    "Rust" : 10
    "其他" : 5
```

## Git图

```mermaid
gitGraph
    commit
    commit
    branch develop
    checkout develop
    commit
    commit
    checkout main
    merge develop
    commit
    branch feature
    checkout feature
    commit
    checkout main
    merge feature
    commit
```

## ER图（实体关系图）

```mermaid
erDiagram
    USER ||--o{ ORDER : places
    USER {
        int id
        string name
        string email
    }
    ORDER ||--|{ ORDER_ITEM : contains
    ORDER {
        int id
        date created_at
        string status
    }
    PRODUCT ||--o{ ORDER_ITEM : "ordered in"
    PRODUCT {
        int id
        string name
        decimal price
    }
    ORDER_ITEM {
        int quantity
        decimal price
    }
```

## 旅程图

```mermaid
journey
    title 用户购物流程
    section 浏览
      浏览商品: 5: 用户
      搜索商品: 4: 用户
    section 购买
      加入购物车: 3: 用户
      结算: 2: 用户, 系统
      支付: 2: 用户, 支付系统
    section 售后
      查看订单: 4: 用户
      收货确认: 5: 用户
```

## 使用方法

在文章的 front matter 中启用 Mermaid：

```yaml
---
mermaid: true
---
```

然后使用代码块语法：

````markdown
```mermaid
graph TD
    A --> B
```
````

## 更多资源

- [Mermaid 官方文档](https://mermaid.js.org/)
- [在线编辑器](https://mermaid.live/)
- [语法参考](https://mermaid.js.org/intro/syntax-reference.html)

---

**用 Mermaid 让图表变得简单！** 📊
