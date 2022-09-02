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

package main

import (
	"strings"

	contribMiddleware "github.com/dapr/components-contrib/middleware"
	dapr "github.com/mcandeia/dapr-components-go-sdk"
	httpMiddleware "github.com/mcandeia/dapr-components-go-sdk/middleware/v1"
)

type myMiddleware struct{}

func (m *myMiddleware) GetHandler(metadata contribMiddleware.Metadata) (func(httpMiddleware.MiddlewareController) error, error) {
	return func(mc httpMiddleware.MiddlewareController) error {
		mc.Next()
		reqBody := string(mc.GetResponseBody())
		mc.SetResponseHeaders(map[string]string{
			"x-marcos-header": "10",
		})
		mc.SetResponseBody([]byte(strings.ReplaceAll(reqBody, "a", "x")))
		return nil
	}, nil
}

func main() {
	dapr.MustRun(dapr.UseHttpMiddleware(&myMiddleware{}))
}
