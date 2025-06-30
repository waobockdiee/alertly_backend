package signup

type User struct {
	ID             int64  `json:"id"`
	Email          string `json:"email" validate:"required,email"`
	FirstName      string `json:"first_name" validate:"required,min=2"`
	LastName       string `json:"last_name" validate:"required,min=2"`
	Password       string `json:"password" validate:"required,min=6"`
	BirthYear      string `json:"birth_year"`
	BirthMonth     string `json:"birth_month"`
	BirthDay       string `json:"birth_day"`
	ActivationCode string
	Nickname       string
}

type ActivationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"activation_code" binding:"required,len=5"`
}
