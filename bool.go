package simpleforce

// Creates an '=' clause for a bool value.
func (c Constraint) EqualsBool(right bool) Constraint {
	c.op = "="
	if right {
		c.right = "TRUE"
	} else {
		c.right = "FALSE"
	}
	return c
}

// Creates a '<>' clause for a bool value.
func (c Constraint) NotEqualsBool(right bool) Constraint {
	c.op = "<>"
	if right {
		c.right = "TRUE"
	} else {
		c.right = "FALSE"
	}
	return c
}
