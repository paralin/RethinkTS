package context

import (
	proto "github.com/paralin/rethinkts/metric"
	r "gopkg.in/dancannon/gorethink.v2"
)

/* Context for each request */
type MetricContext struct {
	MetricSeries *proto.MetricSeries

	RethinkConnection  *r.Session
	RethinkDB          r.Term
	RethinkMetricTable r.Term
	RethinkDataTable   r.Term

	MetricTablePrefix string
}
