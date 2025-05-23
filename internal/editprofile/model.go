package editprofile

type Account struct {
	AccountID        int64  `db:"account_id" json:"account_id"`
	Email            string `db:"email" json:"email"`
	FirstName        string `db:"first_name" json:"first_name"`
	LastName         string `db:"last_name" json:"last_name"`
	Password         string `db:"password" json:"password"`
	NewPassword      string `json:"new_password"`
	Code             string `json:"code"` // without db propert, because this struct will be use in multiples account updates and columns with code context. Example: update_email_code, update_first_name
	IsPremium        bool   `db:"is_premium" json:"is_premium"`
	BirthYear        string `db:"birth_year" json:"birth_year"`
	BirthMonth       string `db:"birth_month" json:"birth_month"`
	BirthDay         string `db:"birth_day" json:"birth_day"`
	IsPrivateProfile bool   `db:"is_private" json:"is_private"`
	PhoneNumber      string `db:"phone_number" json:"phone_number"`
	NickName         string `db:"nickname" json:"nickname"`
}

type Media struct {
	Uri string `json:"uri"`
}
