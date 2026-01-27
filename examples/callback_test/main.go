package main

import (
	"fmt"
	"log"
	"time"

	"github.com/guowei-gong/steamkit-go"
	"github.com/guowei-gong/steamkit-go/steamnet"
)

func main() {
	// 初始化 Steam API
	if err := steamkit.Init(); err != nil {
		log.Fatal(err)
	}
	defer steamkit.Shutdown()

	// 获取当前用户 SteamID
	steamID := steamkit.GetSteamID()
	fmt.Printf("当前用户 SteamID: %d\n", steamID)

	// 设置全局连接状态变化回调
	steamnet.SetConnectionStatusChangedCallback(func(info *steamnet.ConnectionStatusChangedInfo) {
		fmt.Printf("\n[全局回调] 连接状态变化:\n")
		fmt.Printf("  连接: %d\n", info.Connection)
		fmt.Printf("  旧状态: %s\n", info.OldState)
		fmt.Printf("  新状态: %s\n", info.NewState)
		if info.EndReason != 0 {
			fmt.Printf("  结束原因: %d\n", info.EndReason)
			fmt.Printf("  调试信息: %s\n", info.EndDebug)
		}
	})

	// 获取 ISteamNetworkingSockets 接口
	sockets := steamnet.GetSockets()
	if sockets == nil {
		log.Fatal("无法获取 ISteamNetworkingSockets 接口")
	}

	// 创建监听套接字
	fmt.Println("\n创建 P2P 监听套接字...")
	listenSocket, err := sockets.CreateListenSocketP2P(0, nil)
	if err != nil {
		log.Fatalf("创建监听套接字失败: %v", err)
	}
	defer sockets.CloseListenSocket(listenSocket)
	fmt.Printf("监听套接字创建成功: %d\n", listenSocket)

	// 模拟连接到远程对等方
	fmt.Println("\n连接到远程对等方...")
	remoteSteamID := uint64(76561198000000000) // 示例 SteamID
	identity := steamnet.NewIdentityFromSteamID(remoteSteamID)

	conn, err := sockets.ConnectP2P(identity, 0, nil)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	fmt.Printf("连接创建成功: %d\n", conn)

	// 为特定连接设置回调
	steamnet.SetConnectionCallback(conn, func(info *steamnet.ConnectionStatusChangedInfo) {
		fmt.Printf("\n[连接回调] 连接 %d 状态变化:\n", info.Connection)
		fmt.Printf("  旧状态: %s\n", info.OldState)
		fmt.Printf("  新状态: %s\n", info.NewState)

		// 根据状态执行不同操作
		switch info.NewState {
		case steamnet.ConnectionStateConnected:
			fmt.Println("  ✓ 连接已建立")
		case steamnet.ConnectionStateClosedByPeer:
			fmt.Println("  ✗ 对方关闭了连接")
		case steamnet.ConnectionStateProblemDetectedLocally:
			fmt.Println("  ✗ 检测到本地问题")
		}
	})

	// 获取连接信息
	fmt.Println("\n获取连接信息...")
	info, err := sockets.GetConnectionInfo(conn)
	if err != nil {
		log.Printf("获取连接信息失败: %v", err)
	} else {
		fmt.Printf("连接状态: %s\n", info.State)
	}

	// 模拟发送消息
	fmt.Println("\n发送测试消息...")
	testData := []byte("Hello, Steam P2P!")
	err = sockets.SendMessageToConnection(conn, testData, steamnet.SendReliable)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	} else {
		fmt.Println("消息发送成功")
	}

	// 刷新消息
	err = sockets.FlushMessagesOnConnection(conn)
	if err != nil {
		log.Printf("刷新消息失败: %v", err)
	}

	// 运行回调处理循环
	fmt.Println("\n运行回调处理循环（5 秒）...")
	fmt.Println("注意：实际的回调需要 Steam 客户端运行并且有真实的网络事件")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			// 处理 Steam 回调
			steamkit.RunCallbacks()

			// 尝试接收消息
			messages, err := sockets.ReceiveMessagesOnConnection(conn, 32)
			if err != nil {
				// 忽略接收错误（可能没有消息）
				continue
			}

			// 处理接收到的消息
			for _, msg := range messages {
				fmt.Printf("\n收到消息: %s\n", string(msg.Data))
				msg.Release()
			}

		case <-timeout:
			fmt.Println("\n超时，退出循环")
			goto cleanup
		}
	}

cleanup:
	// 清除回调
	fmt.Println("\n清理资源...")
	steamnet.ClearConnectionCallback(conn)

	// 关闭连接
	err = sockets.CloseConnection(conn, 0, "正常关闭", false)
	if err != nil {
		log.Printf("关闭连接失败: %v", err)
	} else {
		fmt.Println("连接已关闭")
	}

	fmt.Println("\n示例程序结束")
}
