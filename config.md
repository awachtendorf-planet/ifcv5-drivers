## global config parameter

| key                 | default | possible values | comment             |
|:--------------------|:--------|:----------------|:--------------------|
| DBSwapMaxConcurrent |         | >= 0            | 0 = disable db swap |
| DBSwapTimeLimit     |         | seconds         | once per timespan   |
| DBSwapSimple        | false   | true/false      | IFCDEV-83           |


## global setup parameter

| key                 | default               | possible values | comment                    |
|:--------------------|:----------------------|:----------------|:---------------------------|
| HealthReport        | false                 | true/false      |                            |
| HealthServer        | http://127.0.0.1:8443 | URL             | eg. http://127.0.0.1:8843  |
| AccessibilityReport | false                 | true/false      |                            |
| AccessibilityServer | http://127.0.0.1:5443 | URL             | eg. https://127.0.0.1:5443 |
| Service             |                       | string          | eg. io.protel.v5_fidserv   |
| Peer                |                       | string          | eg. fidserv-local          |
| APIServerFallback   |                       | URL             | eg. https://127.0.0.1:7000 |
| TelemetryReport     | false                 | true/false      |                            |
| TelemetryTarget     | io.protel.esb         | string          | optional                   |

See [Hugin](https://github.com/weareplanet/projects-IFCV5/repos/tools/browse/hugin) for `Health` settings.  
See [Interface Name Service](https://github.com/weareplanet/projects-IFCV5/repos/tools/browse/ins) for `Accessibility` settings.  
  


### set global config value

```bash
pman --host prod set robobar-instance-1 config 0 DBSwapMaxConcurrent 3
````

### set global setup value

```bash
pman --host prod set bartech-instance-1 setup 0 APIServerFallback https://127.0.0.1:7000
```
  
See [pman](https://github.com/weareplanet/projects-IFCV5/repos/tools/browse/pman) for how to set a value.  

## additional subscribe messages

module `subscribe.inbound`  
- for `InterestedInInfos` 
- messages from pms to ifc
- is created automatically by the ifc
- can be overwritten by an appropriate mapping

module `subscribe.outbound`  
- for `ProvisionInfos`  
- messages from ifc to pms
- must be created manually

### global mapping

Should be in the directory `config/driver/mapping/global`.  

```json
{
    "driver": "driver",
    "module": "subscribe.inbound",
    "station": 0,
    "mapping":
    [
        {
            "key": "Name1",
            "value": "Field1, Field2, Field3"
        },
        {
            "key": "Name2",
            "value": "Field1, Field2, Field3"
        }
    ]
}
```

### station mapping

Should be in the directory `config/driver/mapping`.  

```json
{
    "driver": "driver",
    "module": "subscribe.inbound",
    "station": 1234567890,
    "mapping":
    [
        {
            "key": "Name3",
            "value": "Field1, Field2, Field3"
        },
        {
            "key": "Name2",
            "value": "Field2"
        },
        {
            "key": "Name1",
            "value": "Field1"
        }
    ]
}
```

both mappings are merged together  


```xml
<InterestedInInfo xmlns="http://protel.io/soap" MessageName="Name1">
    <Field xmlns="http://protel.io/soap">Field1</Field>
</InterestedInInfo>
<InterestedInInfo xmlns="http://protel.io/soap" MessageName="Name2">
    <Field xmlns="http://protel.io/soap">Field2</Field>
</InterestedInInfo>
<InterestedInInfo xmlns="http://protel.io/soap" MessageName="Name3">
    <Field xmlns="http://protel.io/soap">Field1</Field>
    <Field xmlns="http://protel.io/soap">Field2</Field>
    <Field xmlns="http://protel.io/soap">Field3</Field>
</InterestedInInfo>
```

See [protel and Indra TMS - record description](https://docs.google.com/document/d/1XP55GCUTzNWpf2V8_ZtwjrWS6w1kfo3r-gwDUdXsQoQ/edit#heading=h.gjdgxs) for possible `MessageName` and `Field` names.  
