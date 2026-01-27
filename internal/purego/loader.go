// Package purego 提供 Steamworks SDK 的 purego 绑定层
package purego

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	steamLib uintptr
)

// loadLib 加载 Steam 库（平台特定实现）
func loadLib() (uintptr, error) {
	return loadLibPlatform()
}

// Init 初始化 purego 绑定层
func Init() error {
	lib, err := loadLib()
	if err != nil {
		return fmt.Errorf("failed to load Steam library: %w", err)
	}
	steamLib = lib

	// 注册所有函数
	if err := registerFunctions(); err != nil {
		return fmt.Errorf("failed to register functions: %w", err)
	}

	return nil
}

// GetLibHandle 返回 Steam 库句柄
func GetLibHandle() uintptr {
	return steamLib
}

// 错误消息缓冲区
type steamErrMsg [1024]byte

func (s *steamErrMsg) String() string {
	for i, b := range s {
		if b == 0 {
			return string(s[:i])
		}
	}
	return ""
}

// 函数指针变量
var (
	// 通用函数
	ptrAPI_RestartAppIfNecessary func(uint32) bool
	ptrAPI_InitFlat              func(uintptr) int32
	ptrAPI_Shutdown              func()
	ptrAPI_RunCallbacks          func()

	// ISteamUser
	ptrAPI_SteamUser             func() uintptr
	ptrAPI_ISteamUser_GetSteamID func(uintptr) uint64

	// ISteamNetworkingSockets
	ptrAPI_SteamNetworkingSockets                      func() uintptr
	ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P func(uintptr, int32, int32, uintptr) uint32
	ptrAPI_ISteamNetworkingSockets_ConnectP2P            func(uintptr, uintptr, int32, int32, uintptr) uint32
	ptrAPI_ISteamNetworkingSockets_AcceptConnection      func(uintptr, uint32) int32
	ptrAPI_ISteamNetworkingSockets_CloseConnection       func(uintptr, uint32, int32, uintptr, bool) bool
	ptrAPI_ISteamNetworkingSockets_CloseListenSocket     func(uintptr, uint32) bool
	ptrAPI_ISteamNetworkingSockets_GetConnectionInfo     func(uintptr, uint32, uintptr) bool
	ptrAPI_ISteamNetworkingSockets_SendMessageToConnection func(uintptr, uint32, uintptr, uint32, int32, uintptr) int32
	ptrAPI_ISteamNetworkingSockets_FlushMessagesOnConnection func(uintptr, uint32) int32
	ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection func(uintptr, uint32, uintptr, int32) int32
	ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnListenSocket func(uintptr, uint32, uintptr, int32) int32
	ptrAPI_ISteamNetworkingSockets_GetConnectionRealTimeStatus func(uintptr, uint32, uintptr, int32, uintptr) int32
	ptrAPI_SteamNetworkingMessage_t_Release func(uintptr)
)

// registerFunctions 注册所有 Steam API 函数
func registerFunctions() error {
	// 通用函数
	purego.RegisterLibFunc(&ptrAPI_RestartAppIfNecessary, steamLib, "SteamAPI_RestartAppIfNecessary")
	purego.RegisterLibFunc(&ptrAPI_InitFlat, steamLib, "SteamAPI_InitFlat")
	purego.RegisterLibFunc(&ptrAPI_Shutdown, steamLib, "SteamAPI_Shutdown")
	purego.RegisterLibFunc(&ptrAPI_RunCallbacks, steamLib, "SteamAPI_RunCallbacks")

	// ISteamUser
	purego.RegisterLibFunc(&ptrAPI_SteamUser, steamLib, "SteamAPI_SteamUser_v023")
	purego.RegisterLibFunc(&ptrAPI_ISteamUser_GetSteamID, steamLib, "SteamAPI_ISteamUser_GetSteamID")

	// ISteamNetworkingSockets
	purego.RegisterLibFunc(&ptrAPI_SteamNetworkingSockets, steamLib, "SteamAPI_SteamNetworkingSockets_SteamAPI_v012")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P, steamLib, "SteamAPI_ISteamNetworkingSockets_CreateListenSocketP2P")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_ConnectP2P, steamLib, "SteamAPI_ISteamNetworkingSockets_ConnectP2P")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_AcceptConnection, steamLib, "SteamAPI_ISteamNetworkingSockets_AcceptConnection")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_CloseConnection, steamLib, "SteamAPI_ISteamNetworkingSockets_CloseConnection")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_CloseListenSocket, steamLib, "SteamAPI_ISteamNetworkingSockets_CloseListenSocket")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_GetConnectionInfo, steamLib, "SteamAPI_ISteamNetworkingSockets_GetConnectionInfo")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_SendMessageToConnection, steamLib, "SteamAPI_ISteamNetworkingSockets_SendMessageToConnection")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_FlushMessagesOnConnection, steamLib, "SteamAPI_ISteamNetworkingSockets_FlushMessagesOnConnection")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection, steamLib, "SteamAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnListenSocket, steamLib, "SteamAPI_ISteamNetworkingSockets_ReceiveMessagesOnListenSocket")
	purego.RegisterLibFunc(&ptrAPI_ISteamNetworkingSockets_GetConnectionRealTimeStatus, steamLib, "SteamAPI_ISteamNetworkingSockets_GetConnectionRealTimeStatus")
	purego.RegisterLibFunc(&ptrAPI_SteamNetworkingMessage_t_Release, steamLib, "SteamAPI_SteamNetworkingMessage_t_Release")

	return nil
}

