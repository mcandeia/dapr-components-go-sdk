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

package dapr

import (
	"errors"
	"net"
	"os"

	"google.golang.org/grpc"
)

var ErrSocketNotDefined = errors.New("socket env `DAPR_COMPONENT_SOCKET_PATH` must be set")

const (
	unixSocketPathEnvVar = "DAPR_COMPONENT_SOCKET_PATH"
)

// Run starts the component server with the given options.
func Run(opts ...Option) error {
	socket, ok := os.LookupEnv(unixSocketPathEnvVar)
	if !ok {
		return ErrSocketNotDefined
	}

	lis, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	svcOpts := &componentOpts{}
	for _, opt := range opts {
		opt(svcOpts)
	}

	server := grpc.NewServer()

	if err = svcOpts.apply(server); err != nil {
		return err
	}

	return server.Serve(lis)
}
