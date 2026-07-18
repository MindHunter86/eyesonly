package users

import "database/sql"

type Store struct {
	db *sql.DB
}

func (m *Store) Create() error { return nil }

func (m *Store) FindByEmail() error { return nil }

func (m *Store) FindByID() error { return nil }

func (m *Store) List() error { return nil }

func (m *Store) UpdateEmail() error { return nil }

func (m *Store) UpdatePassword() error { return nil }
