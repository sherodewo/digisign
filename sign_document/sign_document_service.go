package sign_document

import (
	"github.com/labstack/echo"
	"kpdigisign/model"
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
		return nil, err
	} else {
		return s.signDocumentMapper.MapList(data), nil
	}
}

func (s *service) FindSignDocumentById(id string) (*Mapper, error) {
	data, err := s.signDocumentRepository.FindById(id)
	if err != nil {
		return nil, err
	} else {
		return s.signDocumentMapper.Map(data), nil
	}
}

func (s *service) SaveSignDocument(dto Dto, result string, link string, jsonResponse string) (*Mapper, error) {
	var entity model.SignDocument
	entity.UserID = dto.UserID
	entity.EmailUser = dto.EmailUser
	entity.DocumentID = dto.DocumentID
	entity.SignDocumentResult.Result = result
	entity.SignDocumentResult.Link = link
	entity.SignDocumentResult.JsonResponse = jsonResponse

	data, err := s.signDocumentRepository.Save(entity)
	if err != nil {
		return nil, err
	}
	return s.signDocumentMapper.Map(data), err
}

func (s *service) SaveSignDocumentCallback(documentId string, email string, statusDocument string, result string) (
	interface{}, error) {
	var entity model.SignDocumentCallback
	entity.DocumentID = documentId
	entity.Email = email
	entity.StatusDocument = statusDocument
	entity.Result = result

	data, err := s.signDocumentRepository.SaveCallback(entity)
	if err != nil {
		return nil, err
	}
	return echo.Map{"id": data.ID, "document_id": data.DocumentID, "email": data.Email, "status_document": statusDocument,
		"result": result}, err
}
