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
	"os"
	"reflect"
	"time"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = time.RFC3339Nano
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

func NewFromEnvironment() (Force, error) {
	return NewWithCredentials(os.Getenv("SF_LOGIN_URL"), os.Getenv("SF_CLIENT_ID"), os.Getenv("SF_CLIENT_SECRET"), os.Getenv("SF_USERNAME"), os.Getenv("SF_PASSWORD")+os.Getenv("SF_TOKEN"))
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
func (f Force) Query(query string, dest interface{}) error {
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

func (f Force) Create(interface{}) (interface{}, error) {

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
			if objJson != nil {
				objType := valType.Field(f).Type.Elem()
				objVal, err := unmarshalIndividualObject(objJson, objType)
				if err != nil {
					return val, err
				}
				field.Set(objVal.Addr())
			}
		case reflect.Slice:
			objJson := source.Get(valType.Field(f).Name).Get("records")
			length := source.Get(valType.Field(f).Name).Get("totalSize").MustInt()
			if objJson != nil {
				elemType := field.Type().Elem()
				objSlicePtr := reflect.New(field.Type())
				objSlice := reflect.Indirect(objSlicePtr)
				for i := 0; i < length; i++ {
					o := objJson.GetIndex(i)
					obj, err := unmarshalIndividualObject(o, elemType)
					if err != nil {
						return val, err
					}
					objSlice.Set(reflect.Append(objSlice, obj))
				}
				field.Set(objSlice)
			}
		}
	}
	return val, nil
}
