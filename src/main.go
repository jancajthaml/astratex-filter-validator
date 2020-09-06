// Copyright (c) 2016-2020, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/jancajthaml/astratex-filter-validator/http"
	"github.com/jancajthaml/astratex-filter-validator/html"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)


var baseUri = ""

func getCategories() []html.Category {
	client := http.NewHttpClient()
	parser := html.NewHtmlParser()

	resp, err := client.Get(fmt.Sprintf("https://%s/", baseUri))
	if err != nil {
		return nil
	}

	categories, err := parser.ScrapeCategoriesFrom(resp.Data)
	if err != nil {
		return nil
	}

	return categories
}

func getFilterPropertiesForCategory(category html.Category) []html.FilterProperty {
	client := http.NewHttpClient()
	parser := html.NewHtmlParser()

	resp, err := client.Get(fmt.Sprintf("https://%s%s", baseUri, category.Rel))
	if err != nil {
		return nil
	}

	catId, err := parser.ScrapeCategoryIdFrom(resp.Data)
	if err != nil {
		return nil
	}

	category.Id = catId

	resp, err = client.Get(fmt.Sprintf("https://%s/ajax/commodityContent.aspx?cat=%d&paramFilterInit=1", baseUri, category.Id))
	if err != nil {
		return nil
	}

	properties, err := parser.ScrapeFilterPropertiesFrom(resp.Data)
	if err != nil {
		return nil
	}

	return properties

}

func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		panic("usage: \"<prog> www.astratex.gr\"")
	}

	f := excelize.NewFile()

	baseUri = args[0]
	categories := getCategories()
	for _, category := range categories {
		f.NewSheet(category.Title)
		f.SetSheetFormatPr(category.Title,
			excelize.BaseColWidth(1.0),
			excelize.DefaultColWidth(1.0),
			excelize.DefaultRowHeight(1.0),
			excelize.CustomHeight(false),
			excelize.ZeroHeight(false),
			excelize.ThickTop(true),
			excelize.ThickBottom(true),
		)

		var aWidth = 1
		var bWidth = 1

		properties := getFilterPropertiesForCategory(category)
		for idx, property := range properties {
			uri := fmt.Sprintf("https://%s%s", baseUri, category.Rel)
			if len(property.Value) > aWidth {
				aWidth = len(property.Value)
			}
			if len(uri) > bWidth {
				bWidth = len(uri)
			}
			f.SetCellValue(category.Title, fmt.Sprintf("A%d", idx+1), property.Value)
			f.SetCellValue(category.Title, fmt.Sprintf("B%d", idx+1), uri)
			f.SetCellHyperLink(category.Title, fmt.Sprintf("B%d", idx+1), uri, "External")
		}

		f.SetColWidth(category.Title, "A", "A", float64(aWidth*2))
		f.SetColWidth(category.Title, "B", "B", float64(bWidth*2))
	}

	f.DeleteSheet("Sheet1")

	if err := f.SaveAs(fmt.Sprintf("categories_filters_properties_%s.xlsx", baseUri)); err != nil {
    panic(err.Error())
  }
}
