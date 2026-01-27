package steamnet

import (
	"sync"
	"testing"
)

// 测试全局回调设置
func TestSetConnectionStatusChangedCallback(t *testing.T) {
	called := false
	testInfo := &ConnectionStatusChangedInfo{
		Connection: Connection(1),
		OldState:   ConnectionStateNone,
		NewState:   ConnectionStateConnecting,
	}

	// 设置全局回调
	SetConnectionStatusChangedCallback(func(info *ConnectionStatusChangedInfo) {
		called = true
		if info.Connection != testInfo.Connection {
			t.Errorf("Connection = %v, want %v", info.Connection, testInfo.Connection)
		}
		if info.OldState != testInfo.OldState {
			t.Errorf("OldState = %v, want %v", info.OldState, testInfo.OldState)
		}
		if info.NewState != testInfo.NewState {
			t.Errorf("NewState = %v, want %v", info.NewState, testInfo.NewState)
		}
	})

	// 分发回调
	DispatchConnectionStatusChanged(testInfo)

	if !called {
		t.Error("Global callback was not called")
	}

	// 清除回调
	SetConnectionStatusChangedCallback(nil)
}

// 测试特定连接回调
func TestSetConnectionCallback(t *testing.T) {
	conn1 := Connection(1)
	conn2 := Connection(2)

	called1 := false
	called2 := false
	globalCalled := false

	// 设置全局回调
	SetConnectionStatusChangedCallback(func(info *ConnectionStatusChangedInfo) {
		globalCalled = true
	})

	// 设置连接 1 的回调
	SetConnectionCallback(conn1, func(info *ConnectionStatusChangedInfo) {
		called1 = true
		if info.Connection != conn1 {
			t.Errorf("Connection = %v, want %v", info.Connection, conn1)
		}
	})

	// 设置连接 2 的回调
	SetConnectionCallback(conn2, func(info *ConnectionStatusChangedInfo) {
		called2 = true
		if info.Connection != conn2 {
			t.Errorf("Connection = %v, want %v", info.Connection, conn2)
		}
	})

	// 分发连接 1 的回调
	DispatchConnectionStatusChanged(&ConnectionStatusChangedInfo{
		Connection: conn1,
		OldState:   ConnectionStateNone,
		NewState:   ConnectionStateConnecting,
	})

	if !called1 {
		t.Error("Connection 1 callback was not called")
	}
	if called2 {
		t.Error("Connection 2 callback should not be called")
	}
	if globalCalled {
		t.Error("Global callback should not be called when connection callback exists")
	}

	// 重置标志
	called1 = false
	called2 = false
	globalCalled = false

	// 分发连接 2 的回调
	DispatchConnectionStatusChanged(&ConnectionStatusChangedInfo{
		Connection: conn2,
		OldState:   ConnectionStateNone,
		NewState:   ConnectionStateConnecting,
	})

	if called1 {
		t.Error("Connection 1 callback should not be called")
	}
	if !called2 {
		t.Error("Connection 2 callback was not called")
	}
	if globalCalled {
		t.Error("Global callback should not be called when connection callback exists")
	}

	// 清除回调
	SetConnectionCallback(conn1, nil)
	SetConnectionCallback(conn2, nil)
	SetConnectionStatusChangedCallback(nil)
}

// 测试清除连接回调
func TestClearConnectionCallback(t *testing.T) {
	conn := Connection(1)
	called := false
	globalCalled := false

	// 设置全局回调
	SetConnectionStatusChangedCallback(func(info *ConnectionStatusChangedInfo) {
		globalCalled = true
	})

	// 设置连接回调
	SetConnectionCallback(conn, func(info *ConnectionStatusChangedInfo) {
		called = true
	})

	// 清除连接回调
	ClearConnectionCallback(conn)

	// 分发回调
	DispatchConnectionStatusChanged(&ConnectionStatusChangedInfo{
		Connection: conn,
		OldState:   ConnectionStateNone,
		NewState:   ConnectionStateConnecting,
	})

	if called {
		t.Error("Connection callback should not be called after clearing")
	}
	if !globalCalled {
		t.Error("Global callback should be called after clearing connection callback")
	}

	// 清除回调
	SetConnectionStatusChangedCallback(nil)
}

// 测试并发回调
func TestConcurrentCallbacks(t *testing.T) {
	const numGoroutines = 100
	const numCallbacks = 10

	var wg sync.WaitGroup
	callCount := 0
	var mu sync.Mutex

	// 设置全局回调
	SetConnectionStatusChangedCallback(func(info *ConnectionStatusChangedInfo) {
		mu.Lock()
		callCount++
		mu.Unlock()
	})

	// 并发分发回调
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numCallbacks; j++ {
				DispatchConnectionStatusChanged(&ConnectionStatusChangedInfo{
					Connection: Connection(id),
					OldState:   ConnectionStateNone,
					NewState:   ConnectionStateConnecting,
				})
			}
		}(i)
	}

	wg.Wait()

	expectedCount := numGoroutines * numCallbacks
	if callCount != expectedCount {
		t.Errorf("Call count = %d, want %d", callCount, expectedCount)
	}

	// 清除回调
	SetConnectionStatusChangedCallback(nil)
}

// 测试回调优先级
func TestCallbackPriority(t *testing.T) {
	conn := Connection(1)
	globalCalled := false
	connCalled := false

	// 设置全局回调
	SetConnectionStatusChangedCallback(func(info *ConnectionStatusChangedInfo) {
		globalCalled = true
	})

	// 设置连接回调
	SetConnectionCallback(conn, func(info *ConnectionStatusChangedInfo) {
		connCalled = true
	})

	// 分发回调
	DispatchConnectionStatusChanged(&ConnectionStatusChangedInfo{
		Connection: conn,
		OldState:   ConnectionStateNone,
		NewState:   ConnectionStateConnecting,
	})

	// 连接回调应该被调用，全局回调不应该被调用
	if !connCalled {
		t.Error("Connection callback should be called")
	}
	if globalCalled {
		t.Error("Global callback should not be called when connection callback exists")
	}

	// 清除回调
	SetConnectionCallback(conn, nil)
	SetConnectionStatusChangedCallback(nil)
}

// 测试无回调情况
func TestNoCallback(t *testing.T) {
	// 确保没有设置回调
	SetConnectionStatusChangedCallback(nil)

	// 分发回调（不应该崩溃）
	DispatchConnectionStatusChanged(&ConnectionStatusChangedInfo{
		Connection: Connection(1),
		OldState:   ConnectionStateNone,
		NewState:   ConnectionStateConnecting,
	})

	// 如果没有崩溃，测试通过
}
