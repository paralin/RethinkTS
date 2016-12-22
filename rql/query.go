package rql

import (
	"github.com/paralin/rethinkts/context"
	"github.com/paralin/rethinkts/metric"
	r "gopkg.in/dancannon/gorethink.v2"
)

func BuildDatapointQuery(query *metric.MetricDatapointQuery, ctx *context.MetricContext) r.Term {
	currTerm := ctx.RethinkMetricTable
	timestampTerm := r.Row.Field("timestamp")
	tagTerm := r.Row.Field("tags")

	// First, constrain by time if necessary
	timeConstraint := query.TimeConstraint
	if timeConstraint != nil {
		if timeConstraint.MinTime != 0 {
			currTerm = currTerm.Filter(timestampTerm.Gt(timeConstraint.MinTime))
		}
		if timeConstraint.MaxTime != 0 {
			currTerm = currTerm.Filter(timestampTerm.Lt(timeConstraint.MaxTime))
		}
	}

	tagConstraint := query.TagConstraint
	if tagConstraint != nil {
		tags := tagConstraint.GetTags()
		for k, v := range tags {
			currTerm = currTerm.Filter(tagTerm.Field(k).Eq(v))
		}
	}

	return currTerm
}
