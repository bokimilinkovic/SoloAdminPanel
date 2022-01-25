package postgres

import (
	"github.com/beevik/guid"
	"github.com/jinzhu/gorm"
	"github.com/kolosek/pkg/model/domain"
)

// User is an implementation of repository.User interface.
type User struct {
	db *gorm.DB
}

// NewUser creates a new User repository.
func NewUser(db *gorm.DB) *User {
	return &User{db: db}
}

// FindByID tries to find a user with the given ID. Error is returned only in case of internal database errors.
func (u *User) FindByID(id int64) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("id = ?", id).Find(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail one returns a specific user with whole email sent.
func (u *User) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("email = ?", email).Find(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Exists checks if the user with given email exists in database
func (u *User) Exists(user *domain.User) error {
	existing, err := u.FindByEmail(user.Email)
	if err != nil {
		return err
	}

	// User with given email does not exist in DB, we are sure that email is unique
	if existing == nil {
		return nil
	}

	return nil
}

// Create creates a new user.
func (u *User) Create(user *domain.User) error {
	user.GUID = guid.NewString()
	return u.db.Create(user).Error
}

//Save user
func (u *User) Save(user *domain.User) error {
	user.GUID = guid.NewString()
	return u.db.Save(user).Error
}

// Update updates an existing user.
func (u *User) Update(user *domain.User) error {
	return u.db.Update(user).Error
}

// FindFirstN retrieves first N={10} records.
func (u *User) FindFirstN(count int) ([]*domain.User, error) {
	var users []*domain.User
	if err := u.db.Limit(count).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

//MigrateTables in database
func (a *User) MigrateTables() error {
	if !a.db.HasTable(&domain.User{}) {
		return a.db.AutoMigrate(&domain.User{}).Error
	}

	return nil
}
