## config parameter Inhova

| key                | default | possible values | comment                                            |
|:-------------------|:--------|:----------------|:---------------------------------------------------|
| DriverEncoding     |         | Windows-1252    | if not configured then UTF-8 is used               |
| KeyAnswerTimeout   | 30      |                 | vendor device timeout in seconds per card          |
| SendTrack2         | true    | true/false      | send track2 data                                   |
| SendActivationTime | false   | true/false      | send arrival date and time                         |
| AccessPoints       |         | string          | use this as access points, if empty use KeyOptions |


## module

### accesspoints

```json
{
    "driver": "inhova",
    "module": "accesspoint",
    "station": 1234567890,
    "mapping": [{
        "key": "1",
        "value": "AP1"
    },{
        "key": "2",
        "value": "AP2"
    },{
        "key": "3",
        "value": "AP3"
    },{
        "key": "4",
        "value": "AP4"
    },{
        "key": "5",
        "value": "AP5"
    },{
        "key": "6",
        "value": "AP6"
    },{
        "key": "7",
        "value": "AP7"
    },{
        "key": "8",
        "value": "AP8"
    },{
        "key": "9",
        "value": "AP9"
    }]
}
```

`KeyOptions="10100111000"` will set AP4,AP5,AP6,AP9. (Counting starts with 1 from the right)  