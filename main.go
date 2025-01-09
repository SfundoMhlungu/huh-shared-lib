package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/gofrs/uuid"

	"github.com/charmbracelet/huh/spinner"
)

type Spice int

const (
	Mild Spice = iota + 1
	Medium
	Hot
)

func (s Spice) String() string {
	switch s {
	case Mild:
		return "Mild "
	case Medium:
		return "Medium-Spicy "
	case Hot:
		return "Spicy-Hot "
	default:
		return ""
	}
}

type Order struct {
	Burger       Burger
	Side         string
	Name         string
	Instructions string
	Discount     bool
}

type Burger struct {
	Type     string
	Toppings []string
	Spice    Spice
}

var f *Fields = &Fields{}
var theme *huh.Theme = huh.ThemeBase()

func createMaps() {
	f.M = make(map[string]*Value)
	f.Elems = make(map[string]Fieldsinterface)
	f.Group = make(map[string][]huh.Field)
}

//export SetTheme
func SetTheme(t *C.char) {

	switch C.GoString(t) {
	case "dracula":
		theme = huh.ThemeDracula()
	case "Charm":
		theme = huh.ThemeCharm()
	case "Catppuccin":
		theme = huh.ThemeCatppuccin()
	case "Base16":
		theme = huh.ThemeBase16()
	default:
		theme = huh.ThemeBase()
	}

}

//export CreateInput
func CreateInput(t C.int) *C.char {
	uniqueid, _ := uuid.NewV4()

	if f.Elems == nil {
		createMaps()
	}

	i := &NewInput{}
	i.elem = nil // remove don't need
	i.id = uniqueid.String()
	i.opts = NewInputOpts{}
	i.validators = []string{}
	i.t = int(t)
	i.value = new(string) // remove
	f.Elems[uniqueid.String()] = i
	// // fmt.Println(i.opts)
	// ptr := C.malloc(C.size_t(unsafe.Sizeof(NewInput{})))
	// fmt.Println(ptr)
	// // fmt.Println(*(*NewInputOpts)(opts))
	// if ptr == nil {
	// 	return nil
	// }
	// // // Copy the struct into the allocated memory
	// *(*NewInput)(ptr) = *i
	// fmt.Printf("id: %s\n", i.id)
	// // fmt.Printf("elem:\n%v\n", i.elem)
	// fmt.Printf("opts:\n")
	// fmt.Println(i.opts)
	// return ptr
	return C.CString(uniqueid.String())
}

//export SetInputOptions
func SetInputOptions(id *C.char, title *C.char, desc *C.char, placeholder *C.char, validators *C.char) C.int {

	if id == nil {
		// fmt.Println("Invalid struct pointer")
		return 2
	}
	s, exists := f.Elems[C.GoString(id)]
	if !exists {

		return 1
	}

	// s := f.Elems[C.GoString(id)]
	goStruct := s.(*NewInput)
	goStruct.opts.Description = C.GoString(desc)
	goStruct.opts.Title = C.GoString(title)
	goStruct.opts.Placeholder = C.GoString(placeholder)
	if validators != nil {
		v := strings.Split(C.GoString(validators), ",")
		// fmt.Println(v)
		goStruct.validators = v

	}
	s = nil

	return 0
}

//export RunInput
func RunInput(id *C.char) *C.char {
	s, exists := f.Elems[C.GoString(id)]
	if !exists {

		return C.CString("element does not exists")
	}

	goStruct := s.(*NewInput)
	return C.CString(goStruct.Run())
}

//export Confirm
func Confirm(title *C.char, affirmative *C.char, negative *C.char) *C.char {
	if f.Elems == nil {
		createMaps()
	}
	id, _ := uuid.NewV4()

	conf := &NewConfirm{title: C.GoString(title), affirmative: C.GoString(affirmative), negative: C.GoString(negative)}
	conf.value = false
	f.Elems[id.String()] = conf
	// return conf.Run()
	return C.CString(id.String())
}

