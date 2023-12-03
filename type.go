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
	Paragraf1  string `json:"paragraf1,omitempty" bson:"paragraf1,omitempty"`
	Paragraf2  string `json:"paragraf2,omitempty" bson:"paragraf2,omitempty"`
	Paragraf3  string `json:"paragraf3,omitempty" bson:"paragraf3,omitempty"`
	Paragraf4  string `json:"paragraf4,omitempty" bson:"paragraf4,omitempty"`
	Paragraf5  string `json:"paragraf5,omitempty" bson:"paragraf5,omitempty"`
	Paragraf6  string `json:"paragraf6,omitempty" bson:"paragraf6,omitempty"`
	Paragraf7  string `json:"paragraf7,omitempty" bson:"paragraf7,omitempty"`
	Paragraf8  string `json:"paragraf8,omitempty" bson:"paragraf8,omitempty"`
	Paragraf9  string `json:"paragraf9,omitempty" bson:"paragraf9,omitempty"`
	Paragraf10 string `json:"paragraf10,omitempty" bson:"paragraf10,omitempty"`
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
