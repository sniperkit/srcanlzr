// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/DevMine/repotool/model"
	"github.com/DevMine/srcanlzr/src/ast"
)

// Decode a JSON encoded src.Project read from r.
func Decode(r io.Reader) (*Project, error) {
	dec := newDecoder(r)
	return dec.decode()
}

// Decode a JSON encoded src.Project read from a given file.
func DecodeFile(path string) (*Project, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}

type decoder struct {
	scan *scanner
	buf  []byte
	err  error
}

// newDecoder creates a new JSON decoder that reads from r.
func newDecoder(r io.Reader) *decoder {
	return &decoder{scan: newScanner(r)}
}

// decode decodes JSON input into a src.Project structure.
func (dec *decoder) decode() (*Project, error) {
	prj := dec.decodeProject()
	if dec.err != nil {
		return nil, dec.errorf(dec.err)
	}
	return prj, nil
}

func (dec *decoder) errorf(v interface{}) error {
	return fmt.Errorf("malformed json: %v", v)
}

// decodeProject decodes a project object.
func (dec *decoder) decodeProject() *Project {
	if !dec.assertNewObject() {
		return nil
	}

	prj := Project{}

	if dec.isEmptyObject() {
		return &prj
	}
	if dec.err != nil {
		return nil
	}

	for {
		key, err := dec.scan.nextKey()
		if err != nil {
			if err == io.EOF {
				break
			}
			dec.err = err
			return nil
		}
		if key == "" {
			dec.err = errors.New("empty key")
			return nil
		}

		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "packages":
			prj.Packages = dec.decodePackages()
		case "languages":
			prj.Langs = dec.decodeLanguages()
		case "repository":
			prj.Repo = dec.decodeRepository()
		case "loc":
			if tok != scanIntLit {
				dec.err = fmt.Errorf("expected integer literal, found %v", tok)
				return nil
			}
			prj.LoC, dec.err = dec.unmarshalInt(val)
		case "name":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected string literal, found %v", tok)
				return nil
			}
			prj.Name, dec.err = dec.unmarshalString(val)
		default:
			dec.err = errors.New("unexpected value for project object")
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}
	}
	return &prj
}

// decodePackages decodes a list of package objects.
func (dec *decoder) decodePackages() []*Package {
	if !dec.assertNewArray() {
		return nil
	}

	pkgs := []*Package{}

	if dec.isEmptyArray() {
		return pkgs
	}
	if dec.err != nil {
		return nil
	}

	for {
		pkg := dec.decodePackage()
		if dec.err != nil {
			return nil
		}
		pkgs = append(pkgs, pkg)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}

	return nil
}

// decoderPackage decodes a package object.
func (dec *decoder) decodePackage() *Package {
	if !dec.assertNewObject() {
		return nil
	}

	pkg := Package{}

	if dec.isEmptyObject() {
		return &pkg
	}
	if dec.err != nil {
		return nil
	}

	for {
		var key string
		key, dec.err = dec.scan.nextKey()
		if dec.err != nil {
			return nil
		}

		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "source_files":
			pkg.SrcFiles = dec.decodeSrcFiles()
		case "doc":
			pkg.Doc = dec.decodeStringsList()
		case "loc":
			if tok != scanIntLit {
				dec.err = fmt.Errorf("expected integer literal, found %v", tok)
				return nil
			}
			pkg.LoC, dec.err = dec.unmarshalInt(val)
		case "name":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected string literal, found %v", tok)
				return nil
			}
			pkg.Name, dec.err = dec.unmarshalString(val)
		default:
			dec.err = errors.New("unexpected value for project object")
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}
	}

	return &pkg
}

// decodeSrcFiles decodes a list of source file objects.
func (dec *decoder) decodeSrcFiles() []*SrcFile {
	if !dec.assertNewArray() {
		return nil
	}

	sf := []*SrcFile{}

	if dec.isEmptyArray() {
		return sf
	}
	if dec.err != nil {
		return nil
	}

	for {
		srcFile := dec.decodeSrcFile()
		if dec.err != nil {
			return nil
		}

		sf = append(sf, srcFile)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}

	return sf
}

