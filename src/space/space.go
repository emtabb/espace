package space

import (
	"os"
	"bufio"
	"log"
	"strings"
	"regexp"
	"encoding/csv"
	"errors"
	. "github.com/emtabb/espace"
	. "github.com/emtabb/espace/src/element"
	. "github.com/emtabb/espace/src/constant"
	. "github.com/emtabb/field"
	. "github.com/emtabb/field/src/statistic"
	util "github.com/emtabb/espace/src/space/util"
)

type SpaceStructure struct {
	ESpace
	fieldIndex map[string] int
	numberOfColumns int
	numberOfRows int
	nameFields []string
	isGenerate bool

	inmemField []*Field
	inmemFieldName []string
	StringElements []string
	elements []*Element
	NameSpace string
}

//private Method (s *SpaceStructure)
//=============================================================
func (s *SpaceStructure) generate() {
	s.isGenerate = true
	s.nameFields = make([] string, SIZE_DEFAULT)
	s.fieldIndex = make(map[string] int)
	s.inmemField = make([] *Field, SIZE_DEFAULT)
	s.inmemFieldName = make([] string, SIZE_DEFAULT)
	s.StringElements = make([] string, SIZE_DEFAULT)
	s.elements = make([] *Element, SIZE_DEFAULT)
}

func (s *SpaceStructure) OverRange(numberRows int) bool {
	if numberRows > 20 {
		return true
	}
	return false
}
//public Method (s *SpaceStructure)
//=============================================================

func (s *SpaceStructure) Init() ESpace {
	s.generate()
	return s
}

func (s *SpaceStructure) InitElements(elements []*Element) ESpace {
	if !s.isGenerate {
		s.generate()
	}
	for _, element := range elements {
		s.AppendElement(element)
	}
	return s
}

func (s *SpaceStructure) InitFields(nameFields []string) ESpace {
	if !s.isGenerate {
		s.generate()
	}

	s.nameFields = nameFields[:]
	s.numberOfColumns = len(nameFields)
	
	for i := 0; i < s.numberOfColumns; i++ {
		s.fieldIndex[s.nameFields[i]] = i
	}
	return s
}

func (s *SpaceStructure) InitSpace(space ESpace) ESpace{
	if !s.isGenerate {
		s.generate()
	}
	s.InitFields(space.NameFields()).InitElements(space.Elements())
	return s
}

func (s *SpaceStructure) Element() Element {
	var e Element
	return e
}

func (s *SpaceStructure) Elements() []*Element {
	return s.elements
}

func (s *SpaceStructure) connection(url string) {

}

func (s *SpaceStructure) MongoSpace(ConnectionString string) {

}

func (s *SpaceStructure) CsvSpace(path string) {
	rows, err := s.scanData(path, "")
	s.setupElementByCsv(rows)
	if  err != nil {
        log.Fatal(err)
    }
}

func (s *SpaceStructure) HighEnergySpace(path string) {
	rows, _ := s.scanData(path, "+-------+------------------+")
	s.numberOfColumns = 4
	s.formatToCsv(rows)
	s.setupElementByCsv(rows)
}

func (s *SpaceStructure) formatToCsv(rows int) {
	regexCollectNumber := regexp.MustCompile("[+-]?([0-9]*[.])?[0-9]+")
	for i := 0; i < rows; i++ {
		collectSlice := regexCollectNumber.FindAllString(s.StringElements[i], -1)
		lenCollectionSlice := len(collectSlice)
		NumberColumnsMissing := s.numberOfColumns - lenCollectionSlice
		if NumberColumnsMissing > 0 {
			for NumberAdd := 0; NumberAdd < NumberColumnsMissing; NumberAdd++ {
				collectSlice = append(collectSlice, "0")
			}	
			if lenCollectionSlice == 0 {
				collectSlice = s.nameFields
			}
		}
		s.StringElements[i] = strings.Join(collectSlice, ",")
	}
}

func (s *SpaceStructure) splitData() {

}

