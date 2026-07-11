package services

import (
	"agenda-app/app/internal/repository"
)

type LoginService interface {
}

func NewLoginService(r repository.UserRepository) LoginService {
	return &loginService{repo: r}
}

type loginService struct {
	repo repository.UserRepository
}
