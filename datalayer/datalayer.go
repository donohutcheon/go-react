package datalayer

type DataLayer interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	CreateUser(email, password string) (int64, error)

	CreateContact(name, phone string, userID int64) (int64, error)
	GetContactByID(id int64) (*Contact, error)
	GetContactsByUserID(userID int64) ([]*Contact, error)

	CreateCardTransaction(*CardTransaction) (int64, error)
	GetCardTransactionByID(id int64) (*CardTransaction, error)
	GetCardTransactionsByUserID(userID int64) ([]*CardTransaction, error)
}