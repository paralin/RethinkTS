package rql

import (
	"errors"

	"github.com/fuserobotics/rethinkts/context"
	"github.com/fuserobotics/rethinkts/metric"
)

func CheckDuplicate(dp *metric.MetricDatapoint, ctx *context.MetricContext) error {
	switch ctx.MetricSeries.DedupeStrategy {
	case metric.MetricSeries_NONE:
		return nil
	}
	return errors.New("Deduplication strat not implemented.")
}

func RecordDatapoint(dp *metric.MetricDatapoint, ctx *context.MetricContext) error {
	if err := CheckDuplicate(dp, ctx); err != nil {
		return err
	}

	_, err := ctx.RethinkMetricTable.Insert(dp).RunWrite(ctx.RethinkConnection)
	return err
}
