//go:build linux

package purego

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ebitengine/purego"
)

// loadLibPlatform 加载 Linux 平台的 Steam 库
func loadLibPlatform() (uintptr, error) {
	// 尝试多个可能的路径
	paths := []string{
		"libsteam_api.so",                                    // 当前目录
		filepath.Join(".", "libsteam_api.so"),                // 当前目录（显式）
		filepath.Join(os.Getenv("HOME"), ".steam", "sdk64", "libsteam_api.so"), // Steam SDK 目录
		"/usr/lib/libsteam_api.so",                           // 系统库目录
		"/usr/local/lib/libsteam_api.so",                     // 本地库目录
	}

	var lastErr error
	for _, path := range paths {
		lib, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err == nil {
			return lib, nil
		}
		lastErr = err
	}

	return 0, fmt.Errorf("failed to load libsteam_api.so: %w", lastErr)
}
