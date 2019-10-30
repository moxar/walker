# Walker

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/moxar/walker)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/moxar/walker/master/LICENSE)

Walks around the database's schema, establish relations between tables, and build the `FROM ... JOIN` clauses dynamicaly.

Checkout [test file](./graph_test.go) and the [test samples](./resources/tests) for examples of usage.

## Features

The package builds the `FROM ... JOIN` clause of a SQL query based on the required tables. It needs a schema declaration, that can be hardcoded or fetched from the db (see mysql package).
- Path between tables is established with Dijkstra's algorythm - [Ryan Carrier's implementation](github.com/RyanCarrier/dijkstra)
- Table and relation aliasing supported

## MySQL

The walker/mysql provides a function that loads the schema from the database.

```go
// Fetch the schema from database.
schema, err := mysql.LoadSchema(ctx, db, "my_project_db")
if err != nil {
	// ...
}

// Alias lands table to countries.
schema.Alias("lands", "countries")

// Prepare the graph.
graph, err := walker.NewGraph(schema)
if err != nil {
	// ...
}

// Build the query relating users, cities and countries.
from, err := graph.From("users", "cities, "countries")
if err != nil {
	// ...
}
```
