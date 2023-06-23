package main

import (
    "fmt"
    "os"
    "regexp"
)

//AIが以前話したScriptを保存しておく。
var AiTalkText string  = ""
var AskText    string  = ""

//末尾のテキストは、 AiTalk.go で使用されている。
var OldSaveTextReplace      = regexp.MustCompile(`\\n\\_a\[OnSaveTalk,.*$`)

func OnSaveTalk() string {
    var Value string = ""
    if AiTalkText != "" {
        Text := ( ">>" + AskText + "\n" + AiTalkText + "\n" )
        SaveTalk( Text )
        Value  = DefaultSurfaceText + "保存しました。"
        AiTalkText = ""
        AskText    = ""
    }
    return Value
}


//分けておけば、OnCloseから使用できる。
func SaveTalk( Text string ){
    //fmt.Println( Directory + "log/AskChatGPT.uka" )
    fp,err := os.OpenFile( Directory + "log/AskChatGPT.uka" , os.O_WRONLY|os.O_APPEND|os.O_CREATE , 0665 )
    if err != nil {
        fmt.Println( "保存error" )
        return 
    }
    defer fp.Close()

    Text = OldSaveTextReplace.ReplaceAllString( Text , "" )
    fmt.Fprintln( fp , Text )
}



