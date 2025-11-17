## config parameter Ilco

| key                | default | possible values   | comment                                                 |
|:-------------------|:--------|:------------------|:--------------------------------------------------------|
| DriverEncoding     |         | Windows-1252      | if not configured then UTF-8 is used                    |
| KeyAnswerTimeout   | 30      |                   | vendor device timeout in seconds                        |
| SendTrack2         | true    | true/false        | send track2 data                                        |
| Protocol           | 0       | 0/1               | 0 = direct connection, 1 = gateway                      |
| AuthNumber         |         | string            | '' empty = don't send or send any string as auth number |

