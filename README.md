# simpleforce

Force.com REST API made simple for Go.

## Installation

```bash
go get bitbucket.org/jakebasile/simpleforce
```

## Usage

For basic usage examples, please check the included tests. But, simply:

```go
package main

import (
    "bitbucket.org/jakebasile/simpleforce"
    "fmt"
)

type Account struct {
    Name     string
}

type Contact struct {
	Name      string
	FirstName string
	LastName  string
	Account   *Account
}

func main() {
    f := simpleforce.New("your_session_id", "your_force_url")
    var cs []Contact
    f.Query(`
        SELECT
            Name,
            FirstName,
            LastName,
            Account.Name
        FROM Contact
        WHERE
            FirstName='Jake' OR
            Account.Name='Mutual Mobile'`, &cs)
    for _, c := range cs {
        fmt.Printf("%v Works At %v\n", c.Name, c.Account.Name)
    }
}
```

would output:

    Jake Works At Mutual Mobile
    ...

And so on, based on what data is in your Force.com instance.

## Querygen

The `bitbucket.org/jakebasile/simpleforce/query` package lets you use Go constructs to query Salesforce. It is currently *unfnished but usable*. Beware circular references, as I haven't gotten those working yet.

Here's an example, equivalent to the previous example but using the `query` package.

```go
package main

import (
	"bitbucket.org/jakebasile/simpleforce"
	"bitbucket.org/jakebasile/simpleforce/query"
	"fmt"
)

type Account struct {
	Name string
}

type Contact struct {
	Name      string
	FirstName string
	LastName  string
	Account   *Account
}

func main() {
    f := simpleforce.New("your_session_id", "your_force_url")
	var cs []Contact
	q := query.New(f, &cs)
	q.AddConstraint(query.NewConstraint("FirstName").EqualsString("Jake"))
	q.AddConstraint(query.NewConstraint("Account.Name").EqualsString("Mutual Mobile"))
	q.Run()
	for _, c := range cs {
		fmt.Printf("%v Works At %v\n", c.Name, c.Account.Name)
	}
}
```

## Contributing

Any help would be greatly appreciated! Please **submit a new issue or comment on an existing one** before starting work on something, to make sure there's no overlap and that the new feature/bug fix is consistent.

Be sure to add your name and web address to the CONTRIBUTORS.txt file.

