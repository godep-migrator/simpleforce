package simpleforce

import (
	"bytes"
)

// Creates an '=' clause for a string value.
func (c Constraint) EqualsString(right string) Constraint {
	c.op = "="
	c.right = "'" + right + "'"
	return c
}

// Creates a '<>' clause for a string value.
func (c Constraint) NotEqualsString(right string) Constraint {
	c.op = "<>"
	c.right = "'" + right + "'"
	return c
}

// Creates a '>' clause for a string value.
func (c Constraint) GreaterString(right string) Constraint {
	c.op = ">"
	c.right = "'" + right + "'"
	return c
}

// Creates a '>=' clause for a string value.
func (c Constraint) GreaterEqualsString(right string) Constraint {
	c.op = ">="
	c.right = "'" + right + "'"
	return c
}

// Creates a '<' clause for a string value.
func (c Constraint) LessString(right string) Constraint {
	c.op = "<"
	c.right = "'" + right + "'"
	return c
}

// Creates a '<=' clause for a string value.
func (c Constraint) LessEqualsString(right string) Constraint {
	c.op = "<="
	c.right = "'" + right + "'"
	return c
}

// Creates an IN clause for a string value.
func (c Constraint) InString(in ...string) Constraint {
	c.op = " IN "
	buf := bytes.NewBufferString("(")
	for i, s := range in {
		buf.WriteString("'" + s + "'")
		if i < len(in)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	c.right = buf.String()
	return c
}

// Creates a NOT IN clause for a string value.
func (c Constraint) NotInString(in ...string) Constraint {
	c.op = " NOT IN "
	buf := bytes.NewBufferString("(")
	for i, s := range in {
		buf.WriteString("'" + s + "'")
		if i < len(in)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	c.right = buf.String()
	return c
}

// Creates a LIKE clause for a string value.
func (c Constraint) LikeString(like string) Constraint {
	c.op = " LIKE "
	c.right = "'" + like + "'"
	return c
}
