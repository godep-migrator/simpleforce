// simpleforce is a dead simple wrapper around the Force.com REST API.
package simpleforce

import (
	"bytes"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type Force struct {
	session string
	url     string
}

// Returns a new Force object with the given login credentials.
func New(session, url string) Force {
	return Force{
		session,
		url,
	}
}

func (f Force) NewQuery(dest interface{}) Query {
	return Query{
		f,
		dest,
		make([]Constraint, 0, 0),
	}
}

func (f Force) authorizeRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", "Bearer "+f.session)
	return r, nil
}

func (f Force) RunRawQuery(query string, dest interface{}) error {
	vals := url.Values{}
	vals.Set("q", query)
	url := f.url + "/query?" + vals.Encode()
	req, err := f.authorizeRequest("GET", url, bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respJson, err := simplejson.NewJson(respBytes)
	if err != nil {
		return err
	}
	err = Unmarshal(respJson, dest)
	return err
}

func Unmarshal(source *simplejson.Json, dest interface{}) error {
	sliceValPtr := reflect.ValueOf(dest)
	sliceVal := sliceValPtr.Elem()
	elemType := reflect.TypeOf(dest).Elem().Elem()
	for i := 0; i < source.Get("totalSize").MustInt(); i++ {
		v := source.Get("records").GetIndex(i)
		val, err := unmarshalIndividualObject(v, elemType)
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, val))
	}
	return nil
}

func unmarshalIndividualObject(source *simplejson.Json, valType reflect.Type) (reflect.Value, error) {
	valPtr := reflect.New(valType)
	val := reflect.Indirect(valPtr)
	for f := 0; f < valType.NumField(); f++ {
		// find the field
		field := val.Field(f)
		switch field.Kind() {
		case reflect.String:
			strVal := source.Get(valType.Field(f).Name).MustString()
			field.SetString(strVal)
		case reflect.Ptr:
			objJson := source.Get(valType.Field(f).Name)
			objType := valType.Field(f).Type.Elem()
			objVal, err := unmarshalIndividualObject(objJson, objType)
			if err != nil {
				return val, err
			}
			field.Set(objVal.Addr())
		}
	}
	return val, nil
}

// Force handles authentication, upserting, saving, and deleting. Query handled by Query object.

type Query struct {
	force       Force
	dest        interface{}
	constraints []Constraint
}

func (q *Query) AddConstraint(c Constraint) {
	q.constraints = append(q.constraints, c)
}

func (q *Query) Run() error {
	err := q.force.RunRawQuery(q.Generate(), q.dest)
	if err != nil {
		return err
	}
	return nil
}

func (q *Query) Generate() string {
	sel := q.generateSelect()
	table := reflect.TypeOf(q.dest).Elem().Elem().Name()
	where := q.generateWhere()
	return fmt.Sprintf("SELECT %v FROM %v WHERE %v", sel, table, where)
}

func (q *Query) generateSelect() string {
	return genSelectForType(reflect.TypeOf(q.dest).Elem().Elem(), "")
}

func genSelectForType(t reflect.Type, path string) string {
	buf := bytes.NewBufferString("")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Ptr {
			buf.WriteString(genSelectForType(field.Type.Elem(), field.Name))
			buf.WriteString(",")
		} else if field.Type.Kind() == reflect.Slice {
			// wat do
		} else {
			if len(path) > 0 {
				buf.WriteString(path + "." + field.Name)
			} else {
				buf.WriteString(field.Name)
			}
			buf.WriteString(",")
		}
	}
	s := buf.String()
	// drop last comma.
	return s[:len(s)-1]
}

func (q *Query) generateWhere() string {
	buf := bytes.NewBufferString("")
	for i, c := range q.constraints {
		buf.WriteString(c.Collapse())
		if i < len(q.constraints)-1 {
			buf.WriteString(" AND ")
		}
	}
	return buf.String()
}

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

func (c Constraint) And(right Constraint) Constraint {
	c.op = " AND "
	c.right = right
	return c
}

func (c Constraint) Or(right Constraint) Constraint {
	c.op = " OR "
	c.right = right
	return c
}

func (c Constraint) EqualsString(right string) Constraint {
	c.op = "="
	c.right = "'" + right + "'"
	return c
}

func (c Constraint) NotEqualsString(right string) Constraint {
	c.op = "<>"
	c.right = "'" + right + "'"
	return c
}
