package space

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	. "github.com/emtabb/espace"
	. "github.com/emtabb/espace/api/constant"
	. "github.com/emtabb/espace/api/element"
	"github.com/emtabb/espace/api/space/util"
	. "github.com/emtabb/field"
	. "github.com/emtabb/field/src/statistic"
	"github.com/emtabb/qugo"
	"github.com/emtabb/qugo/qu/impl"
	. "github.com/emtabb/state"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type SpaceStructure struct {
	ESpace
	fieldIndex map[string] int
	numberOfColumns int
	numberOfRows int
	nameFields []string
	isGenerate bool

	// MongoDB Info
	mongodb *mongo.Database
	collection string
	// End Mongodb Info

	inmemField []*Field
	inmemFieldName []string
	StringStates []string
	states []State
	NameSpace []string
}

//private Method (s *SpaceStructure)
//=============================================================
func (s *SpaceStructure) generate() {
	s.isGenerate = true
	s.nameFields = make([] string, SIZE_DEFAULT)
	s.fieldIndex = make(map[string] int)
	s.inmemField = make([] *Field, SIZE_DEFAULT)
	s.inmemFieldName = make([] string, SIZE_DEFAULT)
	s.StringStates = make([] string, SIZE_DEFAULT)
	s.states = make([]State, SIZE_DEFAULT)
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

func (s *SpaceStructure) InitStates(states []State) ESpace {
	if !s.isGenerate {
		s.generate()
	}
	for _, state := range states {
		s.AppendState(state)
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
	s.InitFields(space.NameFields()).InitStates(space.States())
	return s
}

func (s *SpaceStructure) State() State {
	var e State
	return e
}

func (s *SpaceStructure) States(caps ...int) []State {
	limit := int32(500)
	if len(caps) > 0 {
		limit = int32(caps[0])
	}
	return qugo.Operator().InitStates(new(List).ByStates(s.states)).Map(func(state State) State {
		return state.(*Element).GetProperty()
	}).Limit(limit).Collect().ToArray()
}

func (s *SpaceStructure) Elements() []State {
	return s.states
}

func contains(arrays []string, str string) bool {
	for _, v := range arrays {
		if v == str {
			return true
		}
	}

	return false
}

func (s *SpaceStructure) LoadModel(bindModels interface{}) States {

	bindModelsVal := reflect.ValueOf(bindModels)
	if bindModelsVal.Kind() != reflect.Ptr {
		fmt.Errorf("results argument must be a pointer to a slice, but was a %s", bindModelsVal.Kind())
	}

	sliceVal := bindModelsVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		fmt.Errorf("results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind())
	}

	elementType := sliceVal.Type().Elem()
	var index int
	listState := new(List).Generate()
	for _, _state := range s.states {
		mapper := _state.(*Element).GetProperty()
		if sliceVal.Len() == index {
			// slice is full
			newElem := reflect.New(elementType)
			sliceVal = reflect.Append(sliceVal, newElem.Elem())
			sliceVal = sliceVal.Slice(0, sliceVal.Cap())
		}

		newElem := reflect.New(sliceVal.Index(index).Addr().Type()).Elem()
		newElem.Set(sliceVal.Index(index).Addr())

		modelType := reflect.TypeOf(sliceVal.Index(index).Addr().Interface())
		log.Println(modelType.String())
		for i := 0; i < modelType.Elem().NumField(); i++ {
			f := newElem.Elem().Type().Field(i).Tag.Get("json")
			if mapper[f] != nil {
				switch newElem.Elem().Type().Field(i).Type.Kind() {
				case reflect.String :
					newElem.Elem().Field(i).SetString(mapper[f].(string))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64 :
					newElem.Elem().Field(i).SetInt(mapper[f].(int64))
				case reflect.Float32, reflect.Float64:
					newElem.Elem().Field(i).SetFloat(mapper[f].(float64))
				case reflect.Interface:
					newElem.Elem().Field(i).Set(reflect.ValueOf(mapper[f]))
				case reflect.Array:
					if objectId, ok := mapper[f].(primitive.ObjectID); ok {
						newElem.Elem().Field(i).Set(reflect.ValueOf(objectId))
					} else {
						objectIDS, _:= primitive.ObjectIDFromHex(reflect.ValueOf(objectId).String())
						newElem.Elem().Field(i).Set(reflect.ValueOf(objectIDS))
					}

				default:
					newElem.Elem().Field(i).Set(reflect.ValueOf(""))
				}
			}
		}

		sliceVal.Index(index).Set(newElem.Elem())
		listState.Add(sliceVal.Index(index))
		index++
	}
	return listState
}

func (s *SpaceStructure) LoadSpace(collection string) ESpace {

	var mapJsonData []map[string] interface{}
	s.collection = collection
	if s.mongodb != nil {
		cursor, err := (*s.mongodb).Collection(collection).Find(context.TODO(), bson.M {})
		count, err := (*s.mongodb).Collection(collection).CountDocuments(context.TODO(), bson.M {})
		if err != nil {
			panic(err)
		}

		if err = cursor.All(context.TODO(), &mapJsonData); err != nil {
			panic(err)
		}
		s.numberOfRows = int(count)
		s.numberOfColumns = len(mapJsonData[0])

		if !s.isGenerate {
			s.generate()
		}
		s.states = make([]State, s.numberOfRows)
		labels := make([]string, 0)
		for key:= range mapJsonData[0] {
			labels = append(labels, key)
		}
		s.nameFields = labels

		for i := 0; i < s.numberOfRows; i++ {
			s.states[i] = new(Element).Init().Label(labels).Property(mapJsonData[i])
		}
	}
	return s
}

func (s *SpaceStructure) MongoSpace(ConnectionString string, DbName string) ESpace {
	var ctx = context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(ConnectionString))
	if err != nil {
		panic(err)
	}
	s.mongodb = client.Database(DbName)
	return s
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
		collectSlice := regexCollectNumber.FindAllString(s.StringStates[i], -1)
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
		s.StringStates[i] = strings.Join(collectSlice, ",")
	}
}

