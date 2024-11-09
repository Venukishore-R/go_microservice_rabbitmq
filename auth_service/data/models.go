package data

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const dbTimeout = time.Second * 3

var db *gorm.DB

type PostgresRepository struct {
	Conn *gorm.DB
}

func NewPostgresRepository(pool *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		Conn: pool,
	}
}

// type Models struct {
// 	User User
// }

// func New(dbPool *gorm.DB) *Models {
// 	db = dbPool

// 	return &Models{
// 		User: User{},
// 	}
// }

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Email     string    `json:"email" gorm:"unique"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAll returns a slice of all users, sorted by last name
func (u *PostgresRepository) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var users []*User

	err := db.Model(&User{}).WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetByEmail returns one user by email
func (u *PostgresRepository) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user *User
	err := db.Model(&User{}).WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetOne returns one user by id
func (u *PostgresRepository) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user *User
	err := db.Model(&User{}).WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *PostgresRepository) Update(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	err := db.Model(&User{}).WithContext(ctx).Where("id =?", user.ID).UpdateColumns(user).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *PostgresRepository) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	err := db.Model(&User{}).WithContext(ctx).Delete(&User{}, id).Error
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *PostgresRepository) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	insertUser := User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  string(hashedPassword),
		Active:    user.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.Model(&User{}).WithContext(ctx).Create(insertUser).Error
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *PostgresRepository) ResetPassword(password string, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	err = db.Model(&User{}).WithContext(ctx).Where("id = ?", user.ID).UpdateColumn("password", hashedPassword).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *PostgresRepository) PasswordMatches(plainText string, user User) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
