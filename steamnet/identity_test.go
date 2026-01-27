package steamnet

import (
	"testing"
)

func TestNewIdentityFromSteamID(t *testing.T) {
	steamID := uint64(76561198000000000)
	identity := NewIdentityFromSteamID(steamID)

	if identity.Type() != IdentityTypeSteamID {
		t.Errorf("Type() = %v, want %v", identity.Type(), IdentityTypeSteamID)
	}

	if identity.GetSteamID() != steamID {
		t.Errorf("GetSteamID() = %v, want %v", identity.GetSteamID(), steamID)
	}

	if !identity.IsValid() {
		t.Error("IsValid() = false, want true")
	}
}

func TestNewIdentityFromIPAddr(t *testing.T) {
	ip := "192.168.1.1"
	port := uint16(27015)
	identity := NewIdentityFromIPAddr(ip, port)

	if identity.Type() != IdentityTypeIPAddr {
		t.Errorf("Type() = %v, want %v", identity.Type(), IdentityTypeIPAddr)
	}

	gotIP, gotPort := identity.GetIPAddr()
	if gotIP != ip || gotPort != port {
		t.Errorf("GetIPAddr() = (%v, %v), want (%v, %v)", gotIP, gotPort, ip, port)
	}

	if !identity.IsValid() {
		t.Error("IsValid() = false, want true")
	}
}

func TestNewInvalidIdentity(t *testing.T) {
	identity := NewInvalidIdentity()

	if identity.Type() != IdentityTypeInvalid {
		t.Errorf("Type() = %v, want %v", identity.Type(), IdentityTypeInvalid)
	}

	if identity.IsValid() {
		t.Error("IsValid() = true, want false")
	}

	if identity.GetSteamID() != 0 {
		t.Errorf("GetSteamID() = %v, want 0", identity.GetSteamID())
	}
}

func TestIdentity_String(t *testing.T) {
	tests := []struct {
		name     string
		identity Identity
		expected string
	}{
		{
			name:     "SteamID",
			identity: NewIdentityFromSteamID(76561198000000000),
			expected: "SteamID:76561198000000000",
		},
		{
			name:     "IPAddr",
			identity: NewIdentityFromIPAddr("192.168.1.1", 27015),
			expected: "IP:192.168.1.1:27015",
		},
		{
			name:     "Invalid",
			identity: NewInvalidIdentity(),
			expected: "Invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.identity.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIdentity_Equal(t *testing.T) {
	steamID1 := NewIdentityFromSteamID(76561198000000000)
	steamID2 := NewIdentityFromSteamID(76561198000000000)
	steamID3 := NewIdentityFromSteamID(76561198000000001)

	ipAddr1 := NewIdentityFromIPAddr("192.168.1.1", 27015)
	ipAddr2 := NewIdentityFromIPAddr("192.168.1.1", 27015)
	ipAddr3 := NewIdentityFromIPAddr("192.168.1.2", 27015)

	invalid1 := NewInvalidIdentity()
	invalid2 := NewInvalidIdentity()

	tests := []struct {
		name     string
		a        Identity
		b        Identity
		expected bool
	}{
		{"SameS teamID", steamID1, steamID2, true},
		{"DifferentSteamID", steamID1, steamID3, false},
		{"SameIPAddr", ipAddr1, ipAddr2, true},
		{"DifferentIPAddr", ipAddr1, ipAddr3, false},
		{"SameInvalid", invalid1, invalid2, true},
		{"DifferentTypes", steamID1, ipAddr1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.b); got != tt.expected {
				t.Errorf("Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseIdentity(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
		wantType  IdentityType
	}{
		{
			name:      "ValidSteamID",
			input:     "steamid:76561198000000000",
			wantValid: true,
			wantType:  IdentityTypeSteamID,
		},
		{
			name:      "ValidIPAddr",
			input:     "ip:192.168.1.1:27015",
			wantValid: true,
			wantType:  IdentityTypeIPAddr,
		},
		{
			name:      "InvalidFormat",
			input:     "invalid",
			wantValid: false,
			wantType:  IdentityTypeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity, err := ParseIdentity(tt.input)

			if tt.wantValid && err != nil {
				t.Errorf("ParseIdentity() error = %v, want nil", err)
			}

			if !tt.wantValid && err == nil {
				t.Error("ParseIdentity() error = nil, want error")
			}

			if identity.Type() != tt.wantType {
				t.Errorf("Type() = %v, want %v", identity.Type(), tt.wantType)
			}
		})
	}
}
