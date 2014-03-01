package query_test

import (
	"fmt"
	"bitbucket.org/jakebasile/simpleforce"
	"bitbucket.org/jakebasile/simpleforce/query"
	"testing"
)

type Account struct {
	Name string
}

type Contact struct {
	FirstName string
	LastName  string
	Name      string
	Account   *Account
}

func ExampleConstraint() {
	c1 := query.NewConstraint("FirstName").EqualsString("Jake")
	c2 := query.NewConstraint("LastName").NotEqualsString("Basile")
	c3 := query.NewConstraint(c1).Or(c2)
	fmt.Println(c1.Collapse())
	fmt.Println(c2.Collapse())
	fmt.Println(c3.Collapse())
	// Output:
	// (FirstName='Jake')
	// (LastName<>'Basile')
	// ((FirstName='Jake') OR (LastName<>'Basile'))
}

func TestSimpleConstraintCreation(t *testing.T) {
	c := query.NewConstraint("FirstName").EqualsString("Jake")
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(FirstName='Jake')" {
		t.Fail()
	}
}

func TestComplexConstraintCreation(t *testing.T) {
	c1 := query.NewConstraint("FirstName").EqualsString("Jake")
	c2 := query.NewConstraint("LastName").NotEqualsString("Basile")
	ca := query.NewConstraint(c1).And(c2)
	c3 := query.NewConstraint("Account.Name").EqualsString("Mutual Mobile")
	co := query.NewConstraint(ca).Or(c3)
	t.Log(co)
	t.Log(co.Collapse())
	if co.Collapse() != "(((FirstName='Jake') AND (LastName<>'Basile')) OR (Account.Name='Mutual Mobile'))" {
		t.Fail()
	}
}

func TestStringInConstraint(t *testing.T) {
	c := query.NewConstraint("FirstName").InString("Jake", "Kyle")
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(FirstName IN ('Jake','Kyle'))" {
		t.Fail()
	}
}

func TestStringNotInConstraint(t *testing.T) {
	c := query.NewConstraint("FirstName").NotInString("Jake", "Kyle")
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(FirstName NOT IN ('Jake','Kyle'))" {
		t.Fail()
	}
}

func TestStringLikeConstraint(t *testing.T) {
	c := query.NewConstraint("FirstName").LikeString("%K%")
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(FirstName LIKE '%K%')" {
		t.Fail()
	}
}

func TestIntInConstraint(t *testing.T) {
	c := query.NewConstraint("Dummy__c").InInt(1, 2, 3, 4, 5)
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(Dummy__c IN (1,2,3,4,5))" {
		t.Fail()
	}
}

func TestIntNotInConstraint(t *testing.T) {
	c := query.NewConstraint("Dummy__c").NotInInt(1, 2, 3, 4, 5)
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(Dummy__c NOT IN (1,2,3,4,5))" {
		t.Fail()
	}
}

func TestSimpleQueryGeneration(t *testing.T) {
	var cs []Contact
	q := query.New(simpleforce.Force{}, &cs)
	q.AddConstraint(query.NewConstraint("FirstName").EqualsString("Jake"))
	q.AddConstraint(query.NewConstraint("LastName").EqualsString("Basile"))
	q.AddConstraint(query.NewConstraint("Account.Name").EqualsString("Mutual Mobile"))
	t.Log(q.Generate())
	if q.Generate() != "SELECT FirstName,LastName,Name,Account.Name FROM Contact WHERE (FirstName='Jake') AND (LastName='Basile') AND (Account.Name='Mutual Mobile') LIMIT 10" {
		t.Fail()
	}
}
