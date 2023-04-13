package server

type RequestBatch struct {
	TypeCode string `json:"typeCode"`
	FreqCode string `json:"freqCode"`
	ClCode   string `json:"clCode"`

	ReporterCode []string `json:"reporterCode"`
	Period       []string `json:"period"`
	FlowCode     []string `json:"flowCode"`
	PartnerCode  []string `json:"partnerCode"`
	Partner2Code string   `json:"partner2Code"`
	CmdCode      string   `json:"cmdCode"`
	CustomsCode  string   `json:"customsCode"`
	MotCode      string   `json:"motCode"`
	IncludeDesc  string   `json:"includeDesc"`
}

func DefaultBatch() RequestBatch {
	return RequestBatch{
		TypeCode:     "C",
		FreqCode:     "A",
		ClCode:       "HS",
		ReporterCode: []string{""},
		Period:       nil,
		FlowCode:     nil,
		PartnerCode:  []string{"643"},
		Partner2Code: "0",
		CmdCode:      "",
		CustomsCode:  "C00",
		MotCode:      "0",
		IncludeDesc:  "",
	}
}

func (r *RequestBatch) CalcSize() int {
	return len(r.ReporterCode) * len(r.Period) * len(r.FlowCode) * len(r.PartnerCode)
}

type Request struct {
	TypeCode string
	FreqCode string
	ClCode   string

	ReporterCode string
	Period       string
	FlowCode     string
	PartnerCode  string
	Partner2Code string
	CmdCode      string
	CustomsCode  string
	MotCode      string
	IncludeDesc  string
}

func (r *RequestBatch) SplitBatch() []Request {
	requests := make([]Request, 0, r.CalcSize())
	for _, reporterCode := range r.ReporterCode {
		for _, period := range r.Period {
			for _, flowCode := range r.FlowCode {
				for _, partnerCode := range r.PartnerCode {
					currRequest := Request{
						TypeCode:     r.TypeCode,
						FreqCode:     r.FreqCode,
						ClCode:       r.ClCode,
						ReporterCode: reporterCode,
						Period:       period,
						FlowCode:     flowCode,
						PartnerCode:  partnerCode,
						Partner2Code: r.Partner2Code,
						CmdCode:      r.CmdCode,
						CustomsCode:  r.CustomsCode,
						MotCode:      r.MotCode,
						IncludeDesc:  r.IncludeDesc,
					}
					requests = append(requests, currRequest)
				}
			}
		}
	}
	return requests
}
