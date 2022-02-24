package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blinkops/blink-openapi-sdk/consts"
	openapi_sdk "github.com/blinkops/blink-openapi-sdk/plugin"
	"github.com/blinkops/blink-sdk/plugin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
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

func getGetInstalledDevicesParams(request *plugin.ExecuteActionRequest) ([]string, bool, error) {
	params, err := request.GetParameters()
	if err != nil {
		return nil, false, err
	}

	deviceSerials, onlyActiveDevices := params[deviceSerialsParam], params[onlyActiveDevicesParam]
	if deviceSerials == "" {
		return nil, false, errors.New("no device serials provided")
	}
	if onlyActiveDevices == "" {
		return nil, false, errors.New("input for 'return only active devices' not provided")
	}

	onlyActiveDevicesBool, err := strconv.ParseBool(onlyActiveDevices)
	if err != nil {
		return nil, false, errors.New("unable to convert 'return only active devices' to boolean")
	}

	deviceSerials = strings.ReplaceAll(deviceSerials, ", ", ",")
	deviceSerialsList := strings.Split(deviceSerials, ",")

	return deviceSerialsList, onlyActiveDevicesBool, nil
}

func getDeviceIdBySerial(ctx *plugin.ActionContext, requestUrl string, timeout int32, serial string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/devices/queries/devices/v1?filter=serial_number:'%s'", requestUrl, url.QueryEscape(serial)), nil)
	if err != nil {
		return "", err
	}

	resp, err := openapi_sdk.ExecuteRequest(ctx, req, PluginName, nil, nil, timeout, GetCrowdStrikeAccessToken)
	if err != nil {
		return "", err
	}

	var respJson QueryDevicesByFilterResponse
	err = json.Unmarshal(resp.Body, &respJson)
	if err != nil {
		return "", errors.New("failed to unmarshal response json")
	}

	if len(respJson.Resources) > 0 {
		return respJson.Resources[0], nil
	}

	return "", nil
}

func isDeviceActive(ctx *plugin.ActionContext, requestUrl string, timeout int32, deviceId string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/devices/entities/devices/v1?ids=%s", requestUrl, url.QueryEscape(deviceId)), nil)
	if err != nil {
		return false, err
	}

	resp, err := openapi_sdk.ExecuteRequest(ctx, req, PluginName, nil, nil, timeout, GetCrowdStrikeAccessToken)
	if err != nil {
		return false, err
	}

	var respJson GetDeviceResponse
	err = json.Unmarshal(resp.Body, &respJson)
	if err != nil {
		return false, errors.New("failed to unmarshal response json")
	}

	if len(respJson.Resources) > 0 && respJson.Resources[0].Status != "suppressed" {
		return true, nil
	}

	return false, nil
}

func performGetInstalledDevices(ctx *plugin.ActionContext, request *plugin.ExecuteActionRequest, deviceSerials []string, onlyActiveDevices bool) (chan string, error) {
	requestUrl, err := openapi_sdk.GetRequestUrl(ctx, PluginName)
	if err != nil {
		return nil, errors.New("no request url provided")
	}

	activeSerialsChan := make(chan string, len(deviceSerials))
	var errorToReturn error

	defer close(activeSerialsChan)

	var wg sync.WaitGroup
	for _, serial := range deviceSerials {
		serialVal := serial
		wg.Add(1)
		go func() {
			defer wg.Done()
			deviceId, err := getDeviceIdBySerial(ctx, requestUrl, request.Timeout, serialVal)
			if err != nil {
				errorToReturn = err
				return
			}

			// if device is installed
			if deviceId != "" {
				if onlyActiveDevices {
					// add serial to list only if device is active
					deviceActive, err := isDeviceActive(ctx, requestUrl, request.Timeout, deviceId)
					if err != nil {
						errorToReturn = err
						return
					}

					if deviceActive {
						activeSerialsChan <- serialVal
					}
				} else {
					activeSerialsChan <- serialVal
				}
			}
		}()
	}
	wg.Wait()

	return activeSerialsChan, err
}
