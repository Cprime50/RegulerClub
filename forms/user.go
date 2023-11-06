package forms

type Register struct {
	Username        string `json:"user_name" binding:"required,min=3,max=255,alphanum"` //alphanum is to get only alphanuneric chars, meaning just alpabets and numbers without any special chars
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=255"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email,min=3,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type SendEmail struct {
	Email string `json:"email" binding:"required,email,min=3,max=255"`
}

type VerifyEmail struct {
	Email string `json:"email" binding:"required,email,min=3,max=255"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type ResetPassword struct {
	Email           string `json:"email" binding:"required,email,min=3,max=255"`
	OTP             string `json:"otp" binding:"required,len=6"`
	Password        string `json:"password" binding:"required,min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=255"`
}
