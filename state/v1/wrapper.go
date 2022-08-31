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
	"encoding/json"

	contribState "github.com/dapr/components-contrib/state"
	contribQuery "github.com/dapr/components-contrib/state/query"

	"github.com/mcandeia/dapr-components-go-sdk/internal"

	"github.com/dapr/dapr/pkg/proto/common/v1"
	proto "github.com/dapr/dapr/pkg/proto/components/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	consistencyEventual   = "eventual"
	consistencyStrong     = "strong"
	concurrencyLastWrite  = "last-write"
	concurrencyFirstWrite = "first-write"
)

var defaultStore = &store{}

type store struct {
	impl          contribState.Store
	transactional contribState.TransactionalStore
	querier       contribState.Querier
}

//nolint:nosnakecase
var consistencyModels = map[common.StateOptions_StateConsistency]string{
	common.StateOptions_CONSISTENCY_EVENTUAL:    consistencyEventual,
	common.StateOptions_CONSISTENCY_STRONG:      consistencyStrong,
	common.StateOptions_CONSISTENCY_UNSPECIFIED: consistencyStrong,
}

//nolint:nosnakecase
func toConsistency(consistency common.StateOptions_StateConsistency) string {
	c, ok := consistencyModels[consistency]
	if !ok {
		return consistencyStrong
	}
	return c
}

//nolint:nosnakecase
var concurrencyModels = map[common.StateOptions_StateConcurrency]string{
	common.StateOptions_CONCURRENCY_FIRST_WRITE: concurrencyFirstWrite,
	common.StateOptions_CONCURRENCY_LAST_WRITE:  concurrencyLastWrite,
	common.StateOptions_CONCURRENCY_UNSPECIFIED: concurrencyLastWrite,
}

//nolint:nosnakecase
func toConcurrency(concurrency common.StateOptions_StateConcurrency) string {
	c, ok := concurrencyModels[concurrency]
	if !ok {
		return concurrencyLastWrite
	}
	return c
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
				Concurrency: toConcurrency(f.Concurrency),
				Consistency: toConsistency(f.Consistency),
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
		ContentType: internal.IfNotNil(res.ContentType, func(f *string) string {
			return *f
		}),
		Metadata: res.Metadata,
	}
}

func (s *store) Get(_ context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	resp, err := s.impl.Get(toGetRequest(req))
	return internal.IfNotNil(resp, fromGetResponse), err
}

// dataParser is used to parse content by its content type
var dataParser = map[string]func([]byte) (any, error){
	"application/json": func(b []byte) (any, error) {
		var result any
		return result, json.Unmarshal(b, &result)
	},
}

func toSetRequest(req *proto.SetRequest) *contribState.SetRequest {
	var value any = req.Value
	var contentType *string
	if req.ContentType != "" {
		contentType = &req.ContentType
	}
	if ct, ok := req.Metadata["contentType"]; ok {
		contentType = &ct
	}

	if contentType != nil {
		if parser, ok := dataParser[*contentType]; ok {
			v, _ := parser(req.Value)
			value = v
		}
	}

	return &contribState.SetRequest{
		Key:   req.Key,
		Value: value,
		ETag: internal.IfNotNilP(req.Etag, func(f *common.Etag) string {
			return f.Value
		}),
		ContentType: contentType,
		Metadata:    req.Metadata,
		Options: internal.IfNotNil(req.Options, func(f *common.StateOptions) contribState.SetStateOption {
			return contribState.SetStateOption{
				Concurrency: toConcurrency(f.Concurrency),
				Consistency: toConsistency(f.Consistency),
			}
		}),
	}
}

func (s *store) Set(_ context.Context, req *proto.SetRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.impl.Set(toSetRequest(req))
}

func (s *store) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
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
		ContentType: internal.IfNotNil(item.ContentType, func(f *string) string {
			return *f
		}),
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

func toTransactionalStateOperation(op *proto.TransactionalStateOperation) contribState.TransactionalStateOperation {
	var (
		request   any
		operation contribState.OperationType
	)
	if delete := op.GetDelete(); delete != nil {
		request = *toDeleteRequest(delete)
		operation = contribState.Delete
	} else {
		request = *toSetRequest(op.GetSet())
		operation = contribState.Upsert
	}

	return contribState.TransactionalStateOperation{
		Request:   request,
		Operation: operation,
	}
}

func (s *store) Multi(_ context.Context, req *proto.TransactionalStateRequest) (*emptypb.Empty, error) {
	if s.transactional == nil {
		return &emptypb.Empty{}, status.Errorf(codes.Unimplemented, "method Multi not implemented")
	}

	return &emptypb.Empty{}, s.transactional.Multi(&contribState.TransactionalStateRequest{
		Operations: internal.Map(req.Operations, toTransactionalStateOperation),
		Metadata:   req.Metadata,
	})
}

func (s *store) Query(_ context.Context, req *proto.QueryRequest) (*proto.QueryResponse, error) {
	filters, err := internal.MapValuesErr(req.Query.Filter, func(f *anypb.Any) (any, error) {
		var v any
		return v, json.Unmarshal(f.Value, &v)
	})
	if err != nil {
		return nil, err
	}
	resp, err := s.querier.Query(&contribState.QueryRequest{
		Query: contribQuery.Query{
			Filters: filters,
			Sort: internal.Map(req.Query.Sort, func(s *proto.Sorting) contribQuery.Sorting {
				return contribQuery.Sorting{
					Key:   s.Key,
					Order: s.Order.String(),
				}
			}),
			Page: contribQuery.Pagination{
				Limit: int(req.Query.Pagination.Limit),
				Token: req.Query.Pagination.Token,
			},
			Filter: nil,
		},
		Metadata: req.Metadata,
	})
	if err != nil {
		return nil, err
	}

	return &proto.QueryResponse{
		Items: internal.Map(resp.Results, func(item contribState.QueryItem) *proto.QueryItem {
			return &proto.QueryItem{
				Key:  item.Key,
				Data: item.Data,
				Etag: internal.IfNotNil(item.ETag, func(etagValue *string) *common.Etag {
					return &common.Etag{
						Value: *etagValue,
					}
				}),
				Error: item.Error,
				ContentType: internal.IfNotNil(item.ContentType, func(f *string) string {
					return *f
				}),
			}
		}),
		Token:    resp.Token,
		Metadata: resp.Metadata,
	}, nil
}

func Register(server *grpc.Server, store contribState.Store) {
	defaultStore.impl = store
	proto.RegisterStateStoreServer(server, defaultStore)
	if trtnl, ok := store.(contribState.TransactionalStore); ok {
		proto.RegisterTransactionalStateStoreServer(server, defaultStore)
		defaultStore.transactional = trtnl
	}

	if querier, ok := store.(contribState.Querier); ok {
		proto.RegisterQueriableStateStoreServer(server, defaultStore)
		defaultStore.querier = querier
	}
}
