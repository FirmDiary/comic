package main

import (
	"comic/common"
	"comic/datamodels"
	"fmt"
	"testing"
)

func TestDb(t *testing.T)  {
	dbEngine :=  common.NewDbEngine()
	datalist := make([]datamodels.Upload, 0)

	err := dbEngine.Desc("id").Find(&datalist)

	if err != nil {
		panic(err)
	}
	for i := range datalist{
		fmt.Printf("%v \n", datalist[i])
	}
}
