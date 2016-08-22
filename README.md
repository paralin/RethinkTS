RethinkDB Time Series
=====================

RethinkTS is a distributed engine to collect, buffer, and store time series data in RethinkDB.

Storage
=======

A "metric" is a series of "datapoints" indexed by time. Each metric can have "tags" that can be placed on datapoints.

API
===

There is a server in this repo that hosts a few services:

 - GRPC Metric Server (5000)
 - REST GRPC Gateway (8080)

Design
======

Each metric series gets its own table.

These are the types of requests to implement:

 - [ ] Record datapoint
 - [ ] List / stream datapoints with various queries
 - [ ] Aggregate datapoints, cache aggregations (redis?)

When a request comes in, a MetricContext is built. This context contains:

 - RethinkDB session
 - RethinkDB database reference
 - RethinkDB table reference
