package services

import (
	"testing"
)

func TestProfileServices_GetProfile(t *testing.T) {
	type fields struct {
		authData any
		repo     mockAuthFinder
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Positive",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockSuccessAuthFinder,
			},
			want: 200,
		},
		{
			name: "Failed to read authData",
			fields: fields{
				authData: "invalid",
				repo:     mockSuccessAuthFinder,
			},
			want: 400,
		},
		{
			name: "Failed to get data",
			fields: fields{
				authData: mockValidAuthData,
				repo:     mockFailedAuthFinder,
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProfileServices(tt.fields.authData, tt.fields.repo)
			got, _ := svc.GetProfile()
			if got != tt.want {
				t.Errorf("ProfileServices.GetProfile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
