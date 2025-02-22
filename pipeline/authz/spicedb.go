package authz

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/x/httpx"

	"github.com/ory/oathkeeper/driver/configuration"
	"github.com/ory/oathkeeper/pipeline"
	"github.com/ory/oathkeeper/pipeline/authn"
	"github.com/ory/oathkeeper/x"

	"github.com/ory/x/urlx"

	"github.com/pkg/errors"
	"github.com/tomasen/realip"

	"github.com/ory/oathkeeper/helper"
)

type AuthorizerSpiceDBConfiguration struct {
	RequiredAction   string `json:"required_action"`
	RequiredResource string `json:"required_resource"`
	Subject          string `json:"subject"`
	Flavor           string `json:"flavor"`
	BaseURL          string `json:"base_url"`
}

func (c *AuthorizerSpiceDBConfiguration) SubjectTemplateID() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(c.Subject)))
}

func (c *AuthorizerSpiceDBConfiguration) ActionTemplateID() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(c.RequiredAction)))
}

func (c *AuthorizerSpiceDBConfiguration) ResourceTemplateID() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(c.RequiredResource)))
}

type AuthorizerSpiceDB struct {
	c configuration.Provider

	client         *retryablehttp.Client
	contextCreator authorizerSpiceDBWardenContext
	t              *template.Template
}

func NewAuthorizerSpiceDB(c configuration.Provider) *AuthorizerSpiceDB {
	return &AuthorizerSpiceDB{
		c: c,
		client: httpx.NewResilientClient(
			httpx.ResilientClientWithMaxRetryWait(100*time.Millisecond),
			httpx.ResilientClientWithMaxRetry(5),
		),
		contextCreator: func(r *http.Request) map[string]interface{} {
			return map[string]interface{}{
				"remoteIpAddress": realip.RealIP(r),
				"requestedAt":     time.Now().UTC(),
			}
		},
		t: x.NewTemplate("keto_engine_acp_ory"),
	}
}

func (a *AuthorizerSpiceDB) GetID() string {
	return "spice_db"
}

type authorizerSpiceDBWardenContext func(r *http.Request) map[string]interface{}

type AuthorizerSpiceDBRequestBody struct {
	Action   string                 `json:"action"`
	Context  map[string]interface{} `json:"context"`
	Resource string                 `json:"resource"`
	Subject  string                 `json:"subject"`
}

func (a *AuthorizerSpiceDB) WithContextCreator(f authorizerSpiceDBWardenContext) {
	a.contextCreator = f
}

func (a *AuthorizerSpiceDB) Authorize(r *http.Request, session *authn.AuthenticationSession, config json.RawMessage, rule pipeline.Rule) error {
	cf, err := a.Config(config)
	if err != nil {
		return err
	}

	// only Regexp matching strategy is supported for now.
	if !(a.c.AccessRuleMatchingStrategy() == "" || a.c.AccessRuleMatchingStrategy() == configuration.Regexp) {
		return helper.ErrNonRegexpMatchingStrategy
	}

	subject := session.Subject
	if cf.Subject != "" {
		subject, err = a.parseParameter(session, cf.SubjectTemplateID(), cf.Subject)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	action, err := a.parseParameter(session, cf.ActionTemplateID(), cf.RequiredAction)
	if err != nil {
		return errors.WithStack(err)
	}

	resource, err := a.parseParameter(session, cf.ResourceTemplateID(), cf.RequiredResource)
	if err != nil {
		return errors.WithStack(err)
	}

	flavor := "regex"
	if len(cf.Flavor) > 0 {
		flavor = cf.Flavor
	}

	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(&AuthorizerSpiceDBRequestBody{
		Action:   action,
		Resource: resource,
		Context:  a.contextCreator(r),
		Subject:  subject,
	}); err != nil {
		return errors.WithStack(err)
	}

	baseURL, err := url.ParseRequestURI(cf.BaseURL)
	if err != nil {
		return errors.WithStack(err)
	}

	req, err := http.NewRequest("POST", urlx.AppendPaths(baseURL, "/engines/acp/ory", flavor, "/allowed").String(), &b)
	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Add("Content-Type", "application/json")

	retryableReq, err := retryablehttp.FromRequest(req.WithContext(r.Context()))
	if err != nil {
		return errors.WithStack(err)
	}

	res, err := a.client.Do(retryableReq)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return errors.WithStack(helper.ErrForbidden)
	} else if res.StatusCode != http.StatusOK {
		return errors.Errorf("expected status code %d but got %d", http.StatusOK, res.StatusCode)
	}

	var result struct {
		Allowed bool `json:"allowed"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return errors.WithStack(err)
	}

	if !result.Allowed {
		return errors.WithStack(helper.ErrForbidden)
	}

	return nil
}

func (a *AuthorizerSpiceDB) parseParameter(session *authn.AuthenticationSession, templateID, templateString string) (string, error) {

	t := a.t.Lookup(templateID)
	if t == nil {
		var err error
		t, err = a.t.New(templateID).Parse(templateString)
		if err != nil {
			return "", err
		}
	}

	var b bytes.Buffer
	if err := t.Execute(&b, session); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (a *AuthorizerSpiceDB) Validate(config json.RawMessage) error {
	if !a.c.AuthorizerIsEnabled(a.GetID()) {
		return NewErrAuthorizerNotEnabled(a)
	}

	_, err := a.Config(config)
	return err
}

func (a *AuthorizerSpiceDB) Config(config json.RawMessage) (*AuthorizerSpiceDBConfiguration, error) {
	var c AuthorizerSpiceDBConfiguration
	if err := a.c.AuthorizerConfig(a.GetID(), config, &c); err != nil {
		return nil, NewErrAuthorizerMisconfigured(a, err)
	}

	if c.RequiredAction == "" {
		c.RequiredAction = "unset"
	}

	if c.RequiredResource == "" {
		c.RequiredResource = "unset"
	}

	return &c, nil
}
