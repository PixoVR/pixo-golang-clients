package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/rs/zerolog/log"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"

	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
)

var _ Client = (*clientImpl)(nil)

// clientImpl is a struct for the graphql API that contains an abstract client
type clientImpl struct {
	*abstract.ServiceClient
}

// NewClient is a function that returns a clientImpl
func NewClient(config urlfinder.ClientConfig) Client {

	if config.Token == "" && config.APIKey == "" {
		config.APIKey = os.Getenv("PIXO_API_KEY")
	}

	abstractConfig := abstract.Config{
		ServiceConfig: newServiceConfig(config),
		Token:         config.Token,
		APIKey:        config.APIKey,
	}
	abstractClient := abstract.NewClient(abstractConfig)

	return &clientImpl{
		ServiceClient: abstractClient,
	}
}

// NewClientWithBasicAuth is a function that returns a clientImpl with basic auth performed
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (Client, error) {

	client := NewClient(config)

	if err := client.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login to the pixo platform")
		return nil, err
	}

	return client, nil
}

func (p *clientImpl) CheckAuth(ctx context.Context) (User, error) {
	res, err := p.Get(ctx, "auth/check")
	if err != nil {
		return User{}, err
	}

	var resPayload struct {
		Error string
		User  User
	}
	if err = json.NewDecoder(res.Body).Decode(&resPayload); err != nil {
		return User{}, err
	}

	if res.StatusCode != http.StatusOK {
		return User{}, errors.New(resPayload.Error)
	}

	return resPayload.User, nil
}

func (p *clientImpl) ActiveUserID() int {
	if !p.IsAuthenticated() {
		return 0
	}

	token := p.GetToken()

	rawToken, err := auth.ParseJWT(token)
	if err != nil {
		return 0
	}

	return rawToken.UserID
}

func (p *clientImpl) ActiveOrgID() int {
	if !p.IsAuthenticated() {
		return 0
	}

	token := p.GetToken()

	rawToken, err := auth.ParseJWT(token)
	if err != nil {
		return 0
	}

	return rawToken.OrgID
}

func newServiceConfig(config urlfinder.ClientConfig) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:     "v2",
		ServiceName: "primary-api",
		Lifecycle:   config.Lifecycle,
		Region:      config.Region,
		Namespace:   fmt.Sprintf("%s-apex", config.Lifecycle),
		Port:        8000,
	}
}

func (p *clientImpl) Exec(ctx context.Context, query string, v any, variables map[string]interface{}) error {
	req := GraphQLRequestPayload{
		Query:     query,
		Variables: variables,
	}
	reqBody, _ := json.Marshal(req)
	p.SetHeader("Content-Type", "application/json")
	res, err := p.Post(ctx, "query", reqBody)
	if err != nil {
		log.Error().
			Err(err).
			Str("query", query).
			Str("variables", fmt.Sprintf("%v", variables)).
			Msg("Request to platform API failed")
		return err
	}

	resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		var basicRes abstract.Response
		if err = json.Unmarshal(resBody, &basicRes); err == nil {
			return errors.New(basicRes.Error)
		}
		return fmt.Errorf("error: %d", res.StatusCode)
	}

	var gqlRes GraphQLResponse
	if err = json.Unmarshal(resBody, &gqlRes); err != nil {
		return err
	}

	if len(gqlRes.Errors) > 0 {
		return errors.New(gqlRes.Errors[0].Message)
	}

	return json.Unmarshal(gqlRes.Data, v)
}

func (p *clientImpl) ExecWithFile(ctx context.Context, query string, v any, variables map[string]interface{}, filePath, label string) error {
	req := GraphQLRequestPayload{
		Query:     query,
		Variables: variables,
	}
	if filePath == "" {
		return p.Exec(ctx, query, v, variables)
	}

	payload, writer, err := createMultipartWriter(req)

	if err = addFile(writer, filePath, label); err != nil {
		return err
	}

	p.ServiceClient.SetHeader("Content-Type", writer.FormDataContentType())

	res, err := p.Post(ctx, "query", payload.Bytes())
	if err != nil {
		return err
	}

	resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return errors.New(string(resBody))
	}

	var gqlRes GraphQLResponse
	if err = json.Unmarshal(resBody, &gqlRes); err != nil {
		return err
	}

	if len(gqlRes.Errors) > 0 {
		return errors.New(gqlRes.Errors[0].Message)
	}

	return json.Unmarshal(gqlRes.Data, v)
}

func createMultipartWriter(req GraphQLRequestPayload) (*bytes.Buffer, *multipart.Writer, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return nil, nil, err
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("operations", buf.String())

	return payload, writer, nil
}

func addFile(writer *multipart.Writer, filePath, label string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck

	mapData := map[string][]string{}
	mapData["0"] = []string{fmt.Sprintf(`variables.input.%s`, label)}
	jsonData, _ := json.Marshal(mapData)

	_ = writer.WriteField("map", string(jsonData))

	part, err := createFormFile(writer, "0", filepath.Base(filePath))
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	return nil
}

func createFormFile(w *multipart.Writer, fieldName, filename string) (io.Writer, error) {
	fileContentType := mime.TypeByExtension(filepath.Ext(filename))
	if fileContentType == "" {
		fileContentType = "application/octet-stream"
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	h.Set("Content-Type", fileContentType)
	return w.CreatePart(h)
}
