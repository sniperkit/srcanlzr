// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"errors"
	"fmt"
	"reflect"
)

type Trait struct {
	Name       string       `json:"name"`
	Attributes []*Attribute `json:"attributes"`
	Methods    []*Method    `json:"methods"`
	Classes    []*Class     `json:"classes"`
	Traits     []*Trait     `json:"traits"`
}

func newTrait(m map[string]interface{}) (*Trait, error) {
	var err error
	errPrefix := "src/trait"
	trait := Trait{}

	if trait.Name, err = extractStringValue("name", errPrefix, m); err != nil {
		return nil, err
	}

	if trait.Attributes, err = newAttributesSlice("attributes", errPrefix, m); err != nil {
		return nil, err
	}

	if trait.Methods, err = newMethodsSlice("methods", errPrefix, m); err != nil {
		return nil, err
	}

	if trait.Classes, err = newClassesSlice("classes", errPrefix, m); err != nil {
		return nil, err
	}

	if trait.Traits, err = newTraitsSlice("traits", errPrefix, m); err != nil {
		return nil, err
	}

	return &trait, nil
}

func newTraitsSlice(key, errPrefix string, m map[string]interface{}) ([]*Trait, error) {
	var err error
	var s reflect.Value

	traitsMap, ok := m[key]
	if !ok {
		return nil, errNotExist
	}

	if s = reflect.ValueOf(traitsMap); s.Kind() != reflect.Slice {
		return nil, errors.New(fmt.Sprintf("%s: field '%s' is supposed to be a slice",
			errPrefix, key))
	}

	traits := make([]*Trait, s.Len(), s.Len())
	for i := 0; i < s.Len(); i++ {
		trait := s.Index(i).Interface()

		switch trait.(type) {
		case map[string]interface{}:
			if traits[i], err = newTrait(trait.(map[string]interface{})); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New(fmt.Sprintf("%s: '%s' must be a map[string]interface{}",
				errPrefix, key))
		}
	}

	return traits, nil
}
