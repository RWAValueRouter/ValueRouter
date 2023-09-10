package common

type dataResult struct {
	Result string `json:"result"`
}

type response struct {
	Status string      `json:"Status"`
	Tip    string      `json:"Tip"`
	Error  string      `json:"Error"`
	Data   interface{} `json:"Data"`
}
