package user

import "golang.org/x/crypto/bcrypt"

// standarisasi bentuk service
type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
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
