package service

import (
	"context"
	"fmt"

	"github.com/escoutdoor/social/internal/repository"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/hasher"
	"github.com/escoutdoor/social/pkg/validator"
	"github.com/google/uuid"
)

type UserService struct {
	repo      repository.User
	validator *validator.Validator
}

func NewUserService(repo repository.User, validator *validator.Validator) *UserService {
	return &UserService{
		repo:      repo,
		validator: validator,
	}
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*types.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user types.User, input types.UpdateUserReq) (*types.User, error) {
	var err error

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil && !hasher.ComparePw(*input.Password, user.Password) {
		user.Password, err = hasher.HashPw(*input.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
	}
	if input.DOB != nil {
		dbo, err := s.validator.ValidateDate(*input.DOB)
		if err != nil {
			return nil, err
		}
		*user.DOB = types.DOB(dbo)
	}
	if input.Bio != nil {
		user.Bio = input.Bio
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}

	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
