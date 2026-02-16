package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/erp-system/services/auth-service/config"
	"github.com/yourusername/erp-system/services/auth-service/internal/repository"
	"github.com/yourusername/erp-system/shared/models"
	"github.com/yourusername/erp-system/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	userRepo       *repository.UserRepository
	sessionRepo    *repository.SessionRepository
	roleRepo       *repository.RoleRepository
	permissionRepo *repository.PermissionRepository
	deviceRepo     *repository.DeviceRepository
	jwtManager     *utils.JWTManager
	config         *config.Config
	oauthConfig    *oauth2.Config
}

func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	roleRepo *repository.RoleRepository,
	permissionRepo *repository.PermissionRepository,
	deviceRepo *repository.DeviceRepository,
	jwtManager *utils.JWTManager,
	cfg *config.Config,
) *AuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		deviceRepo:     deviceRepo,
		jwtManager:     jwtManager,
		config:         cfg,
		oauthConfig:    oauthConfig,
	}
}

type RegisterRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	Phone          string `json:"phone"`
	OrganizationID string `json:"organization_id"`
	MACAddress     string `json:"mac_address"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RoleResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	IsSystem    bool     `json:"is_system"`
}

type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *models.User  `json:"user"`
	Role         *RoleResponse `json:"role,omitempty"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Resolve organization ID
	var orgID primitive.ObjectID
	if req.OrganizationID != "" {
		// Explicitly provided org ID (admin dashboard flow)
		orgID, err = primitive.ObjectIDFromHex(req.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("invalid organization ID: %w", err)
		}
	} else if req.MACAddress != "" {
		// Look up organization from registered device (mobile app flow)
		device, err := s.deviceRepo.FindByMACAddress(ctx, normalizeMACAddress(req.MACAddress))
		if err != nil {
			return nil, fmt.Errorf("failed to look up device: %w", err)
		}
		if device == nil {
			return nil, fmt.Errorf("device not registered: please contact your administrator to register this device")
		}
		orgID = device.OrganizationID
	}

	// Create user
	user := &models.User{
		OrganizationID:  orgID,
		Email:           req.Email,
		Password:        hashedPassword,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		IsActive:        true,
		IsEmailVerified: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateAuthResponse(ctx, user, "", "")
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Generate tokens
	return s.generateAuthResponse(ctx, user, userAgent, ipAddress)
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	userIDStr, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Find session
	session, err := s.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Delete expired session
		s.sessionRepo.DeleteByRefreshToken(ctx, refreshToken)
		return nil, fmt.Errorf("session expired")
	}

	// Get user
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new access token
	userClaims := models.UserClaims{
		UserID:         user.ID.Hex(),
		OrganizationID: user.OrganizationID.Hex(),
		Email:          user.Email,
		RoleID:         user.RoleID.Hex(),
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(userClaims)
	if err != nil {
		return nil, err
	}

	// Prepare role response
	roleResp, err := s.getRoleResponse(ctx, user.RoleID)
	if err != nil {
		// Log error but proceed
		fmt.Printf("Failed to fetch role data: %v\n", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		Role:         roleResp,
	}, nil
}

// Logout invalidates a user's session
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.sessionRepo.DeleteByRefreshToken(ctx, refreshToken)
}

// GetGoogleLoginURL returns the Google OAuth login URL
func (s *AuthService) GetGoogleLoginURL() string {
	state := uuid.New().String()
	return s.oauthConfig.AuthCodeURL(state)
}

// GoogleLogin authenticates a user via Google OAuth
func (s *AuthService) GoogleLogin(ctx context.Context, code, userAgent, ipAddress string) (*AuthResponse, error) {
	// Exchange code for token
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Google
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Find or create user
	user, err := s.userRepo.FindByGoogleID(ctx, googleUser.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Check if email exists
		user, err = s.userRepo.FindByEmail(ctx, googleUser.Email)
		if err != nil {
			return nil, err
		}

		if user == nil {
			// Create new user
			user = &models.User{
				Email:           googleUser.Email,
				FirstName:       googleUser.GivenName,
				LastName:        googleUser.FamilyName,
				Avatar:          googleUser.Picture,
				GoogleID:        googleUser.ID,
				IsActive:        true,
				IsEmailVerified: googleUser.VerifiedEmail,
			}

			if err := s.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		} else {
			// Link Google account to existing user
			user.GoogleID = googleUser.ID
			if err := s.userRepo.Update(ctx, user); err != nil {
				return nil, err
			}
		}
	}

	// Update last login
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	return s.generateAuthResponse(ctx, user, userAgent, ipAddress)
}

