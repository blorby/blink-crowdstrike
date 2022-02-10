package main

import (
	actions "blink-crowdstrike/custom_actions"
	openapi_sdk "github.com/blinkops/blink-openapi-sdk/plugin"
	plugin_sdk "github.com/blinkops/blink-sdk"
	"github.com/blinkops/blink-sdk/plugin"
	"github.com/blinkops/blink-sdk/plugin/connections"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	connectionTypes := map[string]connections.Connection{
		actions.PluginName: {
			Name:      actions.PluginName,
			Fields:    nil,
			Reference: actions.PluginName,
		},
	}

	metadata := openapi_sdk.PluginMetadata{
		Name:                actions.PluginName,
		Provider:            actions.PluginName,
		MaskFile:            "mask.yaml",
		OpenApiFile:         "crowdstrike-openapi.yaml",
		Tags:                []string{actions.PluginName},
		HeaderValuePrefixes: nil,
		HeaderAlias:         nil,
	}

	checks := openapi_sdk.Callbacks{
		TestCredentialsFunc: func(ctx *plugin.ActionContext) (*plugin.CredentialsValidationResponse, error) {
			value, err := ValidateCredentials(ctx)

			return &plugin.CredentialsValidationResponse{
				AreCredentialsValid:   value,
				RawValidationResponse: err,
			}, nil
		},
		SetCustomAuthHeaders: actions.GetCrowdStrikeAccessToken,
		CustomActions: actions.GetCrowdStrikeCustomActions(),
	}

	crowdstrikePlugin, err := openapi_sdk.NewOpenApiPlugin(connectionTypes, metadata, checks)

	if err != nil {
		log.Error("Failed to create crowdstrike integration: ", err)
		panic(err)
	}

	err = plugin_sdk.Start(crowdstrikePlugin)

	if err != nil {
		log.Error("Failed to start crowdstrike integration: ", err)
		panic(err)
	}
}

// ValidateCredentials test if the provided credentials are correct.
func ValidateCredentials(ctx *plugin.ActionContext) (bool, []byte) {

	req, _ := http.NewRequest(http.MethodGet, "https://api.crowdstrike.com/auth.test", nil)

	connection, err := openapi_sdk.GetCredentials(ctx, actions.PluginName)

	if err != nil {
		return false, []byte("unable to get credentials")
	}

	err = actions.GetCrowdStrikeAccessToken(connection, req)

	if err != nil {
		return false, []byte("invalid credentials")
	}

	return true, nil

}
