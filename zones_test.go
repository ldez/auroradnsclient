package auroradnsclient

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetZones(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	handleAPI(mux, "/zones", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, `[
				{
					"id":   "c56a4180-65aa-42ec-a945-5fd21dec0538",
					"name": "example.com"
				}
			]`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	zones, resp, err := client.GetZones()
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := []Zone{{ID: "c56a4180-65aa-42ec-a945-5fd21dec0538", Name: "example.com"}}
	assert.Equal(t, expected, zones)
}

func TestClient_GetZones_error(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	handleAPI(mux, "/zones", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := fmt.Fprintf(w, `{
  			"error": "AuthenticationRequiredError",
  			"errormsg": "Failed to parse Authorization header"
			}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	zones, resp, err := client.GetZones()
	require.EqualError(t, err, "AuthenticationRequiredError - Failed to parse Authorization header")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Nil(t, zones)
}
