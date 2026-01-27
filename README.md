# SteamKit-Go

Steamworks SDK 的 Go 语言绑定，专注于网络功能。

## 项目状态

🚧 **开发中** - 当前处于早期开发阶段

### 已完成

✅ **阶段 1：基础设施**
- 项目目录结构
- purego 绑定层（跨平台库加载）
- 主包初始化（Init/Shutdown/RunCallbacks）
- 基础示例程序

✅ **阶段 2：核心类型定义**
- steamnet 包的类型定义
- 句柄类型（Connection, ListenSocket, PollGroup）
- 枚举类型（ConnectionState, SendFlags）
- 结构体类型（Message, ConnectionInfo, QuickConnectionStatus）
- Identity 类型和辅助函数
- 错误定义和错误检查函数
- 完整的单元测试（100% 通过）

✅ **阶段 3：连接管理接口**
- ISteamNetworkingSockets 接口定义
- CreateListenSocketP2P - 创建 P2P 监听套接字
- ConnectP2P - 连接到远程对等方
- AcceptConnection - 接受传入连接
- CloseConnection - 关闭连接
- CloseListenSocket - 关闭监听套接字
- GetConnectionInfo - 获取连接信息
- Mock 实现和单元测试
- 连接测试示例程序

✅ **阶段 4：消息收发接口**
- SendMessageToConnection - 发送消息到连接
- FlushMessagesOnConnection - 刷新连接上的消息
- ReceiveMessagesOnConnection - 接收连接上的消息
- ReceiveMessagesOnListenSocket - 接收监听套接字上的消息
- 消息内存管理（ReleaseMessage）
- 消息解析（parseMessage）
- Mock 实现和单元测试（29 个测试全部通过）

✅ **阶段 5：回调处理**
- ConnectionStatusChangedCallback 回调类型
- SetConnectionStatusChangedCallback - 设置全局回调
- SetConnectionCallback - 设置特定连接回调
- ClearConnectionCallback - 清除连接回调
- DispatchConnectionStatusChanged - 回调分发机制
- 回调管理器（支持全局和特定连接回调）
- 完整的单元测试（41 个测试全部通过）
- 回调示例程序

✅ **阶段 6：连接状态轮询**
- GetConnectionRealTimeStatus - 获取连接实时状态
- QuickConnectionStatus 结构体（包含 ping、连接质量、流量统计等）
- Mock 实现和单元测试（42 个测试全部通过）

### 待完成

⏳ **阶段 7：配置选项**
- ConfigValue 处理
- 连接配置选项

⏳ **阶段 8：高级功能**
- Poll Groups
- 其他高级特性

⏳ **阶段 9：性能优化**
- 内存池
- 批量操作优化

⏳ **阶段 10：文档和示例**
- 完整的 API 文档
- 更多示例程序
- 最佳实践指南

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- Steam 客户端（已登录）
- Steamworks SDK v161 或更高版本

### 安装

```bash
go get github.com/guowei-gong/steamkit-go
```

### 基础示例

```go
package main

import (
    "fmt"
    "log"

    "github.com/guowei-gong/steamkit-go"
)

func main() {
    // 初始化 Steam API
    if err := steamkit.Init(); err != nil {
        log.Fatal(err)
    }
    defer steamkit.Shutdown()

    // 获取当前用户 SteamID
    steamID := steamkit.GetSteamID()
    fmt.Printf("SteamID: %d\n", steamID)
}
```

### 运行示例

```bash
# 确保 steam_api64.dll (Windows) 或 libsteam_api.so (Linux) 在当前目录或系统路径中
go run examples/basic_test/main.go
```

## 架构

```
steamkit-go/
├── steamkit.go              # 主包：初始化、RunCallbacks
├── steamnet/                # 网络包（待实现）
│   ├── sockets.go           # ISteamNetworkingSockets 接口
│   ├── types.go             # 类型定义
│   ├── identity.go          # SteamNetworkingIdentity
│   ├── callbacks.go         # 回调处理
│   └── errors.go            # 错误定义
├── internal/purego/         # purego 绑定层
│   ├── loader.go            # 核心加载逻辑
│   ├── loader_windows.go    # Windows 特定
│   ├── loader_linux.go      # Linux 特定
│   └── loader_darwin.go     # macOS 特定
└── examples/                # 示例程序
    └── basic_test/          # 基础测试
```

## 设计文档

详细的设计文档和规范请参考：
- [设计文档](openspec/changes/research-steamnetworkingsockets-binding/design.md)
- [接口规范](specs/steamnet/spec.md)
- [实施路线图](openspec/changes/research-steamnetworkingsockets-binding/design.md#实施路线图)

## 技术栈

- **FFI**: [purego](https://github.com/ebitengine/purego) - 纯 Go 实现的 FFI，无需 cgo
- **SDK**: Steamworks SDK v161

## 许可证

待定

## 贡献

欢迎贡献！请先阅读设计文档了解项目架构。

## 致谢

- [go-steamworks](https://github.com/hajimehoshi/go-steamworks) - 参考项目
- [purego](https://github.com/ebitengine/purego) - FFI 库
