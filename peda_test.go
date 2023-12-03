package peda

import (
	"fmt"
	"testing"
	"time"
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

// ---------------------------------------------------------------------- Tutorial

func TestInsertMongo(t *testing.T) {
	mconn := SetConnectionTest(mongoenv, dbname)
	testPesan := Tutorial{Parameter: "2", Pesan: "Apa yah"}
	testinsert := InsertMongo(mconn, "test", testPesan)
	fmt.Println(testinsert)
}

func TestGetAllMongo(t *testing.T) {
	mconn := SetConnectionTest(mongoenv, dbname)
	testinsert := GetAllMongo(mconn, "test")
	fmt.Println(testinsert)
}

func TestGetOneMongo(t *testing.T) {
	mconn := SetConnectionTest(mongoenv, dbname)
	testPesan := Tutorial{Parameter: "2"}
	testinsert := GetOneMongo(mconn, "test", testPesan)
	fmt.Println(testinsert)
}

func TestUpdateMongo(t *testing.T) {
	mconn := SetConnectionTest(mongoenv, dbname)
	testPesan := Tutorial{Parameter: "1", Pesan: "Berhasil Update"}
	testinsert := UpdateMongo(mconn, "test", testPesan)
	fmt.Println(testinsert)
}

func TestDeleteMongo(t *testing.T) {
	mconn := SetConnectionTest(mongoenv, dbname)
	testPesan := Tutorial{Parameter: "1"}
	testinsert := DeleteMongo(mconn, "test", testPesan)
	fmt.Println(testinsert)
}

func TestTanggal(t *testing.T) {
	wib, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		fmt.Println("Error parsing time location: " + err.Error())
	}

	currentTime := time.Now().In(wib)
	timeStringKomentar := currentTime.Format("January 2, 2006")
	timeStringBerita := currentTime.Format("Monday, 2 January 2006 15:04 MST")

	fmt.Println(timeStringKomentar)
	fmt.Println(timeStringBerita)
}
