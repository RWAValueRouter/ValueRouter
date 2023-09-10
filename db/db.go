package db

import (
	"container/list"
	"database/sql"
	"errors"
	"github.com/RWAValueRouter/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/common"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Conn *Dialect

type Dialect struct {
	db              *sql.DB
	driverSource    string
	maxConns        int
	maxIdles        int
	connMaxLifeTime time.Duration
}

func NewInstance() *Dialect {
	var dia = new(Dialect)
	err := dia.Create(common.Conf.DbConfig.DbDriverSource)
	if err != nil {
		log.Error("init db error: ", " error_message", err.Error())
		os.Exit(0)
	}
	return dia
}

// Create create a db instance .
func (dia *Dialect) Create(driverSource string) error {
	if driverSource == "" {
		return errors.New("driver source can not be blank")
	}
	db, err := sql.Open(common.Conf.DbConfig.DbDriverName,
		driverSource)
	if err != nil {
		return err
	}
	if dia.maxConns == 0 {
		db.SetMaxOpenConns(10)
	} else {
		db.SetMaxOpenConns(dia.maxConns)
	}
	if dia.maxIdles == 0 {
		db.SetMaxIdleConns(5)
	} else {
		db.SetMaxIdleConns(dia.maxIdles)
	}
	if dia.connMaxLifeTime == 0 {
		db.SetConnMaxLifetime(time.Second * 14440)
	} else {
		db.SetConnMaxLifetime(dia.connMaxLifeTime)
	}
	dia.driverSource = driverSource
	dia.db = db
	return nil
}

// IsConnected check if a db connection is still alive
func (dia *Dialect) IsConnected() error {
	err := dia.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// Begin begin a transaction
func (dia *Dialect) Begin() (tx *sql.Tx, err error) {
	tx, err = dia.db.Begin()
	return
}

// Commit commit a transaction
func (dia *Dialect) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

// Rollback rollback a transaction
func (dia *Dialect) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

// Execute data manipulate language already commited no need to commit again
func (dia *Dialect) CommitOneRow(sql string, args ...interface{}) (int64, error) {
	logStr := sql
	for _, v := range args {
		if _, ok := v.(string); ok {
			logStr = strings.Replace(logStr, "?", v.(string), 1)
		}
		if _, ok := v.(int); ok {
			logStr = strings.Replace(logStr, "?", strconv.Itoa(v.(int)), 1)
		}
	}
	log.Info(logStr)
	stmt, err := dia.db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if id == 0 {
		id, _ = result.RowsAffected()
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

// BatchExecute batch data manipulate , need use begin to get a tx and rollback or commit this tx
func BatchExecute(sql string, tx *sql.Tx, params ...interface{}) (int64, error) {
	logStr := sql
	for _, v := range params {
		if _, ok := v.(string); ok {
			logStr = strings.Replace(logStr, "?", v.(string), 1)
		}
		if _, ok := v.(int); ok {
			logStr = strings.Replace(logStr, "?", strconv.Itoa(v.(int)), 1)
		}
	}
	log.Info(logStr)
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(params...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if id == 0 {
		id, _ = result.RowsAffected()
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (dia *Dialect) Query(s string, params ...interface{}) (*list.List, error) {
	logStr := s
	for _, v := range params {
		if _, ok := v.(string); ok {
			logStr = strings.Replace(logStr, "?", v.(string), 1)
		}
		if _, ok := v.(int); ok {
			logStr = strings.Replace(logStr, "?", strconv.Itoa(v.(int)), 1)
		}
	}
	log.Info(logStr)
	rows, err := dia.db.Query(s, params...)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([][]byte, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	retList := list.New()

	// Fetch rows
	for rows.Next() {
		retMap := make(map[string]interface{})
		retList.PushBack(retMap)
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				// skip nil value
				continue
			}
			retMap[columns[i]] = string(col)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retList, nil
}

// GetStructValue convert result to struct , not support boolean type, true in mysql is 1 , false is 0
// struct field need to be exported
func (dia *Dialect) GetStructValue(sql string, strct interface{}, params ...interface{}) ([]interface{}, error) {
	l, err := dia.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	i := 0
	r := make([]interface{}, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		v := reflect.New(reflect.TypeOf(strct))
		src := e.Value.(map[string]interface{})
		vi := v.Interface()
		common.Map2Struct(src, vi)
		r[i] = reflect.ValueOf(vi).Interface()
		i++
	}
	return r, nil
}

func (dia *Dialect) GetIntValue(sql string, params ...interface{}) (int, error) {

	l, err := dia.Query(sql, params...)
	if err != nil || l.Len() == 0 {
		return -1, err
	}
	m := l.Front().Value.(map[string]interface{})
	for _, v := range m {
		_, ok := v.(string)
		if !ok {
			return -1, errors.New("not string type")
		}
		return strconv.Atoi(v.(string))
	}

	return -1, err
}

func (dia *Dialect) GetFloatValue(sql string, params ...interface{}) (float64, error) {

	l, err := dia.Query(sql, params...)
	if err != nil || l.Len() == 0 {
		return -1, err
	}
	m := l.Front().Value.(map[string]interface{})
	for _, v := range m {
		return strconv.ParseFloat(v.(string), 64)
	}

	return -1, err
}

// GetStringValue get string value out of the return value , if it is a int it will convert to a string
func (dia *Dialect) GetStringValue(sql string, params ...interface{}) (string, error) {

	l, err := dia.Query(sql, params...)
	if err != nil || l.Len() == 0 {
		return "", err
	}
	m := l.Front().Value.(map[string]interface{})
	for _, v := range m {
		return v.(string), nil
	}

	return "", err
}

func (dia *Dialect) Close() error {
	return dia.db.Close()
}

func (dia *Dialect) SetMaxOpenConnections(maxConn int) error {
	if maxConn <= 0 {
		return errors.New("maxConn must bigger than 0")
	}
	dia.maxConns = maxConn
	return nil
}

func (dia *Dialect) SetMaxIdles(maxIdles int) error {
	if maxIdles <= 0 {
		return errors.New("maxIdles must bigger than 0")
	}
	dia.maxIdles = maxIdles
	return nil
}

func (dia *Dialect) SetConnMaxLifeTime(connMaxLifeTime int) error {
	if connMaxLifeTime <= 0 {
		return errors.New("connMaxLifeTime must bigger than 0")
	}
	dia.maxIdles = connMaxLifeTime
	return nil
}

func Init() {
	Conn = NewInstance()
}
