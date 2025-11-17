## config parameter

| key               | default | possible values | comment                                          |
|:------------------|:--------|:----------------|:-------------------------------------------------|
| DriverEncoding    |         |                 | if not configured then UTF-8 is used             |
| ExtensionWidth    | 5       | 5 or 7          |                                                  |
| SwapRequest       | true    | true/false      |                                                  |
| AliveRecord       | true    | true/false      | send Alive Record every 10 seconds on idle state |
| RestrictionRecord | false   | true/false      | SX200 false, SX2000 true                         |

## module

### class of service

```json
{
    "driver": "mitel",
    "module": "classofservice",
    "station": 0,
    "default": "0",
    "mapping": [{
        "key": "4",
        "value": "1"
    }]
}
```

### room status

```json
{
    "driver": "mitel",
    "module": "roomstatus",
    "template": "Room Status",
    "alias": "RS",
    "station": 0,
    "mapping": [{
        "key": "clean",
        "value": "1"
    }]
}
```