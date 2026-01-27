package steamnet

import (
	"fmt"
	"unsafe"

	"github.com/guowei-gong/steamkit-go/internal/purego"
)

// ISteamNetworkingSockets 定义网络套接字接口
type ISteamNetworkingSockets interface {
	// P2P 连接管理
	CreateListenSocketP2P(virtualPort int, options []ConfigValue) (ListenSocket, error)
	ConnectP2P(identity Identity, virtualPort int, options []ConfigValue) (Connection, error)
	AcceptConnection(conn Connection) error
	CloseConnection(conn Connection, reason int, debug string, linger bool) error
	CloseListenSocket(socket ListenSocket) error

	// 消息收发
	SendMessageToConnection(conn Connection, data []byte, flags SendFlags) error
	FlushMessagesOnConnection(conn Connection) error
	ReceiveMessagesOnConnection(conn Connection, maxMessages int) ([]*Message, error)
	ReceiveMessagesOnListenSocket(socket ListenSocket, maxMessages int) ([]*Message, error)

	// 连接信息
	GetConnectionInfo(conn Connection) (*ConnectionInfo, error)
	GetConnectionRealTimeStatus(conn Connection) (*QuickConnectionStatus, error)
}

// steamNetworkingSockets 是 ISteamNetworkingSockets 的实现
type steamNetworkingSockets struct {
	handle uintptr
}

// GetSockets 返回 ISteamNetworkingSockets 接口实例
func GetSockets() ISteamNetworkingSockets {
	handle := purego.CallGetSteamNetworkingSockets()
	if handle == 0 {
		return nil
	}
	return &steamNetworkingSockets{
		handle: handle,
	}
}

// CreateListenSocketP2P 创建一个 P2P 监听套接字
func (s *steamNetworkingSockets) CreateListenSocketP2P(virtualPort int, options []ConfigValue) (ListenSocket, error) {
	// TODO: 处理 options 参数
	// 目前传递 0 和 nil
	handle := purego.CallCreateListenSocketP2P(s.handle, int32(virtualPort), 0, 0)
	if handle == 0 {
		return InvalidListenSocket, ErrInvalidSocket
	}
	return ListenSocket(handle), nil
}

// ConnectP2P 连接到远程 P2P 对等方
func (s *steamNetworkingSockets) ConnectP2P(identity Identity, virtualPort int, options []ConfigValue) (Connection, error) {
	if !identity.IsValid() {
		return InvalidConnection, ErrInvalidIdentity
	}

	// TODO: 将 Identity 转换为 SteamNetworkingIdentity 结构体
	// 目前仅支持 SteamID
	if identity.Type() != IdentityTypeSteamID {
		return InvalidConnection, fmt.Errorf("only SteamID identity is supported currently")
	}

	// 创建 SteamNetworkingIdentity 结构体
	// 这是一个简化版本，实际需要完整的结构体定义
	var identityStruct [136]byte // SteamNetworkingIdentity 的大小
	steamID := identity.GetSteamID()
	*(*uint64)(unsafe.Pointer(&identityStruct[0])) = steamID
	*(*int32)(unsafe.Pointer(&identityStruct[128])) = 16 // k_ESteamNetworkingIdentityType_SteamID

	// TODO: 处理 options 参数
	handle := purego.CallConnectP2P(s.handle, uintptr(unsafe.Pointer(&identityStruct[0])), int32(virtualPort), 0, 0)
	if handle == 0 {
		return InvalidConnection, ErrConnectionFailed
	}
	return Connection(handle), nil
}

// AcceptConnection 接受传入连接
func (s *steamNetworkingSockets) AcceptConnection(conn Connection) error {
	if conn == InvalidConnection {
		return ErrInvalidConnection
	}

	result := purego.CallAcceptConnection(s.handle, uint32(conn))
	if result != 0 { // k_EResultOK = 1, 但这里返回 0 表示成功
		return fmt.Errorf("failed to accept connection: result=%d", result)
	}
	return nil
}

// CloseConnection 关闭连接
func (s *steamNetworkingSockets) CloseConnection(conn Connection, reason int, debug string, linger bool) error {
	if conn == InvalidConnection {
		return ErrInvalidConnection
	}

	// 将 debug 字符串转换为 C 字符串
	var debugPtr uintptr
	if debug != "" {
		debugBytes := append([]byte(debug), 0)
		debugPtr = uintptr(unsafe.Pointer(&debugBytes[0]))
	}

	success := purego.CallCloseConnection(s.handle, uint32(conn), int32(reason), debugPtr, linger)
	if !success {
		return fmt.Errorf("failed to close connection")
	}
	return nil
}

// CloseListenSocket 关闭监听套接字
func (s *steamNetworkingSockets) CloseListenSocket(socket ListenSocket) error {
	if socket == InvalidListenSocket {
		return ErrInvalidSocket
	}

	success := purego.CallCloseListenSocket(s.handle, uint32(socket))
	if !success {
		return fmt.Errorf("failed to close listen socket")
	}
	return nil
}

