package controllers_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatus(t *testing.T) {
	url, _ := setup(t)

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url + "/status", nil)
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
