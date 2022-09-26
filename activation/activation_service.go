package activation

import (
	"los-int-digisign/infrastructure/config/digisign"
	"los-int-digisign/model"
	"los-int-digisign/shared/constant"
	"time"

	"github.com/labstack/echo"
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
		tags := map[string]string{
			"app.pkg":    "activation",
			"app.func":   "FindAllActivations",
			"app.action": "read",
			"db.name":    "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	} else {
		return s.activationMapper.MapList(data), nil
	}
}

func (s *service) FindActivationById(id string) (*Mapper, error) {
	data, err := s.activationRepository.FindById(id)
	if err != nil {
		tags := map[string]string{
			"app.pkg":    "activation",
			"app.func":   "FindActivationById",
			"app.action": "read",
			"db.name":    "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	} else {
		return s.activationMapper.Map(data), nil
	}
}

func (s *service) SaveActivation(dto Dto, result string, link string, jsonResponse string, jsonRequest string) (*Mapper, error) {

	entity := model.Activation{
		UserID:     dto.UserID,
		ProspectID: dto.ProspectID,
		EmailUser:  dto.EmailUser,
		ActivationResult: model.ActivationResult{
			Result:       result,
			Link:         link,
			JsonResponse: jsonResponse,
			JsonRequest:  jsonRequest,
		},
	}
	data, err := s.activationRepository.Save(entity)
	if err != nil {
		tags := map[string]string{
			"app.pkg":    "activation",
			"app.func":   "SaveActivation",
			"app.action": "create",
			"db.name":    "di****gn",
		}
		extra := map[string]interface{}{
			"message":     err.Error(),
			"prospect_id": entity.ProspectID,
			"user_id":     entity.UserID,
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
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
		tags := map[string]string{
			"app.pkg":    "activation",
			"app.func":   "SaveActivationCallback",
			"app.action": "create",
			"db.name":    "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	}
	return echo.Map{"email": data.Email, "result": data.Result, "notif": data.Notif}, err
}

func (s *service) SaveTrxDigisign(email string, resp string) error {
	// Get Propsect ID in Customer Personal
	customer, err := s.activationRepository.FindCustomer(email)
	if err != nil {
		return err
	}

	entity := model.TrxDigisign{
		ProspectID: customer.ProspectID,
		Response:   resp,
		Activity:   constant.ACTIVITY_CALLBACK,
		CreatedAt:  time.Now(),
	}

	_, err = s.activationRepository.SaveDigisign(entity)
	if err != nil {
		return err
	}
	return nil
}
