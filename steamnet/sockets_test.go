package steamnet

import (
	"testing"
)

// MockSockets 是 ISteamNetworkingSockets 的 mock 实现
type MockSockets struct {
	CreateListenSocketP2PFunc func(int, []ConfigValue) (ListenSocket, error)
	ConnectP2PFunc            func(Identity, int, []ConfigValue) (Connection, error)
	AcceptConnectionFunc      func(Connection) error
	CloseConnectionFunc       func(Connection, int, string, bool) error
	CloseListenSocketFunc     func(ListenSocket) error
	GetConnectionInfoFunc     func(Connection) (*ConnectionInfo, error)
	GetConnectionRealTimeStatusFunc func(Connection) (*QuickConnectionStatus, error)
	SendMessageToConnectionFunc func(Connection, []byte, SendFlags) error
	FlushMessagesOnConnectionFunc func(Connection) error
	ReceiveMessagesOnConnectionFunc func(Connection, int) ([]*Message, error)
	ReceiveMessagesOnListenSocketFunc func(ListenSocket, int) ([]*Message, error)
}

func (m *MockSockets) CreateListenSocketP2P(virtualPort int, options []ConfigValue) (ListenSocket, error) {
	if m.CreateListenSocketP2PFunc != nil {
		return m.CreateListenSocketP2PFunc(virtualPort, options)
	}
	return ListenSocket(1), nil
}

func (m *MockSockets) ConnectP2P(identity Identity, virtualPort int, options []ConfigValue) (Connection, error) {
	if m.ConnectP2PFunc != nil {
		return m.ConnectP2PFunc(identity, virtualPort, options)
	}
	return Connection(1), nil
}

func (m *MockSockets) AcceptConnection(conn Connection) error {
	if m.AcceptConnectionFunc != nil {
		return m.AcceptConnectionFunc(conn)
	}
	return nil
}

func (m *MockSockets) CloseConnection(conn Connection, reason int, debug string, linger bool) error {
	if m.CloseConnectionFunc != nil {
		return m.CloseConnectionFunc(conn, reason, debug, linger)
	}
	return nil
}

func (m *MockSockets) CloseListenSocket(socket ListenSocket) error {
	if m.CloseListenSocketFunc != nil {
		return m.CloseListenSocketFunc(socket)
	}
	return nil
}

func (m *MockSockets) GetConnectionInfo(conn Connection) (*ConnectionInfo, error) {
	if m.GetConnectionInfoFunc != nil {
		return m.GetConnectionInfoFunc(conn)
	}
	return &ConnectionInfo{}, nil
}

func (m *MockSockets) SendMessageToConnection(conn Connection, data []byte, flags SendFlags) error {
	if m.SendMessageToConnectionFunc != nil {
		return m.SendMessageToConnectionFunc(conn, data, flags)
	}
	return nil
}

func (m *MockSockets) FlushMessagesOnConnection(conn Connection) error {
	if m.FlushMessagesOnConnectionFunc != nil {
		return m.FlushMessagesOnConnectionFunc(conn)
	}
	return nil
}

func (m *MockSockets) ReceiveMessagesOnConnection(conn Connection, maxMessages int) ([]*Message, error) {
	if m.ReceiveMessagesOnConnectionFunc != nil {
		return m.ReceiveMessagesOnConnectionFunc(conn, maxMessages)
	}
	return []*Message{}, nil
}

func (m *MockSockets) ReceiveMessagesOnListenSocket(socket ListenSocket, maxMessages int) ([]*Message, error) {
	if m.ReceiveMessagesOnListenSocketFunc != nil {
		return m.ReceiveMessagesOnListenSocketFunc(socket, maxMessages)
	}
	return []*Message{}, nil
}

func (m *MockSockets) GetConnectionRealTimeStatus(conn Connection) (*QuickConnectionStatus, error) {
	if m.GetConnectionRealTimeStatusFunc != nil {
		return m.GetConnectionRealTimeStatusFunc(conn)
	}
	return &QuickConnectionStatus{
		State: ConnectionStateConnected,
		Ping:  50,
	}, nil
}

