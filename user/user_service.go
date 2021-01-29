package user

import (
	"los-int-digisign/infrastructure/config/digisign"
	"los-int-digisign/model"
	"los-int-digisign/utils"
)

type service struct {
	userRepository Repository
	userMapper     *Mapper
}

func NewUserService(repository Repository, mapper *Mapper) *service {
	return &service{
		userRepository: repository,
		userMapper:     mapper,
	}
}

func (s *service) FindAllUsers() (*[]Mapper, error) {
	data, err := s.userRepository.FindAll()
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "user",
			"app.func": "FindAllUsers",
			"app.action":   "read",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}

		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	} else {
		return s.userMapper.MapList(data), nil
	}
}

func (s *service) FindUserById(id string) (*Mapper, error) {
	data, err := s.userRepository.FindById(id)
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "user",
			"app.func": "FindUserById",
			"app.action":   "read",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}

		digisign.SendToSentry(tags, extra, "DATABASE")
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
		tags := map[string]string{
			"app.pkg":  "user",
			"app.func": "SaveUser",
			"app.action":   "create",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}

		digisign.SendToSentry(tags, extra, "DATABASE")
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
		tags := map[string]string{
			"app.pkg":  "user",
			"app.func": "UpdateUser",
			"app.action":   "update",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}

		digisign.SendToSentry(tags, extra, "DATABASE")
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
		tags := map[string]string{
			"app.pkg":  "user",
			"app.func": "DeleteUser",
			"app.action":   "delete",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}

		digisign.SendToSentry(tags, extra, "DATABASE")
		return err
	} else {
		return nil
	}
}
