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
		// {
		// 	name:   "Positive",
		// 	header: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDEyNTEzMTcsImlkIjoxLCJ1dWlkIjoiNTA0ZDYxYTUtYjI5OC00NWY0LTllOWYtOGNjYTM0ZjQ3OTEyIn0.B2_A7SJvmKcILcwWcPLLe1Y1BlRCBsaz8aEZ2T14lUA",
		// 	want:   200,
		// },
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
