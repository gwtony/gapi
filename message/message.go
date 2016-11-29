package message

// All is sample

// Record is record
type Record struct {
	Content     string
	Disabled    int
}

// Request is Request record
type Request struct {
	Name        string
	Type        string
	DomainId    int
	Ttl         int
	Records     []Record
}

// ResponseData is Response data
type ResponseData Request

// ResponseResult is Result in response
type ResponseResult struct {
	Affected    int
	Data        ResponseData
}

// Response is response
type Response struct {
	Result	ResponseResult
}
