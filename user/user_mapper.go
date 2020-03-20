package user

import "kpdigisign/model"

type Mapper struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserMapper() *Mapper {
	return &Mapper{}
}
func (m *Mapper) Map(model model.User) *Mapper {
	m.ID = model.ID
	m.Username = model.Username
	m.Email = model.Email
	return m
}

func (m *Mapper) MapList(model []model.User) *[]Mapper {
	var length = len(model)
	serialized := make([]Mapper, length)

	for k, v := range model {
		serialized[k] = Mapper{
			ID:    v.ID,
			Username:  v.Username,
			Email: v.Email,
		}
	}
	return &serialized
}
