package mockdatalayer

import (
	"database/sql"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
)

type MockDataLayer struct {
	Accounts []*datalayer.Account
	Contacts []*datalayer.Contact
}

func New() *MockDataLayer {
	m := new(MockDataLayer)
	m.Accounts = accounts
	m.Contacts = contacts
	return m
}

var accounts = []*datalayer.Account {
	{
		Model:    datalayer.Model{
			ID:        1,
			CreatedAt: sql.NullTime{
				Time:  time.Date(2020, 04, 12, 12, 2, 0, 0, time.UTC),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Date(2020, 04, 19, 18, 28, 0, 0, time.UTC),
				Valid: true,
			},
			DeletedAt: sql.NullTime{},
		},
		Email:    sql.NullString{
			String: "subzero@dreamrealm.com",
			Valid:  true,
		},
		Password: sql.NullString{
			String: "$2a$10$NkTUeL6hkTRZ7M13tKYLqOmg7pAQaGPdpch9b5UoTSoO77MHjbPjm",
			Valid:  true,
		},
	},
	{
		Model:    datalayer.Model{
			ID:        2,
			CreatedAt: sql.NullTime{
				Time:  time.Date(2019, 2, 11, 16, 36, 59, 0, time.UTC),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Date(2020, 01, 31, 22, 1, 53, 0, time.UTC),
				Valid: true,
			},
			DeletedAt: sql.NullTime{},
		},
		Email:    sql.NullString{
			String: "reptile@netherrealm.com",
			Valid:  true,
		},
		Password: sql.NullString{
			String: "$2a$10$NkTUeL6hkTRZ7M13tKYLqOmg7pAQaGPdpch9b5UoTSoO77MHjbPjm",
			Valid:  true,
		},
	},
}


var contacts = []*datalayer.Contact {
	{
		Model:  datalayer.Model{
			ID:        1,
			CreatedAt: sql.NullTime{
				Time:  time.Date(2019, 6, 3, 8, 20, 56, 0, time.UTC),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Date(2020, 02, 29, 20, 0, 4, 0, time.UTC),
				Valid: true,
			},
			DeletedAt: sql.NullTime{},
		},
		Name:   "Shao Khan",
		Phone:  "0831111111",
		UserID: 0,
	},
	{
		Model:  datalayer.Model{
			ID:        2,
			CreatedAt: sql.NullTime{
				Time:  time.Date(2019, 6, 3, 8, 20, 56, 0, time.UTC),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Date(2020, 02, 29, 20, 0, 4, 0, time.UTC),
				Valid: true,
			},
			DeletedAt: sql.NullTime{},
		},
		Name:   "Reptile",
		Phone:  "0832222222",
		UserID: 0,
	},
	{
		Model:  datalayer.Model{
			ID:        3,
			CreatedAt: sql.NullTime{
				Time:  time.Date(2019, 6, 3, 8, 20, 56, 0, time.UTC),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Date(2020, 02, 29, 20, 0, 4, 0, time.UTC),
				Valid: true,
			},
			DeletedAt: sql.NullTime{},
		},
		Name:   "Scorpion",
		Phone:  "0832222222",
		UserID: 1,
	},
}