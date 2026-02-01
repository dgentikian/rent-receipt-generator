package service

import (
	"fmt"
	"time"

	"github.com/dgentikian/rent-receipt-generator/internal/models"
	"github.com/dgentikian/rent-receipt-generator/internal/repository"
	"github.com/dgentikian/rent-receipt-generator/pkg/utils"
)

type LandlordService struct {
	repo      repository.LandlordRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewLandlordService(repo repository.LandlordRepository, jwtSecret string, jwtExpiry string) (*LandlordService, error) {
	expiry, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiry duration: %w", err)
	}
	return &LandlordService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: expiry,
	}, nil
}

func (s *LandlordService) Register(req *models.LandlordCreateRequest) (*models.Landlord, error) {
	// Check if email already exists
	existing, err := s.repo.GetByEmail(req.Email)
	if err == nil && existing != nil {
		// No error means user was found
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	landlord := &models.Landlord{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Address:      utils.StringToPtr(req.Address),
		City:         utils.StringToPtr(req.City),
		PostalCode:   utils.StringToPtr(req.PostalCode),
		Phone:        utils.StringToPtr(req.Phone),
	}

	if err := s.repo.Create(landlord); err != nil {
		return nil, fmt.Errorf("failed to create landlord: %w", err)
	}

	return landlord, nil
}

func (s *LandlordService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	landlord, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := utils.CheckPassword(landlord.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT(landlord.ID, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token:    token,
		Landlord: landlord,
	}, nil
}

func (s *LandlordService) GetByID(id int) (*models.Landlord, error) {
	return s.repo.GetByID(id)
}

func (s *LandlordService) Update(id int, req *models.LandlordUpdateRequest) (*models.Landlord, error) {
	landlord, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	landlord.FirstName = req.FirstName
	landlord.LastName = req.LastName
	landlord.Address = utils.StringToPtr(req.Address)
	landlord.City = utils.StringToPtr(req.City)
	landlord.PostalCode = utils.StringToPtr(req.PostalCode)
	landlord.Phone = utils.StringToPtr(req.Phone)

	if err := s.repo.Update(landlord); err != nil {
		return nil, fmt.Errorf("failed to update landlord: %w", err)
	}

	return landlord, nil
}

func (s *LandlordService) UpdateSignature(id int, signatureURL string) error {
	return s.repo.UpdateSignature(id, signatureURL)
}
