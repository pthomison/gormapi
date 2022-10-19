package gormapi

import (
	"encoding/json"
	"net/http"

	"github.com/pthomison/dbutils"
	"github.com/pthomison/errcheck"
)

// , w http.ResponseWriter, r *http.Request

func Index[T any](c dbutils.DBClient) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		ids := dbutils.Query[T](c.DB().Select("id"))

		json, err := json.Marshal(ids)
		errcheck.Check(err)

		_, err = w.Write(json)
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
