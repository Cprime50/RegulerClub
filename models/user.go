package models

import (
	"crypto/rand"
	"errors"
	"log"
	"time"

	"github.com/Cprime50/RegulerClub/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"uniqueIndex;not null; size:255"  json:"user_name"`
	Email    string `gorm:"uniqueIndex;not null; size:255" json:"email"`
	Password string `gorm:"not null; collate:utf8" json:"-"`
	//	Active      bool      `gorm:"not null" json:"active"`
	Verified  bool      `gorm:"not null;DEFAULT:false" json:"verified"`
	Bio       string    `gorm:"size:1024" json:"bio"`
	Image     *string   `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
	RoleID    uint      `gorm:"not null;DEFAULT:4" json:"role_id"`
	Role      Role      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type Role struct {
	gorm.Model
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"size:50;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255;not null" json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OTP struct {
	gorm.Model
	ID        uint
	UserID    uint                `gorm:"not null" json:"user_id"`
	User      *User               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Email     string              `json:"email"`
	Code      string              `json:"-"`
	ExpiresAt time.Time           `json:"expires_at"`
	TypeID    uint                `gorm:"not null" json:"type_id"`
	Type      OTPVerificationType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"type"`
}

type OTPVerificationType struct {
	ID   uint
	Type string
	Ttl  int64
}

// hash password before we save it
func (user *User) SetPassword(password string) string {
	bytePassword := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypting password failed", err)
	}
	user.Password = string(passwordHash)

	return user.Password
}

// create user
func CreateUser(user *User) error {
	if err := database.Db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// confirm user password whether it matches that in database
func (user *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// we can get two types of error here
	if err != nil {
		switch {
		// error when  password doesnt match
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil

		//and error when something unexpected happens
		default:
			return false, err
		}
	}
	return true, nil
}

// Update password
func (user *User) UpdatePassword(email string, hashedpassword string) error {
	err := database.Db.Model(&User{}).Where("email = ?", email).Update("password", hashedpassword).Error
	if err != nil {
		return err
	}
	return nil
}

// Get user by ID
// find the first user with the id
func (user *User) GetUserById(id uint) (*User, error) {
	err := database.Db.Where("id=?", id).Omit("password").First(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// GET FIRST USER BY EMAIL
func (user *User) GetUserByEmail(email string) (*User, error) {
	err := database.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// GET FIRST USER BY Username
func (user *User) GetUserByUsername(username string) (*User, error) {
	err := database.Db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// Get all users
func (user *User) GetAllUsers() ([]*User, error) {
	var users []*User
	if err := database.Db.Omit("password").Find(&user).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update User
func UpdateUser(user *User) error {
	err := database.Db.Omit("password").Updates(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete User
func (user *User) DeleteUser(id uint) error {
	if err := database.Db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

// Verify Account
// Updates users veferied field to true
func (user *User) VerifyAccount(email string) error {
	if err := database.Db.Model(&User{}).Where("email = ?", email).Update("verified", true).Error; err != nil {
		return err
	}
	return nil
}

// ROLE
// create role
func CreateRole(role *Role) error {
	if err := database.Db.Create(&role).Error; err != nil {
		return err
	}
	return nil
}

// Get one role by id
func (role *Role) GetRoleById(id uint) (*Role, error) {
	if err := database.Db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// Gets all roles
func (role *Role) GetAllRoles() ([]*Role, error) {
	var roles []*Role
	if err := database.Db.Find(&role).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// Update role
func UpdateRole(role *Role) error {
	if err := database.Db.Save(&role).Error; err != nil {
		return err
	}
	return nil
}

// Delete role
func (role *Role) DeleteRole(id uint) error {
	if err := database.Db.Where("id = ?", id).Delete(&role).Error; err != nil {
		return err
	}
	return nil
}

// Assign roles to users
func (role *Role) AssignRole(userID uint, roleID uint) error {
	if err := database.Db.Model(&User{}).Where("id = ?", userID).Update("role", roleID).Error; err != nil {
		return err
	}
	return nil
}

// OTP
// Generate random 6 digit token, make token after in 5secs after creation, if user inputs otp wrong 3 times make token expire

const (
	OTPChars  = "0123456789"
	OTPLength = 6
)

// generate OTP
func GenerateOTP() (string, error) {
	buffer := make([]byte, OTPLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	otpCharsLen := len(OTPChars)
	for i := 0; i < OTPLength; i++ {
		buffer[i] = OTPChars[int(buffer[i])%otpCharsLen]
	}
	return string(buffer), nil
}

// Create OTP for email verification
func CreateOTPForEmailVerification(user *User) (string, error) {
	code, err := GenerateOTP()
	if err != nil {
		log.Println("problem generating OTP")
		return "", err
	}
	// ID =1 is for email verification and ID=2 is for password reset
	id := uint(1)
	verificationType, err := GetVerificationType(id)

	otp := &OTP{
		UserID:    user.ID,
		Email:     user.Email,
		Code:      code,
		Type:      *verificationType,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(verificationType.Ttl)),
	}
	// make this to save/update if an otp aready exist instead fo creating new one
	err = database.Db.Model(&OTP{}).Save(otp).Error
	if err != nil {
		return "", err
	}
	return code, err
}

func GetVerificationType(id uint) (*OTPVerificationType, error) {
	var verificationType *OTPVerificationType
	if err := database.Db.Where("id = ?", id).First(&verificationType).Error; err != nil {
		return nil, err
	}
	return verificationType, nil
}

// Create OTP for password reset
func CreateOTPForPasswordReset(user *User) (string, error) {
	code, err := GenerateOTP()
	if err != nil {
		log.Fatal("problem generating OTP")
		return "", err
	}
	// ID =1 is for email verification and ID=2 is for password reset
	id := uint(2)
	verificationType, err := GetVerificationType(id)

	otp := &OTP{
		UserID:    user.ID,
		Email:     user.Email,
		Code:      code,
		Type:      *verificationType,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(verificationType.Ttl)),
	}

	err = database.Db.Model(&OTP{}).Save(otp).Error
	if err != nil {
		return "", err
	}
	return code, err
}

// TODO use email and get otp then in controller compare this otp gotten with otp inputed to know if its wrong,
// the increase wrongg tre by 1 for each try "wrongtrie ++", if wrongtry == 3. delete the otp and craete new one and sned to yser
func GetOTPByEmail(email string) (*OTP, error) {
	var OTP *OTP
	if err := database.Db.Where("email = ?", email).First(&OTP).Error; err != nil {
		return OTP, err
	}
	return OTP, nil
}

func DeleteOTP(otp *OTP) error {
	if err := database.Db.Where("id = ?", otp.ID).Delete(&OTP{}).Error; err != nil {
		return err
	}
	return nil
}
