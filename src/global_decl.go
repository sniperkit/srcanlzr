// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"fmt"
	"reflect"
)

type GlobalDecl struct {
	Doc        []string `json:"doc,omitempty"`
	Name       *Ident   `json:"name"`
	Value      Expr     `json:"value,omitempty"`
	Type       *Ident   `json:"type"`
	Visibility string   `json:"visibility"`
}

func newGlobalDecl(m map[string]interface{}) (*GlobalDecl, error) {
	var err error
	errPrefix := "src/global_decl"
	globaldecl := GlobalDecl{}

	if globaldecl.Doc, err = extractStringSliceValue("doc", errPrefix, m); err != nil && isExist(err) {
		return nil, addDebugInfo(err)
	}

	nameMap, err := extractMapValue("name", errPrefix, m)
	if err != nil {
		return nil, addDebugInfo(err)
	}

	if globaldecl.Name, err = newIdent(nameMap); err != nil {
		return nil, addDebugInfo(err)
	}

	exprMap, err := extractMapValue("value", errPrefix, m)
	if err != nil && isExist(err) {
		return nil, addDebugInfo(err)
	} else if err == nil {
		if globaldecl.Value, err = newExpr(exprMap); err != nil {
			return nil, addDebugInfo(err)
		}
	}

	typeMap, err := extractMapValue("type", errPrefix, m)
	if err != nil {
		return nil, addDebugInfo(err)
	}

	if globaldecl.Type, err = newIdent(typeMap); err != nil {
		return nil, addDebugInfo(err)
	}

	if globaldecl.Visibility, err = extractStringValue("visibility", errPrefix, m); err != nil {
		return nil, addDebugInfo(err)
	}

	return &globaldecl, nil
}

func newGlobalDeclsSlice(key, errPrefix string, m map[string]interface{}) ([]*GlobalDecl, error) {
	var err error
	var s *reflect.Value

	if s, err = reflectSliceValue(key, errPrefix, m); err != nil {
		// XXX It is not possible to add debug info on this error because it is
		// required that this error be en "errNotExist".
		return nil, err
	}

	globaldecls := make([]*GlobalDecl, s.Len(), s.Len())
	for i := 0; i < s.Len(); i++ {
		globaldecl := s.Index(i).Interface()
		if globaldecl == nil {
			continue
		}

		switch globaldecl.(type) {
		case map[string]interface{}:
			if globaldecls[i], err = newGlobalDecl(globaldecl.(map[string]interface{})); err != nil {
				return nil, addDebugInfo(err)
			}
		default:
			return nil, addDebugInfo(fmt.Errorf(
				"%s: '%s' must be a map[string]interface{}, found %v",
				errPrefix, key, reflect.TypeOf(globaldecl)))
		}
	}

	return globaldecls, nil
}
