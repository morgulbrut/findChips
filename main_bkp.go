package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/morgulbrut/findChips/part"
	"github.com/morgulbrut/helferlein"
	"github.com/morgulbrut/soup"
)

func main() {
	prt := os.Args[1]
	ps := getPartInfo(prt)
	for _, p := range ps {
		fmt.Println("\n")
		part.WriteCSV(p)
	}
	//part.Print(p)
	//
}

func getPartInfo(prt string) []part.Part {
	var p part.Part
	var fp []part.Part
	var pns []string

	dists := []string{"Mouser", "Farnell", "Digi-Key", "Avnet", "TTI", "RS"}
	doc := getDistPage(prt)
	p.Distributors = getDistributors(doc, dists)

	// get candidates for parts
	for _, d := range p.Distributors {
		for _, pn := range d.Partnumbers {
			if !helferlein.StringInSlice(pn.ManPartnumber, pns) {
				pns = append(pns, pn.ManPartnumber)
			}
		}
	}

	for _, pn := range pns {
		doc := getDistPage(pn)
		if doc.Error == nil {
			p.Partnumber = pn
			p.Distributors = getDistributors(doc, dists)

			docDet := parseDetails(doc)
			if docDet.Error == nil {
				p.Parameters = getParameters(docDet)
				p.Datasheets = getDatasheets(docDet)
				p.Alternatives = getAlternatives(docDet)
			}
		}

		fp = append(fp, p)
	}

	return fp
}

func getDistPage(prt string) soup.Root {
	url := fmt.Sprintf("https://www.findchips.com/lite/%s", prt)
	resp, _ := soup.Get(url)
	return soup.HTMLParse(resp)
}

func parseDetails(doc soup.Root) soup.Root {
	det := doc.FindAll("a", "class", "sub-header-item")
	var ret soup.Root
	for _, d := range det {
		if strings.Contains(d.Attrs()["href"], "detail") {
			url := fmt.Sprintf("https://www.findchips.com%s", d.Attrs()["href"])
			resp, _ := soup.Get(url)
			ret = soup.HTMLParse(resp)
		}
	}
	return ret
}

func getDistributors(doc soup.Root, dists []string) []part.Distributor {
	var ret []part.Distributor
	var dist part.Distributor

	res := doc.FindAll("div", "class", "distributor-results")
	for _, r := range res {
		d := r.Find("h3", "class", "distributor-title")
		if d.Error == nil {
			ds := d.FullText()
			ds = strings.TrimSpace(ds)
			ds = strings.Split(ds, "\n")[0]
			if helferlein.Contains(ds, dists) {
				dist.Name = ds
				distPartNos := r.FindAll("tr", "class", "row")
				for _, dpn := range distPartNos {
					var pn part.Partnumber
					pn.DistPartnumber = dpn.Attrs()["data-distino"]
					pn.ManPartnumber = dpn.Attrs()["data-mfrpartnumber"]
					pn.Desc = dpn.Children()[5].Find("span", "class", "td-description").Text()
					dist.Partnumbers = append(dist.Partnumbers, pn)
				}
				ret = append(ret, dist)
			}
		}
	}
	return ret
}

func getParameters(doc soup.Root) []part.Parameter {
	var ret []part.Parameter
	det := doc.Find("ul", "class", "part-details-list")
	if det.Error == nil {
		for _, d := range det.FindAll("li") {
			var par part.Parameter
			par.Param = strings.Trim(d.Find("small").Text(), ":")
			par.Val = strings.TrimSpace(d.Find("p").Text())
			ret = append(ret, par)
		}
	}
	return ret
}

func getDatasheets(doc soup.Root) []string {
	var ret []string
	ds := doc.Find("div", "class", "datasheet-item")
	if ds.Error == nil {
		das := ds.Find("a").Attrs()["href"]
		ret = append(ret, das)
	}
	return ret
}

func getAlternatives(doc soup.Root) []string {
	var ret []string
	alts := doc.Find("table", "class", "part-suggestions")
	if alts.Error == nil {
		alt := alts.Find("tbody").FindAll("tr")
		for _, a := range alt {
			link := a.Find("td", "class", "td-col-1")
			if link.Error == nil {
				linkT := link.Find("a").Text()
				ret = append(ret, strings.TrimSpace(linkT))
			}

		}
	}
	return ret
}
