package user

type StatusCode int32

//自增状态码常量
const (
	CodeUserNotLogin StatusCode = 400 + iota
	CodeUserIdError
	CodeTokenInvalid
	CodeTokenExpired
	CodeServerBusy
	CodeInvalidParam
	CodeUserExisted
	CodeUserNotExisted
	CodeRelationTypeError
	CodeRelationActionError
	CodeRelationFollowError
)

const (
	CodeSuccess = 0
)

//使用哈希表装载 状态码 对应的提示语句
var codeMsgMap = map[StatusCode]string{
	CodeSuccess:             "响应成功",
	CodeUserNotLogin:        "用户未登录",
	CodeTokenInvalid:        "Token无效",
	CodeTokenExpired:        "Token已过期",
	CodeServerBusy:          "服务繁忙",
	CodeInvalidParam:        "参数不合法",
	CodeUserExisted:         "用户名已存在",
	CodeUserNotExisted:      "用户不存在",
	CodeUserIdError:         "用户ID解析错误",
	CodeRelationTypeError:   "关系操作参数错误",
	CodeRelationActionError: "关系操作错误",
	CodeRelationFollowError: "关注列表获取错误",
}

func (c StatusCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy] //非定义范围内的状态码，一律返回服务繁忙
	}
	return msg
}
