package data

import (
	"fmt"
)

const (
	NEW_USER = 1
)

var (
	USER_NOT_EXIST   = fmt.Errorf("user not exist")
	PASSWD_INCORRECT = fmt.Errorf("passwd incorrect")
)
