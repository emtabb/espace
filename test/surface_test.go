package test

import (
    "fmt"
)
import "testing"
import . "github.com/emtabb/espace"
import . "github.com/emtabb/espace/surface"
import . "github.com/emtabb/espace/api/element"

func TestSurface(t *testing.T) {
    var st ESpace = Surface()
    x := st.Init()
    x.CsvSpace("simplePhysics.csv")
    fmt.Println(x.FieldTypes())
    testRow10 := make(map[string] interface{})
    testRow10["Particle"] = "sst"
    testRow10["Coordinate"] = "9"
    testRow10["Vector"] = "18"
    testRow10["Time"] = "10"
    ele := x.Elements()
    for _, e := range ele {
        fmt.Println(e.(*Element).GetProperty())
    }
    k := ele[8].(*Element).GetProperty()
    for name := range k {
        success := k[name] == testRow10[name]
        fmt.Println(success)
    }

    y, _ := x.Group([]string{"Coordinate", "Vector"})
    fmt.Println(y)
    eley := y.Elements()
    for _, e := range eley {
        fmt.Println(e.(*Element).GetProperty())
    }
}

func populateData() ESpace {
    dataPath := "TestMoreDateFile.csv"
    var st ESpace = Surface()
    FinancialSpace := st.Init()
    FinancialSpace.CsvSpace(dataPath)
    return FinancialSpace
}

func TestGetNameFieldsBigData(t *testing.T) {
    space := populateData()
    fmt.Println(space.NameFields())
}

func TestReadBigData(t *testing.T) {
    space := populateData()
    FieldGroup := []string {"Year", "Units", "Variable_code", "Variable_name", "Variable_category", "Value"}
    SomeGroup, _ := space.Group(FieldGroup)
    row, _ := SomeGroup.Shape()
    
    for i := 0; i < row; i++ {
        fmt.Println(i + 1, "|", SomeGroup.Elements()[i].(*Element).ToString())
    }
}

func TestSaveBigDataBigData(t *testing.T) {
    space := populateData()
    FieldGroup := []string {"Year", "Units", "Variable_code", "Variable_name", "Variable_category", "Value"}
    SomeGroup, _ := space.Group(FieldGroup)
    SomeGroup.Save("./Data.csv")
}

func TestJoinBigData(t *testing.T) {
    space := populateData()
    dataJoin := space.Fields([]string{"Industry_code_ANZSIC06"})
    dataJoin[0].SetType("STRING")
    FieldGroup := []string {"Year", "Units", "Variable_code", "Variable_name", "Variable_category", "Value"}
    SomeGroup, _ := space.Group(FieldGroup)
    afterJoin, _ := SomeGroup.Join(dataJoin[0])
    
    row, _ := afterJoin.Shape()
    for i := 0; i < row; i++ {
        fmt.Println(i + 1, "|", afterJoin.Elements()[i].(*Element).ToString())
    }
}

func TestDropBigData(t *testing.T) {
    space := populateData()
    FieldGroup := []string {"Year", "Units", "Variable_code", "Variable_name", "Variable_category", "Value"}
    SomeGroup, _ := space.Group(FieldGroup)
    afterDrop, _ := SomeGroup.Drop([]string{"Variable_name", "Value"})
    
    row, col := afterDrop.Shape()
    fmt.Println(row, col)
}

func TestFieldTypes(t *testing.T) {
    space := populateData()
    fmt.Println(space.FieldTypes())
}