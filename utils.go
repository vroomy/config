package config

import (
	"fmt"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/hatchify/errors"
	"github.com/vroomy/common"
	"github.com/vroomy/httpserve"
	"github.com/vroomy/plugins"
)

const (
	// ErrInvalidRoot is returned whe a root is longer than the request path
	ErrInvalidRoot = errors.Error("invalid root, cannot be longer than request path")
	eerrr          = 1
)

func getKeyFromRequestPath(root, requestPath string) (key string, err error) {
	// Clean request path
	requestPath = filepath.Clean(requestPath)

	if len(root) > len(requestPath) {
		err = ErrInvalidRoot
		return
	}

	key = requestPath[len(root):]
	return
}

func trimSlash(in string) (out string) {
	if len(in) == 0 {
		return
	}

	if in[len(in)-1] != '/' {
		return in
	}

	return in[:len(in)-1]
}

func getHandlerParts(handlerKey string) (key, handler string, args []string, err error) {
	spl := strings.SplitN(handlerKey, ".", 2)
	if len(spl) != 2 {
		err = fmt.Errorf("expected key and handler, received \"%s\"", handlerKey)
		return
	}

	key = spl[0]
	handler = spl[1]

	spl = strings.Split(handler, "(")
	if len(spl) == 1 {
		return
	}

	handler = spl[0]
	argsStr := spl[1]

	if argsStr[len(argsStr)-1] != ')' {
		err = ErrExpectedEndParen
		return
	}

	argsStr = argsStr[:len(argsStr)-1]
	args = strings.Split(argsStr, ",")
	return
}

func newPluginHandler(p *plugins.Plugins, handlerKey string) (h httpserve.Handler, err error) {
	var (
		key     string
		handler string
		args    []string
	)

	if key, handler, args, err = getHandlerParts(handlerKey); err != nil {
		return
	}

	var pp *plugin.Plugin
	if pp, err = p.Get(key); err != nil {
		return
	}

	var sym plugin.Symbol
	if sym, err = pp.Lookup(handler); err != nil {
		return
	}

	switch v := sym.(type) {
	case httpserve.Handler:
		h = v
	case common.Handler:
		h = newHandler(v)
	case func(args ...string) (common.Handler, error):
		var ch common.Handler
		if ch, err = v(args...); err != nil {
			return
		}

		h = newHandler(ch)

	default:
		err = fmt.Errorf("invalid handler signature encountered: %T is not supported", sym)
		return
	}

	return
}

func newHandler(c common.Handler) httpserve.Handler {
	return func(ctx *httpserve.Context) httpserve.Response {
		resp := c(ctx)
		switch {
		case resp == nil:
			return nil
		case resp.Adopted:
			return httpserve.NewAdoptResponse()
		}

		switch resp.StatusCode {
		case 204:
			return httpserve.NewNoContentResponse()
		case 301, 302:
			return redirectHandler(resp)
		}

		switch resp.ContentType {
		case "json":
			return httpserve.NewJSONResponse(resp.StatusCode, resp.Value)
		case "jsonp":
			return httpserve.NewJSONPResponse(resp.Callback, resp.Value)
		case "text":
			return textHandler(resp)
		case "xml":
			return xmlHandler(resp)
		}

		return nil
	}
}
