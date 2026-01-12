package models

type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	HTTPStatus int         `json:"-"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrUnauthorized   = &AppError{"UNAUTHORIZED", "Authentication required", 401, nil}
	ErrForbidden      = &AppError{"FORBIDDEN", "Access denied", 403, nil}
	ErrNotFound       = &AppError{"NOT_FOUND", "Resource not found", 404, nil}
	ErrInvalidInput   = &AppError{"INVALID_INPUT", "Invalid input data", 400, nil}
	ErrConflict       = &AppError{"CONFLICT", "Resource conflict", 409, nil}
	ErrInternalServer = &AppError{"INTERNAL_ERROR", "Internal server error", 500, nil}
)

func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}
