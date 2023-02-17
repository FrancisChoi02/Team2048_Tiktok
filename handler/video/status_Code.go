package video

type StatusCode int32

//自增状态码常量
const (
	CodePublishFailed StatusCode = 400 + iota
	CodeServerBusy
	CodeUserIdError
	CodeVideoError
)

const (
	CodeSuccess = 0
)

//使用哈希表装载 状态码 对应的提示语句
var codeMsgMap = map[StatusCode]string{
	CodeSuccess:       "响应成功",
	CodePublishFailed: "视频发布失败",
	CodeServerBusy:    "服务器繁忙",
	CodeUserIdError:   "用户ID解析错误",
	CodeVideoError:    "视频错误",
}

func (c StatusCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy] //非定义范围内的状态码，一律返回服务繁忙
	}
	return msg
}
