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

package authn_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ory/oathkeeper/driver"
	"github.com/ory/x/configx"
	"github.com/ory/x/logrusx"

	"github.com/ory/oathkeeper/driver/configuration"
	"github.com/ory/oathkeeper/pipeline/authn"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticatorAnonymous(t *testing.T) {
	conf, err := configuration.NewViperProvider(context.Background(), logrusx.New("", ""),
		configx.WithValue("log.level", "debug"),
		configx.WithValue(configuration.ViperKeyErrorsJSONIsEnabled, true))
	require.NoError(t, err)

	reg := driver.NewRegistryMemory().WithConfig(conf)

	session := new(authn.AuthenticationSession)

	a, err := reg.PipelineAuthenticator("anonymous")
	require.NoError(t, err)
	assert.Equal(t, "anonymous", a.GetID())

	t.Run("method=authenticate/case=is anonymous user", func(t *testing.T) {
		err := a.Authenticate(
			&http.Request{Header: http.Header{}},
			session,
			json.RawMessage(`{"subject":"anon"}`),
			nil)
		require.NoError(t, err)
		assert.Equal(t, "anon", session.Subject)
	})

	t.Run("method=authenticate/case=has credentials", func(t *testing.T) {
		err := a.Authenticate(
			&http.Request{Header: http.Header{"Authorization": {"foo"}}},
			session,
			json.RawMessage(`{"subject":"anon"}`),
			nil)
		require.Error(t, err)
	})

	t.Run("method=validate", func(t *testing.T) {
		conf.Source().Set(configuration.ViperKeyAuthenticatorAnonymousIsEnabled, true)
		require.NoError(t, a.Validate(json.RawMessage(`{"subject":"foo"}`)))

		conf.Source().Set(configuration.ViperKeyAuthenticatorAnonymousIsEnabled, false)
		require.Error(t, a.Validate(json.RawMessage(`{"subject":"foo"}`)))
	})
}
