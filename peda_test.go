package peda

import (
	"fmt"
	"testing"
	"time"
)

var privatekey = "04c27ac63911f885e270848d41934e5848f6efa0c89665d5135aa4b11208e6a73e806beacb6c90b2005c1d24c7cb98d40964c245f5a082e36a85dbc1e1668503"
var publickey = "3e806beacb6c90b2005c1d24c7cb98d40964c245f5a082e36a85dbc1e1668503"
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
	name := "Test Nama"
	username := "Test Username"
	role := "Test Role"

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
	testPesan := Tutorial{Parameter: "2", Pesan: "testing 2"}
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
	testPesan := Tutorial{Parameter: "1"}
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
	testPesan := Tutorial{Parameter: "2"}
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
