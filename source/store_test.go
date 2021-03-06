/*
Copyright 2017 The Kubernetes Authors.

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

package source

import (
	"testing"

	"github.com/kubernetes-incubator/external-dns/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	t.Run("RegisterAndLookup", testRegisterAndLookup)
	t.Run("LookupMultiple", testLookupMultiple)
}

// testRegisterAndLookup tests that a Source can be registered and looked up by name.
func testRegisterAndLookup(t *testing.T) {
	for _, tc := range []struct {
		title            string
		givenAndExpected map[string]Source
	}{
		{
			"registered source is found by name",
			map[string]Source{
				"foo": &testutils.MockSource{},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			for k, v := range tc.givenAndExpected {
				Register(k, v)
			}

			for k, v := range tc.givenAndExpected {
				assert.Equal(t, v, Lookup(k))
			}
		})
	}
}

// testLookupMultiple tests that Sources can be looked up by providing multiple names.
func testLookupMultiple(t *testing.T) {
	for _, tc := range []struct {
		title       string
		registered  map[string]Source
		names       []string
		expectError bool
	}{
		{
			"multiple registered sources are found by names",
			map[string]Source{
				"foo": &testutils.MockSource{},
				"bar": &testutils.MockSource{},
			},
			[]string{"foo", "bar"},
			false,
		},
		{
			"multiple registered sources, one source not registered",
			map[string]Source{
				"foo": &testutils.MockSource{},
				"bar": &testutils.MockSource{},
			},
			[]string{"foo", "bar", "baz"},
			true,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			for k, v := range tc.registered {
				Register(k, v)
			}

			lookup, err := LookupMultiple(tc.names)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, lookup, len(tc.registered))
				for _, source := range tc.registered {
					assert.Contains(t, lookup, source)
				}
			}
		})
	}
}
