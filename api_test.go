package gormapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pthomison/dbutils"
	"github.com/pthomison/dbutils/sqlite"
	"github.com/pthomison/errcheck"
	"gorm.io/gorm"
)

type gormObject struct {
	gorm.Model

	StringData  string  `json:"string_data,omitempty"`
	IntegerData int     `json:"integer_data,omitempty"`
	FloatData   float64 `json:"float_data,omitempty"`
	BooleanData bool    `json:"boolean_data,omitempty"`
}

func createTestData() *sqlite.SQLiteClient {

	client := sqlite.New(":memory:")
	dbutils.Migrate(client, &gormObject{})

	for i := 0; i < 10; i++ {
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

	// objs := dbutils.SelectAll[gormObject](client)

	// fmt.Printf("%+v\n", objs)

	ts := httptest.NewServer(http.HandlerFunc(Index[gormObject](client)))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	errcheck.CheckTest(err, t)

	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)

}
