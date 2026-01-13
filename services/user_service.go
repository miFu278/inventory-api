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
	// Check if username exists
	existingUser, _ := s.userRepo.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	existingEmail, _ := s.userRepo.GetUserByEmail(email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
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
		return "", nil, errors.New("invalid username or password")
	}

	if !utils.CheckPassword(password, user.PasswordHashed) {
		return "", nil, errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return "", nil, err
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
		return nil, err
	}

	if email != "" {
		user.Email = email
	}
	if phone != "" {
		user.Phone = phone
	}
	if role != "" {
		user.Role = role
	}

	return s.userRepo.UpdateUser(user)
}

func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(oldPassword, user.PasswordHashed) {
		return errors.New("invalid old password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHashed = hashedPassword
	_, err = s.userRepo.UpdateUser(user)
	return err
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.DeleteUser(id)
}
