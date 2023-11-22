package peda

type Properties struct {
	Name string `json:"name" bson:"name"`
}

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

type AuthorizationStruct struct {
	Status  bool     `json:"status" bson:"status"`
	Data    UserAuth `json:"data,omitempty" bson:"data,omitempty"`
	Message string   `json:"message,omitempty" bson:"message,omitempty"`
	Token   string   `json:"token,omitempty" bson:"token,omitempty"`
}

type UserAuth struct {
	Name     string    `json:"name,omitempty" bson:"name,omitempty"`
	Email    string    `json:"email,omitempty" bson:"email,omitempty"`
	Username string    `json:"username,omitempty" bson:"username,omitempty"`
	Role     SemuaRole `json:"role,omitempty" bson:"role,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Name    string `json:"name,omitemptye" bson:"name,omitempty"`
	Role    string `json:"role,omitempty" bson:"role,omitempty"`
}

type ResponseDataUser struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Data    []User `json:"data,omitempty" bson:"data,omitempty"`
}

type Response struct {
	Token string `json:"token,omitempty" bson:"token,omitempty"`
}

type Jaja struct {
	Token string      `json:"token" bson:"token"`
	Data  interface{} `json:"data" bson:"data"`
}

type Berita struct {
	ID       string `json:"id" bson:"id"`
	Kategori string `json:"kategori" bson:"kategori"`
	Judul    string `json:"judul" bson:"judul"`
	Preview  string `json:"preview" bson:"preview"`
	Konten   string `json:"konten" bson:"konten"`
}
