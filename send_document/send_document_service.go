package send_document

import (
	jsoniter "github.com/json-iterator/go"
	"kpdigisign/model"
	"os"
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

func (s *service) SaveSendDocument(dto Dto,result string, notif string, reftrx string,
	jsonResponse string) (*Mapper, error) {
	var entity model.SendDocument

	reqSign, err := jsoniter.Marshal(dto.ReqSign)
	if err != nil {
		return nil, err
	}
	sendTo, err := jsoniter.Marshal(dto.SendTo)
	if err != nil {
		return nil, err
	}
	entity.UserID = os.Getenv("DIGISIGN_USER_ID")
	entity.DocumentID = dto.DocumentID
	entity.Payment = dto.Payment
	entity.SendTo = string(sendTo)
	entity.ReqSign = string(reqSign)

	entity.SendDocumentResult.RefTrx = reftrx
	entity.SendDocumentResult.Notif = notif
	entity.SendDocumentResult.Result = result
	entity.SendDocumentResult.JsonResponse = jsonResponse

	data, err := s.sendDocumentRepository.Save(entity)
	if err != nil {
		return nil, err
	}
	return s.sendDocumentMapper.Map(data), err
}
