RethinkDB Time Series
=====================

RethinkTS is a distributed engine to collect, buffer, and store time series data in RethinkDB.

Note: this codebase has been abandoned in favor of my State Stream project.

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

Each metric series gets its own table. There is a table to list the actual metric series metadata keyed by ID.

When the server is started it must be provided with a RethinkDB database name and optionally an override for the prefix for table names for metrics. It also must be provided with the table name for the list of metrics and metric metadata.

These are the types of requests to implement:

 - [x] Record datapoint
 - [x] List / stream datapoints with various queries
 - [ ] Aggregate datapoints, cache aggregations (redis?)
 - [x] List metrics, create metrics, etc

When a request comes in, a MetricContext is built. This context contains:

 - RethinkDB session
 - RethinkDB database reference
 - RethinkDB table reference

Potential Problems
==================

What happens when:

 - We can't record a datapoint - just return an error and allow the caller to retry?
 - Timestamp is really old - how do we prevent desynced clocks from messing things up?