func (s *SpaceStructure) splitData() {

}

func (s *SpaceStructure) setupElementByCsv(rows int) {
	StringFrame := make([][] string, rows)
	for i := 0; i < rows; i++ {
		StringFrame[i] = make([]string, s.numberOfColumns)
		data := csv.NewReader(strings.NewReader(s.StringStates[i]))
		data.Comma = ','
		StringFrame[i], _ = data.Read()
	}

	field := csv.NewReader(strings.NewReader(s.StringStates[0]))
	field.Comma = ','
	s.nameFields, _ = field.Read()
	s.numberOfColumns = len(s.nameFields)
	for i := 0; i < s.numberOfColumns; i++ {
		s.fieldIndex[s.nameFields[i]] = i
	}

	for i := 0; i < len(StringFrame) - 1; i++ {
		dataTemp := util.StringArrayToInterface(StringFrame[i + 1])
		state := new(Element).Init().Label(s.nameFields).Property(util.MapCsvJson(s.nameFields, dataTemp))
		s.AppendState(state)
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
	isScanNormal := true
	if condition != "" {
		isScanNormal = false
	}
	
	if isScanNormal {
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
		s.StringStates = append(s.StringStates, scanner.Text())
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
				s.StringStates = append(s.StringStates, scanner.Text())
				rows++
				continue
			}
		}
		if beginTable == true {
			if strings.Contains(scanner.Text(), condition) {
				beginTable = false
				continue
			} 
			s.StringStates = append(s.StringStates, scanner.Text())
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
		for i, state := range s.Elements() {
			fieldData[i] = state.(*Element).GetField(fieldName)
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

func (s *SpaceStructure) Head() []State {
	states := make([]State, SIZE_DEFAULT)
	hello := s.States()
	for i := 0; i < 5; i++ {
		states = append(states, hello[i])
	}
	return states
}

func (s *SpaceStructure) NameFields() []string {
	return s.nameFields
}

func (s *SpaceStructure) AppendState(state State) {
	s.states = append(s.states, state)
	s.numberOfRows++
}

// Field Calculation
func (s *SpaceStructure) Fields(nameFields []string) []*Field {
	fieldGroup := make([]*Field, SIZE_DEFAULT)
	for _, name := range nameFields {
		fieldGroup = append(fieldGroup, s.field(name))
	}
	return fieldGroup
	
}

func (s *SpaceStructure) FieldsOfState(nameFields []string) []State {
	initField := nameFields[:]
	fieldStates := make([]State, s.numberOfRows)
	//fieldMapping := s.fieldsDefault[findIndexOfStringArray(s.nameFields, fieldName)]
	for i, state := range s.states {
		initElement := make([]interface {}, SIZE_DEFAULT)
		for _, label := range nameFields {
			initElement = append(initElement, state.(*Element).GetField(label))
		}
		fieldStates[i] = new(Element).Init().Label(initField).Property(util.MapCsvJson(initField, initElement))
	}
	return fieldStates
}

func (s *SpaceStructure) Join(field *Field) (ESpace, error) {
	name := field.GetName()
	states := field.GetData()
	numberMissing := s.numberOfRows - len(states)
	if numberMissing < 0 {
		err := errors.New("Miss match length data join with number of row of SpaceStructure")
		return s, err
	}
	
	if numberMissing > 0 {
		missingRow := make([]interface {}, numberMissing)
		for i := 0; i < numberMissing; i++ {
			missingRow[i] = ""
		}
		states = append(states, missingRow)
	}
	s.nameFields = append(s.nameFields, name)
	
	s.fieldIndex[name] = s.numberOfColumns
	s.numberOfColumns++
	for i := 0; i < s.numberOfRows; i++ {
		s.states[i] = s.states[i].(*Element).Label(s.nameFields).Field(name, states[i])
	}
	return s, nil
}

func (s *SpaceStructure) Group(nameFields []string) (ESpace, error) {
	statesByField := s.FieldsOfState(nameFields)
	var fieldGroup ESpace = new(SpaceStructure).Init().InitFields(nameFields)
	fieldGroup.InitStates(statesByField)
	return fieldGroup, nil
}

func (s *SpaceStructure) Drop(nameFields []string) (ESpace, error) {
	for _, field := range nameFields {
		tempNameFields := make([]string, 0)
		for i := 0; i < s.numberOfRows; i++ {
			if position := util.FindPositionArray(field, s.nameFields); position != -1 {
				s.states[i].(*Element).Drop(field)
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
	for i, state := range s.states {
		sum := state.(*Element).Sum()
		temp := make([]interface {}, s.numberOfColumns)
		for j, label := range state.(*Element).GetLabel() {
			temp[j] = s.states[i].(*Element).GetField(label).(float64) / sum
		}
		tempElement := new(Element).
			Init().
			Label(s.nameFields).
			Property(util.MapCsvJson(s.nameFields, temp))
		s.states[i] = tempElement
	}
	return s, nil
}

func (s *SpaceStructure) FileDocs() []string {
	return s.StringStates
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
func (s *SpaceStructure) SearchState(key string, value interface {}) State {
	searchField := s.field(key)
	position := searchField.Find(value)
	return s.states[position]
}

func (s *SpaceStructure) Search(key string, value interface {}) []interface {} {
	return qugo.Operator().InitStates(new(List).ByStates(s.states)).Filter(func (state State) bool {
		if state.(*Element).GetField(key) == value {
			return true
		}
		return false
	}).CollectInterface()
}

func (s *SpaceStructure) SetState(index int32, state State) {
	s.states[index] = state
}

func (s *SpaceStructure) SetStateKeyValue(index int32, key string, value interface{}) {
	s.states[index] = s.states[index].(*Element).Field(key, value)
}

func (s *SpaceStructure) Save(collection ...string) error {
	saveCollection := s.collection
	if len(collection) > 0 {
		saveCollection = collection[0]
	}
	var err error = nil
	if s.mongodb != nil {
		err = s.smartStoreDatabase(saveCollection)
	} else {
		err = s.smartStoreFile(saveCollection)
	}
	return err
}

func (s *SpaceStructure) smartStoreDatabase(saveCollection string) error {
	var err error = nil
	var ListSave = new(List)
	var ListUpdate = new(List)

	allAvailableStates := s.States()
	for _, availableState := range allAvailableStates {
		if availableState.(map[string] interface{})["_id"] != nil || availableState.(map[string] interface{})["id"] != nil {
			ListUpdate.Add(availableState)
		} else {
			ListSave.Add(availableState)
		}
	}

	_, err = s.mongodb.
		Collection(saveCollection).
		InsertMany(context.TODO(), qugo.Operator().
			InitStates(ListSave).CollectInterface())

	for _, replaceState := range ListUpdate.ToArray() {
		objectId, _ := replaceState.(map[string] interface{})
		_, err = s.mongodb.
			Collection(saveCollection).
			ReplaceOne(context.TODO(), bson.D {{"_id" , objectId["_id"]}}, replaceState.(map[string] interface{}))
	}
	return err
}

func (s *SpaceStructure) smartStoreFile(path string) error {
	//if strings.HasSuffix(path, ".jpge") {
	//	image.storeImage(nil, path)
	//}
	return s.storeFile(path)
}

func (s *SpaceStructure) storeFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		file, err = os.Create(path)
	}

	mapToArray := qugo.Operator().InitStates(new(List).ByStates(s.States(qugo.UN_LIMITED))).Map(func(_state State) State {
		arrayEle := new(List)
		for _, value := range s.NameFields() {
			arrayEle.Add(_state.(map[string] interface{})[value])
		}
		return arrayEle
	}).Limit(impl.UN_LIMITED).Collect().ToArray()

	fileData := ""
	simpleString := make([]string, 0)
	for i := range mapToArray {
		subArray := qugo.Operator().InitStates(mapToArray[i].(States)).Limit(impl.UN_LIMITED).CollectInterface()
		listString := String(subArray)
		simpleString = append(simpleString, strings.Join(listString, ","))

	}
	fileData = strings.Join(simpleString, "\n")
	checkError("Cannot create file", err)
	defer file.Close()

	bytesWrite, err := file.WriteString(fileData)
	fmt.Printf("wrote %d bytesWrite\n", bytesWrite)

	file.Sync()
	return err
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func (s *SpaceStructure) JsonSpace(path string) {
	
}

func (s *SpaceStructure) GetMongoDb() *mongo.Database{
	return s.mongodb
}