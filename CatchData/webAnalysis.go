package main

import (
	"github.com/lxn/walk"
	"sort"
	"strings"
)

type HttpBody struct{
	Deploy struct{
		ListMeta struct{
			Total int `json:"totalItems"`
		}`json:"listMeta"`
		Status struct{
			Run int `json:"running"`
		}`json:"status"`
		Data []struct{
			Meta struct{
				Name string `json:"name"`
				ANNO struct{
					KUBE string `json:"kubectl.kubernetes.io/last-applied-configuration,omitempty"`
				}`json:"annotations,omitempty"`
			}`json:"objectMeta,omitempty"`
		}`json:"deployments,omitempty"`
	}`json:"deploymentList"`
	
}

type KUBE struct{
	SPEC1 struct{
		Template struct{
			SPEC2 struct{
				Container []struct{
					Env []struct{
						Name string `json:"name,omitempty"`
						Value string `json:"value,omitempty"`
					}`json:"env,omitempty"`
				}`json:"containers,omitempty"`
			}`json:"spec,omitempty"`
		}`json:"template,omitempty"`
	} `json:"spec,omitempty"`
}

type Foo struct {
	Index   int
	Name     string
	Port     int
	checked bool
}

type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo
}
func (m *FooModel) ResetRows() {
	m.items = make([]*Foo, 0)
}

// Called by the TableView to retrieve if a given row is checked.
func (m *FooModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FooModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

func NewFooModel() *FooModel {
	m := new(FooModel)
	m.ResetRows()
	return m
}

func (m *FooModel) RowCount() int {
	return len(m.items)
}

func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Name

	case 2:
		return item.Port
	
	}

	panic("unexpected col")
}

func (m *FooModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.Name < b.Name)

		case 2:
			return c(a.Port < b.Port)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}


func selectCall(m *FooModel){
	for i := range m.items {
		if strings.Contains(m.items[i].Name, "usp") ||
		strings.Contains(m.items[i].Name, "pc")||
		strings.Contains(m.items[i].Name, "gc")||
		strings.Contains(m.items[i].Name, "sipgw"){
			m.items[i].checked = true
		}else{
			m.items[i].checked = false
		}
	}
	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)
}


func selectAll(m *FooModel){
	for i := range m.items {
		m.items[i].checked = true
	}
	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)
}

func unSelectAll(m *FooModel){
	for i := range m.items {
		m.items[i].checked = false
	}
	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)
}

