package errors

import "fmt"

// AppError 应用错误
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New 创建新错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	ErrInvalidParams      = New(400, "invalid parameters")
	ErrUnauthorized       = New(401, "unauthorized")
	ErrForbidden          = New(403, "forbidden")
	ErrNotFound           = New(404, "resource not found")
	ErrInternalServer     = New(500, "internal server error")
	ErrDatabaseError      = New(500, "database error")
	ErrUserNotFound       = New(404, "user not found")
	ErrUserExists         = New(400, "user already exists")
	ErrInvalidPassword    = New(401, "invalid password")
	ErrKBNotFound         = New(404, "knowledge base not found")
	ErrDocumentNotFound   = New(404, "document not found")
	ErrSessionNotFound    = New(404, "session not found")
	ErrInvalidToken       = New(401, "invalid token")
	ErrTokenExpired       = New(401, "token expired")
	ErrFileUploadFailed   = New(500, "file upload failed")
	ErrFileParseFailed    = New(500, "file parse failed")
	ErrVectorizeFailed    = New(500, "vectorize failed")
	ErrGRPCCallFailed     = New(500, "grpc call failed")
)