// GetConnectionInfo 获取连接信息
func (s *steamNetworkingSockets) GetConnectionInfo(conn Connection) (*ConnectionInfo, error) {
	if conn == InvalidConnection {
		return nil, ErrInvalidConnection
	}

	// TODO: 定义完整的 SteamNetConnectionInfo_t 结构体
	// 目前返回一个简化版本
	var infoStruct [696]byte // SteamNetConnectionInfo_t 的大小（近似）

	success := purego.CallGetConnectionInfo(s.handle, uint32(conn), uintptr(unsafe.Pointer(&infoStruct[0])))
	if !success {
		return nil, fmt.Errorf("failed to get connection info")
	}

	// TODO: 解析结构体并填充 ConnectionInfo
	// 目前返回一个空的 ConnectionInfo
	info := &ConnectionInfo{
		Identity: NewInvalidIdentity(),
		State:    ConnectionState(*(*int32)(unsafe.Pointer(&infoStruct[0]))),
	}

	return info, nil
}

// SendMessageToConnection 发送消息到连接
func (s *steamNetworkingSockets) SendMessageToConnection(conn Connection, data []byte, flags SendFlags) error {
	if conn == InvalidConnection {
		return ErrInvalidConnection
	}

	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}

	// 调用 Steam API
	result := purego.CallSendMessageToConnection(
		s.handle,
		uint32(conn),
		uintptr(unsafe.Pointer(&data[0])),
		uint32(len(data)),
		int32(flags),
		0, // outMessageNumber (可选)
	)

	// k_EResultOK = 1
	if result != 1 {
		return fmt.Errorf("failed to send message: result=%d", result)
	}

	return nil
}

// FlushMessagesOnConnection 刷新连接上的消息
func (s *steamNetworkingSockets) FlushMessagesOnConnection(conn Connection) error {
	if conn == InvalidConnection {
		return ErrInvalidConnection
	}

	result := purego.CallFlushMessagesOnConnection(s.handle, uint32(conn))

	// k_EResultOK = 1
	if result != 1 {
		return fmt.Errorf("failed to flush messages: result=%d", result)
	}

	return nil
}

// ReceiveMessagesOnConnection 接收连接上的消息
func (s *steamNetworkingSockets) ReceiveMessagesOnConnection(conn Connection, maxMessages int) ([]*Message, error) {
	if conn == InvalidConnection {
		return nil, ErrInvalidConnection
	}

	if maxMessages <= 0 {
		return nil, fmt.Errorf("maxMessages must be positive")
	}

	// 创建消息指针数组
	messagePtrs := make([]uintptr, maxMessages)

	// 调用 Steam API
	numMessages := purego.CallReceiveMessagesOnConnection(
		s.handle,
		uint32(conn),
		uintptr(unsafe.Pointer(&messagePtrs[0])),
		int32(maxMessages),
	)

	if numMessages < 0 {
		return nil, fmt.Errorf("failed to receive messages: result=%d", numMessages)
	}

	if numMessages == 0 {
		return []*Message{}, nil
	}

	// 解析消息
	messages := make([]*Message, numMessages)
	for i := int32(0); i < numMessages; i++ {
		msgPtr := messagePtrs[i]
		if msgPtr == 0 {
			continue
		}

		// 解析 SteamNetworkingMessage_t 结构体
		// 简化版本：只读取数据指针和大小
		msg := parseMessage(msgPtr, conn)
		messages[i] = msg
	}

	return messages, nil
}

// ReceiveMessagesOnListenSocket 接收监听套接字上的消息
func (s *steamNetworkingSockets) ReceiveMessagesOnListenSocket(socket ListenSocket, maxMessages int) ([]*Message, error) {
	if socket == InvalidListenSocket {
		return nil, ErrInvalidSocket
	}

	if maxMessages <= 0 {
		return nil, fmt.Errorf("maxMessages must be positive")
	}

	// 创建消息指针数组
	messagePtrs := make([]uintptr, maxMessages)

	// 调用 Steam API
	numMessages := purego.CallReceiveMessagesOnListenSocket(
		s.handle,
		uint32(socket),
		uintptr(unsafe.Pointer(&messagePtrs[0])),
		int32(maxMessages),
	)

	if numMessages < 0 {
		return nil, fmt.Errorf("failed to receive messages: result=%d", numMessages)
	}

	if numMessages == 0 {
		return []*Message{}, nil
	}

	// 解析消息
	messages := make([]*Message, numMessages)
	for i := int32(0); i < numMessages; i++ {
		msgPtr := messagePtrs[i]
		if msgPtr == 0 {
			continue
		}

		// 解析 SteamNetworkingMessage_t 结构体
		msg := parseMessage(msgPtr, InvalidConnection)
		messages[i] = msg
	}

	return messages, nil
}

