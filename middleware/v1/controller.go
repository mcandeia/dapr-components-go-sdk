/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//nolint:nosnakecase
package http

import (
	"context"
	"errors"

	contribMiddleware "github.com/dapr/components-contrib/middleware"
	proto "github.com/dapr/dapr/pkg/proto/components/v1"
)

var ErrResponseNotMatch = errors.New("received an unexpected response")

type HTTPHeaders struct {
	Headers map[string]string
	URI     string
	Method  string
}

type MiddlewareController interface {
	GetRequestBody() []byte
	SetRequestBody([]byte)
	GetResponseBody() []byte
	SetResponseBody([]byte)
	GetRequestHeaders() HTTPHeaders
	SetRequestHeaders(HTTPHeaders)
	GetResponseHeaders() map[string]string
	SetResponseHeaders(map[string]string)
	GetStatusCode() int
	SetStatusCode(int)
	Next()

	Context() context.Context
}

// Middleware is a http Middleware interface.
type Middleware interface {
	GetHandler(contribMiddleware.Metadata) (func(MiddlewareController) error, error)
}

type middlewareController struct {
	client proto.HttpMiddleware_HandleServer //nolint:nosnakecase
}

func abortOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func asResponse[T any](resp *proto.CommandResponse) (T, error) {
	if r, ok := resp.Response.(T); ok {
		return r, nil
	}
	var zero T
	return zero, ErrResponseNotMatch
}

func sendAndRecv[T any](m *middlewareController, command *proto.Command) T {
	abortOnErr(m.client.Send(command))

	resp, err := m.client.Recv()
	abortOnErr(err)

	commandResp, err := asResponse[T](resp)
	abortOnErr(err)

	return commandResp
}

func (m *middlewareController) GetRequestBody() []byte {
	commandResp := sendAndRecv[*proto.CommandResponse_GetReqBody](m, &proto.Command{
		Command: &proto.Command_GetReqBody{
			GetReqBody: &proto.GetRequestBodyCommand{},
		},
	})
	return commandResp.GetReqBody.Data
}

func (m *middlewareController) SetRequestBody(data []byte) {
	abortOnErr(
		m.client.Send(&proto.Command{
			Command: &proto.Command_SetReqBody{
				SetReqBody: &proto.SetRequestBodyCommand{
					Data: data,
				},
			},
		}),
	)
}

func (m *middlewareController) GetResponseBody() []byte {
	commandResp := sendAndRecv[*proto.CommandResponse_GetRespBody](m, &proto.Command{
		Command: &proto.Command_GetRespBody{
			GetRespBody: &proto.GetResponseBodyCommand{},
		},
	})
	return commandResp.GetRespBody.Data
}

func (m *middlewareController) SetResponseBody(data []byte) {
	abortOnErr(m.client.Send(&proto.Command{
		Command: &proto.Command_SetRespBody{
			SetRespBody: &proto.SetResponseBodyCommand{
				Data: data,
			},
		},
	}))
}

func (m *middlewareController) GetRequestHeaders() HTTPHeaders {
	commandResp := sendAndRecv[*proto.CommandResponse_GetReqHeaders](m, &proto.Command{
		Command: &proto.Command_GetReqHeaders{
			GetReqHeaders: &proto.GetRequestHeadersCommand{},
		},
	})

	return HTTPHeaders{
		Headers: commandResp.GetReqHeaders.Headers,
		URI:     commandResp.GetReqHeaders.Uri,
		Method:  commandResp.GetReqHeaders.Method,
	}
}

func (m *middlewareController) SetRequestHeaders(headers HTTPHeaders) {
	abortOnErr(
		m.client.Send(&proto.Command{
			Command: &proto.Command_SetReqHeaders{
				SetReqHeaders: &proto.SetRequestHeadersCommand{
					Method:  headers.Method,
					Uri:     headers.URI,
					Headers: headers.Headers,
				},
			},
		}),
	)
}

func (m *middlewareController) GetResponseHeaders() map[string]string {
	commandResp := sendAndRecv[*proto.CommandResponse_GetRespHeaders](m, &proto.Command{
		Command: &proto.Command_GetRespHeaders{
			GetRespHeaders: &proto.GetResponseHeadersCommand{},
		},
	})

	return commandResp.GetRespHeaders.Headers
}

func (m *middlewareController) SetResponseHeaders(data map[string]string) {
	abortOnErr(
		m.client.Send(&proto.Command{
			Command: &proto.Command_SetRespHeaders{
				SetRespHeaders: &proto.SetResponseHeadersCommand{
					Headers: data,
				},
			},
		}),
	)
}

func (m *middlewareController) GetStatusCode() int {
	commandResp := sendAndRecv[*proto.CommandResponse_GetRespHeaders](m, &proto.Command{
		Command: &proto.Command_GetRespHeaders{
			GetRespHeaders: &proto.GetResponseHeadersCommand{},
		},
	})

	return int(commandResp.GetRespHeaders.StatusCode)
}

func (m *middlewareController) SetStatusCode(statusCode int) {
	abortOnErr(
		m.client.Send(&proto.Command{
			Command: &proto.Command_SetRespStatus{
				SetRespStatus: &proto.SetResponseStatusCodeCommand{
					StatusCode: int32(statusCode),
				},
			},
		}),
	)
}

func (m *middlewareController) Next() {
	abortOnErr(
		m.client.Send(&proto.Command{
			Command: &proto.Command_ExecNext{
				ExecNext: &proto.ExecNextCommand{},
			},
		}),
	)
}

func (m *middlewareController) Context() context.Context {
	return m.client.Context()
}
