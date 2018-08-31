package auth

import (
	"strings"
	"fmt"
	"crypto/md5"
	"io"
	"encoding/hex"
	"github.com/revel/revel"
	"net/http"
	"errors"
	"math/rand"
	"app/app/lib/dbManager"
	"database/sql"
)

const algorithm = "algorithm"
const nonce = "nonce"
const opaque = "opaque"
const qop = "qop"
const realm = "realm"
const response = "response"
const uri = "uri"
const username = "username"
var wanted = []string{username ,algorithm, nonce, response,  opaque, qop, realm, uri}

func getMD5(texts []string) string {
	h := md5.New()
	io.WriteString(h, strings.Join(texts, ":"))
	return hex.EncodeToString(h.Sum(nil))
}

func getRandomHash()string{
	h := md5.New()
	randNumber := rand.Intn(100)
	io.WriteString(h, string(randNumber))
	return  hex.EncodeToString(h.Sum(nil))
}

func Auth(c *revel.Controller)(revel.Result, error){
	nonceParam 	:= getRandomHash()
	opaqueParam := getRandomHash()

	if auth := c.Request.Header.Get("Authorization"); auth != "" {
		headers := strings.Split(auth, ",")
		parts := make(map[string]string, len(wanted))
		for _, r := range headers {
			for _, w := range wanted {
				if strings.Contains(r, w) {
					parts[w] = strings.Split(r, `"`)[1]
				}
			}
		}

		db, err := dbManager.OpenConnection()
		if err != nil{
			return nil, err
		}
		defer db.Close()

		req := "SELECT c_password FROM airport.users WHERE c_username = $1"
		rows, err := db.Query(req, parts[username])
		if err != nil{
			if err == sql.ErrNoRows{
				return repeatRequest(c, nonceParam, opaqueParam)
			}

			return repeatRequest(c, nonceParam, opaqueParam)
		}
		defer rows.Close()
		var correctPassword string
		for rows.Next(){
			if err := rows.Scan(&correctPassword); err != nil{
				return repeatRequest(c, nonceParam, opaqueParam)
			}
		}
		if strings.Trim(correctPassword," ") == ""{
			return repeatRequest(c, nonceParam, opaqueParam)
		}
		fmt.Println("CORRECT PASSWORD: " + correctPassword)
		ha1 := getMD5([]string{parts[username], parts[realm], correctPassword})
		ha2 := getMD5([]string{c.Request.Method, parts[uri] })
		correctResponse := getMD5([]string{ha1,parts[nonce],ha2})

		fmt.Println(parts)
		fmt.Println("CORRECT: "+correctResponse)
		fmt.Println("RECIVED: "+parts[response])

		if correctResponse != parts[response]{
			return repeatRequest(c, nonceParam, opaqueParam)
		}
		return nil, nil
	} else {
		return repeatRequest(c, nonceParam, opaqueParam)

	}
}

func repeatRequest(c *revel.Controller, nonceParam string, opaqueParam string) (revel.Result,error){
	c.Response.Status = http.StatusUnauthorized

	responseValue := `Digest `
	responseValue += ` realm="`		+ "testrealm@host.com"	+`"`
	responseValue += `, nonce="` 	+ nonceParam			+`" `
	responseValue += `, opaque="`	+ opaqueParam			+`" `
	fmt.Println(responseValue)
	c.Response.Out.Header().Set("WWW-Authenticate", responseValue)
	return c.RenderError(errors.New("401: Not authorized")), nil
}

func LogOut(c *revel.Controller){
	c.Response.Status = http.StatusUnauthorized
	c.Request.Header.Del("Authorization")
	c.Response.Out.Header().Del("WWW-Authenticate")
}
