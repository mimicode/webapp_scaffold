package errno

import (
	"fmt"
	"go.uber.org/zap"
)

type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Errno) Error() string {
	return fmt.Sprintf("code:%d Message:%s", e.Code, e.Message)
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Message
	}
	zap.L().Error("响应错误信息被拦截", zap.Error(err))
	return InternalServerError.Code, InternalServerError.Message
}

func NewErrNo(code int, message string) Errno {
	return Errno{
		Code:    code,
		Message: message,
	}
}
