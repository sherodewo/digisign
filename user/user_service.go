package user

import (
	"los-int-digisign/model"
	"los-int-digisign/utils"
)

type service struct {
	userRepository Repository
	userMapper     *Mapper
}

func NewUserService(repository Repository, mapper *Mapper ) *service {
	return &service{
		userRepository: repository,
		userMapper:mapper,
	}
}

func (s *service) FindAllUsers() (*[]Mapper, error) {
	data, err := s.userRepository.FindAll()
	if err != nil {
		return nil, err
	} else {
		return s.userMapper.MapList(data), nil
	}
}

func (s *service) FindUserById(id string) (*Mapper, error) {
	data, err := s.userRepository.FindById(id)
	if err != nil {
		return nil, err
	} else {
		return s.userMapper.Map(data), nil
	}
}

func (s *service) SaveUser(dto Dto) (*Mapper, error) {
	hashPassword, _ := utils.HashPassword(dto.Password)
	entity := model.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashPassword,
	}
	data, err := s.userRepository.Save(entity)
	if err != nil {
		return nil, err
	} else {
		return s.userMapper.Map(data), nil
	}
}

func (s *service) UpdateUser(id string, dto Dto) (*Mapper, error) {
	hashPassword, _ := utils.HashPassword(dto.Password)
	entity := model.User{
		ID:       id,
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashPassword,
	}
	data, err := s.userRepository.Update(entity)
	if err != nil {
		return nil, err
	} else {
		return s.userMapper.Map(data), nil
	}
}

func (s *service) DeleteUser(id string) error {
	entity := model.User{
		ID: id,
	}
	err := s.userRepository.Delete(entity)
	if err != nil {
		return err
	} else {
		return nil
	}
}
