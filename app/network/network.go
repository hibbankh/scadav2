package network

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

type (
	ErrorResponse struct {
		Error   bool        `json:"error"`
		Message interface{} `json:"message"`
	}

	HttpHeader struct {
		Header string
		Value  string
	}
)

//read json response
func ReadJSONRes(r []byte, a interface{}) error {

	// log.Println("SAMPAI", r.Body)
	// err := json.NewDecoder(r.Body).Decode(&a)
	err := json.Unmarshal(r, &a)
	// fmt.Printf("%+v\n", a)

	// log.Println("SAMPAI 1", r.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	// if r.Body != nil {

	// } else {
	// 	return errors.New("http body not found")
	// }

	return nil
}

func httpClient() http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{Timeout: time.Second * 60}
	client.Transport = tr
	return client
}

func InitHttpRequest(url string, method string, params interface{}, header HttpHeader) (*http.Response, error) {

	values := params
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))

	if header.Value != "" {
		req.Header.Add(header.Header, header.Value)
	}

	if err != nil {
		log.Print(err)
		return nil, err
	}

	client := httpClient()
	resp, err := client.Do(req)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return resp, nil
}

//read incommming json data
//read incommming json data
func ReadJSONData(r *http.Request, a interface{}) error {

	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			log.Println(err)
			return err
		}
		// return true
	}

	return nil
}

//response to json data with for given
func ResponseJSON(w http.ResponseWriter, errorFlag bool, httpErrorCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErrorCode)
	msg := &ErrorResponse{
		Error:   errorFlag,
		Message: payload,
	}

	json.NewEncoder(w).Encode(msg)
}

//response to json data with for given
func GenerateJSON(payload interface{}) (string, error) {

	val, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func GetRemoteIP(req *http.Request) string {
	remoteIP := req.Header.Get("X-Real-IP")
	if remoteIP == "" {
		remoteIP = req.RemoteAddr
	}
	return strings.Split(remoteIP, ":")[0]
}

func IsSecured(req *http.Request) bool {
	if scheme := req.Header.Get("X-Forwarded-Proto"); scheme == "https" {
		return true
	}
	return false
}

func CallApi(r *http.Request) (*simplejson.Json, error) {

	httpclient := &http.Client{}

	resp, err := httpclient.Do(r)
	if err != nil {
		return nil, err

	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {

		log.Printf("got response code %d - %s", resp.StatusCode, body)
		return nil, errors.New("api request failed")
	}

	data, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	return data, nil

}
