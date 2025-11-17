## config parameter

| key                | default | possible values   | comment                                                 |
|:-------------------|:--------|:------------------|:--------------------------------------------------------|
| DriverEncoding     |         | Windows-1252      | if not configured then UTF-8 is used                    |
| KeyAnswerTimeout   | 30      |                   | vendor device timeout in seconds                        |
| Protocol           | 1       | 1/2/3 or 13/14/15 | 13,14,15 parameters                                     |
| SendKeyDelete      | true    | true/false        | send key delete packet                                  |
| SendTrack2         | true    | true/false        | send track2 data                                        |
| AccessPointOverlay |         | string (10x)      | 0 = always denied, 1=always allowed, x = keep pms value |
| KeyID              |         | UDFx              | use UDFx from pms as keyID value (only 14/15)           |

