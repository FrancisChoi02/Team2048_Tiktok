package handler

type StatusCode int32

//自增状态码常量
const (
	CodeSuccess StatusCode = 400 + iota
	CodeUserNotLogin
	CodeTokenInvalid
	CodeTokenExpired
	CodeServerBusy
	CodeInvalidParam
	CodeUserExisted
	CodeUserNotExisted
)

const (
	CodeSignUpSuccess = 0
)

//使用哈希表装载 状态码 对应的提示语句
var codeMsgMap = map[StatusCode]string{
	CodeSuccess:        "success",
	CodeUserNotLogin:   "用户未登录",
	CodeTokenInvalid:   "Token无效",
	CodeTokenExpired:   "Token已过期",
	CodeServerBusy:     "服务繁忙",
	CodeInvalidParam:   "参数不合法",
	CodeUserExisted:    "用户名已存在",
	CodeUserNotExisted: "用户不存在",
	CodeSignUpSuccess:  "注册成功",
}

func (c StatusCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy] //非定义范围内的状态码，一律返回服务繁忙
	}
	return msg
}
