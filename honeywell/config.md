## config parameter

| key      | default | possible values | comment                                        |
|:---------|:--------|:----------------|:-----------------------------------------------|
| Protocol | 1       | 1/2/3           | 1 = Honywell, 2=Alerton PROT1, 3=Alerton PROT2 |


### protocol types

Honeywell  
CI/CO9999CR, no low level ack, wait for replay packet  

Alerton PROT1  
E/V9999?, with LRC, low level ack, no replay packet  

Alerton PROT2  
E/V9999, without LRC, low level ack, no replay packet