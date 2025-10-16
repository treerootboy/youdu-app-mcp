package adapter

import (
	"context"
	"fmt"

	"github.com/addcnos/youdu/v2"
)

// SendTextMessageInput represents input for sending text message
type SendTextMessageInput struct {
	ToUser  string `json:"to_user" jsonschema:"description=Target user ID (use pipe | to separate multiple users)"`
	ToDept  string `json:"to_dept" jsonschema:"description=Target department ID (use pipe | to separate multiple departments)"`
	Content string `json:"content" jsonschema:"description=Message content,required"`
}

// SendTextMessageOutput represents output for sending text message
type SendTextMessageOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the message was sent successfully"`
}

// SendTextMessage sends a text message
func (a *Adapter) SendTextMessage(ctx context.Context, input SendTextMessageInput) (*SendTextMessageOutput, error) {
	// 验证输入
	if input.ToUser == "" && input.ToDept == "" {
		return nil, fmt.Errorf("必须指定接收者：to_user 或 to_dept 至少填写一个")
	}
	if input.Content == "" {
		return nil, fmt.Errorf("消息内容不能为空")
	}

	req := youdu.TextMessageRequest{
		ToUser:  input.ToUser,
		ToDept:  input.ToDept,
		MsgType: "text",
		Text: youdu.MessageText{
			Content: input.Content,
		},
	}

	_, err := a.client.SendTextMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("发送文本消息失败: %w\n提示：请检查用户ID(%s)是否正确，以及应用是否有发送消息的权限", err, input.ToUser)
	}

	return &SendTextMessageOutput{
		Success: true,
	}, nil
}

// SendImageMessageInput represents input for sending image message
type SendImageMessageInput struct {
	ToUser  string `json:"to_user" jsonschema:"description=Target user ID (use pipe | to separate multiple users)"`
	ToDept  string `json:"to_dept" jsonschema:"description=Target department ID (use pipe | to separate multiple departments)"`
	MediaID string `json:"media_id" jsonschema:"description=Media ID of the uploaded image,required"`
}

// SendImageMessageOutput represents output for sending image message
type SendImageMessageOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the message was sent successfully"`
}

// SendImageMessage sends an image message
func (a *Adapter) SendImageMessage(ctx context.Context, input SendImageMessageInput) (*SendImageMessageOutput, error) {
	req := youdu.ImageMessageRequest{
		ToUser:  input.ToUser,
		ToDept:  input.ToDept,
		MsgType: "image",
		Image: youdu.MessageMedia{
			MediaID: input.MediaID,
		},
	}

	_, err := a.client.SendImageMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SendImageMessageOutput{
		Success: true,
	}, nil
}

// SendFileMessageInput represents input for sending file message
type SendFileMessageInput struct {
	ToUser  string `json:"to_user" jsonschema:"description=Target user ID (use pipe | to separate multiple users)"`
	ToDept  string `json:"to_dept" jsonschema:"description=Target department ID (use pipe | to separate multiple departments)"`
	MediaID string `json:"media_id" jsonschema:"description=Media ID of the uploaded file,required"`
}

// SendFileMessageOutput represents output for sending file message
type SendFileMessageOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the message was sent successfully"`
}

// SendFileMessage sends a file message
func (a *Adapter) SendFileMessage(ctx context.Context, input SendFileMessageInput) (*SendFileMessageOutput, error) {
	req := youdu.FileMessageRequest{
		ToUser:  input.ToUser,
		ToDept:  input.ToDept,
		MsgType: "file",
		File: youdu.MessageMedia{
			MediaID: input.MediaID,
		},
	}

	_, err := a.client.SendFileMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SendFileMessageOutput{
		Success: true,
	}, nil
}

// SendLinkMessageInput represents input for sending link message
type SendLinkMessageInput struct {
	ToUser string `json:"to_user" jsonschema:"description=Target user ID (use pipe | to separate multiple users)"`
	ToDept string `json:"to_dept" jsonschema:"description=Target department ID (use pipe | to separate multiple departments)"`
	Title  string `json:"title" jsonschema:"description=Link title,required"`
	URL    string `json:"url" jsonschema:"description=Link URL,required"`
	Action int    `json:"action" jsonschema:"description=Action type (0:webview 1:open external browser),default=0"`
}

// SendLinkMessageOutput represents output for sending link message
type SendLinkMessageOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the message was sent successfully"`
}

// SendLinkMessage sends a link message
func (a *Adapter) SendLinkMessage(ctx context.Context, input SendLinkMessageInput) (*SendLinkMessageOutput, error) {
	req := youdu.LinkMessageRequest{
		ToUser:  input.ToUser,
		ToDept:  input.ToDept,
		MsgType: "link",
		Link: youdu.MessageLink{
			Title:  input.Title,
			URL:    input.URL,
			Action: input.Action,
		},
	}

	_, err := a.client.SendLinkMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SendLinkMessageOutput{
		Success: true,
	}, nil
}

// SendSysMessageInput represents input for sending system message
type SendSysMessageInput struct {
	ToUser      string `json:"to_user" jsonschema:"description=Target user ID (use pipe | to separate multiple users)"`
	ToDept      string `json:"to_dept" jsonschema:"description=Target department ID (use pipe | to separate multiple departments)"`
	Title       string `json:"title" jsonschema:"description=System message title,required"`
	Content     string `json:"content" jsonschema:"description=System message content,required"`
	PopDuration int    `json:"pop_duration" jsonschema:"description=Pop window duration in seconds,default=0"`
}

// SendSysMessageOutput represents output for sending system message
type SendSysMessageOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the message was sent successfully"`
}

// SendSysMessage sends a system message
func (a *Adapter) SendSysMessage(ctx context.Context, input SendSysMessageInput) (*SendSysMessageOutput, error) {
	req := youdu.MessageSysMessageRequest{
		ToUser:  input.ToUser,
		ToDept:  input.ToDept,
		MsgType: "sysMsg",
		SysMsg: youdu.MessageSysMessageSysMsg{
			Title:       input.Title,
			PopDuration: input.PopDuration,
			Msg: []youdu.MessageSysMessageSysMsgMsg{
				{
					Text: youdu.MessageText{
						Content: input.Content,
					},
				},
			},
		},
	}

	_, err := a.client.SendSysMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SendSysMessageOutput{
		Success: true,
	}, nil
}
