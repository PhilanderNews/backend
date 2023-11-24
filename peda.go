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
	var response CredentialUser
	response.Status = false

	mconn := SetConnection(mongoenv, dbname)

	var userdata User
	goblok := r.Header.Get("token")

	if goblok == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

		userdata.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			response.Message = "berhasil decode token"
			datauser := FindUser(mconn, collname, userdata)
			response.Status = true
			response.Data.Username = datauser.Username
			response.Data.Name = datauser.Name
			response.Data.Email = datauser.Email
			response.Data.Role = datauser.Role
		}
	}
	return ReturnStruct(response)
}

func Registrasi(mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialUser
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if usernameExists(mongoenv, dbname, datauser) {
		response.Message = "Username telah dipakai"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			hash, hashErr := HashPassword(datauser.Password)
			if hashErr != nil {
				response.Message = "Gagal Hash Password" + err.Error()
			}
			InsertUserdata(mconn, collname, datauser.Name, datauser.Email, datauser.Username, hash, datauser.Role.Admin, datauser.Role.Author)
			response.Status = true
			response.Message = "Berhasil Input data"
		}
	}
	return ReturnStruct(response)
}

func Login(privatekey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialUser
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

func HapusUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialUser
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var datauser User

	goblok := r.Header.Get("token")

	if goblok == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

		auth.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			auth2 := FindUser(mconn, collname, auth)
			if auth2.Role.Admin != true {
				response.Message = "anda bukan admin jadi tidak diizinkan"
			} else {
				err := json.NewDecoder(r.Body).Decode(&datauser)
				if err != nil {
					response.Message = "error parsing application/json: " + err.Error()
				} else {
					if datauser.Username == "" {
						response.Message = "parameter dari function ini adalah username"
					} else {
						DeleteUser(mconn, collname, datauser)
						response.Status = true
						response.Message = "berhasil hapus " + datauser.Username + " dari database"
					}
				}
			}
		}
	}
	return ReturnStruct(response)
}

func UpdateUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialUser
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var datauser User

	goblok := r.Header.Get("token")

	if goblok == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

		auth.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			auth2 := FindUser(mconn, collname, auth)
			if auth2.Role.Admin != true {
				response.Message = "anda bukan admin jadi tidak diizinkan"
			} else {
				err := json.NewDecoder(r.Body).Decode(&datauser)
				if err != nil {
					response.Message = "error parsing application/json: " + err.Error()
				} else {
					if datauser.Username == "" {
						response.Message = "parameter dari function ini adalah username"
					} else {
						hash, hashErr := HashPassword(datauser.Password)
						if hashErr != nil {
							response.Message = "Gagal Hash Password" + err.Error()
						}
						EditUser(mconn, collname, datauser.Name, datauser.Email, datauser.Username, hash, datauser.Role.Admin, datauser.Role.Author, datauser.Role.User, datauser)
						response.Status = true
						response.Message = "berhasil update " + datauser.Username + " dari database"
					}
				}
			}
		}
	}
	return ReturnStruct(response)
}

// ---------------------------------------------------------------------- Berita

func TambahBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialBerita
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var databerita Berita
	err := json.NewDecoder(r.Body).Decode(&databerita)

	var auth User
	goblok := r.Header.Get("token")

	if goblok == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

		auth.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			auth2 := FindUser(mconn, collname, auth)
			if auth2.Role.Author != true {
				response.Message = "anda bukan author ataupun admin jadi tidak diizinkan"
			} else {
				if err != nil {
					response.Message = "error parsing application/json: " + err.Error()
				} else {
					if idBeritaExists(mongoenv, dbname, databerita) {
						response.Message = "ID telah ada"
					} else {
						response.Status = true
						if err != nil {
							response.Message = "error parsing application/json: " + err.Error()
						} else {
							response.Status = true
							InsertBerita(mconn, collname, databerita)
							response.Message = "Berhasil Input data"
						}
					}
				}
			}
		}
	}
	return ReturnStruct(response)
}

func AmbilDataBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialBerita
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)

	var auth User
	goblok := r.Header.Get("token")

	if goblok == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

		auth.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			auth2 := FindUser(mconn, collname, auth)
			if auth2.Role.User != true {
				response.Message = "akun anda tidak aktif"
			} else {
				databerita := GetAllBerita(mconn, collname)
				return ReturnStruct(databerita)
			}
		}
	}
	return ReturnStruct(response)
}

func AmbilSatuBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialBerita
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)

	var databerita Berita

	var auth User
	goblok := r.Header.Get("token")

	if goblok == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

		auth.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			auth2 := FindUser(mconn, collname, auth)
			if auth2.Role.User != true {
				response.Message = "akun anda tidak aktif"
			} else {
				// idberita := r.URL.Query().Get("page")
				// databerita.ID = idberita
				berita := FindBerita(mconn, collname, databerita)

				return ReturnStruct(berita)
			}
		}
	}
	return ReturnStruct(response)
}

func HapusBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialBerita
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var databerita Berita

	goblok := r.Header.Get("token")

	if goblok == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

		auth.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			auth2 := FindUser(mconn, collname, auth)
			if auth2.Role.Author != true {
				response.Message = "anda bukan author ataupun admin jadi tidak diizinkan"
			} else {
				err := json.NewDecoder(r.Body).Decode(&databerita)
				if err != nil {
					response.Message = "error parsing application/json: " + err.Error()
				} else {
					if databerita.ID == "" {
						response.Message = "parameter dari function ini adalah id"
					} else {
						DeleteBerita(mconn, collname, databerita)
						response.Status = true
						response.Message = "berhasil hapus " + databerita.ID + " dari database"
					}
				}
			}
		}
	}
	return ReturnStruct(response)
}

func UpdateBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialBerita
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var databerita Berita

	goblok := r.Header.Get("token")

	if goblok == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(publickey), goblok)

		auth.Username = checktoken //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			auth2 := FindUser(mconn, collname, auth)
			if auth2.Role.Admin != true {
				response.Message = "anda bukan admin jadi tidak diizinkan"
			} else {
				err := json.NewDecoder(r.Body).Decode(&databerita)
				if err != nil {
					response.Message = "error parsing application/json: " + err.Error()
				} else {
					if databerita.ID == "" {
						response.Message = "parameter dari function ini adalah id"
					} else {
						EditBerita(mconn, collname, databerita)
						response.Status = true
						response.Message = "berhasil update " + databerita.ID + " dari database"
					}
				}
			}
		}
	}
	return ReturnStruct(response)
}
