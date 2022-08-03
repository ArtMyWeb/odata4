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

const (
	CURRENCY           = "Currency"
	USR_TYPE_PAYMENT   = "UsrTypePayment"
	USR_STATUS_PAYMENT = "UsrStatusPayment"
	USR_SOURCE_PAYMENT = "UsrSourcePayment"
)

func getPath(url string) string {
	return c.ODATA4_URL + url
}

func GetListOfPaymentSettings() (ListOfPaymentSettings, error) {
	var err error

	currency, errCurrency := getIdNameData("/0/odata/Currency")
	ShowError("Can't get Currency List", errCurrency)
	if errCurrency != nil {
		err = errCurrency
	}

	usrTypePayment, errUsrTypePayment := getIdNameData("/0/odata/UsrTypePayment")
	ShowError("Can't get UsrTypePayment List", errUsrTypePayment)
	if errUsrTypePayment != nil {
		err = errUsrTypePayment
	}

	usrStatusPayment, errUsrStatusPayment := getIdNameData("/0/odata/UsrStatusPayment")
	ShowError("Can't get UsrStatusPayment List", errUsrStatusPayment)
	if errUsrStatusPayment != nil {
		err = errUsrStatusPayment
	}

	usrSourcePayment, errUsrSourcePayment := getIdNameData("/0/odata/UsrSourcePayment")
	ShowError("Can't get UsrSourcePayment List", errUsrSourcePayment)
	if errUsrSourcePayment != nil {
		err = errUsrSourcePayment
	}

	return ListOfPaymentSettings{
		Currency:         currency,
		UsrTypePayment:   usrTypePayment,
		UsrStatusPayment: usrStatusPayment,
		UsrSourcePayment: usrSourcePayment,
	}, err

}

func (data *ListOfPaymentSettings) GetIdByTypeAndName(t, name string) (result string) {

	switch t {
	case CURRENCY:
		result = getIdByName(data.Currency, name)
	case USR_STATUS_PAYMENT:
		result = getIdByName(data.UsrStatusPayment, name)
	case USR_TYPE_PAYMENT:
		result = getIdByName(data.UsrTypePayment, name)
	case USR_SOURCE_PAYMENT:
		result = getIdByName(data.UsrSourcePayment, name)
	}
	return result
}

func getIdByName(items []IdNameData, name string) string {
	for _, item := range items {
		if item.Name == name {
			return item.ID
		}
	}
	return ""
}

func getIdNameData(route string) ([]IdNameData, error) {
	var data IdNameDataResponse

	url := getPath(route)
	Cookies = GetConnectionCookies()

	request := gorequest.New()
	resp, _, errs := request.Get(url).AddCookies(Cookies).EndStruct(&data)
	if errs != nil {
		ShowError("Get Data By URL: "+url, errs[0])
		return []IdNameData{}, errs[0]
	}

	if resp.StatusCode == http.StatusOK {
		return data.Value, nil
	}

	return []IdNameData{}, nil
}

func getODataIDByName(url, name string) string {
	items, err := getIdNameData(url)

	if err != nil {
		return ""
	}

	return getIdByName(items, name)
}

func CheckError(message string, err error) {
	if err != nil {
		log.Fatalf(
			"[%s] %s: %s \n",
			time.Now().Format("Mon, 02 Jan 2006 15:04:05 "),
			message,
			err.Error(),
		)
	}
}

func ShowError(message string, err error) {
	if err != nil {
		log.Printf(
			"[%s] %s: %s \n",
			time.Now().Format("Mon, 02 Jan 2006 15:04:05 "),
			message,
			err.Error())
	}
}
