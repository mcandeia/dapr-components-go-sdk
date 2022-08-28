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

package state

import (
	"context"

	contribState "github.com/dapr/components-contrib/state"

	"github.com/mcandeia/dapr-components-go-sdk/internal"

	"github.com/dapr/dapr/pkg/proto/common/v1"
	proto "github.com/dapr/dapr/pkg/proto/components/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

type store struct {
	impl contribState.Store
}

func (s *store) Init(ctx context.Context, metadata *proto.MetadataRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.impl.Init(contribState.Metadata{
		Properties: metadata.Properties,
	})
}

func (s *store) Features(context.Context, *emptypb.Empty) (*proto.FeaturesResponse, error) {
	features := &proto.FeaturesResponse{
		Feature: internal.Map(s.impl.Features(), func(f contribState.Feature) string {
			return string(f)
		}),
	}

	return features, nil
}

func toDeleteRequest(req *proto.DeleteRequest) *contribState.DeleteRequest {
	return &contribState.DeleteRequest{
		Key: req.Key,
		ETag: internal.IfNotNilP(req.Etag, func(f *common.Etag) string {
			return f.Value
		}),
		Metadata: req.Metadata,
		Options: internal.IfNotNil(req.Options, func(f *common.StateOptions) contribState.DeleteStateOption {
			return contribState.DeleteStateOption{
				Concurrency: string(f.Concurrency),
				Consistency: string(f.Consistency),
			}
		}),
	}
}

func (s *store) Delete(_ context.Context, req *proto.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.impl.Delete(toDeleteRequest(req))
}

func toGetRequest(req *proto.GetRequest) *contribState.GetRequest {
	return &contribState.GetRequest{
		Key:      req.Key,
		Metadata: req.Metadata,
		Options: contribState.GetStateOption{
			Consistency: req.Consistency.String(),
		},
	}
}

func fromGetResponse(res *contribState.GetResponse) *proto.GetResponse {
	return &proto.GetResponse{
		Data: res.Data,
		Etag: internal.IfNotNil(res.ETag, func(etagValue *string) *common.Etag {
			return &common.Etag{
				Value: *etagValue,
			}
		}),
		Metadata: res.Metadata,
	}
}

func (s *store) Get(_ context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	resp, err := s.impl.Get(toGetRequest(req))
	return internal.IfNotNil(resp, fromGetResponse), err
}

func toSetRequest(req *proto.SetRequest) *contribState.SetRequest {
	return &contribState.SetRequest{
		Key:   req.Key,
		Value: req.Value,
		ETag: internal.IfNotNilP(req.Etag, func(f *common.Etag) string {
			return f.Value
		}),
		Metadata: req.Metadata,
		Options: internal.IfNotNil(req.Options, func(f *common.StateOptions) contribState.SetStateOption {
			return contribState.SetStateOption{
				Concurrency: string(f.Concurrency),
				Consistency: string(f.Consistency),
			}
		}),
	}
}

func (s *store) Set(_ context.Context, req *proto.SetRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.impl.Set(toSetRequest(req))
}

func (s *store) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, contribState.Ping(s.impl)
}

func (s *store) BulkDelete(_ context.Context, req *proto.BulkDeleteRequest) (*emptypb.Empty, error) {
	return nil, s.impl.BulkDelete(internal.Map(req.Items, func(delReq *proto.DeleteRequest) contribState.DeleteRequest {
		return *toDeleteRequest(delReq)
	}))
}

func fromBulkGetResponse(item contribState.BulkGetResponse) *proto.BulkStateItem {
	return &proto.BulkStateItem{
		Key:  item.Key,
		Data: item.Data,
		Etag: internal.IfNotNil(item.ETag, func(etagValue *string) *common.Etag {
			return &common.Etag{
				Value: *etagValue,
			}
		}),
		Error:    item.Error,
		Metadata: item.Metadata,
	}
}

func (s *store) BulkGet(_ context.Context, req *proto.BulkGetRequest) (*proto.BulkGetResponse, error) {
	got, items, err := s.impl.BulkGet(internal.Map(req.Items, func(getReq *proto.GetRequest) contribState.GetRequest {
		return *toGetRequest(getReq)
	}))
	return &proto.BulkGetResponse{
		Got:   got,
		Items: internal.Map(items, fromBulkGetResponse),
	}, err
}

func (s *store) BulkSet(_ context.Context, req *proto.BulkSetRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.impl.BulkSet(internal.Map(req.Items, func(setReq *proto.SetRequest) contribState.SetRequest {
		return *toSetRequest(setReq)
	}))
}

func NewWrapper(impl contribState.Store) proto.StateStoreServer {
	return &store{
		impl: impl,
	}
}
