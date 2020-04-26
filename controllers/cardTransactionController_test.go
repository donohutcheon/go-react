package controllers_test

import (
	"bytes"
	"context"
	"database/sql"
	"time"


	"encoding/json"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

type CreateCardTransactionControllerResponse struct {
	Message         string                 `json:"message"`
	Status          bool                   `json:"status"`
	CardTransaction models.CardTransaction `json:"cardTransaction"`
}

type GetCardTransactionControllerResponse struct {
	Message          string                   `json:"message"`
	Status           bool                     `json:"status"`
	CardTransactions []models.CardTransaction `json:"cardTransactions"`
}

type CreateCardTransactionParameters struct {
	skip          bool
	request       models.CardTransaction
	expResponse   CreateCardTransactionControllerResponse
	expHTTPStatus int
}

type GetCardTransactionParameters struct {
	skip          bool
	expResponse   GetCardTransactionControllerResponse
	expHTTPStatus int
}

func TestCardTransactions(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name                        string
		authParameters              AuthParameters
		createCardTransactionParams CreateCardTransactionParameters
		getCardTransactionParams GetCardTransactionParameters
	}{
		{
			name: "Golden",
			authParameters: AuthParameters{
				authRequest: models.User{
					Email:    "subzero@dreamrealm.com",
					Password: "secret",
				},
				expHTTPStatus: http.StatusOK,
				expLoginResp : AuthResponse{
					Message: "Logged In",
					Status:  true,
				},
			},
			createCardTransactionParams: CreateCardTransactionParameters{
				request: models.CardTransaction{
					Model:    datalayer.Model{},
					DateTime: time.Date(2020, 04, 25, 19, 46, 23, 33, time.UTC),
					Amount: models.CurrencyValue{
						Value: 400,
						Scale: 2,
					},
					CurrencyCode:         "ZAR",
					Reference:            "simulation",
					MerchantName:         "Dwelms en Dinges",
					MerchantCity:         "Hillbrow",
					MerchantCountryCode:  "ZA",
					MerchantCountryName:  "South Africa",
					MerchantCategoryCode: "contraband",
					MerchantCategoryName: "Contraband",
				},
				expResponse: CreateCardTransactionControllerResponse{
					Message: "success",
					Status:  true,
					CardTransaction: models.CardTransaction{
						Model: datalayer.Model{
							ID: 0,
							CreatedAt: datalayer.JsonNullTime{
								NullTime: sql.NullTime{
									Time:  testTime,
									Valid: true,
								},
							},
						},
						DateTime: time.Date(2020, 04, 25, 19, 46, 23, 33, time.UTC),
						Amount: models.CurrencyValue{
							Value: 400,
							Scale: 2,
						},
						CurrencyCode:         "ZAR",
						Reference:            "simulation",
						MerchantName:         "Dwelms en Dinges",
						MerchantCity:         "Hillbrow",
						MerchantCountryCode:  "ZA",
						MerchantCountryName:  "South Africa",
						MerchantCategoryCode: "contraband",
						MerchantCategoryName: "Contraband",
					},
				},
				expHTTPStatus: http.StatusOK,
			},
			getCardTransactionParams: GetCardTransactionParameters{
				expResponse: GetCardTransactionControllerResponse{
					Message: "success",
					Status:  true,
					CardTransactions: []models.CardTransaction{
						{
							Model: datalayer.Model{
								ID: 0,
								CreatedAt: datalayer.JsonNullTime{
									NullTime: sql.NullTime{
										Time:  testTime,
										Valid: true,
									},
								},
							},
							DateTime: time.Date(2020, 04, 25, 19, 46, 23, 33, time.UTC),
							Amount: models.CurrencyValue{
								Value: 400,
								Scale: 2,
							},
							CurrencyCode:         "ZAR",
							Reference:            "simulation",
							MerchantName:         "Dwelms en Dinges",
							MerchantCity:         "Hillbrow",
							MerchantCountryCode:  "ZA",
							MerchantCountryName:  "South Africa",
							MerchantCategoryCode: "contraband",
							MerchantCategoryName: "Contraband",
						},
					},
				},
				expHTTPStatus: http.StatusOK,
			},
		},
		{
			name: "No data",
			authParameters: AuthParameters{
				authRequest: models.User{
					Email:    "subzero@dreamrealm.com",
					Password: "secret",
				},
				expHTTPStatus: http.StatusOK,
				expLoginResp : AuthResponse{
					Message: "Logged In",
					Status:  true,
				},
			},
			createCardTransactionParams: CreateCardTransactionParameters{
				skip: true,
			},
			getCardTransactionParams: GetCardTransactionParameters{
				expResponse: GetCardTransactionControllerResponse{
					Message: "success",
					Status:  true,
					CardTransactions: []models.CardTransaction{},
				},
				expHTTPStatus: http.StatusOK,
			},
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := new(http.Client)
			url, _ := setup(t)
			gotAuthResp := login(t, ctx, cl, url, test.authParameters)
			createCardTransaction(t, ctx, cl, url, gotAuthResp, &test.createCardTransactionParams)
			getCardTransactions(t, ctx, cl, url, gotAuthResp, &test.getCardTransactionParams)
		})
	}
}

func createCardTransaction(t *testing.T, ctx context.Context, cl *http.Client,
	url string, auth *AuthResponse, params *CreateCardTransactionParameters) {
	if params.skip {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/card-transactions/new", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "Bearer " + auth.Token.AccessToken)

	b, err := json.Marshal(params.request)
	require.NoError(t, err)

	req.Body = ioutil.NopCloser(bytes.NewReader(b))
	defer req.Body.Close()

	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	gotResp := new(CreateCardTransactionControllerResponse)
	err = json.Unmarshal(body, gotResp)
	require.NoError(t, err)

	assert.Equal(t, params.expResponse.Status, gotResp.Status)
	assert.Equal(t, params.expResponse.Message, gotResp.Message)
	assert.Equal(t, params.expResponse.CardTransaction.Amount.Value, gotResp.CardTransaction.Amount.Value)
}

func getCardTransactions(t *testing.T, ctx context.Context, cl *http.Client,
	url string, auth *AuthResponse, params *GetCardTransactionParameters) {
	if params.skip {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/me/card-transactions", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "Bearer " + auth.Token.AccessToken)

	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	gotResp := new(GetCardTransactionControllerResponse)
	err = json.Unmarshal(body, gotResp)
	require.NoError(t, err)

	assert.Equal(t, params.expResponse.Status, gotResp.Status)
	assert.Equal(t, params.expResponse.Message, gotResp.Message)
	require.Equal(t, len(params.expResponse.CardTransactions), len(gotResp.CardTransactions))
	for i, x := range params.expResponse.CardTransactions {
		// Negate datetime fields we can't control.
		x.CreatedAt = datalayer.JsonNullTime{}
		gotResp.CardTransactions[i].CreatedAt = datalayer.JsonNullTime{}

		assert.Equal(t, x, gotResp.CardTransactions[i])
	}
}