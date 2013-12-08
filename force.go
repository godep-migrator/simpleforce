/*
Package simpleforce is a dead simple wrapper around the Force.com REST API.

It allows you to query for Force.com objects by using idiomatic Go constructs, or you can short
circuit the query engine and qrite your own SOQL. In either case, data is returned to you via
structs of your own creation, allowing you full control over what data is returned.
*/
package simpleforce

import (
	"bytes"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type Force struct {
	session string
	url     string
}

// Returns a new Force object with the given login credentials. This object is the main
// point of entry for all your Force.com needs.
func New(session, url string) Force {
	return Force{
		session,
		url,
	}
}

func NewWithCredentials(loginUrl, consumerKey, consumerSecret, username, password string) (Force, error) {
	resp, err := http.PostForm(loginUrl+"/services/oauth2/token", url.Values{
		"grant_type":    {"password"},
		"client_id":     {consumerKey},
		"client_secret": {consumerSecret},
		"username":      {username},
		"password":      {password},
	})
	if err != nil {
		return Force{}, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Force{}, err
	}
	respJson, err := simplejson.NewJson(respBytes)
	if err != nil {
		return Force{}, err
	}
	session := respJson.Get("access_token").MustString()
	url := respJson.Get("instance_url").MustString() + "/services/data/v27.0"
	return New(session, url), err
}

// Creates a new query for you to customize. When executed, this query will fill the given destination
// slice with the results of the query.
func (f Force) NewQuery(dest interface{}) Query {
	return Query{
		f,
		dest,
		make([]Constraint, 0, 0),
		10,
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

// Run a raw SOQL query string. This will fill the given destination slice with the results of your query.
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
	err = unmarshal(respJson, dest)
	return err
}

func unmarshal(source *simplejson.Json, dest interface{}) error {
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
		field := val.Field(f)
		switch field.Kind() {
		case reflect.Bool:
			boolVal := source.Get(valType.Field(f).Name).MustBool()
			field.SetBool(boolVal)
		case reflect.Int:
			intVal := source.Get(valType.Field(f).Name).MustInt64()
			field.SetInt(intVal)
		case reflect.Int64:
			intVal := source.Get(valType.Field(f).Name).MustInt64()
			field.SetInt(intVal)
		case reflect.Float32:
			floatVal := source.Get(valType.Field(f).Name).MustFloat64()
			field.SetFloat(floatVal)
		case reflect.Float64:
			floatVal := source.Get(valType.Field(f).Name).MustFloat64()
			field.SetFloat(floatVal)
		case reflect.String:
			strVal := source.Get(valType.Field(f).Name).MustString()
			field.SetString(strVal)
		case reflect.Struct:
			strVal := source.Get(valType.Field(f).Name).MustString()
			if valType.Field(f).Type.Name() == "Time" {
				if t, err := time.Parse(DateTimeFormat, strVal); err == nil {
					// it's a datetime string, probably!
					field.Set(reflect.ValueOf(t))
				} else if t, err = time.Parse(DateFormat, strVal); err == nil {
					// nope, it's a date string!
					field.Set(reflect.ValueOf(t))
				} else {
					return val, err
				}
			}
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
