package peda

type Properties struct {
	Name string `json:"name" bson:"name"`
}

type User struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	Private  string `json:"private,omitempty" bson:"private,omitempty"`
	Publick  string `json:"publick,omitempty" bson:"publick,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Name    string `json:"nam,omitemptye" bson:"name,omitempty"`
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
