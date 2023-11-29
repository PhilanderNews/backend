package peda

import (
	"fmt"
	"testing"
)

var privatekey2 = "085151de76c61415e40ee4440b3848d0317a28e58fc7f9ddff8455cca31fe499c78ba61336ef2a8ac88e3126e0c3ad0e4cdccc271a11240db3f61e91bb68a936"
var publickey2 = "c78ba61336ef2a8ac88e3126e0c3ad0e4cdccc271a11240db3f61e91bb68a936"
var encode = "v4.public.eyJleHAiOiIyMDIzLTExLTMwVDAxOjE4OjIyKzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yOVQyMzoxODoyMiswNzowMCIsIm5hbWUiOiJJYnJvaGltIE11YmFyb2siLCJuYmYiOiIyMDIzLTExLTI5VDIzOjE4OjIyKzA3OjAwIiwicm9sZSI6ImFkbWluIiwidXNlcm5hbWUiOiJpYnJvaGltIn06oaydO1Z0ok4O8y2JSX918IcdHc_ThArPS8wxX86ay8BixM2Ugzn7DQ5t03yoeVH0vV1JC-qAwIdvgfpTWncP"

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
