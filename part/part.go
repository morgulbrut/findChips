package part

import (
	"github.com/morgulbrut/color"
)

type Parameter struct {
	Param string
	Val   string
}

type Partnumber struct {
	ManPartnumber  string
	DistPartnumber string
	Desc           string
}

type Distributor struct {
	Name        string
	Partnumbers []Partnumber
}

type Part struct {
	Partnumber   string
	Parameters   []Parameter
	Datasheets   []string
	Alternatives []string
	Distributors []Distributor
}

func PrintPart(p Part) {
	color.Red(p.Partnumber)
	color.Yellow("Distributors")
	for _, d := range p.Distributors {
		color.Green(d.Name)
		for _, pn := range d.Partnumbers {
			color.Cyan("%s : %s : %s\n", pn.ManPartnumber, pn.DistPartnumber, pn.Desc)
		}
	}
	color.Yellow("Details")
	for _, p := range p.Parameters {
		color.Cyan("%s : %s", p.Param, p.Val)
	}
	color.Yellow("Alternatives")
	for _, a := range p.Alternatives {
		color.Cyan(a)
	}
	color.Yellow("Datasheets")
	for _, a := range p.Datasheets {
		color.Cyan(a)
	}
}
