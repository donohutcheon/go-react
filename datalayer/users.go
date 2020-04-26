package datalayer

import (
	"database/sql"
)

type User struct {
	Model
	Email    sql.NullString `db:"email"`
	Password sql.NullString `db:"password"`
	Role     sql.NullString `db:"role"`
	LoggedOutAt JsonNullTime `db:"logged_out_at"`
}

func (p *PersistenceDataLayer) GetUserByEmail(email string) (*User, error) {
	user := new(User)
	row := p.GetConn().QueryRowx(`select * from users where email = ?`, email)
	err := row.StructScan(user)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PersistenceDataLayer) GetUserByID(id int64) (*User, error) {
	user := new(User)
	row := p.GetConn().QueryRowx(`SELECT * FROM users WHERE id=?`, id)
	err := row.StructScan(user)
	if err == sql.ErrNoRows {
		return nil, ErrNoData
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PersistenceDataLayer) CreateUser(email, password string) (int64, error){
	result, err := p.GetConn().Exec("insert into users(email, password) values (?, ?)", email, password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}