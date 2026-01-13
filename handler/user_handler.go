package handler

import (
	"context"
	"inventory-api/dtos"
	"inventory-api/middleware"
	"inventory-api/services"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(api huma.API) {
	// Auth routes (public)
	huma.Register(api, huma.Operation{
		OperationID: "register-user",
		Method:      http.MethodPost,
		Path:        "/auth/register",
		Summary:     "Register a new user",
		Tags:        []string{"Authentication"},
	}, h.Register)

	huma.Register(api, huma.Operation{
		OperationID: "login-user",
		Method:      http.MethodPost,
		Path:        "/auth/login",
		Summary:     "Login user",
		Tags:        []string{"Authentication"},
	}, h.Login)

	// User routes (protected - will need middleware)
	huma.Register(api, huma.Operation{
		OperationID: "get-profile",
		Method:      http.MethodGet,
		Path:        "/users/profile",
		Summary:     "Get current user profile",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.GetProfile)

	huma.Register(api, huma.Operation{
		OperationID: "change-password",
		Method:      http.MethodPost,
		Path:        "/users/change-password",
		Summary:     "Change password",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.ChangePassword)

	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      http.MethodGet,
		Path:        "/users",
		Summary:     "Get all users (admin only)",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.GetAllUsers)

	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "Get user by ID (admin only)",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.GetUserByID)

	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPut,
		Path:        "/users/{id}",
		Summary:     "Update user (admin only)",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.UpdateUser)

	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Method:      http.MethodDelete,
		Path:        "/users/{id}",
		Summary:     "Delete user (admin only)",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.DeleteUser)
}

func (h *UserHandler) Register(ctx context.Context, input *dtos.RegisterRequest) (*dtos.SingleUserResponse, error) {
	user, err := h.userService.Register(
		input.Body.Username,
		input.Body.Password,
		input.Body.Email,
		input.Body.Phone,
		input.Body.Role,
	)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	return &dtos.SingleUserResponse{Body: dtos.ToUserResponse(user)}, nil
}

func (h *UserHandler) Login(ctx context.Context, input *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	token, user, err := h.userService.Login(input.Body.Username, input.Body.Password)
	if err != nil {
		return nil, huma.Error401Unauthorized(err.Error())
	}

	resp := &dtos.LoginResponse{}
	resp.Body.Token = token
	resp.Body.User = dtos.ToUserResponse(user)
	return resp, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, input *struct{}) (*dtos.SingleUserResponse, error) {
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("User not authenticated")
	}

	user, err := h.userService.GetUserByID(auth.UserID)
	if err != nil {
		return nil, huma.Error404NotFound("User not found")
	}

	return &dtos.SingleUserResponse{Body: dtos.ToUserResponse(user)}, nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, input *dtos.ChangePasswordRequest) (*struct {
	Body struct {
		Message string `json:"message"`
	}
}, error) {
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("User not authenticated")
	}

	err := h.userService.ChangePassword(auth.UserID, input.Body.OldPassword, input.Body.NewPassword)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	resp := &struct {
		Body struct {
			Message string `json:"message"`
		}
	}{}
	resp.Body.Message = "Password changed successfully"
	return resp, nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, input *dtos.PaginationQuery) (*dtos.UserListResponse, error) {
	// Check if user is admin
	if !middleware.IsAdmin(ctx) {
		return nil, huma.Error403Forbidden("Only admins can list all users")
	}

	users, err := h.userService.GetAllUsers(input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	resp := &dtos.UserListResponse{}
	resp.Body.Users = dtos.ToUserResponseList(users)
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset
	return resp, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, input *dtos.IDParam) (*dtos.SingleUserResponse, error) {
	// Check if user is admin or accessing their own profile
	if !middleware.IsOwnerOrAdmin(ctx, input.ID) {
		return nil, huma.Error403Forbidden("Access denied")
	}

	user, err := h.userService.GetUserByID(input.ID)
	if err != nil {
		return nil, huma.Error404NotFound("User not found")
	}

	return &dtos.SingleUserResponse{Body: dtos.ToUserResponse(user)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *dtos.UpdateUserRequest) (*dtos.SingleUserResponse, error) {
	auth := middleware.GetAuthContext(ctx)
	if auth == nil {
		return nil, huma.Error401Unauthorized("Authentication required")
	}

	// Check permissions
	isOwner := auth.UserID == input.ID
	isAdmin := middleware.IsAdmin(ctx)

	if !isOwner && !isAdmin {
		return nil, huma.Error403Forbidden("Access denied")
	}

	// Only admins can change roles
	if input.Body.Role != nil && !isAdmin {
		return nil, huma.Error403Forbidden("Only admins can change user roles")
	}

	email := ""
	phone := ""
	role := ""
	if input.Body.Email != nil {
		email = *input.Body.Email
	}
	if input.Body.Phone != nil {
		phone = *input.Body.Phone
	}
	if input.Body.Role != nil {
		role = *input.Body.Role
	}

	user, err := h.userService.UpdateUser(input.ID, email, phone, role)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	return &dtos.SingleUserResponse{Body: dtos.ToUserResponse(user)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *dtos.IDParam) (*dtos.EmptyResponse, error) {
	// Only admins can delete users
	if !middleware.IsAdmin(ctx) {
		return nil, huma.Error403Forbidden("Only admins can delete users")
	}

	auth := middleware.GetAuthContext(ctx)
	// Prevent self-deletion
	if auth != nil && auth.UserID == input.ID {
		return nil, huma.Error400BadRequest("Cannot delete your own account")
	}

	err := h.userService.DeleteUser(input.ID)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	return &dtos.EmptyResponse{}, nil
}
