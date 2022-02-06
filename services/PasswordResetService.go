package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/rand"
	"net/url"
	"path"
	"strconv"
	"time"
)

type resetPayload struct {
	Signature []byte    `json:"signature"`
	Data      resetData `json:"data"`
}
type resetData struct {
	ID       int           `json:"id"`
	Code     int           `json:"code"`
	Created  int64         `json:"created"`
	Validity time.Duration `json:"validity"` // validity is in minutes
}

const validity = time.Hour / 2
const codeKey = "c"
const payloadKey = "p"

// PasswordResetMailer is tha mailer  interface required by password reset
// service
type PasswordResetMailer interface {
	SendResetCode(name, email string, code int) error
	SendResetLink(name, email, link string) error
}

// PasswordResetService ...
type PasswordResetService struct {
	appSecret []byte
	verifPath string
	mailer    PasswordResetMailer
}

// NewPasswordResetService construct the password reset service
func NewPasswordResetService(verifPath, secret string,
	mailer PasswordResetMailer) (PasswordResetService, error) {
	return PasswordResetService{[]byte(secret), verifPath, mailer}, nil
}

// VerifyAndGetUserID verifies the payload and the supplied code,
// returning the userId if the payload and code are valid
func (prc PasswordResetService) VerifyAndGetUserID(code int,
	payload string) (int, error) {
	jsonDTO, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return 0, err
	}
	var data resetPayload
	err = json.Unmarshal(jsonDTO, &data)
	if err != nil {
		return 0, err
	}
	if prc.verifyData(code, data) == false {
		return 0, errors.New("Invalid or expired data")
	}
	return data.Data.ID, nil
}

func (prc PasswordResetService) verifyData(code int, data resetPayload) bool {
	signature := data.Signature
	content := data.Data
	content.Code = code
	expected := prc.signData(content)
	equal := hmac.Equal(signature, expected)
	if equal == false {
		return false
	}
	return prc.checkDateValid(content)
}

func (prc PasswordResetService) checkDateValid(data resetData) bool {
	createdAt := time.Unix(data.Created, 0)
	difference := createdAt.Sub(time.Now())
	return difference <= data.Validity
}

// GetURLClaims return the code and
func (prc PasswordResetService) GetURLClaims(requestURL string) (int, string) {
	reqURL, _ := url.ParseRequestURI(requestURL)
	payload := reqURL.Query().Get(payloadKey)
	code := reqURL.Query().Get(codeKey)
	intCode, _ := strconv.Atoi(code)
	return intCode, payload
}

// SendPasswordResetLink ...
func (prc PasswordResetService) SendPasswordResetLink(name string,
	id int, email, baseURL string) error {
	code, payload := prc.generatePayload(id)
	verifURL := path.Join(baseURL, prc.verifPath)
	base, _ := url.Parse(verifURL)
	params := url.Values{}
	params.Add(payloadKey, payload)
	params.Add(codeKey, strconv.Itoa(code))
	base.RawQuery = params.Encode()
	link := base.String()
	return prc.mailer.SendResetLink(name, email, link)
}

// SendPasswordResetCode ...
func (prc PasswordResetService) SendPasswordResetCode(name string,
	id int, email string) (string, error) {
	code, payload := prc.generatePayload(id)
	err := prc.mailer.SendResetCode(name, email, code)
	return payload, err
}

func (prc PasswordResetService) generatePayload(id int) (int, string) {
	data := prc.generateData(id)
	signature := prc.signData(data)
	code := data.Code
	data.Code = 0 // remove the code from the payload  struct
	dto := resetPayload{signature, data}
	dtoJSON, _ := json.Marshal(dto)
	payload := base64.StdEncoding.EncodeToString(dtoJSON)
	return code, payload
}

func (prc PasswordResetService) generateData(id int) resetData {
	created := time.Now().Unix()
	code := generateCode()
	return resetData{
		Created:  created,
		ID:       id,
		Validity: validity,
		Code:     code,
	}
}

func (prc PasswordResetService) signData(d resetData) []byte {
	jsonString, _ := json.Marshal(d)
	hm := hmac.New(sha256.New, prc.appSecret)
	hm.Write(jsonString)
	return hm.Sum(nil)
}

func generateCode() int {
	rand.Seed(time.Now().UnixNano())
	return 100000 + rand.Intn(999999-100000)
}
