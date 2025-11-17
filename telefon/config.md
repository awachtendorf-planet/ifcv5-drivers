## config parameter

| key           | default | possible values | comment                                       |
|:--------------|:--------|:----------------|:----------------------------------------------|
| Decimals      | 2       | numeric         | 0 = amount remains the same                   |
| Protocol      | empty   |                 | must be matched with the template vendor name |
| IgnoreLRC     | false   | true/false      | overwrite template settings                   |
| IgnoreENQ     | false   | true/false      | overwrite template settings                   |
| IgnorePolling | false   | true/false      | overwrite template settings                   |


### struct definition

```golang
type Config struct {
    Template struct {
        Driver   string   `json:"driver"`            // must matched 'telefon'
        Vendor   string   `json:"vendor"`            // some useful vendor description
        Disabled bool     `json:"disabled"`          // set to true to remove the layout from the parser
        TryRun   bool     `json:"tryrun"`            // set to true to test the meta compiler
        Framing  Framing  `json:"framing,omitempty"` // optional, defined the byte stream framing
        Protocol Protocol `json:"protocol,omitempty"`// optional, defined the byte stream framing
        Layout   []Layout `json:"layout"`            //
    } `json:"template"`
}

type Layout struct {
    Name     string `json:"name,omitempty"`    // auto generated, overwrite only if absolutely necessary
    Hint     string `json:"hint,omitempty"`    // auto generated, overwrite only if absolutely necessary
    Garbage  bool   `json:"garbage,omitempty"` // auto generated, overwrite only if absolutely necessary
    Rewind   int    `json:"rewind,omitempty"`  // auto generated, overwrite only if absolutely necessary
    Field    []Field `json:"field"`            // defined the parser commands
}

type Framing struct {
    Start string `json:"start,omitempty"` // meta compiler generate garbage filter and insert 'Start' to each field
    End   string `json:"end,omitempty"`   // meta compiler generate garbage filter and append 'overread until End' to each field
}

type Protocol struct {
    Ack     string  `json:"ack,omitempty"` // defines incoming and outgoing low level packet
    Nak     string  `json:"nak,omitempty"` // defines incoming and outgoing low level packet
    Enq     string  `json:"enq,omitempty"` // defines incoming and outgoing low level packet
    LRC     LRC     `json:"lrc,omitempty"`
    Reply   Reply   `json:"reply,omitempty"`
    Polling Polling `json:"polling,omitempty"`
}

type LRC struct {
    Type   string `json:"type,omitempty"`   // defines LRC method
    Len    int    `json:"len,omitempty"`    // length of LRC
    Seed   int    `json:"seed,omitempty"`   // seed of LRC
    Inside bool   `json:"inside,omitempty"` // before or after framing
}

type Reply struct {
    Enq string `json:"enq,omitempty"` // send enq reply
}

type Polling struct {
    Char     string `json:"char,omitempty"`     // send char byte or sequence as polling packet
    Interval int    `json:"interval,omitempty"` // polling interval in seconds, 0 = disable polling (internal min 3)
}

type Field struct {
    Name         string `json:"name,omitempty"`         // auto generated, overwrite only if absolutely necessary
    Type         string `json:"type,omitempty"`         // auto generated, overwrite only if absolutely necessary (int, byte, []byte)
    Len          int    `json:"len,omitempty"`          // 0 = overread, auto generated for 'Equal' value
    Equal        string `json:"equal,omitempty"`        // match a specifically byte or sequence of bytes (hexadecimal representation)
    Endian       string `json:"endian,omitempty"`       // defined the int representation (little, big)
    Overread     bool   `json:"overread,omitempty"`     // ! describes an unneeded field, can be combined with 'Len'
    Extension    bool   `json:"extension,omitempty"`    // ! extractable value
    DialedNumber bool   `json:"dialednumber,omitempty"` // ! extractable value
    Duration     bool   `json:"duration,omitempty"`     // ! extractable value, format eg "hh:mm:ss" or "mmmss"
    Units        bool   `json:"units,omitempty"`        // ! extractable value
    Amount       bool   `json:"amount,omitempty"`       // ! extractable value, optional format eg "2" (decimals)
    CallDate     bool   `json:"calldate,omitempty"`     // ! extractable value, format eg "2006/01/02" or "2006/01/02 15:04:05"
    CallTime     bool   `json:"calltime,omitempty"`     // ! extractable value, optional if 'CallDate' only match the day, format eg "15:04:05" or "3:04pm"
    CallType     bool   `json:"calltype,omitempty"`     // ! extractable value
    User         bool   `json:"user,omitempty"`         // ! extractable value
    Outlet       bool   `json:"outlet,omitempty"`       // ! extractable value
    RoomStatus   bool   `json:"roomstatus,omitempty"`   // ! extractable value
    Article      bool   `json:"article,omitempty"`      // ! extractable value
    Quantity     bool   `json:"quantity,omitempty"`     // ! extractable value
    Format       string `json:"format,omitempty"`       // ! available for CallDate, CallTime, Duration, Amount
}
```

