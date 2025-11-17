## config parameter

| key              | default | possible values | comment                                  |
|:-----------------|:--------|:----------------|:-----------------------------------------|
| DriverEncoding   |         |                 | if not configured then UTF-8 is used     |
| KeyAnswerTimeout | 30      |                 | vendor device timeout in seconds         |
| Track4PP         | false   | true/false      | fill PP with Track4 data                 |
| Track4PF         | false   | true/false      | fill PF with Track4 data                 |
| Track2Origin     | 0       | 0/1/2           | 1 = PMS, 2 = Vendor provides Track2 data |
| Track2Trim       | false   | true/false      | remove trailing ff                       |


## module

### udf

```json
{
    "driver": "visionline",
    "module": "udf",
    "station": 0,
    "mapping": [{
        "key": "DNF",
        "value": "Email"
    }, {
        "key": "DNG",
        "value": "Phone"
    }, {
        "key": "FNF",
        "value": "Email"
    }, {
        "key": "FNG",
        "value": "Phone"
    }]
}
```

### accesspoints

```json
{
    "driver": "visionline",
    "module": "accesspoint",
    "template": "Code Card", 
    "alias": "CR",
    "station": 1234567890,
    "mapping": [{
        "key": "1",
        "value": "AP1"
    },{
        "key": "2",
        "value": "AP2"
    },{
        "key": "3",
        "value": "AP3a,AP3b"
    }]
}
```

