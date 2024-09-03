package app_error

type AppErrorCode uint32

const (
	Internal         AppErrorCode = 0
	InvalidData      AppErrorCode = 1
	IllegalOperation AppErrorCode = 2
	ResourceNotFound AppErrorCode = 3
	Unauthorized     AppErrorCode = 4
)

type AppError struct {
	Code  AppErrorCode
	Error error
}

func InvalidDataError(err error) *AppError {
	return &AppError{Code: InvalidData, Error: err}
}

func InternalError(err error) *AppError {
	return &AppError{Code: Internal, Error: err}
}

func ResourceNotFoundError(err error) *AppError {
	return &AppError{Code: ResourceNotFound, Error: err}
}

func UnauthorizedError(err error) *AppError {
	return &AppError{Code: Unauthorized, Error: err}
}

func IllegalOperationError(err error) *AppError {
	return &AppError{Code: IllegalOperation, Error: err}
}