// parseMessage 解析 SteamNetworkingMessage_t 结构体
func parseMessage(msgPtr uintptr, conn Connection) *Message {
	// SteamNetworkingMessage_t 结构体布局（简化）：
	// offset 0: void* m_pData
	// offset 8: int m_cbSize
	// offset 16: HSteamNetConnection m_conn
	// offset 24: SteamNetworkingIdentity m_identityPeer
	// offset 160: int64 m_nConnUserData
	// offset 168: SteamNetworkingMicroseconds m_usecTimeReceived

	dataPtr := *(*uintptr)(unsafe.Pointer(msgPtr))
	dataSize := *(*int32)(unsafe.Pointer(msgPtr + 8))
	msgConn := *(*uint32)(unsafe.Pointer(msgPtr + 16))
	timeReceived := *(*int64)(unsafe.Pointer(msgPtr + 168))

	// 复制数据到 Go 切片
	data := make([]byte, dataSize)
	if dataSize > 0 && dataPtr != 0 {
		src := unsafe.Slice((*byte)(unsafe.Pointer(dataPtr)), dataSize)
		copy(data, src)
	}

	// 如果连接无效，使用消息中的连接
	if conn == InvalidConnection {
		conn = Connection(msgConn)
	}

	return &Message{
		Data:         data,
		Connection:   conn,
		Identity:     NewInvalidIdentity(), // TODO: 解析 m_identityPeer
		UserData:     0,                     // TODO: 解析 m_nConnUserData
		TimeReceived: timeReceived,
		cPtr:         msgPtr,
		released:     false,
	}
}

// ReleaseMessage 释放消息内存（辅助函数）
func ReleaseMessage(msg *Message) {
	if msg == nil || msg.released {
		return
	}
	if msg.cPtr != 0 {
		purego.CallReleaseMessage(msg.cPtr)
	}
	msg.released = true
}

// GetConnectionRealTimeStatus 获取连接的实时状态
func (s *steamNetworkingSockets) GetConnectionRealTimeStatus(conn Connection) (*QuickConnectionStatus, error) {
	if conn == InvalidConnection {
		return nil, ErrInvalidConnection
	}

	// SteamNetConnectionRealTimeStatus_t 结构体（简化版本）
	// 实际结构体大小约为 312 字节
	var statusStruct [312]byte

	result := purego.CallGetConnectionRealTimeStatus(
		s.handle,
		uint32(conn),
		uintptr(unsafe.Pointer(&statusStruct[0])),
		0, // numLanes
		0, // lanes
	)

	// k_EResultOK = 1
	if result != 1 {
		return nil, fmt.Errorf("failed to get connection real-time status: result=%d", result)
	}

	// 解析 SteamNetConnectionRealTimeStatus_t 结构体
	// offset 0: ESteamNetworkingConnectionState m_eState (int32)
	// offset 4: int m_nPing (int32)
	// offset 8: float m_flConnectionQualityLocal (float32)
	// offset 12: float m_flConnectionQualityRemote (float32)
	// offset 16: float m_flOutPacketsPerSec (float32)
	// offset 20: float m_flOutBytesPerSec (float32)
	// offset 24: float m_flInPacketsPerSec (float32)
	// offset 28: float m_flInBytesPerSec (float32)
	// offset 32: int m_nSendRateBytesPerSecond (int32)
	// offset 36: int m_cbPendingUnreliable (int32)
	// offset 40: int m_cbPendingReliable (int32)
	// offset 44: int m_cbSentUnackedReliable (int32)

	status := &QuickConnectionStatus{
		State:               ConnectionState(*(*int32)(unsafe.Pointer(&statusStruct[0]))),
		Ping:                int(*(*int32)(unsafe.Pointer(&statusStruct[4]))),
		ConnectionQuality:   *(*float32)(unsafe.Pointer(&statusStruct[8])),
		OutPacketsPerSec:    *(*float32)(unsafe.Pointer(&statusStruct[16])),
		OutBytesPerSec:      *(*float32)(unsafe.Pointer(&statusStruct[20])),
		InPacketsPerSec:     *(*float32)(unsafe.Pointer(&statusStruct[24])),
		InBytesPerSec:       *(*float32)(unsafe.Pointer(&statusStruct[28])),
		SendRateBytesPerSec: int(*(*int32)(unsafe.Pointer(&statusStruct[32]))),
		PendingUnreliable:   int(*(*int32)(unsafe.Pointer(&statusStruct[36]))),
		PendingReliable:     int(*(*int32)(unsafe.Pointer(&statusStruct[40]))),
		SentUnackedReliable: int(*(*int32)(unsafe.Pointer(&statusStruct[44]))),
	}

	return status, nil
}
