package services

import (
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/herdiansc/go-cms/models"
)

type mockJsonDecoder struct {
	err error
}

func (m mockJsonDecoder) Decode(v any) error {
	json.Unmarshal([]byte(`{"msisdn":"08221876567","name":"herdian","password":"adalah","username":"herdiansc"}`), &v)
	return m.err
}

type mockRequestValidator struct {
	err error
}

func (m mockRequestValidator) Struct(s interface{}) error {
	return m.err
}

type mockPasswordHasher struct {
	h   string
	err error
}

func (m mockPasswordHasher) HashPassword(password string) (string, error) {
	return m.h, m.err
}

type mockAuthCreator struct {
	err error
}

func (m mockAuthCreator) Create(auth models.Auth) error {
	return m.err
}

type mockRequestDecoder struct {
	r io.Reader
	e error
}

func newMockRequestDecoder(r io.Reader, e error) mockRequestDecoder {
	return mockRequestDecoder{
		r: r,
		e: e,
	}
}

func (m mockRequestDecoder) Decode(v any) error {
	return m.e
}

var (
	mockSuccessJsonDecoder = mockJsonDecoder{
		err: nil,
	}
	mockFailedJsonDecoder = mockJsonDecoder{
		err: errors.New("error"),
	}
	mockSuccessRequestValidator = mockRequestValidator{
		err: nil,
	}
	mockFailedRequestValidator = mockRequestValidator{
		err: errors.New("error"),
	}
	mockSuccessPasswordHasher = mockPasswordHasher{
		h:   "hashed",
		err: nil,
	}
	mockFailedPasswordHasher = mockPasswordHasher{
		h:   "",
		err: errors.New("error"),
	}
	mockSuccessAuthCreator = mockAuthCreator{
		err: nil,
	}
	mockFailedAuthCreator = mockAuthCreator{
		err: errors.New("error"),
	}
)

func TestRegistrationService_Register(t *testing.T) {
	cases := []struct {
		name      string
		dec       mockJsonDecoder
		validator mockRequestValidator
		hasher    mockPasswordHasher
		repo      mockAuthCreator
		want      int
	}{
		{
			name:      "Positive",
			dec:       mockSuccessJsonDecoder,
			validator: mockSuccessRequestValidator,
			hasher:    mockSuccessPasswordHasher,
			repo:      mockSuccessAuthCreator,
			want:      201,
		},
		{
			name:      "Negative: Failed decode body",
			dec:       mockFailedJsonDecoder,
			validator: mockSuccessRequestValidator,
			hasher:    mockSuccessPasswordHasher,
			repo:      mockSuccessAuthCreator,
			want:      400,
		},
		{
			name:      "Negative: Failed to validate request",
			dec:       mockSuccessJsonDecoder,
			validator: mockFailedRequestValidator,
			hasher:    mockSuccessPasswordHasher,
			repo:      mockSuccessAuthCreator,
			want:      400,
		},
		{
			name:      "Negative: Failed hash password",
			dec:       mockSuccessJsonDecoder,
			validator: mockSuccessRequestValidator,
			hasher:    mockFailedPasswordHasher,
			repo:      mockSuccessAuthCreator,
			want:      500,
		},
		{
			name:      "Negative: Failed to save data",
			dec:       mockSuccessJsonDecoder,
			validator: mockSuccessRequestValidator,
			hasher:    mockSuccessPasswordHasher,
			repo:      mockFailedAuthCreator,
			want:      500,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewRegistrationServices(tt.dec, tt.validator, tt.hasher, tt.repo)
			code, _ := svc.Register()
			if code != tt.want {
				t.Errorf("Expected resp to be %q but it was %q", tt.want, code)
			}
		})
	}
}
