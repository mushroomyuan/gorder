我来帮您总结一下启动 gRPC 的逻辑链条。让我先查看相关的代码文件来理解完整的流程。
Search files...
Ran tool
Ran tool
现在我来为您总结 gRPC 服务启动的完整逻辑链条：

## gRPC 服务启动逻辑链条总结

### 1. **配置初始化阶段**
```go
// internal/order/main.go 或 internal/stock/main.go
func init() {
    if err := config.NewViperConfig(); err != nil {
        logrus.Fatal(err)
    }
}
```
- 加载配置文件 `internal/common/config/global.yaml`
- 初始化 Viper 配置管理器

### 2. **应用上下文创建**
```go
// main 函数中
serviceName := viper.Sub("order").GetString("service-name")  // 获取服务名
ctx, cancel := context.WithCancel(context.Background())      // 创建上下文
defer cancel()

application, cleanup := service.NewApplication(ctx)         // 创建应用实例
defer cleanup()
```

### 3. **应用层初始化**
```go
// internal/order/service/application.go
func NewApplication(ctx context.Context) (app app.Application, cleanup func()) {
    // 创建 gRPC 客户端连接其他服务
    stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
    stockGRPC := grpc.NewStockGRPC(stockClient)
    
    // 初始化仓储、日志、指标等
    orderInmemRepo := adapters.NewMemoryOrderRepository()
    logger := logrus.NewEntry(logrus.StandardLogger())
    metricsClient := metrics.TodoMetrics{}
    
    // 组装应用层（CQRS 模式）
    return app.Application{
        Commands: app.Commands{
            CreateOrder: command.NewCreateOrderHandler(...),
            UpdateOrder: command.NewUpdateOrderHandler(...),
        },
        Queries: app.Queries{
            GetCustomerOrder: query.NewGetCustomerOrderHandler(...),
        },
    }
}
```

### 4. **gRPC 服务器启动**
```go
// main 函数中
go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
    svc := ports.NewGRPCServer(application)
    orderpb.RegisterOrderServiceServer(server, svc)
})
```

### 5. **端口适配器创建**
```go
// internal/order/ports/grpc.go
func NewGRPCServer(app app.Application) *GRPCServer {
    return &GRPCServer{app: app}
}

// 实现 protobuf 生成的服务接口
func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
    // 调用应用层处理业务逻辑
}
```

### 6. **服务注册与启动**
```go
// internal/common/server/grpc.go
func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
    // 获取配置的地址
    addr := viper.Sub(serviceName).GetString("grpc-addr")
    if addr == "" {
        addr = viper.GetString("fullback-grpc-addr")
    }
    RunGRPCServerOnAddr(addr, registerServer)
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
    // 创建 gRPC 服务器实例
    grpcServer := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            grpc_tags.UnaryServerInterceptor(...),
            grpc_logrus.UnaryServerInterceptor(...),
        ),
    )
    
    // 注册服务处理器
    registerServer(grpcServer)
    
    // 启动网络监听
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        logrus.Panic(err)
    }
    
    logrus.Infof("Starting gRPC Server,Listening:%s", addr)
    if err := grpcServer.Serve(listener); err != nil {
        logrus.Panic(err)
    }
}
```

### 7. **Protobuf 生成的注册函数**
```go
// internal/common/genproto/orderpb/order_grpc.pb.go
func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
    // 验证服务器实现
    if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
        t.testEmbeddedByValue()
    }
    // 注册服务描述符
    s.RegisterService(&OrderService_ServiceDesc, srv)
}
```

## 完整流程图

```
配置初始化 → 应用上下文创建 → 应用层初始化 → gRPC服务器启动
     ↓              ↓              ↓              ↓
  加载配置文件    创建Context     初始化仓储/客户端   创建gRPC Server
     ↓              ↓              ↓              ↓
  Viper配置      取消函数        CQRS架构组装     添加拦截器
                                                  ↓
                                            端口适配器创建
                                                  ↓
                                            实现业务接口
                                                  ↓
                                            服务注册与启动
                                                  ↓
                                            网络监听启动
```

## 关键特点

1. **分层架构**: 采用六边形架构，业务逻辑与技术实现分离
2. **CQRS 模式**: 命令查询职责分离
3. **依赖注入**: 通过构造函数注入依赖
4. **中间件支持**: 支持日志、标签等拦截器
5. **配置驱动**: 通过配置文件管理服务地址和参数
6. **优雅关闭**: 通过 context 和 defer 实现资源清理

这个设计使得 gRPC 服务具有良好的可测试性、可维护性和扩展性。