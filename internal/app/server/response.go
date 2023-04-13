package server

import (
	"fmt"
	"reflect"
	"strings"
)

type Report struct {
	CmdCode  string
	FlowCode string

	Period       string
	ReporterCode float64
	PartnerCode  float64
	Partner2Code float64
	PrimaryValue float64
}

type Response struct {
	FlowCode string

	Count float64
	Data  []Report
}

func (r *Response) String() string {
	var builder strings.Builder
	v := reflect.ValueOf(r)

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i).Interface()
		fieldName := v.Type().Field(i).Name
		builder.WriteString(fmt.Sprintf("%s: %v\n", fieldName, fieldVal))
	}
	return builder.String()
}