func (s *SpaceStructure) setupElementByCsv(rows int) {
	StringFrame := make([][] string, rows)
	for i := 0; i < rows; i++ {
		StringFrame[i] = make([]string, s.numberOfColumns)
		data := csv.NewReader(strings.NewReader(s.StringElements[i]))
		data.Comma = ','
		StringFrame[i], _ = data.Read()
	}

	field := csv.NewReader(strings.NewReader(s.StringElements[0]))
	field.Comma = ','
	s.nameFields, _ = field.Read()
	s.numberOfColumns = len(s.nameFields)
	//s.setFieldsDefault(s.numberOfColumns)
	for i := 0; i < s.numberOfColumns; i++ {
		s.fieldIndex[s.nameFields[i]] = i
	}

	for i := 0; i < len(StringFrame) - 1; i++ {
		dataTemp := util.StringArrayToInterface(StringFrame[i + 1])
		element := new(Element).Init().InitLabel(s.nameFields)
		element.Set(util.MapCsvJson(s.nameFields, dataTemp))
		s.AppendElement(element)
	}
}

/*
*	Read file line by line from the path,
*	Have the condition to get only useful line.
*
*/

func (s *SpaceStructure) scanData(path string, condition string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rows := 0
	isScanNornal := true
	if condition != "" {
		isScanNornal = false
	}
	
	if (isScanNornal) {
		rows, err = s.scanNormal(file)
	} else {
		rows, err = s.scanTable(file, condition)
	}
	return rows, err
}

func (s *SpaceStructure) scanNormal(file *os.File) (int, error) {
	rows := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.StringElements = append(s.StringElements, scanner.Text())
		rows++
	}
	err := scanner.Err()
	return rows, err
}

func (s *SpaceStructure) scanTable(file *os.File, condition string) (int, error) {
	rows := 0
	scanner := bufio.NewScanner(file)
	beginTable := false
	for scanner.Scan() {
		if beginTable == false {
			if strings.Contains(scanner.Text(), condition) {
				beginTable = true
				s.StringElements = append(s.StringElements, scanner.Text())
				rows++
				continue
			}
		}
		if beginTable == true {
			if strings.Contains(scanner.Text(), condition) {
				beginTable = false
				continue
			} 
			s.StringElements = append(s.StringElements, scanner.Text())
			rows++
		}
	}
	err := scanner.Err()
	return rows, err
}

func (s *SpaceStructure) Field() Field {
	var f Field
	return f
}

func (s *SpaceStructure) field(fieldName string) *Field {
	if !ContainString(fieldName, s.inmemFieldName) {
		s.inmemFieldName = append(s.inmemFieldName, fieldName)
		immemField := new(Field).Init().Name(fieldName)
		fieldData := make([]interface {}, s.numberOfRows)
		for i, element := range s.Elements() {
			fieldData[i] = element.Field(fieldName)
		}
		immemField.Data(fieldData)
		s.inmemField = append(s.inmemField, immemField)
	}
	return s.inmemField[util.FindPositionArray(fieldName, s.inmemFieldName)]
}

func (s *SpaceStructure) Float(fieldName string) []float64 {
	return s.field(fieldName).Double()
}

func (s *SpaceStructure) Shape() (int, int) {
	return s.numberOfRows, s.numberOfColumns
}

func (s *SpaceStructure) getFieldDefault() string {
	return "FIELD"
}

func (s *SpaceStructure) Head() []*Element {
	elements := make([]*Element, SIZE_DEFAULT);
	for i := 0; i < 5; i++ {
		elements = append(elements, s.Elements()[i]);
	}
	return elements
}

func (s *SpaceStructure) NameFields() []string {
	return s.nameFields
}

func (s *SpaceStructure) AppendElement(element *Element) {
	s.elements = append(s.elements, element)
	s.numberOfRows++;
}

