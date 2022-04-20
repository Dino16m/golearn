package services

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strconv"
	"testing"
	"time"

	"golearn-api-template/services/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PasswordResetServiceTest struct {
	suite.Suite
	service PasswordResetService
	mailer  *mocks.PasswordResetMailer
}

const appURL = "http://localhost"
const verifPath = "/reset-password"

var name = "dummy"
var email = "dummy@test.com"
var id = 10

func (p *PasswordResetServiceTest) SetupTest() {
	mailer := new(mocks.PasswordResetMailer)
	p.mailer = mailer
	p.service, _ = NewPasswordResetService(
		verifPath, "secret", mailer,
	)
}

func (p *PasswordResetServiceTest) TestSendPasswordResetLink() {
	p.mailer.On("SendResetLink", name, email, mock.MatchedBy(
		func(link string) bool {
			_, err := url.Parse(link)
			return err == nil
		}),
	).Return(nil)
	p.service.SendPasswordResetLink(name, id, email, appURL)
}
func (p *PasswordResetServiceTest) TestSendPasswordResetCode() {
	p.mailer.On(
		"SendResetCode", name, email, mock.AnythingOfType("int"),
	).Return(nil)
	p.service.SendPasswordResetCode(name, id, email)
}

func (p *PasswordResetServiceTest) TestGetURLClaims() {
	resetLink := p.getLink()
	base, _ := url.Parse(resetLink)
	queryCode := base.Query().Get("c")
	expectedCode, _ := strconv.Atoi(queryCode)
	expectedPayload := base.Query().Get("p")

	code, payload := p.service.GetURLClaims(resetLink)

	p.Equal(expectedCode, code)
	p.Equal(expectedPayload, payload)
}

func (p *PasswordResetServiceTest) getLink() string {
	var resetLink string
	p.mailer.On("SendResetLink", name, email, mock.MatchedBy(
		func(link string) bool {
			resetLink = link
			_, err := url.Parse(link)
			return err == nil
		}),
	).Return(nil)
	p.service.SendPasswordResetLink(name, id, email, appURL)
	return resetLink
}
func (p *PasswordResetServiceTest) TestVerifyDataRejectedWithWrongCode() {
	resetLink := p.getLink()
	code, payload := p.service.GetURLClaims(resetLink)
	modifiedCode := 666666
	p.NotEqual(code, modifiedCode)
	id, err := p.service.VerifyAndGetUserID(modifiedCode, payload)
	p.Equal(id, 0)
	p.NotNil(err)
}
func (p *PasswordResetServiceTest) TestVerifyDataRejectedWithWrongPayload() {
	resetLink := p.getLink()
	code, _ := p.service.GetURLClaims(resetLink)
	modifiedPayload := "helloworld"
	id, err := p.service.VerifyAndGetUserID(code, modifiedPayload)
	p.Equal(id, 0)
	p.NotNil(err)
}

func (p *PasswordResetServiceTest) TestVerifyDataRejectedWithExpiredPayload() {
	resetLink := p.getLink()
	code, payload := p.service.GetURLClaims(resetLink)
	decodedPayload, _ := base64.StdEncoding.DecodeString(payload)
	var data resetPayload
	json.Unmarshal(decodedPayload, &data)
	created, valid := data.Data.Created, data.Data.Validity
	oneMinute := time.Minute
	later := time.Unix(created, 0).Add(valid + oneMinute).Unix()
	data.Data.Created = later

	dto, _ := json.Marshal(data)
	modifiedPayload := base64.StdEncoding.EncodeToString(dto)

	id, err := p.service.VerifyAndGetUserID(code, modifiedPayload)
	p.Equal(id, 0)
	p.NotNil(err)
}
func (p *PasswordResetServiceTest) TestVerifyDataAcceptedWithValidPayload() {
	resetLink := p.getLink()
	code, payload := p.service.GetURLClaims(resetLink)
	userID, err := p.service.VerifyAndGetUserID(code, payload)
	p.Nil(err)
	p.Equal(id, userID)
}

func TestPasswordResetService(t *testing.T) {
	suite.Run(t, new(PasswordResetServiceTest))
}