func (dec *decoder) decodeSrcFile() *SrcFile {
	if !dec.assertNewObject() {
		return nil
	}

	sf := SrcFile{}

	if dec.isEmptyObject() {
		return &sf
	}
	if dec.err != nil {
		return nil
	}

	for {
		key, err := dec.scan.nextKey()
		if err != nil {
			if err == io.EOF {
				break
			}
			dec.err = err
			return nil
		}
		if key == "" {
			dec.err = errors.New("empty key")
			return nil
		}

		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "path":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected 'string literal', found '%v'", tok)
			}
			sf.Path, dec.err = dec.unmarshalString(val)
		case "language":
			dec.scan.back()
			sf.Lang = dec.decodeLanguage()
		case "imports":
			dec.scan.back()
			sf.Imports = dec.decodeStringsList()
		case "type_specifiers":
			dec.scan.back()
			sf.TypeSpecs = dec.decodeTypeSpecs()
		case "structs":
			dec.scan.back()
			sf.Structs = dec.decodeStructs()
		case "constants":
			dec.scan.back()
			sf.Constants = dec.decodeGlobalDecls()
		case "variables":
			dec.scan.back()
			sf.Vars = dec.decodeGlobalDecls()
		case "functions":
			dec.scan.back()
			sf.Funcs = dec.decodeFuncs()
		case "interfaces":
			dec.scan.back()
			sf.Interfaces = dec.decodeInterfaces()
		case "clases":
			dec.scan.back()
			sf.Classes = dec.decodeClassDecls()
		case "enums":
			dec.scan.back()
			sf.Enums = dec.decodeEnumDecls()
		case "traits":
			dec.scan.back()
			sf.Traits = dec.decodeTraits()
		case "loc":
			if tok != scanIntLit {
				dec.err = fmt.Errorf("expected integer literal, found %v", tok)
				return nil
			}
			sf.LoC, dec.err = dec.unmarshalInt(val)
		default:
			dec.err = fmt.Errorf("unexpected value for the key '%s' of a source file object", key)
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}

	}

	return &sf
}

// decodeTypeSpecs decodes a list of types specifiers objects.
func (dec *decoder) decodeTypeSpecs() []*ast.TypeSpec {
	if !dec.assertNewArray() {
		return nil
	}

	ts := []*ast.TypeSpec{}

	if dec.isEmptyArray() {
		return ts
	}
	if dec.err != nil {
		return nil
	}

	for {
		typeSpec := dec.decodeTypeSpec()
		if dec.err != nil {
			return nil
		}

		ts = append(ts, typeSpec)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}

	return ts
}

// decodeTypeSpec decodes a type specifier object.
func (dec *decoder) decodeTypeSpec() *ast.TypeSpec {
	if !dec.assertNewObject() {
		return nil
	}

	typeSpec := ast.TypeSpec{}

	if dec.isEmptyObject() {
		return &typeSpec
	}
	if dec.err != nil {
		return nil
	}

	for {
		key, err := dec.scan.nextKey()
		if err != nil {
			if err == io.EOF {
				break
			}
			dec.err = err
			return nil
		}
		if key == "" {
			dec.err = errors.New("empty key")
			return nil
		}

		_, _, err = dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "doc":
			dec.scan.back()
			typeSpec.Doc = dec.decodeStringsList()
		case "name":
			dec.scan.back()
			typeSpec.Name = dec.decodeIdent()
		case "type":
			dec.scan.back()
			typeSpec.Type = dec.decodeExpr()
		default:
			dec.err = fmt.Errorf("unexpected value for the key '%s' of a type specifier object", key)
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}
	}
	return &typeSpec
}

// TODO: implement
func (dec *decoder) decodeStructs() []*ast.StructType {
	return nil
}

// TODO: implement
func (dec *decoder) decodeGlobalDecls() []*ast.GlobalDecl {
	return nil
}

// TODO: implement
func (dec *decoder) decodeFuncs() []*ast.FuncDecl {
	return nil
}

// TODO: implement
func (dec *decoder) decodeInterfaces() []*ast.Interface {
	return nil
}

// TODO: implement
func (dec *decoder) decodeClassDecls() []*ast.ClassDecl {
	return nil
}

// TODO: implement
func (dec *decoder) decodeEnumDecls() []*ast.EnumDecl {
	return nil
}

// TODO: implement
func (dec *decoder) decodeTraits() []*ast.Trait {
	return nil
}

// TODO: implement
func (dec *decoder) decodeIdent() *ast.Ident {
	return nil
}

// TODO: implement
func (dec *decoder) decodeExpr() ast.Expr {
	return nil
}

// decoderStringsList decodes a list of strings.
func (dec *decoder) decodeStringsList() []string {
	if !dec.assertNewArray() {
		return nil
	}

	sl := []string{}

	if dec.isEmptyArray() {
		return sl
	}
	if dec.err != nil {
		return nil
	}

	for {
		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}
		if tok != scanStringLit {
			dec.err = fmt.Errorf("expected string, found %v", tok)
			return nil
		}
		sl = append(sl, string(val))

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}
	return sl
}

