package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (dbCfg *DbConfig) insertUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		dataParam := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&dataParam); err != nil {
			responseWithError(w, 403, err.Error())
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dataParam.Password), bcrypt.DefaultCost)

		if err != nil {
			responseWithError(w, 303, "somthing went wrong while storing password")
			return
		}

		useUid := uuid.New()
		_, err = insertUser(dbCfg.sqlDb, &User{
			Uid:      useUid,
			Name:     dataParam.Name,
			Password: string(hashedPassword),
		})

		if err != nil {
			responseWithError(w, 300, err.Error())
		}

		responseWithJson(w, 200, struct {
			UserId string `json:"userId"`
		}{
			UserId: useUid.String(),
		})
	}

}

func (dbCfg *DbConfig) getUserDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		dataParam := struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		}{}

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&dataParam)

		if err != nil {
			responseWithError(w, 301, err.Error())
			return
		}

		id, name, pas, err := getUserDetails(dbCfg.sqlDb, dataParam.UserName)

		if err != nil {
			responseWithError(w, 301, err.Error())
			return
		}

		same := bcrypt.CompareHashAndPassword([]byte(pas), []byte(dataParam.Password))

		if same != nil {
			responseWithError(w, 203, "Password is not same")
			return
		} else {
			responseWithJson(w, 202, struct {
				Uid      uuid.UUID `json:"uid"`
				Username string    `json:"usename"`
			}{
				Uid:      uuid.MustParse(id),
				Username: name,
			})
		}

	}
}

// func (dbCfg *DbConfig) getAllData(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		res, err := getUsers(dbCfg.sqlDb)

// 		if err != nil {
// 			responseWithError(w, 305, err.Error())
// 			return
// 		}

// 		responseWithJson(w, 201, res)

// 	}
// }
