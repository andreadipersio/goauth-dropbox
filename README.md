goauth-dropbox
==============

A Dropbox Core API authentication library.

[![GoDoc](https://godoc.org/github.com/andreadipersio/goauth-dropbox?status.png)](http://godoc.org/github.com/andreadipersio/goauth-dropbox)

### Usage

```Go

package main

import (
    "fmt"
    "net/http"

    "github.com/andreadipersio/goauth-dropbox/dropbox"
)

func main() {
    dropboxHandler := &dropbox.OAuth2Handler{
        Key: "my app key",
        Secret: "my app secret",

        RedirectURI: "http://localhost:8001/oauth/dropbox",

        ErrorCallback: func(w http.ResponseWriter, r *http.Request, err error) {
            http.Error(w, fmt.Sprintf("OAuth error - %v", err), 500)
        },

        SuccessCallback: func(w http.ResponseWriter, r *http.Request, token *dropbox.Token) {
            http.SetCookie(w, &http.Cookie{
                Name: "dropbox_token",
                Value: token.Token,
            })

            http.SetCookie(w, &http.Cookie{
                Name: "dropbox_uid",
                Value: token.UID,
            })
        },
    }

    http.Handle("/oauth/dropbox", dropboxHandler)
    http.ListenAndServe(":8001", nil)
}

```
