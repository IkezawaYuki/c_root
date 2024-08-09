package entity

type UserRose int

const (
	UserRoleCustomer UserRose = 0
	UserRoleAdmin    UserRose = 1
)

type UserSession struct {
	UserID string
	Role   UserRose
}
