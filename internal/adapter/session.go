package adapter

import (
	"context"

	"github.com/addcnos/youdu/v2"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// CreateSessionInput represents input for creating a session
type CreateSessionInput struct {
	Title   string   `json:"title" jsonschema:"description=Session title,required"`
	Creator string   `json:"creator" jsonschema:"description=Session creator user ID,required"`
	Members []string `json:"members" jsonschema:"description=List of member user IDs"`
	Type    string   `json:"type" jsonschema:"description=Session type (single/group),required"`
}

// CreateSessionOutput represents output for creating a session
type CreateSessionOutput struct {
	SessionID string `json:"session_id" jsonschema:"description=Created session ID"`
}

// CreateSession creates a new session
func (a *Adapter) CreateSession(ctx context.Context, input CreateSessionInput) (*CreateSessionOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceSession, permission.ActionCreate); err != nil {
		return nil, err
	}

	req := youdu.CreateSessionRequest{
		Title:   input.Title,
		Creator: input.Creator,
		Member:  input.Members,
		Type:    youdu.SessionType(input.Type),
	}

	resp, err := a.client.CreateSession(ctx, req)
	if err != nil {
		return nil, err
	}

	return &CreateSessionOutput{
		SessionID: resp.SessionID,
	}, nil
}

// GetSessionInput represents input for getting session information
type GetSessionInput struct {
	SessionID string `json:"session_id" jsonschema:"description=Session ID,required"`
}

// GetSessionOutput represents output for session information
type GetSessionOutput struct {
	Session youdu.SessionResponse `json:"session" jsonschema:"description=Session information"`
}

// GetSession retrieves session information
func (a *Adapter) GetSession(ctx context.Context, input GetSessionInput) (*GetSessionOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceSession, permission.ActionRead, input.SessionID); err != nil {
		return nil, err
	}

	resp, err := a.client.GetSession(ctx, input.SessionID)
	if err != nil {
		return nil, err
	}

	return &GetSessionOutput{
		Session: resp,
	}, nil
}

// UpdateSessionInput represents input for updating a session
type UpdateSessionInput struct {
	SessionID string `json:"session_id" jsonschema:"description=Session ID,required"`
	Title     string `json:"title" jsonschema:"description=New session title"`
	OpUser    string `json:"op_user" jsonschema:"description=Operator user ID,required"`
}

// UpdateSessionOutput represents output for updating a session
type UpdateSessionOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// UpdateSession updates an existing session
func (a *Adapter) UpdateSession(ctx context.Context, input UpdateSessionInput) (*UpdateSessionOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceSession, permission.ActionUpdate, input.SessionID); err != nil {
		return nil, err
	}

	req := youdu.UpdateSessionRequest{
		SessionID: input.SessionID,
		Title:     input.Title,
		OpUser:    input.OpUser,
	}

	_, err := a.client.UpdateSession(ctx, req)
	if err != nil {
		return nil, err
	}

	return &UpdateSessionOutput{
		Success: true,
	}, nil
}

// SendTextSessionMessageInput represents input for sending text session message
type SendTextSessionMessageInput struct {
	SessionID string `json:"session_id" jsonschema:"description=Session ID,required"`
	Content   string `json:"content" jsonschema:"description=Message content,required"`
	Sender    string `json:"sender" jsonschema:"description=Sender user ID,required"`
}

// SendTextSessionMessageOutput represents output for sending text session message
type SendTextSessionMessageOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the message was sent successfully"`
}

// SendTextSessionMessage sends a text message to a session
func (a *Adapter) SendTextSessionMessage(ctx context.Context, input SendTextSessionMessageInput) (*SendTextSessionMessageOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceSession, permission.ActionUpdate, input.SessionID); err != nil {
		return nil, err
	}

	req := youdu.TextSessionMessageRequest{
		SessionID: input.SessionID,
		Sender:    input.Sender,
		MsgType:   "text",
		Text: youdu.MessageText{
			Content: input.Content,
		},
	}

	_, err := a.client.SendTextSessionMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SendTextSessionMessageOutput{
		Success: true,
	}, nil
}

// SendImageSessionMessageInput represents input for sending image session message
type SendImageSessionMessageInput struct {
	SessionID string `json:"session_id" jsonschema:"description=Session ID,required"`
	MediaID   string `json:"media_id" jsonschema:"description=Media ID of the uploaded image,required"`
	Sender    string `json:"sender" jsonschema:"description=Sender user ID,required"`
}

// SendImageSessionMessageOutput represents output for sending image session message
type SendImageSessionMessageOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the message was sent successfully"`
}

// SendImageSessionMessage sends an image message to a session
func (a *Adapter) SendImageSessionMessage(ctx context.Context, input SendImageSessionMessageInput) (*SendImageSessionMessageOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceSession, permission.ActionUpdate, input.SessionID); err != nil {
		return nil, err
	}

	req := youdu.ImageSessionMessageRequest{
		SessionID: input.SessionID,
		Sender:    input.Sender,
		MsgType:   "image",
		Image: youdu.MessageMedia{
			MediaID: input.MediaID,
		},
	}

	_, err := a.client.SendImageSessionMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SendImageSessionMessageOutput{
		Success: true,
	}, nil
}

// SendFileSessionMessageInput represents input for sending file session message
type SendFileSessionMessageInput struct {
	SessionID string `json:"session_id" jsonschema:"description=Session ID,required"`
	MediaID   string `json:"media_id" jsonschema:"description=Media ID of the uploaded file,required"`
	Sender    string `json:"sender" jsonschema:"description=Sender user ID,required"`
}

// SendFileSessionMessageOutput represents output for sending file session message
type SendFileSessionMessageOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the message was sent successfully"`
}

// SendFileSessionMessage sends a file message to a session
func (a *Adapter) SendFileSessionMessage(ctx context.Context, input SendFileSessionMessageInput) (*SendFileSessionMessageOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceSession, permission.ActionUpdate, input.SessionID); err != nil {
		return nil, err
	}

	req := youdu.FileSessionMessageRequest{
		SessionID: input.SessionID,
		Sender:    input.Sender,
		MsgType:   "file",
		File: youdu.MessageFile{
			MediaID: input.MediaID,
		},
	}

	_, err := a.client.SendFileSessionMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SendFileSessionMessageOutput{
		Success: true,
	}, nil
}
