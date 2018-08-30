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
)

const algorithm = "algorithm"
const nonce = "nonce"
const opaque = "opaque"
const qop = "qop"
const realm = "realm"
const response = "response"
const uri = "uri"
var wanted = []string{algorithm, nonce, response,  opaque, qop, realm, uri}

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

func Auth(c *revel.Controller)revel.Result{
	correctUsername := "admin"
	correctPassword := "admin"
	realm := "testrealm@host.com"


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

		ha1 := getMD5([]string{correctUsername, realm, correctPassword})
		ha2 := getMD5([]string{"GET", parts[uri] })
		correctResponse := getMD5([]string{ha1,parts[nonce],ha2})

		fmt.Println(parts)
		fmt.Println("NONCE: ", nonce)
		fmt.Println("CORRECT: "+correctResponse)
		fmt.Println("RECIVED: "+parts[response])

		if correctResponse != parts[response]{
			c.Response.Status = http.StatusUnauthorized
			c.Response.Out.Header().Set("WWW-Authenticate", `Basic realm="revel"`)
			return c.RenderError(errors.New("401: Not authorized"))
		}
		return nil
	} else {

		c.Response.Status = http.StatusUnauthorized

		responseValue := `Digest `
		responseValue += ` realm="`		+ realm 			+`"`
		responseValue += `, nonce="` 	+ getRandomHash()	+`" `
		responseValue += `, opaque="`	+ getRandomHash()	+`" `
		fmt.Println(responseValue)
		c.Response.Out.Header().Set("WWW-Authenticate", responseValue)
		return c.RenderError(errors.New("401: Not authorized"))
	}
}

func LogOut(c *revel.Controller){

	c.Response.Status = http.StatusUnauthorized
	c.Response.Out.Header().Set("WWW-Authenticate", "")
}
