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
)

const (
	PluginName             = "crowdstrike"
	hostAgentIdParam       = "Host Agent ID"
	deviceSerialsParam     = "Device Serials"
	onlyActiveDevicesParam = "Return only active devices"
)

type QueryDevicesByFilterResponse struct {
	Resources []string `json:"resources"`
}

type GetDeviceResponse struct {
	Resources []struct {
		Status string `json:"detection_suppression_status,omitempty"`
	} `json:"resources"`
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
	deviceSerials, onlyActiveDevices, err := getGetInstalledDevicesParams(request)

	activeSerialsChan, err := performGetInstalledDevices(ctx, request, deviceSerials, onlyActiveDevices)
	if err != nil {
		return nil, err
	}

	var activeSerials []string
	for len(activeSerialsChan) > 0 {
		activeSerials = append(activeSerials, <-activeSerialsChan)
	}

	activeSerialsStr, err := json.Marshal(activeSerials)
	if err != nil {
		return nil, errors.New("failed to marshal active serials json")
	}

	return &plugin.ExecuteActionResponse{ErrorCode: consts.OK, Result: activeSerialsStr}, nil
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
