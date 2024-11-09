package data

import (
	"time"

	"gorm.io/gorm"
)

type PostgresTestRepository struct {
	Conn *gorm.DB
}

func NewPostgresTestRepository(db *gorm.DB) *PostgresTestRepository {
	return &PostgresTestRepository{
		Conn: db,
	}
}

func (u *PostgresTestRepository) GetAll() ([]*User, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()

	// var users []*User

	// err := db.Model(&User{}).WithContext(ctx).Find(&users).Error
	// if err != nil {
	// 	return nil, err
	// }

	users := []*User{}

	return users, nil
}

// GetByEmail returns one user by email
func (u *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()

	// var user *User
	// err := db.Model(&User{}).WithContext(ctx).Where("email = ?", email).First(&user).Error
	// if err != nil {
	// 	return nil, err
	// }

	user := &User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return user, nil
}

// GetOne returns one user by id
func (u *PostgresTestRepository) GetOne(id int) (*User, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()

	// var user *User
	// err := db.Model(&User{}).WithContext(ctx).Where("id = ?", id).First(&user).Error
	// if err != nil {
	// 	return nil, err
	// }

	user := &User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

func (u *PostgresTestRepository) Update(user User) error {
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()

	// err := db.Model(&User{}).WithContext(ctx).Where("id =?", user.ID).UpdateColumns(user).Error
	// if err != nil {
	// 	return err
	// }

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *PostgresTestRepository) DeleteByID(id int) error {
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()

	// err := db.Model(&User{}).WithContext(ctx).Delete(&User{}, id).Error
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *PostgresTestRepository) Insert(user User) (int, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	// if err != nil {
	// 	return 0, err
	// }

	// var newID int
	// insertUser := User{
	// 	Email:     user.Email,
	// 	FirstName: user.FirstName,
	// 	LastName:  user.LastName,
	// 	Password:  string(hashedPassword),
	// 	Active:    user.Active,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }

	// err = db.Model(&User{}).WithContext(ctx).Create(insertUser).Error
	// if err != nil {
	// 	return 0, err
	// }

	newID := 101

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *PostgresTestRepository) ResetPassword(password string, user User) error {
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	// if err != nil {
	// 	return err
	// }

	// err = db.Model(&User{}).WithContext(ctx).Where("id = ?", user.ID).UpdateColumn("password", hashedPassword).Error
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (u *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
	// 		return false, nil
	// 	default:
	// 		return false, err
	// 	}
	// }

	return true, nil
}
