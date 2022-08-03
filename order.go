package odata4

import (
	"net/http"

	"github.com/parnurzeal/gorequest"
)

type (
	NewOrderData struct {
		UsrName              string  `json:"UsrName,omitempty"`
		UsrClientId          string  `json:"UsrClientId,omitempty"`
		UsrSum               float32 `json:"UsrSum"`
		UsrTypePaymentId     string  `json:"UsrTypePaymentId"`
		UsrDatePayment       string  `json:"UsrDatePayment"`
		UsrStatusPaymentId   string  `json:"UsrStatusPaymentId"`
		UsrIDPayment         string  `json:"UsrIDPayment"`
		UsrReferal           string  `json:"UsrReferal,omitempty"`
		UsrBank              string  `json:"UsrBank,omitempty"`
		UsrUnsubscribed      string  `json:"UsrUnsubscribed,omitempty"`
		UsrReasonUnsubscribe string  `json:"UsrReasonUnsubscribe,omitempty"`
		UsrRegion            string  `json:"UsrRegion,omitempty"`
		UsrCity              string  `json:"UsrCity,omitempty"`
		UsrSourcePaymentId   string  `json:"UsrSourcePaymentId"`
		UsrCurrencyId        string  `json:"UsrCurrencyId"`
	}

	OrderData struct {
		Id                   string  `json:"Id"`
		CreatedOn            string  `json:"CreatedOn"`
		CreatedById          string  `json:"CreatedById"`
		ModifiedOn           string  `json:"ModifiedOn"`
		ModifiedById         string  `json:"ModifiedById"`
		UsrName              string  `json:"UsrName"`
		UsrNotes             string  `json:"UsrNotes"`
		UsrIDClient          string  `json:"UsrIDClient"`
		UsrClientId          string  `json:"UsrClientId"`
		UsrSum               float32 `json:"UsrSum"`
		UsrCurrencyId        string  `json:"UsrCurrencyId"`
		UsrTypePaymentId     string  `json:"UsrTypePaymentId"`
		UsrDatePayment       string  `json:"UsrDatePayment"`
		UsrStatusPaymentId   string  `json:"UsrStatusPaymentId"`
		UsrIDPayment         string  `json:"UsrIDPayment"`
		UsrReferal           string  `json:"UsrReferal"`
		UsrUnsubscribed      string  `json:"UsrUnsubscribed"`
		UsrReasonUnsubscribe string  `json:"UsrReasonUnsubscribe"`
		UsrRegion            string  `json:"UsrRegion"`
		UsrCity              string  `json:"UsrCity"`
		UsrSourcePaymentId   string  `json:"UsrSourcePaymentId"`
	}

	UpdateOrderDataParams struct {
		UsrStatusPaymentId   string `json:"UsrStatusPaymentId"`
		UsrUnsubscribed      string `json:"UsrUnsubscribed,omitempty"`
		UsrReasonUnsubscribe string `json:"UsrReasonUnsubscribe,omitempty"`
		UsrBank              string `json:"UsrBank,omitempty"`
	}

	ListOfPaymentSettings struct {
		Currency         []IdNameData `json:"currency"`
		UsrTypePayment   []IdNameData `json:"usrTypePayment"`
		UsrStatusPayment []IdNameData `json:"usrStatusPayment"`
		UsrSourcePayment []IdNameData `json"usrSourcePayment"`
	}
)

const (
	PAY       = "Одноразовий донат"
	PAYDANATE = "Регулярний донат"
	SUBSCRIBE = "Оформлена підписка"
	REVERSE   = "Повернення платежу"
)

const (
	STATUS_NEW          = "Новий платіж"
	STATUS_SUCCESS      = "Успішний платіж"
	STATUS_SUBSCRIBED   = "Підписку оформлено"
	STATUS_UNSUBSCRIBED = "Підписку скасовано"
	STATUS_FAILURE      = "Невдалий платіж"
	STATUS_PROCESSING   = "Платіж в обробці"
	STATUS_REVERSED     = "Повернений платіж"
)

func (o *NewOrderData) SetCurrecyIDByName(name string) {
	o.UsrCurrencyId = getODataIDByName("/0/odata/Currency", name)
}

func (o *NewOrderData) SetTypePaymentIDbyName(name string) {
	o.UsrTypePaymentId = getODataIDByName("/0/odata/UsrTypePayment", name)
}

func (o *NewOrderData) SetStatusPaymentIDByName(name string) {
	o.UsrStatusPaymentId = getODataIDByName("/0/odata/UsrStatusPayment", name)
}

func (o *NewOrderData) SetSourcePayment(name string) {
	o.UsrSourcePaymentId = getODataIDByName("/0/odata/UsrSourcePayment", name)
}

func CreateNewOrderAndGetID(o NewOrderData) string {
	order, err := CreateNewOrder(o)
	ShowError("Create New Order in CRM", err)

	if err == nil {
		return order.Id
	}

	return ""
}

func CreateNewOrder(o NewOrderData) (OrderData, error) {
	var order OrderData

	url := getPath("/0/odata/UsrDonate")
	Cookies = GetConnectionCookies()
	hParam := getParamsFromCookies(COOKIE_BPMCSRF)

	request := gorequest.New()
	resp, _, errs := request.Post(url).
		Set(COOKIE_BPMCSRF, hParam).
		AddCookies(Cookies).
		Send(o).
		EndStruct(&order)
	if errs != nil {
		ShowError("Create New Contact", errs[0])
	}

	if resp.StatusCode == http.StatusCreated {
		return order, nil
	}

	return order, errs[0]
}

func UpdateOrder(params UpdateOrderDataParams, orderID string) error {
	url := getPath("/0/odata/UsrDonate(" + orderID + ")")
	Cookies = GetConnectionCookies()
	hParam := getParamsFromCookies(COOKIE_BPMCSRF)

	request := gorequest.New()
	_, _, errs := request.Patch(url).
		Set(COOKIE_BPMCSRF, hParam).
		AddCookies(Cookies).
		Send(params).
		End()
	if errs != nil {
		return errs[0]
	}

	return nil
}
