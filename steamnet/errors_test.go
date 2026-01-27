package steamnet

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	err := &Error{
		Code:    ErrCodeInvalidConnection,
		Message: "test error",
	}

	expected := "steamnet: test error (code: 1)"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %v, want %v", got, expected)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		code int
	}{
		{"InvalidConnection", ErrInvalidConnection, ErrCodeInvalidConnection},
		{"InvalidSocket", ErrInvalidSocket, ErrCodeInvalidSocket},
		{"ConnectionFailed", ErrConnectionFailed, ErrCodeConnectionFailed},
		{"NotConnected", ErrNotConnected, ErrCodeNotConnected},
		{"InvalidIdentity", ErrInvalidIdentity, ErrCodeInvalidIdentity},
		{"AuthFailed", ErrAuthFailed, ErrCodeAuthFailed},
		{"SendFailed", ErrSendFailed, ErrCodeSendFailed},
		{"ReceiveFailed", ErrReceiveFailed, ErrCodeReceiveFailed},
		{"InvalidPollGroup", ErrInvalidPollGroup, ErrCodeInvalidPollGroup},
		{"InvalidMessage", ErrInvalidMessage, ErrCodeInvalidMessage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Code = %v, want %v", tt.err.Code, tt.code)
			}
		})
	}
}

func TestIsInvalidConnection(t *testing.T) {
	if !IsInvalidConnection(ErrInvalidConnection) {
		t.Error("IsInvalidConnection(ErrInvalidConnection) = false, want true")
	}

	if IsInvalidConnection(ErrInvalidSocket) {
		t.Error("IsInvalidConnection(ErrInvalidSocket) = true, want false")
	}

	if IsInvalidConnection(errors.New("other error")) {
		t.Error("IsInvalidConnection(other error) = true, want false")
	}
}

func TestIsInvalidSocket(t *testing.T) {
	if !IsInvalidSocket(ErrInvalidSocket) {
		t.Error("IsInvalidSocket(ErrInvalidSocket) = false, want true")
	}

	if IsInvalidSocket(ErrInvalidConnection) {
		t.Error("IsInvalidSocket(ErrInvalidConnection) = true, want false")
	}
}

func TestIsConnectionFailed(t *testing.T) {
	if !IsConnectionFailed(ErrConnectionFailed) {
		t.Error("IsConnectionFailed(ErrConnectionFailed) = false, want true")
	}

	if IsConnectionFailed(ErrInvalidConnection) {
		t.Error("IsConnectionFailed(ErrInvalidConnection) = true, want false")
	}
}

func TestNewError(t *testing.T) {
	code := 999
	message := "custom error"
	err := NewError(code, message)

	if err.Code != code {
		t.Errorf("Code = %v, want %v", err.Code, code)
	}

	if err.Message != message {
		t.Errorf("Message = %v, want %v", err.Message, message)
	}
}

func TestWrapError(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := WrapError(baseErr, "wrapped")

	if wrappedErr == nil {
		t.Fatal("WrapError() = nil, want error")
	}

	if !errors.Is(wrappedErr, baseErr) {
		t.Error("WrapError() should wrap the base error")
	}

	// Test with nil error
	if WrapError(nil, "message") != nil {
		t.Error("WrapError(nil) should return nil")
	}
}
