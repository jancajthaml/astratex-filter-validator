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

package validator

import (
	"fmt"
	"unicode"
	"strings"
)

type PropertiesValidator struct {
	Locale string
}

func NewPropertiesValidator(locale string) PropertiesValidator {
  return PropertiesValidator{
  	Locale: locale,
  }
}

func (validator *PropertiesValidator) Validate(value string) (bool, error) {
  if validator == nil {
    return false, fmt.Errorf("cannot call methods on nil reference")
  }

  canon := strings.Replace(value, "DEN", "", -1)
  if strings.ToLower(canon) != canon {
  	return false, nil
  }

  if strings.Contains(canon, "?") {
  	return false, nil
  }

  if validator.Locale == "el-GR" {
  	for _, r := range canon {
	  	if unicode.IsDigit(r) || unicode.IsSpace(r) || unicode.IsSymbol(r) || r == '/' {
	  		continue
	  	}
	  	if !unicode.In(r, unicode.Scripts["Greek"]) {
	  		return false, nil
	  	}
	  }
  }

	return true, nil
}
