/*
Copyright 2011 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jsonconfig

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestIncludes(t *testing.T) {
	obj, err := ReadFile("testdata/include1.json")
	if err != nil {
		t.Fatal(err)
	}
	two := obj.RequiredObject("two")
	if err := obj.Validate(); err != nil {
		t.Error(err)
	}
	if g, e := two.RequiredString("key"), "value"; g != e {
		t.Errorf("sub object key = %q; want %q", g, e)
	}
}

func TestIncludeLoop(t *testing.T) {
	_, err := ReadFile("testdata/loop1.json")
	if err == nil {
		t.Fatal("expected an error about import cycles.")
	}
	if !strings.Contains(err.Error(), "include cycle detected"){
		t.Fatal("expected an error about import cycles; got: %v", err)
	}
}

func TestBoolEnvs(t *testing.T) {
	os.Setenv("TEST_EMPTY", "")
	os.Setenv("TEST_TRUE", "true")
	obj, err := ReadFile("testdata/boolenv.json")
	if err != nil {
		t.Fatal(err)
	}
	if str := obj.RequiredString("str"); str != "" {
		t.Errorf("str = %q, want empty", str)
	}
	if v := obj.RequiredBool("false"); v != false {
		t.Error("key 'false' is true")
	}
	if v := obj.RequiredBool("true"); v != true {
		t.Error("key 'true' is false")
	}
	if err := obj.Validate(); err != nil {
		t.Error(err)
	}
}

func TestListExpansion(t *testing.T) {
	os.Setenv("TEST_BAR", "bar")
	obj, err := ReadFile("testdata/listexpand.json")
	if err != nil {
		t.Fatal(err)
	}
	s := obj.RequiredString("str")
	l := obj.RequiredList("list")
	if err := obj.Validate(); err != nil {
                t.Error(err)
        }
	want := []string{"foo", "bar"}
	if !reflect.DeepEqual(l, want) {
		t.Errorf("got = %#v\nwant = %#v", l, want)
	}
	if s != "bar" {
		t.Errorf("str = %q, want %q", s, "bar")
	}
}
