package passwd

import (
	"strings"

	"github.com/matthxwpavin/ticketing/random"
	"golang.org/x/crypto/bcrypt"
)

func Generate(password string) (string, error) {
	salt := random.RandStringMaskImprSrcUnsafe(8)
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{string(hashedPasswd), salt}, "."), nil
}

func Compare(sltPasswd string, passwd string) error {
	i := strings.LastIndex(sltPasswd, ".")
	hashed, salt := sltPasswd[:i], sltPasswd[i+1:]
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(passwd+salt))
}
