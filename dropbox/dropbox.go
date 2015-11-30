// Package dropbox implement OAuth2 authentication for Dropbox
// providing an handler (OAuth2Handler) which perform OAuth token authorization
// and exchange.
package dropbox

import (
    "fmt"

    "encoding/json"

    "net/url"
    "net/http"
)

const (
    authorizationURL = "https://www.dropbox.com/1/oauth2/authorize?"  +
                       "client_id=%v&response_type=code&redirect_uri=%v"

    tokenExchangeURL = "https://api.dropbox.com/1/oauth2/token"
)

type Token struct {
    // Dropbox user id
    UID   string `json:"uid"`

    // Dropbox bearer access token. Can be used to make 
    // API call to Dropbox CORE API.
    // https://www.dropbox.com/developers/core/docs
    Token string `json:"access_token"`
	Error *string `json:"error"`
}

type OAuth2Handler struct {
    // App Key
    Key,

    // App Secret
    Secret,

    // OAuth redirect URL
    RedirectURI string

    // UID and access token
    Token *Token

    // SuccessCallback is executed when TokenExchange succeed
    SuccessCallback func(http.ResponseWriter, *http.Request, *Token)

    // ErrorCallback is executed when any of the OAuth step fails
    ErrorCallback   func(http.ResponseWriter, *http.Request, error)
}

func (h *OAuth2Handler) AuthorizeURL() string {
    return fmt.Sprintf(authorizationURL, h.Key, h.RedirectURI)
}

// TokenExchange method convert an auth code to a bearer token
// https://www.dropbox.com/developers/core/docs#oa2-token
func (h *OAuth2Handler) TokenExchange(authcode string) (*Token, error) {
    data := url.Values{}

    data.Add("code", authcode)
    data.Add("grant_type", "authorization_code")
    data.Add("client_id", h.Key)
    data.Add("client_secret", h.Secret)
    data.Add("redirect_uri", h.RedirectURI)

    resp, err := http.PostForm(tokenExchangeURL, data)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    token := &Token{}

    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(token)

    return token, nil
}

// If no auth code is found, then redirect to dropbox authorization endpoint,
// otherwise try to exchange the auth code with a bearer token, by invoking
// OAuth2Handler.TokenExchange.
// On success token is passed to OAuth2Handler.SuccessCallback,
// otherwise error is passed to OAuth2Handler.ErrorCallback
// (error is a string - error_code: error_description).
func (h *OAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    authcode := r.FormValue("code")
    oauthErrCode, oauthErrMsg := r.FormValue("error"), r.FormValue("error_description")

    // oauthErrCode --- http://tools.ietf.org/html/rfc6749#section-4.1.2.1
    if oauthErrCode != "" {
        h.ErrorCallback(w, r, fmt.Errorf("%v: %v", oauthErrCode, oauthErrMsg))
        return
    }

    if authcode == "" {
        http.Redirect(w, r, h.AuthorizeURL(), 302)
        return
    }

    if token, err := h.TokenExchange(authcode); err != nil {
        h.ErrorCallback(w, r, err)
    } else {
        h.SuccessCallback(w, r, token)
    }
}
