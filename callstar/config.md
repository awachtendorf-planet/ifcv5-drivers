## config parameter

| key            | default | possible values | comment                              |
|:---------------|:--------|:----------------|:-------------------------------------|
| DriverEncoding |         |                 | if not configured then UTF-8 is used |
| Decimals       | 2       |                 | 0 = amount remains the same          |
| SwapRequest    | true    | true/false      |                                      |
| SwapLabel      | true    | true/false      | send Z+/Z-                           |
| MinibarMode    | 0       | 0/1             | 1 = Minibar Ticket as Sales Outlet   |



## module

### room status

```json
{
    "driver": "callstar",
    "module": "roomstatus",
    "template": "Room Status",
    "alias": "RS",
    "station": 0,
    "mapping": [{
        "key": "clean",
        "value": "1"
    },{
        "key": "dirty",
        "value": "2"
    }]
}

```

### sales outlet

```json
{
  "driver": "callstar",
  "module": "outlet",
  "alias": "SO",
  "template": "Generic Outlet",
  "station": 0,
  "mapping": [
    {
      "key": "101",
      "value": "0001"
    },
    {
      "key": "202",
      "value": "0002"
    }
  ]
}

```

### vip status

```json
{
    "driver": "callstar",
    "module": "vipstatus",
    "station": 0,
    "default": "VIP",
    "mapping": [{
        "key": "1",
        "value": "VIP1"
    }, {
        "key": "2",
        "value": "VIP2"
    }]
}

```

### language code

```json
{
    "driver": "callstar",
    "module": "languagecode",
    "station": 0,
    "default": "EA",
    "mapping": [{
        "key": "en_US",
        "value": "US"
    }, {
        "key": "en_EN",
        "value": "UK"
    }, {
        "key": "de_DE",
        "value": "GE"
    }, {
        "key": "fr_FR",
        "value": "FR"
    }, {
        "key": "US",
        "value": "US"
    }, {
        "key": "EN",
        "value": "UK"
    }, {
        "key": "DE",
        "value": "GE"
    }, {
        "key": "FR",
        "value": "FR"
    }]
}

```

