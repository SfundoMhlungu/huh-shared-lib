package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"os"
	"strconv"
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
func Note(title *C.char, desc *C.char, label *C.char, next C.int) {
	if f.Elems == nil {
		createMaps()
	}
	n := huh.NewNote()
	n.Title(C.GoString(title))
	n.Description(C.GoString(desc))

	if label != nil && C.GoString(label) != "" {
		n.NextLabel("\n" + C.GoString(label))
		n.Next(next != 0)
	}
	n.Run()
}

//export Test
func Test() {
	var burger Burger
	var order = Order{Burger: burger}
	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("Charmburger").
			Description("Welcome to _Charmburger™_.\n\nHow may we take your order?\n\n").
			Next(true).
			NextLabel("Next"),
		),

		// Choose a burger.
		// We'll need to know what topping to add too.
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("Charmburger Classic", "Chickwich", "Fishburger", "Charmpossible™ Burger")...).
				Title("Choose your burger").
				Description("At Charm we truly have a burger for everyone.").
				Validate(func(t string) error {
					if t == "Fishburger" {
						return fmt.Errorf("no fish today, sorry")
					}
					return nil
				}).
				Value(&order.Burger.Type),

			huh.NewMultiSelect[string]().
				Title("Toppings").
				Description("Choose up to 4.").
				Options(
					huh.NewOption("Lettuce", "Lettuce").Selected(true),
					huh.NewOption("Tomatoes", "Tomatoes").Selected(true),
					huh.NewOption("Charm Sauce", "Charm Sauce"),
					huh.NewOption("Jalapeños", "Jalapeños"),
					huh.NewOption("Cheese", "Cheese"),
					huh.NewOption("Vegan Cheese", "Vegan Cheese"),
					huh.NewOption("Nutella", "Nutella"),
				).
				Validate(func(t []string) error {
					if len(t) <= 0 {
						return fmt.Errorf("at least one topping is required")
					}
					return nil
				}).
				Value(&order.Burger.Toppings).
				Filterable(true).
				Limit(4),
		),

		// Prompt for toppings and special instructions.
		// The customer can ask for up to 4 toppings.
		huh.NewGroup(
			huh.NewSelect[Spice]().
				Title("Spice level").
				Options(
					huh.NewOption("Mild", Mild).Selected(true),
					huh.NewOption("Medium", Medium),
					huh.NewOption("Hot", Hot),
				).
				Value(&order.Burger.Spice),

			huh.NewSelect[string]().
				Options(huh.NewOptions("Fries", "Disco Fries", "R&B Fries", "Carrots")...).
				Value(&order.Side).
				Title("Sides").
				Description("You get one free side with this order."),
		),

		// Gather final details for the order.
		huh.NewGroup(
			huh.NewInput().
				Value(&order.Name).
				Title("What's your name?").
				Placeholder("Margaret Thatcher").
				Validate(func(s string) error {
					if s == "Frank" {
						return errors.New("no franks, sorry")
					}
					return nil
				}).
				Description("For when your order is ready."),

			huh.NewText().
				Value(&order.Instructions).
				Placeholder("Just put it in the mailbox please").
				Title("Special Instructions").
				Description("Anything we should know?").
				CharLimit(400).
				Lines(5),

			huh.NewConfirm().
				Title("Would you like 15% off?").
				Value(&order.Discount).
				Affirmative("Yes!").
				Negative("No."),
		),
	).WithAccessible(accessible)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	prepareBurger := func() {
		time.Sleep(2 * time.Second)
	}

	_ = spinner.New().Title("Preparing your burger...").Accessible(accessible).Action(prepareBurger).Run()
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
