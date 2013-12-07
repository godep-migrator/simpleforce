# simpleforce

Force.com REST API made simple for Go.

## Installation

```bash
go get github.com/jakebasile/simpleforce
```

## Usage

For basic usage examples, please check the included tests. But, simply:

```go
import (
    "fmt"
    "github.com/jakebasile/simpleforce"
)

type Account struct {
    Name     string
}

type Contact struct {
    FirstName  string
    LastName   string
    Account    *Account
}

func main() {
    f := simpleforce.New("your_session_id", "your_force_url")
    var cs []Account
    q := f.NewQuery(a)
    c1 := simpleforce.NewConstraint("FirstName").EqualsString("Jake")
    c2 := simpleforce.NewConstraint("Account.Name").EqualsString("Mutual Mobile")
    q.AddConstraint(simpleforce.NewConstraint(c1).Or(c2))
    q.Run()
    for _, c := range as {
        fmt.Printf("%v Works At %v", c.FirstName, c.Account.Name)
    }
}
```

would output:

    Jake Works At Mutual Mobile
    ...

And so on, based on what data is in your Force.com instance.

## Contributing

Please **submit a new issue or comment on an existing one** before starting work on something, to make sure there's no overlap and that the new feature/bug fix is consistent.

Be sure to add your name and web address to the CONTRIBUTORS.txt file.

