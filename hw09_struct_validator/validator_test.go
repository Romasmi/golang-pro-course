package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	Account struct {
		ID    string `validate:"type:uuid"`
		Email string `validate:"type:email"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Stats struct {
		Values []int `validate:"min:0|max:10"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          any
		expectedErr error
	}{
		{
			name: "valid user",
			in: User{
				ID:     "3c255929-3f89-4d4d-b10a-444c499d1a9b",
				Name:   "John",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "09876543210"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid age min",
			in: User{
				ID:    "3c255929-3f89-4d4d-b10a-444c499d1a9b",
				Age:   10,
				Email: "test@example.com",
				Role:  "admin",
			},
			expectedErr: ErrMinValue,
		},
		{
			name: "invalid age max",
			in: User{
				ID:    "12345678-1234-1234-1234-123456789012",
				Age:   60,
				Email: "test@example.com",
				Role:  "admin",
			},
			expectedErr: ErrMaxValue,
		},
		{
			name: "invalid email regexp",
			in: User{
				ID:    "3c255929-3f89-4d4d-b10a-444c499d1a9b",
				Age:   25,
				Email: "invalid-email",
				Role:  "admin",
			},
			expectedErr: ErrInvalidRegexp,
		},
		{
			name: "invalid email type",
			in: Account{
				ID:    "3c255929-3f89-4d4d-b10a-444c499d1a9b",
				Email: "invalid-email",
			},
			expectedErr: ErrInvalidEmail,
		},
		{
			name: "invalid uuid type",
			in: Account{
				ID:    "invalid-uuid",
				Email: "test@example.com",
			},
			expectedErr: ErrInvalidUUID,
		},
		{
			name: "valid account",
			in: Account{
				ID:    "3c255929-3f89-4d4d-b10a-444c499d1a9b",
				Email: "test@example.com",
			},
			expectedErr: nil,
		},
		{
			name: "valid app",
			in: App{
				Version: "1.2.3",
			},
			expectedErr: nil,
		},
		{
			name: "invalid app version",
			in: App{
				Version: "1.2",
			},
			expectedErr: ErrLengthMismatch,
		},
		{
			name: "token without tags",
			in: Token{
				Header:    []byte("h"),
				Payload:   []byte("p"),
				Signature: []byte("s"),
			},
			expectedErr: nil,
		},
		{
			name: "invalid role in",
			in: User{
				ID:    "3c255929-3f89-4d4d-b10a-444c499d1a9b",
				Age:   25,
				Email: "test@example.com",
				Role:  "guest",
			},
			expectedErr: ErrNotInList,
		},
		{
			name: "invalid phones len",
			in: User{
				ID:     "3c255929-3f89-4d4d-b10a-444c499d1a9b",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"123"},
			},
			expectedErr: ErrLengthMismatch,
		},
		{
			name: "pointer to struct",
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "John",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "09876543210"},
			},
			expectedErr: nil,
		},
		{
			name: "in for numbers",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "invalid in for numbers",
			in: Response{
				Code: 403,
				Body: "Forbidden",
			},
			expectedErr: ErrNotInList,
		},
		{
			name: "multiple errors",
			in: User{
				Age:   10,    // ErrMinValue
				Email: "bad", // ErrInvalidRegexp
			},
			expectedErr: ErrMinValue,
		},
		{
			name: "combined rules valid",
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			name: "combined rules invalid",
			in: User{
				Age: 60, // max is 50
			},
			expectedErr: ErrMaxValue,
		},
		{
			name: "slice of ints valid",
			in: Stats{
				Values: []int{0, 5, 10},
			},
			expectedErr: nil,
		},
		{
			name: "slice of ints invalid",
			in: Stats{
				Values: []int{0, 15, 10},
			},
			expectedErr: ErrMaxValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				if err != nil && tt.name != "not a struct" {
					t.Errorf("Validate() error = %v, expectedErr nil", err)
				}
				if tt.name == "not a struct" && err == nil {
					t.Errorf("Validate() should return error for non-struct")
				}
			} else {
				if err == nil {
					t.Errorf("Validate() error nil, expectedErr %v", tt.expectedErr)
				} else if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Validate() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			}
		})
	}

	err := Validate("not a struct")
	if err == nil {
		t.Errorf("Validate() error nil, expectedErr %v", ErrNotStruct)
	} else if !errors.Is(err, ErrNotStruct) {
		t.Errorf("Validate() error = %v, expectedErr %v", err, ErrNotStruct)
	}
}
