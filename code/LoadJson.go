package main

import (
    "fmt"
    "os"
    "encoding/json"
)

func LoadJson(){
	JsonChatGPTConfig, err := os.Open( Directory + "/Config.json")
	if err != nil {
        fmt.Println( err )
	}
	defer JsonChatGPTConfig.Close()
    decoder := json.NewDecoder( JsonChatGPTConfig )
    err     = decoder.Decode( &Config )
	if err != nil {
        fmt.Println(  err  )
    }
}






