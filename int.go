package simpleforce

import (
	"bytes"
	"fmt"
)

// Creates an '=' clause for a int value.
func (c Constraint) EqualsInt(right int) Constraint {
	c.op = "="
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '<>' clause for a int value.
func (c Constraint) NotEqualsInt(right int) Constraint {
	c.op = "<>"
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '>' clause for a int value.
func (c Constraint) GreaterInt(right int) Constraint {
	c.op = ">"
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '>=' clause for a int value.
func (c Constraint) GreaterEqualsInt(right int) Constraint {
	c.op = ">="
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '<' clause for a int value.
func (c Constraint) LessInt(right int) Constraint {
	c.op = "<"
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates a '<=' clause for a int value.
func (c Constraint) LessEqualsInt(right int) Constraint {
	c.op = "<="
	c.right = fmt.Sprintf("%v", right)
	return c
}

// Creates an IN clause for a int value.
func (c Constraint) InInt(in ...int) Constraint {
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
func (c Constraint) NotInInt(in ...int) Constraint {
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
