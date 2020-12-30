package send_document

import (
	"los-int-digisign/model"

	jsoniter "github.com/json-iterator/go"
)

type service struct {
	sendDocumentRepository Repository
	sendDocumentMapper     *Mapper
}

func NewSendDocumentService(repository Repository, mapper *Mapper) *service {
	return &service{
		sendDocumentRepository: repository,
		sendDocumentMapper:     mapper,
	}
}

func (s *service) FindAllSendDocuments() (*[]Mapper, error) {
	data, err := s.sendDocumentRepository.FindAll()
	if err != nil {
		return nil, err
	} else {
		return s.sendDocumentMapper.MapList(data), nil
	}
}

func (s *service) FindSendDocumentById(id string) (*Mapper, error) {
	data, err := s.sendDocumentRepository.FindById(id)
	if err != nil {
		return nil, err
	} else {
		return s.sendDocumentMapper.Map(data), nil
	}
}

func (s *service) SaveSendDocument(dto Dto, result string, notif string, reftrx string,
	jsonResponse string, jsonRequest string) (*Mapper, error) {
	var entity model.SendDocument

	reqSign, err := jsoniter.Marshal(dto.ReqSign)
	if err != nil {
		return nil, err
	}
	sendTo, err := jsoniter.Marshal(dto.SendTo)
	if err != nil {
		return nil, err
	}
	entity.UserID = dto.UserID
	entity.DocumentID = dto.DocumentID
	entity.Payment = dto.Payment
	entity.SendTo = string(sendTo)
	entity.ReqSign = string(reqSign)

	entity.SendDocumentResult.RefTrx = reftrx
	entity.SendDocumentResult.Notif = notif
	entity.SendDocumentResult.Result = result
	entity.SendDocumentResult.JsonResponse = jsonResponse
	entity.SendDocumentResult.JsonRequest = jsonRequest

	data, err := s.sendDocumentRepository.Save(entity)
	if err != nil {
		return nil, err
	}
	return s.sendDocumentMapper.Map(data), err
}
