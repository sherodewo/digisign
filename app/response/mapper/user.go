package mapper

import "kpdigisign/app/models"

type userMapper struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserMapper() *userMapper {
	return &userMapper{}
}
func (us *userMapper) Map(user models.User) *userMapper {
	us.ID = user.ID
	us.Username = user.Username
	us.Email = user.Email

	return us
}

func (us *userMapper) MapList(user []models.User) interface{} {
	var length = len(user)
	serialized := make([]userMapper, length, length)

	for k, v := range user {
		serialized[k] = userMapper{
			ID:       v.ID,
			Username: v.Username,
			Email:    v.Email,
		}
	}
	return serialized
}
