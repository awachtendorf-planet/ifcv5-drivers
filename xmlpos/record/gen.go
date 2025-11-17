package vendor

//go:generate xsdgen -pkg xsd -o record.go xmlpos.xsd

/*

uncomment the generated function in record.go because we use our own version (date.go) to match date and time format from xmlpos protocol

type xsdDate time.Time
type xsdTime time.Time

*/
