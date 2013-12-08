package simpleforce

import (
	"bytes"
	"time"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = time.RFC3339Nano
)

// Creates an '=' clause for a time.Time value.
func (c Constraint) EqualsTime(includeTime bool, right time.Time) Constraint {
	c.op = "="
	if includeTime {
		c.right = right.Format(DateTimeFormat)
	} else {
		c.right = right.Format(DateFormat)
	}
	return c
}

// Creates a '<>' clause for a time.Time value.
func (c Constraint) NotEqualsTime(includeTime bool, right time.Time) Constraint {
	c.op = "<>"
	if includeTime {
		c.right = right.Format(DateTimeFormat)
	} else {
		c.right = right.Format(DateFormat)
	}
	return c
}

// Creates a '>' clause for a time.Time value.
func (c Constraint) GreaterTime(includeTime bool, right time.Time) Constraint {
	c.op = ">"
	if includeTime {
		c.right = right.Format(DateTimeFormat)
	} else {
		c.right = right.Format(DateFormat)
	}
	return c
}

// Creates a '>=' clause for a time.Time value.
func (c Constraint) GreaterEqualsTime(includeTime bool, right time.Time) Constraint {
	c.op = ">="
	if includeTime {
		c.right = right.Format(DateTimeFormat)
	} else {
		c.right = right.Format(DateFormat)
	}
	return c
}

// Creates a '<' clause for a time.Time value.
func (c Constraint) LessTime(includeTime bool, right time.Time) Constraint {
	c.op = "<"
	if includeTime {
		c.right = right.Format(DateTimeFormat)
	} else {
		c.right = right.Format(DateFormat)
	}
	return c
}

// Creates a '<=' clause for a time.Time value.
func (c Constraint) LessEqualsTime(includeTime bool, right time.Time) Constraint {
	c.op = "<="
	if includeTime {
		c.right = right.Format(DateTimeFormat)
	} else {
		c.right = right.Format(DateFormat)
	}
	return c
}

// Creates an IN clause for a time.Time value.
func (c Constraint) InTime(includeTime bool, in ...time.Time) Constraint {
	c.op = " IN "
	buf := bytes.NewBufferString("(")
	for i, s := range in {
		if includeTime {
			buf.WriteString(s.Format(DateTimeFormat))
		} else {
			buf.WriteString(s.Format(DateFormat))
		}
		if i < len(in)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	c.right = buf.String()
	return c
}

// Creates a NOT IN clause for a time.Time value.
func (c Constraint) NotInTime(includeTime bool, in ...time.Time) Constraint {
	c.op = " NOT IN "
	buf := bytes.NewBufferString("(")
	for i, s := range in {
		if includeTime {
			buf.WriteString(s.Format(DateTimeFormat))
		} else {
			buf.WriteString(s.Format(DateFormat))
		}
		if i < len(in)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	c.right = buf.String()
	return c
}
