package odata4

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/parnurzeal/gorequest"
)

type (
	Config struct {
		ODATA4_URL  string `env:"ODATA4_URL"`
		ODATA4_USER string `env:"ODATA4_USER"`
		ODATA4_PASS string `env:"ODATA4_PASS"`
	}

	LoginData struct {
		UserName     string `josn:"UserName"`
		UserPassword string `json:"UserPassword"`
	}

	LoginResponse struct {
		Code              int     `json:"Code"`
		Message           string  `json:"Message"`
		Exception         *string `json:"Exception"`
		PasswordChangeUrl *string `json:"PasswordChangeUrl"`
		RedirectUrl       *string `json:"RedirectUrl"`
	}
)

var c Config
var Cookies []*http.Cookie
var lock = &sync.Mutex{}

const COOKIE_BPMCSRF = "BPMCSRF"

func init() {
	env.Parse(&c)
	Cookies = GetConnectionCookies(true)

}

func GetConnectionCookies(forceRefresh bool) []*http.Cookie {
	lock.Lock()
	defer lock.Unlock()

	if forceRefresh {
		Cookies = getCookies()
		return Cookies
	}

	if Cookies == nil {
		log.Println("Cookies not found")
		Cookies = getCookies()
	} else if !checkCookies(Cookies) {
		log.Println("Recreate Cookies")
		Cookies = getCookies()
	}

	return Cookies
}

func getCookies() []*http.Cookie {
	log.Println("Getting cookies")

	loginData := LoginData{
		UserName:     c.ODATA4_USER,
		UserPassword: c.ODATA4_PASS,
	}

	var loginResponse LoginResponse
	url := getPath("/ServiceModel/AuthService.svc/Login")

	request := gorequest.New()
	resp, _, err := request.Post(url).
		Send(loginData).
		EndStruct(&loginResponse)

	if err != nil {
		CheckError("Try to Login", err[0])
	}

	checkLogin(loginResponse)

	Cookies = (*http.Response)(resp).Cookies()
	if !checkCookies(Cookies) {
		log.Fatalln("Cann't get Cookies")
	}

	log.Println("Cookies created")

	return Cookies
}

func checkCookies(c []*http.Cookie) bool {
	return checkIfCookieByNameExists(c, COOKIE_BPMCSRF) && isCookieNotExpired(c)
}

func isCookieNotExpired(c []*http.Cookie) bool {
	for _, item := range c {
		if item.RawExpires != "" && item.Expires.After(time.Now().UTC()) {
			return true
		}

	}

	return false
}

func checkIfCookieByNameExists(c []*http.Cookie, name string) bool {
	for _, item := range c {
		if item.Name == name && item.Value != "" {
			return true
		}
	}

	return false
}

func checkLogin(lr LoginResponse) {
	if lr.Code != 0 {
		log.Fatalln(time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), "Code: ", lr.Code, "Message: ", lr.Message)
	}
}

func getParamsFromCookies(param string) string {
	for _, cookie := range Cookies {
		if cookie.Name == param {
			return cookie.Value
		}
	}

	return ""
}
