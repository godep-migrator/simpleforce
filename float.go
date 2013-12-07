package simpleforce

import (
	"bytes"
	"fmt"
)

// Creates an '=' clause for a int value.
func (c Constraint) EqualsFloat(right float64) Constraint {
	c.op = "="
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '<>' clause for a int value.
func (c Constraint) NotEqualsFloat(right float64) Constraint {
	c.op = "<>"
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '>' clause for a int value.
func (c Constraint) GreaterFloat(right float64) Constraint {
	c.op = ">"
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '>=' clause for a int value.
func (c Constraint) GreaterEqualsFloat(right float64) Constraint {
	c.op = ">="
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '<' clause for a int value.
func (c Constraint) LessFloat(right float64) Constraint {
	c.op = "<"
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '<=' clause for a int value.
func (c Constraint) LessEqualsFloat(right float64) Constraint {
	c.op = "<="
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates an IN clause for a int value.
func (c Constraint) InFloat(in ...float64) Constraint {
	c.op = " IN "
	buf := bytes.NewBufferString("(")
	for i, s := range in {
		buf.WriteString(fmt.Sprintf("%v", s))
		if i < len(in)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	c.right = buf.String()
	return c
}

// Creates a NOT IN clause for a int value.
func (c Constraint) NotInFloat(in ...float64) Constraint {
	c.op = " NOT IN "
	buf := bytes.NewBufferString("(")
	for i, s := range in {
		buf.WriteString(fmt.Sprintf("%v", s))
		if i < len(in)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	c.right = buf.String()
	return c
}
