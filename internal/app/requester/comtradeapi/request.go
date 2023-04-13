package comtradeapi

import (
	"fmt"
	"proj/internal/app/server"
	"reflect"
	"strings"
)

type ApiRequest server.Request

func (r ApiRequest) Url(apiUrl string) string {
	builder := strings.Builder{}
	builder.WriteString(apiUrl)
	builder.WriteString(fmt.Sprintf("%s/%s/%s?", r.TypeCode, r.FreqCode, r.ClCode))

	v := reflect.ValueOf(r)
	for i := 3; i < v.NumField(); i++ {
		fieldVal := v.Field(i).String()
		fieldName := v.Type().Field(i).Name
		if fieldVal != "" {
			builder.WriteString(fmt.Sprintf("%s=%s&", fieldName, fieldVal))
		}
	}
	return builder.String()
}

// TODO: fix split overhead

func (r ApiRequest) SplitByPartner2(splitSize int) []ApiRequest {
	requests := make([]ApiRequest, 0, splitSize)
	builders := make([]strings.Builder, splitSize)
	cnt := 0

	if r.Partner2Code == "" {
		for code := range CountryName {
			idx := cnt % splitSize
			builders[idx].WriteString(fmt.Sprintf("%d,", code))
			cnt++
		}
	} else {
		partners := strings.Split(r.Partner2Code, ",")
		for _, code := range partners {
			idx := cnt % splitSize
			builders[idx].WriteString(fmt.Sprintf("%s,", code))
			cnt++
		}
	}

	for _, partner := range builders {
		req := r
		partnerStr := partner.String()
		if len(partnerStr) == 0 {
			continue
		}
		req.Partner2Code = partnerStr[:len(partnerStr)-1]
		requests = append(requests, req)
	}
	return requests
}

func (r ApiRequest) SplitByReporter(splitSize int) []ApiRequest {
	requests := make([]ApiRequest, 0, splitSize)
	builders := make([]strings.Builder, splitSize)
	cnt := 0

	if r.ReporterCode == "" {
		for code := range CountryName {
			idx := cnt % splitSize
			builders[idx].WriteString(fmt.Sprintf("%d,", code))
			cnt++
		}
	} else {
		reporters := strings.Split(r.ReporterCode, ",")
		for _, code := range reporters {
			idx := cnt % splitSize
			builders[idx].WriteString(fmt.Sprintf("%s,", code))
			cnt++
		}
	}

	for _, reporter := range builders {
		req := r
		reporterStr := reporter.String()
		if len(reporterStr) == 0 {
			continue
		}
		req.ReporterCode = reporterStr[:len(reporterStr)-1]
		requests = append(requests, req)
	}
	return requests
}
