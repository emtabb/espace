# ESpace - Data Management 

### I. Installation.
- Require: go runtime environment 1.14 or higher.

```terminal
go get github.com/emtabb/espace.
```

### II. Document the method

##### 1. Generate Methods
- Init is the method for data struct implement ESpace generate all variables for its properties.
- InitElements() provice feature set the arrays with standard map type of this framework. Element type 
    is imported from /src/element
- InitFields() also calls Init NameField method, provide to set name fields properties of class, namefield is
the name of one column or one key work in element map.
- InitSpace() Init all properties of other space.
- Save() save all information of properties to data storage. In this version itnot working.

##### 2. Processing data format methods.
- LoadSpace(string) is mechanism of processing to handle data not clean. No format, after processing data have map structure and storaged in mem by Elements 
- CsvSpace(string) is read all data have Csv structure from the file path. All data is converted to Elements Type.
- AppendElement(*Element) add new one element to the space.
- FileDocs() return arrays contain each element to string.

##### 3. Handle data with element type.
- Head(): Return 5 Element in Space.
- Element(): Return Element Type .
- ElementField([]string): return all new Elements from Space with []string is sub namefields.
- SeachElement(string, interface {}): Return Element with key is string and data is interface {}.
- SetElement(int, *Element): Set element for integer index with new Element.
- SetElementKeyValue(int, string, interface {}): Set key of integer index element with new value is interface {}.

##### 4. Handle data with field type.
- Field(): Return Field Type
- Fields([]string): Return array of fields with namefields init()
- FieldTypes(): Return all type of data in space for each key of element.

##### 5. Get properties methods.
- Search(string, interface{}): search data 
- Float(string): return array of data with types is float64 for one namefields.
- Shape(): return size of Space: row - column.
- NameFields(): return array data name fields.

##### 6. Contruct dynamic space.
- Join(*Field): Add new field for space, return itself Space.
- Group([]string): Group subfields in all fields in the Space by the nameFields init, return the new space.
- Drop([]string): Delete subfields in all fields in the Space by the nameFields init, return itself Space.
- Reshape(): Return new space with data in each element reshape from (-1, 1).

### III. Introduction.
- Import on program:
```golang
import "github.com/emtabb/espace"
```
- Generate the struct implementation with this example:
```golang
import "github.com/emtabb/espace"
import "github.com/emtabb/espace/surface"

func main() {
    var someSpace espace.ESpace = surface.Surface()
    // or
    someSpace := surface.Surface()
}
```
