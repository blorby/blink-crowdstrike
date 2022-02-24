package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/blinkops/blink-openapi-sdk/consts"
	openapi_sdk "github.com/blinkops/blink-openapi-sdk/plugin"
	customact "github.com/blinkops/blink-openapi-sdk/plugin/custom_actions"
	"github.com/blinkops/blink-sdk/plugin"
	"net/http"
	"strings"
	"sync"
)

const (
	PluginName         = "crowdstrike"
	hostAgentIdParam   = "Host Agent ID"
	deviceSerialsParam = "Device Serials"
)

type QueryDevicesByFilterResponse struct {
	Meta struct {
		Pagination struct {
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

func GetCrowdStrikeCustomActions() customact.CustomActions {
	actions := map[string]customact.ActionHandler{
		"GetInstalledDevices": getInstalledDevices,
		"DeleteDevice":        deleteDevice,
	}

	return customact.CustomActions{
		Actions:           actions,
		ActionsFolderPath: "custom_actions/actions",
	}
}

func getInstalledDevices(ctx *plugin.ActionContext, request *plugin.ExecuteActionRequest) (*plugin.ExecuteActionResponse, error) {
	requestUrl, err := openapi_sdk.GetRequestUrl(ctx, PluginName)
	if err != nil {
		return nil, errors.New("no request url provided")
	}

	params, err := request.GetParameters()
	if err != nil {
		return nil, err
	}

	deviceSerials := params[deviceSerialsParam]
	if deviceSerials == "" {
		return nil, errors.New("no device serials provided")
	}

	deviceSerials = strings.ReplaceAll(deviceSerials, ", ", ",")
	deviceSerialsList := strings.Split(deviceSerials, ",")

	installedSerialsChan := make(chan string, len(deviceSerialsList))
	var errorToReturn error

	var wg sync.WaitGroup
	for _, serial := range deviceSerialsList {
		serialVal := serial
		wg.Add(1)
		go func() {
			defer wg.Done()
			deviceInstalled, err := isDeviceInstalled(ctx, requestUrl, request.Timeout, serialVal)
			if err != nil {
				errorToReturn = err
				return
			}

			if deviceInstalled {
				installedSerialsChan <- serialVal
			}
		}()
	}
	wg.Wait()

	close(installedSerialsChan)

	if err != nil {
		return nil, err
	}

	var installedSerials []string
	for len(installedSerialsChan) > 0 {
		installedSerials = append(installedSerials, <-installedSerialsChan)
	}

	installedSerialsStr, err := json.Marshal(installedSerials)
	if err != nil {
		return nil, errors.New("failed to marshal installed serials json")
	}

	return &plugin.ExecuteActionResponse{ErrorCode: consts.OK, Result: installedSerialsStr}, nil
}

func deleteDevice(ctx *plugin.ActionContext, request *plugin.ExecuteActionRequest) (*plugin.ExecuteActionResponse, error) {
	requestUrl, err := openapi_sdk.GetRequestUrl(ctx, PluginName)
	if err != nil {
		return nil, errors.New("no request url provided")
	}

	params, err := request.GetParameters()
	if err != nil {
		return nil, err
	}

	id := params[hostAgentIdParam]

	reqBody := json.RawMessage(`{
  "ids": [
    "` + id + `"
  ]
}`)
	marshalledReqBody, err := json.Marshal(&reqBody)
	if err != nil {
		return nil, errors.New("failed to marshal json")
	}

	url := requestUrl + "/devices/entities/devices-actions/v2?action_name=hide_host"
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshalledReqBody))
	req.Header.Set("Content-Type", "application/json")

	return execRequest(ctx, req, request.Timeout)
}
