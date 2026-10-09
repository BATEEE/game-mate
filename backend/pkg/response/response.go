package response

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Success[T any](data T, message string) APIResponse[T] {
	return APIResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, err error) APIResponse[any] {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	return APIResponse[any]{
		Success: false,
		Message: message,
		Error:   errStr,
	}
}
