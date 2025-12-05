# 贡献指南 / Contributing Guide

感谢你对本项目的关注！/ Thank you for your interest in this project!

## 中文

### 如何贡献

我们欢迎所有形式的贡献，包括但不限于：

- 🐛 报告 Bug
- ✨ 提出新功能
- 📝 改进文档
- 💻 提交代码
- 🌐 翻译

### 贡献流程

1. **Fork 仓库**
   - 点击右上角的 Fork 按钮

2. **克隆仓库**
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-async-python-asyncio.git
   cd go-async-python-asyncio
   ```

3. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **进行修改**
   - 遵循现有的代码风格
   - 添加必要的注释
   - 更新相关文档

5. **测试代码**
   
   Go 示例：
   ```bash
   cd go-async/your-example
   go run main.go
   ```
   
   Python 示例：
   ```bash
   cd python-asyncio/your-example
   python your_script.py
   ```

6. **提交更改**
   ```bash
   git add .
   git commit -m "描述你的更改"
   ```

7. **推送分支**
   ```bash
   git push origin feature/your-feature-name
   ```

8. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 清楚地描述你的更改
   - 关联相关的 Issue（如有）

### 代码风格

#### Go

- 遵循 [Effective Go](https://golang.org/doc/effective_go) 规范
- 使用 `gofmt` 格式化代码
- 添加适当的注释（中英文）
- 示例代码应该简洁易懂

```go
// 好的示例
func goodExample() {
    // 清晰的注释
    result := doSomething()
    if result != nil {
        handleResult(result)
    }
}

// 避免
func bad_example(){
result:=do_something()
if result!=nil{handle_result(result)}}
```

#### Python

- 遵循 PEP 8 规范
- 使用 4 个空格缩进
- 添加类型提示（Python 3.7+）
- 添加文档字符串

```python
# 好的示例
async def good_example(param: str) -> Dict:
    """
    清晰的文档字符串
    
    Args:
        param: 参数说明
    
    Returns:
        返回值说明
    """
    result = await do_something(param)
    return result

# 避免
async def bad_example(param):
    return await do_something(param)
```

### 添加新示例

当添加新示例时，请：

1. **创建清晰的目录结构**
   ```
   go-async/XX-your-example/
   ├── main.go
   └── README.md (可选)
   ```

2. **包含完整的注释**
   - 中文和英文双语注释
   - 解释代码的目的和原理

3. **提供可运行的代码**
   - 代码应该能够独立运行
   - 不依赖外部服务（或提供 mock）

4. **更新主 README**
   - 在项目结构中添加你的示例
   - 如果需要，添加到快速开始部分

### 报告 Bug

报告 Bug 时，请包含：

- 📝 详细的问题描述
- 💻 复现步骤
- 🖥️ 环境信息（OS, Go/Python 版本）
- 📷 截图或错误日志（如适用）

使用 Issue 模板：

```markdown
**描述问题 / Describe the bug**
清晰简洁地描述问题

**复现步骤 / Steps to reproduce**
1. 执行 '...'
2. 运行 '...'
3. 看到错误

**期望行为 / Expected behavior**
应该发生什么

**环境 / Environment**
- OS: [例如 macOS 13.0]
- Go 版本: [例如 1.21]
- Python 版本: [例如 3.11]

**额外信息 / Additional context**
其他相关信息
```

### 提出新功能

提出新功能时，请：

1. 检查是否已有类似的 Issue
2. 清楚地描述功能需求
3. 解释为什么需要这个功能
4. 提供示例代码或伪代码（如可能）

### 改进文档

文档改进包括：

- 修正拼写或语法错误
- 改进说明的清晰度
- 添加缺失的示例
- 改进代码注释

### 行为准则

- ✅ 尊重他人
- ✅ 欢迎新人
- ✅ 专注于问题，而非个人
- ✅ 接受建设性批评
- ❌ 使用不当语言
- ❌ 人身攻击
- ❌ 骚扰行为

---

## English

### How to Contribute

We welcome all forms of contributions:

- 🐛 Bug reports
- ✨ Feature requests
- 📝 Documentation improvements
- 💻 Code contributions
- 🌐 Translations

### Contribution Process

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-async-python-asyncio.git
   ```

3. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make changes**
   - Follow existing code style
   - Add necessary comments
   - Update relevant documentation

5. **Test your code**
   ```bash
   # Go
   go run main.go
   
   # Python
   python script.py
   ```

6. **Commit changes**
   ```bash
   git commit -m "Description of changes"
   ```

7. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

8. **Create Pull Request**

### Code Style

See Chinese section above for detailed style guidelines.

### Reporting Bugs

Include:
- Clear description
- Steps to reproduce
- Environment details
- Screenshots/logs if applicable

### Feature Requests

- Check for existing issues
- Describe the feature clearly
- Explain why it's needed
- Provide examples if possible

### Code of Conduct

- ✅ Be respectful
- ✅ Welcome newcomers
- ✅ Focus on issues, not people
- ✅ Accept constructive criticism
- ❌ No inappropriate language
- ❌ No personal attacks
- ❌ No harassment

---

## 联系方式 / Contact

有问题？欢迎：
- 创建 Issue
- 发起讨论

Questions? Feel free to:
- Open an Issue
- Start a Discussion

---

**感谢你的贡献！/ Thank you for your contribution!** ⭐
