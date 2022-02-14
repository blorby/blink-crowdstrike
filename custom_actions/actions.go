package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	openapi_sdk "github.com/blinkops/blink-openapi-sdk/plugin"
	customact "github.com/blinkops/blink-openapi-sdk/plugin/custom_actions"
	"github.com/blinkops/blink-sdk/plugin"
	"net/http"
)

const (
	PluginName = "crowdstrike"
	hostAgentIdParam = "Host Agent ID"
)

func GetCrowdStrikeCustomActions() customact.CustomActions {
	actions := map[string]customact.ActionHandler{
		"DeleteDevice": deleteDevice,
	}

	return customact.CustomActions{
		Actions:           actions,
		ActionsFolderPath: "custom_actions/actions",
	}
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
    "`+ id +`"
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

