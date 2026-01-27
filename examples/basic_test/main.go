package main

import (
	"fmt"
	"log"

	"github.com/guowei-gong/steamkit-go"
)

func main() {
	fmt.Println("SteamKit-Go 基础测试")
	fmt.Println("====================")

	// 初始化 Steam API
	fmt.Println("正在初始化 Steam API...")
	if err := steamkit.Init(); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	defer steamkit.Shutdown()

	fmt.Println("✓ Steam API 初始化成功")

	// 获取 SteamID
	steamID := steamkit.GetSteamID()
	if steamID == 0 {
		log.Fatal("无法获取 SteamID（可能未登录 Steam）")
	}

	fmt.Printf("✓ 当前用户 SteamID: %d\n", steamID)
	fmt.Println("\n基础设施测试通过！")
}
