package config

import (
	"fmt"

	"github.com/vroomy/common"
	"github.com/vroomy/httpserve"
)

const (
	errInvalidRedirectValueFmt = "invalid redirect value type, expected %T and received %T"
	errInvalidTextValueFmt     = "invalid text value type, expected %T or %T and received %T"
	errInvalidXMLValueFmt      = "invalid XML value type, expected %T and received %T"
)

var exampleCommonHandler common.Handler = func(common.Context) *common.Response { return nil }

func redirectHandler(resp *common.Response) httpserve.Response {
	var (
		url string
		ok  bool
	)

	if url, ok = resp.Value.(string); !ok {
		err := fmt.Errorf(errInvalidRedirectValueFmt, url, resp.Value)
		out.ErrorWithData(err.Error(), resp.Value)
		return httpserve.NewTextResponse(500, []byte(err.Error()))
	}

	return httpserve.NewRedirectResponse(resp.StatusCode, url)
}

func textHandler(resp *common.Response) httpserve.Response {
	var bs []byte
	switch n := resp.Value.(type) {
	case string:
		bs = []byte(n)
	case []byte:
		bs = n

	default:
		err := fmt.Errorf(errInvalidTextValueFmt, "foo", bs, resp.Value)
		out.ErrorWithData(err.Error(), resp.Value)
		return httpserve.NewTextResponse(500, []byte(err.Error()))
	}

	return httpserve.NewTextResponse(resp.StatusCode, bs)
}

func xmlHandler(resp *common.Response) httpserve.Response {
	var (
		data []byte
		ok   bool
	)

	if data, ok = resp.Value.([]byte); !ok {
		err := fmt.Errorf(errInvalidXMLValueFmt, data, resp.Value)
		out.ErrorWithData(err.Error(), resp.Value)
		return httpserve.NewTextResponse(500, []byte(err.Error()))
	}

	return httpserve.NewXMLResponse(resp.StatusCode, data)
}
