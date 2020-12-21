package activation

import (
	"github.com/labstack/echo"
	"los-int-digisign/model"
)

type service struct {
	activationRepository Repository
	activationMapper     *Mapper
}

func NewActivationService(repository Repository, mapper *Mapper) *service {
	return &service{
		activationRepository: repository,
		activationMapper:     mapper,
	}
}

func (s *service) FindAllActivations() (*[]Mapper, error) {
	data, err := s.activationRepository.FindAll()
	if err != nil {
		return nil, err
	} else {
		return s.activationMapper.MapList(data), nil
	}
}

func (s *service) FindActivationById(id string) (*Mapper, error) {
	data, err := s.activationRepository.FindById(id)
	if err != nil {
		return nil, err
	} else {
		return s.activationMapper.Map(data), nil
	}
}

func (s *service) SaveActivation(dto Dto, result string, link string, jsonResponse string) (*Mapper, error) {

	entity := model.Activation{
		UserID:     dto.UserID,
		ProspectID: dto.ProspectID,
		EmailUser:  dto.EmailUser,
		ActivationResult: model.ActivationResult{
			Result:       result,
			Link:         link,
			JsonResponse: jsonResponse,
		},
	}
	data, err := s.activationRepository.Save(entity)
	if err != nil {
		return nil, err
	}
	return s.activationMapper.Map(data), err
}

func (s *service) SaveActivationCallback(email string, result string, notif string) (interface{}, error) {
	var entity model.ActivationCallback
	entity.Email = email
	entity.Result = result
	entity.Notif = notif

	data, err := s.activationRepository.SaveCallback(entity)
	if err != nil {
		return nil, err
	}
	return echo.Map{"email": data.Email, "result": data.Result, "notif": data.Notif}, err
}
