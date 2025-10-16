package adapter

import (
	"context"

	"github.com/addcnos/youdu/v2"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// GetGroupListInput represents input for getting group list
type GetGroupListInput struct {
	UserID string `json:"user_id" jsonschema:"description=User ID to get groups for,required"`
}

// GetGroupListOutput represents output for group list
type GetGroupListOutput struct {
	Groups []youdu.GroupItem `json:"groups" jsonschema:"description=List of groups"`
}

// GetGroupList retrieves the list of groups for a user
func (a *Adapter) GetGroupList(ctx context.Context, input GetGroupListInput) (*GetGroupListOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceGroup, permission.ActionRead); err != nil {
		return nil, err
	}

	resp, err := a.client.GetGroupList(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &GetGroupListOutput{
		Groups: resp.GroupList,
	}, nil
}

// GetGroupInfoInput represents input for getting group information
type GetGroupInfoInput struct {
	GroupID string `json:"group_id" jsonschema:"description=Group ID,required"`
}

// GetGroupInfoOutput represents output for group information
type GetGroupInfoOutput struct {
	Group youdu.GroupInfoResponse `json:"group" jsonschema:"description=Group information"`
}

// GetGroupInfo retrieves information about a specific group
func (a *Adapter) GetGroupInfo(ctx context.Context, input GetGroupInfoInput) (*GetGroupInfoOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceGroup, permission.ActionRead); err != nil {
		return nil, err
	}

	resp, err := a.client.GetGroupInfo(ctx, input.GroupID)
	if err != nil {
		return nil, err
	}

	return &GetGroupInfoOutput{
		Group: resp,
	}, nil
}

// CreateGroupInput represents input for creating a group
type CreateGroupInput struct {
	Name string `json:"name" jsonschema:"description=Group name,required"`
}

// CreateGroupOutput represents output for creating a group
type CreateGroupOutput struct {
	GroupID string `json:"group_id" jsonschema:"description=Created group ID"`
}

// CreateGroup creates a new group
func (a *Adapter) CreateGroup(ctx context.Context, input CreateGroupInput) (*CreateGroupOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceGroup, permission.ActionCreate); err != nil {
		return nil, err
	}

	req := youdu.CreateGroupRequest{
		Name: input.Name,
	}

	resp, err := a.client.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return &CreateGroupOutput{
		GroupID: resp.ID,
	}, nil
}

// UpdateGroupInput represents input for updating a group
type UpdateGroupInput struct {
	GroupID string `json:"group_id" jsonschema:"description=Group ID,required"`
	Name    string `json:"name" jsonschema:"description=New group name"`
}

// UpdateGroupOutput represents output for updating a group
type UpdateGroupOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// UpdateGroup updates an existing group
func (a *Adapter) UpdateGroup(ctx context.Context, input UpdateGroupInput) (*UpdateGroupOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceGroup, permission.ActionUpdate); err != nil {
		return nil, err
	}

	req := youdu.UpdateGroupRequest{
		ID:   input.GroupID,
		Name: input.Name,
	}

	_, err := a.client.UpdateGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return &UpdateGroupOutput{
		Success: true,
	}, nil
}

// DeleteGroupInput represents input for deleting a group
type DeleteGroupInput struct {
	GroupID string `json:"group_id" jsonschema:"description=Group ID,required"`
}

// DeleteGroupOutput represents output for deleting a group
type DeleteGroupOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// DeleteGroup deletes a group
func (a *Adapter) DeleteGroup(ctx context.Context, input DeleteGroupInput) (*DeleteGroupOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceGroup, permission.ActionDelete); err != nil {
		return nil, err
	}

	_, err := a.client.DeleteGroup(ctx, input.GroupID)
	if err != nil {
		return nil, err
	}

	return &DeleteGroupOutput{
		Success: true,
	}, nil
}

// AddGroupMemberInput represents input for adding group members
type AddGroupMemberInput struct {
	GroupID string   `json:"group_id" jsonschema:"description=Group ID,required"`
	Members []string `json:"members" jsonschema:"description=List of member user IDs to add,required"`
}

// AddGroupMemberOutput represents output for adding group members
type AddGroupMemberOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// AddGroupMember adds members to a group
func (a *Adapter) AddGroupMember(ctx context.Context, input AddGroupMemberInput) (*AddGroupMemberOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceGroup, permission.ActionUpdate); err != nil {
		return nil, err
	}

	req := youdu.GroupUpdateMemberRequest{
		ID:       input.GroupID,
		UserList: input.Members,
	}

	_, err := a.client.AddGroupMember(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AddGroupMemberOutput{
		Success: true,
	}, nil
}

// DelGroupMemberInput represents input for deleting group members
type DelGroupMemberInput struct {
	GroupID string   `json:"group_id" jsonschema:"description=Group ID,required"`
	Members []string `json:"members" jsonschema:"description=List of member user IDs to remove,required"`
}

// DelGroupMemberOutput represents output for deleting group members
type DelGroupMemberOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// DelGroupMember removes members from a group
func (a *Adapter) DelGroupMember(ctx context.Context, input DelGroupMemberInput) (*DelGroupMemberOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceGroup, permission.ActionUpdate); err != nil {
		return nil, err
	}

	req := youdu.GroupUpdateMemberRequest{
		ID:       input.GroupID,
		UserList: input.Members,
	}

	_, err := a.client.DelGroupMember(ctx, req)
	if err != nil {
		return nil, err
	}

	return &DelGroupMemberOutput{
		Success: true,
	}, nil
}
