package simpleforce_test

import (
	"fmt"
	"github.com/jakebasile/simpleforce"
	"os"
	"testing"
)

var (
	force simpleforce.Force
)

func init() {
	session := os.Getenv("FORCE_SESSION")
	url := os.Getenv("FORCE_URL")
	force = simpleforce.New(session, url)
}

type Account struct {
	Name string
}

type Contact struct {
	FirstName string
	LastName  string
	Name      string
	Account   *Account
}

func BenchmarkQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var cs []Contact
		q := force.NewQuery(&cs)
		q.AddConstraint(simpleforce.NewConstraint("Account.Name").NotEqualsString(""))
		q.AddConstraint(simpleforce.NewConstraint("FirstName").NotEqualsString(""))
		q.AddConstraint(simpleforce.NewConstraint("LastName").NotEqualsString(""))
		q.Limit(1000)
		q.Run()
	}
}

func Example() {
	var cs []Contact
	q := force.NewQuery(&cs)
	q.AddConstraint(simpleforce.NewConstraint("Name").EqualsString("Jake Basile"))
	q.Run()
	for _, c := range cs {
		fmt.Printf("%v %v Is From %v", c.FirstName, c.LastName, c.Account.Name)
	}
}

func ExampleForce_RunRawQuery() {
	var cs []Contact
	force.RunRawQuery("SELECT Name FROM Contact WHERE FirstName='Jake' AND LastName='Basile'", &cs)
	for _, c := range cs {
		fmt.Println(c.Name)
	}
}

func ExampleConstraint() {
	c1 := simpleforce.NewConstraint("FirstName").EqualsString("Jake")
	c2 := simpleforce.NewConstraint("LastName").NotEqualsString("Basile")
	c3 := simpleforce.NewConstraint(c1).Or(c2)
	fmt.Println(c1.Collapse())
	fmt.Println(c2.Collapse())
	fmt.Println(c3.Collapse())
}

func TestQueryCreation(t *testing.T) {
	var as []Account
	q := force.NewQuery(&as)
	t.Log(q)
}

func TestSimpleConstraintCreation(t *testing.T) {
	c := simpleforce.NewConstraint("FirstName").EqualsString("Jake")
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(FirstName='Jake')" {
		t.Fail()
	}
}

func TestComplexConstraintCreation(t *testing.T) {
	c1 := simpleforce.NewConstraint("FirstName").EqualsString("Jake")
	c2 := simpleforce.NewConstraint("LastName").NotEqualsString("Basile")
	ca := simpleforce.NewConstraint(c1).And(c2)
	c3 := simpleforce.NewConstraint("Account.Name").EqualsString("Mutual Mobile")
	co := simpleforce.NewConstraint(ca).Or(c3)
	t.Log(co)
	t.Log(co.Collapse())
	if co.Collapse() != "(((FirstName='Jake') AND (LastName<>'Basile')) OR (Account.Name='Mutual Mobile'))" {
		t.Fail()
	}
}

func TestStringInConstraint(t *testing.T) {
	c := simpleforce.NewConstraint("FirstName").InString("Jake", "Kyle")
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(FirstName IN ('Jake','Kyle'))" {
		t.Fail()
	}
}

func TestStringNotInConstraint(t *testing.T) {
	c := simpleforce.NewConstraint("FirstName").NotInString("Jake", "Kyle")
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(FirstName NOT IN ('Jake','Kyle'))" {
		t.Fail()
	}
}

func TestStringLikeConstraint(t *testing.T) {
	c := simpleforce.NewConstraint("FirstName").LikeString("%K%")
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(FirstName LIKE '%K%')" {
		t.Fail()
	}
}

func TestIntInConstraint(t *testing.T) {
	c := simpleforce.NewConstraint("Dummy__c").InInt(1, 2, 3, 4, 5)
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(Dummy__c IN (1,2,3,4,5))" {
		t.Fail()
	}
}

func TestIntNotInConstraint(t *testing.T) {
	c := simpleforce.NewConstraint("Dummy__c").NotInInt(1, 2, 3, 4, 5)
	t.Log(c)
	t.Log(c.Collapse())
	if c.Collapse() != "(Dummy__c NOT IN (1,2,3,4,5))" {
		t.Fail()
	}
}

func TestSimpleQueryGeneration(t *testing.T) {
	var cs []Contact
	q := force.NewQuery(&cs)
	q.AddConstraint(simpleforce.NewConstraint("FirstName").EqualsString("Jake"))
	q.AddConstraint(simpleforce.NewConstraint("LastName").EqualsString("Basile"))
	q.AddConstraint(simpleforce.NewConstraint("Account.Name").EqualsString("Mutual Mobile"))
	t.Log(q.Generate())
	if q.Generate() != "SELECT FirstName,LastName,Name,Account.Name FROM Contact WHERE (FirstName='Jake') AND (LastName='Basile') AND (Account.Name='Mutual Mobile') LIMIT 10" {
		t.Fail()
	}
}

func TestSimpleQueryRun(t *testing.T) {
	var cs []Contact
	q := force.NewQuery(&cs)
	q.AddConstraint(simpleforce.NewConstraint("Account.Name").NotEqualsString(""))
	q.AddConstraint(simpleforce.NewConstraint("FirstName").NotEqualsString(""))
	q.AddConstraint(simpleforce.NewConstraint("LastName").NotEqualsString(""))
	q.Limit(1000)
	t.Log(q.Generate())
	q.Run()
	if len(cs) != 1000 {
		t.Fail()
	}
	for _, c := range cs {
		t.Log(c.FirstName, c.LastName, c.Account.Name)
		if c.FirstName == "" || c.LastName == "" || c.Account == nil || c.Account.Name == "" {
			t.Fail()
		}
	}
}

func TestRawQuery(t *testing.T) {
	var cs []Contact
	force.RunRawQuery("SELECT FirstName FROM Contact WHERE FirstName<>'' LIMIT 1", &cs)
	for _, c := range cs {
		t.Log(c)
		if c.FirstName == "" || c.LastName != "" || c.Account == nil || c.Account.Name != "" {
			t.Fail()
		}
	}
}
