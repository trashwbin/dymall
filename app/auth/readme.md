# 认证服务(Auth Service)

## 服务介绍

认证服务是整个电商系统的核心基础服务之一，负责处理所有与用户身份认证、授权和权限管理相关的功能。本服务基于Kitex框架开发，提供高性能、可靠的认证机制。

### 核心功能

1. 统一身份认证
   - 用户登录认证
   - Token的生成与验证
   - 会话管理

2. 权限管理
   - 基于RBAC的权限控制
   - 用户角色管理
   - 权限策略配置

3. Token服务
   - JWT Token的签发与验证
   - Token刷新机制
   - Token黑名单管理

## 技术实现

- 使用[Kitex](https://github.com/cloudwego/kitex/)框架作为RPC通信基础
- 采用JWT(JSON Web Token)实现无状态的用户认证
- 实现基于RBAC(Role-Based Access Control)的权限管理
- 集成单元测试框架，确保代码质量
- 提供SDK方便其他服务集成认证功能

## 目录结构

```
├── biz/                # 业务逻辑层
│   ├── service/        # 核心服务实现
│   │   ├── auth.go     # 认证相关业务逻辑
│   │   └── token.go    # Token处理逻辑
│   └── dal/           # 数据访问层
│       └── user.go     # 用户数据操作
├── conf/              # 配置文件目录
│   ├── dev/          # 开发环境配置
│   └── prod/         # 生产环境配置
├── middleware/       # 中间件
│   └── auth.go       # 认证中间件
├── sdk/              # 服务SDK
│   └── auth.go       # 认证服务客户端SDK
├── handler.go        # 请求处理器
├── main.go           # 服务入口
└── resources/        # 资源文件
    └── rbac_model.conf # RBAC模型配置
```

## 关键组件说明

- **handler.go**: 处理RPC请求，进行参数验证和响应封装
- **biz/service**: 实现核心的认证和授权逻辑
- **middleware**: 提供认证中间件，供其他服务使用
- **sdk**: 提供客户端SDK，简化服务间认证集成

## 如何运行

1. 安装依赖
```shell
go mod tidy
```

2. 编译运行
```shell
sh build.sh
sh output/bootstrap.sh
```

## 服务集成

其他服务如需集成认证功能，可以通过以下步骤：

1. 引入认证服务SDK
2. 配置认证中间件
3. 在需要认证的接口上使用中间件

示例代码：
```go
// 初始化认证SDK
authSDK := auth.NewAuthSDK()

// 使用认证中间件
endpoint = middleware.AuthMiddleware(true, false)(endpoint)
```
