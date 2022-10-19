package gormapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pthomison/dbutils"
	"github.com/pthomison/errcheck"
)

type IDView struct {
	ID uint
}

func Index[T any](c dbutils.DBClient) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		var t T

		ids := dbutils.Query[IDView](c.DB().Model(t).Select("ID"))

		jsonBytes, err := json.Marshal(ids)
		errcheck.Check(err)

		_, err = w.Write(jsonBytes)
		errcheck.Check(err)
	}
}

func All[T any](c dbutils.DBClient) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		objs := dbutils.SelectAll[T](c)

		jsonBytes, err := json.Marshal(objs)
		errcheck.Check(err)

		_, err = w.Write(jsonBytes)
		errcheck.Check(err)
	}
}

func ID[T any](c dbutils.DBClient) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		id_str := strings.TrimPrefix(r.URL.Path, "/id/")
		id, err := strconv.Atoi(id_str)
		errcheck.Check(err)

		obj := dbutils.Query[T](c.DB().Where("id = ?", id))[0]

		jsonBytes, err := json.Marshal(obj)
		errcheck.Check(err)

		_, err = w.Write(jsonBytes)
		errcheck.Check(err)
	}
}
