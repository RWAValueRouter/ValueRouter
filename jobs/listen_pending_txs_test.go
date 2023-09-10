package jobs

import (
	"github.com/RWAValueRouter/ValueRouter/common"
	"github.com/RWAValueRouter/ValueRouter/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ListenPendingTxsSuite struct {
	suite.Suite
}

func (suite *ListenPendingTxsSuite) SetupTest() {
	common.Conf.DbConfig.DbDriverSource = "gotest:12345678@tcp(127.0.0.1:3306)/smw"
	common.Conf.DbConfig.DbDriverName = "mysql"
	db.Init()
	jobs.Stop()
}

func (suite *ListenPendingTxsSuite) TestListenPendingTxs() {
	if db.Conn.IsConnected() != nil {
		return
	}
	exist, err := db.Conn.GetIntValue("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'smw'  AND table_name = 'signs_detail'")
	var affected int64
	if assert.NoError(suite.T(), err) {
		if exist == 0 {
			affected, err = db.Conn.CommitOneRow("CREATE TABLE `signs_detail` (\n  `id` int NOT NULL AUTO_INCREMENT,\n  `key_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc sign key ID',\n  `user_account` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户account',\n  `group_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '组ID',\n  `enode` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '节点对应的enode',\n  `threshold` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '门限制',\n  `msg_hash` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'msg signed hash values , using | to separate multiple hashes',\n  `msg_context` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'msg_context list separated by |',\n  `rsv` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT 'rsv list separated by |',\n  `public_key` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc address public key',\n  `mpc_address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'mpc address',\n  `reply_initializer` tinyint DEFAULT NULL COMMENT '0:not initializer ,1: initializer',\n  `reply_status` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'reply status of creating mpc wallet',\n  `reply_timestamp` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'reply timestamp',\n  `reply_enode` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'reply enode',\n  `initiator_public_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'initial public key',\n  `key_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'key类型',\n  `mode` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'Mode模式',\n  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0:pending , 1 SUCCESS , 2 FAIL, 3 Timeout',\n  `error` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'error message',\n  `tip` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'tip message',\n  `local_system_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `ip_port` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'gid对应用户对应的ipport地址',\n  `txid` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'signing tx id',\n  `chain_type` int DEFAULT '0' COMMENT 'chain type',\n  `chain_id` int DEFAULT '0' COMMENT 'chain id', \n `sign_type` int DEFAULT '0' COMMENT 'sign type', \n PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin")
			if assert.NoError(suite.T(), err) {
				assert.True(suite.T(), affected == 0)
			}
		}
	}
	_, err = db.Conn.CommitOneRow("INSERT INTO `signs_detail` VALUES (12000001,'0xd9fc83b9d5149c03ec8067cd76121221c864771e383c68a6f1d1ae79154c7633','0xd17831dd9db4ce9a8d331c807329e93015ca2bcb','1f29316aef29dad0f2b1a4a5a53943318a65df56116a836200f8ed45f70c2f1b2965e1ac53d2c938be530af74c6f19e60bd612524cdc14264af48f7776c87aeb','enode://08ba43b0715bb27e03911592d3fed49a22f49ecaf1628d44c4c7d4a8914b86d423716d7316d98a80f4d280b852a33ba7972f0a744e386375e6d58a12fab96752@127.0.0.1:30825','2/2','0x1db35952a0b9a883b0b0beb91e54c2c15265bdb1cf00a69cb1d81779a8bc6985','{\\\"from\\\":\\\"0x5cdf2041218AB47407EcfCAd357FEEd5175C7259\\\",\\\"to\\\":\\\"0xa3eDE283D61d5f100a95099D78fd07566bB6c7F0\\\",\\\"chainId\\\":\\\"0x5\\\",\\\"value\\\":\\\"0x16BCC41E90000\\\",\\\"nonce\\\":1,\\\"gas\\\":48016,\\\"gasPrice\\\":67727197933,\\\"data\\\":\\\"\\\",\\\"originValue\\\":\\\"0.0004\\\",\\\"name\\\":\\\"ETH Goerli\\\"}','2EF54A55295A3A6C67B320235DA515C65C75AFE717966A0B23F14F237FC6CF744C892A6AFE9A23276C3EB044FB87659CEA1F1CE29604203324EC73292E15BEFF00','04f32c7733eea8e4bc236dbc26346432fdef6e11bd6bc4e8c9abd306b3e1c727b3e3fc0fe96bd82bcdfe043ee5083cf00c0c0f7d56b114fa8e309cfb9d2846ad3a','0x5cdf2041218AB47407EcfCAd357FEEd5175C7259',0,'AGREE','1678781066993','08ba43b0715bb27e03911592d3fed49a22f49ecaf1628d44c4c7d4a8914b86d423716d7316d98a80f4d280b852a33ba7972f0a744e386375e6d58a12fab96752','2e2b74160a62114e8901668022ab8df0d30ae9c69a48100ab70d50da4713ca6d71ca1bee30bd60a505077ea1c1c2b67b423ed75d535599c3be2b46f397de1a96','EC256K1','2',1,'','','2023-03-14 08:03:09','127.0.0.1:3792','0x8439fe6676bb6f262bcfc5bce53459060110e5c58eaad59b0732cdcb301e9a88',0,5,0)")
	assert.NoError(suite.T(), err)
	CC[0] = map[int]*ChainConfig{
		5: {
			Rpc_list:       []string{"https://eth-goerli.public.blastapi.io"},
			Chain_name:     "goerli testnet",
			Chain_currency: "ETH",
			Is_evm:         "YES",
			Chain_id:       5,
			Chain_type:     0,
		},
	}
	listenTransactions()
	time.Sleep(5 * time.Second)
	_, err = db.Conn.CommitOneRow("delete from `signs_detail` where id = ?", 12000001)
	assert.NoError(suite.T(), err)
}

func TestListenPendingTxs(t *testing.T) {
	suite.Run(t, new(ListenPendingTxsSuite))
}
