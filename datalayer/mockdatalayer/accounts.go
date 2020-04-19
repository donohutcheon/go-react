package mockdatalayer

import (
	"database/sql"
	"math"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
)


func (m *MockDataLayer) GetAccountByEmail(email string) (*datalayer.Account, error) {
	for _, account := range m.Accounts {
		if account.Email.Valid && email == account.Email.String {
			return account, nil
		}
	}

	return nil, datalayer.ErrNoData
}

func (m *MockDataLayer) GetAccountByID(id int64) (*datalayer.Account, error) {
	for _, account := range m.Accounts {
		if id == account.ID {
			return account, nil
		}
	}

	return nil, datalayer.ErrNoData
}

func (m *MockDataLayer) getNextAccountID() int64 {
	var maxID int64 = math.MaxInt64
	for _, account := range m.Accounts {
		if account.ID > maxID {
			maxID = account.ID
		}
	}

	return maxID + 1
}

func (m *MockDataLayer) CreateAccount(email, password string) (int64, error){
	account, err := m.GetAccountByEmail(email)
	if err != datalayer.ErrNoData {
		return 0, err
	}

	account = &datalayer.Account{
		Model:    datalayer.Model{
			ID:        m.getNextAccountID(),
			CreatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{},
			DeletedAt: sql.NullTime{},
		},
		Email:    sql.NullString{
			String: email,
			Valid:  true,
		},
		Password: sql.NullString{
			String: password,
			Valid:  true,
		},
	}

	m.Accounts = append(m.Accounts, account)

	return account.ID, nil
}