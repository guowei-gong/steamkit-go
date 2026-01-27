//go:build windows

package purego

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// loadLibPlatform 加载 Windows 平台的 Steam 库
func loadLibPlatform() (uintptr, error) {
	// 尝试多个可能的路径
	paths := []string{
		"steam_api64.dll",                                    // 当前目录
		filepath.Join(".", "steam_api64.dll"),                // 当前目录（显式）
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Steam", "steam_api64.dll"), // Steam 安装目录
	}

	var lastErr error
	for _, path := range paths {
		handle, err := syscall.LoadLibrary(path)
		if err == nil {
			return uintptr(handle), nil
		}
		lastErr = err
	}

	return 0, fmt.Errorf("failed to load steam_api64.dll: %w", lastErr)
}
