package user

import (
	"fmt"
	"testing"
)

func TestAuth(t *testing.T) {

	pwd, err := GenerateHashPwd("Foobar")
	if err != nil {
		fmt.Println("加密错误", err)
		return
	}
	fmt.Println("pwd", pwd)
}

func TestCompare(t *testing.T) {

	pwd, _ := GenerateHashPwd("123456")

	fmt.Println("pwd", pwd)

	err := Compare(pwd, "123456")
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println("pass 123456")

	err = Compare(pwd, "1234567")
	if err != nil {
		fmt.Println("err2", err)
		return
	}

}
