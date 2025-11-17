## config parameter

| key                | default | possible values | comment                                                    |
|:-------------------|:--------|:----------------|:-----------------------------------------------------------|
| DriverEncoding     |         |                 | if not configured then UTF-8 is used                       |
| KeyAnswerTimeout   | 30      |                 | vendor device timeout in seconds                           |
| LicenseCode        |         |                 | TCP only                                                   |
| Radisson           | false   | true/false      | UserType/Usergroup from AccessPoints                       |
| 2800               | false   | true/false      | 2800/Vision Protocol                                       |
| AccessPoints       |         |                 | use this as access points, if key not exist use KeyOptions |
| UseRmtForKeyDelete | false   | true/false      | remote coder/ local coder                                  |

## module

### user type

key is the 1 position from KeyOptions if Radisson is true  

```json
{
    "driver": "vc3000",
    "module": "usertype",
    "station": 0,
    "template": "User Type",
    "alias": "T",
    "default": "SINGLE ROOM",
    "mapping": [{
        "key": "1",
        "value": "ROOM"
    },{
        "key": "3",
        "value": "Single Room"
    }]
}
```

### user group

key is the 2 position from KeyOptions if Radisson is true  

```json
{
    "driver": "vc3000",
    "module": "usergroup",
    "station": 0,
    "template": "User Group",
    "alias": "U",
    "default": "GUEST",
    "mapping": [{
        "key": "1",
        "value": "GUEST"
    },{
        "key": "2",
        "value": "Regular Guest"
    }]
}
```

### udf

```json
{
    "driver": "vc3000",
    "module": "udf",
    "station": 1234567890,
    "mapping": [
        {
            "key": "E?",
            "value": "P"
        }
        
    ]
}
```

