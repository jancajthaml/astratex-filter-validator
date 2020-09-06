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
	"github.com/jancajthaml/astratex-filter-validator/validator"
)

//var baseUri = "www.astratex.gr"

func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		panic("usage: \"<prog> www.astratex.gr\"")
	}

	baseUri := args[0]

	//fmt.Println("Hello")

	client := http.NewHttpClient()
	parser := html.NewHtmlParser()


	resp, err := client.Get(fmt.Sprintf("https://%s/", baseUri))
	if err != nil {
		panic(err.Error())
	}

	categories, err := parser.ScrapeCategoriesFrom(resp.Data)
	if err != nil {
		panic(err.Error())
	}

	locale, err := parser.ScrapeLocaleFrom(resp.Data)
	if err != nil {
		panic(err.Error())
	}

	check := validator.NewPropertiesValidator(locale)

	for _, category := range categories {
		//fmt.Printf("Validating category %s at %s\n", category.Title, category.Rel)

		resp, err = client.Get(fmt.Sprintf("https://%s%s", baseUri, category.Rel))
		if err != nil {
			panic(err.Error())
		}

		catId, err := parser.ScrapeCategoryIdFrom(resp.Data)
		if err != nil {
			panic(err.Error())
		}

		category.Id = catId

		resp, err = client.Get(fmt.Sprintf("https://%s/ajax/commodityContent.aspx?cat=%d&paramFilterInit=1", baseUri, category.Id))
		if err != nil {
			panic(err.Error())
		}

		properties, err := parser.ScrapeFilterPropertiesFrom(resp.Data)
		if err != nil {
			panic(err.Error())
		}

		//fmt.Printf("Found %d filter properties\n", len(properties))

		for _, property := range properties {
			ok, err := check.Validate(property.Value)
			if err != nil {
				panic(err.Error())
			}

			if !ok {
				fmt.Printf("Invalid filter property \"%s\" of category %s at https://%s%s\n", property.Value, category.Title, baseUri, category.Rel)
			}
		}

	}

	//fmt.Println("Bye")
}
