package models

import (
	"context"
	"errors"

	"github.com/MindHunter86/eyesonly/internal/database/sqlite"
	"github.com/MindHunter86/eyesonly/internal/utils"
)

type Account struct {
	ID         uint32
	Email      string
	IsVerified bool
}

type accountPassword struct {
	s256 []byte
}

func (m *accountPassword) release() {
	// todo - add full slice clean
	m.s256[1] = 0
	m.s256[2] = 0
	m.s256[3] = 0

	m.s256 = m.s256[:0]
}

var ErrNoInput = errors.New("there were empty values found")
var ErrCreateAccountAlrdExists = errors.New("received email has been already registered")
var ErrNoAccFound = errors.New("there is no account with received email/login pair")
var ErrAccIsBlocked = errors.New("the account has been blocked, call support")

func CreateAccount(c context.Context, email, password []byte) (_ *Account, e error) {
	if len(email) == 0 || len(password) == 0 {
		return nil, ErrNoInput
	}

	// todo : check password strength

	// prepare context
	db := utils.ContextValueExtract[*sqlite.SqliteDB](c, utils.CtxCliContext)

	// check email existance in DB
	var ex bool
	if ex, e = db.IsEmailExists(utils.UnsafeString(email)); e != nil {
		return
	} else if ex {
		return nil, ErrCreateAccountAlrdExists
	}

	// todo : all salt for sha256. Migrate to sha512 or bcrypt
	// hash password
	pwd := acquirePooledType[accountPassword]()
	defer releasePooledType(pwd)

	if pwd.s256, e = hashSha256(pwd.s256, password); e != nil {
		return
	}

	// create new account
	if e = db.CreateAccount(c, utils.UnsafeString(email), utils.UnsafeString(pwd.s256)); e != nil {
		return
	}

	// todo send verification email

	var sacc *sqlite.SchemaAccount
	if sacc, e = db.LoadAccount(c, utils.UnsafeString(email), utils.UnsafeString(pwd.s256)); e != nil {
		return
	} else if sacc == nil {
		return nil, ErrNoAccFound
	}

	return &Account{
		ID:         sacc.ID,
		Email:      sacc.Email,
		IsVerified: sacc.IsVerified,
	}, nil
}

func LoginAccount(c context.Context, email, password []byte) (_ *Account, e error) {
	if len(email) == 0 || len(password) == 0 {
		return nil, ErrNoInput
	}

	// todo : all salt for sha256. Migrate to sha512 or bcrypt
	// hash password
	pwd := acquirePooledType[accountPassword]()
	defer releasePooledType(pwd)

	if pwd.s256, e = hashSha256(pwd.s256, password); e != nil {
		return
	}

	// prepare context
	db := utils.ContextValueExtract[*sqlite.SqliteDB](c, utils.CtxCliContext)

	// find account
	var sacc *sqlite.SchemaAccount
	if sacc, e = db.LoadAccount(c, utils.UnsafeString(email), utils.UnsafeString(pwd.s256)); e != nil {
		return
	} else if sacc == nil {
		return nil, ErrNoAccFound
	}

	if sacc.IsBlocked {
		return nil, ErrAccIsBlocked
	}

	return &Account{
		ID:         sacc.ID,
		Email:      sacc.Email,
		IsVerified: sacc.IsVerified,
	}, nil
}
