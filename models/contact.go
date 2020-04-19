package models

import (
	"log"

	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/datalayer"
)

type Contact struct {
	datalayer.Model
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	UserID    int64  `json:"user_id"` //The user that this contact belongs to
	dataLayer datalayer.DataLayer
}

func NewContact(dataLayer *datalayer.DataLayer) *Contact {
	contact := new(Contact)
	contact.dataLayer = *dataLayer
	return contact
}

func (c *Contact) convert(contact datalayer.Contact) {
	c.ID = contact.ID
	c.CreatedAt = contact.CreatedAt
	c.UpdatedAt = contact.UpdatedAt
	c.DeletedAt = contact.DeletedAt
	c.Name = contact.Name
	c.Phone = contact.Phone
}

// TODO: return errors
func (c *Contact) validate() (response.Response, bool) {
	if c.Name == "" {
		return response.New(false, "Contact name should be on the payload"), false
	}

	if c.Phone == "" {
		return response.New(false, "Phone number should be on the payload"), false
	}

	if c.UserID <= 0 {
		return response.New(false, "User is not recognized"), false
	}

	//All the required parameters are present
	return response.New(true, "success"), true
}

func (c *Contact) Create() (response.Response, error) {
	// TODO: c.Validate() to return an error
	if resp, ok := c.validate(); !ok {
		return resp, nil
	}

	dl := c.dataLayer
	id, err := dl.CreateContact(c.Name, c.Phone, c.UserID)
	if err != nil {
		// TODO: remove logging
		log.Fatal(err)
		return nil, err
	}

	contact, err := dl.GetContactByID(id)
	if err != nil {
		return nil, err
	}

	resp := response.New(true, "success")
	err = resp.Set("contact", contact)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Contact) GetContact(id int64) (*Contact, error) {
	dl := c.dataLayer
	dbContact, err := dl.GetContactByID(id)
	if err == datalayer.ErrNoData {
		return nil, err // TODO: return proper error with code
	}

	contact := new(Contact)
	contact.convert(*dbContact)

	return contact, nil
}

func (c *Contact) GetContacts(userID int64) ([]*Contact, error) {
	dl := c.dataLayer
	contacts := make([]*Contact, 0)

	dbContacts, err := dl.GetContactsByUserID(userID)
	if err == datalayer.ErrNoData {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	for _, dbContact := range dbContacts {
		contact := new(Contact)
		contact.convert(*dbContact)
		contacts = append(contacts, contact)
	}

	return contacts, err
}
