package errno

// 10 - 系统错误 00  - 01
// 20 - 用户输入 00  - 01

var (
	OK                  = Errno{Code: 0, Message: "ok"}
	InternalServerError = Errno{Code: 100001, Message: "Internal server error."}
)
