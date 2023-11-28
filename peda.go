package peda

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/whatsauth/wa"
)

func ReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func Authorization(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialUser
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		checktoken := DecodeGetUsername(os.Getenv(publickey), header)
		auth.Username = checktoken
		if checktoken == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				response.Message = "berhasil decode token"
				datauser := FindUser(mconn, collname, auth)
				response.Status = true
				response.Data.Username = datauser.Username
				response.Data.Name = datauser.Name
				response.Data.Email = datauser.Email
				response.Data.Role = datauser.Role
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

func Registrasi(token, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if usernameExists(mongoenv, dbname, datauser) {
		response.Message = "username telah dipakai"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			hash, hashErr := HashPassword(datauser.Password)
			if hashErr != nil {
				response.Message = "gagal Hash Password" + err.Error()
			}
			if datauser.No_whatsapp == "" {
				response.Message = "nomor whatsapp wajib diisi"
			}
			InsertUserdata(mconn, collname, datauser.Name, datauser.Email, datauser.No_whatsapp, datauser.Username, hash, datauser.Role)
			response.Status = true
			response.Message = "berhasil Input data"

			var username = datauser.Username
			var password = datauser.Password
			var nohp = datauser.No_whatsapp

			dt := &wa.TextMessage{
				To:       nohp,
				IsGroup:  false,
				Messages: "Selamat anda berhasil registrasi, berikut adalah username anda: " + username + " dan ini adalah password anda: " + password + "\nDisimpan baik baik ya",
			}

			atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")
		}
	}
	return ReturnStruct(response)
}

func Login(token, privatekey, mongoenv, dbname, collname string, r *http.Request) string {
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
			var nama = user.Name
			var nohp = user.No_whatsapp
			tokenstring, err := Encode(user.Username, user.Role, os.Getenv(privatekey))
			if err != nil {
				return ReturnStruct(response.Message == "Gagal Encode Token :"+err.Error())
			} else {
				response.Status = true
				response.Data.Name = nama
				response.Data.Email = user.Email
				response.Data.Username = user.Username
				response.Data.Role = user.Role
				response.Token = tokenstring
				response.Message = "berhasil login"

				dt := &wa.TextMessage{
					To:       nohp,
					IsGroup:  false,
					Messages: nama + " berhasil login\nPerlu diingat sesi login hanya berlaku 2 jam",
				}

				atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")

				return ReturnStruct(response)
			}
		} else {
			response.Message = "Password Salah"
		}
	}
	return ReturnStruct(response)
}

func HapusUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var datauser User

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" {
					err := json.NewDecoder(r.Body).Decode(&datauser)
					if err != nil {
						response.Message = "error parsing application/json: " + err.Error()
					} else {
						if datauser.Username == "" {
							response.Message = "parameter dari function ini adalah username"
						} else {
							if usernameExists(mongoenv, dbname, datauser) {
								DeleteUser(mconn, collname, datauser)
								response.Status = true
								response.Message = "berhasil hapus " + datauser.Username + " dari database"
							} else {
								response.Message = "akun yang ingin dihapus tidak ditemukan"
							}
						}
					}
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

func UpdateUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var datauser User

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername //userdata.Username dibuat menjadi checktoken agar userdata.Username dapat digunakan sebagai filter untuk menggunakan function FindUser

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" {
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
							if usernameExists(mongoenv, dbname, datauser) {
								EditUser(mconn, collname, datauser.Name, datauser.Email, datauser.No_whatsapp, datauser.Username, hash, datauser.Role)
								response.Status = true
								response.Message = "berhasil update " + datauser.Username + " dari database"
							} else {
								response.Message = "akun yang ingin diedit tidak ditemukan"
							}
						}
					}
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

func AmbilSemuaUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var datauser User

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" {
					err := json.NewDecoder(r.Body).Decode(&datauser)
					if err != nil {
						response.Message = "error parsing application/json: " + err.Error()
					} else {
						datauser := GetAllUser(mconn, collname)
						return ReturnStruct(datauser)
					}
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

func AmbilSatuUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var datauser User

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" {
					err := json.NewDecoder(r.Body).Decode(&datauser)
					if err != nil {
						response.Message = "error parsing application/json: " + err.Error()
					} else {
						if usernameExists(mongoenv, dbname, datauser) {
							user := FindUser(mconn, collname, datauser)
							return ReturnStruct(user)
						} else {
							response.Message = "berita tidak ditemukan"
						}
					}
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

// ---------------------------------------------------------------------- Berita

func TambahBerita(publickey, mongoenv, dbname, colluser, collberita string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var databerita Berita
	err := json.NewDecoder(r.Body).Decode(&databerita)

	var auth User
	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" || tokenrole == "author" {
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
								InsertBerita(mconn, collberita, databerita)
								response.Message = "berhasil Input data"
							}
						}
					}
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

func AmbilSemuaBerita(publickey, mongoenv, dbname, colluser, collberita string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)

	var auth User
	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" || tokenrole == "author" || tokenrole == "user" {
					databerita := GetAllBerita(mconn, collberita)
					return ReturnStruct(databerita)
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

func AmbilSatuBerita(publickey, mongoenv, dbname, colluser, collberita string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)

	var databerita Berita

	var auth User
	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" || tokenrole == "author" || tokenrole == "user" {
					idberita := r.URL.Query().Get("page")
					databerita.ID = idberita
					if idBeritaExists(mongoenv, dbname, databerita) {
						berita := FindBerita(mconn, collberita, databerita)
						return ReturnStruct(berita)
					} else {
						response.Message = "berita tidak ditemukan"
					}
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

func HapusBerita(publickey, mongoenv, dbname, colluser, collberita string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var databerita Berita

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" {
					err := json.NewDecoder(r.Body).Decode(&databerita)
					if err != nil {
						response.Message = "error parsing application/json: " + err.Error()
					} else {
						if databerita.ID == "" {
							response.Message = "parameter dari function ini adalah id"
						} else {
							if idBeritaExists(mongoenv, dbname, databerita) {
								DeleteBerita(mconn, collberita, databerita)
								response.Status = true
								response.Message = "berhasil hapus " + databerita.ID + " dari database"
							} else {
								response.Message = "berita tidak ditemukan"
							}
						}
					}
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}

func UpdateBerita(publickey, mongoenv, dbname, colluser, collberita string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var databerita Berita

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)

		auth.Username = tokenusername

		if tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				if tokenrole == "admin" {
					err := json.NewDecoder(r.Body).Decode(&databerita)
					if err != nil {
						response.Message = "error parsing application/json: " + err.Error()
					} else {
						if databerita.ID == "" {
							response.Message = "parameter dari function ini adalah id"
						} else {
							if idBeritaExists(mongoenv, dbname, databerita) {
								EditBerita(mconn, collberita, databerita)
								response.Status = true
								response.Message = "berhasil update " + databerita.ID + " dari database"
							} else {
								response.Message = "berita tidak ditemukan"
							}
						}
					}
				} else {
					response.Message = "anda tidak memiliki akses"
				}
			} else {
				response.Message = "akun tidak ditemukan"
			}
		}
	}
	return ReturnStruct(response)
}
