package peda

import (
	"fmt"
	"testing"
)

var privatekey = "34853f421ac024ef0e9ea957600a4eefb725af2e5ddda3cdd0c2d688ba35d245f963664dfce8aea29a5c9178b4ca3b3aac47d02f2d4093b8ea9675f001f7bb98"
var publickey = "f963664dfce8aea29a5c9178b4ca3b3aac47d02f2d4093b8ea9675f001f7bb98"
var encode = ""
var mongoenv = ""
var dbname = ""
var collname = ""

func TestGeneratePaseto(t *testing.T) {
	privateKey, publicKey := GenerateKey()
	fmt.Println("privatekey: " + privateKey)
	fmt.Println("publickey: " + publicKey)
}

func TestEncode(t *testing.T) {
	name := "Ibrohim Mubarok"
	username := "asal"
	role := "admin"

	tokenstring, err := Encode(name, username, role, privatekey)
	fmt.Println("error : ", err)
	fmt.Println("token : ", tokenstring)
}

func TestDecode(t *testing.T) {
	pay, err := Decode(publickey, encode)
	name := DecodeGetName(publickey, encode)
	username := DecodeGetUsername(publickey, encode)
	role := DecodeGetRole(publickey, encode)

	fmt.Println("name :", name)
	fmt.Println("username :", username)
	fmt.Println("role :", role)
	fmt.Println("err : ", err)
	fmt.Println("payload : ", pay)
}

func TestUsernameExists(t *testing.T) {
	userdata := User{
		Username: "ibrohim",
		// Lengkapi dengan data pengguna lainnya jika diperlukan
	}
	hasil := !usernameExists(mongoenv, dbname, userdata)
	fmt.Println(hasil)
}
