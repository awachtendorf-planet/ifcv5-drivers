## config parameter

| key                         | default | possible values | comment                              |
|:----------------------------|:--------|:----------------|:-------------------------------------|
| DriverEncoding              |         |                 | if not configured then UTF-8 is used |
| DBSwapTimeLimit             | 60      | seconds         | overwrite db swap limit              |
| Decimals                    | 2       |                 | 0 = amount remains the same          |
| KeyAnswerTimeout            | 30      |                 | vendor device timeout in seconds     |
| IndividualRoomStatusPackets | false   | true/false      |                                      |
| GuestMessageHandling        | 0       | 0/1/2           | IFCDEV-19                            |
|                             |         |                 | 0 use RE, discard XL                 |
|                             |         |                 | 1 discard ML from RE, use XL         |
|                             |         |                 | 2 use both RE and XL                 |

## module

### language code

```json
{
    "driver": "fidserv",
    "module": "languagecode",
    "station": 0,
    "alias": "GL",
    "default": "EA",
    "mapping": [{
        "key": "en_US",
        "value": "EA"
    }, {
        "key": "en_EN",
        "value": "EA"
    }, {
        "key": "de_DE",
        "value": "GE"
    }, {
        "key": "fr_FR",
        "value": "FR"
    }, {
        "key": "it_IT",
        "value": "IT"
    }, {
        "key": "ja_JA",
        "value": "JA"
    }, {
        "key": "es_ES",
        "value": "SP"
    }, {
        "key": "US",
        "value": "EA"
    }, {
        "key": "EN",
        "value": "EA"
    }, {
        "key": "DE",
        "value": "GE"
    }, {
        "key": "FR",
        "value": "FR"
    }, {
        "key": "IT",
        "value": "IT"
    }, {
        "key": "JA",
        "value": "JA"
    }, {
        "key": "ES",
        "value": "SP"
    }]
}
```

### sales outlet

```json
{
  "driver": "fidserv",
  "module": "outlet",
  "alias": "SO",
  "template": "Generic Outlet",
  "station": 1234567890,
  "mapping": [
    {
      "key": "666",
      "value": "123"
    },
    {
      "key": "789",
      "value": "ABC"
    }
  ]
}
```

### room status

```json
{
    "driver": "fidserv",
    "module": "roomstatus",
    "template": "Room Equipment Status",
    "alias": "RS",
    "station": 0,
    "mapping": [{
        "key": "clean",
        "value": "1"
    }]
}
```

### class of service

```json
{
    "driver": "fidserv",
    "module": "classofservice",
    "station": 0,
    "alias": "CS",
    "default": "0",
    "mapping": [{
        "key": "4",
        "value": "1"
    }]
}
```

### map inbound fields in scope of record id

eg. map CT to X1 for record PS/PR  

```json
{
    "driver": "fidserv",
    "module": "map.inbound",
    "station": 1234567890,
    "mapping": [
        {
            "key": "PSX1",
            "value": "CT"
        },
        {
            "key": "PRX1",
            "value": "CT"
        }
    ]
}
```


### udf per station

By default A0 is mapped to UDF1, A1 to UDF2 etc. Use the per station mapping to overwrite the default settings.


```json
{
    "driver": "fidserv",
    "module": "udf",
    "station": 1234567890,
    "mapping": [
        {
            "key": "GIA0",
            "value": "Email"
        },
        {
            "key": "GIA1",
            "value": "UDF1"
        },
        {
            "key": "GCA8",
            "value": "GN"
        },
        {
            "key": "GCA9",
            "value": "GN"
        },{
            "key": "GIGX",
            "value": "DisplayName"
        },{
            "key": "KR$1",
            "value": "$4"
        }
    ]
}
```

### minibar right (auto generated)

```json
{
    "driver": "fidserv",
    "module": "minibarright",
    "alias": "MR",
    "station": 0,
    "default": "MN",
    "mapping": [{
        "key": "0",
        "value": "ML"
    }, {
        "key": "2",
        "value": "MU"
    }]
}
```

### pay tv right (auto generated)

```json
{
    "driver": "fidserv",
    "module": "paytvright",
    "alias": "TV",
    "station": 0,
    "default": "TU",
    "mapping": [{
        "key": "0",
        "value": "TN"
    }, {
        "key": "1",
        "value": "TX"
    }, {
        "key": "2",
        "value": "TM"
    }, {
        "key": "3",
        "value": "TU"
    }]
}
```

### video right (auto generated)

```json
{
    "driver": "fidserv",
    "module": "videoright",
    "alias": "VR",
    "station": 0,
    "default": "VN",
    "mapping": [{
        "key": "0",
        "value": "VN"
    }, {
        "key": "1",
        "value": "VB"
    }, {
        "key": "2",
        "value": "VA"
    }]
}
```