// Field Calculation
func (s *SpaceStructure) Fields(nameFields []string) []*Field {
	fieldGroup := make([]*Field, SIZE_DEFAULT)
	for _, name := range nameFields {
		fieldGroup = append(fieldGroup, s.field(name))
	}
	return fieldGroup
	
}

func (s *SpaceStructure) ElementField(nameFields []string) []*Element {
	initField := nameFields[:]
	fieldElements := make([]*Element, s.numberOfRows)
	//fieldMapping := s.fieldsDefault[findIndexOfStringArray(s.nameFields, fieldName)]
	for i, element := range s.elements {
		initElement := make([]interface {}, SIZE_DEFAULT)
		for _, label := range nameFields {
			initElement = append(initElement, element.Field(label))
		}
		fieldElements[i] = new(Element).Init().InitLabel(initField)
		fieldElements[i].Set(util.MapCsvJson(initField, initElement))
	}
	return fieldElements
}

func (s *SpaceStructure) Join(field *Field) (ESpace, error) {
	name := field.GetName()
	elements := field.GetData()
	numberMissing := s.numberOfRows - len(elements)
	if numberMissing < 0 {
		err := errors.New("Miss match length data join with number of row of SpaceStructure")
		return s, err
	}
	
	if numberMissing > 0 {
		missingRow := make([]interface {}, numberMissing)
		for i := 0; i < numberMissing; i++ {
			missingRow[i] = ""
		}
		elements = append(elements, missingRow)
	}
	s.nameFields = append(s.nameFields, name);
	
	s.fieldIndex[name] = s.numberOfColumns;
	s.numberOfColumns++;
	for i := 0; i < s.numberOfRows; i++ {
		s.elements[i] = s.elements[i].InitLabel(s.nameFields)
		s.elements[i].SetField(name, elements[i])
	}
	return s, nil
}

func (s *SpaceStructure) Group(nameFields []string) (ESpace, error) {
	elementsByField := s.ElementField(nameFields)
	var fieldGroup ESpace = new(SpaceStructure).Init().InitFields(nameFields)
	fieldGroup.InitElements(elementsByField)
	return fieldGroup, nil
}

func (s *SpaceStructure) Drop(nameFields []string) (ESpace, error) {
	for _, field := range nameFields {
		tempNameFields := make([]string, 0)
		for i := 0; i < s.numberOfRows; i++ {
			if position := util.FindPositionArray(field, s.nameFields); position != -1 {
				s.elements[i].Drop(field)
				for i := 0; i < s.numberOfColumns; i++ {
					if i != position {
						tempNameFields = append(tempNameFields, s.nameFields[i])
					}
				}
				
			}
		}
		s.nameFields = tempNameFields
		s.numberOfColumns--
	}
	
	return s, nil
}

func (s *SpaceStructure) Reshape() (ESpace, error) {
	for i, element := range s.elements {
		sum := element.Sum()
		temp := make([]interface {}, s.numberOfColumns)
		for j, label := range element.Label() {
			temp[j] = s.elements[i].Field(label).(float64) / sum
		}
		tempElement := new(Element).Init().InitLabel(s.nameFields)
		tempElement.Set(util.MapCsvJson(s.nameFields, temp))
		s.elements[i] = tempElement
	}
	return s, nil
}

func (s *SpaceStructure) FileDocs() []string {
	return s.StringElements
}

// Statistic
func (s *SpaceStructure) FieldTypes() []string {
	strFieldTypes := make([]string, SIZE_DEFAULT)
	inmemField := s.Fields(s.NameFields())
	for _, field := range inmemField {
		strFieldTypes = append(strFieldTypes, field.GetType())
	}
	return strFieldTypes
}

// Search Engine
func (s *SpaceStructure) SearchElement(key string, value interface {}) *Element {
	searchField := s.field(key)
	position := searchField.Find(value)
	return s.elements[position]
}

func (s *SpaceStructure) Search(key string, value interface {}) []interface {} {
	return nil
}

func (s *SpaceStructure) Save() {

}

func (s *SpaceStructure) JsonSpace(path string) {
	
}