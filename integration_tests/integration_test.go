package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/herdiansc/go-cms/handlers"
	"github.com/herdiansc/go-cms/middlewares"
	"github.com/herdiansc/go-cms/models"
	"github.com/hhkbp2/testify/require"
	"gorm.io/gorm"
)

var testDBInstance *gorm.DB

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDBInstance = testDB.DB
	defer testDB.TearDown()

	os.Exit(m.Run())
}

func TestAuthRegister(t *testing.T) {
	body := []byte(`{
		"password": "123",
		"role": "WRITER",
		"username": "hdn"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.NewAuthHandler(testDBInstance).Register)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}
}

func TestAuthLogin(t *testing.T) {
	body := []byte(`{
		"password": "123",
		"username": "hdn"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.NewAuthHandler(testDBInstance).Login)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestAuthProfile(t *testing.T) {
	body := []byte(`{
		"password": "123",
		"username": "hdn"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.NewAuthHandler(testDBInstance).Login)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)

	token := response.Data.(map[string]interface{})["token"]

	req, err = http.NewRequest(http.MethodGet, "/auth/profile", bytes.NewBuffer(body))
	require.NoError(t, err)

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))

	rr = httptest.NewRecorder()
	privateHandler := middlewares.Authenticate(http.HandlerFunc(handlers.NewAuthHandler(testDBInstance).GetProfile))

	privateHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

}

func TestArticleCreate(t *testing.T) {
	body := []byte(`{
		"password": "123",
		"username": "hdn"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.NewAuthHandler(testDBInstance).Login)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)

	token := response.Data.(map[string]interface{})["token"]

	articleBody := []byte(`{
		"content": "content",
		"status": "DRAFT",
		"tags": [
			"testing"
		],
		"title": "title"
	}`)
	req, err = http.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(articleBody))
	require.NoError(t, err)

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))

	rr = httptest.NewRecorder()
	privateHandler := middlewares.Authenticate(http.HandlerFunc(handlers.NewArticleHandler(testDBInstance).Create))

	privateHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestArticleList(t *testing.T) {
	body := []byte(`{
		"password": "123",
		"username": "hdn"
	}`)

	req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.NewAuthHandler(testDBInstance).Login)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)

	token := response.Data.(map[string]interface{})["token"]

	req, err = http.NewRequest(http.MethodGet, "/articles", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))

	rr = httptest.NewRecorder()
	privateHandler := middlewares.Authenticate(http.HandlerFunc(handlers.NewArticleHandler(testDBInstance).List))

	privateHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}
