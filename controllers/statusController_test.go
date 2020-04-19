package controllers_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/donohutcheon/gowebserver/controllers"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatus(t *testing.T) {
	route := "/status"
	url, closer := setup(route, controllers.Status)
	defer closer()
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url + route, nil)
	assert.NoError(t, err)

	cl := new(http.Client)
	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)

	expResp := response.New(true, "Service is up")
	expected, err := json.Marshal(expResp)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(body))
}
