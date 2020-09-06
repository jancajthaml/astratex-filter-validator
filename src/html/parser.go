// Copyright (c) 2020, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package html

import (
	"fmt"
	"bytes"
	"strconv"
	"strings"
	"regexp"

  "github.com/PuerkitoBio/goquery"
)

type HtmlParser struct {

}

func NewHtmlParser() HtmlParser {
  return HtmlParser{
  }
}

func (parser *HtmlParser) ScrapeCategoriesFrom(data []byte) ([]Category, error) {
  if parser == nil {
    return nil, fmt.Errorf("cannot call methods on nil reference")
  }

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	categories := make([]Category, 0)

	doc.Find(".CategoryTreeHorizontal .categoryTree a").Each(func(i int, s *goquery.Selection) {
		cat := Category{
			Title: strings.TrimSpace(s.Text()),
			Rel: s.AttrOr("href", ""),
		}
		if cat.Title != "" && cat.Rel != "" {
			categories = append(categories, cat)
		}
	})

	return categories, nil
}

func (parser *HtmlParser) ScrapeCategoryIdFrom(data []byte) (int, error) {
  if parser == nil {
    return -1, fmt.Errorf("cannot call methods on nil reference")
  }

	pattern := regexp.MustCompile(`gCat\s?=\s?([0-9]{1,10})\;`)

	groups := pattern.FindStringSubmatch(string(data))
	if len(groups) > 1 {
    i, err := strconv.Atoi(groups[1])
    if err == nil {
    	return i, nil
    }
	}

	return -1, fmt.Errorf("not found")
}

func (parser *HtmlParser) ScrapeFilterPropertiesFrom(data []byte) ([]FilterProperty, error) {
  if parser == nil {
    return nil, fmt.Errorf("cannot call methods on nil reference")
  }

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	properties := make([]FilterProperty, 0)

	doc.Find(".n-vlastnosti div a").Each(func(i int, s *goquery.Selection) {
		prop := FilterProperty{
			Value: strings.TrimSpace(s.Text()),
		}
		properties = append(properties, prop)
	})

	return properties, nil
}
