package mysql

import "errors"

//mysql包下的错误码
var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("ID无效")
	ErrorInserFaied      = errors.New("用户插入失败")
	ErrorVideoNotExist   = errors.New("视频不存在")
)
