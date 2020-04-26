package mockdatalayer

import (
	"encoding/json"
	"github.com/donohutcheon/gowebserver/datalayer"
	"io/ioutil"
)

type MockDataLayer struct {
	Users    []*datalayer.User
	Contacts []*datalayer.Contact
	CardTransactions []*datalayer.CardTransaction
	usersFilename string
	contactsFilename string
}

func New() *MockDataLayer {
	m := new(MockDataLayer)

	return m
}

func (m *MockDataLayer) ResetAndReload() error {
	m.Users = m.Users[:0]
	err := m.LoadUserTestData(m.usersFilename)
	if err != nil {
		return err
	}

	m.Contacts = m.Contacts[:0]
	m.LoadContactTestData(m.contactsFilename)
	if err != nil {
		return err
	}

	m.CardTransactions = m.CardTransactions[:0]

	return nil
}

func (m *MockDataLayer) LoadUserTestData(filename string) error{
	m.usersFilename = filename
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &m.Users)
	if err != nil {
		return err
	}

	return nil
}

func (m *MockDataLayer) LoadContactTestData(filename string) error{
	m.contactsFilename = filename
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &m.Contacts)
	if err != nil {
		return err
	}

	return nil
}

func (m *MockDataLayer) LoadCardTransactionTestData(filename string) error{
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &m.CardTransactions)
	if err != nil {
		return err
	}

	return nil
}