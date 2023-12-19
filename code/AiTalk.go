package main


import (
    "fmt"
    "regexp"
    "context"
    "os"
    "os/signal"
    openai "github.com/sashabaranov/go-openai"
)


//ChatHistoryArray = append( ChatHistoryArray, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleSystem   , Content: "" } )
//ChatHistoryArray = append( ChatHistoryArray, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleUser     , Content: "" } )
//ChatHistoryArray = append( ChatHistoryArray, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleAssistant, Content: "" } )
//msgs = append( msgs, openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleAssistant, Content: resp.Choices[0].Message.Content } )

type ChatGPTConfig struct {
    API_KEY         string
    //API初期値
    ChargeAPI       int 
    //一分間あたりの使用可能数
    ChargeAPIMax    int
    //補充クールタイム。
    ChargeAPISec int
    ChatHistoryMax  int 
}
var Config ChatGPTConfig
var ThreadUse     int = 0

//使用する履歴数を確保する。
var ChatHistoryArray []openai.ChatCompletionMessage


//FunctionTalkAI
func AiTalk( Text string , ID string , FileNameAdd string ) {
    NewMessage := openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleUser, Content: Text } 

    fmt.Println( "-------------" )
    fmt.Println( "Talk AI :" + ID )
    if ThreadUse == 1 {
        fmt.Println( "Thread SKIP" )
        return
    }
    if Config.ChargeAPI <= 0 {
        fmt.Println( "API STOP" )
        return
    }

    Config.ChargeAPI = Config.ChargeAPI - 1
    ThreadUse = 1
    client := openai.NewClient( Config.API_KEY )

    msgs := []openai.ChatCompletionMessage{}
    //溢れている分を古いものから削除
    for( len( ChatHistoryArray ) > Config.ChatHistoryMax ) {
        ChatHistoryArray   = ChatHistoryArray[1:]
    }

    //会話履歴を注入する。
    n := 0
    for ( len( ChatHistoryArray ) != n ){
        //古いものから挿入する。
        msgs = append( msgs, ChatHistoryArray[ n ] )
        n++
    }


    //今回のセリフ。
    msgs = append( msgs, NewMessage )

    Sig, cancel := signal.NotifyContext( context.Background() , os.Interrupt )
    defer cancel()
    resp, err := client.CreateChatCompletion(
        Sig,
        openai.ChatCompletionRequest{ 
            Model   : openai.GPT4TurboPreview , 
            Messages: msgs  ,
            //MaxTokens        : 600  ,
            //元->設定なし。
            //まともなコードが欲しい。
            //Temperature      : 0.1,
            //まともだけど、振れ幅がない。
            Temperature      : 0.5,
            //0.6からお仕事お疲れさまでしたか?という謎文章になる。
            //Temperature      : 0.6,
        },)


    if err != nil {
        fmt.Println( "ChatGPT関係のエラー" )
        //rate limit抵触で発生した。
        //トークンの消費数超過でも発生した。4097トークン。
        fmt.Fprintln(os.Stderr, err)
        ThreadUse = 0
        return 
    }

    ChatHistoryArray = append( ChatHistoryArray , NewMessage )
    ChatHistoryArray = append( ChatHistoryArray , openai.ChatCompletionMessage{ Role: openai.ChatMessageRoleAssistant, Content: resp.Choices[0].Message.Content } )


    TextChatGPT := resp.Choices[0].Message.Content
    fmt.Println( TextChatGPT )
    //保存
    //OldSaveTextReplaceのテキストも変更すること
    TextChatGPT  = "\\0\\b[2]" + TextChatGPT + "\\n\\_a[OnSaveTalk," + ID + "_" + FileNameAdd + "]〇\\_a"

    ThreadUse = 0
    NextTalk    = TextChatGPT
    AiTalkText  = TextChatGPT
    AskText     = Text

    //検閲配列
    CheckArray := []string{ "AI" , "人工知能" }
    for _,v := range CheckArray {
        //AIという単語が入っていた場合履歴に残さない。
        var CheckAI = regexp.MustCompile( v )
        if CheckAI.MatchString( TextChatGPT )  {
            //送信と受信で二つ分。
            ChatHistoryArray   = ChatHistoryArray[:len( ChatHistoryArray ) - 2 ]
            fmt.Println( "AIワードを削除。" )
            break
        }
    }
}










