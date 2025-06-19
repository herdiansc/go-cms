package services

// HashGenerator defines hash generator func
type HashGenerator func(password []byte, cost int) ([]byte, error)

// HashingService defines hash service struct
type HashingService struct {
	hash HashGenerator
}

// NewHashingService inits HashingService
func NewHashingService(h HashGenerator) HashingService {
	return HashingService{
		hash: h,
	}
}

// HashPassword generates a bcrypt hash for the given password.
func (svc HashingService) HashPassword(password string) (string, error) {
	bytes, err := svc.hash([]byte(password), 14)
	return string(bytes), err
}

// HashComparer defines hash comparer func
type HashComparer func(hashedPassword, password []byte) error

// HashingCompareService defines hash comparer service
type HashingCompareService struct {
	compare HashComparer
}

// NewHashingCompareService inits HashingCompareService
func NewHashingCompareService(hc HashComparer) HashingCompareService {
	return HashingCompareService{
		compare: hc,
	}
}

// VerifyPassword verifies if the given password matches the stored hash.
func (svc HashingCompareService) VerifyPassword(password, hash string) bool {
	err := svc.compare([]byte(hash), []byte(password))
	return err == nil
}
