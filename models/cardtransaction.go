package models

import (
	"log"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
)

type CurrencyValue struct {
	Value int64 `json:"value"`
	Scale int   `json:"scale"`
}

type CardTransaction struct {
	datalayer.Model
	DateTime             time.Time     `json:"dateTime"`
	Amount               CurrencyValue `json:"amount"`
	CurrencyCode         string        `json:"currencyCode"`
	Reference            string        `json:"reference"`
	MerchantName         string        `json:"merchantName"`
	MerchantCity         string        `json:"merchantCity"`
	MerchantCountryCode  string        `json:"merchantCountryCode"`
	MerchantCountryName  string        `json:"merchantCountryName"`
	MerchantCategoryCode string        `json:"merchantCategoryCode"`
	MerchantCategoryName string        `json:"merchantCategoryName"`
	UserID               int64         `json:"user_id"`
	dataLayer            datalayer.DataLayer
}

func NewCardTransaction(dataLayer *datalayer.DataLayer) *CardTransaction {
	cardTransaction := new(CardTransaction)
	cardTransaction.dataLayer = *dataLayer
	return cardTransaction
}

func newFromDBCardTransaction(cardTransaction *datalayer.CardTransaction) *CardTransaction{
	c := new(CardTransaction)
	c.ID = cardTransaction.ID
	c.CreatedAt = cardTransaction.CreatedAt
	c.UpdatedAt = cardTransaction.UpdatedAt
	c.DeletedAt = cardTransaction.DeletedAt
	c.DateTime = cardTransaction.DateTime
	c.Amount.Value = cardTransaction.Amount
	c.Amount.Scale = cardTransaction.CurrencyScale
	c.CurrencyCode = cardTransaction.CurrencyCode
	c.Reference = cardTransaction.Reference
	c.MerchantName = cardTransaction.MerchantName
	c.MerchantCity = cardTransaction.MerchantCity
	c.MerchantCountryCode = cardTransaction.MerchantCountryCode
	c.MerchantCountryName = cardTransaction.MerchantCountryName
	c.MerchantCategoryCode = cardTransaction.MerchantCategoryCode
	c.MerchantCategoryName = cardTransaction.MerchantCategoryName
	return c
}

func (c *CardTransaction) convertToDB() *datalayer.CardTransaction {
	cardTransaction := new(datalayer.CardTransaction)
	cardTransaction.ID = c.ID
	cardTransaction.CreatedAt = c.CreatedAt
	cardTransaction.UpdatedAt = c.UpdatedAt
	cardTransaction.DeletedAt = c.DeletedAt
	cardTransaction.DateTime = c.DateTime
	cardTransaction.Amount = c.Amount.Value
	cardTransaction.CurrencyScale = c.Amount.Scale
	cardTransaction.CurrencyCode = c.CurrencyCode
	cardTransaction.Reference = c.Reference
	cardTransaction.MerchantName = c.MerchantName
	cardTransaction.MerchantCity = c.MerchantCity
	cardTransaction.MerchantCountryCode = c.MerchantCountryCode
	cardTransaction.MerchantCountryName = c.MerchantCountryName
	cardTransaction.MerchantCategoryCode = c.MerchantCategoryCode
	cardTransaction.MerchantCategoryName = c.MerchantCategoryName
	cardTransaction.UserID = c.UserID
	return cardTransaction
}

// TODO: return errors
func (c *CardTransaction) validate() error {
	// TODO: Validate no empty fields

	if c.UserID <= 0 {
		return ErrUserDoesNotExist
	}

	if len(c.CurrencyCode) == 0 {
		return ErrValidationFailed
	}

	if len(c.MerchantName) == 0 {
		return ErrValidationFailed
	}

	//All the required parameters are present
	return nil
}

func (c *CardTransaction) CreateCardTransaction() (*CardTransaction, error) {
	// TODO: c.Validate() to return an error
	err := c.validate()
	if err != nil {
		return nil, err
	}

	dbCardTransaction := c.convertToDB()

	dl := c.dataLayer
	id, err := dl.CreateCardTransaction(dbCardTransaction)
	if err != nil {
		// TODO: remove logging
		log.Fatal(err)
		return nil, err
	}

	dbCardTransaction, err = dl.GetCardTransactionByID(id)
	if err != nil {
		return nil, err
	}

	data := newFromDBCardTransaction(dbCardTransaction)

	return data, nil
}

func (c *CardTransaction) GetCardTransaction(id int64) (*CardTransaction, error) {
	dl := c.dataLayer
	dbCardTransaction, err := dl.GetCardTransactionByID(id)
	if err == datalayer.ErrNoData {
		return nil, err // TODO: return proper error with code
	}

	cardTransaction := newFromDBCardTransaction(dbCardTransaction)

	return cardTransaction, nil
}

func (c *CardTransaction) GetCardTransactionsByUserID(userID int64) ([]*CardTransaction, error) {
	dl := c.dataLayer
	cardTransactions := make([]*CardTransaction, 0)

	dbCardTransactions, err := dl.GetCardTransactionsByUserID(userID)
	if err == datalayer.ErrNoData {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	for _, dbCardTransaction := range dbCardTransactions {
		cardTransaction := newFromDBCardTransaction(dbCardTransaction)
		cardTransactions = append(cardTransactions, cardTransaction)
	}

	return cardTransactions, err
}
