# Release patch series

The October 21 packaging commit pins Go 1.26.10 and the existing distroless
Debian 12 runtime by digest. The runtime receives the builder's clean CA bundle
because distroless still carries the affected roots. This complements Go's
eMudhra intermediate blocklist without changing application code or dependencies.
Drop this isolated commit when upstream supplies equivalent trust handling.

This fork follows released upstream Oathkeeper tags. The current base is
`v26.2.0` (`c84dbe07ecbf6f10154f04ec49b137a115155289`). Historical branches and
published tags are retained unchanged.

GitHub Actions are disabled for this fork. Release verification runs locally:
build the images, run the tests and vulnerability scans, then validate the
candidate before publishing it. Workflow files remain as upstream reference.

Use `release/v<upstream-version>-custom.<revision>` branches. Keep fork CI setup
in the first commit, followed by isolated behavior fixes with their regression
tests. Feature PRs target that release branch; the default branch does not carry
the release patch series. Upstream-only publishing, newsletter, and
project-board jobs require upstream credentials and are skipped in forks.

Keep each behavior fix independently suitable for upstream submission. Avoid
deployment-specific configuration, bundled patches, dependency churn, or
unrelated formatting. Preserve upstream security fixes during rebases.

Before publishing a candidate:

```sh
go test ./rule ./proxy ./api ./pipeline/...
go test -race ./rule ./proxy -run '^(TestRulePathNormalizationPreservesTrailingSlash|TestRulePathTraversal|TestConfigureBackendURL|TestInitializeSession)$'
go vet ./rule ./proxy ./api ./pipeline/...
```

The full race run also exposes pre-existing upstream watcher/logger and
authorizer races. Keep those findings separate from the focused regression gate;
do not describe the entire upstream race suite as passing.

The legacy `go-acc` coverage step can fail the unchanged asynchronous
introspection-cache hit assertion. Fork CI runs the rule/proxy/API unit tests
and focused race checks instead; linting and the upstream reload, end-to-end,
forwarded-header and bearer-token suites remain enabled. Full affected-package
tests above remain part of local release acceptance.

The upstream repository-wide Format workflow rewrites unrelated vendored source
on this base and is restricted to upstream. Format changed Go files with
`gofmt`/`goimports` and changed Markdown/YAML with the pinned Prettier version;
do not include whole-repository formatting churn in behavior fixes.

The initial candidate inherited dependency advisories from the upstream release.
The separate security update upgrades the affected JOSE, OpenTelemetry, gRPC,
MongoDB and Go extension modules, including their required dependencies. Run
`govulncheck ./...` and scan each built image before release; retain those
reports with the test results.

CI stamps the version and commit explicitly and disables automatic VCS version
inference. The unchanged module path otherwise derives an obsolete v0.x
pseudo-version from Git history, which scanners mistake for code predating the
security fixes in the release base. Dependency build information is retained.
The containerized Kubescape scanner receives a Docker archive of the built image
because it cannot access the runner's Docker daemon. No advisory exceptions are
applied. Kubescape reports the module-only advisory for unused OpenPGP code as
critical; use govulncheck to distinguish module presence from imported or called
vulnerable code. Hosted workflows remain disabled.

Run HTTP integration coverage for both decision and proxy modes, including
trailing slashes, path traversal normalization, authentication and rejection. An
image build must pass its regression gate before producing a final image.

For an upstream update, create a new release branch and rebase the isolated
commits. Review `git range-diff` against the previous base and series, rerun the
tests, and drop fixes that upstream now provides. Publish immutable signed
`v<upstream-version>-custom.<revision>` source tags; candidates use an `-rc.N`
suffix. Never move a published source or image tag.
