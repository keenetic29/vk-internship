package password

import (
	"testing"
)

func TestHashAndCheck(t *testing.T) {
	t.Run("Valid password", func(t *testing.T) {
		password := "securePassword123"
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword failed: %v", err)
		}

		if err := CheckPassword(password, hash); err != nil {
			t.Error("Correct password should pass")
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		hash, _ := HashPassword("rightPassword")
		if err := CheckPassword("wrongPassword", hash); err == nil {
			t.Error("Wrong password should fail")
		}
	})

}