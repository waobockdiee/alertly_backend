package signup

type User struct {
	ID             int64  `json:"id"`
	Email          string `json:"email" validate:"required,email"`
	FirstName      string `json:"first_name" validate:"required,min=2"`
	LastName       string `json:"last_name" validate:"required,min=2"`
	Password       string `json:"password" validate:"required,min=6"`
	BirthYear      string `json:"birth_year" validate:"required"`
	BirthMonth     string `json:"birth_month" validate:"required"`
	BirthDay       string `json:"birth_day" validate:"required"`
	ActivationCode string
	Nickname       string
}
