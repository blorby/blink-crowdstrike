package actions

import (
	"encoding/json"
	"errors"
	"github.com/blinkops/blink-openapi-sdk/consts"
	openapi_sdk "github.com/blinkops/blink-openapi-sdk/plugin"
	"github.com/blinkops/blink-sdk/plugin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	requestUrlParam   = "REQUEST_URL"
	clientIdParam     = "Client ID"
	clientSecretParam = "Client Secret"
)

func GetCrowdStrikeAccessToken(connection map[string]string, request *http.Request) error {
	requestUrl := connection[requestUrlParam]

	queryParams := url.Values{
		"client_id":     {connection[clientIdParam]},
		"client_secret": {connection[clientSecretParam]},
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, requestUrl+"/oauth2/token", strings.NewReader(queryParams.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return errors.New("invalid credentials")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var responseBody struct {
		AccessToken string `json:"access_token"`
	}

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return err
	}
	request.Header.Set("AUTHORIZATION", "Bearer " + responseBody.AccessToken)
	return nil
}

func execRequest(ctx *plugin.ActionContext, request *http.Request, timeout int32) (*plugin.ExecuteActionResponse, error) {
	res := &plugin.ExecuteActionResponse{ErrorCode: consts.OK}
	response, err := openapi_sdk.ExecuteRequest(ctx, request, PluginName, nil, nil, timeout, GetCrowdStrikeAccessToken)
	if err != nil {
		return nil, err
	}
	res.Result = response.Body
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return res, nil
	}
	res.ErrorCode = consts.Error

	return res, nil
}
