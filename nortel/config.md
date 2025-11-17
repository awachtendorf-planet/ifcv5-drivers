## config parameter

| key                 | default | possible values | comment                                                                                                |
|:--------------------|:--------|:----------------|:-------------------------------------------------------------------------------------------------------|
| DriverEncoding      |         | ISO 8859-1      | if not configured then UTF-8 is used                                                                   |
| Decimals            | 2       |                 | 0 = amount remains the same                                                                            |
| ManagmentDevice     | 0       |                 | device number of the configured managment device                                                       |
| ExtensionWidth      | 7       | > 2             |                                                                                                        |
| GuestNameLength     | 23      |                 | 23 recommended, 27 maximum                                                                             |
| SendGuestName       | true    | true/false      | send guest name at check-in and data change                                                            |
| SendGuestLanguage   | true    | true/false      | send guest language at check-in and data change                                                        |
| SendGuestVIPState   | true    | true/false      | send guest VIP state at check-in and data change                                                       |
| SendGuestNameLength | true    | true/false      | send name format "Guest Name" 12                                                                       |
| SendRoomVacant      | false   | true/false      | send "Room Vacant" at checkout                                                                         |
| SendClassOfService  | true    | true/false      | send CCRS/ECC1/ECC2 packets                                                                            |
| HandleWakeup        | true    | true/false      | disable wakeup  from pms                                                                               |
| PostAllRecords      | true    | true/false      | S/X/E create each record for pms (if true MergeRecords will be ignored)                                |
| PostLastRecordOnly  | false   | true/false      | only create a pms record from E record                                                                 |
| MergeRecords        | false   | true/false      | create a single pms record by merge S/X/E records (only if PostAllRecords/PostLastRecordOnly is false) |
| PostInternalCall    | false   | true/false      |                                                                                                        |
| DialledNumberRegexp |         | regexp          | define a regexp to filter the dialled number as the first group match , eg "A\\d{3}(.*)"               |


## module

### room status

```json
{
  "driver": "nortel",
  "module": "roomstatus",
  "template": "Room Status",
  "default": "RE",
  "station": 0,
  "mapping": [
    {
      "key": "dirty",
      "value": "RE"
    },
    {
      "key": "progress",
      "value": "PR"
    },
    {
      "key": "clean",
      "value": "CL"
    },
    {
      "key": "passed",
      "value": "PA"
    },
    {
      "key": "failed",
      "value": "FA"
    },
    {
      "key": "skipped",
      "value": "SK"
    },
    {
      "key": "out of order",
      "value": "NS"
    }
  ]
}
```

### language code

```json
{
  "driver": "nortel",
  "module": "languagecode",
  "station": 0,
  "mapping": [
    {
      "key": "FR",
      "value": "3"
    },
    {
      "key": "DE",
      "value": "1"
    }
  ]
}
```

