package service

import (
	"github.com/anyswap/FastMulThreshold-DSA/log"
)

type ServiceMiddleWare struct{}

func (service *ServiceMiddleWare) GetGroupId(threshold string, userAccountsAndIpPortAddr []string) map[string]interface{} {
	if data, err := getGroupId(threshold, userAccountsAndIpPortAddr); err != nil {
		log.Error("getGroupId", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetGroupIdByRawData(raw string) map[string]interface{} {
	if data, err := getGroupIdByRawData(raw); err != nil {
		log.Error("GetGroupIdByRawData", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) KeyGen(rsv string, msg string) map[string]interface{} {
	if data, err := doKeyGen(rsv, msg); err != nil {
		log.Error("KeyGen", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) KeyGenByRawData(raw string) map[string]interface{} {
	if data, err := doKeyGenByRawData(raw); err != nil {
		log.Error("KeyGenByRawData", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetReqAddrStatus(keyId string) map[string]interface{} {
	if data, err := getReqAddrStatus(keyId); err != nil {
		log.Error("getReqAddrStatus", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetAccountList(userAccount string) map[string]interface{} {
	if data, err := getAccountList(userAccount); err != nil {
		log.Error("getAccountList", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) Sign(rsv string, msg string) map[string]interface{} {
	if data, err := doSign(rsv, msg); err != nil {
		log.Error("doSign", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetUnsigedTransactionHash(unsignedTx string, chain int) map[string]interface{} {
	if data, err := getUnsigedTransactionHash(unsignedTx, chain); err != nil {
		log.Error("getUnsigedTransactionHash", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetApprovalList(userAccount string) map[string]interface{} {
	if data, err := getApprovalList(userAccount); err != nil {
		log.Error("getApprovalList", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetApprovalListByPagination(userAccount string, status, page, pageSize int) map[string]interface{} {
	if data, err := getApprovalListByPagination(userAccount, status, page, pageSize); err != nil {
		log.Error("getApprovalList", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetSignHistoryByPagination(userAccount string, page, pageSize int) map[string]interface{} {
	if data, err := getSignHistoryByPagination(userAccount, page, pageSize); err != nil {
		log.Error("getSignHistory", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetSignHistory(userAccount string) map[string]interface{} {
	if data, err := getSignHistory(userAccount); err != nil {
		log.Error("getSignHistory", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) AcceptSign(rsv string, msg string) map[string]interface{} {
	if data, err := acceptSign(rsv, msg); err != nil {
		log.Error("acceptSign", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetTxHashByKeyId(keyId string) map[string]interface{} {
	if data, err := getTxHashByKeyId(keyId); err != nil {
		log.Error("getTxHashByKeyId", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetTxStatusByKeyId(keyId string) map[string]interface{} {
	if data, err := getTxStatusByKeyId(keyId); err != nil {
		log.Error("getTxStatusByKeyId", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) AddAsset(rsv string, msg string) map[string]interface{} {
	if data, err := doAddAsset(rsv, msg); err != nil {
		log.Error("doAddAsset", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) AddAssetForMpcAddress(rsv string, msg string) map[string]interface{} {
	if data, err := doAddAssetForMpcAddress(rsv, msg); err != nil {
		log.Error("AddAssetForMpcAddress", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetMpcAddressDetail(mpc_account string) map[string]interface{} {
	if data, err := doGetMpcAddressDetail(mpc_account); err != nil {
		log.Error("doGetMpcAddressDetail", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetAsset(account string, chain_id, chain_type int) map[string]interface{} {
	if data, err := doGetAsset(account, chain_id, chain_type); err != nil {
		log.Error("doGetAsset", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetApprovalListByKeyId(key_id string) map[string]interface{} {
	if data, err := doGetApprovalListByKeyId(key_id); err != nil {
		log.Error("doGetApprovalListByKeyId", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetNodesNumber() map[string]interface{} {
	if data, err := doGetNodesNumber(); err != nil {
		log.Error("doGetNodesNumber", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}

func (service *ServiceMiddleWare) GetLatestMpcAddressStatus(mpc_address string, chain_id, chain_type int) map[string]interface{} {
	if data, err := doGetLatestMpcAddressStatus(mpc_address, chain_id, chain_type); err != nil {
		log.Error("doGetNodesNumber", "error", err.Error())
		return map[string]interface{}{
			"Status": "error",
			"Tip":    "something unexpected happen",
			"Error":  err.Error(),
			"Data":   "",
		}
	} else {
		return map[string]interface{}{
			"Status": "success",
			"Tip":    "",
			"Error":  "",
			"Data":   data,
		}
	}
}
