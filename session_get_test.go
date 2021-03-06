// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetVar(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type GetVar struct {
		Id      int64  `xorm:"autoincr pk"`
		Msg     string `xorm:"varchar(255)"`
		Age     int
		Money   float32
		Created time.Time `xorm:"created"`
	}

	assert.NoError(t, testEngine.Sync2(new(GetVar)))

	var data = GetVar{
		Msg:   "hi",
		Age:   28,
		Money: 1.5,
	}
	_, err := testEngine.InsertOne(data)
	assert.NoError(t, err)

	var msg string
	has, err := testEngine.Table("get_var").Cols("msg").Get(&msg)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, "hi", msg)

	var age int
	has, err = testEngine.Table("get_var").Cols("age").Get(&age)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, 28, age)

	var money float64
	has, err = testEngine.Table("get_var").Cols("money").Get(&money)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, "1.5", fmt.Sprintf("%.1f", money))

	var valuesString = make(map[string]string)
	has, err = testEngine.Table("get_var").Get(&valuesString)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, 5, len(valuesString))
	assert.Equal(t, "1", valuesString["id"])
	assert.Equal(t, "hi", valuesString["msg"])
	assert.Equal(t, "28", valuesString["age"])
	assert.Equal(t, "1.5", valuesString["money"])

	var valuesInter = make(map[string]interface{})
	has, err = testEngine.Table("get_var").Where("id = ?", 1).Select("*").Get(&valuesInter)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, 5, len(valuesInter))
	assert.EqualValues(t, 1, valuesInter["id"])
	assert.Equal(t, "hi", fmt.Sprintf("%s", valuesInter["msg"]))
	assert.EqualValues(t, 28, valuesInter["age"])
	assert.Equal(t, "1.5", fmt.Sprintf("%v", valuesInter["money"]))

	var valuesSliceString = make([]string, 5)
	has, err = testEngine.Table("get_var").Get(&valuesSliceString)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, "1", valuesSliceString[0])
	assert.Equal(t, "hi", valuesSliceString[1])
	assert.Equal(t, "28", valuesSliceString[2])
	assert.Equal(t, "1.5", valuesSliceString[3])

	var valuesSliceInter = make([]interface{}, 5)
	has, err = testEngine.Table("get_var").Get(&valuesSliceInter)
	assert.NoError(t, err)
	assert.Equal(t, true, has)

	v1, err := convertInt(valuesSliceInter[0])
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v1)

	assert.Equal(t, "hi", fmt.Sprintf("%s", valuesSliceInter[1]))

	v3, err := convertInt(valuesSliceInter[2])
	assert.NoError(t, err)
	assert.EqualValues(t, 28, v3)

	v4, err := convertFloat(valuesSliceInter[3])
	assert.NoError(t, err)
	assert.Equal(t, "1.5", fmt.Sprintf("%v", v4))
}

func convertFloat(v interface{}) (float64, error) {
	switch v.(type) {
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	case string:
		i, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	case []byte:
		i, err := strconv.ParseFloat(string(v.([]byte)), 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("unsupported type: %v", v)
}

func convertInt(v interface{}) (int64, error) {
	switch v.(type) {
	case int:
		return int64(v.(int)), nil
	case int8:
		return int64(v.(int8)), nil
	case int16:
		return int64(v.(int16)), nil
	case int32:
		return int64(v.(int32)), nil
	case int64:
		return v.(int64), nil
	case []byte:
		i, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	case string:
		i, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("unsupported type: %v", v)
}
