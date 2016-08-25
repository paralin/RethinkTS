package metric_impl

import (
	"errors"
	"io"

	"github.com/fuserobotics/rethinkts/context"
	"github.com/fuserobotics/rethinkts/metric"
	"github.com/fuserobotics/rethinkts/rql"
	"github.com/mitchellh/mapstructure"
	netctx "golang.org/x/net/context"
	"google.golang.org/grpc"
	r "gopkg.in/dancannon/gorethink.v2"
)

type metricServer struct {
	BaseContext context.MetricContext
}

func (ms *metricServer) doRecordDatapoint(in *metric.RecordDatapointRequest) error {
	ctx, err := context.BuildContextFromRequest(&ms.BaseContext, in.Context)
	if err != nil {
		return err
	}

	if err := in.Datapoint.Validate(); err != nil {
		return err
	}

	if err := rql.RecordDatapoint(in.Datapoint, &ctx); err != nil {
		context.PurgeContextFromCache(in.Context.Identifier.Id)
		return err
	}
	return nil
}

func (ms *metricServer) RecordDatapoint(nctx netctx.Context, in *metric.RecordDatapointRequest) (*metric.RecordDatapointResponse, error) {
	ms.doRecordDatapoint(in)
	return &metric.RecordDatapointResponse{NumRecorded: 1}, nil
}

func (ms *metricServer) RecordDatapointStream(stream metric.MetricService_RecordDatapointStreamServer) error {
	rejectedDatapoints := make([]*metric.MetricDatapoint, 0)
	var numRecorded int32
	for {
		dp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := ms.doRecordDatapoint(dp); err != nil {
			if dp.Datapoint != nil {
				rejectedDatapoints = append(rejectedDatapoints, dp.Datapoint)
			}
		} else {
			numRecorded++
		}
	}
	return stream.SendAndClose(&metric.RecordDatapointResponse{
		NumRecorded: numRecorded,
		Rejected:    rejectedDatapoints,
	})
}

func (ms *metricServer) ListDatapoint(in *metric.ListDatapointRequest, stream metric.MetricService_ListDatapointServer) error {
	// Build context
	ctx, err := context.BuildContextFromRequest(&ms.BaseContext, in.Context)
	if err != nil {
		return err
	}

	if !in.IncludeInitial && !in.Tail {
		return errors.New("Include initial AND tail cannot both be false.")
	}

	if in.Query == nil {
		return errors.New("Query cannot be nil.")
	}

	{
		err := stream.Send(&metric.ListDatapointResponse{
			ResponseType: metric.ListDatapointResponse_LIST_DATAPOINT_SERIES_DETAILS,
			Series:       ctx.MetricSeries,
		})
		if err != nil {
			return err
		}
	}

	// Query
	query := rql.BuildDatapointQuery(in.Query, &ctx)
	if in.Tail {
		query = query.Changes(r.ChangesOpts{
			IncludeInitial: in.IncludeInitial,
		})

		cursor, err := query.Run(ctx.RethinkConnection)
		if err != nil {
			return err
		}
		defer cursor.Close()

		change := make(chan r.ChangeResponse)
		cursor.Listen(change)

		for {
			select {
			case <-stream.Context().Done():
				return nil
			case ch, ok := <-change:
				if !ok {
					return nil
				}
				var updateType metric.ListDatapointResponse_ListDatapointResponseType
				dp := &metric.ListDatapointResponse{}
				dp.Datapoint = &metric.MetricDatapoint{}
				if ch.NewValue != nil {
					if ch.OldValue != nil {
						updateType = metric.ListDatapointResponse_LIST_DATAPOINT_REPLACE
					} else {
						updateType = metric.ListDatapointResponse_LIST_DATAPOINT_ADD
					}
					newValueMap := ch.NewValue.(map[string]interface{})
					if err := mapstructure.Decode(newValueMap, dp.Datapoint); err != nil {
						return err
					}
				} else if ch.OldValue != nil {
					updateType = metric.ListDatapointResponse_LIST_DATAPOINT_DEL
					newValueMap := ch.OldValue.(map[string]interface{})
					if err := mapstructure.Decode(newValueMap, dp.Datapoint); err != nil {
						return err
					}
				} else {
					break
				}
				dp.ResponseType = updateType

				if err := stream.Send(dp); err != nil {
					return err
				}
			}
		}

	} else {
		cursor, err := query.Run(ctx.RethinkConnection)
		if err != nil {
			return err
		}
		defer cursor.Close()

		for {
			item := &metric.MetricDatapoint{}

			if !cursor.Next(item) {
				break
			}

			err := stream.Send(&metric.ListDatapointResponse{
				ResponseType: metric.ListDatapointResponse_LIST_DATAPOINT_ADD,
				Datapoint:    item,
			})

			if err != nil {
				return err
			}
		}

		if err := cursor.Err(); err != nil {
			return err
		}
	}

	return nil
}

func (ms *metricServer) ListMetric(nctx netctx.Context, in *metric.ListMetricRequest) (*metric.ListMetricResponse, error) {
	query := ms.BaseContext.RethinkDataTable
	cursor, err := query.Run(ms.BaseContext.RethinkConnection)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var result []*metric.MetricSeries
	if err := cursor.All(&result); err != nil {
		return nil, err
	}

	return &metric.ListMetricResponse{
		Metric: result,
	}, nil
}

func (ms *metricServer) CreateMetric(nctx netctx.Context, in *metric.CreateMetricRequest) (*metric.CreateMetricResponse, error) {
	if err := in.Metric.Validate(); err != nil {
		return nil, err
	}

	_, err := ms.BaseContext.RethinkDataTable.Insert(in.Metric).RunWrite(ms.BaseContext.RethinkConnection)
	if err != nil {
		return nil, err
	}
	return &metric.CreateMetricResponse{
		Metric: in.Metric,
	}, nil
}

func RegisterServer(ctx *context.MetricContext, grpcServer *grpc.Server) {
	metric.RegisterMetricServiceServer(grpcServer, &metricServer{
		BaseContext: *ctx,
	})
}
