package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)

	http.ListenAndServe(":8080", nil)
	fmt.Println("HTTP Server started")
}

func index(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.FormValue("errormsg")
	fmt.Fprintf(w, `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    <center> <h1> Login Form </h1> </center>  
    <form action="/register" method="POST">  
        <div class="container">   
            <label>Email : </label>   
            <input type="email" placeholder="Enter Email" name="email" required>  
            <label>Password : </label>   
            <input type="password" placeholder="Enter Password" name="password" required>  
            <button type="submit">Login</button>   
			<center> <h3> %s </h3> </center>  
        </div>   
    </form>     
</body>
</html>
	`, errorMsg)
}

var db = map[string][]byte{} //to store email and pass

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		errorMsg := url.QueryEscape("method not a POST")
		http.Redirect(w, r, "/&errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	e := r.FormValue("email")
	if e == "" {
		errorMsg := url.QueryEscape("email is empty")
		http.Redirect(w, r, "/&errormsg="+errorMsg, http.StatusSeeOther)
		return
	}
	p := r.FormValue("password")
	if p == "" {
		errorMsg := url.QueryEscape("password is empty")
		http.Redirect(w, r, "/&errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	bsp, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		errorMsg := url.QueryEscape("internal erver error")
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}
	log.Println("email", e)
	log.Println("pass", p)
	log.Println("bcrypted pass", bsp)

	db[e] = bsp
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
