package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/morgulbrut/findChips/part"
	"github.com/morgulbrut/soup"
)

func main() {
	prt := os.Args[1]
	p := getPartInfo(prt)
	part.PrintPart(p)
}

func contains(input string, words []string) bool {
	for _, word := range words {
		if strings.Index(input, word) > -1 {
			return true
		}
	}
	return false
}

func getPartInfo(prt string) part.Part {

	var p part.Part
	p.Partnumber = prt

	url := fmt.Sprintf("https://www.findchips.com/lite/%s", prt)
	resp, _ := soup.Get(url)
	doc := soup.HTMLParse(resp)

	dists := []string{"Mouser", "Farnel", "Digi-Key", "Avnet", "TTI", "RS"}

	res := doc.FindAll("div", "class", "distributor-results")
	for _, r := range res {
		distributor := r.Find("h3", "class", "distributor-title").FullText()
		distributor = strings.TrimSpace(distributor)
		distributor = strings.Split(distributor, "\n")[0]
		if contains(distributor, dists) {
			var dist part.Distributor
			dist.Name = distributor
			distPartNos := r.FindAll("tr", "class", "row")
			for _, dpn := range distPartNos {
				var pn part.Partnumber
				pn.DistPartnumber = dpn.Attrs()["data-distino"]
				pn.ManPartnumber = dpn.Attrs()["data-mfrpartnumber"]
				pn.Desc = dpn.Children()[5].Find("span", "class", "td-description").Text()
				dist.Partnumbers = append(dist.Partnumbers, pn)
			}
			p.Distributors = append(p.Distributors, dist)
		}
	}

	det := doc.FindAll("a", "class", "sub-header-item")
	for _, d := range det {
		if strings.Contains(d.Attrs()["href"], "detail") {
			url = fmt.Sprintf("https://www.findchips.com%s", d.Attrs()["href"])
			resp, _ = soup.Get(url)
			doc = soup.HTMLParse(resp)
			det := doc.Find("ul", "class", "part-details-list")
			for _, d := range det.FindAll("li") {
				var par part.Parameter
				par.Param = strings.Trim(d.Find("small").Text(), ":")
				par.Val = strings.TrimSpace(d.Find("p").Text())
				p.Parameters = append(p.Parameters, par)
			}
			ds := doc.Find("div", "class", "datasheet-item").Find("a").Attrs()["href"]
			p.Datasheets = append(p.Datasheets, ds)
			alts := doc.Find("table", "class", "part-suggestions").Find("tbody").FindAll("tr")
			for _, a := range alts {
				link := a.Find("td", "class", "td-col-1").Find("a")
				p.Alternatives = append(p.Alternatives, strings.TrimSpace(link.Text()))
			}
		}
	}
	return p
}
