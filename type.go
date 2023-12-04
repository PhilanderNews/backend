package peda

import (
	"time"
)

type Payload struct {
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Exp      time.Time `json:"exp"`
	Iat      time.Time `json:"iat"`
	Nbf      time.Time `json:"nbf"`
}

type User struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	No_whatsapp string `json:"no_whatsapp,omitempty" bson:"no_whatsapp"`
	Username    string `json:"username" bson:"username"`
	Password    string `json:"password,omitempty" bson:"password"`
	Role        string `json:"role,omitempty" bson:"role,omitempty"`
}

type CredentialUser struct {
	Status  bool   `json:"status" bson:"status"`
	Data    User   `json:"data,omitempty" bson:"data,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Pesan struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
}

type Berita struct {
	ID       string   `json:"id" bson:"id"`
	Kategori string   `json:"kategori" bson:"kategori"`
	Judul    string   `json:"judul" bson:"judul"`
	Preview  string   `json:"preview" bson:"preview"`
	Konten   Paragraf `json:"konten" bson:"konten"`
	Penulis  string   `json:"penulis" bson:"penulis"`
	Sumber   string   `json:"sumber" bson:"sumber"`
	Image    string   `json:"image" bson:"image"`
	Waktu    string   `json:"waktu" bson:"waktu"`
}

type Paragraf struct {
	Paragraf1  string `json:"paragraf1" bson:"paragraf1"`
	Paragraf2  string `json:"paragraf2" bson:"paragraf2"`
	Paragraf3  string `json:"paragraf3" bson:"paragraf3"`
	Paragraf4  string `json:"paragraf4" bson:"paragraf4"`
	Paragraf5  string `json:"paragraf5" bson:"paragraf5"`
	Paragraf6  string `json:"paragraf6" bson:"paragraf6"`
	Paragraf7  string `json:"paragraf7" bson:"paragraf7"`
	Paragraf8  string `json:"paragraf8" bson:"paragraf8"`
	Paragraf9  string `json:"paragraf9" bson:"paragraf9"`
	Paragraf10 string `json:"paragraf10" bson:"paragraf10"`
}

type Komentar struct {
	ID        string `json:"id" bson:"id"`
	ID_berita string `json:"id_berita" bson:"id_berita"`
	Name      string `json:"name" bson:"name"`
	Tanggal   string `json:"tanggal" bson:"tanggal"`
	Komentar  string `json:"komentar" bson:"komentar"`
}

// ---------------------------------------------------------------------- Tutorial

type Tutorial struct {
	Parameter string `json:"parameter" bson:"parameter"`
	Pesan     string `json:"pesan" bson:"pesan"`
}
