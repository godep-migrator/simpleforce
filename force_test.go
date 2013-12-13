package simpleforce_test

import (
	"fmt"
	"github.com/jakebasile/simpleforce"
	"os"
	"testing"
	"time"
)

var (
	force simpleforce.Force
)

func init() {
	if len(os.Getenv("SF_PASSWORD")) > 0 {
		un := os.Getenv("SF_USERNAME")
		pw := os.Getenv("SF_PASSWORD")
		secret := os.Getenv("SF_CLIENT_SECRET")
		id := os.Getenv("SF_CLIENT_ID")
		login := os.Getenv("SF_LOGIN_URL")
		nforce, err := simpleforce.NewWithCredentials(login, id, secret, un, pw)
		if err != nil {
			panic(err)
		}
		force = nforce
	} else {
		session := os.Getenv("FORCE_SESSION")
		url := os.Getenv("FORCE_URL")
		force = simpleforce.New(session, url)
	}
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
		force.Query(`
			SELECT 
				Account.Name,
				FirstName,
				LastName,
				Name
			FROM Contact
			WHERE 
				Account.Name <> '' AND
				FirstName <> '' AND
				LastName <> ''
			LIMIT 1000`, &cs)
	}
}

func ExampleForce_Query() {
	var cs []Contact
	force.Query("SELECT Name FROM Contact WHERE FirstName='Jake' AND LastName='Basile'", &cs)
	for _, c := range cs {
		fmt.Println(c.Name)
	}
}

func TestDate(t *testing.T) {
	type Contact struct {
		Birthdate time.Time
	}

	var cs []Contact
	force.Query("SELECT Birthdate FROM Contact WHERE Birthdate <> NULL LIMIT 1", &cs)
	t.Log(cs)
	for _, c := range cs {
		if c.Birthdate.IsZero() {
			t.Fail()
		}
	}
}
func TestRawQuery(t *testing.T) {
	var cs []Contact
	force.Query("SELECT FirstName FROM Contact WHERE FirstName<>'' LIMIT 1", &cs)
	for _, c := range cs {
		t.Log(c)
		if c.FirstName == "" || c.LastName != "" || c.Account == nil || c.Account.Name != "" {
			t.Fail()
		}
	}
}

func TestChildObjects(t *testing.T) {
	type Contact struct {
		Name string
	}
	type Account struct {
		Name     string
		Contacts []Contact
	}

	var as []Account
	force.Query("SELECT Name, (SELECT Name FROM Contacts) FROM Account WHERE Name='Mutual Mobile' LIMIT 1", &as)
	t.Log(as)
	for _, a := range as {
		if len(a.Contacts) == 0 {
			t.Fail()
		}
		for _, c := range a.Contacts {
			if c.Name == "" {
				t.Fail()
			}
		}
	}
}
