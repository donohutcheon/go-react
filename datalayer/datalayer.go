package datalayer

type DataLayer interface {
	GetAccountByEmail(email string) (*Account, error)
	GetAccountByID(id int64) (*Account, error)
	CreateAccount(email, password string) (int64, error)

	CreateContact(name, phone string, userID int64) (int64, error)
	GetContactByID(id int64) (*Contact, error)
	GetContactsByUserID(userID int64) ([]*Contact, error)
}