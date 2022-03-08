//go:build sass || all

// Copyright (c) 2022  The Go-Enjin Authors
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

package sass

import (
	"strings"

	"github.com/tdewolff/parse/v2/buffer"
	"github.com/urfave/cli/v2"
	"github.com/wellington/go-libsass"

	"github.com/go-enjin/be/pkg/feature"
	"github.com/go-enjin/be/pkg/log"
)

var _sass *Feature

var _ feature.Feature = (*Feature)(nil)

var _ feature.OutputTranslator = (*Feature)(nil)

const Tag feature.Tag = "OutputSass"

type Feature struct {
	feature.CFeature

	includePaths []string
}

type MakeFeature interface {
	feature.MakeFeature

	IncludePaths(paths ...string) MakeFeature
}

func New() MakeFeature {
	if _sass == nil {
		_sass = new(Feature)
		_sass.Init(_sass)
	}
	return _sass
}

func (f *Feature) IncludePaths(paths ...string) MakeFeature {
	f.includePaths = append(f.includePaths, paths...)
	return f
}

func (f *Feature) Tag() (tag feature.Tag) {
	tag = Tag
	return
}

func (f *Feature) Build(b feature.Buildable) (err error) {
	return
}

func (f *Feature) Startup(ctx *cli.Context) (err error) {
	return
}

func (f *Feature) CanTranslate(mime string) (ok bool) {
	ok = strings.Contains(mime, "text/x-scss")
	return
}

func (f *Feature) TranslateOutput(s feature.Service, input []byte, _ string) (output []byte, mime string, err error) {
	o := buffer.NewWriter([]byte{})
	r := buffer.NewReader(input)
	var comp libsass.Compiler
	if comp, err = libsass.New(o, r); err != nil {
		return
	}
	if len(f.includePaths) > 0 {
		log.DebugF("using sass include paths: %v", f.includePaths)
		if err = comp.Option(libsass.IncludePaths(f.includePaths)); err != nil {
			return
		}
	}
	if err = comp.Run(); err == nil {
		output = o.Bytes()
		mime = "text/css; charset=utf-8"
	}
	return
}