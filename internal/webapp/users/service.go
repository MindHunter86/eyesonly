package users

import "context"

type Service struct {
}

func (m *Service) CreateUser(c context.Context) (_ *User, e error) {
	return
}

func (m *Service) Get(c context.Context, id string) (_ *User, e error) {
	return
}

func (m *Service) List(c context.Context, limit int) (_ []User, e error) {
	return
}

func (m *Service) ChangeEmail(c context.Context, uid, email, pass string) (_ *User, e error) {
	return
}

func (m *Service) ChangePassword(c context.Context, uid, curpass, newpass string) (e error) {
	return
}
