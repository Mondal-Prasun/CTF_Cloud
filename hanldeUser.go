package main

import (
	"encoding/json"
	"net/http"
)

func (dbCfg *dbConfig) insertUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		dataParam := struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&dataParam); err != nil {
			responseWithError(w, 403, err.Error())
			return
		}

		res, err := insertUser(dbCfg.sqlDb, &User{
			Name: dataParam.Name,
			Age:  dataParam.Age,
		})

		if err != nil {
			responseWithError(w, 300, err.Error())
		}

		responseWithJson(w, 200, res)
	}

}

func (dbCfg *dbConfig) getAllData(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		res, err := getUsers(dbCfg.sqlDb)

		if err != nil {
			responseWithError(w, 305, err.Error())
			return
		}

		responseWithJson(w, 201, res)

	}
}
