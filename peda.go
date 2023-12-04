package peda

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

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
	var auth User
	response.Status = false

	// Extract token from the request header
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode token values
	tokenname := DecodeGetName(os.Getenv(publickey), header)
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	// Create User struct with the decoded username
	auth.Username = tokenusername

	// Check if decoding results are valid
	if tokenname == "" || tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Successful token decoding and user validation
	response.Message = "Berhasil decode token"
	response.Status = true
	response.Data.Name = tokenname
	response.Data.Username = tokenusername
	response.Data.Role = tokenrole

	return ReturnStruct(response)
}

func Registrasi(token, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode user data from the request body
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Check if the username already exists
	if usernameExists(mongoenv, dbname, datauser) {
		response.Message = "Username telah dipakai"
		return ReturnStruct(response)
	}

	// Hash the user's password
	hash, hashErr := HashPassword(datauser.Password)
	if hashErr != nil {
		response.Message = "Gagal hash password: " + hashErr.Error()
		return ReturnStruct(response)
	}

	// Check if the 'No_whatsapp' field is empty
	if datauser.No_whatsapp == "" {
		response.Message = "Nomor WhatsApp wajib diisi"
		return ReturnStruct(response)
	}

	// Insert user data into the database
	InsertUserdata(mconn, collname, datauser.Name, datauser.Email, datauser.No_whatsapp, datauser.Username, hash, datauser.Role)
	response.Status = true
	response.Message = "Berhasil input data"

	// Prepare and send a WhatsApp message with registration details
	var username = datauser.Username
	var password = datauser.Password
	var nohp = datauser.No_whatsapp

	dt := &wa.TextMessage{
		To:       nohp,
		IsGroup:  false,
		Messages: "Selamat anda berhasil registrasi, berikut adalah username anda: " + username + "\nDan ini adalah password anda: " + password + "\nDisimpan baik baik ya",
	}

	// Make an API call to send WhatsApp message
	atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")

	return ReturnStruct(response)
}

func Login(token, privatekey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode user data from the request body
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(mongoenv, dbname, datauser) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the entered password is not valid
	if !IsPasswordValid(mconn, collname, datauser) {
		response.Message = "Password Salah"
		return ReturnStruct(response)
	}

	// Retrieve user details
	user := FindUser(mconn, collname, datauser)

	// Prepare and encode token
	tokenstring, tokenerr := Encode(user.Name, user.Username, user.Role, os.Getenv(privatekey))
	if tokenerr != nil {
		response.Message = "Gagal encode token: " + tokenerr.Error()
		return ReturnStruct(response)
	}

	// Successful login
	response.Status = true
	response.Token = tokenstring
	response.Message = "Berhasil login"

	// Send a WhatsApp message notifying the user about the successful login
	var nama = user.Name
	var nohp = user.No_whatsapp
	dt := &wa.TextMessage{
		To:       nohp,
		IsGroup:  false,
		Messages: nama + " berhasil login\nPerlu diingat sesi login hanya berlaku 2 jam",
	}
	atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv(token), dt, "https://api.wa.my.id/api/send/message/text")

	return ReturnStruct(response)
}

func AmbilSatuUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode user data from the request body
	var auth User
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode token to get username and role
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	// Check if decoding was successful
	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user has admin privileges
	if tokenrole != "admin" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Check if the user to be fetched exists
	if usernameExists(mongoenv, dbname, datauser) {
		user := FindUser(mconn, collname, datauser)
		return ReturnStruct(user)
	} else {
		response.Message = "User tidak ditemukan"
		return ReturnStruct(response)
	}
}

func AmbilSemuaUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode token to get username and role
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	// Check if decoding was successful
	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(mongoenv, dbname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user has admin privileges
	if tokenrole != "admin" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Get all users if the user is an admin
	datauser := GetAllUser(mconn, collname)
	return ReturnStruct(datauser)
}

func UpdateUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode user data from the request body
	var auth User
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode token to get username and role
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	// Check if decoding was successful
	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user has admin privileges
	if tokenrole != "admin" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Check if the username parameter is provided
	if datauser.Username == "" {
		response.Message = "Parameter dari function ini adalah username"
		return ReturnStruct(response)
	}

	// Check if the user to be edited exists
	if !usernameExists(mongoenv, dbname, datauser) {
		response.Message = "Akun yang ingin diedit tidak ditemukan"
		return ReturnStruct(response)
	}

	// Hash the user's password if provided
	if datauser.Password != "" {
		hash, hashErr := HashPassword(datauser.Password)
		if hashErr != nil {
			response.Message = "Gagal Hash Password: " + hashErr.Error()
			return ReturnStruct(response)
		}
		datauser.Password = hash
	} else {
		// Retrieve user details
		user := FindUser(mconn, collname, datauser)
		datauser.Password = user.Password
	}

	// Perform user update
	EditUser(mconn, collname, datauser)

	response.Status = true
	response.Message = "Berhasil update " + datauser.Username + " dari database"
	return ReturnStruct(response)
}

func HapusUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode user data from the request body
	var auth User
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode token to get username and role
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	// Check if decoding was successful
	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user has admin privileges
	if tokenrole != "admin" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Check if the username parameter is provided
	if datauser.Username == "" {
		response.Message = "Parameter dari function ini adalah username"
		return ReturnStruct(response)
	}

	// Check if the user to be deleted exists
	if !usernameExists(mongoenv, dbname, datauser) {
		response.Message = "Akun yang ingin dihapus tidak ditemukan"
		return ReturnStruct(response)
	}

	// Perform user deletion
	DeleteUser(mconn, collname, datauser)

	response.Status = true
	response.Message = "Berhasil hapus " + datauser.Username + " dari database"
	return ReturnStruct(response)
}

// ---------------------------------------------------------------------- Berita

func TambahBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode berita data from the request body
	var databerita Berita
	err := json.NewDecoder(r.Body).Decode(&databerita)

	//Define waktu
	wib, timeErr := time.LoadLocation("Asia/Jakarta")

	if timeErr != nil {
		response.Message = "Error parsing time location: " + err.Error()
		return ReturnStruct(response)
	}

	currentTime := time.Now().In(wib)
	timeStringBerita := currentTime.Format("Monday, 2 January 2006 15:04 MST")

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Get token and perform basic token validation
	var auth User
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode token to get user details
	tokenname := DecodeGetName(os.Getenv(publickey), header)
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)
	auth.Username = tokenusername

	// Check if decoding was successful
	if tokenname == "" || tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user has admin or author privileges
	if tokenrole != "admin" && tokenrole != "author" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Check if the berita ID already exists
	if idBeritaExists(mongoenv, dbname, databerita) {
		response.Message = "ID telah ada"
		return ReturnStruct(response)
	}

	// Insert berita data into the database
	response.Status = true
	databerita.Penulis = tokenname
	databerita.Waktu = timeStringBerita
	InsertBerita(mconn, collname, databerita)
	response.Message = "Berhasil input data"

	return ReturnStruct(response)
}

func AmbilSatuBerita(mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode berita data from the request body
	var databerita Berita
	err := json.NewDecoder(r.Body).Decode(&databerita)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Check if the berita ID parameter is provided
	if databerita.ID == "" {
		response.Message = "Parameter dari function ini adalah ID"
		return ReturnStruct(response)
	}

	// Check if the berita exists
	if idBeritaExists(mongoenv, dbname, databerita) {
		// Fetch berita data from the database
		berita := FindBerita(mconn, collname, databerita)
		return ReturnStruct(berita)
	} else {
		response.Message = "Berita tidak ditemukan"
	}

	return ReturnStruct(response)
}

func AmbilSemuaBerita(mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Fetch all berita data from the database
	databerita := GetAllBerita(mconn, collname)

	return ReturnStruct(databerita)
}

func UpdateBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode berita data from the request body
	var auth User
	var databerita Berita
	err := json.NewDecoder(r.Body).Decode(&databerita)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode token to get user details
	tokenname := DecodeGetName(os.Getenv(publickey), header)
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	auth.Username = tokenusername

	// Check if decoding was successful
	if tokenname == "" || tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Fetch berita data from the database
	namapenulis := FindBerita(mconn, collname, databerita)

	// Check if the user is not an admin or not the author of the berita
	if !(tokenrole == "admin" || tokenname == namapenulis.Penulis) {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Check if the berita ID parameter is provided
	if databerita.ID == "" {
		response.Message = "Parameter dari function ini adalah ID"
		return ReturnStruct(response)
	}

	// Check if the berita exists
	if idBeritaExists(mongoenv, dbname, databerita) {
		databerita.Penulis = tokenname
		databerita.Waktu = namapenulis.Waktu
		EditBerita(mconn, collname, databerita)
		response.Status = true
		response.Message = "Berhasil update " + databerita.ID + " dari database"
	} else {
		response.Message = "Berita tidak ditemukan"
	}

	return ReturnStruct(response)
}

func HapusBerita(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Decode berita data from the request body
	var auth User
	var databerita Berita
	err := json.NewDecoder(r.Body).Decode(&databerita)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode token to get user details
	tokenname := DecodeGetName(os.Getenv(publickey), header)
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	auth.Username = tokenusername

	// Check if decoding was successful
	if tokenname == "" || tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user account exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Fetch berita data from the database
	namapenulis := FindBerita(mconn, collname, databerita)

	// Check if the user has admin or author privileges
	if !(tokenrole == "admin" || tokenname == namapenulis.Penulis) {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Check if the berita ID parameter is provided
	if databerita.ID == "" {
		response.Message = "Parameter dari function ini adalah ID"
		return ReturnStruct(response)
	}

	// Check if the berita exists
	if idBeritaExists(mongoenv, dbname, databerita) {
		DeleteBerita(mconn, collname, databerita)
		response.Status = true
		response.Message = "Berhasil hapus " + databerita.ID + " dari database"
	} else {
		response.Message = "Berita tidak ditemukan"
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

	// Decode JSON request body into datakomentar
	var auth User
	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	// Define time
	currentTime := time.Now()
	timeStringKomentar := currentTime.Format("January 2, 2006")

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Get token and perform basic token validation
	header := r.Header.Get("token")
	if header == "" {
		response.Status = true
		response.Message = "Berhasil Input data tanpa login"
		datakomentar.Name = "Anonymous"
		datakomentar.Tanggal = timeStringKomentar
		InsertKomentar(mconn, collname, datakomentar)
		return ReturnStruct(response)
	}

	// Decode token to get user details
	tokenname := DecodeGetName(os.Getenv(publickey), header)
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	auth.Username = tokenusername

	if tokenname == "" || tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the komentar ID parameter is provided
	if datakomentar.ID == "" {
		response.Message = "Parameter dari function ini adalah ID"
		return ReturnStruct(response)
	}

	// Check if the komentar ID exists
	if idKomentarExists(mongoenv, dbname, datakomentar) {
		response.Message = "ID telah ada"
		return ReturnStruct(response)
	}

	// Check if the berita ID parameter is provided
	if datakomentar.ID_berita == "" {
		response.Message = "Parameter dari function ini adalah ID Berita"
		return ReturnStruct(response)
	}

	// Set berita ID from komentar data
	databerita.ID = datakomentar.ID_berita

	// Check if the berita exists
	if !idBeritaExists(mongoenv, dbname, databerita) {
		response.Message = "Berita tidak ditemukan"
		return ReturnStruct(response)
	}

	// Insert the komentar data
	response.Status = true
	datakomentar.Name = tokenname
	datakomentar.Tanggal = timeStringKomentar
	InsertKomentar(mconn, collname, datakomentar)
	response.Message = "Berhasil Input data"

	return ReturnStruct(response)
}

func AmbilSatuKomentar(mongoenv, dbname, collname string, r *http.Request) string {
	// Initialize response
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Initialize datakomentar and auth
	var datakomentar Komentar

	// Decode JSON request body into datakomentar
	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Check if the komentar ID parameter is provided
	if datakomentar.ID == "" {
		response.Message = "Parameter dari function ini adalah ID"
		return ReturnStruct(response)
	}

	// Check if the komentar exists
	if !idKomentarExists(mongoenv, dbname, datakomentar) {
		response.Message = "Komentar tidak ditemukan"
		return ReturnStruct(response)
	}

	// Find and return the komentar
	komentar := FindKomentar(mconn, collname, datakomentar)
	return ReturnStruct(komentar)
}

func AmbilSemuaKomentar(mongoenv, dbname, collname string, r *http.Request) string {
	// Initialize response
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Get all komentar data
	datakomentar := GetAllKomentar(mconn, collname)
	return ReturnStruct(datakomentar)
}

func UpdateKomentar(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	// Inisialisasi respons dengan status awal false
	var response Pesan
	response.Status = false

	// Set up koneksi MongoDB
	mconn := SetConnection(mongoenv, dbname)

	// Inisialisasi struktur User dan Komentar
	var auth User
	var datakomentar Komentar

	// Decode body request menjadi struktur Komentar
	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Ambil token dari header request
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode informasi user dari token
	tokenname := DecodeGetName(os.Getenv(publickey), header)
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	// Set informasi user untuk validasi
	auth.Username = tokenusername

	// Validasi informasi user kosong
	if tokenname == "" || tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Validasi keberadaan user di database
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Validasi parameter yang diperlukan
	if datakomentar.ID == "" {
		response.Message = "Parameter dari function ini adalah id"
		return ReturnStruct(response)
	}

	// Validasi keberadaan komentar di database
	if !idKomentarExists(mongoenv, dbname, datakomentar) {
		response.Message = "Komentar tidak ditemukan"
		return ReturnStruct(response)
	}

	// Temukan informasi komentator dari database
	namakomentator := FindKomentar(mconn, collname, datakomentar)

	// Validasi apakah user memiliki akses (admin atau pemilik komentar)
	if tokenname != namakomentator.Name {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Lakukan edit pada komentar
	datakomentar.Name = tokenname
	datakomentar.ID_berita = namakomentator.ID_berita
	datakomentar.Tanggal = namakomentator.Tanggal
	EditKomentar(mconn, collname, datakomentar)

	// Set status respons menjadi true dan tambahkan informasi pada pesan
	response.Status = true
	response.Message = "Berhasil update " + datakomentar.ID + " dari database"

	return ReturnStruct(response)
}

func HapusKomentar(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	// Initialize response
	var response Pesan
	response.Status = false

	// Establish MongoDB connection
	mconn := SetConnection(mongoenv, dbname)

	// Initialize auth and datakomentar
	var auth User
	var datakomentar Komentar

	// Decode JSON request body into datakomentar
	err := json.NewDecoder(r.Body).Decode(&datakomentar)

	// Check for JSON decoding errors
	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	// Get token from request header
	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	// Decode user information from the token
	tokenname := DecodeGetName(os.Getenv(publickey), header)
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	auth.Username = tokenusername

	if tokenname == "" || tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	// Check if the user exists
	if !usernameExists(mongoenv, dbname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	// Find namakomentator based on ID
	namakomentator := FindKomentar(mconn, collname, datakomentar)

	// Check user role for authorization
	if !(tokenrole == "admin" || tokenname == namakomentator.Name) {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	// Check if ID is provided
	if datakomentar.ID == "" {
		response.Message = "Parameter dari function ini adalah id"
		return ReturnStruct(response)
	}

	// Check if the komentar exists
	if !idKomentarExists(mongoenv, dbname, datakomentar) {
		response.Message = "Komentar tidak ditemukan"
		return ReturnStruct(response)
	}

	// Delete the komentar
	DeleteKomentar(mconn, collname, datakomentar)

	// Set response status and message
	response.Status = true
	response.Message = "Berhasil hapus " + datakomentar.ID + " dari database"

	return ReturnStruct(response)
}

// ---------------------------------------------------------------------- Tutorial

func TutorialGCFInsertMongo(mongoenv, dbname, collname string, r *http.Request) string {
	var pesan Pesan
	mconn := SetConnection(mongoenv, dbname)
	var datatest Tutorial

	err := json.NewDecoder(r.Body).Decode(&datatest)

	// Check for JSON decoding errors
	if err != nil {
		pesan.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(pesan)
	}
	InsertMongo(mconn, collname, datatest)
	return ReturnStruct(datatest)
}
