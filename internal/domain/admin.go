package domain

type Admin struct {
	ID       uint   `json:"id"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
