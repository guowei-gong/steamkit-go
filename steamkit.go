// Package steamkit 提供 Steamworks SDK 的 Go 语言绑定
package steamkit

import (
	"fmt"

	"github.com/guowei-gong/steamkit-go/internal/purego"
)

// ESteamAPIInitResult 表示 Steam API 初始化结果
type ESteamAPIInitResult int32

const (
	ESteamAPIInitResult_OK              ESteamAPIInitResult = 0
	ESteamAPIInitResult_FailedGeneric   ESteamAPIInitResult = 1
	ESteamAPIInitResult_NoSteamClient   ESteamAPIInitResult = 2
	ESteamAPIInitResult_VersionMismatch ESteamAPIInitResult = 3
)

// RestartAppIfNecessary 检查是否需要通过 Steam 重启应用
// 如果返回 true，应用应该立即退出
func RestartAppIfNecessary(appID uint32) bool {
	return purego.CallRestartAppIfNecessary(appID)
}

// Init 初始化 Steam API
// 必须在使用任何其他 Steam API 之前调用
func Init() error {
	// 初始化 purego 绑定层
	if err := purego.Init(); err != nil {
		return fmt.Errorf("failed to initialize purego: %w", err)
	}

	// 调用 SteamAPI_InitFlat
	result, errMsg := purego.CallInitFlat()
	if ESteamAPIInitResult(result) != ESteamAPIInitResult_OK {
		if errMsg != "" {
			return fmt.Errorf("SteamAPI_InitFlat failed: %s (code: %d)", errMsg, result)
		}
		return fmt.Errorf("SteamAPI_InitFlat failed with code: %d", result)
	}

	return nil
}

// Shutdown 关闭 Steam API
// 应该在程序退出前调用
func Shutdown() {
	purego.CallShutdown()
}

// RunCallbacks 处理 Steam 回调
// 应该在主循环中定期调用（建议每 10-50ms）
func RunCallbacks() {
	purego.CallRunCallbacks()
}

// GetSteamID 获取当前用户的 SteamID
func GetSteamID() uint64 {
	return purego.CallGetSteamID()
}
