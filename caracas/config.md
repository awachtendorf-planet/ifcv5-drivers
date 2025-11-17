## config parameter

| key            | default    | possible values | comment                                 |
|:---------------|:-----------|:----------------|:----------------------------------------|
| DriverEncoding | ISO 8859-1 |                 | ASCII character set                     |
| BindMode       | true       | bool            | send bind on start up                   |
| NewVersionMode | true       | bool            | Caracas V4                              |
| SendTime       | 0          | uint64          | Milliseconds                            |
| TimeZone       |            | IANA standard   | IANA timezone e.g Europe/Berlin         |


### language code

```json
{
    "driver": "caracas",
    "module": "languagecode",
    "station": 0,
    "default": "01",
    "mapping": [
    {
        "key": "US",
        "value": "01"
    }, {
        "key": "EN",
        "value": "01"
    }, {
        "key": "DE",
        "value": "02"
    }, {
        "key": "FR",
        "value": "03"
    }, {
        "key": "RU",
        "value": "04"
    }, {
        "key": "JA",
        "value": "05"
    }, {
        "key": "ES",
        "value": "06"
    }, {
        "key": "EL",
        "value": "07"
    }, {
        "key": "CMN",
        "value": "08"
    }, {
        "key": "PT",
        "value": "09"
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