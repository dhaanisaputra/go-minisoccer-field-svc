package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSQLError            = errors.New("database server failed to execute query")
	ErrToManyRequests      = errors.New("too many requests")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("invalid token")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidUploadFile   = errors.New("invalid upload file")
	ErrSizeTooBig          = errors.New("size too big")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrToManyRequests,
	ErrUnauthorized,
	ErrInvalidToken,
	ErrForbidden,
}
