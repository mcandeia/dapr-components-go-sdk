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

package http

import (
	"context"
	"errors"

	contribMetadata "github.com/dapr/components-contrib/metadata"
	contribMiddleware "github.com/dapr/components-contrib/middleware"

	proto "github.com/dapr/dapr/pkg/proto/components/v1"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var defaultMiddleware = &middleware{}

type middleware struct {
	middlewareFactory Middleware
	handler           func(MiddlewareController) error
}

//nolint:nosnakecase
func (m *middleware) Handle(svc proto.HttpMiddleware_HandleServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch sourceErr := r.(type) {
			case string:
				err = errors.New(sourceErr)
			case error:
				err = sourceErr
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
		}
	}()
	controller := &middlewareController{client: svc}

	if handlerErr := m.handler(controller); handlerErr != nil {
		return handlerErr
	}
	return err
}

func (m *middleware) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *middleware) Init(ctx context.Context, metadata *proto.MetadataRequest) (*emptypb.Empty, error) {
	handler, err := s.middlewareFactory.GetHandler(contribMiddleware.Metadata{
		Base: contribMetadata.Base{Properties: metadata.Properties},
	})
	if err != nil {
		return nil, err
	}
	s.handler = handler

	return &emptypb.Empty{}, nil
}

func Register(server *grpc.Server, middleware Middleware) {
	defaultMiddleware.middlewareFactory = middleware
	proto.RegisterHttpMiddlewareServer(server, defaultMiddleware)
}
