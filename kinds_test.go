package errx

import (
	"fmt"
	"strconv"
	"testing"
)

var (
	NotFound  = Kind("notfound")
	FileOpen  = DataKind[string]("fileopen")
	PageLoad  = DataKind[string]("pageload")
	UserLogin = DataKind[UserInfo]("userlogin")
)

type UserInfo map[string]string

func (u UserInfo) Id() int {
	v, err := strconv.Atoi(u["id"])
	if err != nil {
		panic(err)
	}

	return v
}

func (u UserInfo) Email() string {
	return u["email"]
}

func TestKinds(t *testing.T) {
	err := NewKind(1746358247, UserLogin(UserInfo{"id": "1234400114", "email": "master@mailer4.com"}), "something wrong")
	err = WrapKind(1746358802, PageLoad("test.com"), err)
	err = Wrap(1746358817, err)
	err = WrapKind(1746358827, NotFound, err)
	err = WrapKind(1746358853, FileOpen("users.txt"), err)

	fmt.Println(err)

	if d, ok := FindData(err, PageLoad); ok {
		fmt.Println("data:", *d)
	}

	fmt.Println("done")
}
