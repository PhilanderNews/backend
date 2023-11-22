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

func Authorization(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response AuthorizationStruct
	response.Status = false

	mconn := SetConnection(mongoenv, dbname)

	var userdata User
	goblok := r.Header.Get("token")

	checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

	userdata.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

	if checktoken == "" {
		response.Message = "hasil decode tidak ada"
	} else {
		response.Message = "hasil decode token"
		datauser := FindUser(mconn, collname, userdata)
		response.Status = true
		response.Data.Username = datauser.Username
		response.Data.Name = datauser.Name
		response.Data.Email = datauser.Email
		response.Data.Role = datauser.Role
	}

	return ReturnStruct(response)
}

func Registrasi(mongoenv, dbname, collname string, r *http.Request) string {
	var response AuthorizationStruct
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
			InsertUserdata(mconn, collname, datauser.Name, datauser.Email, datauser.Username, hash, datauser.Role.Admin, datauser.Role.Author)
			response.Message = "Berhasil Input data"
		}
	}
	return ReturnStruct(response)
}

func Login(privatekey, mongoenv, dbname, collname string, r *http.Request) string {
	var response AuthorizationStruct
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collname, datauser) {
			user := FindUser(mconn, collname, datauser)
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(privatekey))
			if err != nil {
				return ReturnStruct(response.Message == "Gagal Encode Token :"+err.Error())
			} else {
				response.Status = true
				response.Data.Name = user.Name
				response.Data.Email = user.Email
				response.Data.Username = user.Username
				response.Data.Role = user.Role
				response.Message = "User berhasil login"
				response.Token = tokenstring
				return ReturnStruct(response)
			}
		} else {
			response.Message = "Password Salah"
		}
	}
	return ReturnStruct(response)
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

// ---------------------------------------------------------------------- Berita

func TambahBerita(mongoenv, dbname, collname string, r *http.Request) string {
	var response Credential
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var databerita Berita
	err := json.NewDecoder(r.Body).Decode(&databerita)
	if idBeritaExists(mongoenv, dbname, databerita) {
		response.Status = false
		response.Message = "ID telah ada"
	} else {
		response.Status = true
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			response.Status = true
			InsertBerita(mconn, collname, databerita.ID, databerita.Kategori, databerita.Judul, databerita.Preview, databerita.Konten)
			response.Message = "Berhasil Input data"
		}
	}
	return ReturnStruct(response)
}

func AmbilDataBerita(mongoenv, dbname, collname string) string {
	mconn := SetConnection(mongoenv, dbname)
	databerita := GetAllBerita(mconn, collname)
	return ReturnStruct(databerita)
}

func AmbilSatuBerita(mongoenv, dbname, collname string, r *http.Request) string {
	mconn := SetConnection(mongoenv, dbname)
	var databerita Berita
	berita := FindBerita(mconn, collname, databerita)

	return ReturnStruct(berita)
}
