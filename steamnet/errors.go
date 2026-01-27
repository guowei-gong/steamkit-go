package steamnet

import (
	"errors"
	"fmt"
)

// Error 表示 Steam 网络错误
type Error struct {
	Code    int    // Steam 错误码
	Message string // 错误消息
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return fmt.Sprintf("steamnet: %s (code: %d)", e.Message, e.Code)
}

// 预定义错误码
const (
	ErrCodeInvalidConnection  = 1
	ErrCodeInvalidSocket      = 2
	ErrCodeConnectionFailed   = 3
	ErrCodeNotConnected       = 4
	ErrCodeInvalidIdentity    = 5
	ErrCodeAuthFailed         = 6
	ErrCodeSendFailed         = 7
	ErrCodeReceiveFailed      = 8
	ErrCodeInvalidPollGroup   = 9
	ErrCodeInvalidMessage     = 10
)

// 预定义错误
var (
	ErrInvalidConnection  = &Error{Code: ErrCodeInvalidConnection, Message: "invalid connection"}
	ErrInvalidSocket      = &Error{Code: ErrCodeInvalidSocket, Message: "invalid listen socket"}
	ErrConnectionFailed   = &Error{Code: ErrCodeConnectionFailed, Message: "connection failed"}
	ErrNotConnected       = &Error{Code: ErrCodeNotConnected, Message: "not connected"}
	ErrInvalidIdentity    = &Error{Code: ErrCodeInvalidIdentity, Message: "invalid identity"}
	ErrAuthFailed         = &Error{Code: ErrCodeAuthFailed, Message: "authentication failed"}
	ErrSendFailed         = &Error{Code: ErrCodeSendFailed, Message: "send failed"}
	ErrReceiveFailed      = &Error{Code: ErrCodeReceiveFailed, Message: "receive failed"}
	ErrInvalidPollGroup   = &Error{Code: ErrCodeInvalidPollGroup, Message: "invalid poll group"}
	ErrInvalidMessage     = &Error{Code: ErrCodeInvalidMessage, Message: "invalid message"}
)

// IsInvalidConnection 检查是否为无效连接错误
func IsInvalidConnection(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeInvalidConnection
}

// IsInvalidSocket 检查是否为无效套接字错误
func IsInvalidSocket(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeInvalidSocket
}

// IsConnectionFailed 检查是否为连接失败错误
func IsConnectionFailed(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeConnectionFailed
}

// IsNotConnected 检查是否为未连接错误
func IsNotConnected(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeNotConnected
}

// IsInvalidIdentity 检查是否为无效身份错误
func IsInvalidIdentity(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeInvalidIdentity
}

// IsAuthFailed 检查是否为认证失败错误
func IsAuthFailed(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeAuthFailed
}

// IsSendFailed 检查是否为发送失败错误
func IsSendFailed(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeSendFailed
}

// IsReceiveFailed 检查是否为接收失败错误
func IsReceiveFailed(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeReceiveFailed
}

// IsInvalidPollGroup 检查是否为无效轮询组错误
func IsInvalidPollGroup(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeInvalidPollGroup
}

// IsInvalidMessage 检查是否为无效消息错误
func IsInvalidMessage(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == ErrCodeInvalidMessage
}

// NewError 创建新的错误
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// WrapError 包装错误
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