// decodeLanguages decodes a list of languages.
func (dec *decoder) decodeLanguages() []*Language {
	if !dec.assertNewArray() {
		return nil
	}

	ls := []*Language{}

	if dec.isEmptyArray() {
		return ls
	}
	if dec.err != nil {
		return nil
	}

	for {
		lang := dec.decodeLanguage()
		if dec.err != nil {
			return nil
		}

		ls = append(ls, lang)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}

	return ls
}

// decodeLanguage decode a src.Language object.
func (dec *decoder) decodeLanguage() *Language {
	if !dec.assertNewObject() {
		return nil
	}

	lang := Language{}

	if dec.isEmptyObject() {
		return &lang
	}
	if dec.err != nil {
		return nil
	}

	for {
		key, err := dec.scan.nextKey()
		if err != nil {
			if err == io.EOF {
				break
			}
			dec.err = err
			return nil
		}
		if key == "" {
			dec.err = errors.New("empty key")
			return nil
		}

		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "paradigms":
			// Since the '[' character has been consumed, we need to step back
			// brefore calling decodeStringsList.
			dec.scan.back()
			lang.Paradigms = dec.decodeStringsList()
		case "language":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected 'string literal', found '%v'", tok)
				return nil
			}
			lang.Lang, dec.err = dec.unmarshalString(val)
		default:
			dec.err = fmt.Errorf("unexpected value for the key '%s' of a language object", key)
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}
	}
	return &lang
}

// TODO: implement
func (dec *decoder) decodeRepository() *model.Repository {
	return nil
}

// TODO: implement
func (dec *decoder) unmarshalInt(data []byte) (int64, error) {
	return 0, nil
}

// unmarshalString unmarshals a bytes slice into a string.
func (dec *decoder) unmarshalString(data []byte) (string, error) {
	if data == nil {
		return "", errors.New("unable to unmarshal string: data is nil")
	}
	return string(data), nil
}

// assertNewObject makes sure that the next value is a new object. In other
// words, the next value must begin with a '{'. If it is not, it will set
// dec.err and return false.
func (dec *decoder) assertNewObject() bool {
	// Since Language is a JSON Object, we expect to find a '{' character.
	_, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return false
	}
	if tok != scanBeginObject {
		dec.err = fmt.Errorf("expected object, found '%v'", tok)
		return false
	}
	return true
}

// assertNewArray makes sure that the next value is a new array. In order
// words, the next value must begin with a '['. If it is not, it will set
// dec.err and return false.
func (dec *decoder) assertNewArray() bool {
	_, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return false
	}
	if tok != scanBeginArray {
		dec.err = fmt.Errorf("expected array, found '%v'", tok)
		return false
	}
	return true
}

// isEndObject returns true if the next value marks the end of the object
// ('}') and false otherwise. If it is false, the next value must be a
// comma. If not, it will set dec.err accordingly.
func (dec *decoder) isEndObject() bool {
	_, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return false
	}
	if tok == scanEndObject {
		return true
	}
	if tok != scanComma {
		dec.err = fmt.Errorf("expected 'comma', found '%v'", tok)
	}
	return false
}

// isEndArray returns true if the next value marks the end of the array
// (']') and false otherwise. If it is false, the next value must be a
// comma. If not, it will set dec.err accordingly.
func (dec *decoder) isEndArray() bool {
	_, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return false
	}
	if tok == scanEndArray {
		return true
	}
	if tok != scanComma {
		dec.err = fmt.Errorf("expected 'comma', found '%s'", tok)
	}
	return false
}

// isEmptyObject tests if the object is empty (no key/value pairs inside).
//
// This method does not consume any byte.
//
// If an error occurs, it returns false and set dec.err.
func (dec *decoder) isEmptyObject() bool {
	// The object can be empty, so we have to check for that and without
	// consuming the next byte.
	if b, err := dec.scan.peek(); err != nil {
		if err == io.EOF {
			dec.err = errors.New("unexpected EOF")
		} else {
			dec.err = err
		}
		return false
	} else if b == '}' {
		return true
	}
	return false
}

// isEmptyArray tests if the object is empty (no values inside).
//
// This method does not consume any byte.
//
// If an error occurs, it returns false and set dec.err.
func (dec *decoder) isEmptyArray() bool {
	if b, err := dec.scan.peek(); err != nil {
		if err == io.EOF {
			dec.err = errors.New("unexpected EOF")
		} else {
			dec.err = err
		}
		return false
	} else if b == ']' {
		return true
	}
	return false
}
