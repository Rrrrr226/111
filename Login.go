package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"net/http"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/dblogin?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
	/*u0 := User{Username: "Admin", Password: "hdu123"}
	result := db.Create(&u0)
	if result.Error != nil {
		fmt.Println(result.Error)
	} else {
		fmt.Printf("User created with ID: %v\n", u0.ID)
	}*/
}

// 检查是否匹配
func findUser(Username string) *User {
	var user User
	db.Where("username=?", Username).First(&user)
	if user.Username != "" {
		return &user
	}
	return nil
}

func authenticate(Username, Password string) bool {
	user := findUser(Username)
	return user != nil && user.Password == Password
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.PostFormValue("Username")
	password := r.PostFormValue("Password")
	if username == "" || password == "" {
		http.Error(w, "Username or password is empty", http.StatusBadRequest)
		return
	}
	if authenticate(username, password) {
		http.Redirect(w, r, "/hello", http.StatusFound)
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./Hello.txt")
	if err != nil {
		http.Error(w, "Failed to hello message", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(b))
}

func registerhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	username := r.PostFormValue("Username")
	password := r.PostFormValue("Password")
	if username == "" || password == "" {
		http.Error(w, "Username or password is empty", http.StatusBadRequest)
		return
	}
	if findUser(username) != nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	newuser := User{Username: username, Password: password}
	result := db.Create(&newuser)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/register", registerhandler)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
