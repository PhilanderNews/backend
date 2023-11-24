package peda

type User struct {
	Name     string    `json:"name,omitempty" bson:"name,omitempty"`
	Email    string    `json:"email,omitempty" bson:"email,omitempty"`
	Username string    `json:"username" bson:"username"`
	Password string    `json:"password" bson:"password"`
	Role     SemuaRole `json:"role,omitempty" bson:"role,omitempty"`
}

type SemuaRole struct {
	Admin  bool `json:"admin" bson:"admin"`
	Author bool `json:"author" bson:"author"`
	User   bool `json:"user" bson:"user"`
}

type CredentialUser struct {
	Status  bool   `json:"status" bson:"status"`
	Data    User   `json:"data,omitempty" bson:"data,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
}
type CredentialBerita struct {
	Status  bool   `json:"status" bson:"status"`
	Data    Berita `json:"data,omitempty" bson:"data,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
}

type Berita struct {
	ID       string   `json:"id" bson:"id"`
	Kategori string   `json:"kategori" bson:"kategori"`
	Judul    string   `json:"judul" bson:"judul"`
	Preview  string   `json:"preview" bson:"preview"`
	Konten   Paragraf `json:"konten" bson:"konten"`
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
