package mw

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"
)

type PushEvent struct {
	PushID *int64  `json:"push_id,omitempty"`
	Ref    *string `json:"ref,omitempty"`
}

func AuthorizePushEvent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		body, err := io.ReadAll(ctx.Request().Body) // get all variables needed to check signature
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "body read error")
		}
		signature := ctx.Request().Header.Get("X-Hub-Signature-256")
		secret, err := hex.DecodeString(os.Getenv("WEBHOOK_SECRET"))
		if err == nil && len(signature) > 0 {
			messageMAC, hashFunc, err := messageMAC(signature) // check signature
			if err != nil || !checkMAC(body, messageMAC, secret, hashFunc) {
				return echo.NewHTTPError(http.StatusForbidden, "wrong signature")
			}
		} else {
			return echo.NewHTTPError(http.StatusForbidden, "")
		}

		if ctx.Request().Header.Get("X-Github-Event") != "push" { // check event type
			return echo.NewHTTPError(http.StatusForbidden, "wrong event")
		}
		payload := &PushEvent{}
		err = json.Unmarshal(body, &payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "json unmarshal error")
		}
		if *payload.Ref != "refs/heads/master" { // check if master branch
			return echo.NewHTTPError(http.StatusOK, "not master branch")
		}

		// go on
		return next(ctx)
	}
}

func genMAC(message, key []byte, hashFunc func() hash.Hash) []byte {
	mac := hmac.New(hashFunc, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func checkMAC(message, messageMAC, key []byte, hashFunc func() hash.Hash) bool {
	expectedMAC := genMAC(message, key, hashFunc)
	return hmac.Equal(messageMAC, expectedMAC)
}

func messageMAC(signature string) ([]byte, func() hash.Hash, error) {
	if signature == "" {
		return nil, nil, errors.New("missing signature")
	}
	sigParts := strings.SplitN(signature, "=", 2)
	if len(sigParts) != 2 {
		return nil, nil, fmt.Errorf("error parsing signature %q", signature)
	}

	var hashFunc func() hash.Hash
	switch sigParts[0] {
	case "sha256":
		hashFunc = sha256.New
	default:
		return nil, nil, fmt.Errorf("unknown hash type prefix: %q", sigParts[0])
	}

	buf, err := hex.DecodeString(sigParts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding signature %q: %v", signature, err)
	}
	return buf, hashFunc, nil
}
