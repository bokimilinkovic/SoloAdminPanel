package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	cl := &http.Client{}
	cookie := obtainCookie(cl)
	if cookie != nil {
		infos, err := sendSettings(cl, cookie)
		if err != nil {
			panic(err)
		}

		fmt.Println("Response : ", infos)
	}
}

func obtainCookie(cl *http.Client) *http.Cookie {
	req, err := http.NewRequest("POST", "http://localhost:8080/v1/authenticate/google", strings.NewReader(`{"auth_code":"4/0AY0e-g4qvY4DfX0Ft-FmI7lQ0zTR0gqE8BCIdUBtq8LbETsUYZa-9Zwg9Zeyi7W1zyzJtQ"}`))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-type", "application/json")
	res, err := cl.Do(req)
	if err != nil {
		return nil
	}
	if res.StatusCode == 200 {
		return res.Cookies()[0]
	}

	return nil
}

func sendSettings(cl *http.Client, cookie *http.Cookie) (string, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/v1/settings", nil)
	if err != nil {
		panic(err)
	}
	req.AddCookie(cookie)
	res, err := cl.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}

	return "", nil
}
