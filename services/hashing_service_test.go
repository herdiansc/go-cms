package services

import "testing"

func mockSuccessHashGenerator(password []byte, cost int) ([]byte, error) {
	return []byte{}, nil
}

func TestHashingService_HashPassword(t *testing.T) {
	cases := []struct {
		name string
		hash HashGenerator
		want error
	}{
		{
			name: "Positive",
			hash: mockSuccessHashGenerator,
			want: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewHashingService(tt.hash)
			_, err := svc.HashPassword("a")
			if err != tt.want {
				t.Errorf("Expected err to be %q but it was %q", tt.want, err)
			}
		})
	}
}

func mockSuccessHashComparer(hashedPassword, password []byte) error {
	return nil
}

func TestHashingCompareService_VerifyPassword(t *testing.T) {
	cases := []struct {
		name    string
		compare HashComparer
		want    bool
	}{
		{
			name:    "Positive",
			compare: mockSuccessHashComparer,
			want:    true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewHashingCompareService(tt.compare)
			result := svc.VerifyPassword("a", "hash-a")
			if result != tt.want {
				t.Errorf("Expected result to be %v but it was %v", tt.want, result)
			}
		})
	}
}
