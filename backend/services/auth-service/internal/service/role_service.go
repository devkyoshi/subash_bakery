package service

import (
	"context"
	"fmt"

	"github.com/yourusername/erp-system/services/auth-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleService struct {
	roleRepo       *repository.RoleRepository
	permissionRepo *repository.PermissionRepository
	userRepo       *repository.UserRepository
}

func NewRoleService(
	roleRepo *repository.RoleRepository,
	permissionRepo *repository.PermissionRepository,
	userRepo *repository.UserRepository,
) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(ctx context.Context, req CreateRoleRequest) (*models.Role, error) {
	// Check if organization ID is valid if provided
	var orgID primitive.ObjectID
	var err error
	if req.OrganizationID != "" {
		orgID, err = primitive.ObjectIDFromHex(req.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("invalid organization ID: %w", err)
		}
	}

	// Parse permission IDs
	permissionIDs := make([]primitive.ObjectID, 0, len(req.PermissionIDs))
	for _, idStr := range req.PermissionIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid permission ID %s: %w", idStr, err)
		}
		permissionIDs = append(permissionIDs, id)
	}

	role := &models.Role{
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Description:    req.Description,
		OrganizationID: orgID,
		Permissions:    permissionIDs,
		IsSystem:       false,
		IsActive:       true,
		Priority:       req.Priority,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// ListRoles retrieves roles, optionally filtered by organization
func (s *RoleService) ListRoles(ctx context.Context, orgIDStr string) ([]*models.Role, error) {
	var orgID *primitive.ObjectID
	if orgIDStr != "" {
		id, err := primitive.ObjectIDFromHex(orgIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid organization ID: %w", err)
		}
		orgID = &id
	}

	return s.roleRepo.FindAll(ctx, orgID)
}

// AssignRole assigns a role to a user
func (s *RoleService) AssignRole(ctx context.Context, userIDStr, roleIDStr string) error {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	roleID, err := primitive.ObjectIDFromHex(roleIDStr)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	// Verify role exists
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Update user role
	user.RoleID = roleID
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// GetRole retrieves a role by ID
func (s *RoleService) GetRole(ctx context.Context, roleIDStr string) (*models.Role, error) {
	roleID, err := primitive.ObjectIDFromHex(roleIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID: %w", err)
	}

	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role not found")
	}

	return role, nil
}

type UpdateRoleRequest struct {
	Name          string   `json:"name"`
	DisplayName   string   `json:"display_name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permissions"`
	Priority      int      `json:"priority"`
	IsActive      *bool    `json:"is_active"`
}

// UpdateRole updates an existing role
func (s *RoleService) UpdateRole(ctx context.Context, roleIDStr string, req UpdateRoleRequest) (*models.Role, error) {
	roleID, err := primitive.ObjectIDFromHex(roleIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID: %w", err)
	}

	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role not found")
	}

	// Update fields if provided
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Priority != 0 {
		role.Priority = req.Priority
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	// Update permissions if provided
	if req.PermissionIDs != nil {
		permissionIDs := make([]primitive.ObjectID, 0, len(req.PermissionIDs))
		for _, idStr := range req.PermissionIDs {
			id, err := primitive.ObjectIDFromHex(idStr)
			if err != nil {
				return nil, fmt.Errorf("invalid permission ID %s: %w", idStr, err)
			}
			permissionIDs = append(permissionIDs, id)
		}
		role.Permissions = permissionIDs
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// DeleteRole deletes a role
func (s *RoleService) DeleteRole(ctx context.Context, roleIDStr string) error {
	roleID, err := primitive.ObjectIDFromHex(roleIDStr)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	// Check if role is system role
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role != nil && role.IsSystem {
		return fmt.Errorf("cannot delete system role")
	}

	return s.roleRepo.Delete(ctx, roleID)
}

// ListPermissions retrieves all permissions
func (s *RoleService) ListPermissions(ctx context.Context) ([]*models.Permission, error) {
	return s.permissionRepo.FindAll(ctx)
}

type CreateRoleRequest struct {
	Name           string   `json:"name" binding:"required"`
	DisplayName    string   `json:"display_name" binding:"required"`
	Description    string   `json:"description"`
	OrganizationID string   `json:"organization_id"`
	PermissionIDs  []string `json:"permissions"`
	Priority       int      `json:"priority"`
}
