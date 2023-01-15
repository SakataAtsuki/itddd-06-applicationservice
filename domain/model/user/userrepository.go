package user

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type UserRepositorier interface {
	FindByUserName(name *UserName) (*User, error)
	FindByUserId(userId *UserId) (*User, error)
	Save(user *User) error
	Update(user *User) error
	Delete(user *User) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

func (ur *UserRepository) FindByUserName(name *UserName) (user *User, err error) {
	tx, err := ur.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	rows, err := tx.Query("SELECT id, name FROM users WHERE name = $1", name.value)
	if err != nil {
		return nil, &FindByUserNameQueryError{UserName: *name, Message: fmt.Sprintf("userrepository.FindByUserName err: %s", err), Err: err}
	}
	defer rows.Close()

	userId := &UserId{}
	userName := &UserName{}
	for rows.Next() {
		err := rows.Scan(&userId.value, &userName.value)
		if err != nil {
			return nil, err
		}
		user = &User{id: *userId, name: *userName}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

type FindByUserNameQueryError struct {
	UserName UserName
	Message  string
	Err      error
}

func (err *FindByUserNameQueryError) Error() string {
	return err.Message
}

func (ur *UserRepository) Save(user *User) (err error) {
	tx, err := ur.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("INSERT INTO users(id, name) VALUES ($1, $2)", user.id.value, user.name.value)
	if err != nil {
		return &SaveQueryRowError{UserName: user.name, Message: fmt.Sprintf("userrepository.Save err: %s", err), Err: err}
	}
	return nil
}

type SaveQueryRowError struct {
	UserName UserName
	Message  string
	Err      error
}

func (err *SaveQueryRowError) Error() string {
	return err.Message
}

func (ur *UserRepository) FindByUserId(id *UserId) (user *User, err error) {
	tx, err := ur.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	rows, err := tx.Query("SELECT id, name FROM users WHERE id = $1", id.value)
	if err != nil {
		return nil, &FindByUserIdQueryError{UserId: *id, Message: fmt.Sprintf("userrepository.FindByUserId err: %s", err), Err: err}
	}
	defer rows.Close()

	userId := &UserId{}
	userName := &UserName{}
	for rows.Next() {
		err := rows.Scan(&userId.value, &userName.value)
		if err != nil {
			return nil, err
		}
		user = &User{id: *userId, name: *userName}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

type FindByUserIdQueryError struct {
	UserId  UserId
	Message string
	Err     error
}

func (err *FindByUserIdQueryError) Error() string {
	return err.Message
}

type UpdateQueryError struct {
	User    User
	Message string
	Err     error
}

func (err *UpdateQueryError) Error() string {
	return err.Message
}

func (ur *UserRepository) Update(user *User) (err error) {
	tx, err := ur.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("UPDATE users SET name=$2 WHERE id=$1", user.id.value, user.name.value)
	if err != nil {
		return &UpdateQueryError{User: *user, Message: fmt.Sprintf("userrepository.Delete err: %v", err), Err: err}
	}
	return nil
}

type DeleteQueryError struct {
	User    User
	Message string
	Err     error
}

func (err *DeleteQueryError) Error() string {
	return err.Message
}

func (ur *UserRepository) Delete(user *User) (err error) {
	tx, err := ur.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("DELETE FROM users WHERE id=$1", user.id.value)
	if err != nil {
		return &DeleteQueryError{User: *user, Message: fmt.Sprintf("userrepository.Delete err: %v", err), Err: err}
	}
	return nil
}
