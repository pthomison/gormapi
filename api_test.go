package gormapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pthomison/dbutils"
	"github.com/pthomison/dbutils/sqlite"
	"github.com/pthomison/errcheck"
)

const (
	testObjCount = 10
)

type gormObject struct {
	ID uint

	StringData  string  `json:"string_data,omitempty"`
	IntegerData int     `json:"integer_data,omitempty"`
	FloatData   float64 `json:"float_data,omitempty"`
	BooleanData bool    `json:"boolean_data,omitempty"`
}

func createTestData() *sqlite.SQLiteClient {

	client := sqlite.New(":memory:")
	dbutils.Migrate(client, &gormObject{})

	for i := 0; i < testObjCount; i++ {
		dbutils.Create(client, []gormObject{
			gormObject{
				StringData:  fmt.Sprintf("stringdata %v", i),
				IntegerData: i,
				FloatData:   float64(i),
				BooleanData: i%2 == 0,
			},
		})
	}

	return client
}

func TestIndex(t *testing.T) {

	client := createTestData()

	objs := dbutils.SelectAll[gormObject](client)

	ts := httptest.NewServer(http.HandlerFunc(Index(client, &gormObject{})))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	errcheck.CheckTest(err, t)

	indexBytes, err := io.ReadAll(res.Body)
	res.Body.Close()
	errcheck.CheckTest(err, t)

	index := make([]IDView, testObjCount)

	err = json.Unmarshal(indexBytes, &index)
	errcheck.CheckTest(err, t)

	if len(objs) != len(index) {
		errcheck.CheckTest(errors.New("mismatch in index(func) data lengths"), t)
	}

	checkData := make([]IDView, testObjCount)

	for i := 0; i < testObjCount; i++ {
		checkData[i].ID = uint(i + 1)
	}

	checkJson, err := json.Marshal(checkData)
	errcheck.CheckTest(err, t)

	if string(checkJson) != string(indexBytes) {
		fmt.Println(string(checkJson))
		fmt.Println(string(indexBytes))

		errcheck.CheckTest(errors.New("mismatch in index(func) return data"), t)
	}

}

func TestAll(t *testing.T) {

	client := createTestData()

	objs := dbutils.SelectAll[gormObject](client)
	objsJSON, err := json.Marshal(objs)
	errcheck.CheckTest(err, t)

	ts := httptest.NewServer(http.HandlerFunc(All[gormObject](client)))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	errcheck.CheckTest(err, t)

	resBytes, err := io.ReadAll(res.Body)
	res.Body.Close()
	errcheck.CheckTest(err, t)

	resultObjs := make([]gormObject, testObjCount)
	err = json.Unmarshal(resBytes, &resultObjs)
	errcheck.CheckTest(err, t)

	if len(objs) != len(resultObjs) {
		errcheck.CheckTest(errors.New("mismatch in all(func) data lengths"), t)
	}

	if bytes.Compare(resBytes, objsJSON) != 0 {
		fmt.Println(string(resBytes))
		fmt.Println(string(objsJSON))

		errcheck.CheckTest(errors.New("mismatch in all(func) return data"), t)
	}

}

func TestID(t *testing.T) {

	testIndex := 4

	client := createTestData()

	objs := dbutils.SelectAll[gormObject](client)
	obj := objs[testIndex-1]
	objJSON, err := json.Marshal(obj)
	errcheck.CheckTest(err, t)

	ts := httptest.NewServer(http.HandlerFunc(ID[gormObject](client)))
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%v/id/%v", ts.URL, testIndex))
	errcheck.CheckTest(err, t)

	resBytes, err := io.ReadAll(res.Body)
	res.Body.Close()
	errcheck.CheckTest(err, t)

	if bytes.Compare(resBytes, objJSON) != 0 {
		fmt.Println(string(resBytes))
		fmt.Println(string(objJSON))

		errcheck.CheckTest(errors.New("mismatch in id(func) return data"), t)
	}

}
