# Flatpages 使用文档

> Flatpages 是一个简单而强大的静态页面管理系统，支持 Markdown 格式的内容编写，自动生成目录，并提供美观的阅读界面。

## 功能特点

- 支持 Markdown 格式编写内容
- 自动生成文章目录
- 代码高亮显示
- 阅读进度指示
- 响应式设计，支持移动端
- 支持文章导航（上一篇/下一篇）
- 支持文章搜索
- 国际化支持
- 自动计算阅读时间
- 支持多目录配置
- 分页显示

## 使用方法

### 1. 配置启用

在配置文件中启用 flatpages 功能：

```toml
[flatpages]
    # 是否启用flatpages
    enable = true
    # 支持配置多个flatpage目录
    [[flatpages.dirs]]
        # 导航名称，用于显示在导航栏
        nav_name = "帮助文档"
        # 导航路径，用于URL路径，如果不指定则使用文件夹名称
        nav_path = "docs"
        # 文件路径，存放markdown文件的目录
        file_path = "statics/flatpages/docs"
        # 每页显示的条目数，可选，默认为10
        page_size = 20
        # 页面描述，用于SEO
        meta_desc = "帮助文档中心"

    # 可以继续添加更多目录配置...
```

### 2. 创建文章

在配置的 `file_path` 目录下创建 `.md` 文件，按照以下格式编写：

```markdown
# 文章标题

> 文章描述（会显示在列表页）

正文内容...

## 二级标题

### 三级标题

正文内容...
```

### 3. Markdown 格式说明

Flatpages 支持标准 Markdown 语法，包括：

- 标题（H1-H4）
- 列表（有序和无序）
- 代码块（支持语法高亮）
- 引用块
- 链接
- 图片
- 行内代码

代码块示例：

```python
def hello():
    print("Hello, World!")
```

### 4. 特殊功能

#### 代码复制

所有代码块右上角都会自动添加复制按钮，方便用户复制代码。

#### 目录导航

系统会自动根据文章的标题（H2-H4）生成目录，并在右侧显示。目录支持：

- 自动高亮当前阅读位置
- 点击跳转
- 滚动同步

#### 阅读进度

页面顶部会显示阅读进度条，直观展示阅读位置。

#### 阅读时间

系统会自动计算文章的阅读时间（基于每分钟300字的阅读速度）。

## 实现原理

### 1. 文件系统

Flatpages 使用 Go 的标准文件系统操作来管理静态文件。

### 2. 路由注册

系统在启动时通过 `InitFlatpages` 函数注册相关路由：

- `/{nav_path}/` - 文章列表页
- `/{nav_path}/:slug` - 文章详情页

### 3. Markdown 解析

系统会解析每个 Markdown 文件：

1. 从文件名生成 URL slug
2. 提取文章标题（第一个 H1 标题）
3. 提取文章描述（第一个引用块）
4. 计算阅读时间
5. 记录文件修改时间作为更新时间

### 4. 分页实现

列表页使用基于 offset 的分页系统：

- 通过 `offset` 参数控制分页
- 每页显示数量由配置中的 `page_size` 控制
- 支持上一页/下一页导航

### 5. 搜索实现

列表页的搜索功能使用 JavaScript 实现，支持对标题和描述的实时搜索。

### 6. 国际化支持

系统集成了 i18n 支持，可以通过配置启用多语言支持：

```toml
[i18n]
enable = true
```

## 最佳实践

1. **文件命名**
   - 使用有意义的文件名，它将作为 URL 的一部分
   - 避免使用特殊字符和空格
   - 推荐使用小写字母和连字符
2. **内容组织**
   - 每个文件必须有一个 H1 标题
   - 使用引用块来添加文章描述
   - 合理使用二级和三级标题来组织内容
   - 控制单个文件的大小，建议不超过 1000 行
3. **图片处理**
   - 图片建议存放在 `statics/img` 目录下
   - 使用相对路径引用图片
   - 压缩图片以提高加载速度
4. **代码展示**
   - 指定代码块的语言以获得正确的语法高亮
   - 为重要的代码添加注释
   - 确保代码块的缩进正确

## 故障排除

1. **页面未显示**
   - 检查配置文件中 `flatpages.enable` 是否为 true
   - 确认 Markdown 文件位于正确的目录
   - 检查 `nav_path` 配置是否正确
2. **目录未生成**
   - 检查文章是否包含二级或三级标题
   - 确认标题格式正确（## 或 ###）
3. **样式异常**
   - 检查 Markdown 语法是否正确
   - 确认文件编码为 UTF-8
4. **搜索无效**
   - 检查浏览器控制台是否有 JavaScript 错误
   - 确认页面 JavaScript 正确加载
5. **分页问题**
   - 检查 `page_size` 配置是否正确
   - 确认 URL 中的 `offset` 参数格式正确
