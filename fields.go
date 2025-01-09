package main

import "C"
import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// remove absolete, value earguly returned on run
type Fieldsinterface interface {
	GetValue() string
	Run() string
	ForGroup() huh.Field
}
type Fields struct {
	M     mapWithValues
	Elems map[string]Fieldsinterface
	Group map[string][]huh.Field
}

type NewInputOpts struct {
	Description string
	Title       string
	Placeholder string
}
type NewInput struct {
	elem       *huh.Input
	value      *string
	id         string
	t          int
	opts       NewInputOpts
	validators []string
}

func (i *NewInput) GetValue() string {
	return *i.value
}

func (i *NewInput) Run() string {

	if i.t == 1 {
		h := huh.NewText().Value(i.value).Description(i.opts.Description).
			Title(i.opts.Title).
			Validate(func(s string) error {
				if i.validators != nil {

					for _, vali := range i.validators {
						f, e := Validators[vali]
						fmt.Println(vali, e)

						if e {
							err := f(s)
							if err != nil {
								return err
							}
						}

					}
				}

				return nil
			}).
			Placeholder(i.opts.Placeholder).WithTheme(theme)
		h.Run()
		return *i.value
	}
	// fmt.Println(i.opts, "opts")
	h := huh.NewInput().Value(i.value).
		Description(i.opts.Description).
		Title(i.opts.Title).Validate(func(s string) error {
		if i.validators != nil {
			for _, vali := range i.validators {
				f, e := Validators[vali]

				if e {
					err := f(s)
					if err != nil {
						return err
					}
				}

			}
		}

		return nil
	})

	h.Run()
	h.Focus()

	return *i.value
}

func (i *NewInput) ForGroup() huh.Field {

	if i.t == 1 {
		h := huh.NewText().Value(i.value).Description(i.opts.Description).
			Title(i.opts.Title).
			Placeholder(i.opts.Placeholder).WithTheme(theme)

		return h
	}
	// fmt.Println(i.opts, "opts")
	h := huh.NewInput().Value(i.value).
		Description(i.opts.Description).
		Title(i.opts.Title)
	return h
}

type NewConfirm struct {
	title       string
	affirmative string
	negative    string
	value       bool
}

// absolete
func (i *NewConfirm) GetValue() string {
	if i.value {
		return "1"
	}

	return "0"
}

func (i *NewConfirm) Run() string {
	var value bool
	// fmt.Println(i.opts, "opts")
	huh.NewConfirm().
		Title(i.title).
		Value(&value).
		Affirmative(i.affirmative).
		Negative(i.negative).WithTheme(theme).Run()

	if value {
		return "1"
	}

	return "0"
}

func (i *NewConfirm) ForGroup() huh.Field {

	h := huh.NewConfirm().
		Title(i.title).
		Value(&i.value).
		Affirmative(i.affirmative).
		Negative(i.negative).WithTheme(theme)

	return h
}

type NewSelect struct {
	options []string
	title   string
	value   *string
}

func (s *NewSelect) GetValue() string {
	return *s.value
}

func (s *NewSelect) Run() string {

	huh.NewSelect[string]().
		Title(s.title).
		Options(
			huh.NewOptions(s.options...)...,
		).
		Value(s.value).WithTheme(theme).Run()

	return *s.value
}

func (s *NewSelect) ForGroup() huh.Field {

	sel := huh.NewSelect[string]().
		Title(s.title).
		Options(
			huh.NewOptions(s.options...)...,
		).
		Value(s.value).WithTheme(theme)

	return sel
}

type NewMultiSelect struct {
	options []string
	title   string
	value   []string
}

func (s *NewMultiSelect) GetValue() string {
	return strings.Join(s.value, ",")
}

func (s *NewMultiSelect) Run() string {

	// fmt.Println("muli select", s.options)
	huh.NewMultiSelect[string]().
		Title(s.title).
		Options(
			huh.NewOptions(s.options...)...,
		).
		Value(&s.value).Limit(2).WithTheme(theme).Run()

	return strings.Join(s.value, ",")
}

func (s *NewMultiSelect) ForGroup() huh.Field {
	var value string
	sel := huh.NewSelect[string]().
		Title(s.title).
		Options(
			huh.NewOptions(s.options...)...,
		).
		Value(&value).WithTheme(theme)

	return sel
}

type NewNote struct {
	title string
	desc  string
	label *string
	next  bool
}

func (n *NewNote) GetValue() string {
	return ""
}

func (nt *NewNote) Run() string {
	n := huh.NewNote()
	n.Title(nt.title)
	n.Description(nt.desc)

	if nt.label != nil {
		n.NextLabel("\n" + *nt.label)
		n.Next(nt.next)
	}
	n.Run()

	return ""
}

func (nt *NewNote) ForGroup() huh.Field {
	n := huh.NewNote()
	n.Title(nt.title)
	n.Description(nt.desc)

	if nt.label != nil {
		n.NextLabel("\n" + *nt.label)
		n.Next(nt.next)
	}

	return n
}
