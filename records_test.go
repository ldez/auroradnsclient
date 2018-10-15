package auroradnsclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetRecords(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	zoneID := "c56a4180-zone-1"

	handleAPI(mux, "/zones/c56a4180-zone-1/records", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, `[
        {
          "id": "aaa",
          "type": "TXT",
          "name": "foo.com",
          "ttl": 300
        },
        {
          "id": "bbb",
          "type": "TXT",
          "name": "bar.com",
          "ttl": 600
        }
      ]`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	records, resp, err := client.GetRecords(zoneID)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := []Record{
		{ID: "aaa", RecordType: "TXT", Name: "foo.com", Content: "", TTL: 300},
		{ID: "bbb", RecordType: "TXT", Name: "bar.com", Content: "", TTL: 600},
	}
	assert.Equal(t, expected, records)
}

func TestClient_GetRecords_error(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	zoneID := "c56a4180-zone-1"

	handleAPI(mux, "/zones/c56a4180-zone-1/records", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
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

	records, resp, err := client.GetRecords(zoneID)
	require.EqualError(t, err, "AuthenticationRequiredError - Failed to parse Authorization header")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Nil(t, records)
}

func TestClient_CreateRecord(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	zoneID := "c56a4180-zone-2"

	handleAPI(mux, "/zones/c56a4180-zone-2/records", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if string(reqBody) != `{"type":"TXT","name":"foo","content":"w6uP8Tcg6K2QR905Rms8iXTlksL6OD1KOWBxTK7wxPI","ttl":300}` {
			http.Error(w, fmt.Sprintf("invalid request body: %s", string(reqBody)), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprintf(w, `{
				"id":   "c56a4180-record-1",
				"type": "TXT",
				"name": "foo",
				"ttl":  300
			}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	record := Record{
		RecordType: "TXT",
		Name:       "foo",
		Content:    "w6uP8Tcg6K2QR905Rms8iXTlksL6OD1KOWBxTK7wxPI",
		TTL:        300,
	}

	newRecord, resp, err := client.CreateRecord(zoneID, record)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	expected := &Record{
		ID:         "c56a4180-record-1",
		RecordType: "TXT",
		Name:       "foo",
		Content:    "",
		TTL:        300,
	}
	assert.Equal(t, expected, newRecord)
}

func TestClient_CreateRecord_error(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	zoneID := "c56a4180-zone-2"

	handleAPI(mux, "/zones/c56a4180-zone-2/records", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if string(reqBody) != `{"type":"TXT","name":"foo","content":"w6uP8Tcg6K2QR905Rms8iXTlksL6OD1KOWBxTK7wxPI","ttl":300}` {
			http.Error(w, fmt.Sprintf("invalid request body: %s", string(reqBody)), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		_, err = fmt.Fprintf(w, `{
  			"error": "AuthenticationRequiredError",
  			"errormsg": "Failed to parse Authorization header"
			}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	record := Record{
		RecordType: "TXT",
		Name:       "foo",
		Content:    "w6uP8Tcg6K2QR905Rms8iXTlksL6OD1KOWBxTK7wxPI",
		TTL:        300,
	}

	newRecord, resp, err := client.CreateRecord(zoneID, record)
	require.EqualError(t, err, "AuthenticationRequiredError - Failed to parse Authorization header")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Nil(t, newRecord)
}

func TestClient_RemoveRecord(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	zoneID := "c56a4180-zone-3"
	recordID := "c56a4180-record-2"

	handleAPI(mux, "/zones/c56a4180-zone-3/records/c56a4180-record-2", http.MethodDelete, nil)

	result, resp, err := client.RemoveRecord(zoneID, recordID)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result)
}

func TestClient_RemoveRecord_error(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	zoneID := "c56a4180-zone-3"
	recordID := "c56a4180-record-2"

	handleAPI(mux, "/zones/c56a4180-zone-3/records/c56a4180-record-2", http.MethodDelete, func(w http.ResponseWriter, r *http.Request) {
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

	result, resp, err := client.RemoveRecord(zoneID, recordID)
	require.EqualError(t, err, "AuthenticationRequiredError - Failed to parse Authorization header")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.False(t, result)
}
