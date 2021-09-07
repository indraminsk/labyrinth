package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"greenjade/config"
	"io/ioutil"
	"os"
	"testing"
)

func fetchJsonData(t *testing.T, pathToJson string) (level LevelType, status error) {
	var (
		err error
		jsonFile   *os.File
		jsonData []byte
	)

	jsonFile, err = os.Open(pathToJson)
	if err != nil {
		return level, errors.New(fmt.Sprint("[error] can't open file:", pathToJson))
	}

	defer func() {
		if err = jsonFile.Close(); err != nil {
			fmt.Println("[error] clear memory file")
		}
	}()

	jsonData, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		return level, errors.New("[error] read json data")
	}

	err = json.Unmarshal(jsonData, &level)
	if err != nil {
		return level, errors.New("[error] unmarshal json data")
	}


	return level, status
}

func TestValidateDataAllOk(t *testing.T) {
	var (
		err, status error

		cfg *config.Config
		level LevelType
	)

	cfg = config.Conf("../")

	level, err = fetchJsonData(t, "../testdata/data_all_ok_2.json")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	status = level.Validate(cfg.Constraints)
	if status != nil {
		t.Error(status.Error())
	}
}

func TestValidateDataPointValueLessThanMin(t *testing.T) {
	var (
		err, status error

		cfg *config.Config
		level LevelType
	)

	cfg = config.Conf("../")

	level, err = fetchJsonData(t, "../testdata/data_point_value_less_than_min.json")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	status = level.Validate(cfg.Constraints)
	if status == nil {
		t.Error("unexpected success")
	}
}

func TestValidateDataPointValueGreaterThanMax(t *testing.T) {
	var (
		err, status error

		cfg *config.Config
		level LevelType
	)

	cfg = config.Conf("../")

	level, err = fetchJsonData(t, "../testdata/data_point_value_greater_than_max.json")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	status = level.Validate(cfg.Constraints)
	if status == nil {
		t.Error("unexpected success")
	}
}

func TestValidateDataNotRectangle(t *testing.T) {
	var (
		err, status error

		cfg *config.Config
		level LevelType
	)

	cfg = config.Conf("../")

	level, err = fetchJsonData(t, "../testdata/data_not_rectangle.json")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	status = level.Validate(cfg.Constraints)
	if status == nil {
		t.Error("unexpected success")
	}
}

func TestValidateDataTooManyX(t *testing.T) {
	var (
		err, status error

		cfg *config.Config
		level LevelType
	)

	cfg = config.Conf("../")

	level, err = fetchJsonData(t, "../testdata/data_too_many_x.json")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	status = level.Validate(cfg.Constraints)
	if status == nil {
		t.Error("unexpected success")
	}
}

func TestValidateDataTooManyY(t *testing.T) {
	var (
		err, status error

		cfg *config.Config
		level LevelType
	)

	cfg = config.Conf("../")

	level, err = fetchJsonData(t, "../testdata/data_too_many_y.json")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	status = level.Validate(cfg.Constraints)
	if status == nil {
		t.Error("unexpected success")
	}
}