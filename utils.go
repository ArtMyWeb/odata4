package odata4

import (
	"log"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

type (
	IdNameData struct {
		ID   string `json:"Id,omitempty"`
		Name string `json:"Name,omitempty"`
	}

	IdNameDataResponse struct {
		Value []IdNameData `json:"value"`
	}
)

func getPath(url string) string {
	return c.ODATA4_URL + url
}

func getIdNameData(route string) []IdNameData {
	var data IdNameDataResponse

	url := getPath(route)
	Cookies = GetConnectionCookies()

	request := gorequest.New()
	resp, _, errs := request.Get(url).AddCookies(Cookies).EndStruct(&data)
	if errs != nil {
		CheckError("Get Data By URL: "+url, errs[0])
	}

	if resp.StatusCode == http.StatusOK {
		return data.Value
	}

	return []IdNameData{}
}

func getODataIDByName(url, name string) string {
	items := getIdNameData(url)

	for _, item := range items {
		if item.Name == name {
			return item.ID
		}
	}

	return ""
}

func CheckError(message string, err error) {
	if err != nil {
		log.Fatalln(time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), message, err.Error())
	}
}
