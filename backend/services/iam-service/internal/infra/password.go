package infra

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher bcrypt 密码哈希器
type BcryptPasswordHasher struct {
	cost int
}

// NewBcryptPasswordHasher 创建 bcrypt 密码哈希器
func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptPasswordHasher{cost: cost}
}

func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *BcryptPasswordHasher) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
