package main

import (
    "fmt"
    "log"
	"net/http"
	"encoding/json"
    "os"
    "net/url"
    "io/ioutil"
)

var welcomeMessage = `{
    "name" : "App Search Converter",
    "cluster_name" : "elasticsearch",
    "cluster_uuid" : "saO0mdn8SNSSkVEchTSm3A",
    "version" : {
      "number" : "7.00",
      "build_flavor" : "default",
      "build_type" : "tar",
      "build_hash" : "ef48eb35cf30adf4db14086e8aabd07ef6fb113f",
      "build_date" : "2020-03-26T06:34:37.794943Z",
      "build_snapshot" : false,
      "lucene_version" : "8.4.0",
      "minimum_wire_compatibility_version" : "6.8.0",
      "minimum_index_compatibility_version" : "6.0.0-beta1"
    },
    "tagline" : "You Know, for converting to App Search"
  }`

//Configuration is a structure to hold the global configuration data read from a file
type Configuration struct {
    Port int `json:"port"`
    KEY string `json:"API_KEY"`
    ASEndpoint  string `json:"App_Search_endpoint"`
}

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, welcomeMessage)
    fmt.Println("Endpoint Hit: homePage")
}

func buildAppSearchAddress(endpoint string) string{
    u, err := url.Parse(endpoint)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("AS Endpoint: ", configuration.ASEndpoint)
    base, err := url.Parse(configuration.ASEndpoint)
    if err != nil {
        log.Fatal(err)
    }
    fullAddress := base.ResolveReference(u).String()
    return fullAddress
}

func getAppSearch(endpoint string)string{
    address := buildAppSearchAddress(endpoint)
    fmt.Println("Fetching: ", address)
    // Create a Bearer string by appending string access token
    var bearer = "Bearer " + configuration.KEY

    // Create a new request using http
    req, err := http.NewRequest("GET", address, nil)

    // add authorization header to the req
    req.Header.Add("Authorization", bearer)

    // Send req using http Client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error on response.\n[ERRO] -", err)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    strBody := string(body)
    return strBody
}

func getEngines(w http.ResponseWriter, r *http.Request){

    enginesEndpoint := "/api/as/v1/engines"
    
    body := getAppSearch(enginesEndpoint)

    fmt.Fprintf(w, body)
    fmt.Println("Endpoint Hit: Engines")
}

func handleRequests(port int) {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/engines", getEngines)
    address := fmt.Sprintf(":%d", port)
    fmt.Println("Starting server with address= ",address)
    log.Fatal(http.ListenAndServe(address, nil))
}

var configuration = Configuration{}
func main() {
    file, _ := os.Open("config.json")
    defer file.Close()
    decoder := json.NewDecoder(file)
    err := decoder.Decode(&configuration)
    if err != nil {
        fmt.Println("error:", err)
    }
    handleRequests(configuration.Port)
}