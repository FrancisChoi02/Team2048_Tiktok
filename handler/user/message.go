package user

import (
	logic "Team2048_Tiktok/logic/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// MessageHandler  发送聊天信息
// @Summary 发送聊天信息
// @Description 发送聊天信息接口
// @Tags 聊天相关接口
// @Accept application/json
// @Produce application/json
// @Param to_user_id query int true "接收方用户 ID"
// @Param action_type query string true "消息操作类型"
// @Param content query string true "消息内容"
// @Security ApiKeyAuth
// @Success 200 {string} string "消息发送成功"
// @Router /douyin/message/action/ [post]
func MessageHandler(c *gin.Context) {
	// 1.获取请求中的参数和消息文本
	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseMessage(c, CodeUserIdError)
		return
	}
	rawToUser := c.Query("to_user_id") //获取发送方
	tmpToUser, err := strconv.Atoi(rawToUser)
	if err != nil {
		ResponseMessage(c, CodeUserIdError)
		return
	}
	toUserId := int64(tmpToUser)

	actionType := c.Query("action_type") //获取消息操作的类型

	fmt.Println("THIS IS actionType", actionType)
	if actionType != "1" {
		zap.L().Error("ActionType invalid, it has to be 1", zap.Error(err))
		ResponseMessage(c, CodeRelationTypeError)
		return
	}

	content := c.Query("content") //获取消息内容

	// 2.处理消息发送逻辑
	if err := logic.SendMessage(userId, toUserId, content); err != nil {
		zap.L().Error("logic.SendMessage() failed", zap.Error(err))
		ResponseMessage(c, CodeRelationTypeError)
		return
	}

	// 3.返回消息成功发送的响应
	ResponseMessage(c, CodeSuccess)
}

// ChatHistoryHandler  获取聊天记录
// @Summary 获取聊天记录
// @Description 获取聊天记录接口
// @Tags 聊天相关接口
// @Accept application/json
// @Produce application/json
// @Param to_user_id query int true "接收方用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.Message "消息发送成功"
// @Router /chat/history [get]
func ChatHistoryHandler(c *gin.Context) {
	// 1.获取请求中的参数
	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseChatHistoryError(c, CodeUserIdError)
		return
	}
	rawToUser := c.Query("to_user_id") //获取发送方
	tmpToUser, err := strconv.Atoi(rawToUser)
	if err != nil {
		ResponseChatHistoryError(c, CodeUserIdError)
		return
	}
	toUserId := int64(tmpToUser)

	// 2.处理聊天记录逻辑
	messageList, err := logic.GetMessageList(userId, toUserId)
	if err != nil {
		zap.L().Error("logic.GetMessageList failer()", zap.Error(err))
		ResponseChatHistoryError(c, CodeRelationTypeError)
		return
	}

	// 3.返回聊天记录列表与响应
	ResponseChatHistory(c, CodeSuccess, messageList)

}
