## config parameter

| key            | default | possible values | comment                              |
|:---------------|:--------|:----------------|:-------------------------------------|
| DriverEncoding |         | ISO 8859-1      | if not configured then UTF-8 is used |
| Decimals       | 2       |                 | 0 = amount remains the same          |
| Protocol       | 1       | 1-8             |                                      |

### protocol types

* 1 = Guestlink
* 2 = Eclipse
* 3 = Sonifi
* 4 = Movielink
* 5 = OnCommand
* 6 = Quadriga
* 7 = TripleGuest
* 8 = MagiNet
* 9 = MagiNetEnhanced (curently not supported)

Eclipse, Sonifi, OnCommand, TripleGuest, MagiNet  
- Reservierungsnummer 6 stellig
- Reservierungsnummer links ausgerichtet  

Movelink, Quadriga  
- Reservierungsnummer 6 stellig
- Reservierungsnummer rechts ausgerichtet  

MagiNetEnhanced  
- Reservierungsnummer 8 stellig

Eclipse, TripleGuest  
- mit "Helo" Paket
  

## module

### language code

```json
{
    "driver": "guestlink",
    "module": "languagecode",
    "station": 0,
    "default": " ",
    "mapping": [{
        "key": "en_US",
        "value": "E"
    }, {
        "key": "en_EN",
        "value": "EA"
    }, {
        "key": "de_DE",
        "value": "G"
    }, {
        "key": "fr_FR",
        "value": "F"
    }, {
        "key": "it_IT",
        "value": "I"
    }, {
        "key": "us",
        "value": "E"
    }, {
        "key": "en",
        "value": "E"
    }, {
        "key": "de",
        "value": "G"
    }, {
        "key": "fr",
        "value": "F"
    }, {
        "key": "it",
        "value": "I"
    }]
}
```

### pay tv right

```json
{
    "driver": "guestlink",
    "module": "paytvright",
    "alias": "TV",
    "station": 0,
    "default": "STD",
    "mapping": [{
        "key": "0",
        "value": "CLO"
    }, {
        "key": "1",
        "value": "KID"
    }, {
        "key": "2",
        "value": "STD"
    }, {
        "key": "3",
        "value": "VIP"
    }]
}

```