// GetMe returns the current user's profile with role
func (s *AuthService) GetMe(ctx context.Context, userID string) (*AuthResponse, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	roleResp, err := s.getRoleResponse(ctx, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &AuthResponse{
		User: user,
		Role: roleResp,
	}, nil
}

// GetUser returns a user's profile by ID
func (s *AuthService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// Helper to fetch role and permissions
func (s *AuthService) getRoleResponse(ctx context.Context, roleID primitive.ObjectID) (*RoleResponse, error) {
	if roleID.IsZero() {
		return nil, nil
	}

	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil || role == nil {
		return nil, err
	}

	perms, err := s.permissionRepo.FindByIDs(ctx, role.Permissions)
	if err != nil {
		return nil, err
	}

	permNames := make([]string, len(perms))
	for i, p := range perms {
		permNames[i] = p.Name
	}

	return &RoleResponse{
		ID:          role.ID.Hex(),
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		Permissions: permNames,
		IsSystem:    role.IsSystem,
	}, nil
}

// generateAuthResponse generates access and refresh tokens
func (s *AuthService) generateAuthResponse(ctx context.Context, user *models.User, userAgent, ipAddress string) (*AuthResponse, error) {
	// Generate access token
	userClaims := models.UserClaims{
		UserID:         user.ID.Hex(),
		OrganizationID: user.OrganizationID.Hex(),
		Email:          user.Email,
		RoleID:         user.RoleID.Hex(),
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(userClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session
	session := &models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(s.config.RefreshTokenExpiry),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Fetch Role Data
	roleResp, err := s.getRoleResponse(ctx, user.RoleID)
	if err != nil {
		// Log but don't fail, maybe they just don't have a role yet
		fmt.Printf("Warning: Failed to fetch role for user %s: %v\n", user.ID.Hex(), err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		Role:         roleResp,
	}, nil
}

// UpdateUserOrganization updates a user's organization
func (s *AuthService) UpdateUserOrganization(ctx context.Context, userID, orgID string) (*models.User, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	orgObjID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization ID: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, userObjID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	user.OrganizationID = orgObjID
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers retrieves users with pagination and filtering
func (s *AuthService) ListUsers(ctx context.Context, orgID string, filters map[string]interface{}, page, limit int) ([]*models.User, int64, error) {
	orgObjID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid organization ID: %w", err)
	}

	users, total, err := s.userRepo.FindAll(ctx, orgObjID, filters, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Collect unique role IDs
	roleIDs := make([]primitive.ObjectID, 0)
	roleIDMap := make(map[primitive.ObjectID]bool)
	for _, user := range users {
		if !user.RoleID.IsZero() && !roleIDMap[user.RoleID] {
			roleIDs = append(roleIDs, user.RoleID)
			roleIDMap[user.RoleID] = true
		}
	}

	// Fetch roles
	roles, err := s.roleRepo.FindByIDs(ctx, roleIDs)
	if err != nil {
		// Log error but proceed with users
		fmt.Printf("Failed to fetch roles: %v\n", err)
	} else {
		// Create a map for faster lookup
		rolesMap := make(map[primitive.ObjectID]*models.Role)
		for _, role := range roles {
			rolesMap[role.ID] = role
		}

		// Assign roles to users
		for _, user := range users {
			if role, ok := rolesMap[user.RoleID]; ok {
				user.Role = role
			}
		}
	}

	return users, total, nil
}

// CreateUser creates a new user (admin function)
func (s *AuthService) CreateUser(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Parse organization ID if provided
	var orgID primitive.ObjectID
	if req.OrganizationID != "" {
		orgID, err = primitive.ObjectIDFromHex(req.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("invalid organization ID: %w", err)
		}
	}

	// Create user
	user := &models.User{
		OrganizationID:  orgID,
		Email:           req.Email,
		Password:        hashedPassword,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		IsActive:        true,
		IsEmailVerified: false, // Admin created users might need verification or auto-verified
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	RoleID    string `json:"role_id"`
	IsActive  *bool  `json:"is_active"`
}

// UpdateUser updates an existing user
func (s *AuthService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*models.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.RoleID != "" {
		roleID, err := primitive.ObjectIDFromHex(req.RoleID)
		if err == nil {
			user.RoleID = roleID
		}
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *AuthService) DeleteUser(ctx context.Context, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	return s.userRepo.Delete(ctx, id)
}

// normalizeMACAddress normalizes a MAC address to uppercase colon-separated format
func normalizeMACAddress(mac string) string {
	mac = strings.TrimSpace(mac)
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, ".", "")
	mac = strings.ToUpper(mac)

	if len(mac) != 12 {
		return mac // Return as-is; device repo will handle no-match
	}

	parts := make([]string, 6)
	for i := 0; i < 6; i++ {
		parts[i] = mac[i*2 : i*2+2]
	}
	return strings.Join(parts, ":")
}