### date/time format

- 2006  year long
- 06    year short
- 01    month numeric
- Jan   month string
- 02    day
- 15    hour
- 3     hour
- 04    minute
- 05    second

[Golang definition](https://golang.org/pkg/time/#example_Parse)

### duration format

- h Hour
- m minute
- s second

### lrc handler

- xor


### template sample - call packet

```json
{
    "template": {

        "driver": "telefon",
        "vendor": "PBX test",
        "tryrun": true,

        "framing": {
            "start": "",
            "end": "0d0a"
        },

        "layout": [

            {
                "sample": "2017/11/09 16:28:27,00:00:38,9,12,O,99040,00540299040,,0,1009687,0,E12,Zim12,T9001,Line 1.1,0,0,,,Zim12,   0.0618,,0001.45,1,0,618,1.00,U,Zim12,,192.168.0.130,38411,192.168.0.130,38414\r\n",

                "field": [{
                    "format": "2006/01/02",
                    "calldate": true
                }, {
                    "equal": "20"
                }, {
                    "format": "15:04:05",
                    "calltime": true
                }, {
                    "equal": "2c"
                }, {
                    "format": "hh:mm:ss",
                    "duration": true
                }, {
                    "equal": "2c"
                }, {
                    "units": true
                }, {
                    "equal": "2c"
                }, {
                    "overread": true
                }, {
                    "equal": "4f2c"
                }, {
                    "overread": true
                }, {
                    "equal": "2c"
                }, {
                    "dialednumber": true
                }, {
                    "equal": "2c"
                }, {
                    "overread": true
                }, {
                    "equal": "2c45"
                }, {
                    "extension": true
                }, {
                    "equal": "2c"
                }, {
                    "overread": true
                }, {
                    "equal": "2c2c30"
                }, {
                    "amount": true
                }, {
                    "equal": "2c"
                }]
            }

        ]
    }
}
```

### template sample - protocol with LRC/Polling

```json
{
    "template": {

        "driver": "telefon",
        "vendor": "PBX test2",
        "tryrun": true,

        "framing": {
            "start": "02",
            "end": "03"
        },

        "protocol": {

            "ack": "06",
            "nak": "15",
            "enq": "05",

            "lrc": {
                "type": "xor",
                "len": 1,
                "seed": 0,
                "inside": false
            },

            "reply": {
                "enq": "06"
            },

            "polling": {
                "char": "05",
                "interval": 30
            }

        },

        "layout": [

            {
                "name": "test fias ls",
                "field": [{
                        "name": "cmd",
                        "equal": "4c53"
                    }, {
                        "equal": "7c"
                    }, {
                        "name": "data"
                    }, {
                        "equal": "7c"
                    }

                ]
            }

        ]
    }
}
```