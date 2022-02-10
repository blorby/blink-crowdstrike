package actions

import (
"github.com/blinkops/blink-openapi-sdk/consts"
openapi_sdk "github.com/blinkops/blink-openapi-sdk/plugin"
"github.com/blinkops/blink-sdk/plugin"
"net/http"
)

func execRequest(ctx *plugin.ActionContext, request *http.Request, timeout int32) (*plugin.ExecuteActionResponse, error) {
	res := &plugin.ExecuteActionResponse{ErrorCode: consts.OK}
	response, err := openapi_sdk.ExecuteRequest(ctx, request, PluginName, Prefixes, HeaderAlias, timeout, nil)
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
