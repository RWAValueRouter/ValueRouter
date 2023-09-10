package jobs

type ReqAddressStatus int

var ReqAddressStatusMap = make(map[string]ReqAddressStatus)

const (
	Pending = iota
	Success
	Failure
	Timeout
)

func (bp ReqAddressStatus) String() string {
	return []string{"Pending", "Success", "Failure", "Timeout"}[bp]
}

func init() {
	ReqAddressStatusMap["pending"] = Pending
	ReqAddressStatusMap["success"] = Success
	ReqAddressStatusMap["failure"] = Failure
	ReqAddressStatusMap["timeout"] = Timeout
}
