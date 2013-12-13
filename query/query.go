package query

import (
	"bytes"
	"fmt"
	"reflect"
	"github.com/jakebasile/simpleforce"
)

// A Force.com query that constructs SOQL for you.
type Query struct {
	force       simpleforce.Force
	dest        interface{}
	constraints []Constraint
	limit       int
}

// Creates a new query for you to customize. When executed, this query will fill the given destination
// slice with the results of the query.
func New(f simpleforce.Force, dest interface{}) Query {
	return Query{
		f,
		dest,
		make([]Constraint, 0, 0),
		10,
	}
}

// Adds a Constraint to the query. All constraints added in this way are ANDed together.
func (q *Query) AddConstraint(c Constraint) {
	q.constraints = append(q.constraints, c)
}

func (q *Query) Limit(l int) {
	q.limit = l
}

// Runs the query, depositing results in the destination given on query creation.
func (q *Query) Run() error {
	err := q.force.Query(q.Generate(), q.dest)
	if err != nil {
		return err
	}
	return nil
}

// Constructs the SOQL that this query represents.
func (q *Query) Generate() string {
	sel := q.generateSelect()
	table := reflect.TypeOf(q.dest).Elem().Elem().Name()
	where := q.generateWhere()
	var limit string
	if q.limit > 0 {
		limit = fmt.Sprintf(" LIMIT %v", q.limit)
	} else {
		limit = ""
	}
	return fmt.Sprintf("SELECT %v FROM %v WHERE %v%v", sel, table, where, limit)
}

func (q *Query) generateSelect() string {
	return genSelectForType(reflect.TypeOf(q.dest).Elem().Elem(), "")
}

func genSelectForType(t reflect.Type, path string) string {
	buf := bytes.NewBufferString("")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Ptr {
			buf.WriteString(genSelectForType(field.Type.Elem(), field.Name))
			buf.WriteString(",")
		} else if field.Type.Kind() == reflect.Slice {
			// wat do
		} else {
			if len(path) > 0 {
				buf.WriteString(path + "." + field.Name)
			} else {
				buf.WriteString(field.Name)
			}
			buf.WriteString(",")
		}
	}
	s := buf.String()
	// drop last comma.
	return s[:len(s)-1]
}

func (q *Query) generateWhere() string {
	buf := bytes.NewBufferString("")
	for i, c := range q.constraints {
		buf.WriteString(c.Collapse())
		if i < len(q.constraints)-1 {
			buf.WriteString(" AND ")
		}
	}
	return buf.String()
}
