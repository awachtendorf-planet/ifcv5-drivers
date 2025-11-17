## config parameter

| key              | default | possible values | comment                                            |
|:-----------------|:--------|:----------------|:---------------------------------------------------|
| DriverEncoding   |         |                 | if not configured then UTF-8 is used               |
| Decimals         | 2       |                 | 0 = amount remains the same                        |
| ExtensionWidth   | 6       | 4 or 6          |                                                    |
| SimpleCheckin    | false   | true/false      |                                                    |
| SendCheckoutDate | false   | true/false      | filled with blanks                                 |
| SendGuestName    | false   | true/false      | filled with blanks                                 |
| SendHappyHour    | false   | true/false      | removed from packet, extend guestname length to 32 |


## module

### room status

```json
{
    "driver": "bartech",
    "module": "roomstatus",
    "template": "Room Status",
    "alias": "RS",
    "station": 0,
    "mapping": [{
        "key": "vacant, cleaned",
        "value": "01"
    },{
        "key": "occupied, clean",
        "value": "02"
    },{
        "key": "out of order",
        "value": "03"
    },{
        "key": "bulb to be changed",
        "value": "04"
    },{
        "key": "tap to be repaired",
        "value": "05"
    }]
}

```

