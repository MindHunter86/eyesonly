package users

import "time"

type (
	User struct {
		ID        string    `json:"id"`
		Email     string    `json:"email"`
		Password  string    `json:"-"`
		Role      string    `json:"role"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	PublicUser struct {
		ID     string
		Email  string
		Role   string
		Status string
	}
)

func Public(u *User) *PublicUser {
	return &PublicUser{
		ID:     u.ID,
		Email:  u.Email,
		Role:   u.Role,
		Status: u.Status,
	}
}
