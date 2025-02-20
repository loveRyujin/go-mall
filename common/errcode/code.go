package errcode

// 此处为公共的错误码, 预留 10000000 ~ 10000099 间的 100 个错误码
var (
	Success            = newError(0, "success")
	ErrServer          = newError(10000000, "服务器内部错误")
	ErrParams          = newError(10000001, "参数错误，请检查")
	ErrNotFound        = newError(10000002, "资源不存在")
	ErrPanic           = newError(10000003, "系统开小差了，请稍后再试")
	ErrToken           = newError(10000004, "token无效")
	ErrForbidden       = newError(10000005, "未授权")
	ErrTooManyRequests = newError(10000006, "请求过多")
)

// 各个业务模块自定义的错误码, 从 10000100 开始, 可以按照不同的业务模块划分不同的号段
// Example:
//var (
//	ErrOrderClosed  = NewError(10000100, "订单已关闭")
//)

var (
	ErrUserInvalid      = newError(10000101, "用户异常")
	ErrUserNameOccupied = newError(10000102, "用户名已被占用")
	ErrUserNotRight     = newError(10000103, "用户名或密码不正确")
)
