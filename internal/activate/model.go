package activate

type ActivateAccountRequest struct {
	Email          string `json:"email" validate:"required,email"`
	ActivationCode string `json:"activation_code" validate:"required"`
}
