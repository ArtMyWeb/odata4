package odata4

import (
	"log"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

type (
	ContactData struct {
		ID                         string `json:"Id"`
		CreatedOn                  string `json:"CreatedOn,omitempty"`
		CreatedById                string `json:"CreatedById,omitempty"`
		ModifiedOn                 string `json:"ModifiedOn,omitempty"`
		ModifiedById               string `json:"ModifiedById,omitempty"`
		CommunicationTypeId        string `json:"CommunicationTypeId,omitempty"`
		Number                     string `json:"Number,omitempty"`
		ContactId                  string `json:"ContactId,omitempty"`
		Position                   int    `json:"Position,omitempty"`
		SocialMediaId              string `json:"SocialMediaId,omitempty"`
		SearchNumber               string `json:"SearchNumber,omitempty"`
		ProcessListeners           int    `json:"ProcessListeners,omitempty"`
		IsCreatedBySynchronization bool   `json:"IsCreatedBySynchronization,omitempty"`
		ExternalCommunicationType  string `json:"ExternalCommunicationType,omitempty"`
		Name                       string `json:"Name"`
	}

	NewContact struct {
		Name             string `json:"Name"`
		MobilePhone      string `json:"MobilePhone,omitempty"`
		Email            string `json:"Email,omitempty"`
		UsrTypeUserId    string `json:"UsrTypeUserId"`
		UsriduserGO      string `json:"UsriduserGO,omitempty"`
		UsrIdUserBF      string `json:"UsrIdUserBF,omitempty"`
		UsrIdUserFLI     string `json:"UsrIdUserFLI,omitempty"`
		UsrIdUserDodatok string `json:"UsrIdUserDodatok,omitempty"`
	}

	ConstactResponse struct {
		Value []ContactData `json:"value"`
	}

	UpdateContactParams struct {
		UsriduserGO  string `json:"UsriduserGO,omitempty"`
		UsrIdUserBF  string `json:"UsrIdUserBF,omitempty"`
		UsrIdUserFLI string `json:"UsrIdUserFLI,omitempty"`
	}
)

const (
	CLIENTTYPE_CITIZEN   = "Громадянин"
	CLIENTTYPE_SUPPORTER = "Прихильник"
)

func (c *NewContact) SetContactType(cType string) {
	c.UsrTypeUserId = getODataIDByName("/0/odata/UsrTypeUser", cType)
}

func createFilterByPhoneAndNumberForContact(phone, email string) string {
	filter := ""

	if email != "" && phone == "" {
		filter = "?$filter=contains(Number, '" + email + "')"
	} else if email == "" && phone != "" {
		filter = "?$filter=contains(Number, '" + phone + "')"
	} else if email != "" && phone != "" {
		filter = "?$filter=contains(Number, '" + email + "') or contains(Number, '" + phone + "')"
	}

	return filter
}

func GetContactIdByNameAndPhone(phone, email string) string {

	if phone == "" && email == "" {
		return ""
	}

	var contacts ConstactResponse

	filter := createFilterByPhoneAndNumberForContact(phone, email)
	url := getPath("/0/odata/ContactCommunication" + filter)
	Cookies = GetConnectionCookies(true)

	request := gorequest.New()
	resp, _, errs := request.Get(url).AddCookies(Cookies).EndStruct(&contacts)

	if errs != nil {
		log.Println("Cet Contact Data", errs[0])
		return ""
	}
	if resp.StatusCode == http.StatusOK && len(contacts.Value) > 0 {
		return contacts.Value[0].ContactId
	}

	return ""
}

func CreateNewContactAndGetID(c NewContact) (string, error) {
	var contact ContactData

	url := getPath("/0/odata/Contact")
	Cookies = GetConnectionCookies(false)
	hParam := getParamsFromCookies(COOKIE_BPMCSRF)

	request := gorequest.New()
	_, _, errs := request.Post(url).
		Set(COOKIE_BPMCSRF, hParam).
		AddCookies(Cookies).
		Send(c).
		EndStruct(&contact)

	if errs != nil {
		ShowError("Create New Contact", errs[0])
		return "", errs[0]

	}

	return contact.ID, nil
}

func UpdateContact(params UpdateContactParams, userID string) error {
	url := getPath("/0/odata/Contact(" + userID + ")")
	Cookies = GetConnectionCookies(false)
	hParams := getParamsFromCookies(COOKIE_BPMCSRF)

	request := gorequest.New()
	_, _, errs := request.Patch(url).
		Set(COOKIE_BPMCSRF, hParams).
		AddCookies(Cookies).
		Send(params).
		End()

	if errs != nil {
		return errs[0]
	}

	return nil
}
