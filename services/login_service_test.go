package services

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/herdiansc/go-cms/models"
)

type mockPasswordComparer struct {
	r bool
}

func (m mockPasswordComparer) VerifyPassword(password, hash string) bool {
	return m.r
}

type mockAuthFinder struct {
	a models.Auth
	e error
}

func (m mockAuthFinder) FindByUsername(username string) (models.Auth, error) {
	return m.a, m.e
}

type mockSigner func(method jwt.SigningMethod, claim jwt.Claims, opts ...jwt.TokenOption) *jwt.Token

var (
	mockSuccessPasswordComparer = mockPasswordComparer{
		r: true,
	}
	mockFailedPasswordComparer = mockPasswordComparer{
		r: false,
	}
	mockSuccessAuthFinder = mockAuthFinder{
		a: models.Auth{},
		e: nil,
	}
	mockFailedAuthFinder = mockAuthFinder{
		a: models.Auth{},
		e: errors.New("error"),
	}
)

func TestLoginService_Login(t *testing.T) {
	cases := []struct {
		name      string
		dec       mockJsonDecoder
		validator mockRequestValidator
		comparer  mockPasswordComparer
		repo      mockAuthFinder
		signer    Signer
		want      int
	}{
		{
			name:      "Positive",
			dec:       mockSuccessJsonDecoder,
			validator: mockSuccessRequestValidator,
			comparer:  mockSuccessPasswordComparer,
			repo:      mockSuccessAuthFinder,
			signer:    jwt.NewWithClaims,
			want:      200,
		},
		{
			name:      "Negative: Failed decode body",
			dec:       mockFailedJsonDecoder,
			validator: mockSuccessRequestValidator,
			comparer:  mockSuccessPasswordComparer,
			repo:      mockSuccessAuthFinder,
			signer:    jwt.NewWithClaims,
			want:      400,
		},
		{
			name:      "Negative: Failed to validate request",
			dec:       mockSuccessJsonDecoder,
			validator: mockFailedRequestValidator,
			comparer:  mockSuccessPasswordComparer,
			repo:      mockSuccessAuthFinder,
			signer:    jwt.NewWithClaims,
			want:      400,
		},
		{
			name:      "Negative: Failed find auth",
			dec:       mockSuccessJsonDecoder,
			validator: mockSuccessRequestValidator,
			comparer:  mockSuccessPasswordComparer,
			repo:      mockFailedAuthFinder,
			signer:    jwt.NewWithClaims,
			want:      404,
		},
		{
			name:      "Negative: Failed to compare password",
			dec:       mockSuccessJsonDecoder,
			validator: mockSuccessRequestValidator,
			comparer:  mockFailedPasswordComparer,
			repo:      mockSuccessAuthFinder,
			signer:    jwt.NewWithClaims,
			want:      401,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewLoginServices(tt.dec, tt.validator, tt.comparer, tt.signer, tt.repo)
			code, _ := svc.Login()
			if code != tt.want {
				t.Errorf("Expected resp to be %q but it was %q", tt.want, code)
			}
		})
	}
}
