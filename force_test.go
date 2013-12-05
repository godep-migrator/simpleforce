package simpleforce

import (
	"fmt"
	"os"
	"testing"
)

var (
	force Force
)

func init() {
	session := os.Getenv("FORCE_SESSION")
	url := os.Getenv("FORCE_URL")
	force = New(session, url)
}

type Account struct {
	Name string
}

type Contact struct {
	FirstName string
	LastName  string
	Account   *Account
}

func Example() {
	var cs []Contact
	q := force.NewQuery(&cs)
	q.AddConstraint(NewConstraint("Name").EqualsString("Jake Basile"))
	q.Run()
	for _, c := range cs {
		fmt.Printf("%v %v Is From %v", c.FirstName, c.LastName, c.Account.Name)
	}
	// Output:
	// Jake Basile Is From Mutual Mobile
}

func TestQueryCreation(t *testing.T) {
	var as []Account
	q := force.NewQuery(&as)
	t.Log(q)
}

func TestSimpleConstraintCreation(t *testing.T) {
	c := NewConstraint("FirstName").EqualsString("Jake")
	t.Log(c)
}

func TestComplexConstraintCreation(t *testing.T) {
	c1 := NewConstraint("FirstName").EqualsString("Jake")
	c2 := NewConstraint("LastName").NotEqualsString("Basile")
	ca := NewConstraint(c1).And(c2)
	c3 := NewConstraint("Account.Name").EqualsString("Mutual Mobile")
	co := NewConstraint(ca).Or(c3)
	t.Log(co)
	t.Log(co.Collapse())
}

func TestSimpleQueryGeneration(t *testing.T) {
	var cs []Contact
	q := force.NewQuery(&cs)
	q.AddConstraint(NewConstraint("FirstName").EqualsString("Jake"))
	q.AddConstraint(NewConstraint("LastName").EqualsString("Basile"))
	q.AddConstraint(NewConstraint("Account.Name").EqualsString("Mutual Mobile"))
	t.Log(q.Generate())
}

func TestSimpleQueryRun(t *testing.T) {
	var cs []Contact
	q := force.NewQuery(&cs)
	q.AddConstraint(NewConstraint("FirstName").EqualsString("Jake"))
	t.Log(q.Generate())
	q.Run()
	for _, c := range cs {
		t.Log(c.FirstName, c.LastName, c.Account.Name)
	}
}

func TestRawQuery(t *testing.T) {
	var cs []Contact
	force.RunRawQuery("SELECT FirstName FROM Contact LIMIT 1", &cs)
	for _, c := range cs {
		t.Log(c.FirstName, c.LastName, c.Account.Name)
	}
}
