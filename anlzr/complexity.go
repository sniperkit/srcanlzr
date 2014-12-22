// Copyright 2014 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anlzr

import (
	"fmt"

	"github.com/DevMine/srcanlzr/src"
)

type Complexity struct{}

func (c Complexity) Analyze(p *src.Project, r *Result) error {
	cm := ComplexityMetrics{}

	var totalFuncs int64
	var totalFiles int64

	var totalComplexityPerFunc int64
	var totalComplexityPerFile float32

	for _, pkg := range p.Packages {
		for _, sf := range pkg.SourceFiles {
			var fileComplexity int64
			var numFuncs int64

			for _, f := range sf.Functions {
				numFuncs++
				fileComplexity += functionCyclomaticComplexity(&f)

				for _, stmt := range f.StmtList {
					switch stmt.(type) {
					case src.IfStatement, src.ForStatement:
						fileComplexity++
					}
				}

				fileComplexity += int64(len(f.Return))
			}

			for _, cls := range sf.Classes {
				for _, m := range cls.Methods {
					numFuncs++
					fileComplexity += methodCyclomaticComplexity(m)

					for _, stmt := range m.StmtList {
						switch stmt.(type) {
						case src.IfStatement, src.ForStatement:
							fileComplexity++
						}
					}

					fileComplexity += int64(len(m.Return))

				}
			}

			for _, mod := range sf.Modules {
				for _, m := range mod.Methods {
					numFuncs++
					fileComplexity += methodCyclomaticComplexity(m)

					for _, stmt := range m.StmtList {
						switch stmt.(type) {
						case src.IfStatement, src.ForStatement:
							fileComplexity++
						}
					}

					fileComplexity += int64(len(m.Return))

				}
			}

			if numFuncs > 0 {
				totalFiles++
				totalFuncs += numFuncs
				totalComplexityPerFunc += fileComplexity
				totalComplexityPerFile += float32(fileComplexity) / float32(numFuncs)
			}
		}
	}

	cm.AveragePerFunc = float32(totalComplexityPerFunc) / float32(totalFuncs)
	cm.AveragePerFile = totalComplexityPerFile / float32(totalFiles)

	r.Complexity = cm

	return nil
}

func functionCyclomaticComplexity(f *src.Function) int64 {
	cc := int64(1) // cyclomatic complexity

	for _, s := range f.StmtList {
		cc += statementComplexity(&s)
	}

	return cc
}

func methodCyclomaticComplexity(m *src.Method) int64 {
	cc := int64(1) // cyclomatic complexity

	for _, s := range m.StmtList {
		cc += statementComplexity(&s)
	}

	return cc
}

func statementComplexity(s src.Statement) int64 {
	var c int64

	switch s.(type) {
	case src.IfStatement:
		fmt.Println("foo")
		c++

		stmt := s.(src.IfStatement)

		for _, s := range stmt.StmtList {
			c += statementComplexity(s)
		}
	case src.ForStatement:
		fmt.Println("bar")
		c++

		stmt := s.(src.ForStatement)

		for _, s := range stmt.StmtList {
			c += statementComplexity(s)
		}
	case src.CallStatement:
		fmt.Println("plop")
		c++
	}

	return c
}