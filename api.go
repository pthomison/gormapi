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

func Index(c dbutils.DBClient, t interface{}) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		ids := dbutils.Query[IDView](c.DB().Model(&t).Select("ID"))

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

// func All(w http.ResponseWriter, r *http.Request) {
// 	objs := []APIObject{}

// 	objs = utils.SelectAll[APIObject](dbClient, []string{})

// 	json, err := json.Marshal(objs)
// 	utils.Check(err)

// 	_, err = w.Write(json)
// 	utils.Check(err)
// }

// func ID(w http.ResponseWriter, r *http.Request) {
// 	id_str := strings.TrimPrefix(r.URL.Path, "/id/")
// 	id, err := strconv.Atoi(id_str)
// 	utils.Check(err)

// 	fmt.Println(id)

// 	objs := []APIObject{}
// 	objs = utils.SelectWhere[APIObject](dbClient, []string{}, "id = ?", id)

// 	json, err := json.Marshal(objs)
// 	utils.Check(err)

// 	_, err = w.Write(json)
// 	utils.Check(err)
// }
