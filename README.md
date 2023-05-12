# touka
### Simple service to work with [UN global trade statistics database](https://comtradeplus.un.org/)
Touka is intended to work together with [Midoriya](https://github.com/G0P0Ke/midoriya) (for data visalization) via using same database and communicating via REST API and RabbitMQ.<br><br>

Service was created as course-work for HSE Department of Applied Economics to collect russian foreign trade data using mirror statistics method since this data publication was stopped by the Central Bank of Russian Federation in 2022.

## Usage
1. Build service
2. Run it specifying configuration files location with command flags: `--requesterCfg`, `--storageCfg`, `--serverCfg`. By default they must be stored in `/etc/touka/<requester | server | storage>/production`
3. After that you can use handler `/get` with POST method to request data:

```
$ curl -X POST  "<HOST>:<PORT>/get" -d '{
  "typeCode": "C",
  "freqCode": "M",
  "clCode": "HS",
  "reporterCode": ["156"],
  "period": ["201801"],
  "flowCode": ["X"],
  "partnerCode": ["643"],
  "partner2Code": "0",
  "cmdCode": "TOTAL",
  "customsCode": "C00",
  "motCode": "0",
  "includeDesc": ""
}'

HTTP/1.1 200 OK
```
If server's response is 200, touka has started to handle it. After handling request service will save data in database and notify Midoriya via RabbitMQ exchange. 

### Request structure:
| Field name | Type | Description |
| --- | --- | --- |
| typeCode | string | Type of trade: C for commodities and S for service |
| freqCode | string | Trade frequency: A for annual and M for monthly |
| clCode | string | Trade (IMTS) classifications: HS, SITC, BEC or EBOPS |
| reporterCode | [ ]string | Reporter code |
| period | [ ]string | Year or month in format YYYY for year or YYYYMM for year and month |
| flowCode | [ ]string | Trade flow code |
| partnerCode | [ ]string | Partner code |
| partner2Code | string | Second partner/consignment code |
| cmdCode | string | Commodity code |
| customsCode | string | Customs code |
| motCode | string | Mode of transport code |
| includeDesc | string | Include descriptions of data variables |
