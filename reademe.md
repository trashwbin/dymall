3.1 技术选型与相关开发文档
开发文档-抖音电商
技术栈
- 微服务框架：Kitex - 字节跳动开源的高性能RPC框架
- 服务间通信：Protocol Buffers (protobuf) - 高效的数据序列化协议
- 日志管理：使用klog记录日志
系统规模预估
存储需求
- 用户数据：预计千万级用户，需要~500GB存储空间
- 商品数据：预计百万级SKU，需要~1TB存储空间
- 订单数据：日订单量预计10万级，需要~2TB存储空间/年
- 总体存储需求：~5TB（含冗余备份）
服务器需求
- 应用服务器：
  - 核心服务(用户、商品、订单)：各2-4台
  - 辅助服务(购物车、支付等)：各1-2台
  - 总计：约15-20台应用服务器
- 数据库服务器：主从架构，总计4-6台
- 缓存服务器：2-4台
3.2 架构设计
系统架构
采用微服务架构，主要包含以下核心服务：
1. 用户服务(user)
  - 用户账户管理
    - 用户注册：支持用户名、密码、邮箱等基本信息注册
    - 用户登录：提供账号密码登录接口
    - 用户登出：支持用户安全退出
  - 用户信息管理
    - 信息查询：获取用户基本信息，包括ID、用户名、邮箱等
    - 信息更新：支持更新用户个人信息
  - 技术实现
    - 基于Kitex框架实现RPC服务
    - 使用Protocol Buffers定义服务接口
    - 集成Consul实现服务注册与发现
  - 服务交互
    - 与认证服务：配合进行用户认证和授权
2. 认证服务(auth)
  - 统一身份认证
    - 用户登录认证
    - 会话管理
  - Token管理
    - JWT Token的签发与验证
    - Token刷新机制
    - Token黑名单管理
  - 权限管理
    - 基于RBAC的权限控制
    - 用户角色管理
    - 权限策略配置

3. 商品服务(product)
  - 商品管理
    - 商品CRUD：支持商品的创建、查询、更新和删除操作
    - 批量操作：支持批量获取商品信息
    - 商品搜索：支持基于关键词的商品搜索
    - 分类管理：支持商品分类的灵活配置
  - 数据模型
    - Product：商品基本信息，包含ID、名称、描述、图片、价格等
    - Category：商品分类信息
  - 技术实现
    - 使用Kitex框架实现RPC服务
    - 采用分层架构：handler层处理请求，service层实现业务逻辑，dal层处理数据访问
    - 支持配置化管理：区分开发、测试、生产环境
    - 集成认证中间件：实现接口访问控制
  - 服务交互
    - 与订单服务：提供商品信息和库存查询
    - 与购物车服务：提供商品详情
    - 与结算服务：确认商品价格信息

4. 购物车服务(cart)
  - 购物车管理
    - 添加商品(AddItem): 支持添加指定数量的商品到用户购物车
    - 查看购物车(GetCart): 获取用户购物车中的所有商品信息
    - 清空购物车(EmptyCart): 清除用户购物车中的所有商品
  - 数据模型
    - CartItem: 购物车商品项，包含商品ID和数量
    - Cart: 用户购物车，包含用户ID和商品项列表

5. 订单服务(order)
  - 订单管理
    - 创建订单：支持从购物车创建订单，包含商品信息、收货地址、用户币种等
    - 更新订单：支持修改订单收货地址
    - 取消订单：支持取消未支付订单，可选择是否级联取消关联的支付单
    - 查询订单：支持按订单ID查询详情和批量查询订单列表
  - 订单状态管理
    - 支持多种订单状态：待支付、已支付、已取消、已过期
    - 提供订单支付状态更新接口
    - 自动处理订单过期
  - 数据模型
    - Order：订单主体信息，包含订单ID、用户信息、商品列表、支付信息等
    - OrderItem：订单商品项，包含商品信息和实际成交价格
    - Address：订单收货地址信息
  - 服务交互
    - 与支付服务交互：处理订单支付状态
    - 与商品服务交互：确认商品信息和价格
    - 与用户服务交互：验证用户信息

6. 支付服务(payment)
  - 支付处理
    - 创建支付单：基于订单信息创建支付单，包含支付金额、支付方式等
    - 支付确认：处理支付网关的回调，确认支付状态
    - 支付查询：查询支付单状态和支付详情
    - 支付退款：支持订单退款处理
  - 支付状态管理
    - 支持多种支付状态：待支付、支付中、支付成功、支付失败、已退款
    - 提供支付状态变更通知机制
    - 自动处理支付超时
  - 数据模型
    - Payment：支付单信息，包含支付ID、订单ID、支付金额、支付方式等
    - PaymentLog：支付操作日志，记录支付状态变更历史
  - 服务交互
    - 与订单服务交互：同步支付状态，触发订单状态更新
    - 与用户服务交互：验证支付账户信息

