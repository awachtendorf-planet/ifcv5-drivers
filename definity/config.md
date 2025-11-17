## config parameter

| key            | default    | possible values | comment                                 |
|:---------------|:-----------|:----------------|:----------------------------------------|
| DriverEncoding | ISO 8859-1 |                 | ASCII character set                     |
| Protocol       | 1          | 1               | currently only ascii mode are supported |
| RecordFormat   | 1          | 1-2             | standard/extended                       |
| CoveragePath   | 0000       | 4 digit         |                                         |

### protocol types

* 1 = ascii mode
* 2 = transparent mode (not supported yet)


### record format types

* 1 = standard format, extension 5 digit, name 15 digit
* 2 = extended format, extension 7 digit, name 30 digit

## module

### language code

```json
{
    "driver": "definity",
    "module": "languagecode",
    "station": 0,
    "default": "20",
    "mapping": [
    {
        "key": "US",
        "value": "20"
    }, {
        "key": "EN",
        "value": "26"
    }, {
        "key": "DE",
        "value": "2a"
    }, {
        "key": "FR",
        "value": "2b"
    }, {
        "key": "RU",
        "value": "34"
    }, {
        "key": "JA",
        "value": "21"
    }, {
        "key": "ES",
        "value": "22"
    }, {
        "key": "EL",
        "value": "23"
    }, {
        "key": "CMN",
        "value": "24"
    }, {
        "key": "PT",
        "value": "28"
    }, {
        "key": "AR",
        "value": "2d"
    }
]
}
```

### room status

```json
{
    "driver": "definity",
    "module": "roomstatus",
    "template": "Room Status",
    "alias": "RS",
    "station": 0,
    "mapping": [{
        "key": "housekeeper in room",
        "value": "1"
    },{
        "key": "room clean - vacant",
        "value": "2"
    },{
        "key": "room clean - occupied",
        "value": "3"
    },{
        "key": "room dirty - vacant",
        "value": "4"
    },{
        "key": "room dirty - vacant",
        "value": "5"
    },{
        "key": "room clean - needs inspection",
        "value": "6"
    }]
}

```