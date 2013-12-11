package simpleforce

// Part of a WHERE clause for a SOQL query.
type Constraint struct {
	left  interface{}
	op    string
	right interface{}
}

func NewConstraint(left interface{}) Constraint {
	return Constraint{
		left,
		"",
		nil,
	}
}

// Turns a Constraint into a WHERE clause.
func (c *Constraint) Collapse() string {
	var leftString string
	switch c.left.(type) {
	case string:
		leftString = c.left.(string)
	case Constraint:
		leftCon := c.left.(Constraint)
		leftString = leftCon.Collapse()
	}
	var rightString string
	switch c.right.(type) {
	case string:
		rightString = c.right.(string)
	case Constraint:
		rightCon := c.right.(Constraint)
		rightString = rightCon.Collapse()
	}
	// Handles root constraint case.
	if len(leftString) > 2 {
		return "(" + leftString + c.op + rightString + ")"
	} else {
		return rightString
	}
}

// Combines two Constraints with AND.
func (c Constraint) And(right Constraint) Constraint {
	c.op = " AND "
	c.right = right
	return c
}

// Combines two Constraints with OR.
func (c Constraint) Or(right Constraint) Constraint {
	c.op = " OR "
	c.right = right
	return c
}

func (c Constraint) EqualsNull() Constraint {
	c.op = "="
	c.right = "NULL"
	return c
}

func (c Constraint) NotEqualsNull() Constraint {
	c.op = "<>"
	c.right = "NULL"
	return c
}