//export Select
func Select(title *C.char, opts *C.char) *C.char {
	if f.Elems == nil {
		createMaps()
	}
	id, _ := uuid.NewV4()
	s := strings.Split(C.GoString(opts), ",")

	element := &NewSelect{title: C.GoString(title), options: s}
	element.value = new(string)
	f.Elems[id.String()] = element
	// return C.CString(element.Run())
	return C.CString(id.String())
}

//export MultiSelect
func MultiSelect(id *C.char, title *C.char, opts *C.char) C.int {
	if f.Elems == nil {
		createMaps()
	}
	s := strings.Split(C.GoString(opts), ",")

	element := &NewMultiSelect{title: C.GoString(title), options: s}
	element.value = []string{}
	// huh.NewForm(huh.NewGroup(element.ForGroup()))
	f.Elems[C.GoString(id)] = element
	// return C.CString(strings.Join(element.Run(), ","))
	return 0
}

//export CreateGroup
func CreateGroup(ids *C.char) *C.char {
	if f.Elems == nil {
		createMaps()
	}
	allids := strings.Split(C.GoString(ids), ",")
	// fmt.Println(allids, "all ids")
	// fmt.Printf("elements in cache:  %#v\n", f.Elems)
	var elems []huh.Field

	for _, eleId := range allids {
		s, exists := f.Elems[eleId]

		if exists {
			elems = append(elems, s.ForGroup())

		}
	}
	// fmt.Printf("elements before group:%#v\n", elems) // ✅ works on free the struct is not added

	groudid, _ := uuid.NewV4()

	f.Group[groudid.String()] = elems

	// fmt.Printf("elements in group:  %#v\n", f.Group)
	return C.CString(groudid.String())
}

//export CreateForm
func CreateForm(groudids *C.char) {
	if f.Elems == nil {
		createMaps()
	}
	s := strings.Split(C.GoString(groudids), ",")

	g := make([]*huh.Group, 0, len(s))
	for _, id := range s {
		if el, exists := f.Group[id]; exists {
			// fmt.Println(el, id)
			g = append(g, huh.NewGroup(el...))
		}
	}

	huh.NewForm(g...).Run()
}

//export GetValue
func GetValue(id *C.char) *C.char {
	if el, exists := f.Elems[C.GoString(id)]; exists {
		return C.CString(el.GetValue())
	}

	return C.CString("")
}

//export Run
func Run(id *C.char) *C.char {
	elem, exists := f.Elems[C.GoString(id)]
	if !exists {
		return nil

	}

	return C.CString(elem.Run())

	// switch C.GoString(t) {
	// case "select":
	// 	return C.CString(elem.Run())
	// default:
	// 	return C.CString(elem.Run())
	// }

}

//export Spinner
func Spinner(seconds C.int, title *C.char) {

	act := func() {
		time.Sleep(time.Duration(seconds) * time.Second)
	}

	_ = spinner.New().Title(C.GoString(title)).Accessible(true).Action(act).Run()
}

//export Note
func Note(title *C.char, desc *C.char, label *C.char, next C.int) *C.char {
	if f.Elems == nil {
		createMaps()
	}
	nt := &NewNote{}
	nt.desc = C.GoString(desc)
	nt.title = C.GoString(title)
	// n := huh.NewNote()
	// n.Title(C.GoString(title))
	// n.Description(C.GoString(desc))

	if label != nil && C.GoString(label) != "" {
		s := C.GoString(label)
		nt.label = &s

	}

	id, _ := uuid.NewV4()

	f.Elems[id.String()] = nt

	return C.CString(id.String())

}

//export FreeStruct
func FreeStruct(id *C.char) C.int {
	_, exists := f.Elems[C.GoString(id)]
	if exists {
		// fmt.Println("Before delete:", f.Elems)
		i := C.GoString(id)
		f.Elems[i] = nil
		delete(f.Elems, i)
		// fmt.Println("After delete:", f.Elems)
		return 0
	} else {
		return 1
	}
}

func main() {
}
