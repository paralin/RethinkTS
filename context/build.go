package context

import (
	"errors"

	"github.com/paralin/rethinkts/metric"
	"github.com/paralin/rethinkts/util"
	r "gopkg.in/dancannon/gorethink.v2"
)

var ContextCache map[string]MetricContext = make(map[string]MetricContext)

func BuildBaseContext(rctx *r.Session, db, metaTable, metricTablePrefix string) (MetricContext, error) {
	res := MetricContext{
		RethinkConnection: rctx,
		RethinkDB:         r.DB(db),
		RethinkDataTable:  r.DB(db).Table(metaTable),
		MetricTablePrefix: metricTablePrefix,
	}

	// Check DB exists
	_, err := res.RethinkDB.Info().Run(rctx)
	if err != nil {
		return res, err
	}

	// Check meta table exists
	_, err = res.RethinkDataTable.Info().Run(rctx)
	if err != nil {
		return res, err
	}

	return res, nil
}

func BuildContextFromRequest(baseCtx *MetricContext, rctx *metric.RequestContext) (MetricContext, error) {
	res := *baseCtx

	if rctx == nil {
		return res, errors.New("Request context cannot be null.")
	}

	if rctx.Identifier == nil {
		return res, errors.New("Request context identifier cannot be null.")
	}

	if rctx.Identifier.Id == "" {
		return res, errors.New("Request context metric id cannot be null.")
	}

	rctxId := rctx.Identifier.Id
	if res, ok := ContextCache[rctxId]; ok {
		return res, nil
	}

	// Find metric by identifier
	cursor, err := baseCtx.RethinkDataTable.Get(rctx.Identifier.Id).Run(baseCtx.RethinkConnection)
	if err != nil {
		return res, err
	}

	metricSeries := &metric.MetricSeries{}
	if err := cursor.One(metricSeries); err != nil {
		return res, err
	}

	res.MetricSeries = metricSeries

	// Ask for info on metric table
	res.RethinkMetricTable = res.RethinkDB.Table(res.MetricTablePrefix + util.TableNameForMetricId(metricSeries.Id))
	_, err = res.RethinkMetricTable.Info().Run(res.RethinkConnection)
	if err == nil {
		ContextCache[rctxId] = res
	}
	return res, err
}

func PurgeContextFromCache(id string) {
	delete(ContextCache, id)
}
