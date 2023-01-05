package models

import (
	"crypto/rand"
	"database/sql"
	"io"

	"github.com/google/uuid"
	"github.com/hexcraft-biz/model"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const (
	PW_SALT_BYTES = 16
)

// ================================================================
// Data Struct
// ================================================================
type EntityUser struct {
	*model.Prototype `dive:""`
	Identity         string `db:"identity"`
	Password         []byte `db:"password"`
	Salt             []byte `db:"salt"`
	Status           string `db:"status"`
}

func (u *EntityUser) GetAbsUser() (*AbsUser, error) {
	return &AbsUser{
		ID:        *u.ID,
		Identity:  u.Identity,
		Password:  string(u.Password),
		Salt:      string(u.Salt),
		Status:    u.Status,
		CreatedAt: u.Ctime.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.Mtime.Format("2006-01-02 15:04:05"),
	}, nil
}

type AbsUser struct {
	ID        uuid.UUID `json:"id"`
	Identity  string    `json:"identity"`
	Password  string    `json:"_"`
	Salt      string    `json:"_"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// ================================================================
// Engine
// ================================================================
type UsersTableEngine struct {
	*model.Engine
}

func NewUsersTableEngine(db *sqlx.DB) *UsersTableEngine {
	return &UsersTableEngine{
		Engine: model.NewEngine(db, "users"),
	}
}

func (e *UsersTableEngine) Insert(identity string, password string, status string) (*EntityUser, error) {
	saltBytes := make([]byte, PW_SALT_BYTES)
	if _, err := io.ReadFull(rand.Reader, saltBytes); err != nil {
		return nil, err
	}
	salt := string(saltBytes)

	pwdBytes := []byte(password + salt)

	hashBytes, hashErr := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.DefaultCost)
	if hashErr != nil {
		return nil, hashErr
	}

	u := &EntityUser{
		Prototype: model.NewPrototype(),
		Identity:  identity,
		Password:  hashBytes,
		Salt:      saltBytes,
		Status:    status,
	}

	_, err := e.Engine.Insert(u)
	return u, err
}

func (e *UsersTableEngine) GetByID(id string) (*EntityUser, error) {
	row := EntityUser{}
	q := `SELECT * FROM ` + e.TblName + ` WHERE id = UUID_TO_BIN(?);`
	if err := e.Engine.Get(&row, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &row, nil
}

func (e *UsersTableEngine) GetByIdentity(identity string) (*EntityUser, error) {
	row := EntityUser{}
	q := `SELECT * FROM ` + e.TblName + ` WHERE identity = ?;`
	if err := e.Engine.Get(&row, q, identity); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &row, nil
}

func (e *UsersTableEngine) ResetPwd(id *uuid.UUID, password string, saltBytes []byte) (int64, error) {
	salt := string(saltBytes)

	pwdBytes := []byte(password + salt)

	hashBytes, hashErr := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.DefaultCost)
	if hashErr != nil {
		return 0, hashErr
	}

	q := `UPDATE ` + e.TblName + ` SET password = ? WHERE id = UUID_TO_BIN(?);`
	if rst, err := e.Exec(q, hashBytes, &id); err != nil {
		return 0, err
	} else {
		return rst.RowsAffected()
	}
}

func (e *UsersTableEngine) UpdateStatus(id *uuid.UUID, status string) (int64, error) {
	q := `UPDATE ` + e.TblName + ` SET status = ? WHERE id = UUID_TO_BIN(?);`
	if rst, err := e.Exec(q, status, &id); err != nil {
		return 0, err
	} else {
		return rst.RowsAffected()
	}
}
