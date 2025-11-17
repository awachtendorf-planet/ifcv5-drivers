## config parameter

| key                | default | possible values   | comment                                                 |
|:-------------------|:--------|:------------------|:--------------------------------------------------------|
| DriverEncoding     |         | Windows-1252      | if not configured then UTF-8 is used                    |
| GDXMode            | true    | true/false        | Send name records                                       |
| SendTG67           | false   | true/false        | send telegram 67                                        |
| SendTG80           | false   | true/false        | send telegram 80                                        |
| DisplayType        | 0       | 0/1               | DisplayType in Telegram 41, these are the actual values |
| DirectWakeupMode   | false   | true/false        | false: PBX controls the wakeup, true: PMS got control   |

