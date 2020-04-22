package mockdatalayer

import (
	"encoding/json"
	"github.com/donohutcheon/gowebserver/datalayer"
	"io/ioutil"
)

type MockDataLayer struct {
	Accounts []*datalayer.Account
	Contacts []*datalayer.Contact
}

func New() *MockDataLayer {
	m := new(MockDataLayer)

	return m
}

func (m *MockDataLayer) LoadAccountTestData(filename string) error{
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &m.Accounts)
	if err != nil {
		return err
	}

	return nil
}

func (m *MockDataLayer) LoadContactTestData(filename string) error{
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