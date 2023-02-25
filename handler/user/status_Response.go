package user

import (
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseErrorWithMsg 返回带指定字符串的错误消息
func ResponseErrorWithMsg(c *gin.Context, code StatusCode, tmpString string) {
	res := &model.CommonResponse{
		StatusCode: int32(code),
		StatusMsg:  tmpString,
	}

	c.JSON(http.StatusOK, res)
}

// ResponseSignUp 返回注册结果
func ResponseSignUp(c *gin.Context, res model.UserDetailResponse, code StatusCode) {
	res.Code = int32(code)
	res.Msg = code.Msg()
	c.JSON(http.StatusOK, res)
}

// ResponseLogin 返回登录结果
func ResponseLogin(c *gin.Context, res model.UserDetailResponse, code StatusCode) {
	res.Code = int32(code)
	res.Msg = code.Msg()
	c.JSON(http.StatusOK, res)
}

// ResponseInfo 返回用户信息获取结果
func ResponseInfo(c *gin.Context, res model.UserInfoResponse, code StatusCode) {
	res.Code = int32(code)
	res.Msg = code.Msg()
	c.JSON(http.StatusOK, res)
}

// ResponseRelation 返回关系操作的结果
func ResponseRelation(c *gin.Context, code StatusCode) {
	res := &model.CommonResponse{
		StatusCode: int32(code),
		StatusMsg:  code.Msg(),
	}
	c.JSON(http.StatusOK, res)
}

// ResponseRelation 返回用户列表获取错误
func ResponseRelationListError(c *gin.Context, code StatusCode) {
	userList := new([]model.UserResponse)
	res := &model.FollowRelationResponse{
		Code:     int32(code),
		Msg:      code.Msg(),
		UserList: *userList,
	}
	c.JSON(http.StatusOK, res)
}

// ResponseRelation 返回用户列表获取成功
func ResponseRelationListSuccess(c *gin.Context, code StatusCode, userList *[]model.UserResponse) {
	res := &model.FollowRelationResponse{
		Code:     int32(code),
		Msg:      code.Msg(),
		UserList: *userList,
	}
	c.JSON(http.StatusOK, res)
}

// ResponseFriendListSuccess 返回聊天好友列表获取成功
func ResponseFriendListSuccess(c *gin.Context, code StatusCode, friendList *[]model.FriendResponse) {
	res := &model.FriendListResponse{
		Code:       int32(code),
		Msg:        code.Msg(),
		FriendList: *friendList,
	}
	c.JSON(http.StatusOK, res)
}

// ResponseFriendListError 返回聊天好友列表获取失败
func ResponseFriendListError(c *gin.Context, code StatusCode) {
	friendList := new([]model.FriendResponse)
	res := &model.FriendListResponse{
		Code:       int32(code),
		Msg:        code.Msg(),
		FriendList: *friendList,
	}
	c.JSON(http.StatusOK, res)
}

// ResponseMessage 返回消息操作的响应
func ResponseMessage(c *gin.Context, code StatusCode) {
	res := &model.CommonResponse{
		StatusCode: int32(code),
		StatusMsg:  code.Msg(),
	}
	c.JSON(http.StatusOK, res)
}

// ResponseChatHistory  返回聊天记录响应
func ResponseChatHistory(c *gin.Context, code StatusCode, messageList *[]model.Message) {
	res := &model.ChatHistoryResponse{
		Code:        int32(code),
		Msg:         code.Msg(),
		MessageList: *messageList,
	}
	c.JSON(http.StatusOK, res)
}

// ResponseChatHistoryError 返回聊天记录错误
func ResponseChatHistoryError(c *gin.Context, code StatusCode) {
	messageList := new([]model.Message)
	res := &model.ChatHistoryResponse{
		Code:        int32(code),
		Msg:         code.Msg(),
		MessageList: *messageList,
	}
	c.JSON(http.StatusOK, res)
}
