package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"strconv"
)

type data struct {
	ID        int    `json:"id"`
	Signature []byte `json:"signature"`
}

const paramKey = "u"

// VerificationMailer ...
type VerificationMailer interface {
	SendEmailVerification(name string, email string, link string) error
}

// EmailVerifService generates links to be sent to users as well as verifies
// users when they return with the link
type EmailVerifService struct {
	appSecret []byte
	mailer    VerificationMailer
	verifPath string
}

// NewEmailVerifService ...
func NewEmailVerifService(verifPath, secret string,
	mailer VerificationMailer) (EmailVerifService, error) {
	return EmailVerifService{[]byte(secret), mailer, verifPath}, nil
}

// VerifyAndGetUserID checks that the payload contained in a link belong to a
// particular user and that the link was issued by the app
// it returns an error if the url is malformed or if the payload is invalid
func (ev EmailVerifService) VerifyAndGetUserID(requestURL string) (int, error) {
	reqURL, err := url.ParseRequestURI(requestURL)
	if err != nil {
		return 0, err
	}
	payload := reqURL.Query().Get(paramKey)
	id, match := ev.verifyAndGetID(payload)
	if match == false {
		return id, errors.New("Invalid verification payload")
	}
	return id, nil
}

// SendVerificationMail ...
func (ev EmailVerifService) SendVerificationMail(
	name string, id int, email, base string) error {
	verificationURL := path.Join(base, ev.verifPath)
	link := ev.generateVerificationLink(id, verificationURL)
	return ev.mailer.SendEmailVerification(name, email, link)
}

func (ev EmailVerifService) generateVerificationLink(id int, baseURL string) string {
	base, _ := url.Parse(baseURL)
	params := url.Values{}
	token := ev.generateToken(id)
	params.Add(paramKey, token)
	base.RawQuery = params.Encode()
	return base.String()
}

func (ev EmailVerifService) generateToken(id int) string {
	signedID := ev.signID(id)

	dataDTO := data{ID: id, Signature: signedID}
	jsonString, _ := json.Marshal(dataDTO)
	return base64.StdEncoding.EncodeToString(jsonString)
}

func (ev EmailVerifService) signID(id int) []byte {
	hm := hmac.New(sha256.New, ev.appSecret)
	hm.Write([]byte(strconv.Itoa(id)))
	return hm.Sum(nil)
}

func (ev EmailVerifService) verifyAndGetID(payload string) (int, bool) {
	var dataDTO data
	decodedPayload, _ := base64.StdEncoding.DecodeString(payload)
	if err := json.Unmarshal([]byte(decodedPayload), &dataDTO); err != nil {
		return 0, false
	}

	expected := ev.signID(dataDTO.ID)
	provided := dataDTO.Signature
	return dataDTO.ID, hmac.Equal(expected, provided)
}
