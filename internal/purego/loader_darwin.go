//go:build darwin

package purego

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ebitengine/purego"
)

// loadLibPlatform 加载 macOS 平台的 Steam 库
func loadLibPlatform() (uintptr, error) {
	// 尝试多个可能的路径
	paths := []string{
		"libsteam_api.dylib",                                 // 当前目录
		filepath.Join(".", "libsteam_api.dylib"),             // 当前目录（显式）
		filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Steam", "Steam.AppBundle", "Steam", "Contents", "MacOS", "libsteam_api.dylib"), // Steam 安装目录
		"/usr/local/lib/libsteam_api.dylib",                  // 本地库目录
	}

	var lastErr error
	for _, path := range paths {
		lib, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err == nil {
			return lib, nil
		}
		lastErr = err
	}

	return 0, fmt.Errorf("failed to load libsteam_api.dylib: %w", lastErr)
}
