package steamnet

import (
	"sync"
)

// ConnectionStatusChangedCallback 是连接状态变化的回调函数类型
type ConnectionStatusChangedCallback func(info *ConnectionStatusChangedInfo)

// callbackManager 管理所有回调
type callbackManager struct {
	mu                  sync.RWMutex
	connectionCallbacks map[Connection]ConnectionStatusChangedCallback
	globalCallback      ConnectionStatusChangedCallback
}

var (
	globalCallbackManager = &callbackManager{
		connectionCallbacks: make(map[Connection]ConnectionStatusChangedCallback),
	}
)

// SetConnectionStatusChangedCallback 设置全局连接状态变化回调
// 当任何连接的状态发生变化时，都会调用此回调
func SetConnectionStatusChangedCallback(callback ConnectionStatusChangedCallback) {
	globalCallbackManager.mu.Lock()
	defer globalCallbackManager.mu.Unlock()
	globalCallbackManager.globalCallback = callback
}

// SetConnectionCallback 为特定连接设置状态变化回调
// 当指定连接的状态发生变化时，会优先调用此回调而不是全局回调
func SetConnectionCallback(conn Connection, callback ConnectionStatusChangedCallback) {
	globalCallbackManager.mu.Lock()
	defer globalCallbackManager.mu.Unlock()
	if callback == nil {
		delete(globalCallbackManager.connectionCallbacks, conn)
	} else {
		globalCallbackManager.connectionCallbacks[conn] = callback
	}
}

// ClearConnectionCallback 清除特定连接的回调
func ClearConnectionCallback(conn Connection) {
	SetConnectionCallback(conn, nil)
}

// DispatchConnectionStatusChanged 分发连接状态变化回调
// 这个函数由内部调用，用户不应直接调用
func DispatchConnectionStatusChanged(info *ConnectionStatusChangedInfo) {
	globalCallbackManager.mu.RLock()
	defer globalCallbackManager.mu.RUnlock()

	// 先尝试调用特定连接的回调
	if callback, ok := globalCallbackManager.connectionCallbacks[info.Connection]; ok {
		callback(info)
		return
	}

	// 否则调用全局回调
	if globalCallbackManager.globalCallback != nil {
		globalCallbackManager.globalCallback(info)
	}
}

