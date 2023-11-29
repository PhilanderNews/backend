package peda

import (
	"fmt"
	"testing"
)

var privatekey2 = ""
var publickey2 = ""
var encode = ""

func TestGeneratePaseto(t *testing.T) {
	privateKey, publicKey := GenerateKey()
	fmt.Println("privateKey: " + privateKey)
	fmt.Println("publicKey: " + publicKey)
}

func TestEncode(t *testing.T) {
	name := "Ibrohim Mubarok"
	username := "ibrohim"
	role := "admin"

	tokenstring, err := Encode(name, username, role, privatekey2)
	fmt.Println("error : ", err)
	fmt.Println("token : ", tokenstring)
}

func TestDecode(t *testing.T) {
	pay, err := Decode(publickey2, encode)
	name := DecodeGetName(publickey2, encode)
	username := DecodeGetUsername(publickey2, encode)
	role := DecodeGetRole(publickey2, encode)

	fmt.Println("name :", name)
	fmt.Println("username :", username)
	fmt.Println("role :", role)
	fmt.Println("err : ", err)
	fmt.Println("payload : ", pay)
}
