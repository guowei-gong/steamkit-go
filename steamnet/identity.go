package steamnet

import (
	"fmt"
	"net"
)

// IdentityType 表示身份类型
type IdentityType int

const (
	IdentityTypeInvalid IdentityType = 0 // 无效身份
	IdentityTypeSteamID IdentityType = 1 // SteamID
	IdentityTypeIPAddr  IdentityType = 2 // IP 地址
)

// Identity 表示网络端点的身份
type Identity struct {
	identityType IdentityType
	steamID      uint64
	ipAddr       string
	port         uint16
}

// NewIdentityFromSteamID 从 SteamID 创建身份
func NewIdentityFromSteamID(steamID uint64) Identity {
	return Identity{
		identityType: IdentityTypeSteamID,
		steamID:      steamID,
	}
}

// NewIdentityFromIPAddr 从 IP 地址创建身份
func NewIdentityFromIPAddr(ip string, port uint16) Identity {
	return Identity{
		identityType: IdentityTypeIPAddr,
		ipAddr:       ip,
		port:         port,
	}
}

// NewInvalidIdentity 创建无效身份
func NewInvalidIdentity() Identity {
	return Identity{
		identityType: IdentityTypeInvalid,
	}
}

// GetSteamID 获取 SteamID
// 如果身份类型不是 SteamID，返回 0
func (i Identity) GetSteamID() uint64 {
	if i.identityType == IdentityTypeSteamID {
		return i.steamID
	}
	return 0
}

// GetIPAddr 获取 IP 地址和端口
// 如果身份类型不是 IP 地址，返回空字符串和 0
func (i Identity) GetIPAddr() (string, uint16) {
	if i.identityType == IdentityTypeIPAddr {
		return i.ipAddr, i.port
	}
	return "", 0
}

// Type 返回身份类型
func (i Identity) Type() IdentityType {
	return i.identityType
}

// IsValid 检查身份是否有效
func (i Identity) IsValid() bool {
	return i.identityType != IdentityTypeInvalid
}

// String 返回身份的字符串表示
func (i Identity) String() string {
	switch i.identityType {
	case IdentityTypeSteamID:
		return fmt.Sprintf("SteamID:%d", i.steamID)
	case IdentityTypeIPAddr:
		return fmt.Sprintf("IP:%s:%d", i.ipAddr, i.port)
	case IdentityTypeInvalid:
		return "Invalid"
	default:
		return "Unknown"
	}
}

// Equal 比较两个身份是否相等
func (i Identity) Equal(other Identity) bool {
	if i.identityType != other.identityType {
		return false
	}

	switch i.identityType {
	case IdentityTypeSteamID:
		return i.steamID == other.steamID
	case IdentityTypeIPAddr:
		return i.ipAddr == other.ipAddr && i.port == other.port
	case IdentityTypeInvalid:
		return true
	default:
		return false
	}
}

// ParseIdentity 从字符串解析身份
// 支持格式：
//   - "steamid:76561198000000000" - SteamID
//   - "ip:192.168.1.1:27015" - IP 地址
func ParseIdentity(s string) (Identity, error) {
	// 尝试解析为 SteamID
	var steamID uint64
	if n, _ := fmt.Sscanf(s, "steamid:%d", &steamID); n == 1 {
		return NewIdentityFromSteamID(steamID), nil
	}

	// 尝试解析为 IP 地址
	var ip string
	var port uint16
	if n, _ := fmt.Sscanf(s, "ip:%s", &ip); n == 1 {
		// 解析 IP:Port
		host, portStr, err := net.SplitHostPort(ip)
		if err != nil {
			return NewInvalidIdentity(), fmt.Errorf("invalid IP address format: %w", err)
		}
		if n, _ := fmt.Sscanf(portStr, "%d", &port); n == 1 {
			return NewIdentityFromIPAddr(host, port), nil
		}
	}

	return NewInvalidIdentity(), fmt.Errorf("invalid identity format: %s", s)
}
