package element

import "errors"
import "log"
import "github.com/emtabb/espace/src/space/util"

type Element struct {
	element map[string] interface{}
	label []string
	dimension int
}

func (ele *Element) Init() *Element {
	const SIZE_DEFAULT = 0
	ele.label = make([]string, SIZE_DEFAULT)
	ele.element = make(map[string] interface {})
	return ele
}

func (ele *Element) InitLabel(labels []string) *Element {
	ele.label = labels[:]
	return ele
}

func (ele *Element) Set(element map[string] interface{}) {
	for label := range element {
		ele.element[label] = element[label]
		ele.dimension++
	}
}

func (ele *Element) SetField(fieldName string, field interface{}) {
	ele.element[fieldName] = field
}

func (ele *Element) SetElement(element map[string] interface {}) {
	if !ele.validElementLabel(element) {
		return
	}
	for _, key := range ele.Label() {
		ele.element[key] = element[key]
	}
}

func (ele *Element) validElementLabel(element map[string] interface {}) bool {
	for _, label := range ele.Label() {
		if _, ok := element[label]; !ok {
			return false
		}
	}
	return true
}

func (ele *Element) isKeyExist(key string) bool {
	log.Println("")
	for _, rootKey := range ele.Label() {
		if key == rootKey {
			return true
		}
	}
	return false
}

func (ele *Element) Get() map[string] interface{} {
	return ele.element
}

func (ele *Element) Label() []string {
	return ele.label
}

func (ele *Element) Field(fieldName string) interface {} {
	return ele.element[fieldName]
}

func (ele *Element) Array() ([]*Element) {
	return []*Element { ele }
}

func (ele *Element) Sum() float64 {
	sum := 0.0
	for label := range ele.element {
		data := ele.Field(label).(float64);
		sum += data
	}
	return sum
}

func (ele *Element) Drop(field string) error {
	if _, ok := ele.element[field]; ok {
		tempLabel := make([]string, 0)
		if position := util.FindPositionArray(field, ele.label); position != -1 {
			delete(ele.element, field)
			for i := 0; i < len(ele.label); i++ {
				if i != position {
					tempLabel = append(tempLabel, ele.label[i])
				}	
			}
			ele.label = tempLabel
			return nil
		}
	}
	return errors.New("Have errors")
}

func (ele *Element) Cache() {
	
}

func (ele *Element) ReadCache() {}

func (ele *Element) ToArray() []interface {} {
	arrays := make([]interface {}, ele.dimension)
	for i, name := range ele.Label() {
		arrays[i] = ele.Field(name)
	}
	return arrays
}

func (ele *Element) ToString() string {
	stringElement := ""
	for _, name := range ele.Label() {
		if str, ok := ele.Field(name).(string); ok {
			stringElement += str + " | "
		} else {
			stringElement += util.DoubleToString(ele.Field(name).(float64)) + " | "
		}
		
	}
	return stringElement
}

func (ele *Element) JsonString() (string, error) {
	return util.JsonStringify(ele.element)
}

func (ele *Element) ByJsonString(strJson string) error {
	err := errors.New("")
	ele.element, err = util.JsonParseInterface(strJson)
	return err
}

func (ele *Element) ByJson(element map[string] interface{}) {
	ele.element = element
}