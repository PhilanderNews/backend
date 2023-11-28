package peda

import (
	"fmt"
	"testing"
)

var privatekey2 = "085151de76c61415e40ee4440b3848d0317a28e58fc7f9ddff8455cca31fe499c78ba61336ef2a8ac88e3126e0c3ad0e4cdccc271a11240db3f61e91bb68a936"
var publickey2 = "c78ba61336ef2a8ac88e3126e0c3ad0e4cdccc271a11240db3f61e91bb68a936"
var encode = "v4.public.eyJleHAiOiIyMDIzLTExLTI5VDAwOjIyOjA1KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yOFQyMjoyMjowNSswNzowMCIsIm5iZiI6IjIwMjMtMTEtMjhUMjI6MjI6MDUrMDc6MDAiLCJyb2xlIjoiYWRtaW4iLCJ1c2VybmFtZSI6Imlicm9oaW0ifXUeIuOthFir7tNFpkCsdb41IKBSDpYoAZlRhD3wgryojOf-e13zZQga_uxvrvv7aXdnX6PJWXWr4NJE_RXLyAE"

func TestGeneratePaseto(t *testing.T) {
	privateKey, publicKey := GenerateKey()
	fmt.Println("privateKey: " + privateKey)
	fmt.Println("publicKey: " + publicKey)
}

func TestEncode(t *testing.T) {
	username := "ibrohim"
	role := "admin"

	tokenstring, err := Encode(username, role, privatekey2)
	fmt.Println("error : ", err)
	fmt.Println("token : ", tokenstring)
}

func TestDecode(t *testing.T) {
	pay, err := Decode(publickey2, encode)
	user := DecodeGetUsername(publickey2, encode)
	role := DecodeGetRole(publickey2, encode)

	fmt.Println("user :", user)
	fmt.Println("role :", role)
	fmt.Println("err : ", err)
	fmt.Println("payload : ", pay)
}
