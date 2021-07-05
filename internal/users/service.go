package users

import (
	"github.com/google/uuid"
	"github.com/robertscherbarth/go-service-skeleton/internal/users/ports"
	"go.uber.org/zap"
)

type Store interface {
	Create(user User) error
	Delete(id string) error
	FindByID(id string) (User, error)
	FindAll() ([]User, error)
}

type Service struct {
	logger *zap.Logger
	store  Store
}

func NewService(logger *zap.Logger, store Store) *Service {
	return &Service{
		logger: logger,
		store:  store,
	}
}

func (s *Service) Add(user ports.User) error {
	id := uuid.New()
	s.logger.Info("create user", zap.String("id", id.String()))

	var tag string
	if user.Tag != nil {
		tag = *user.Tag
	}

	u := User{
		ID:   id,
		Name: user.Name,
		Tag:  tag,
	}
	return s.store.Create(u)
}

func (s *Service) Delete(id string) error {
	s.logger.Info("delete user", zap.String("id", id))
	return s.Delete(id)
}

func (s *Service) FindAll() ([]ports.User, error) {
	s.logger.Info("find all users")
	users, err := s.store.FindAll()
	if err != nil {
		return nil, err
	}

	var responseUsers []ports.User
	for _, v := range users {
		responseUsers = append(responseUsers, s.convert(v))
	}

	return responseUsers, nil
}

func (s *Service) FindByID(id string) (ports.User, error) {
	s.logger.Info("find user by id", zap.String("id", id))

	user, err := s.store.FindByID(id)
	if err != nil {
		return ports.User{}, err
	}

	responseUser := s.convert(user)
	return responseUser, nil
}

func (s *Service) convert(user User) ports.User {
	return ports.User{
		NewUser: ports.NewUser{
			Name: user.Name,
			Tag:  &user.Tag,
		},
		Id: user.ID.String(),
	}
}
