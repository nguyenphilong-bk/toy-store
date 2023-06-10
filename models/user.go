package models

import (
	"errors"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Phone     string    `json:"phone"`
	Birthday  string    `json:"birthday"`
}

// UserModel ...
type UserModel struct{}

var authModel = new(AuthModel)

// Login ...
func (m UserModel) Login(form forms.LoginForm) (user User, token Token, err error) {

	err = db.GetDB().Raw("SELECT id, email, password, name, updated_at, created_at FROM public.users WHERE email=LOWER(?) LIMIT 1", form.Email).Scan(&user).Error

	if err != nil {
		return user, token, err
	}

	//Compare the password form and database if match
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	if err != nil {
		return user, token, err
	}

	//Generate the JWT auth token
	tokenDetails, err := authModel.CreateToken(user.ID.String())
	if err != nil {
		return user, token, err
	}

	// saveErr := authModel.CreateAuth(user.ID.String(), tokenDetails)
	// if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	// }

	return user, token, nil
}

// Register ...
func (m UserModel) Register(form forms.RegisterForm) (user User, err error) {
	getDb := db.GetDB()

	//Check if the user exists in database
	var checkUser int
	err = getDb.Raw("SELECT count(id) FROM public.users WHERE email=LOWER(?) LIMIT 1", form.Email).Scan(&checkUser).Error
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	if checkUser > 0 {
		return user, errors.New("email already exists")
	}

	bytePassword := []byte(form.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	//Create the user and return back the user ID
	
	userID := ""
	err = getDb.Raw("INSERT INTO public.users(email, password, name, phone) VALUES(?, ?, ?, ?) RETURNING id ", form.Email, string(hashedPassword), form.Name, form.Phone).Scan(&userID).Error
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	user.ID = uuid.MustParse(userID)
	user.Name = form.Name
	user.Email = form.Email
	user.Phone = form.Phone

	return user, err
}

// One ...
func (m UserModel) One(userID int64) (user User, err error) {
	err = db.GetDB().First(&user, userID).Error
	return user, err
}
