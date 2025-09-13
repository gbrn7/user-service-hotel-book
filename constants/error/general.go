package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSQLError            = errors.New("database server failed to execute query")
	ErrTooManyRequests     = errors.New("too many request")
	ErrUnautorized         = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("forbidden")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrTooManyRequests,
	ErrUnautorized,
	ErrInvalidToken,
}
