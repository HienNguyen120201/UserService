package util

import "golang.org/x/crypto/bcrypt"

func HashPwd(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(bytes), err
}

func CheckPwdHash(hashedPwd, pwd string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPwd),
		[]byte(pwd),
	)
}
