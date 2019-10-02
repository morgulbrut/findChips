package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/morgulbrut/color"
	"github.com/morgulbrut/soup"
)

func main() {
	part := os.Args[1]
	getDistributorInfo(part)
	getPartInfo(part)
}

func getPartInfo(part string) {
	url := fmt.Sprintf("https://www.findchips.com/lite/%s", part)
	resp, _ := soup.Get(url)
	doc := soup.HTMLParse(resp)

	det := doc.FindAll("a", "class", "sub-header-item")
	color.Yellow("############## Details ##############")
	for _, d := range det {
		if strings.Contains(d.Attrs()["href"], "detail") {
			url = fmt.Sprintf("https://www.findchips.com%s", d.Attrs()["href"])
			resp, _ = soup.Get(url)
			doc = soup.HTMLParse(resp)
			det := doc.Find("ul", "class", "part-details-list")
			for _, d := range det.FindAll("li") {
				parm := d.Find("small").Text()
				val := strings.TrimSpace(d.Find("p").Text())
				fmt.Printf("\"%v\",\"%v\"\n", parm, val)
			}
			color.Yellow("############## Datasheet Link ##############")
			ds := doc.Find("div", "class", "datasheet-item").Find("a").Attrs()["href"]
			fmt.Println(ds)
			color.Yellow("############## Alternativen ##############")
			alts := doc.Find("table", "class", "part-suggestions").Find("tbody").FindAll("tr")
			for _, a := range alts {
				link := a.Find("td", "class", "td-col-1").Find("a")
				linkTxt := strings.TrimSpace(link.Text())
				fmt.Println(linkTxt)
			}
		}
	}
}

func getDistributorInfo(part string) {

	dists := []string{"Mouser", "Farnel", "Digi-Key", "Avnet", "TTI", "RS"}

	url := fmt.Sprintf("https://www.findchips.com/lite/%s", part)
	resp, _ := soup.Get(url)
	doc := soup.HTMLParse(resp)

	res := doc.FindAll("div", "class", "distributor-results")
	for _, r := range res {
		distributor := r.Find("h3", "class", "distributor-title").FullText()
		distributor = strings.TrimSpace(distributor)
		distributor = strings.Split(distributor, "\n")[0]

		if contains(distributor, dists) {

			distPartNos := r.FindAll("tr", "class", "row")
			color.Yellow("############## %s ##############", distributor)
			for _, p := range distPartNos {
				distributorPartNo := p.Attrs()["data-distino"]
				partNo := p.Attrs()["data-mfrpartnumber"]
				desc := p.Children()[5].Find("span", "class", "td-description").Text()
				fmt.Printf("\"%v\",\"%v\",\"%v\",\"%v\"\n", distributor, partNo, distributorPartNo, desc)
			}
		}
	}
}

func contains(input string, words []string) bool {
	for _, word := range words {
		if strings.Index(input, word) > -1 {
			return true
		}
	}
	return false
}
