package main

import (
	"fmt"
	"log"

	"github.com/guowei-gong/steamkit-go"
	"github.com/guowei-gong/steamkit-go/steamnet"
)

func main() {
	fmt.Println("SteamKit-Go 连接管理测试")
	fmt.Println("========================")

	// 初始化 Steam API
	fmt.Println("\n1. 正在初始化 Steam API...")
	if err := steamkit.Init(); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	defer steamkit.Shutdown()
	fmt.Println("   ✓ Steam API 初始化成功")

	// 获取当前用户 SteamID
	steamID := steamkit.GetSteamID()
	if steamID == 0 {
		log.Fatal("   ✗ 无法获取 SteamID（可能未登录 Steam）")
	}
	fmt.Printf("   ✓ 当前用户 SteamID: %d\n", steamID)

	// 获取 ISteamNetworkingSockets 实例
	fmt.Println("\n2. 获取 ISteamNetworkingSockets 实例...")
	sockets := steamnet.GetSockets()
	if sockets == nil {
		log.Fatal("   ✗ 无法获取 ISteamNetworkingSockets 实例")
	}
	fmt.Println("   ✓ ISteamNetworkingSockets 实例获取成功")

	// 测试创建监听套接字
	fmt.Println("\n3. 测试创建 P2P 监听套接字...")
	listenSocket, err := sockets.CreateListenSocketP2P(0, nil)
	if err != nil {
		log.Printf("   ✗ 创建监听套接字失败: %v", err)
	} else {
		fmt.Printf("   ✓ 监听套接字创建成功: %d\n", listenSocket)

		// 关闭监听套接字
		fmt.Println("\n4. 关闭监听套接字...")
		if err := sockets.CloseListenSocket(listenSocket); err != nil {
			log.Printf("   ✗ 关闭监听套接字失败: %v", err)
		} else {
			fmt.Println("   ✓ 监听套接字关闭成功")
		}
	}

	// 测试连接到自己（仅用于演示）
	fmt.Println("\n5. 测试创建 P2P 连接...")
	identity := steamnet.NewIdentityFromSteamID(steamID)
	fmt.Printf("   目标身份: %s\n", identity.String())

	connection, err := sockets.ConnectP2P(identity, 0, nil)
	if err != nil {
		log.Printf("   ✗ 创建连接失败: %v", err)
	} else {
		fmt.Printf("   ✓ 连接创建成功: %d\n", connection)

		// 获取连接信息
		fmt.Println("\n6. 获取连接信息...")
		info, err := sockets.GetConnectionInfo(connection)
		if err != nil {
			log.Printf("   ✗ 获取连接信息失败: %v", err)
		} else {
			fmt.Printf("   ✓ 连接状态: %s\n", info.State.String())
		}

		// 关闭连接
		fmt.Println("\n7. 关闭连接...")
		if err := sockets.CloseConnection(connection, 0, "测试完成", false); err != nil {
			log.Printf("   ✗ 关闭连接失败: %v", err)
		} else {
			fmt.Println("   ✓ 连接关闭成功")
		}
	}

	fmt.Println("\n========================")
	fmt.Println("连接管理测试完成！")
}
