package db

import (
	"github.com/RWAValueRouter/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/common"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCommitOneRow(t *testing.T) {
	if Conn.IsConnected() != nil {
		return
	}
	exist, err := Conn.GetIntValue("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'smw'  AND table_name = 'test'")
	var affected int64
	if assert.NoError(t, err) {
		if exist == 0 {
			affected, err = Conn.CommitOneRow("CREATE TABLE `test` (\n  `id` int NOT NULL AUTO_INCREMENT,\n  `Name` varchar(256) COLLATE utf8mb4_bin DEFAULT NULL,\n  `Age` int DEFAULT NULL,\n  `Flag` tinyint(1) DEFAULT NULL,\n  `What` varchar(256) COLLATE utf8mb4_bin DEFAULT NULL,\n  `Floatvalue` float DEFAULT NULL,\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin")
			if assert.NoError(t, err) {
				assert.True(t, affected == 0)
			}
			affected, err = Conn.CommitOneRow("INSERT INTO `test` VALUES (1,'test111',111,1,NULL,NULL),(2,'test111',111,1,NULL,NULL),(3,'test111',111,1,NULL,NULL),(4,'test111',111,1,NULL,NULL),(5,'test111',111,1,NULL,NULL),(6,'test111',111,1,NULL,NULL),(7,'test111',111,1,NULL,NULL),(8,'test111',111,1,NULL,NULL)")
			if assert.NoError(t, err) {
				assert.True(t, affected > 0)
			}
		}
	}
}

func TestDialect_GetStructValue(t *testing.T) {
	type test struct {
		Id         int
		Name       string
		Age        int
		Flag       bool
		What       string
		Floatvalue float64
	}
	if Conn.IsConnected() != nil {
		return
	}
	ret, err := Conn.GetStructValue("select * from test", test{})
	assert.NoError(t, err)
	for _, v := range ret {
		_, ok := v.(*test)
		assert.Equal(t, true, ok)
	}

	ret, err = Conn.GetStructValue("select * from test where id > ? ", test{}, 1)
	assert.NoError(t, err)
	for _, v := range ret {
		_, ok := v.(*test)
		assert.Equal(t, true, ok)
	}
}

func TestDialect_GetStringValue(t *testing.T) {
	if Conn.IsConnected() != nil {
		return
	}
	_, err := Conn.GetStringValue("select name from test where id = ?", 2)
	assert.NoError(t, err)

	_, err = Conn.GetStringValue("select age from test where id = ?", 2)
	assert.NoError(t, err)
}

func TestDialect_GetIntValue(t *testing.T) {
	if Conn.IsConnected() != nil {
		return
	}
	_, err := Conn.GetIntValue("select age from test where id = ?", 2)
	assert.NoError(t, err)

	_, err = Conn.GetIntValue("select name from test where id = ?", 2)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "strconv.Atoi: parsing"))
	}
}

func TestBatchBatch(t *testing.T) {
	if Conn.IsConnected() != nil {
		return
	}
	tx, err := Conn.Begin()
	assert.NoError(t, err)
	affected, err := BatchExecute("insert into test(name,age,flag) values(?,?,?)", tx, []interface{}{"test111", 111, 1}...)
	assert.NoError(t, err)
	log.Info("BatchExecute", "affected", affected)
	err = Conn.Commit(tx)
	assert.NoError(t, err)
}

func TestGetFloatValue(t *testing.T) {
	if Conn.IsConnected() != nil {
		return
	}
	f1, err := Conn.GetFloatValue("select floatvalue from test where id = ?", 1)
	assert.NoError(t, err)
	log.Info("getfloatvalue", "value", f1)

	f1, err = Conn.GetFloatValue("select name from test where id = ?", 1)
	if assert.Error(t, err) {
		assert.True(t, strings.Contains(err.Error(), "strconv.ParseFloat: parsing"))
	}
}

func init() {
	common.Conf.DbConfig.DbDriverSource = "gotest:12345678@tcp(127.0.0.1:3306)/smw"
	common.Conf.DbConfig.DbDriverName = "mysql"
	Init()
}
