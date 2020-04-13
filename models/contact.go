package models

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.com/donohutcheon/gowebserver/controllers/response"
)

type Contact struct {
	Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserID int64  `json:"user_id"` //The user that this contact belongs to
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

	result, err := GetConn().Exec("insert into contacts(name, phone, user_id) values (?, ?, ?)", c.Name, c.Phone, c.UserID)
	if err != nil {
		// TODO: remove logging
		log.Fatal(err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		// TODO: remove logging
		log.Fatal(err)
		return nil, err
	}
	contact := GetContact(id)

	resp := response.New(true, "success")
	err = resp.Set("contact", contact)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetContact(id int64) *Contact {
	var createdAt, updatedAt, deletedAt sql.NullTime
	contact := &Contact{}

	err := GetConn().QueryRow(`SELECT id, user_id, name, phone, created_at, updated_at, deleted_at FROM contacts WHERE id=?`, id).Scan(&contact.ID, &contact.UserID, &contact.Name, &contact.Phone, &createdAt, &updatedAt, &deletedAt)
	if err == sql.ErrNoRows {
		fmt.Println(false, "Contact does not exist.")
		return nil
	} else if err != nil {
		fmt.Printf("Failed to query contacts for user ID [%d] from database", id)
		return nil
	}

	if createdAt.Valid {
		contact.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		contact.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		contact.DeletedAt = &deletedAt.Time
	}

	return contact
}

func GetContacts(userID int64) []*Contact {
	contacts := make([]*Contact, 0)
	var createdAt, updatedAt, deletedAt sql.NullTime

	rows, err := GetConn().Query(`SELECT id, name, phone, created_at, updated_at, deleted_at FROM contacts WHERE user_id=?`, userID)
	if err == sql.ErrNoRows {
		fmt.Println(false, "User account does not exist. Please re-login")
		return nil
	} else if err != nil {
		fmt.Printf("Failed to query contacts for user ID [%d] from database", userID)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var contact Contact
		contact.UserID = userID
		rows.Scan(&contact.ID, &contact.Name, &contact.Phone, &createdAt, &updatedAt, &deletedAt)
		if createdAt.Valid {
			contact.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			contact.UpdatedAt = &updatedAt.Time
		}
		if deletedAt.Valid {
			contact.DeletedAt = &deletedAt.Time
		}
		contacts = append(contacts, &contact)
	}

	return contacts
}
