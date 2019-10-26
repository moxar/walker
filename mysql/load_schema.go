package mysql

import (
	"context"
	"fmt"

	"github.com/moxar/walker"
)

// Selecter is implemented by types capable of running SELECT queries.
type Selecter interface {
	SelectContext(ctx context.Context, dts interface{}, query string, args ...interface{}) error
}

// foreignKey definition.
type foreignKey struct {
	TableName            string
	ColumnName           string
	ReferencedTableName  string
	ReferencedColumnName string
}

func (fk foreignKey) join() string {
	return fmt.Sprintf("%s.%s = %s.%s", fk.TableName, fk.ColumnName, fk.ReferencedTableName, fk.ReferencedColumnName)
}

// LoadSchema for the given table.
func LoadSchema(ctx context.Context, sel Selecter, db string) (*walker.Schema, error) {
	var tables []walker.Vertex
	if err := sel.SelectContext(ctx, &tables, "SHOW TABLES;"); err != nil {
		return nil, err
	}

	var indexes []foreignKey
	if err := sel.SelectContext(ctx, &indexes, `
		SELECT table_name tablename, 
			column_name columnname, 
			referenced_table_name referencedtablename, 
			referenced_column_name referencedcolumnname 
		FROM information_schema.key_column_usage 
		WHERE table_schema = ? AND referenced_table_schema = ?
	`, db, db); err != nil {
		return nil, err
	}

	s := walker.NewSchema()
	s.Verticies = make([]walker.Vertex, 0, len(tables))
	for _, t := range tables {
		s.Verticies = append(s.Verticies, t)
	}

	uniq := make(map[[2]string]*walker.Arc)
	s.Arcs = make([]walker.Arc, 0, len(indexes))
	for _, index := range indexes {
		pair := [2]string{index.TableName, index.ReferencedTableName}
		if arc, ok := uniq[pair]; ok {
			arc.Link += " AND " + index.join()
			continue
		}
		arc := walker.Arc{
			Source:      index.TableName,
			Destination: index.ReferencedTableName,
			Link:        index.join(),
		}
		uniq[pair] = &arc
		s.Arcs = append(s.Arcs, arc)
	}

	return s, nil
}
