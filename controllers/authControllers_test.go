package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/donohutcheon/gowebserver/app"
	"github.com/donohutcheon/gowebserver/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	type AuthResponse struct {
		Message string            `json:"message"`
		Status  bool              `json:"status"`
		Token   app.TokenResponse `json:"token"`
	}
	tests := []struct {
		name       string
		request    []byte
		expResp    AuthResponse
		expStatus  int
		expTokenValid bool
	}{
		{
			name: "Success",
			request: []byte(`{"email": "subzero@dreamrealm.com", "password": "secret"}`),
			expResp: AuthResponse{
				Message: "Logged In",
				Status:  true,
			},
			expStatus: http.StatusOK,
			expTokenValid: true,
		},
		{
			name: "Non-existent User",
			request: []byte(`{"email": "skeletor@eternia.com", "password": "secret"}`),
			expResp: AuthResponse{
				Message: "Invalid login credentials",
				Status:  false,
			},
			expStatus: http.StatusForbidden,
		},
		{
			name: "Wrong Password",
			request: []byte(`{"email": "subzero@dreamrealm.com", "password": "wrong"}`),
			expResp: AuthResponse{
				Message: "Invalid login credentials",
				Status:  false,
			},
			expStatus: http.StatusForbidden,
		},
		{
			name: "Garbage Request",
			request: []byte(`garbage`),
			expResp: AuthResponse{
				Message: "Invalid request format",
				Status:  false,
			},
			expStatus: http.StatusBadRequest,
		},
	}

	route := "/auth/login"
	url, closer := setup(route, controllers.Authenticate)
	defer closer()
	ctx := context.Background()
	now := time.Now()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url + route, nil)
			assert.NoError(t, err)

			req.Body = ioutil.NopCloser(bytes.NewReader(test.request))
			defer req.Body.Close()

			cl := new(http.Client)
			res, err := cl.Do(req)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			gotResp := new(AuthResponse)
			err = json.Unmarshal(body, gotResp)
			require.NoError(t, err)

			assert.Equal(t, test.expStatus, res.StatusCode)
			assert.Equal(t, test.expResp.Message, gotResp.Message)
			assert.Equal(t, test.expResp.Status, gotResp.Status)
			if test.expTokenValid {
				assert.NotEmpty(t, gotResp.Token.AccessToken)
				assert.NotEmpty(t, gotResp.Token.RefreshToken)
				assert.Less(t, now.Unix(), gotResp.Token.ExpiresIn)
			} else {
				assert.Empty(t, gotResp.Token)
			}
		})
	}
}
