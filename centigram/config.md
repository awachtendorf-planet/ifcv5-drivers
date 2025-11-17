## config parameter

| key               | default | possible values | comment                              |
|:------------------|:--------|:----------------|:-------------------------------------|
| DriverEncoding    |         |                 | if not configured then UTF-8 is used |
| Protocol          | 1       | 1/2             | 1 = HIS, 2 = Encore                  |
| SwapRequest       | true    | true/false      |                                      |
| SendGuestName     | false   | true/false      | Record Type D                        |
| SendGuestLanguage | false   | true/false      | Record Type 8                        |


## module

### language code

```json
{
    "driver": "centigram",
    "module": "languagecode",
    "station": 0,
    "default": "0",
    "mapping": [ {
        "key": "US",
        "value": "1"
    }, {
        "key": "EN",
        "value": "1"
    }, {
        "key": "ES",
        "value": "2"
    }, {
        "key": "DE",
        "value": "3"
    }]
}
```
