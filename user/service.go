package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// standarisasi bentuk service
type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

type service struct {
	repository Repository
}

// fungsi untuk instansiasi service dengan parameter repository
func NewService(repository Repository) *service {
	return &service{repository}
}

// fungsi untuk mendapatkan input dari struct input lalu akan dikirimkan ke repository
func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	// instansiasi object dari struct User lalu assign dengan data dari
	// struct input yang diterima
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	// encrypt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	// pemanggilan atribut dari struct service yaitu repository dan memanggil
	// Save untuk menyimpan objek user kedalam repository lalu ke DB
	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No User found with that email.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	// cari user dengan ID
	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}

	// rubah nilai struct user
	user.AvatarFileName = fileLocation
	//simpan ke DB
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetUserByID(ID int) (User, error) {
	// panggil repo find by ID
	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}
	// jika tidak terdapat user dengan ID tersebut
	if user.ID == 0 {
		return user, errors.New("No user found with that ID")
	}

	return user, nil
}
