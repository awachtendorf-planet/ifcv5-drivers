## config parameter

| key                | default     | possible values | comment                                                          |
|:-------------------|:------------|:----------------|:-----------------------------------------------------------------|
| DriverEncoding     |             | ISO 8859-1      | if not configured then UTF-8 is used                             |
| Decimals           | 2           |                 | 0 = amount remains the same, must be configured with 0           |
| Protocol           | 2           | 1-2             | extension width 5 or 8                                           |
| DataTransfer       | 0           | 0-1             | data as article or total amount                                  |
| WaitForReplyPacket | true        | true/false      | only on serial layer, tcp always expected a reply packet         |
| OutletRegexp       | ".{4}(\\d)" | regexp          | outlet if DataTransfer=1, default 5th position from Dialled Code |


### protocol types

* 1 = AHL4400_5 (extension width 5)
* 2 = AHL4400_8 (extension width 8)

### datatransfer types

* 0 = article(\*amount)
* 1 = total amount

## module

### language code

```json
{
    "driver": "ahl",
    "module": "languagecode",
    "station": 0,
    "default": " ",
    "mapping": [{
        "key": "en_US",
        "value": "1"
    }, {
        "key": "en_EN",
        "value": "1"
    }, {
        "key": "de_DE",
        "value": "2"
    }, {
        "key": "fr_FR",
        "value": "3"
    }, {
        "key": "it_IT",
        "value": "4"
    }, {
        "key": "ja_JA",
        "value": "5"
    }, {
        "key": "es_ES",
        "value": "6"
    }, {
        "key": "US",
        "value": "1"
    }, {
        "key": "EN",
        "value": "1"
    }, {
        "key": "DE",
        "value": "2"
    }, {
        "key": "FR",
        "value": "3"
    }, {
        "key": "IT",
        "value": "4"
    }, {
        "key": "JA",
        "value": "5"
    }, {
        "key": "ES",
        "value": "6"
    }]
}
```

### class of service

```json
{
    "driver": "ahl",
    "module": "classofservice",
    "station": 0,
    "default": "  ",
    "mapping": [{
        "key": "1",
        "value": "01"
    }, {
        "key": "2",
        "value": "02"
    }, {
        "key": "3",
        "value": "03"
    }]
}

```

### room status

```json
{
    "driver": "ahl",
    "module": "roomstatus",
    "template": "Room Status",
    "alias": "RS",
    "station": 0,
    "mapping": [{
        "key": "clean",
        "value": "0"
    },{
        "key": "dirty",
        "value": "1"
    }]
}

```