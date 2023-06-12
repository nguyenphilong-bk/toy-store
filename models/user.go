package models

import (
	"errors"
	"time"

	"toy-store/db"
	"toy-store/forms"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID        uuid.UUID `db:"id, primarykey" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Name      string    `db:"name" json:"name"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Phone     string    `db:"phone" json:"phone"`
	Birthday  string    `db:"birthday" json:"birthday"`
	Cart      Cart      `json:"cart"`
}

// UserModel ...
type UserModel struct{}

var authModel = new(AuthModel)

// Login ...
func (m UserModel) Login(form forms.LoginForm) (user User, token Token, err error) {

	err = db.GetDB().SelectOne(&user, "SELECT id, email, password, name, updated_at, created_at FROM public.users WHERE email=LOWER($1) LIMIT 1", form.Email)

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
	checkUser, err := getDb.SelectInt("SELECT count(id) FROM public.users WHERE email=LOWER($1) LIMIT 1", form.Email)
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
	err = getDb.QueryRow("INSERT INTO public.users(email, password, name, phone) VALUES($1, $2, $3, $4) RETURNING id ", form.Email, string(hashedPassword), form.Name, form.Phone).Scan(&user.ID)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	user.Name = form.Name
	user.Email = form.Email
	user.Phone = form.Phone

	return user, err
}

// One ...
func (m UserModel) One(userID string) (user User, err error) {
	err = db.GetDB().SelectOne(&user, "SELECT id, name, email, phone FROM public.users WHERE id=$1 LIMIT 1", userID)
	return user, err
}
