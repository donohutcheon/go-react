package controllers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserControllerResponse struct {
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	User    models.User `json:"user"`
}

func TestGetCurrentUser(t *testing.T) {
	tests := []struct {
		name          string
		authRequest   []byte
		expLoginResp  AuthResponse
		expResp       UserControllerResponse
		expStatus     int
		expTokenValid bool
	}{
		{
			name: "Success",
			authRequest: []byte(`{"email": "subzero@dreamrealm.com", "password": "secret"}`),
			expLoginResp : AuthResponse{
				Message: "Logged In",
				Status:  true,
			},
			expResp: UserControllerResponse{
				Message: "success",
				Status:  true,
				User: models.User{
					Model:        datalayer.Model{
						ID:        0,
						CreatedAt: datalayer.JsonNullTime{
							NullTime : sql.NullTime{
								Time:  time.Now(),
								Valid: true,
							},
						},
					},
					Email:        "nightwolf@earthrealm.com",
					Roles:        []string{"ADMIN", "USER"},
					Settings:     models.Settings{
						ID:        0,
						ThemeName: "default",
					},
				},
			},
			expStatus: http.StatusOK,
			expTokenValid: true,
		},
	}

	url, _ := setup(t)
	ctx := context.Background()
	now := time.Now()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := new(http.Client)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/auth/login", nil)
			assert.NoError(t, err)
			req.Body = ioutil.NopCloser(bytes.NewReader(test.authRequest))
			defer req.Body.Close()
			res, err := cl.Do(req)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			gotAuthResp := new(AuthResponse)
			err = json.Unmarshal(body, gotAuthResp)
			require.NoError(t, err)

			assert.Equal(t, test.expStatus, res.StatusCode)
			assert.Equal(t, test.expLoginResp.Message, gotAuthResp.Message)
			assert.Equal(t, test.expLoginResp.Status, gotAuthResp.Status)
			if test.expTokenValid {
				require.NotEmpty(t, gotAuthResp.Token.AccessToken)
				require.NotEmpty(t, gotAuthResp.Token.RefreshToken)
				require.Less(t, now.Unix(), gotAuthResp.Token.ExpiresIn)
			} else {
				assert.Empty(t, gotAuthResp.Token)
			}

			req, err = http.NewRequestWithContext(ctx, http.MethodGet, url+"/users/current", nil)
			assert.NoError(t, err)
			req.Header.Add("Authorization", "Bearer "+gotAuthResp.Token.AccessToken)
			res, err = cl.Do(req)
			require.NoError(t, err)

			body, err = ioutil.ReadAll(res.Body)
			gotResp := new(UserControllerResponse)
			err = json.Unmarshal(body, gotResp)
			require.NoError(t, err)

			assert.Equal(t, test.expResp.Status, gotResp.Status)
			assert.Equal(t, test.expResp.Message, gotResp.Message)
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name              string
		authParameters    AuthParameters
		createUserReq     models.User
		expCreateUserResp UserControllerResponse
		expStatus         int
	}{
		{
			name: "Golden",
			authParameters: AuthParameters{
				authRequest: models.User{
					Email:    "jade@edenia.com",
					Password: "secret",
				},
				expHTTPStatus: http.StatusOK,
				expLoginResp: AuthResponse{
					Message: "Logged In",
					Status:  true,
				},
			},
			createUserReq: models.User{
				Email:     "jade@edenia.com",
				Password:  "secret",
			},
			expCreateUserResp : UserControllerResponse{
				Message: "User has been created",
				Status:  true,
				User: models.User{
					Email: "jade@edenia.com",
				},
			},
			expStatus: http.StatusOK,
		},
		{
			name: "Incomplete Email",
			createUserReq: models.User{
				Model:     datalayer.Model{},
				Email:     "",
				Password:  "secret",
			},
			expCreateUserResp : UserControllerResponse{
				Message: "Email address is required",
				Status:  false,
			},
			expStatus: http.StatusBadRequest,
		},
		{
			name: "Incomplete Password",
			createUserReq: models.User{
				Model:     datalayer.Model{},
				Email:     "sindel@dreamrealm.com",
			},
			expCreateUserResp : UserControllerResponse{
				Message: "Password is required",
				Status:  false,
			},
			expStatus: http.StatusBadRequest,
		},
	}
	
	url, _ := setup(t)
	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := new(http.Client)
			// Create new user.
			{
				req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/auth/sign-up", nil)

				assert.NoError(t, err)

				b, err := json.Marshal(test.createUserReq)
				require.NoError(t, err)

				req.Body = ioutil.NopCloser(bytes.NewReader(b))
				defer req.Body.Close()

				res, err := cl.Do(req)
				require.NoError(t, err)

				body, err := ioutil.ReadAll(res.Body)
				gotCreateUserResp := new(UserControllerResponse)
				err = json.Unmarshal(body, gotCreateUserResp)
				require.NoError(t, err)
				assert.Equal(t, test.expStatus, res.StatusCode)
				assert.Equal(t, test.expCreateUserResp.Message, gotCreateUserResp.Message)
				assert.Equal(t, test.expCreateUserResp.Status, gotCreateUserResp.Status)
				assert.Equal(t, test.expCreateUserResp.User.Email, gotCreateUserResp.User.Email)
			}

			// Exit test if the previous call failed.
			if !test.expCreateUserResp.Status {
				return
			}

			login(t, ctx, cl, url, test.authParameters)
		})
	}
}