// 测试 Mock 实现
func TestMockSockets(t *testing.T) {
	mock := &MockSockets{}

	// 测试 CreateListenSocketP2P
	socket, err := mock.CreateListenSocketP2P(0, nil)
	if err != nil {
		t.Errorf("CreateListenSocketP2P() error = %v", err)
	}
	if socket == InvalidListenSocket {
		t.Error("CreateListenSocketP2P() returned invalid socket")
	}

	// 测试 ConnectP2P
	identity := NewIdentityFromSteamID(76561198000000000)
	conn, err := mock.ConnectP2P(identity, 0, nil)
	if err != nil {
		t.Errorf("ConnectP2P() error = %v", err)
	}
	if conn == InvalidConnection {
		t.Error("ConnectP2P() returned invalid connection")
	}

	// 测试 AcceptConnection
	if err := mock.AcceptConnection(conn); err != nil {
		t.Errorf("AcceptConnection() error = %v", err)
	}

	// 测试 GetConnectionInfo
	info, err := mock.GetConnectionInfo(conn)
	if err != nil {
		t.Errorf("GetConnectionInfo() error = %v", err)
	}
	if info == nil {
		t.Error("GetConnectionInfo() returned nil")
	}

	// 测试 CloseConnection
	if err := mock.CloseConnection(conn, 0, "test", false); err != nil {
		t.Errorf("CloseConnection() error = %v", err)
	}

	// 测试 CloseListenSocket
	if err := mock.CloseListenSocket(socket); err != nil {
		t.Errorf("CloseListenSocket() error = %v", err)
	}
}

// 测试错误处理
func TestSocketsErrorHandling(t *testing.T) {
	mock := &MockSockets{
		ConnectP2PFunc: func(identity Identity, virtualPort int, options []ConfigValue) (Connection, error) {
			if !identity.IsValid() {
				return InvalidConnection, ErrInvalidIdentity
			}
			return Connection(1), nil
		},
		AcceptConnectionFunc: func(conn Connection) error {
			if conn == InvalidConnection {
				return ErrInvalidConnection
			}
			return nil
		},
		CloseConnectionFunc: func(conn Connection, reason int, debug string, linger bool) error {
			if conn == InvalidConnection {
				return ErrInvalidConnection
			}
			return nil
		},
	}

	// 测试无效身份
	_, err := mock.ConnectP2P(NewInvalidIdentity(), 0, nil)
	if !IsInvalidIdentity(err) {
		t.Errorf("ConnectP2P() with invalid identity should return ErrInvalidIdentity, got %v", err)
	}

	// 测试无效连接
	err = mock.AcceptConnection(InvalidConnection)
	if !IsInvalidConnection(err) {
		t.Errorf("AcceptConnection() with invalid connection should return ErrInvalidConnection, got %v", err)
	}

	err = mock.CloseConnection(InvalidConnection, 0, "", false)
	if !IsInvalidConnection(err) {
		t.Errorf("CloseConnection() with invalid connection should return ErrInvalidConnection, got %v", err)
	}
}

// 测试消息收发
func TestMessageSendReceive(t *testing.T) {
	testData := []byte("Hello, Steam!")
	testConn := Connection(1)

	mock := &MockSockets{
		SendMessageToConnectionFunc: func(conn Connection, data []byte, flags SendFlags) error {
			if conn == InvalidConnection {
				return ErrInvalidConnection
			}
			if len(data) == 0 {
				return ErrSendFailed
			}
			return nil
		},
		ReceiveMessagesOnConnectionFunc: func(conn Connection, maxMessages int) ([]*Message, error) {
			if conn == InvalidConnection {
				return nil, ErrInvalidConnection
			}
			if maxMessages <= 0 {
				return nil, ErrReceiveFailed
			}
			// 返回测试消息
			return []*Message{
				{
					Data:       testData,
					Connection: conn,
					Identity:   NewInvalidIdentity(),
				},
			}, nil
		},
	}

	// 测试发送消息
	err := mock.SendMessageToConnection(testConn, testData, SendReliable)
	if err != nil {
		t.Errorf("SendMessageToConnection() error = %v", err)
	}

	// 测试发送空消息
	err = mock.SendMessageToConnection(testConn, []byte{}, SendReliable)
	if !IsSendFailed(err) {
		t.Errorf("SendMessageToConnection() with empty data should return ErrSendFailed, got %v", err)
	}

	// 测试接收消息
	messages, err := mock.ReceiveMessagesOnConnection(testConn, 32)
	if err != nil {
		t.Errorf("ReceiveMessagesOnConnection() error = %v", err)
	}
	if len(messages) != 1 {
		t.Errorf("ReceiveMessagesOnConnection() returned %d messages, want 1", len(messages))
	}
	if len(messages) > 0 && string(messages[0].Data) != string(testData) {
		t.Errorf("ReceiveMessagesOnConnection() data = %s, want %s", string(messages[0].Data), string(testData))
	}

	// 测试接收消息（无效参数）
	_, err = mock.ReceiveMessagesOnConnection(testConn, 0)
	if !IsReceiveFailed(err) {
		t.Errorf("ReceiveMessagesOnConnection() with maxMessages=0 should return ErrReceiveFailed, got %v", err)
	}
}

