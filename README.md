# simpleforce

Force.com REST API made simple for Go.

## Installation

```bash
go get bitbucket.org/jakebasile/simpleforce
```

## Usage

For basic usage examples, please check the included tests. But, simply:

```go
import (
    "fmt"
    "bitbucket.org/jakebasile/simpleforce"
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
    f.Query(`
        SELECT
            FirstName,
            LastName,
            Account.Name
        FROM Contact
        WHERE
            FirstName='Jake' OR
            Account.Name='Mutual Mobile'`, &cs)
    for _, c := range cs {
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

