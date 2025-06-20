package services

import (
	"testing"
)

func TestTokenVerifyServices_Verify(t *testing.T) {
	cases := []struct {
		name   string
		header string
		want   int
	}{
		{
			name:   "Positive",
			header: "Bearer invalid",
			want:   400,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTokenVerifyServices()
			code, _ := svc.Verify(tt.header)
			if code != tt.want {
				t.Errorf("Expected resp to be %q but it was %q", tt.want, code)
			}
		})
	}
}
