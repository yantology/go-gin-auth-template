package services

import (
	"errors"
	"log"

	"github.com/yantology/go-gin-auth-template/internal/models"
	"github.com/yantology/go-gin-auth-template/internal/repository"
	"github.com/yantology/go-gin-auth-template/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtUtil  *utils.JWTUtil
}

func NewAuthService(userRepo *repository.UserRepository, jwtUtil *utils.JWTUtil) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) error {
	log.Println("Registering user")
	// Check if user exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already registered")
	}

	log.Println("Creating user")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.TokenResponse, error) {
	// Get user
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.jwtUtil.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) ChangePassword(userID int, req *models.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetById(userID)
	if err != nil {
		return err
	}

	// Check old password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	return s.userRepo.UpdatePassword(userID, string(hashedPassword))
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	// Validate refresh token
	token, err := s.jwtUtil.ValidateToken(refreshToken, true)
	if err != nil {
		return nil, err
	}

	// Extract user ID
	userID, err := s.jwtUtil.ExtractUserID(token)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.jwtUtil.GenerateTokens(userID)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
