package adapter

import (
	"context"

	"github.com/addcnos/youdu/v2"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// GetUserInput represents input for getting user information
type GetUserInput struct {
	UserID string `json:"user_id" jsonschema:"description=User ID,required"`
}

// GetUserOutput represents output for user information
type GetUserOutput struct {
	User youdu.UserResponse `json:"user" jsonschema:"description=User information"`
}

// GetUser retrieves user information
func (a *Adapter) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceUser, permission.ActionRead, input.UserID); err != nil {
		return nil, err
	}

	resp, err := a.client.GetUser(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{
		User: resp,
	}, nil
}

// CreateUserInput represents input for creating a user
type CreateUserInput struct {
	UserID   string `json:"user_id" jsonschema:"description=User ID,required"`
	Name     string `json:"name" jsonschema:"description=User name,required"`
	Gender   int    `json:"gender" jsonschema:"description=Gender (0:Unknown 1:Male 2:Female),default=0"`
	Mobile   string `json:"mobile" jsonschema:"description=Mobile phone number"`
	Phone    string `json:"phone" jsonschema:"description=Phone number"`
	Email    string `json:"email" jsonschema:"description=Email address"`
	DeptID   int    `json:"dept_id" jsonschema:"description=Department ID,required"`
	Password string `json:"password" jsonschema:"description=User password"`
}

// CreateUserOutput represents output for creating a user
type CreateUserOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// CreateUser creates a new user
func (a *Adapter) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceUser, permission.ActionCreate); err != nil {
		return nil, err
	}

	req := youdu.CreateUserRequest{
		UserID: input.UserID,
		Name:   input.Name,
		Gender: input.Gender,
		Mobile: input.Mobile,
		Phone:  input.Phone,
		Email:  input.Email,
		Dept:   []int{input.DeptID},
	}

	_, err := a.client.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		Success: true,
	}, nil
}

// UpdateUserInput represents input for updating a user
type UpdateUserInput struct {
	UserID string `json:"user_id" jsonschema:"description=User ID,required"`
	Name   string `json:"name" jsonschema:"description=User name"`
	Gender int    `json:"gender" jsonschema:"description=Gender (0:Unknown 1:Male 2:Female)"`
	Mobile string `json:"mobile" jsonschema:"description=Mobile phone number"`
	Phone  string `json:"phone" jsonschema:"description=Phone number"`
	Email  string `json:"email" jsonschema:"description=Email address"`
}

// UpdateUserOutput represents output for updating a user
type UpdateUserOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// UpdateUser updates an existing user
func (a *Adapter) UpdateUser(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceUser, permission.ActionUpdate, input.UserID); err != nil {
		return nil, err
	}

	req := youdu.UpdateUserRequest{
		UserID: input.UserID,
		Name:   input.Name,
		Gender: input.Gender,
		Mobile: input.Mobile,
		Phone:  input.Phone,
		Email:  input.Email,
	}

	_, err := a.client.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &UpdateUserOutput{
		Success: true,
	}, nil
}

// DeleteUserInput represents input for deleting a user
type DeleteUserInput struct {
	UserID string `json:"user_id" jsonschema:"description=User ID,required"`
}

// DeleteUserOutput represents output for deleting a user
type DeleteUserOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// DeleteUser deletes a user
func (a *Adapter) DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceUser, permission.ActionDelete, input.UserID); err != nil {
		return nil, err
	}

	_, err := a.client.DeleteUser(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &DeleteUserOutput{
		Success: true,
	}, nil
}
