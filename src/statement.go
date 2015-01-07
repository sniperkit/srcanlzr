// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	IfStmtName     = "IF"
	LoopStmtName   = "LOOP"
	AssignStmtName = "ASSIGN"
	CallStmtName   = "CALL"
	OtherStmtName  = "OTHER"
)

type Statement interface{}

func newStatement(m map[string]interface{}) (Statement, error) {
	errPrefix := "src/statement"

	typ, ok := m["type"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s: field 'type' does not exist", errPrefix))
	}

	switch typ {
	case IfStmtName:
		return newIfStatement(m)
	case LoopStmtName:
		return newLoopStatement(m)
	case AssignStmtName:
		return newAssignStatement(m)
	case CallStmtName:
		return newCallStatement(m)
	case OtherStmtName:
		return newOtherStatement(m)
	}

	return nil, errors.New("unknown statement type")
}

func newStatementsSlice(key, errPrefix string, m map[string]interface{}) ([]Statement, error) {
	var err error
	var s reflect.Value

	stmtsMap, ok := m[key]
	if !ok || stmtsMap == nil {
		return nil, errNotExist
	}

	if s = reflect.ValueOf(stmtsMap); s.Kind() != reflect.Slice {
		return nil, errors.New(fmt.Sprintf("%s: field '%s' is supposed to be a slice",
			errPrefix, key))
	}

	stmts := make([]Statement, s.Len(), s.Len())
	for i := 0; i < s.Len(); i++ {
		stmt := s.Index(i).Interface()

		switch stmt.(type) {
		case map[string]interface{}:
			if stmts[i], err = newStatement(stmt.(map[string]interface{})); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New(fmt.Sprintf("%s: '%s' must be a map[string]interface{}",
				errPrefix, key))
		}
	}

	return stmts, nil
}
