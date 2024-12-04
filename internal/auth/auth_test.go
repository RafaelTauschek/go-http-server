package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headerValue   string
		expectedToken string
		expectedError bool
	}{
		{
			name:          "valid bearer token",
			headerValue:   "Bearer mytoken",
			expectedToken: "mytoken",
			expectedError: false,
		},
		{
			name:          "invalid bearer token",
			headerValue:   "Authentication header",
			expectedToken: "",
			expectedError: true,
		},
		{
			name:          "No token provided",
			headerValue:   "",
			expectedToken: "",
			expectedError: true,
		},
		{
			name:          "Too much arguments provded",
			headerValue:   "Bearer mytoken extra",
			expectedToken: "",
			expectedError: true,
		},
		{
			name:          "Too little arguments provided",
			headerValue:   "Bearer",
			expectedToken: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.headerValue != "" {
				headers.Add("Authorization", tt.headerValue)
			}
			token, err := GetBearerToken(headers)

			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if token != tt.expectedToken {
				t.Errorf("got token %q, want %q", token, tt.expectedToken)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	userID := uuid.New()

	tests := []struct {
		name          string
		setupToken    func() string
		expectedError bool
		expectedID    uuid.UUID
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, err := MakeJWT(userID, secret)
				if err != nil {
					t.Fatalf("failed to create test token: %v", err)
				}
				return token
			},
			expectedError: false,
			expectedID:    userID,
		},
		{
			name: "expired token",
			setupToken: func() string {
				token, err := MakeJWT(userID, secret)
				if err != nil {
					t.Fatalf("failed to create test token: %v", err)
				}
				return token
			},
			expectedError: true,
			expectedID:    uuid.Nil,
		},
		{
			name: "wrong secret",
			setupToken: func() string {
				token, err := MakeJWT(userID, wrongSecret)
				if err != nil {
					t.Fatalf("failed to create test token: %v", err)
				}
				return token
			},
			expectedError: true,
			expectedID:    uuid.Nil,
		},
		{
			name: "malformed token",
			setupToken: func() string {
				return "not.a.jwt"
			},
			expectedError: true,
			expectedID:    uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ValidateJWT(tt.setupToken(), secret)

			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if id != tt.expectedID {
				t.Errorf("got id %v, but expected %v", id, tt.expectedID)
			}
		})
	}
}
