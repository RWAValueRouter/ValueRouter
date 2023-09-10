package jobs

import (
	"github.com/RWAValueRouter/ValueRouter/common"
	"github.com/RWAValueRouter/ValueRouter/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type ListenConfigSuite struct {
	suite.Suite
}

func (suite *ListenConfigSuite) SetupTest() {
	common.Conf.DbConfig.DbDriverSource = "gotest:12345678@tcp(127.0.0.1:3306)/smw"
	common.Conf.DbConfig.DbDriverName = "mysql"
	db.Init()
	jobs.Stop()
}

func (suite *ListenConfigSuite) TestListenConfig() {
	if db.Conn.IsConnected() != nil {
		return
	}
	exist, err := db.Conn.GetIntValue("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'smw'  AND table_name = 'chain_config'")
	var affected int64
	if assert.NoError(suite.T(), err) {
		if exist == 0 {
			affected, err = db.Conn.CommitOneRow("CREATE TABLE `chain_config` (\n  `id` int NOT NULL AUTO_INCREMENT,\n  `rpc_list` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'rpc endpoints separated with |',\n  `chain_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '链名称',\n  `chain_currency` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '链的原生代币名称',\n  `is_evm` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '是否是evm链',\n  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0:normal ,-1: invalid',\n  `local_system_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `chain_type` int DEFAULT '0' COMMENT 'chain type defined in service',\n  `chain_id` int DEFAULT '0' COMMENT '只有是evm链时候存在chainId，其他链默认是0',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=100002 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin")
			if assert.NoError(suite.T(), err) {
				assert.True(suite.T(), affected == 0)
			}
		}
	}
	_, err = db.Conn.CommitOneRow("insert into chain_config(id, rpc_list, chain_name, chain_currency, is_evm, chain_id, chain_type) values(100000,'https://rpc.ankr.com/fantom_testnet_test_case', 'Fantom Testnet special', 'FTM_TEST', 'YES', 40021111, 0)")
	assert.NoError(suite.T(), err)
	_, err = db.Conn.CommitOneRow("insert into chain_config(id, rpc_list, chain_name, chain_currency, is_evm, chain_id, chain_type) values(100001,'https://rpc.ankr.com/FNMtom_testnet_test_case', 'FFNMtom Testnet special', 'FNM_TEST', 'YES', 40021222, 1)")
	assert.NoError(suite.T(), err)
	listenConfig()
	time.Sleep(5 * time.Second)
	var wait sync.WaitGroup
	for i := 0; i < 100; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			Cm.RLock()
			defer Cm.RUnlock()
			v, ok := CC[0]
			assert.True(suite.T(), ok)
			vv, ok := v[40021111]
			assert.True(suite.T(), ok)
			assert.Equal(suite.T(), "Fantom Testnet special", vv.Chain_name)

			v, ok = CC[1]
			assert.True(suite.T(), ok)
			vv, ok = v[40021222]
			assert.True(suite.T(), ok)
			assert.Equal(suite.T(), "FFNMtom Testnet special", vv.Chain_name)
		}()
	}
	wait.Wait()

	affected, err = db.Conn.CommitOneRow("delete from chain_config where id in (?,?)", 100000, 100001)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), 2, int(affected))
	}

}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ListenConfigSuite))
}
