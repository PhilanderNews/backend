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
	var auth User
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		tokenname := DecodeGetName(os.Getenv(publickey), header)
		tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
		tokenrole := DecodeGetRole(os.Getenv(publickey), header)
		auth.Username = tokenusername
		if tokenname == "" || tokenusername == "" || tokenrole == "" {
			response.Message = "hasil decode tidak ditemukan"
		} else {
			if usernameExists(mongoenv, dbname, auth) {
				response.Message = "berhasil decode token"
				response.Status = true
				response.Data.Username = tokenname
				response.Data.Name = tokenusername
				response.Data.Role = tokenrole
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
			} else {
				if datauser.No_whatsapp == "" {
					response.Message = "nomor whatsapp wajib diisi"
				} else {
					InsertUserdata(mconn, collname, datauser.Name, datauser.Email, datauser.No_whatsapp, datauser.Username, hash, datauser.Role)
					response.Status = true
					response.Message = "berhasil Input data"

					var username = datauser.Username
					var password = datauser.Password
					var nohp = datauser.No_whatsapp

					dt := &wa.TextMessage{
						To:       nohp,
						IsGroup:  false,
						Messages: "Selamat anda berhasil registrasi, berikut adalah username anda: " + username + "\nDan ini adalah password anda: " + password + "\nDisimpan baik baik ya",
					}

					atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")
				}
			}

		}
	}
	return ReturnStruct(response)
}

func Login(token, privatekey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
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
			tokenstring, tokenerr := Encode(user.Name, user.Username, user.Role, os.Getenv(privatekey))
			if tokenerr != nil {
				return ReturnStruct(response.Message == "gagal encode token :"+err.Error())
			} else {
				response.Status = true
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

	err := json.NewDecoder(r.Body).Decode(&datauser)

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if tokenrole == "admin" {
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
					} else {
						response.Message = "anda tidak memiliki akses"
					}
				} else {
					response.Message = "akun tidak ditemukan"
				}
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

	err := json.NewDecoder(r.Body).Decode(&datauser)

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if tokenrole == "admin" {
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
					} else {
						response.Message = "anda tidak memiliki akses"
					}
				} else {
					response.Message = "akun tidak ditemukan"
				}
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
					datauser := GetAllUser(mconn, collname)
					return ReturnStruct(datauser)
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

	err := json.NewDecoder(r.Body).Decode(&datauser)

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if tokenrole == "admin" {
						if usernameExists(mongoenv, dbname, datauser) {
							user := FindUser(mconn, collname, datauser)
							return ReturnStruct(user)
						} else {
							response.Message = "user tidak ditemukan"
						}
					} else {
						response.Message = "anda tidak memiliki akses"
					}
				} else {
					response.Message = "akun tidak ditemukan"
				}
			}
		}
	}
	return ReturnStruct(response)
}

// ---------------------------------------------------------------------- Berita

func TambahBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
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
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenname := DecodeGetName(os.Getenv(publickey), header)
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenname == "" || tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if tokenrole == "admin" || tokenrole == "author" {
						if idBeritaExists(mongoenv, dbname, databerita) {
							response.Message = "ID telah ada"
						} else {
							response.Status = true
							databerita.Penulis = tokenname
							InsertBerita(mconn, collname, databerita)
							response.Message = "berhasil Input data"
						}
					} else {
						response.Message = "anda tidak memiliki akses"
					}
				} else {
					response.Message = "akun tidak ditemukan"
				}
			}
		}

	}
	return ReturnStruct(response)
}

func AmbilSemuaBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
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
					databerita := GetAllBerita(mconn, collname)
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

func AmbilSatuBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
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
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if tokenrole == "admin" || tokenrole == "author" || tokenrole == "user" {
						if idBeritaExists(mongoenv, dbname, databerita) {
							berita := FindBerita(mconn, collname, databerita)
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
	}
	return ReturnStruct(response)
}

func HapusBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var databerita Berita

	err := json.NewDecoder(r.Body).Decode(&databerita)

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenname := DecodeGetName(os.Getenv(publickey), header)
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenname == "" || tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if databerita.ID == "" {
						response.Message = "parameter dari function ini adalah id"
					} else {
						namapenulis := FindBerita(mconn, collname, databerita)
						if tokenrole == "admin" || tokenname == namapenulis.Penulis {
							if idBeritaExists(mongoenv, dbname, databerita) {
								DeleteBerita(mconn, collname, databerita)
								response.Status = true
								response.Message = "berhasil hapus " + databerita.ID + " dari database"
							} else {
								response.Message = "berita tidak ditemukan"
							}
						} else {
							response.Message = "anda tidak memiliki akses"
						}
					}
				} else {
					response.Message = "akun tidak ditemukan"
				}
			}
		}
	}
	return ReturnStruct(response)
}

func UpdateBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var databerita Berita

	err := json.NewDecoder(r.Body).Decode(&databerita)

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenname := DecodeGetName(os.Getenv(publickey), header)
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenname == "" || tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if databerita.ID == "" {
						response.Message = "parameter dari function ini adalah id"
					} else {
						namapenulis := FindBerita(mconn, collname, databerita)
						if tokenrole == "admin" || tokenname == namapenulis.Penulis {
							if idBeritaExists(mongoenv, dbname, databerita) {
								databerita.Penulis = tokenname
								EditBerita(mconn, collname, databerita)
								response.Status = true
								response.Message = "berhasil update " + databerita.ID + " dari database"
							} else {
								response.Message = "berita tidak ditemukan"
							}
						} else {
							response.Message = "anda tidak memiliki akses"
						}
					}
				} else {
					response.Message = "akun tidak ditemukan"
				}
			}
		}
	}
	return ReturnStruct(response)
}

// ---------------------------------------------------------------------- Komentar

func TambahKomentar(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datakomentar Komentar
	var databerita Berita
	databerita.ID = datakomentar.ID_berita

	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	var auth User
	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenname := DecodeGetName(os.Getenv(publickey), header)
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenname == "" || tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if tokenrole == "admin" || tokenrole == "author" || tokenrole == "user" {
						if idBeritaExists(mongoenv, dbname, databerita) {
							if idKomentarExists(mongoenv, dbname, datakomentar) {
								response.Message = "ID telah ada"
							} else {
								response.Status = true
								datakomentar.Name = tokenname
								InsertKomentar(mconn, collname, datakomentar)
								response.Message = "berhasil Input data"
							}
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
	}
	return ReturnStruct(response)
}

func AmbilSemuaKomentar(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datakomentar Komentar
	var databerita Berita
	databerita.ID = datakomentar.ID_berita

	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	var auth User
	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if tokenrole == "admin" || tokenrole == "author" || tokenrole == "user" {
						if idBeritaExists(mongoenv, dbname, databerita) {
							datakomentar := GetAllKomentar(mconn, collname)
							return ReturnStruct(datakomentar)
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
	}
	return ReturnStruct(response)
}

func AmbilSatuKomentar(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datakomentar Komentar
	var databerita Berita
	databerita.ID = datakomentar.ID_berita
	var auth User

	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if tokenrole == "admin" || tokenrole == "author" || tokenrole == "user" {
						if idBeritaExists(mongoenv, dbname, databerita) {
							if idKomentarExists(mongoenv, dbname, datakomentar) {
								komentar := FindKomentar(mconn, collname, datakomentar)
								return ReturnStruct(komentar)
							} else {
								response.Message = "komentar tidak ditemukan"
							}
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
	}
	return ReturnStruct(response)
}

func HapusKomentar(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var datakomentar Komentar
	var databerita Berita
	databerita.ID = datakomentar.ID_berita

	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenname := DecodeGetName(os.Getenv(publickey), header)
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenname == "" || tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if idBeritaExists(mongoenv, dbname, databerita) {
						if datakomentar.ID == "" {
							response.Message = "parameter dari function ini adalah id"
						} else {
							namakomentator := FindKomentar(mconn, collname, datakomentar)
							if tokenrole == "admin" || tokenname == namakomentator.Name {
								if idKomentarExists(mongoenv, dbname, datakomentar) {
									DeleteKomentar(mconn, collname, datakomentar)
									response.Status = true
									response.Message = "berhasil hapus " + datakomentar.ID + " dari database"
								} else {
									response.Message = "komentar tidak ditemukan"
								}
							} else {
								response.Message = "anda tidak memiliki akses"
							}
						}
					} else {
						response.Message = "berita tidak ditemukan"
					}
				} else {
					response.Message = "akun tidak ditemukan"
				}
			}
		}
	}
	return ReturnStruct(response)
}

func UpdateKomentar(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var auth User
	var datakomentar Komentar
	var databerita Berita
	databerita.ID = datakomentar.ID_berita

	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	header := r.Header.Get("token")

	if header == "" {
		response.Message = "header login tidak ditemukan"
	} else {
		if err != nil {
			response.Message = "error parsing application/json: " + err.Error()
		} else {
			tokenname := DecodeGetName(os.Getenv(publickey), header)
			tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
			tokenrole := DecodeGetRole(os.Getenv(publickey), header)

			auth.Username = tokenusername

			if tokenname == "" || tokenusername == "" || tokenrole == "" {
				response.Message = "hasil decode tidak ditemukan"
			} else {
				if usernameExists(mongoenv, dbname, auth) {
					if datakomentar.ID == "" || datakomentar.Name == "" {
						response.Message = "parameter dari function ini adalah id"
					} else {
						if idBeritaExists(mongoenv, dbname, databerita) {
							namakomentator := FindKomentar(mconn, collname, datakomentar)
							if tokenrole == "admin" || tokenname == namakomentator.Name {
								if idKomentarExists(mongoenv, dbname, datakomentar) {
									EditKomentar(mconn, collname, datakomentar)
									response.Status = true
									datakomentar.Name = tokenname
									response.Message = "berhasil update " + datakomentar.ID + " dari database"
								} else {
									response.Message = "komentar tidak ditemukan"
								}
							} else {
								response.Message = "anda tidak memiliki akses"
							}
						} else {
							response.Message = "berita tidak ditemukan"
						}
					}
				} else {
					response.Message = "akun tidak ditemukan"
				}
			}
		}
	}
	return ReturnStruct(response)
}
