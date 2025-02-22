/*
 * Copyright © 2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author       Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright  2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license  	   Apache-2.0
 */

package authz_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/tidwall/sjson"

	"github.com/ory/viper"

	"github.com/ory/oathkeeper/driver/configuration"
	"github.com/ory/oathkeeper/internal"
	"github.com/ory/oathkeeper/x"

	"github.com/ory/oathkeeper/pipeline/authn"
	. "github.com/ory/oathkeeper/pipeline/authz"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/oathkeeper/rule"
)

func TestAuthorizerSpiceDBWarden(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistry(conf)

	rule := &rule.Rule{ID: "TestAuthorizer"}

	a, err := reg.PipelineAuthorizer("spice_db")
	require.NoError(t, err)
	assert.Equal(t, "spice_db", a.GetID())

	for k, tc := range []struct {
		setup     func(t *testing.T) *httptest.Server
		r         *http.Request
		session   *authn.AuthenticationSession
		config    json.RawMessage
		expectErr bool
	}{
		{
			expectErr: true,
		},
		{
			config:    []byte(`{ "required_action": "action", "required_resource": "resource" }`),
			r:         &http.Request{URL: &url.URL{}},
			session:   new(authn.AuthenticationSession),
			expectErr: true,
		},
		{
			config: []byte(`{ "required_action": "action", "required_resource": "resource", "flavor": "regex" }`),
			r:      &http.Request{URL: &url.URL{}},
			setup: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusForbidden)
				}))
			},
			session:   new(authn.AuthenticationSession),
			expectErr: true,
		},
		{
			config: []byte(`{ "required_action": "action", "required_resource": "resource", "flavor": "exact" }`),
			r:      &http.Request{URL: &url.URL{}},
			setup: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Contains(t, r.Header, "Content-Type")
					assert.Contains(t, r.Header["Content-Type"], "application/json")
					assert.Contains(t, r.URL.Path, "exact")
					w.Write([]byte(`{"allowed":false}`))
				}))
			},
			session:   new(authn.AuthenticationSession),
			expectErr: true,
		},
		{
			config: []byte(`{ "required_action": "action:{{ printIndex .MatchContext.RegexpCaptureGroups (sub 1 1 | int)}}:{{ index .MatchContext.RegexpCaptureGroups (sub 2 1 | int)}}", "required_resource": "resource:{{ index .MatchContext.RegexpCaptureGroups 0}}:{{ index .MatchContext.RegexpCaptureGroups 1}}" }`),
			r:      &http.Request{URL: x.ParseURLOrPanic("https://localhost/api/users/1234/abcde")},
			setup: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var ki AuthorizerSpiceDBRequestBody
					require.NoError(t, json.NewDecoder(r.Body).Decode(&ki))
					assert.EqualValues(t, AuthorizerSpiceDBRequestBody{
						Action:   "action:1234:abcde",
						Resource: "resource:1234:abcde",
						Context:  map[string]interface{}{},
						Subject:  "peter",
					}, ki)
					assert.Contains(t, r.URL.Path, "regex")
					w.Write([]byte(`{"allowed":true}`))
				}))
			},
			session: &authn.AuthenticationSession{
				Subject: "peter",
				MatchContext: authn.MatchContext{
					RegexpCaptureGroups: []string{"1234", "abcde"},
				},
			},
			expectErr: false,
		},
		{
			config: []byte(`{ "required_action": "action:{{ index .MatchContext.RegexpCaptureGroups 0}}:{{ index .MatchContext.RegexpCaptureGroups 1}}", "required_resource": "resource:{{ index .MatchContext.RegexpCaptureGroups 0}}:{{ index .MatchContext.RegexpCaptureGroups 1}}", "subject": "{{ .Extra.name }}" }`),
			r:      &http.Request{URL: x.ParseURLOrPanic("https://localhost/api/users/1234/abcde")},
			setup: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var ki AuthorizerSpiceDBRequestBody
					require.NoError(t, json.NewDecoder(r.Body).Decode(&ki))
					assert.EqualValues(t, AuthorizerSpiceDBRequestBody{
						Action:   "action:1234:abcde",
						Resource: "resource:1234:abcde",
						Context:  map[string]interface{}{},
						Subject:  "peter",
					}, ki)
					assert.Contains(t, r.URL.Path, "regex")
					w.Write([]byte(`{"allowed":true}`))
				}))
			},
			session: &authn.AuthenticationSession{
				Extra: map[string]interface{}{"name": "peter"},
				MatchContext: authn.MatchContext{
					RegexpCaptureGroups: []string{"1234", "abcde"},
				}},
			expectErr: false,
		},
		{
			config: []byte(`{ "required_action": "action:{{ index .MatchContext.RegexpCaptureGroups 0 }}:{{ .Extra.name }}", "required_resource": "resource:{{ index .MatchContext.RegexpCaptureGroups 0}}:{{ .Extra.apiVersion }}", "subject": "{{ .Extra.name }}" }`),
			r:      &http.Request{URL: x.ParseURLOrPanic("https://localhost/api/users/1234/abcde?limit=10")},
			setup: func(t *testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var ki AuthorizerSpiceDBRequestBody
					require.NoError(t, json.NewDecoder(r.Body).Decode(&ki))
					assert.EqualValues(t, AuthorizerSpiceDBRequestBody{
						Action:   "action:1234:peter",
						Resource: "resource:1234:1.0",
						Context:  map[string]interface{}{},
						Subject:  "peter",
					}, ki)
					assert.Contains(t, r.URL.Path, "regex")
					w.Write([]byte(`{"allowed":true}`))
				}))
			},
			session: &authn.AuthenticationSession{
				Extra: map[string]interface{}{
					"name":       "peter",
					"apiVersion": "1.0"},
				MatchContext: authn.MatchContext{RegexpCaptureGroups: []string{"1234"}},
			},
			expectErr: false,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			baseURL := "http://73fa403f-7e9c-48ef-870f-d21b2c34fc80c6cb6404-bb36-4e70-8b90-45155657fda6/"
			if tc.setup != nil {
				ts := tc.setup(t)
				defer ts.Close()
				baseURL = ts.URL
			}

			a.(*AuthorizerSpiceDB).WithContextCreator(func(r *http.Request) map[string]interface{} {
				return map[string]interface{}{}
			})

			tc.config, _ = sjson.SetBytes(tc.config, "base_url", baseURL)
			err := a.Authorize(tc.r, tc.session, tc.config, rule)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	t.Run("method=validate", func(t *testing.T) {
		viper.Set(configuration.AuthorizerSpiceDBIsEnabled, false)
		require.Error(t, a.Validate(json.RawMessage(`{"base_url":"","required_action":"foo","required_resource":"bar"}`)))

		viper.Set(configuration.AuthorizerSpiceDBIsEnabled, false)
		require.Error(t, a.Validate(json.RawMessage(`{"base_url":"http://foo/bar","required_action":"foo","required_resource":"bar"}`)))

		viper.Reset()
		viper.Set(configuration.AuthorizerSpiceDBIsEnabled, true)
		require.Error(t, a.Validate(json.RawMessage(`{"base_url":"","required_action":"foo","required_resource":"bar"}`)))

		viper.Set(configuration.AuthorizerSpiceDBIsEnabled, true)
		require.NoError(t, a.Validate(json.RawMessage(`{"base_url":"http://foo/bar","required_action":"foo","required_resource":"bar"}`)))
	})
}
