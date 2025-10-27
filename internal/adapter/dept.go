package adapter

import (
	"context"
	"fmt"

	"github.com/addcnos/youdu/v2"
	"github.com/yourusername/youdu-app-mcp/internal/permission"
)

// DeptListInput 获取部门列表的输入参数
type DeptListInput struct {
	DeptID int `json:"dept_id" jsonschema:"description=部门ID（0表示根部门）,default=0"`
}

// DeptListOutput 部门列表的输出结果
type DeptListOutput struct {
	Departments []youdu.DeptItem `json:"departments" jsonschema:"description=部门列表"`
}

// GetDeptList 获取部门列表
func (a *Adapter) GetDeptList(ctx context.Context, input DeptListInput) (*DeptListOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceDept, permission.ActionRead, fmt.Sprintf("%d", input.DeptID)); err != nil {
		return nil, err
	}

	resp, err := a.client.GetDeptList(ctx, input.DeptID)
	if err != nil {
		return nil, err
	}

	return &DeptListOutput{
		Departments: resp.DeptList,
	}, nil
}

// DeptUserListInput 获取部门用户列表的输入参数
type DeptUserListInput struct {
	DeptID int `json:"dept_id" jsonschema:"description=部门ID,required"`
}

// DeptUserListOutput 部门用户列表的输出结果
type DeptUserListOutput struct {
	Users []youdu.UserItem `json:"users" jsonschema:"description=部门用户列表"`
}

// GetDeptUserList 获取部门中的用户列表
func (a *Adapter) GetDeptUserList(ctx context.Context, input DeptUserListInput) (*DeptUserListOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceDept, permission.ActionRead, fmt.Sprintf("%d", input.DeptID)); err != nil {
		return nil, err
	}

	resp, err := a.client.GetDeptUserList(ctx, input.DeptID)
	if err != nil {
		return nil, err
	}

	return &DeptUserListOutput{
		Users: resp.UserList,
	}, nil
}

// DeptAliasListInput 获取部门别名列表的输入参数
type DeptAliasListInput struct {
	DeptID int `json:"dept_id" jsonschema:"description=部门ID（0表示所有部门）,default=0"`
}

// DeptAliasListOutput 部门别名列表的输出结果
type DeptAliasListOutput struct {
	Aliases []youdu.DeptAliasItem `json:"aliases" jsonschema:"description=部门别名列表"`
}

// GetDeptAliasList 获取部门别名列表
func (a *Adapter) GetDeptAliasList(ctx context.Context, input DeptAliasListInput) (*DeptAliasListOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceDept, permission.ActionRead); err != nil {
		return nil, err
	}

	resp, err := a.client.GetDeptAliasList(ctx)
	if err != nil {
		return nil, err
	}

	return &DeptAliasListOutput{
		Aliases: resp.AliasList,
	}, nil
}

// CreateDeptInput 创建部门的输入参数
type CreateDeptInput struct {
	Name     string `json:"name" jsonschema:"description=部门名称,required"`
	ParentID int    `json:"parent_id" jsonschema:"description=父部门ID,required"`
	SortID   int    `json:"sort_id" jsonschema:"description=排序顺序,default=0"`
	Alias    string `json:"alias" jsonschema:"description=部门别名"`
}

// CreateDeptOutput 创建部门的输出结果
type CreateDeptOutput struct {
	DeptID int `json:"dept_id" jsonschema:"description=创建的部门ID"`
}

// CreateDept 创建新部门
func (a *Adapter) CreateDept(ctx context.Context, input CreateDeptInput) (*CreateDeptOutput, error) {
	// 权限检查
	if err := a.checkPermission(permission.ResourceDept, permission.ActionCreate); err != nil {
		return nil, err
	}

	req := youdu.CreateDeptRequest{
		Name:     input.Name,
		ParentID: input.ParentID,
		SortID:   input.SortID,
		Alias:    input.Alias,
	}

	resp, err := a.client.CreateDept(ctx, req)
	if err != nil {
		return nil, err
	}

	return &CreateDeptOutput{
		DeptID: resp.ID,
	}, nil
}

// UpdateDeptInput 更新部门的输入参数
type UpdateDeptInput struct {
	DeptID   int    `json:"dept_id" jsonschema:"description=部门ID,required"`
	Name     string `json:"name" jsonschema:"description=部门名称"`
	ParentID int    `json:"parent_id" jsonschema:"description=父部门ID"`
	SortID   int    `json:"sort_id" jsonschema:"description=排序顺序"`
	Alias    string `json:"alias" jsonschema:"description=部门别名"`
}

// UpdateDeptOutput 更新部门的输出结果
type UpdateDeptOutput struct {
	Success bool `json:"success" jsonschema:"description=操作是否成功"`
}

// UpdateDept 更新现有部门
func (a *Adapter) UpdateDept(ctx context.Context, input UpdateDeptInput) (*UpdateDeptOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceDept, permission.ActionUpdate, fmt.Sprintf("%d", input.DeptID)); err != nil {
		return nil, err
	}

	req := youdu.UpdateDeptRequest{
		ID:       input.DeptID,
		Name:     input.Name,
		ParentID: input.ParentID,
		SortID:   input.SortID,
		Alias:    input.Alias,
	}

	_, err := a.client.UpdateDept(ctx, req)
	if err != nil {
		return nil, err
	}

	return &UpdateDeptOutput{
		Success: true,
	}, nil
}

// DeleteDeptInput 删除部门的输入参数
type DeleteDeptInput struct {
	DeptID int `json:"dept_id" jsonschema:"description=部门ID,required"`
}

// DeleteDeptOutput 删除部门的输出结果
type DeleteDeptOutput struct {
	Success bool `json:"success" jsonschema:"description=操作是否成功"`
}

// DeleteDept 删除部门
func (a *Adapter) DeleteDept(ctx context.Context, input DeleteDeptInput) (*DeleteDeptOutput, error) {
	// 权限检查（包含行级权限）
	if err := a.checkPermissionWithID(permission.ResourceDept, permission.ActionDelete, fmt.Sprintf("%d", input.DeptID)); err != nil {
		return nil, err
	}

	_, err := a.client.DeleteDept(ctx, input.DeptID)
	if err != nil {
		return nil, err
	}

	return &DeleteDeptOutput{
		Success: true,
	}, nil
}
