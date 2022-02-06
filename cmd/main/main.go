package main

import (
	"encoding/json"
	openapi_sdk "github.com/blinkops/blink-openapi-sdk/plugin"
	plugin_sdk "github.com/blinkops/blink-sdk"
	"github.com/blinkops/blink-sdk/plugin"
	"github.com/blinkops/blink-sdk/plugin/connections"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	PluginName = "EXAMPLE"
)

var prefixes = openapi_sdk.HeaderValuePrefixes{"AUTHORIZATION": "bearer "}
var headerAlias = openapi_sdk.HeaderAlias{"TOKEN": "AUTHORIZATION"}

func main() {
	connectionTypes := map[string]connections.Connection{
		PluginName: {
			Name:      PluginName,
			Fields:    nil,
			Reference: PluginName,
		},
	}

	metadata := openapi_sdk.PluginMetadata{
		Name:                PluginName,
		Provider:            PluginName,
		MaskFile:            "mask.yaml",
		OpenApiFile:         "EXAMPLE-openapi.yaml",
		Tags:                []string{PluginName},
		HeaderValuePrefixes: prefixes,
		HeaderAlias:         headerAlias,
	}

	checks := openapi_sdk.Callbacks{
		TestCredentialsFunc:
		func(ctx *plugin.ActionContext) (*plugin.CredentialsValidationResponse, error) {
			value, err := ValidateCredentials(ctx)

			return &plugin.CredentialsValidationResponse{
				AreCredentialsValid:   value,
				RawValidationResponse: err,
			}, nil
		},
		ValidateResponse:         Validate,
		GetTokenFromCrendentials: nil,
	}

	EXAMPLEPlugin, err := openapi_sdk.NewOpenApiPlugin(connectionTypes, metadata, checks)

	if err != nil {
		log.Error("Failed to create EXAMPLE integration: ", err)
		panic(err)
	}

	err = plugin_sdk.Start(EXAMPLEPlugin)

	if err != nil {
		log.Error("Failed to start EXAMPLE integration: ", err)
		panic(err)
	}
}

// ValidateCredentials test if the provided credentials are correct.
func ValidateCredentials(ctx *plugin.ActionContext) (bool, []byte) {

	req, _ := http.NewRequest(http.MethodGet, "https://api.EXAMPLE.com/auth.test", nil)

	response, err := openapi_sdk.ExecuteRequest(ctx, req, PluginName, prefixes, headerAlias, 30, nil)

	if err != nil {
		return false, []byte(err.Error())
	}

	return Validate(response)

}

func Validate(response openapi_sdk.Result) (bool, []byte) {
	var data map[string]interface{}

	err := json.Unmarshal(response.Body, &data)
	if err != nil {
		return false, response.Body
	}
	// validate the json

	// return false, and provide a message if the json contains an error.

	return true, nil
}
