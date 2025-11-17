## config parameter

| key                  | default | possible values | comment                                    |
|:---------------------|:--------|:----------------|:-------------------------------------------|
| DriverEncoding       |         |                 | if not configured then UTF-8 is used       |
| KeyAnswerTimeout     | 30      |                 | vendor device timeout in seconds           |
| Password             | empty   | max 7 chars     |                                            |
| PasswordFromUser     | false   | true/false      | use UserID as password                     |
| PassNumberAlwaysNull | false   | true/false      | ignore key options (PassNumberOption)      |
| KeyLevel             | 1       | 1-4             |                                            |
| KeyDelete            | true    | true/false      | send key delete                            |
| Protocol             | 1       | 1/2             | 1 = Short room length(5), 2 = Extended(15) |
| LeadingZeroes        | false   | true/false      | keep leading zeroes or not                 |

### key level

* 1 = guest
* 2 = connector
* 3 = multi-connector
* 4 = limited-use

