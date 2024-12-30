package hash

import "golang.org/x/crypto/bcrypt"

// GenerateHash 비밀번호 해싱
func GenerateHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CompareHash 해시된 비밀번호 비교
func CompareHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
