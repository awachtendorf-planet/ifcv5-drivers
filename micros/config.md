## config parameter

| key                         | default | possible values | comment                              |
|:----------------------------|:--------|:----------------|:-------------------------------------|
| Decimals                    | 2       |                 | 0 = amount remains the same          |

### paymentmethod

```json
{
    "driver": "micros",
    "module": "paymentmethod",
    "station": 123456789,
    "mapping": [{
        "key": "1",
        "value": "10"
    }]
}
```