## config parameter

| key                | default | possible values   | comment                                                 |
|:-------------------|:--------|:------------------|:--------------------------------------------------------|
| DriverEncoding     |         | Windows-1252      | if not configured then UTF-8 is used                    |
| AdjustnameLength   | false   | true/false        | if yes use 25 instead of 23 characters for displayname  |
| UseNameDisplay     | true    | true/false        | send guestname                                          |
| MessageOffCharacter| 0       | 0 / 2             | any other value than 0 or 2 will be interpreted as 0    |
| ShortRoomname      | true    | true/false        | true = 5 digits roomname length, false = 6 digits       |
| SupportDND         | true    | true/false        |                                                         |
| Protocol           | 3       | 0/1/2/3           | 0,1,2 are equal to TIG_Protocol, 3 is the Tiger spec    |

