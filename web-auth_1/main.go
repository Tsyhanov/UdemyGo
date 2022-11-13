package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	password []byte
	First    string
}

//for jwt
type customClaims struct {
	jwt.StandardClaims
	SID string
}

var key = []byte("my_secret_string")
var db = map[string]user{}        //to store email and pass
var session = map[string]string{} //sessionid and email

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)

	http.ListenAndServe(":8080", nil)
	fmt.Println("HTTP Server started")
}

func index(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("sessionID")
	if err != nil {
		c = &http.Cookie{
			Name:  "SessionID",
			Value: "",
		}
	}
	//check cookie
	s, err := parseToken(c.Value)
	if err != nil {
		log.Println("parseToken error", err)
	}
	var e string
	if s != "" {
		e = session[s]
	}

	var f string
	if user, ok := db[e]; ok {
		f = user.First
	}

	errMsg := r.FormValue("errormsg")
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
		<center> <h1> Your name: %s </h1> </center>  
		<center> <h1> Your email: %s </h1> </center>  
		<center> <h1> The error: %s </h1> </center>  
		<center> <h1> Register Form </h1> </center>  
		<form action="/register" method="POST">  
			<div class="container">   
				<label>Email : </label>   
				<input type="email" placeholder="Enter Email" name="email" required> 
				<label>First : </label>   
				<input type="text" placeholder="Enter First Name" name="first" required>  
				<label>Password : </label>   
				<input type="password" placeholder="Enter Password" name="password" required>  
				<button type="submit">Login</button>   
			</div>   
		</form> 
		<h1>LOG IN</h1>
		<form action="/login" method="POST">
			<div class="container">   
				<label>Email : </label>   
				<input type="email" placeholder="Enter Email" name="email" required>  
				<label>Password : </label>   
				<input type="password" placeholder="Enter Password" name="password" required>  
				<button type="submit">Login</button>   
		</form>>    
	</body>
	</html>
	`, f, e, errMsg)
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		errorMsg := url.QueryEscape("method not a POST")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	e := r.FormValue("email")
	if e == "" {
		errorMsg := url.QueryEscape("email is empty")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}
	p := r.FormValue("password")
	if p == "" {
		errorMsg := url.QueryEscape("password is empty")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	f := r.FormValue("first")
	if f == "" {
		errorMsg := url.QueryEscape("first name is empty")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
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

	db[e] = user{
		password: bsp,
		First:    f,
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//login method
func login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		errorMsg := url.QueryEscape("method not a POST")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	e := r.FormValue("email")
	if e == "" {
		errorMsg := url.QueryEscape("email is empty")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}
	p := r.FormValue("password")
	if p == "" {
		errorMsg := url.QueryEscape("password is empty")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	//check if email is in use
	if _, ok := db[e]; !ok {
		errorMsg := url.QueryEscape("email is new")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	err := bcrypt.CompareHashAndPassword(db[e].password, []byte(p))
	if err != nil {
		errorMsg := url.QueryEscape("Your email or password is incorrect")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}

	//create session id and store it by cookie
	sUUID := uuid.New().String()
	session[sUUID] = e
	token, err := createToken(sUUID)
	if err != nil {
		errorMsg := url.QueryEscape("Could nor create token")
		http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
		return
	}
	//create cookie
	c := http.Cookie{
		Name:  "sessionID",
		Value: token,
	}
	//set cookie
	http.SetCookie(w, &c)

	errorMsg := url.QueryEscape("You logged in " + e)
	http.Redirect(w, r, "/?errormsg="+errorMsg, http.StatusSeeOther)
}

func createToken(sid string) (string, error) {

	cc := customClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		SID: sid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	st, err := token.SignedString(key) //signed token
	if err != nil {
		return "", fmt.Errorf("jwt signing error")

	}
	return st, nil

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(sid))
	//signature
	//must be in printable format (hex or base64)
	//hex
	//signedMac := fmt.Sprintf("%x", mac.Sum(nil))
	//base64
	signedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signedMac + "|" + sid, nil
}

//take signed string and get session id from there
func parseToken(ss string) (string, error) {
	//jwt
	token, err := jwt.ParseWithClaims(ss, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("parseWithClaims has different algorithms")
		}
		return key, nil
	})

	if err != nil {
		return "", fmt.Errorf("could not parse token %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid %w", err)
	}

	return token.Claims.(*customClaims).SID, nil

	xs := strings.SplitN(ss, "|", 2)
	if len(xs) != 2 {
		return "", fmt.Errorf("wrong number in string parse token")
	}
	b64 := xs[0]
	xb, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("could not parse token and decode string %w", err)
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(xs[1]))

	ok := hmac.Equal(xb, mac.Sum(nil))
	if !ok {
		return "", fmt.Errorf("not equal signed sid and session id")
	}

	return xs[1], nil
}
