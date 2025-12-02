# Swagger 使用指南

## 访问 Swagger 文档

启动服务后，访问以下地址：

- **Swagger UI**: http://localhost:8888/swagger/index.html
- **API JSON**: http://localhost:8888/swagger/doc.json
- **API YAML**: http://localhost:8888/swagger/doc.yaml

## 如何添加新接口文档

### 1. 在 Handler 方法上添加注释

在 `internal/service/` 目录下的服务文件中，为每个 Handler 方法添加 Swagger 注释：

```go
// Login 管理员登录
// @Summary 管理员登录
// @Description 使用用户名和密码登录，返回 JWT Token
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} response.Response "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Router /auth/login [post]
func (s *AuthService) Login(c *gin.Context) {
    // ... 实现代码
}
```

### 2. 需要认证的接口添加 Security

```go
// @Security BearerAuth
```

### 3. 常用注释说明

| 注释 | 说明 | 示例 |
|------|------|------|
| @Summary | 接口摘要 | `@Summary 管理员登录` |
| @Description | 详细描述 | `@Description 使用用户名和密码登录` |
| @Tags | 接口分组 | `@Tags 认证管理` |
| @Accept | 接受的内容类型 | `@Accept json` |
| @Produce | 返回的内容类型 | `@Produce json` |
| @Param | 参数说明 | `@Param id path int true "用户ID"` |
| @Success | 成功响应 | `@Success 200 {object} response.Response` |
| @Failure | 失败响应 | `@Failure 400 {object} response.Response` |
| @Router | 路由路径 | `@Router /auth/login [post]` |
| @Security | 需要的认证 | `@Security BearerAuth` |

### 4. 参数类型说明

**Path 参数（路径参数）**:
```go
// @Param id path int true "用户ID"
```

**Query 参数（查询参数）**:
```go
// @Param page query int false "页码" default(1)
// @Param keyword query string false "搜索关键词"
```

**Body 参数（请求体）**:
```go
// @Param request body dto.LoginRequest true "登录信息"
```

**Header 参数**:
```go
// @Param Authorization header string true "Bearer Token"
```

### 5. 重新生成文档

修改代码后，需要重新生成 Swagger 文档：

```bash
# 在项目根目录执行
swag init

# 或使用 make 命令（如果配置了 Makefile）
make swagger
```

### 6. 查看生成的文档

生成的文档位于 `docs/` 目录：
- `docs/docs.go` - Go 代码
- `docs/swagger.json` - JSON 格式
- `docs/swagger.yaml` - YAML 格式

## 在 Swagger UI 中测试接口

### 1. 不需要认证的接口

直接点击 "Try it out"，填写参数，点击 "Execute" 即可。

### 2. 需要认证的接口

**步骤**：
1. 先调用 `/auth/login` 接口获取 Token
2. 点击页面右上角的 "Authorize" 按钮
3. 输入 `Bearer <你的token>`（注意 Bearer 后面有空格）
4. 点击 "Authorize"
5. 现在可以调用需要认证的接口了

### 3. 示例

假设登录接口返回的 token 是：
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

在 Authorize 弹窗中输入：
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 常见问题

### 1. 修改代码后文档没更新

重新运行 `swag init` 生成文档，然后重启服务。

### 2. 接口在 Swagger UI 中看不到

检查：
- 注释格式是否正确（注意空格和换行）
- 是否运行了 `swag init`
- 路由路径是否正确

### 3. 参数类型显示不正确

确保 DTO 结构体有正确的 JSON tag：

```go
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

### 4. 导入 docs 包报错

确保：
1. 已运行 `swag init` 生成了 docs 包
2. 在 `http.go` 中导入了 docs 包：
   ```go
   _ "github.com/ydcloud-dy/leaf-api/docs"
   ```

## 提示

- 建议为所有公开 API 添加 Swagger 注释
- 保持注释和代码同步更新
- 使用有意义的 Tags 对接口进行分组
- 详细描述参数和返回值，方便前端开发对接
