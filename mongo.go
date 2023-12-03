package peda

import (
	"context"
	"os"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetConnection(mongoenv, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(mongoenv),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func SetConnectionTest(mongoenv, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: mongoenv,
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

// ---------------------------------------------------------------------- User

// Create

func InsertUserdata(mongoenv *mongo.Database, collname, name, email, no_whatsapp, username, password, role string) (InsertedID interface{}) {
	req := new(User)
	req.Name = name
	req.Email = email
	req.No_whatsapp = no_whatsapp
	req.Username = username
	req.Password = password
	req.Role = role
	return atdb.InsertOneDoc(mongoenv, collname, req)
}

// Read

func GetAllUser(mongoenv *mongo.Database, collname string) []User {
	user := atdb.GetAllDoc[[]User](mongoenv, collname)
	return user
}

func FindUser(mongoenv *mongo.Database, collname string, userdata User) User {
	filter := bson.M{"username": userdata.Username}
	return atdb.GetOneDoc[User](mongoenv, collname, filter)
}

func IsPasswordValid(mongoenv *mongo.Database, collname string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mongoenv, collname, filter)
	hashChecker := CheckPasswordHash(userdata.Password, res.Password)
	return hashChecker
}

func usernameExists(mongoenv, dbname string, userdata User) bool {
	mconn := SetConnection(mongoenv, dbname).Collection("user")
	filter := bson.M{"username": userdata.Username}

	var user User
	err := mconn.FindOne(context.Background(), filter).Decode(&user)
	return err == nil
}

// Update

func EditUser(mongoenv *mongo.Database, collname string, datauser User) interface{} {
	filter := bson.M{"username": datauser.Username}
	return atdb.ReplaceOneDoc(mongoenv, collname, filter, datauser)
}

// Delete

func DeleteUser(mongoenv *mongo.Database, collname string, userdata User) interface{} {
	filter := bson.M{"username": userdata.Username}
	return atdb.DeleteOneDoc(mongoenv, collname, filter)
}

//-------------------------------------------------------------------- Berita

// Create

func InsertBerita(mongoenv *mongo.Database, collname string, databerita Berita) interface{} {
	return atdb.InsertOneDoc(mongoenv, collname, databerita)
}

// Read

func GetAllBerita(mongoenv *mongo.Database, collname string) []Berita {
	berita := atdb.GetAllDoc[[]Berita](mongoenv, collname)
	return berita
}

func FindBerita(mongoenv *mongo.Database, collname string, databerita Berita) Berita {
	filter := bson.M{"id": databerita.ID}
	return atdb.GetOneDoc[Berita](mongoenv, collname, filter)
}

func idBeritaExists(mongoenv, dbname string, databerita Berita) bool {
	mconn := SetConnection(mongoenv, dbname).Collection("berita")
	filter := bson.M{"id": databerita.ID}

	var berita Berita
	err := mconn.FindOne(context.Background(), filter).Decode(&berita)
	return err == nil
}

// Update

func EditBerita(mongoenv *mongo.Database, collname string, databerita Berita) interface{} {
	filter := bson.M{"id": databerita.ID}
	return atdb.ReplaceOneDoc(mongoenv, collname, filter, databerita)
}

// Delete

func DeleteBerita(mongoenv *mongo.Database, collname string, databerita Berita) interface{} {
	filter := bson.M{"id": databerita.ID}
	return atdb.DeleteOneDoc(mongoenv, collname, filter)
}

//-------------------------------------------------------------------- Komentar

// Create

func InsertKomentar(mongoenv *mongo.Database, collname string, datakomentar Komentar) interface{} {
	return atdb.InsertOneDoc(mongoenv, collname, datakomentar)
}

// Read

func GetAllKomentar(mongoenv *mongo.Database, collname string) []Komentar {
	komentar := atdb.GetAllDoc[[]Komentar](mongoenv, collname)
	return komentar
}

func FindKomentar(mongoenv *mongo.Database, collname string, datakomentar Komentar) Komentar {
	filter := bson.M{"id": datakomentar.ID}
	return atdb.GetOneDoc[Komentar](mongoenv, collname, filter)
}

func idKomentarExists(mongoenv, dbname string, datakomentar Komentar) bool {
	mconn := SetConnection(mongoenv, dbname).Collection("komentar")
	filter := bson.M{"id": datakomentar.ID}

	var komentar Komentar
	err := mconn.FindOne(context.Background(), filter).Decode(&komentar)
	return err == nil
}

// Update

func EditKomentar(mongoenv *mongo.Database, collname string, datakomentar Komentar) interface{} {
	filter := bson.M{"id": datakomentar.ID}
	return atdb.ReplaceOneDoc(mongoenv, collname, filter, datakomentar)
}

// Delete

func DeleteKomentar(mongoenv *mongo.Database, collname string, datakomentar Komentar) interface{} {
	filter := bson.M{"id": datakomentar.ID}
	return atdb.DeleteOneDoc(mongoenv, collname, filter)
}

// ---------------------------------------------------------------------- Tutorial

func InsertMongo(mongoenv *mongo.Database, collname string, pesan Tutorial) interface{} {
	return atdb.InsertOneDoc(mongoenv, collname, pesan)
}

func GetAllMongo(mongoenv *mongo.Database, collname string) []Tutorial {
	data := atdb.GetAllDoc[[]Tutorial](mongoenv, collname)
	return data
}

func GetOneMongo(mongoenv *mongo.Database, collname string, datapesan Tutorial) Tutorial {
	filter := bson.M{"parameter": datapesan.Parameter}
	return atdb.GetOneDoc[Tutorial](mongoenv, collname, filter)
}

func UpdateMongo(mongoenv *mongo.Database, collname string, datapesan Tutorial) interface{} {
	filter := bson.M{"parameter": datapesan.Parameter}
	return atdb.ReplaceOneDoc(mongoenv, collname, filter, datapesan)
}

func DeleteMongo(mongoenv *mongo.Database, collname string, datapesan Tutorial) interface{} {
	filter := bson.M{"parameter": datapesan.Parameter}
	return atdb.DeleteOneDoc(mongoenv, collname, filter)
}
