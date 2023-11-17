package peda

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
)

func ReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func Registrasi(mongoenv, dbname, collname string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if usernameExists(mongoenv, dbname, datauser) {
		response.Status = false
		response.Message = "Username telah dipakai"
	} else {
		response.Status = true
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			response.Status = true
			hash, hashErr := HashPassword(datauser.Password)
			if hashErr != nil {
				response.Message = "Gagal Hash Password" + err.Error()
			}
			InsertUserdata(mconn, collname, datauser.Name, datauser.Email, datauser.Username, datauser.Role, hash)
			response.Message = "Berhasil Input data"
		}
	}
	return ReturnStruct(response)
}

func Login(privatekey, mongoenv, dbname, collname string, r *http.Request) string {
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return ReturnStruct(CreateToken("error parsing application/json: ", err.Error()))
	} else {
		if IsPasswordValid(mconn, collname, datauser) {
			user := FindUser(mconn, collname, datauser)
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(privatekey))
			if err != nil {
				return ReturnStruct(CreateToken("Gagal Encode Token :", err.Error()))
			} else {
				return ReturnStruct(CreateToken(tokenstring, user))
			}
		} else {
			return ReturnStruct(CreateToken("Password Salah", datauser))
		}
	}
}

func HapusUser(mongoenv, dbname, collname string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
	} else {
		DeleteUser(mconn, collname, datauser)
		response.Message = "Berhasil Delete data"
	}
	return ReturnStruct(response)
}
