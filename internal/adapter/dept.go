package adapter

import (
	"context"

	"github.com/addcnos/youdu/v2"
)

// DeptListInput represents input for getting department list
type DeptListInput struct {
	DeptID int `json:"dept_id" jsonschema:"description=Department ID (0 for root department),default=0"`
}

// DeptListOutput represents output for department list
type DeptListOutput struct {
	Departments []youdu.DeptItem `json:"departments" jsonschema:"description=List of departments"`
}

// GetDeptList retrieves the list of departments
func (a *Adapter) GetDeptList(ctx context.Context, input DeptListInput) (*DeptListOutput, error) {
	resp, err := a.client.GetDeptList(ctx, input.DeptID)
	if err != nil {
		return nil, err
	}

	return &DeptListOutput{
		Departments: resp.DeptList,
	}, nil
}

// DeptUserListInput represents input for getting department user list
type DeptUserListInput struct {
	DeptID int `json:"dept_id" jsonschema:"description=Department ID,required"`
}

// DeptUserListOutput represents output for department user list
type DeptUserListOutput struct {
	Users []youdu.UserItem `json:"users" jsonschema:"description=List of users in department"`
}

// GetDeptUserList retrieves users in a department
func (a *Adapter) GetDeptUserList(ctx context.Context, input DeptUserListInput) (*DeptUserListOutput, error) {
	resp, err := a.client.GetDeptUserList(ctx, input.DeptID)
	if err != nil {
		return nil, err
	}

	return &DeptUserListOutput{
		Users: resp.UserList,
	}, nil
}

// DeptAliasListInput represents input for getting department alias list
type DeptAliasListInput struct {
	DeptID int `json:"dept_id" jsonschema:"description=Department ID (0 for all departments),default=0"`
}

// DeptAliasListOutput represents output for department alias list
type DeptAliasListOutput struct {
	Aliases []youdu.DeptAliasItem `json:"aliases" jsonschema:"description=List of department aliases"`
}

// GetDeptAliasList retrieves department aliases
func (a *Adapter) GetDeptAliasList(ctx context.Context, input DeptAliasListInput) (*DeptAliasListOutput, error) {
	resp, err := a.client.GetDeptAliasList(ctx)
	if err != nil {
		return nil, err
	}

	return &DeptAliasListOutput{
		Aliases: resp.AliasList,
	}, nil
}

// CreateDeptInput represents input for creating a department
type CreateDeptInput struct {
	Name     string `json:"name" jsonschema:"description=Department name,required"`
	ParentID int    `json:"parent_id" jsonschema:"description=Parent department ID,required"`
	SortID   int    `json:"sort_id" jsonschema:"description=Sort order,default=0"`
	Alias    string `json:"alias" jsonschema:"description=Department alias"`
}

// CreateDeptOutput represents output for creating a department
type CreateDeptOutput struct {
	DeptID int `json:"dept_id" jsonschema:"description=Created department ID"`
}

// CreateDept creates a new department
func (a *Adapter) CreateDept(ctx context.Context, input CreateDeptInput) (*CreateDeptOutput, error) {
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

// UpdateDeptInput represents input for updating a department
type UpdateDeptInput struct {
	DeptID   int    `json:"dept_id" jsonschema:"description=Department ID,required"`
	Name     string `json:"name" jsonschema:"description=Department name"`
	ParentID int    `json:"parent_id" jsonschema:"description=Parent department ID"`
	SortID   int    `json:"sort_id" jsonschema:"description=Sort order"`
	Alias    string `json:"alias" jsonschema:"description=Department alias"`
}

// UpdateDeptOutput represents output for updating a department
type UpdateDeptOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// UpdateDept updates an existing department
func (a *Adapter) UpdateDept(ctx context.Context, input UpdateDeptInput) (*UpdateDeptOutput, error) {
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

// DeleteDeptInput represents input for deleting a department
type DeleteDeptInput struct {
	DeptID int `json:"dept_id" jsonschema:"description=Department ID,required"`
}

// DeleteDeptOutput represents output for deleting a department
type DeleteDeptOutput struct {
	Success bool `json:"success" jsonschema:"description=Whether the operation succeeded"`
}

// DeleteDept deletes a department
func (a *Adapter) DeleteDept(ctx context.Context, input DeleteDeptInput) (*DeleteDeptOutput, error) {
	_, err := a.client.DeleteDept(ctx, input.DeptID)
	if err != nil {
		return nil, err
	}

	return &DeleteDeptOutput{
		Success: true,
	}, nil
}
