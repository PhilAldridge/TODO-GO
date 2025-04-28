package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PhilAldridge/TODO-GO/lib"
	"github.com/PhilAldridge/TODO-GO/router"
)

var (
	url      string
	username string
	password string
	jwt      string
)

func login() {
	fmt.Printf("username: %s, password: %s\n", username, password)
	if username != "" {
		url = fmt.Sprintf("%s%s/TodosV2",lib.BaseUrl, lib.PortNo)
		body, _ := json.Marshal(router.UserPutBody{
			Username: username,
			Password: password,
		})

		usersUrl:= fmt.Sprintf("%s%s/Users",lib.BaseUrl, lib.PortNo)
		res := sendAndReceive(http.MethodPost, body, usersUrl)
		if len(res) > 0 {
			jwt = string(res)
		} else {
			fmt.Println("Login failed")
			os.Exit(1)
		}
	} else {
		url = fmt.Sprintf("%s%s/Todos",lib.BaseUrl, lib.PortNo)
	}
}

func sendAndReceive(method string, body []byte, urlString string) []byte {
	req, err := http.NewRequest(method, urlString, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return []byte{}
	}

	req.Header.Set("Authorization","Bearer "+jwt)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return []byte{}
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return []byte{}
	}

	if res.StatusCode == http.StatusOK {
		return resBody
	}
	return []byte{}
}