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

func getMD5(texts []string) string {//получить контрольную сумму от массива строк в соответствии с алгоритмом дайджест
									//аутентификации
	h := md5.New()
	io.WriteString(h, strings.Join(texts, ":"))//объединение строк с разделителем ":"
	return hex.EncodeToString(h.Sum(nil))
}

func getRandomHash()string{//получение случайного хэша
	h := md5.New()
	randNumber := rand.Intn(100)
	io.WriteString(h, string(randNumber))
	return  hex.EncodeToString(h.Sum(nil))
}

func Auth(c *revel.Controller)(revel.Result, error){//проверка аутентификации
	nonceParam 	:= getRandomHash()
	opaqueParam := getRandomHash()

	if auth := c.Request.Header.Get("Authorization"); auth != "" {//если заголовок аутентификации не пуст
	//проверяем пароль
		headers := strings.Split(auth, ",")
		parts := make(map[string]string, len(wanted))
		for _, r := range headers {
			for _, w := range wanted {
				if strings.Contains(r, w) {//достаём информацию из заголовка в цикле
					parts[w] = strings.Split(r, `"`)[1]
				}
			}
		}

		db, err := dbManager.OpenConnection()//открытие базы данных
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
		for rows.Next(){//получение пароля из базы данных по имени пользователя
			if err := rows.Scan(&correctPassword); err != nil{
				return repeatRequest(c, nonceParam, opaqueParam)
			}
		}
		if strings.Trim(correctPassword," ") == ""{//проверка на пустоту пароля
			return repeatRequest(c, nonceParam, opaqueParam)//если пароль пуст - повторяем запрос на пароль
		}
		fmt.Println("CORRECT PASSWORD: " + correctPassword)
		ha1 := getMD5([]string{parts[username], parts[realm], correctPassword})//получение хэша правильного пароля
		ha2 := getMD5([]string{c.Request.Method, parts[uri] })
		correctResponse := getMD5([]string{ha1,parts[nonce],ha2})

		fmt.Println(parts)
		fmt.Println("CORRECT: "+correctResponse)
		fmt.Println("RECIVED: "+parts[response])

		if correctResponse != parts[response]{//если пароль неправилный
			return repeatRequest(c, nonceParam, opaqueParam)//повторяе запрос на пароль
		}
		return nil, nil//есои пароль правильный - авторизация завершена
	} else {//если заголовок аутентификации пуст - повторяем запрос на пароль
		return repeatRequest(c, nonceParam, opaqueParam)

	}
}

func repeatRequest(c *revel.Controller, nonceParam string, opaqueParam string) (revel.Result,error){//запрос аутентификации
	c.Response.Status = http.StatusUnauthorized

	responseValue := `Digest `
	responseValue += ` realm="`		+ "testrealm@host.com"	+`"`
	responseValue += `, nonce="` 	+ nonceParam			+`" `
	responseValue += `, opaque="`	+ opaqueParam			+`" `
	fmt.Println(responseValue)
	c.Response.Out.Header().Set("WWW-Authenticate", responseValue)
	return c.RenderError(errors.New("401: Not authorized")), nil
}

func LogOut(c *revel.Controller){//снятие авторизированности
	c.Response.Status = http.StatusUnauthorized//установка статуса "не аутентифицирован"
	c.Request.Header.Del("Authorization")//удаление заголовков авторизации
	c.Response.Out.Header().Del("WWW-Authenticate")//
}
