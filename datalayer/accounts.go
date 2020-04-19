package datalayer

import (
	"database/sql"
)

type Account struct {
	Model
	Email    sql.NullString `db:"email"`
	Password sql.NullString `db:"password"`
}

func (p *PersistenceDataLayer) GetAccountByEmail(email string) (*Account, error) {
	account := new(Account)
	row := p.GetConn().QueryRowx(`select id, email, password, created_at, updated_at, deleted_at from accounts where email = ?`,
		email)
	err := row.StructScan(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (p *PersistenceDataLayer) GetAccountByID(id int64) (*Account, error) {
	account := new(Account)
	err := p.GetConn().QueryRow(`SELECT id, email, password, created_at, updated_at, deleted_at FROM accounts WHERE id=?`, id).Scan(&account)
	 if err != nil {
		return nil, err
	}

	return account, nil
}

func (p *PersistenceDataLayer) CreateAccount(email, password string) (int64, error){
	result, err := p.GetConn().Exec("insert into accounts(email, password) values (?, ?)", email, password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}