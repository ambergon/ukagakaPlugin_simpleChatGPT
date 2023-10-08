package main

/*
   #include <windows.h>
   #include <stdlib.h>
   #include <string.h>
*/
import "C"

import (
    "fmt"
    "unsafe"
    "strings"
    "regexp"
)

func main() {
    fmt.Println( "test" )
}


var Directory string
var References []string 
var CheckID         = regexp.MustCompile("^ID: ")
var CheckReference  = regexp.MustCompile("^Reference.+?: ")

//反応確認用テキスト
var DefaultSurfaceText string = "\\0\\s[6]"

type ResponseStruct struct {
    Shiori  string
    Sender  string
    Charset string
    Marker  string
    Value   string
}
func GetResponse( r *ResponseStruct ) string {
    V := ""
    if r.Value  != "" { V = "Value: "  + r.Value     + "\r\n" }
    res :=  r.Shiori    + "\r\n" + 
            r.Sender    + "\r\n" + 
            r.Charset   + "\r\n" + 
            V + "\r\n\r\n"
    return res
}

var NextTalk        string  = ""
var ChargeAPISec    int     = 0


//export load
func load(h C.HGLOBAL, length C.long ) C.BOOL {
    fmt.Println( "load simpleChatGPT" )
    Directory = C.GoStringN(( *C.char )( unsafe.Pointer( h )), ( C.int )( length ))
    fmt.Println( Directory  )

    //設定読み込み。
    LoadJson()

	C.GlobalFree( h )
	return C.TRUE
}


//export unload
func unload() bool {
    fmt.Println( "unload simpleChatGPT" )
	return true
}


//export request
func request( h C.HGLOBAL, length *C.long ) C.HGLOBAL {
	RequestText := C.GoStringN(( *C.char )( unsafe.Pointer( h )), ( C.int )( *length ))
	C.GlobalFree( h )


    Value           := ""
    Marker          := ""
    ID              := ""
    References      = []string{}
    var NOTIFY bool = false

    Response := new( ResponseStruct )
    Response.Sender  = "Sender: GolangAI"
    Response.Charset = "Charset: UTF-8"

    //IDとReference
    //必要な情報を分解する。
    RequestLines := strings.Split( RequestText , "\r\n" )
    for _ , line := range RequestLines {
        if( line == "NOTIFY PLUGIN/2.0" ){
            //"GET PLUGIN/2.0";
            NOTIFY = true

        } else if CheckID.MatchString( line )  {
            //fmt.Println( line )
            ID = CheckID.ReplaceAllString( line , "" )

        } else if CheckReference.MatchString( line )  {
            //fmt.Println( line )
            ref := CheckReference.ReplaceAllString( line , "" )
            References = append( References , ref )

        } else {
            //fmt.Println( line )
        }
    }

    //実行関数
    if ID == "OnOtherGhostTalk" {
        //AiTalkModeとして、OnCommnicateが処理されるが何も発言されない\\![]を使用する。
        if References[3] == "OnCommunicate" && References[4] == "\\![]"{
            IdReferences  := strings.Split( References[5] , "\u0001" )
            if IdReferences[0] == "user" {
                fmt.Print( "\n>> " )
                fmt.Println( IdReferences[1] )
                go AiTalk( IdReferences[1] , ID , IdReferences[1] )
            }
        }

    } else if ID == "OnSaveTalk"  {
        Value = OnSaveTalk()

    } else if ID == "OnSecondChange"  {
        //ChargeAPISecごとに補充する。
        if Config.ChargeAPIMax > Config.ChargeAPI {
            ChargeAPISec++
            if ChargeAPISec >= Config.ChargeAPISec {
                Config.ChargeAPI++
                ChargeAPISec = 0
            }
        }
        if NextTalk != "" && NOTIFY == false {
            Value       = NextTalk
            NextTalk    = ""
        }

    } else {
        //fmt.Println( "no touch :" + ID )
        //fmt.Print( "NOTIFY : " )
        //fmt.Println( NOTIFY )
        //fmt.Print( "References : " )
        //fmt.Println( References )
        //fmt.Println( "" )
    }


    if Value == "" {
        Response.Shiori  = "PLUGIN/2.0 204 No Content"
    } else {
        Response.Shiori = "PLUGIN/2.0 200 OK"
        Response.Value  = Value
    }
    if Marker != "" {
        Response.Marker  = Marker
    }

    res_buf := C.CString( GetResponse( Response ))
    defer C.free( unsafe.Pointer( res_buf ))

	res_size := C.strlen( res_buf )
	ret      := C.GlobalAlloc( C.GPTR , ( C.SIZE_T )( res_size ))
	C.memcpy(( unsafe.Pointer )( ret ) , ( unsafe.Pointer )( res_buf ) , res_size )
	*length = ( C.long )( res_size )
	return ret
}






















