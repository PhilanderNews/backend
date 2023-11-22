package peda

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpdateGetData(t *testing.T) {
	mconn := SetConnection("mongoenv", "unistore")
	datagedung := GetAllUser(mconn, "user")
	fmt.Println(datagedung)
}

// 	result := GCFCreateHandler(MONGOCONNSTRINGENV, dbname, collectionname, datauser)
// 	fmt.Println(result)
// 	// You can add assertions here to validate the result, or check the database for the created user.
// }

func TestCreateNewUserRole(t *testing.T) {
	var userdata User
	userdata.Username = "unistore"
	userdata.Password = "unistore"
	mconn := SetConnection("MONGOCONNSTRINGENV", "unistore")
	CreateNewUserRole(mconn, "user", userdata)
}

func TestDeleteUser(t *testing.T) {

	mconn := SetConnection("mongoenv", "unistore")
	var userdata User
	userdata.Username = "unistore"
	DeleteUser(mconn, "user", userdata)
}

func TestGFCPostHandlerUser(t *testing.T) {
	mconn := SetConnection("mongoenv", "unistore")
	var userdata User
	userdata.Username = "unistore"
	userdata.Password = "unistore"
	CreateNewUserRole(mconn, "user", userdata)
}

func TestFunciionUser(t *testing.T) {
	mconn := SetConnection("mongoenv", "unistore")
	var userdata User
	userdata.Username = "unistore"
	userdata.Password = "unistore"
	CreateNewUserRole(mconn, "user", userdata)
}

func TestGeneratePasswordHashh(t *testing.T) {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestHashFunctionn(t *testing.T) {
	mconn := SetConnection("mongoenv", "unistore")
	var userdata User
	userdata.Username = "unistore"
	userdata.Password = "unistore"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}
func TestFindUser(t *testing.T) {
	var userdata User
	userdata.Username = "unistore"
	mconn := SetConnection("mongoenv", "unistore")
	res := FindUser(mconn, "user", userdata)
	fmt.Println(res)
}

func TestGeneratePasswordHash(t *testing.T) {
	password := "unistore"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)
	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("unistore", privateKey)
	fmt.Println(hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("mongoenv", "unistore")
	var userdata User
	userdata.Username = "unistore"
	userdata.Password = "unistore"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("mongoenv", "unistore")
	var userdata User
	userdata.Username = "unistore"
	userdata.Password = "unistore"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}

func TestRole(t *testing.T) {
	role := "admin"
	if role == "admin" {
		fmt.Println("anu")
	} else {
		fmt.Println("2sad")
	}

}
