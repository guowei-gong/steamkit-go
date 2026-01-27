package steamnet

import (
	"testing"
)

func TestConnectionState_String(t *testing.T) {
	tests := []struct {
		state    ConnectionState
		expected string
	}{
		{ConnectionStateNone, "None"},
		{ConnectionStateConnecting, "Connecting"},
		{ConnectionStateFindingRoute, "FindingRoute"},
		{ConnectionStateConnected, "Connected"},
		{ConnectionStateClosedByPeer, "ClosedByPeer"},
		{ConnectionStateProblemDetectedLocally, "ProblemDetectedLocally"},
		{ConnectionState(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("ConnectionState.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSendFlags_String(t *testing.T) {
	tests := []struct {
		flag     SendFlags
		expected string
	}{
		{SendUnreliable, "Unreliable"},
		{SendUnreliableNoNagle, "UnreliableNoNagle"},
		{SendReliable, "Reliable"},
		{SendReliableNoNagle, "ReliableNoNagle"},
		{SendFlags(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.flag.String(); got != tt.expected {
				t.Errorf("SendFlags.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInvalidHandles(t *testing.T) {
	if InvalidConnection != 0 {
		t.Errorf("InvalidConnection = %v, want 0", InvalidConnection)
	}
	if InvalidListenSocket != 0 {
		t.Errorf("InvalidListenSocket = %v, want 0", InvalidListenSocket)
	}
	if InvalidPollGroup != 0 {
		t.Errorf("InvalidPollGroup = %v, want 0", InvalidPollGroup)
	}
}

func TestMessage_Release(t *testing.T) {
	msg := &Message{
		Data:       []byte("test"),
		Connection: Connection(1),
		released:   false,
	}

	// 第一次释放
	msg.Release()
	if !msg.released {
		t.Error("Message should be marked as released")
	}

	// 第二次释放（应该是安全的）
	msg.Release()
	if !msg.released {
		t.Error("Message should still be marked as released")
	}
}
