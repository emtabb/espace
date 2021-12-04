package espace

import (
	. "github.com/emtabb/field"
	. "github.com/emtabb/state"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
*
*
*
*/
type ESpace interface {
	/*
	*	- Init is the method for data struct implement ESpace generate all variables for its properties.
	*	- InitElements() provice feature set the arrays with standard map type of this framework. Element type 
	*     is imported from /src/element
	*	- InitFields() also calls Init NameField method, provide to set name fields properties of class, namefield is
	*	the name of one column or one key work in element map.
	*	- InitSpace() Init all properties of other space.
	*	- Save() save all information of properties to data storage. In this version itnot working.
	*/
	Init() ESpace
	InitStates([]State) ESpace
	InitFields([]string) ESpace
	InitSpace(ESpace) ESpace
	Save(...string) error

	/*
	*	- LoadSpace(string) is mechanism of processing to handle data not clean.
	*	No format, after processing data have map structure and stored in mem by States
	*	- CsvSpace(string) is read all data have Csv structure from the file path. All data is converted to States Type.
	*	- AppendState(State) add new one element to the space.
	*	- FileDocs() return arrays contain each element to string.
	*/
	LoadSpace(string) ESpace
	InitSchema(State) ESpace
	LoadModel(interface{}) States
	MongoSpace(string, string) ESpace
	GetMongoDb() *mongo.Database
	CsvSpace(string)
	AppendState(State)
	FileDocs() []string

	/*
	*	- Head(): Return 5 State in Space.
	*	- State(): Return State Type .
	*	- FieldsOfState([]string): return all new States from Space with []string is sub namefields.
	*	- SearchState(string, interface {}): Return State with key is string and data is interface {}.
	*	- SetState(int, State): Set state for integer index with new State.
	*	- SetStateKeyValue(int, string, interface {}): Set key of integer index element with new value is interface {}.
	*/
	Head() []State //
	State() State //
	States(...int) []State
	Elements() []State
	FieldsOfState([]string) []State
	SearchState(string, interface {}) State
	SetState(int32, State)
	SetStateKeyValue(int32, string, interface {})
	

	/*
	*	- Field(): Return Field Type
	*	- Fields([]string): Return array of fields with namefields init()
	*	- FieldTypes(): Return all type of data in space for each key of element.
	*/
	Field() Field
	Fields([]string) []*Field
	FieldTypes() []string

	/*
	*	- Search(string, interface{}): search data 
	*	- Float(string): return array of data with types is float64 for one namefields.
	*	- Shape(): return size of Space: row - column.
	*	- NameFields(): return array data name fields.
	*/
	Search(string, interface {}) []interface {}
	Float(string) []float64
	Shape() (int, int)
	NameFields() []string

	/*
	*	- Join(*Field): Add new field for space, return itself Space.
	*	- Group([]string): Group subfields in all fields in the Space by the nameFields init, return the new space.
	*	- Drop([]string): Delete subfields in all fields in the Space by the nameFields init, return itself Space.
	*	- Reshape(): Return new space with data in each element reshape from (-1, 1)
	*/
	Join(*Field) (ESpace, error)
	Group([]string) (ESpace, error)
	Drop([]string) (ESpace, error)
	Reshape() (ESpace, error)
}