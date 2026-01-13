package services

import (
	"errors"
	"inventory-api/models"
	"inventory-api/repo"
	"inventory-api/utils"
)

type UserService struct {
	userRepo  repo.UserRepositoryInterface
	jwtSecret string
}

func NewUserService(userRepo repo.UserRepositoryInterface, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) Register(username, password, email, phone, role string) (*models.User, error) {
	// Validate role
	if !models.IsValidRole(role) {
		return nil, errors.New("invalid role")
	}

	// Check if username exists
	existingUser, _ := s.userRepo.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errors.New("registration failed")
	}

	// Check if email exists
	existingEmail, _ := s.userRepo.GetUserByEmail(email)
	if existingEmail != nil {
		return nil, errors.New("registration failed")
	}

	// Validate password strength
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	// Create user
	user := &models.User{
		Username:       username,
		PasswordHashed: hashedPassword,
		Email:          email,
		Phone:          phone,
		Role:           role,
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserService) Login(username, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(password, user.PasswordHashed) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(username)
}

func (s *UserService) GetAllUsers(limit, offset int) ([]models.User, error) {
	return s.userRepo.GetAllUsers(limit, offset)
}

func (s *UserService) UpdateUser(id uint, email, phone, role string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate and update email
	if email != "" {
		// Check if email is already taken by another user
		existingEmail, _ := s.userRepo.GetUserByEmail(email)
		if existingEmail != nil && existingEmail.ID != id {
			return nil, errors.New("email already in use")
		}
		user.Email = email
	}

	// Update phone
	if phone != "" {
		user.Phone = phone
	}

	// Validate and update role
	if role != "" {
		if !models.IsValidRole(role) {
			return nil, errors.New("invalid role")
		}
		user.Role = role
	}

	return s.userRepo.UpdateUser(user)
}

func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPassword(oldPassword, user.PasswordHashed) {
		return errors.New("current password is incorrect")
	}

	// Validate new password strength
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	// Check if new password is same as old
	if oldPassword == newPassword {
		return errors.New("new password must be different from current password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to process password")
	}

	user.PasswordHashed = hashedPassword
	_, err = s.userRepo.UpdateUser(user)
	if err != nil {
		return errors.New("failed to update password")
	}
	return nil
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.DeleteUser(id)
}