// 测试消息释放
func TestMessageRelease(t *testing.T) {
	msg := &Message{
		Data:     []byte("test"),
		cPtr:     0,
		released: false,
	}

	// 测试释放
	ReleaseMessage(msg)
	if !msg.released {
		t.Error("Message should be marked as released")
	}

	// 测试重复释放（应该是安全的）
	ReleaseMessage(msg)
	if !msg.released {
		t.Error("Message should still be marked as released")
	}

	// 测试释放 nil 消息（应该是安全的）
	ReleaseMessage(nil)
}

// 测试 FlushMessages
func TestFlushMessages(t *testing.T) {
	testConn := Connection(1)

	mock := &MockSockets{
		FlushMessagesOnConnectionFunc: func(conn Connection) error {
			if conn == InvalidConnection {
				return ErrInvalidConnection
			}
			return nil
		},
	}

	// 测试刷新消息
	err := mock.FlushMessagesOnConnection(testConn)
	if err != nil {
		t.Errorf("FlushMessagesOnConnection() error = %v", err)
	}

	// 测试无效连接
	err = mock.FlushMessagesOnConnection(InvalidConnection)
	if !IsInvalidConnection(err) {
		t.Errorf("FlushMessagesOnConnection() with invalid connection should return ErrInvalidConnection, got %v", err)
	}
}

// 测试 GetConnectionRealTimeStatus
func TestGetConnectionRealTimeStatus(t *testing.T) {
	testConn := Connection(1)

	mock := &MockSockets{
		GetConnectionRealTimeStatusFunc: func(conn Connection) (*QuickConnectionStatus, error) {
			if conn == InvalidConnection {
				return nil, ErrInvalidConnection
			}
			return &QuickConnectionStatus{
				State:             ConnectionStateConnected,
				Ping:              50,
				ConnectionQuality: 0.95,
				OutPacketsPerSec:  100.0,
				OutBytesPerSec:    10000.0,
				InPacketsPerSec:   100.0,
				InBytesPerSec:     10000.0,
			}, nil
		},
	}

	// 测试获取状态
	status, err := mock.GetConnectionRealTimeStatus(testConn)
	if err != nil {
		t.Errorf("GetConnectionRealTimeStatus() error = %v", err)
	}
	if status == nil {
		t.Fatal("GetConnectionRealTimeStatus() returned nil status")
	}
	if status.State != ConnectionStateConnected {
		t.Errorf("State = %v, want %v", status.State, ConnectionStateConnected)
	}
	if status.Ping != 50 {
		t.Errorf("Ping = %d, want 50", status.Ping)
	}
	if status.ConnectionQuality != 0.95 {
		t.Errorf("ConnectionQuality = %f, want 0.95", status.ConnectionQuality)
	}

	// 测试无效连接
	_, err = mock.GetConnectionRealTimeStatus(InvalidConnection)
	if !IsInvalidConnection(err) {
		t.Errorf("GetConnectionRealTimeStatus() with invalid connection should return ErrInvalidConnection, got %v", err)
	}
}
