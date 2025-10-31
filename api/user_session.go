// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	policies "github.com/openstor/console/api/policy"
	"github.com/openstor/madmin-go/v4"

	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/openstor/pkg/v3/policy/condition"

	minioIAMPolicy "github.com/openstor/pkg/v3/policy"

	"github.com/go-openapi/runtime/middleware"
	"github.com/openstor/console/api/operations"
	authApi "github.com/openstor/console/api/operations/auth"
	"github.com/openstor/console/models"
	"github.com/openstor/console/pkg/auth/idp/oauth2"
	"github.com/openstor/console/pkg/auth/ldap"
)

type Conditions struct {
	S3Prefix []string `json:"s3:prefix"`
}

func registerSessionHandlers(api *operations.ConsoleAPI) {
	// session check
	api.AuthSessionCheckHandler = authApi.SessionCheckHandlerFunc(func(params authApi.SessionCheckParams, session *models.Principal) middleware.Responder {
		sessionResp, err := getSessionResponse(params.HTTPRequest.Context(), session)
		if err != nil {
			return authApi.NewSessionCheckDefault(err.Code).WithPayload(err.APIError)
		}
		return authApi.NewSessionCheckOK().WithPayload(sessionResp)
	})
}

func getClaimsFromToken(sessionToken string) (map[string]interface{}, error) {
	jp := jwtgo.NewParser()
	var claims jwtgo.MapClaims
	_, _, err := jp.ParseUnverified(sessionToken, &claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// getSessionResponse parse the token of the current session and returns a list of allowed actions to render in the UI
func getSessionResponse(ctx context.Context, session *models.Principal) (*models.SessionResponse, *CodedAPIError) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// serialize output
	if session == nil {
		return nil, ErrorWithContext(ctx, ErrInvalidSession)
	}
	tokenClaims, _ := getClaimsFromToken(session.STSSessionToken)
	// initialize admin client
	mAdminClient, err := NewMinioAdminClient(ctx, &models.Principal{
		STSAccessKeyID:     session.STSAccessKeyID,
		STSSecretAccessKey: session.STSSecretAccessKey,
		STSSessionToken:    session.STSSessionToken,
	})
	if err != nil {
		return nil, ErrorWithContext(ctx, err, ErrInvalidSession)
	}
	userAdminClient := AdminClient{Client: mAdminClient}
	// Obtain the current policy assigned to this user
	// necessary for generating the list of allowed endpoints
	accountInfo, err := getAccountInfo(ctx, userAdminClient)
	if err != nil {
		return nil, ErrorWithContext(ctx, err, ErrInvalidSession)
	}
	erasure := accountInfo.Server.Type == madmin.Erasure
	rawPolicy := policies.ReplacePolicyVariables(tokenClaims, accountInfo)
	policy, err := minioIAMPolicy.ParseConfig(bytes.NewReader(rawPolicy))
	if err != nil {
		return nil, ErrorWithContext(ctx, err, ErrInvalidSession)
	}
	currTime := time.Now().UTC()

	customStyles := session.CustomStyleOb
	// This actions will be global, meaning has to be attached to all resources
	conditionValues := map[string][]string{
		condition.AWSUsername.Name(): {session.AccountAccessKey},
		// All calls to MinIO from console use temporary credentials.
		condition.AWSPrincipalType.Name():   {"AssumeRole"},
		condition.AWSSecureTransport.Name(): {strconv.FormatBool(getMinIOEndpointIsSecure())},
		condition.AWSCurrentTime.Name():     {currTime.Format(time.RFC3339)},
		condition.AWSEpochTime.Name():       {strconv.FormatInt(currTime.Unix(), 10)},

		// All calls from console are signature v4.
		condition.S3SignatureVersion.Name(): {"AWS4-HMAC-SHA256"},
		// All calls from console use header-based authentication
		condition.S3AuthType.Name(): {"REST-HEADER"},
		// This is usually empty, may be set some times (rare).
		condition.S3LocationConstraint.Name(): {GetMinIORegion()},
	}

	claims, err := getClaimsFromToken(session.STSSessionToken)
	if err != nil {
		return nil, ErrorWithContext(ctx, err, ErrInvalidSession)
	}

	// Support all LDAP, JWT variables
	for k, v := range claims {
		vstr, ok := v.(string)
		if !ok {
			// skip all non-strings
			continue
		}
		// store all claims from sessionToken
		conditionValues[k] = []string{vstr}
	}

	defaultActions := policy.IsAllowedActions("", "", conditionValues)

	// Allow Create Access Key when admin:CreateServiceAccount is provided with a condition
	for _, statement := range policy.Statements {
		if statement.Effect == "Deny" && len(statement.Conditions) > 0 &&
			statement.Actions.Contains(minioIAMPolicy.CreateServiceAccountAdminAction) {
			defaultActions.Add(minioIAMPolicy.Action(minioIAMPolicy.CreateServiceAccountAdminAction))
		}
	}

	permissions := map[string]minioIAMPolicy.ActionSet{
		ConsoleResourceName: defaultActions,
	}
	deniedActions := map[string]minioIAMPolicy.ActionSet{}

	var allowResources []*models.PermissionResource

	for _, statement := range policy.Statements {
		for _, resource := range statement.Resources.ToSlice() {
			resourceName := resource.String()
			statementActions := statement.Actions.ToSlice()
			var prefixes []string

			if statement.Effect == "Allow" {
				// check if val are denied before adding them to the map
				var allowedActions []minioIAMPolicy.Action
				if dActions, ok := deniedActions[resourceName]; ok {
					for _, action := range statementActions {
						if len(dActions.Intersection(minioIAMPolicy.NewActionSet(action))) == 0 {
							// It's ok to allow this action
							allowedActions = append(allowedActions, action)
						}
					}
				} else {
					allowedActions = statementActions
				}

				// Add validated actions
				if resourceActions, ok := permissions[resourceName]; ok {
					mergedActions := append(resourceActions.ToSlice(), allowedActions...)
					permissions[resourceName] = minioIAMPolicy.NewActionSet(mergedActions...)
				} else {
					mergedActions := append(defaultActions.ToSlice(), allowedActions...)
					permissions[resourceName] = minioIAMPolicy.NewActionSet(mergedActions...)
				}

				// Allow Permissions request
				conditions, err := statement.Conditions.MarshalJSON()
				if err != nil {
					return nil, ErrorWithContext(ctx, err)
				}

				var wrapper map[string]Conditions

				if err := json.Unmarshal(conditions, &wrapper); err != nil {
					return nil, ErrorWithContext(ctx, err)
				}

				for condition, elements := range wrapper {
					prefixes = elements.S3Prefix

					resourceElement := models.PermissionResource{
						Resource:          resourceName,
						Prefixes:          prefixes,
						ConditionOperator: condition,
					}

					allowResources = append(allowResources, &resourceElement)
				}
			} else {
				// Add new banned actions to the map
				if resourceActions, ok := deniedActions[resourceName]; ok {
					mergedActions := append(resourceActions.ToSlice(), statementActions...)
					deniedActions[resourceName] = minioIAMPolicy.NewActionSet(mergedActions...)
				} else {
					deniedActions[resourceName] = statement.Actions
				}
				// Remove existing val from key if necessary
				if currentResourceActions, ok := permissions[resourceName]; ok {
					var newAllowedActions []minioIAMPolicy.Action
					for _, action := range currentResourceActions.ToSlice() {
						if len(deniedActions[resourceName].Intersection(minioIAMPolicy.NewActionSet(action))) == 0 {
							// It's ok to allow this action
							newAllowedActions = append(newAllowedActions, action)
						}
					}
					permissions[resourceName] = minioIAMPolicy.NewActionSet(newAllowedActions...)
				}
			}
		}
	}
	resourcePermissions := map[string][]string{}
	for key, val := range permissions {
		var resourceActions []string
		for _, action := range val.ToSlice() {
			resourceActions = append(resourceActions, string(action))
		}
		resourcePermissions[key] = resourceActions

	}

	// environment constants
	var envConstants models.EnvironmentConstants

	envConstants.MaxConcurrentUploads = getMaxConcurrentUploadsLimit()
	envConstants.MaxConcurrentDownloads = getMaxConcurrentDownloadsLimit()

	sessionResp := &models.SessionResponse{
		Features:        getListOfEnabledFeatures(ctx, userAdminClient, session),
		Status:          models.SessionResponseStatusOk,
		Operator:        false,
		DistributedMode: erasure,
		Permissions:     resourcePermissions,
		AllowResources:  allowResources,
		CustomStyles:    customStyles,
		EnvConstants:    &envConstants,
		ServerEndPoint:  getMinIOServer(),
	}
	return sessionResp, nil
}

// getListOfEnabledFeatures returns a list of features
func getListOfEnabledFeatures(ctx context.Context, minioClient MinioAdmin, session *models.Principal) []string {
	features := []string{}
	logSearchURL := getLogSearchURL()
	oidcEnabled := oauth2.IsIDPEnabled()
	ldapEnabled := ldap.GetLDAPEnabled()

	if logSearchURL != "" {
		features = append(features, "log-search")
	}
	if oidcEnabled {
		features = append(features, "oidc-idp", "external-idp")
	}
	if ldapEnabled {
		features = append(features, "ldap-idp", "external-idp")
	}

	if session.Hm {
		features = append(features, "hide-menu")
	}
	if session.Ob {
		features = append(features, "object-browser-only")
	}
	if minioClient != nil {
		_, err := minioClient.kmsStatus(ctx)
		if err == nil {
			features = append(features, "kms")
		}
	}

	return features
}
