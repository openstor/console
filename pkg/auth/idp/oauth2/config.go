// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package oauth2

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"

	"github.com/openstor/console/pkg/auth/token"
	"github.com/openstor/openstor-go/v7/pkg/set"
	"github.com/openstor/pkg/v3/env"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/oauth2"
	xoauth2 "golang.org/x/oauth2"
)

// ProviderConfig - OpenID IDP Configuration for console.
type ProviderConfig struct {
	URL                      string
	DisplayName              string // user-provided - can be empty
	ClientID, ClientSecret   string
	HMACSalt, HMACPassphrase string
	Scopes                   string
	Userinfo                 bool
	RedirectCallbackDynamic  bool
	RedirectCallback         string
	EndSessionEndpoint       string
	RoleArn                  string // can be empty
}

// GetOauth2Provider instantiates a new oauth2 client using the configured credentials
// it returns a *Provider object that contains the necessary configuration to initiate an
// oauth2 authentication flow.
//
// We only support Authentication with the Authorization Code Flow - spec:
// https://openid.net/specs/openid-connect-core-1_0.html#CodeFlowAuth
func (pc ProviderConfig) GetOauth2Provider(name string, scopes []string, r *http.Request, clnt *http.Client) (provider *Provider, err error) {
	var ddoc DiscoveryDoc
	ddoc, err = parseDiscoveryDoc(r.Context(), pc.URL, clnt)
	if err != nil {
		return nil, err
	}

	supportedResponseTypes := set.NewStringSet()
	for _, responseType := range ddoc.ResponseTypesSupported {
		// FIXME: ResponseTypesSupported is a JSON array of strings - it
		// may not actually have strings with spaces inside them -
		// making the following code unnecessary.
		for _, s := range strings.Fields(responseType) {
			supportedResponseTypes.Add(s)
		}
	}

	isSupported := requiredResponseTypes.Difference(supportedResponseTypes).IsEmpty()
	if !isSupported {
		return nil, fmt.Errorf("expected 'code' response type - got %s, login not allowed", ddoc.ResponseTypesSupported)
	}

	// If provided scopes are empty we use the user configured list or a default list.
	if len(scopes) == 0 {
		for _, s := range strings.Split(pc.Scopes, ",") {
			w := strings.TrimSpace(s)
			if w == "" {
				continue
			}
			scopes = append(scopes, w)
		}
		if len(scopes) == 0 {
			scopes = defaultScopes
		}
	}

	redirectURL := pc.RedirectCallback
	if pc.RedirectCallbackDynamic {
		// dynamic redirect if set, will generate redirect URLs
		// dynamically based on incoming requests.
		redirectURL = getLoginCallbackURL(r)
	}

	// add "openid" scope always.
	scopes = append(scopes, "openid")

	client := new(Provider)
	client.oauth2Config = &xoauth2.Config{
		ClientID:     pc.ClientID,
		ClientSecret: pc.ClientSecret,
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  ddoc.AuthEndpoint,
			TokenURL: ddoc.TokenEndpoint,
		},
		Scopes: scopes,
	}

	client.IDPName = name
	client.UserInfo = pc.Userinfo
	client.client = clnt

	return client, nil
}

// GetStateKeyFunc - return the key function used to generate the authorization
// code flow state parameter.

func (pc ProviderConfig) GetStateKeyFunc() StateKeyFunc {
	return func() []byte {
		return pbkdf2.Key([]byte(pc.HMACPassphrase), []byte(pc.HMACSalt), 4096, 32, sha1.New)
	}
}

func (pc ProviderConfig) GetARNInf() string {
	return pc.RoleArn
}

type OpenIDPCfg map[string]ProviderConfig

func GetSTSEndpoint() string {
	return strings.TrimSpace(env.Get(ConsoleMinIOServer, "http://localhost:9000"))
}

func GetIDPURL() string {
	return env.Get(ConsoleIDPURL, "")
}

func GetIDPClientID() string {
	return env.Get(ConsoleIDPClientID, "")
}

func GetIDPUserInfo() bool {
	return env.Get(ConsoleIDPUserInfo, "") == "on"
}

func GetIDPSecret() string {
	return env.Get(ConsoleIDPSecret, "")
}

// Public endpoint used by the identity oidcProvider when redirecting
// the user after identity verification
func GetIDPCallbackURL() string {
	return env.Get(ConsoleIDPCallbackURL, "")
}

func GetIDPCallbackURLDynamic() bool {
	return env.Get(ConsoleIDPCallbackURLDynamic, "") == "on"
}

func IsIDPEnabled() bool {
	return GetIDPURL() != "" &&
		GetIDPClientID() != ""
}

// GetPassphraseForIDPHmac returns passphrase for the pbkdf2 function used to sign the oauth2 state parameter
func getPassphraseForIDPHmac() string {
	return env.Get(ConsoleIDPHmacPassphrase, token.GetPBKDFPassphrase())
}

// GetSaltForIDPHmac returns salt for the pbkdf2 function used to sign the oauth2 state parameter
func getSaltForIDPHmac() string {
	return env.Get(ConsoleIDPHmacSalt, token.GetPBKDFSalt())
}

// getIDPScopes return default scopes during the IDP login request
func getIDPScopes() string {
	return env.Get(ConsoleIDPScopes, "openid,profile,email")
}
