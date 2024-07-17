package crypto

import (
	"strings"

	"github.com/matthxwpavin/ticketing/random"
	"golang.org/x/crypto/bcrypt"
)

func GenerateWithSalt(password string) (SaltedPassword, error) {
	salt := random.RandStringMaskImprSrcUnsafe(8)
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return SaltedPassword(strings.Join([]string{string(hashedPasswd), salt}, ".")), nil
}

func Compare(sltPasswd string, passwd string) error {
	hashedPasswd, salt := SaltedPassword(sltPasswd).Split()
	return bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(passwd+salt))
}

type SaltedPassword string

func (s SaltedPassword) Split() (passwd string, salt string) {
	sp := string(s)
	i := strings.LastIndex(sp, ".")
	return sp[:i], sp[i+1:]
}

func (s SaltedPassword) String() string {
	return string(s)
}
