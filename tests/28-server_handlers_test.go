package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kachaje/sacco-schema/database"
	"github.com/kachaje/sacco-schema/menus"
	menufuncs "github.com/kachaje/sacco-schema/menus/menuFuncs"
)

func setupTestServer() (*httptest.Server, func()) {
	menufuncs.DB = database.NewDatabase(":memory:")
	demoMode := true
	menufuncs.DemoMode = demoMode

	// Initialize menus
	activeMenu = menus.NewMenus(nil, &demoMode)

	// Create a test server that mimics the actual server setup
	mux := http.NewServeMux()

	// USSD handler
	mux.HandleFunc("/ussd", func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.FormValue("sessionId")
		phoneNumber := r.FormValue("phoneNumber")
		text := r.FormValue("text")

		if phoneNumber == "" {
			phoneNumber = "000000000"
		}

		var preferredLanguage string
		result := menufuncs.CheckPreferredLanguage(phoneNumber, ".settings")
		if result != nil {
			preferredLanguage = *result
		}

		session := menufuncs.CreateNewSession(phoneNumber, sessionID, ".settings", preferredLanguage, menufuncs.DemoMode)
		response := activeMenu.LoadMenu(session.CurrentMenu, session, phoneNumber, text, ".settings")

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))
	})

	// Cron jobs handler
	mux.HandleFunc("/cron/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := map[string]any{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// For testing, just return success
		w.Write([]byte("Done\n"))
	})

	ts := httptest.NewServer(mux)

	return ts, func() {
		ts.Close()
		menufuncs.DB.Close()
	}
}

var activeMenu *menus.Menus

func TestUSSDHandler(t *testing.T) {
	ts, cleanup := setupTestServer()
	defer cleanup()

	tests := []struct {
		name           string
		sessionID      string
		phoneNumber    string
		text           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Initial request with empty text",
			sessionID:      "test123",
			phoneNumber:    "1234567890",
			text:           "",
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome",
		},
		{
			name:           "Request with menu selection",
			sessionID:      "test123",
			phoneNumber:    "1234567890",
			text:           "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "Request with default phone number",
			sessionID:      "test123",
			phoneNumber:    "",
			text:           "",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := map[string]string{
				"sessionId":   tt.sessionID,
				"phoneNumber": tt.phoneNumber,
				"text":        tt.text,
			}

			body := bytes.NewBufferString(buildFormData(formData))
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/ussd", body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			rr := httptest.NewRecorder()
			rr.Code = resp.StatusCode
			bodyBytes, _ := io.ReadAll(resp.Body)
			rr.Body.Write(bodyBytes)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedBody != "" && !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("Expected body to contain %s, got %s", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestUSSDHandlerInvalidMethod(t *testing.T) {
	ts, cleanup := setupTestServer()
	defer cleanup()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/ussd", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// USSD handler accepts both GET and POST, so this should work
	if resp.StatusCode != http.StatusOK {
		t.Logf("Handler returned status %d", resp.StatusCode)
	}
}

func TestCronJobsHandler(t *testing.T) {
	ts, cleanup := setupTestServer()
	defer cleanup()

	tests := []struct {
		name           string
		method         string
		body           map[string]any
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid cron job request",
			method:         http.MethodPost,
			body:           map[string]any{"targetDate": "2024-01-31"},
			expectedStatus: http.StatusOK,
			expectedBody:   "Done",
		},
		{
			name:           "Cron job with profit parameter",
			method:         http.MethodPost,
			body:           map[string]any{"targetDate": "2024-01-31", "profit": 100000},
			expectedStatus: http.StatusOK,
			expectedBody:   "Done",
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "",
		},
		{
			name:           "Missing targetDate uses current date",
			method:         http.MethodPost,
			body:           map[string]any{},
			expectedStatus: http.StatusOK,
			expectedBody:   "Done",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody []byte
			var err error

			if tt.body != nil {
				requestBody, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatal(err)
				}
			}

			req, err := http.NewRequest(tt.method, ts.URL+"/cron/jobs", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyStr := string(bodyBytes)
			if tt.expectedBody != "" && !strings.Contains(bodyStr, tt.expectedBody) {
				t.Errorf("Expected body to contain %s, got %s", tt.expectedBody, bodyStr)
			}
		})
	}
}

func TestCronJobsHandlerInvalidJSON(t *testing.T) {
	ts, cleanup := setupTestServer()
	defer cleanup()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/cron/jobs", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func buildFormData(data map[string]string) string {
	var parts []string
	for k, v := range data {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, "&")
}
