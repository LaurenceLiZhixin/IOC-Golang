/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package inject

import (
	"fmt"
	"io"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
)

// writeHeaderWithoutConstrain writes out the build tag, package declaration, and imports
func writeHeaderWithoutConstrain(pkg *loader.Package, out io.Writer, packageName string, imports *importsList, headerText string) {
	// NB(directxman12): blank line after build tags to distinguish them from comments
	_, err := fmt.Fprintf(out, `// Code generated by iocli, run 'iocli gen' to re-generate

package %[1]s

import (
%[2]s
)

`, packageName, strings.Join(imports.ImportSpecs(), "\n"), headerText)
	if err != nil {
		pkg.AddError(err)
	}

}

// writeHeader writes out the build tag, package declaration, and imports
func writeHeader(pkg *loader.Package, out io.Writer, packageName string, imports *importsList, headerText string) {
	// NB(directxman12): blank line after build tags to distinguish them from comments
	_, err := fmt.Fprintf(out, `//go:build !ignore_autogenerated
// +build !ignore_autogenerated

%[3]s

// Code generated by iocli, run 'iocli gen' to re-generate

package %[1]s

import (
%[2]s
)

`, packageName, strings.Join(imports.ImportSpecs(), "\n"), headerText)
	if err != nil {
		pkg.AddError(err)
	}

}

// writeMethods writes each method to the file, sorted by type name.
func writeMethods(pkg *loader.Package, out io.Writer, outBuffer []byte) {
	_, err := out.Write(outBuffer)
	if err != nil {
		pkg.AddError(err)
	}
}

// writeFormatted outputs the given code, after gofmt-ing it.  If we couldn't gofmt,
// we write the unformatted code for debugging purposes.
func writeOut(ctx *genall.GenerationContext, outputFile io.WriteCloser, root *loader.Package, outBytes []byte) {
	if outputFile == nil {
		var err error
		name := "zz_generated.ioc.go"
		outputFile, err = ctx.Open(root, name)
		if err != nil {
			root.AddError(err)
			return
		}
	}
	n, err := outputFile.Write(outBytes)
	if err != nil {
		root.AddError(err)
		return
	}
	if n < len(outBytes) {
		root.AddError(io.ErrShortWrite)
	}
}
