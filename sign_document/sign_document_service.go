package sign_document

import (
	"los-int-digisign/infrastructure/config/digisign"
	"los-int-digisign/model"
	"los-int-digisign/shared/constant"
	"time"

	"github.com/labstack/echo"
)

type service struct {
	signDocumentRepository Repository
	signDocumentMapper     *Mapper
}

func NewSignDocumentService(repository Repository, mapper *Mapper) *service {
	return &service{
		signDocumentRepository: repository,
		signDocumentMapper:     mapper,
	}
}

func (s *service) FindAllSignDocuments() (*[]Mapper, error) {
	data, err := s.signDocumentRepository.FindAll()
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "sign_document",
			"app.func": "FindAllSignDocuments",
			"app.action":   "create",
			"db.name":  "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	} else {
		return s.signDocumentMapper.MapList(data), nil
	}
}

func (s *service) FindSignDocumentById(id string) (*Mapper, error) {
	data, err := s.signDocumentRepository.FindById(id)
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "sign_document",
			"app.func": "FindSignDocumentById",
			"app.action":   "create",
			"db.name":  "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	} else {
		return s.signDocumentMapper.Map(data), nil
	}
}

func (s *service) SaveSignDocument(dto Dto, result string, link string, jsonResponse string, jsonRequest string) (*Mapper, error) {
	var entity model.SignDocument
	entity.UserID = dto.UserID
	entity.EmailUser = dto.EmailUser
	entity.DocumentID = dto.DocumentID
	entity.SignDocumentResult.Result = result
	entity.SignDocumentResult.Link = link
	entity.SignDocumentResult.JsonResponse = jsonResponse
	entity.SignDocumentResult.JsonRequest = jsonRequest

	data, err := s.signDocumentRepository.Save(entity)
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "sign_document",
			"app.func": "SaveSignDocument",
			"app.action":   "create",
			"db.name":  "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
			"user_id": entity.UserID,
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	}
	return s.signDocumentMapper.Map(data), err
}

func (s *service) SaveSignDocumentCallback(documentId string, email string, statusDocument string, result string, notif string) (
	interface{}, error) {
	var entity model.SignDocumentCallback
	entity.DocumentID = documentId
	entity.Email = email
	entity.StatusDocument = statusDocument
	entity.Result = result
	entity.Notif = notif

	data, err := s.signDocumentRepository.SaveCallback(entity)
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "sign_document",
			"app.func": "SaveSignDocumentCallback",
			"app.action":   "create",
			"db.name":  "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	}
	return echo.Map{"id": data.ID, "document_id": data.DocumentID, "email": data.Email, "status_document": statusDocument,
		"result": result}, err
}

func (s *service) SaveTrxDigisign(docID string, resp string ) error{
	// Get Propsect ID in Customer Personal
	customer, err := s.signDocumentRepository.FindCustomer(docID)
	if err != nil {
		return err
	}

	entity := model.TrxDigisign{
		ProspectID: customer.ProspectID,
		Response:   resp,
		Activity:   constant.ACTIVITY_SEND_DOC,
		CreatedAt:  time.Now(),
	}

	_, err = s.signDocumentRepository.SaveDigisign(entity)
	if err != nil {
		return err
	}
	return nil
}