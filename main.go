package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    response, err := http.Get("https://httpbin.org/ip")
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
        return
    }

    data, _ := ioutil.ReadAll(response.Body)
    var objmap map[string]*json.RawMessage
    json.Unmarshal(data, &objmap)
    var origin string
    json.Unmarshal(*objmap["origin"], &origin)

    fmt.Println("Your public IP address is:", origin)
}

