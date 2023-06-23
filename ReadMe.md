# ukagakaPlugin_SimpleChatGPT
OnOtherGhostTalkを利用したChatGPTとの対話用プラグイン。<br >
感情表現等は取り払ったシンプルな対話用。<br >


## Usage
Communicateボックスからの送信に対して、ゴースト側で`\![]`のみを返した場合のみ起爆するようにしている。<br >
これだとOnOtherGhostTalkは起爆するが、余計なバルーンが発生しないため。<br >
この機能を通して発言したモノには最後に〇がつく。<br >
これをクリックするとプラグインディレクトリに対話内容が保存される。<br >

```.yaya
OnCommunicate {
    "\![]"
}
```



## 動作環境
Windows 10<br >
SSP 2.6.48<br >


## 必要なモノ。
ChatGPTのAPI.<br >
Config.jsonの中のAPI_KEYの項目に挿入すること。<br >


## 他
Golangで書いたDLLの為、プラグインの再読み込み等(FreeLibrary)をするとフリーズする。<br >


## Author
ambergon
