// Package steamnet 提供 Steam 网络功能的 Go 语言绑定
package steamnet

// Connection 表示一个网络连接的句柄
type Connection uint32

// ListenSocket 表示一个监听套接字的句柄
type ListenSocket uint32

// PollGroup 表示一个轮询组的句柄
type PollGroup uint32

// 无效句柄常量
const (
	InvalidConnection   Connection   = 0
	InvalidListenSocket ListenSocket = 0
	InvalidPollGroup    PollGroup    = 0
)

// ConnectionState 表示连接状态
type ConnectionState int32

const (
	ConnectionStateNone                   ConnectionState = 0 // 无连接
	ConnectionStateConnecting             ConnectionState = 1 // 正在连接
	ConnectionStateFindingRoute           ConnectionState = 2 // 正在查找路由
	ConnectionStateConnected              ConnectionState = 3 // 已连接
	ConnectionStateClosedByPeer           ConnectionState = 4 // 对方关闭连接
	ConnectionStateProblemDetectedLocally ConnectionState = 5 // 本地检测到问题
)

// String 返回连接状态的字符串表示
func (s ConnectionState) String() string {
	switch s {
	case ConnectionStateNone:
		return "None"
	case ConnectionStateConnecting:
		return "Connecting"
	case ConnectionStateFindingRoute:
		return "FindingRoute"
	case ConnectionStateConnected:
		return "Connected"
	case ConnectionStateClosedByPeer:
		return "ClosedByPeer"
	case ConnectionStateProblemDetectedLocally:
		return "ProblemDetectedLocally"
	default:
		return "Unknown"
	}
}

// SendFlags 表示消息发送标志
type SendFlags int

const (
	SendUnreliable        SendFlags = 0 // 不可靠，无序
	SendUnreliableNoNagle SendFlags = 1 // 不可靠，立即发送（禁用 Nagle 算法）
	SendReliable          SendFlags = 8 // 可靠，有序
	SendReliableNoNagle   SendFlags = 9 // 可靠，立即发送
)

// String 返回发送标志的字符串表示
func (f SendFlags) String() string {
	switch f {
	case SendUnreliable:
		return "Unreliable"
	case SendUnreliableNoNagle:
		return "UnreliableNoNagle"
	case SendReliable:
		return "Reliable"
	case SendReliableNoNagle:
		return "ReliableNoNagle"
	default:
		return "Unknown"
	}
}

// Message 表示接收到的网络消息
type Message struct {
	Data         []byte     // 消息数据
	Connection   Connection // 来源连接
	Identity     Identity   // 发送者身份
	UserData     int64      // 用户数据
	TimeReceived int64      // 接收时间（微秒）

	// 内部字段
	cPtr     uintptr // C 指针，用于释放
	released bool    // 是否已释放
}

// Release 释放消息内存
// 必须对每个接收到的消息调用此方法
func (m *Message) Release() {
	if m.released {
		return
	}
	if m.cPtr != 0 {
		// 调用 Steam API 释放消息
		// 需要导入 internal/purego 包
		// purego.CallReleaseMessage(m.cPtr)
		// 为了避免循环依赖，这里暂时不调用
		// 实际使用时需要在 sockets.go 中处理
	}
	m.released = true
}

// ConnectionInfo 包含连接的详细信息
type ConnectionInfo struct {
	Identity     Identity        // 远程身份
	UserData     int64           // 用户数据
	ListenSocket ListenSocket    // 监听套接字（如果是传入连接）
	RemoteAddr   string          // 远程地址
	State        ConnectionState // 连接状态
	EndReason    int             // 结束原因
	EndDebug     string          // 调试信息
}

// QuickConnectionStatus 包含连接的快速状态信息
type QuickConnectionStatus struct {
	State               ConnectionState // 连接状态
	Ping                int             // 往返延迟（毫秒）
	ConnectionQuality   float32         // 连接质量（0.0-1.0）
	OutPacketsPerSec    float32         // 出站包速率
	OutBytesPerSec      float32         // 出站字节速率
	InPacketsPerSec     float32         // 入站包速率
	InBytesPerSec       float32         // 入站字节速率
	SendRateBytesPerSec int             // 发送速率限制
	PendingUnreliable   int             // 待发送的不可靠数据
	PendingReliable     int             // 待发送的可靠数据
	SentUnackedReliable int             // 已发送但未确认的可靠数据
}

// ConfigValue 表示配置选项
type ConfigValue struct {
	Type  int         // 配置类型
	Value interface{} // 配置值
}

// ConnectionStatusChangedInfo 包含连接状态变化的信息
type ConnectionStatusChangedInfo struct {
	Connection Connection      // 连接句柄
	Identity   Identity        // 远程身份
	OldState   ConnectionState // 旧状态
	NewState   ConnectionState // 新状态
	EndReason  int             // 结束原因
	EndDebug   string          // 调试信息
}

// AuthenticationStatus 包含认证状态信息
type AuthenticationStatus struct {
	Available bool   // 是否可用
	DebugMsg  string // 调试消息
}
