package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/erp-system/services/auth-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Seeder handles database seeding operations
type Seeder struct {
	permissionRepo *repository.PermissionRepository
	roleRepo       *repository.RoleRepository
}

// NewSeeder creates a new seeder instance
func NewSeeder(permissionRepo *repository.PermissionRepository, roleRepo *repository.RoleRepository) *Seeder {
	return &Seeder{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
	}
}

// SeedAll runs all seed operations in the correct order
func (s *Seeder) SeedAll(ctx context.Context) error {
	log.Println("Starting database seeding...")

	// Seed permissions first (roles depend on permissions)
	if err := s.SeedPermissions(ctx); err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	// Seed roles
	if err := s.SeedRoles(ctx); err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// SeedPermissions creates system permissions if they don't exist
func (s *Seeder) SeedPermissions(ctx context.Context) error {
	log.Println("Seeding permissions...")

	// Define system permissions
	systemPermissions := []struct {
		Name        string
		DisplayName string
		Description string
		Resource    string
		Action      string
		Category    string
		Scope       models.PermissionScope
	}{
		{
			Name:        "users.read",
			DisplayName: "Read Users",
			Description: "Can view users",
			Resource:    "users",
			Action:      "read",
			Category:    "User Management",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "users.write",
			DisplayName: "Write Users",
			Description: "Can create/edit users",
			Resource:    "users",
			Action:      "write",
			Category:    "User Management",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "roles.read",
			DisplayName: "Read Roles",
			Description: "Can view roles",
			Resource:    "roles",
			Action:      "read",
			Category:    "User Management",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "roles.write",
			DisplayName: "Write Roles",
			Description: "Can create/edit roles",
			Resource:    "roles",
			Action:      "write",
			Category:    "User Management",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "inventory.read",
			DisplayName: "Read Inventory",
			Description: "Can view inventory",
			Resource:    "inventory",
			Action:      "read",
			Category:    "Inventory",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "inventory.write",
			DisplayName: "Write Inventory",
			Description: "Can manage inventory",
			Resource:    "inventory",
			Action:      "write",
			Category:    "Inventory",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "products.read",
			DisplayName: "Read Products",
			Description: "Can view products",
			Resource:    "products",
			Action:      "read",
			Category:    "Products",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "products.write",
			DisplayName: "Write Products",
			Description: "Can manage products",
			Resource:    "products",
			Action:      "write",
			Category:    "Products",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "orders.read",
			DisplayName: "Read Orders",
			Description: "Can view orders",
			Resource:    "orders",
			Action:      "read",
			Category:    "Orders",
			Scope:       models.PermissionScopeOrganization,
		},
		{
			Name:        "orders.write",
			DisplayName: "Write Orders",
			Description: "Can manage orders",
			Resource:    "orders",
			Action:      "write",
			Category:    "Orders",
			Scope:       models.PermissionScopeOrganization,
		},
	}

	// Check which permissions already exist
	var permissionsToCreate []*models.Permission
	for _, perm := range systemPermissions {
		existing, err := s.permissionRepo.FindByName(ctx, perm.Name)
		if err != nil {
			return fmt.Errorf("failed to check permission %s: %w", perm.Name, err)
		}

		if existing == nil {
			permissionsToCreate = append(permissionsToCreate, &models.Permission{
				BaseModel: models.BaseModel{
					Version: 0,
				},
				Name:        perm.Name,
				DisplayName: perm.DisplayName,
				Description: perm.Description,
				Resource:    perm.Resource,
				Action:      perm.Action,
				Category:    perm.Category,
				Scope:       perm.Scope,
				IsSystem:    true,
			})
		}
	}

	if len(permissionsToCreate) > 0 {
		if err := s.permissionRepo.BulkCreate(ctx, permissionsToCreate); err != nil {
			return fmt.Errorf("failed to create permissions: %w", err)
		}
		log.Printf("Created %d new permissions", len(permissionsToCreate))
	} else {
		log.Println("All permissions already exist, skipping creation")
	}

	return nil
}

// SeedRoles creates system roles if they don't exist
func (s *Seeder) SeedRoles(ctx context.Context) error {
	log.Println("Seeding roles...")

	// Get permission IDs for STANDARD_USER role
	standardUserPermissionIDs := []primitive.ObjectID{}
	permissionNames := []string{"users.read", "inventory.read", "products.read", "orders.read"}

	for _, name := range permissionNames {
		perm, err := s.permissionRepo.FindByName(ctx, name)
		if err != nil {
			return fmt.Errorf("failed to find permission %s: %w", name, err)
		}
		if perm != nil {
			standardUserPermissionIDs = append(standardUserPermissionIDs, perm.ID)
		}
	}

	// Define system roles
	systemRoles := []struct {
		Name        string
		DisplayName string
		Description string
		Permissions []primitive.ObjectID
		IsDefault   bool
	}{
		{
			Name:        "ADMIN",
			DisplayName: "Administrator",
			Description: "System Administrator with full access",
			Permissions: nil, // null permissions means full access
			IsDefault:   false,
		},
		{
			Name:        "STANDARD_USER",
			DisplayName: "Standard User",
			Description: "Standard user with read-only access",
			Permissions: standardUserPermissionIDs,
			IsDefault:   true,
		},
	}

	// Create roles if they don't exist
	zeroOrgID, _ := primitive.ObjectIDFromHex("000000000000000000000000")

	for _, role := range systemRoles {
		existing, err := s.roleRepo.FindByName(ctx, role.Name)
		if err != nil {
			return fmt.Errorf("failed to check role %s: %w", role.Name, err)
		}

		if existing == nil {
			newRole := &models.Role{
				BaseModel: models.BaseModel{
					Version: 0,
				},
				OrganizationID: zeroOrgID,
				Name:           role.Name,
				DisplayName:    role.DisplayName,
				Description:    role.Description,
				Permissions:    role.Permissions,
				IsSystem:       true,
				IsDefault:      role.IsDefault,
				IsActive:       true,
				Priority:       0,
			}

			if err := s.roleRepo.Create(ctx, newRole); err != nil {
				return fmt.Errorf("failed to create role %s: %w", role.Name, err)
			}
			log.Printf("Created role: %s", role.Name)
		} else {
			log.Printf("Role %s already exists, skipping creation", role.Name)
		}
	}

	return nil
}
