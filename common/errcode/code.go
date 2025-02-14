package errcode

var (
	Success            = newError(0, "success")
	ErrServer          = newError(10000000, "服务器内部错误")
	ErrParams          = newError(10000001, "参数错误，请检查")
	ErrNotFound        = newError(10000002, "资源不存在")
	ErrPanic           = newError(10000003, "系统开小差了，请稍后再试")
	ErrToken           = newError(10000004, "token无效")
	ErrForbidden       = newError(10000005, "未授权")
	ErrTooManyRequests = newError(10000006, "请求过多")
	ErrUserInvalid     = newError(10000007, "用户无效")
)