7. 结算服务(checkout)
  - 结算单管理
    - 创建结算单：基于购物车商品创建结算单
    - 查询结算单：获取结算单详细信息，包含商品列表
    - 更新结算单：支持更新结算单状态和商品信息
    - 删除结算单：清理过期或无效的结算单
  - 数据存储
    - MySQL：持久化存储结算单基本信息
    - Redis：缓存结算单和商品信息，提升查询性能
  - 业务逻辑
    - 价格计算：计算商品总价、优惠金额和实付金额
    - 订单生成：调用订单服务，创建正式订单

8. 调度服务(scheduler)
  - 任务调度管理
    - 定时任务调度：支持基于cron表达式的定时任务配置和执行
    - 一次性任务：支持单次执行的任务调度
    - 任务优先级：支持设置任务的优先级，确保关键任务优先执行
  - 异步任务处理
    - 任务队列管理：维护不同类型的任务队列
    - 任务状态追踪：实时监控任务执行状态和进度
    - 失败重试机制：支持配置任务失败重试策略
  - 系统集成
    - 与订单服务集成：处理订单超时和状态更新
    - 与支付服务集成：处理支付超时和退款任务

## 3.3 项目代码介绍

### 目录结构

```
├── app/                # 应用服务目录
│   ├── auth/          # 认证服务
│   ├── cart/          # 购物车服务
│   ├── checkout/      # 结算服务
│   ├── order/         # 订单服务
│   ├── payment/       # 支付服务
│   ├── product/       # 商品服务
│   ├── scheduler/     # 调度服务
│   └── user/          # 用户服务
├── db/                # 数据库相关
│   └── sql/          # SQL脚本
├── idl/               # 接口定义
│   └── */            # 各服务的proto文件
└── docker-compose.yaml # 容器编排配置
```

### 服务结构

每个服务遵循相同的结构：

```
├── biz/              # 业务逻辑
│   ├── service/      # 服务实现
│   └── dal/         # 数据访问
├── conf/            # 配置文件
│   ├── dev/        # 开发环境配置
│   └── prod/       # 生产环境配置
├── middleware/     # 中间件
│   └── auth.go     # 认证中间件
├── sdk/            # 服务SDK
│   └── auth.go     # 认证服务客户端SDK
├── resources/      # 资源文件
│   └── rbac_model.conf # RBAC模型配置
├── handler.go      # 请求处理
├── main.go        # 启动入口
└── utils/         # 工具函数
```

### RPC通信结构

每个服务的RPC通信相关代码组织如下：

```
├── kitex_gen/           # Kitex自动生成的代码
│   └── */              # 各个服务的协议代码
│       ├── *.pb.go     # Protocol Buffers生成的Go代码
│       └── */*.go      # Kitex生成的服务接口代码
└── infra/
    └── rpc/           # RPC客户端实现
        └── client.go  # RPC客户端初始化和配置
```

kitex_gen目录包含了由Kitex工具自动生成的RPC通信代码：
- 基于Proto文件生成的Go结构体和接口定义
- 序列化/反序列化相关代码
- RPC服务端和客户端的基础实现

infra/rpc目录包含了各个服务的RPC客户端实现：
- 客户端连接池的初始化和管理
- 服务发现和负载均衡配置

## 3.4 项目总结与反思

### 3.4.1 目前仍存在的问题

1. 性能优化
   - 缺乏系统的性能监控和分析机制
   - 数据库查询和缓存策略需要优化
   - 服务间通信开销较大

2. 代码质量
   - 单元测试覆盖率不足
   - 部分业务逻辑耦合度较高
   - 错误处理机制不够统一

3. 服务治理
   - 缺乏完整的服务降级和熔断机制
   - 日志收集和分析系统不够完善
   - 监控告警体系需要加强

### 3.4.2 已识别出的优化项

1. 技术架构优化
   - 引入服务网格(Service Mesh)提升服务治理能力
   - 完善链路追踪系统，提升问题定位效率
   - 优化数据库分库分表方案

2. 性能提升
   - 优化缓存策略，引入多级缓存
   - 实现批量接口减少网络请求
   - 优化数据库索引和查询语句

3. 可用性提升
   - 完善服务容错机制
   - 实现更细粒度的限流策略
   - 优化服务发现和负载均衡

### 3.4.3 架构演进的可能性

1. 微服务架构优化
   - 服务粒度进一步优化，提取公共模块
   - 引入DDD领域驱动设计理念
   - 采用事件驱动架构处理异步场景

2. 技术栈升级
   - 评估新版本Kitex框架的特性
   - 引入云原生技术栈
   - 探索Serverless架构的应用

3. 部署架构演进
   - 完善容器化部署方案
   - 引入服务网格实现更细粒度的流量控制
   - 探索多集群部署方案

### 3.4.4 项目过程中的反思与总结

1. 技术选型
   - Kitex框架在微服务实现上表现良好
   - Protocol Buffers在服务通信上效率高
   - 分层架构便于代码维护和扩展

2. 开发流程
   - 需要加强需求分析和系统设计
   - 代码审查流程需要进一步规范
   - 文档更新需要及时跟进

3. 团队协作
   - 加强技术方案讨论和知识分享
   - 建立更完善的开发规范
   - 提升团队的技术能力

4. 经验总结
   - 合理的技术选型是项目成功的基础
   - 微服务架构需要完善的治理体系
   - 持续优化和改进是项目长期发展的关键