// CallRestartAppIfNecessary 调用 SteamAPI_RestartAppIfNecessary
func CallRestartAppIfNecessary(appID uint32) bool {
	return ptrAPI_RestartAppIfNecessary(appID)
}

// CallInitFlat 调用 SteamAPI_InitFlat
func CallInitFlat() (int32, string) {
	var msg steamErrMsg
	result := ptrAPI_InitFlat(uintptr(unsafe.Pointer(&msg)))
	return result, msg.String()
}

// CallShutdown 调用 SteamAPI_Shutdown
func CallShutdown() {
	ptrAPI_Shutdown()
}

// CallRunCallbacks 调用 SteamAPI_RunCallbacks
func CallRunCallbacks() {
	ptrAPI_RunCallbacks()
}

// CallGetSteamID 获取当前用户的 SteamID
func CallGetSteamID() uint64 {
	userPtr := ptrAPI_SteamUser()
	if userPtr == 0 {
		return 0
	}
	return ptrAPI_ISteamUser_GetSteamID(userPtr)
}

// CallGetSteamNetworkingSockets 获取 ISteamNetworkingSockets 接口指针
func CallGetSteamNetworkingSockets() uintptr {
	return ptrAPI_SteamNetworkingSockets()
}

// CallCreateListenSocketP2P 创建 P2P 监听套接字
func CallCreateListenSocketP2P(handle uintptr, virtualPort int32, numOptions int32, options uintptr) uint32 {
	return ptrAPI_ISteamNetworkingSockets_CreateListenSocketP2P(handle, virtualPort, numOptions, options)
}

// CallConnectP2P 连接到远程 P2P 对等方
func CallConnectP2P(handle uintptr, identityRemote uintptr, virtualPort int32, numOptions int32, options uintptr) uint32 {
	return ptrAPI_ISteamNetworkingSockets_ConnectP2P(handle, identityRemote, virtualPort, numOptions, options)
}

// CallAcceptConnection 接受传入连接
func CallAcceptConnection(handle uintptr, conn uint32) int32 {
	return ptrAPI_ISteamNetworkingSockets_AcceptConnection(handle, conn)
}

// CallCloseConnection 关闭连接
func CallCloseConnection(handle uintptr, peer uint32, reason int32, debug uintptr, linger bool) bool {
	return ptrAPI_ISteamNetworkingSockets_CloseConnection(handle, peer, reason, debug, linger)
}

// CallCloseListenSocket 关闭监听套接字
func CallCloseListenSocket(handle uintptr, socket uint32) bool {
	return ptrAPI_ISteamNetworkingSockets_CloseListenSocket(handle, socket)
}

// CallGetConnectionInfo 获取连接信息
func CallGetConnectionInfo(handle uintptr, conn uint32, info uintptr) bool {
	return ptrAPI_ISteamNetworkingSockets_GetConnectionInfo(handle, conn, info)
}

// CallSendMessageToConnection 发送消息到连接
func CallSendMessageToConnection(handle uintptr, conn uint32, data uintptr, dataSize uint32, flags int32, outMessageNumber uintptr) int32 {
	return ptrAPI_ISteamNetworkingSockets_SendMessageToConnection(handle, conn, data, dataSize, flags, outMessageNumber)
}

// CallFlushMessagesOnConnection 刷新连接上的消息
func CallFlushMessagesOnConnection(handle uintptr, conn uint32) int32 {
	return ptrAPI_ISteamNetworkingSockets_FlushMessagesOnConnection(handle, conn)
}

// CallReceiveMessagesOnConnection 接收连接上的消息
func CallReceiveMessagesOnConnection(handle uintptr, conn uint32, messages uintptr, maxMessages int32) int32 {
	return ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnConnection(handle, conn, messages, maxMessages)
}

// CallReceiveMessagesOnListenSocket 接收监听套接字上的消息
func CallReceiveMessagesOnListenSocket(handle uintptr, socket uint32, messages uintptr, maxMessages int32) int32 {
	return ptrAPI_ISteamNetworkingSockets_ReceiveMessagesOnListenSocket(handle, socket, messages, maxMessages)
}

// CallReleaseMessage 释放消息内存
func CallReleaseMessage(messagePtr uintptr) {
	ptrAPI_SteamNetworkingMessage_t_Release(messagePtr)
}

// CallGetConnectionRealTimeStatus 获取连接的实时状态
func CallGetConnectionRealTimeStatus(handle uintptr, conn uint32, status uintptr, numLanes int32, lanes uintptr) int32 {
	return ptrAPI_ISteamNetworkingSockets_GetConnectionRealTimeStatus(handle, conn, status, numLanes, lanes)
}
