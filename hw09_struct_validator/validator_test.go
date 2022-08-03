package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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

	MinMax struct {
		Age int `validate:"min:18|max:50"`
	}

	In struct {
		Role   UserRole `validate:"in:admin,stuff"`
		Status string   `validate:"in:new,created"`
	}
)

func TestValidateError(t *testing.T) {
	uRole := UserRole("www")
	tests := []struct {
		in       interface{}
		expected error
	}{
		{
			in:       "ErrNotStruct",
			expected: ErrNotStruct,
		},
		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expected: nil,
		},
		{
			in:       Response{Code: 201},
			expected: ValidationErrors{{Field: "Code", Err: ErrValidationContains}},
		},
		{
			in:       App{Version: ""},
			expected: ValidationErrors{{Field: "Version", Err: ErrValidationLength}},
		},
		{
			in:       MinMax{Age: 17},
			expected: ValidationErrors{{Field: "Age", Err: ErrValidationMinimum}},
		},
		{
			in:       MinMax{Age: 51},
			expected: ValidationErrors{{Field: "Age", Err: ErrValidationMaximum}},
		},
		{
			in:       In{Role: uRole, Status: "new"},
			expected: ValidationErrors{{Field: "Role", Err: ErrValidationContains}},
		},
		{
			in: In{Role: uRole, Status: "new1"},
			expected: ValidationErrors{
				{Field: "Role", Err: ErrValidationContains},
				{Field: "Status", Err: ErrValidationContains},
			},
		},
		{
			in: User{
				ID:     "asd",
				Name:   "Name",
				Age:    17,
				Email:  "asd.asd",
				Role:   uRole,
				Phones: []string{"8 8005553535 "},
				meta:   []byte("{}"),
			},
			expected: ValidationErrors{
				{Field: "ID", Err: ErrValidationLength},
				{Field: "Age", Err: ErrValidationMinimum},
				{Field: "Email", Err: ErrValidationRegexp},
				{Field: "Role", Err: ErrValidationContains},
				{Field: "Phones", Err: ErrValidationLength},
			},
		},
	}

	for i, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, Validate(tt.in))
		})
	}
}

func TestValidateSuccess(t *testing.T) {
	uRole := UserRole("admin")
	tests := []struct {
		in interface{}
	}{
		{
			in: App{Version: "0.0.1"},
		},
		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
		},
		{
			in: Response{Code: 200},
		},
		{
			in: MinMax{Age: 19},
		},
		{
			in: MinMax{Age: 50},
		},
		{
			in: In{Role: uRole, Status: "created"},
		},
		{
			in: User{
				ID:     "5f04797b-e4ea-4ede-91c7-576a42d1f764",
				Name:   "Name",
				Age:    21,
				Email:  "test@test.ru",
				Role:   uRole,
				Phones: []string{"88005553535", "89995553535"},
				meta:   []byte("{}"),
			},
		},
	}

	for i, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			require.NoError(t, Validate(tt.in))
		})
	}
}
