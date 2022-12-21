# Monorepo

Imagine we are software engineers at a devtool SaaS organization.
The software we write to support the organization
will take the shape of multiple programs.

We work incrementally,
choose technologies suited to their tasks that we also enjoy,
and are unafraid to write custom software
that improves our happiness and productivity.

## New repo

We initialize a monorepo on GitHub: `org/repo`:

```
.
└── README.md
```

We write a `README.md`:

```markdown
Add to your shell profile:

# Set environment variable to monorepo path

export ORG="$HOME/org"

# Prepend monorepo scripts

export PATH="$ORG/bin:$PATH"

Clone:

git clone https://github.com/org/repo.git $ORG
cd $ORG
```

The `$ORG` environment variable will be used throughout the codebase.

## Design

We design an architecture like this:

- a "Dashboard" web interface for customers (Svelte, TypeScript)
- SDKs for customers (Go, Node, Ruby)
- an HTTP server (Go)
- a Postgres database backing the HTTP API

```
.
├── ...
├── cmd
│   └── serverd
│       └── main.go
├── dashboard
├── sdk
│   ├── go
│   ├── node
│   └── ruby
└── server
    └── http_handler.go
```

## Commit messages

We use commit message [conventions](https://chris.beams.io/posts/git-commit/)
with a subject line prefix:

```
$ git log
dashboard: redirect to new project form after sign up
sdk/node: enable Keep-Alive in HTTP agent
server: document access to prod db
```

The prefix refers to the module that is changing.
It often matches a filesystem directory.
Sometimes the change will touch files across modules,
but the prefix should be "where the action is."

We try to keep the subject line under 50 characters and
[err on the side of brevity](https://golang.org/doc/effective_go.html#package-names)
for module names.

## Dependencies

To aid local development and help onboard teammates,
we provide a `setup.sh` script for each module.
In turn, each `setup.sh` script might invoke more general `$ORG/bin/setup-*` scripts.

Since everyone has `$ORG/bin` on their `$PATH`,
scripts in `bin/` are available in everyone's shells
without a directory prefix.

[ASDF](/asdf-version-manager)'s `.tool-versions` can pin versions for dependencies:

```
.
├── ...
├── .tool-versions
├── bin
│   ├── setup-go
│   ├── setup-node
│   ├── setup-postgres
│   └── setup-ruby
├── dashboard
│   └── setup.sh
├── sdk
│  ├── go
│  │   └── setup.sh
│  ├── node
│  │   └── setup.sh
│  └── ruby
│      └── setup.sh
└── server
    └── setup.sh
```

## Formatting

Our Go code is already formatted with [gofmt]
but we could be more consistent in our TypeScript and Ruby code.
We add config files for [Prettier] and [Rubocop]
at the top of the file hierarchy:

[gofmt]: https://golang.org/cmd/gofmt/
[prettier]: https://prettier.io/
[rubocop]: https://rubocop.readthedocs.io/en/latest/

```
.
├── ...
├── .rubocop.yml
└── .prettierrc.yml
```

## CI

To move quickly without breaking things,
we write automated tests
and run those tests continuously
when we integrate our changes.

We want our Continuous Integration (CI) service to:

- begin running the tests in < 5s after opening or editing a pull request
- run tests in an environment with [parity] to production
- run only the tests relevant to the change

[parity]: https://12factor.net/dev-prod-parity

Common causes of slow-starting tests on hosted CI services
are multi-tenant queues and containerization:
noisy neighbors backlog the queues
and containers have cache misses.

To meet our design goals, we write our own CI service, `testbot`:

```
.
├── ...
├── cmd
│   └── testbot
│       └── main.go
├── dashboard
│   └── Testfile
├── sdk
│   ├── go
│   │   └── Testfile
│   ├── node
│   │   └── Testfile
│   └── ruby
│       └── Testfile
├── server
│   └── Testfile
└── testbot
    ├── Testfile
    ├── farmer
    │   ├── Procfile
    │   └── main.go
    └── worker
        └── main.go
```

We open a GitHub pull request to add a new feature to the SDKs:

```
sdk/go/account.go             | 10 +++++-----
sdk/go/account_test.go        | 10 +++++-----
sdk/node/src/account.ts       | 10 +++++-----
sdk/node/test/account.ts      | 10 +++++-----
sdk/ruby/lib/account.rb       | 10 +++++-----
sdk/ruby/spec/account_spec.rb | 10 +++++-----
6 files changed, 60 insertions(+), 60 deletions(-)
```

A `testbot farmer` process on a server
responds to the GitHub webhook by:

- identifying the directories containing files that have changed
- walking up the file hierarchy to find `Testfile`s for changed directories
- saving test jobs to its backing Postgres database

Each `Testfile` defines test jobs for its directory. Ours might be:

```
$ cat $ORG/sdk/go/Testfile
tests: cd $ORG/sdk && go test -cover ./...
$ cat $ORG/sdk/node/Testfile
tests: cd $ORG/sdk/node && npm install && npm test
$ cat $ORG/sdk/ruby/Testfile
tests: $ORG/sdk/ruby/test.sh
```

A single `Testfile` can define multiple test jobs.
As test scripts become lengthier,
it is convenient to extract them to their own script.

Each line contains a name and command
separated by a colon.
The name appears in GitHub checks.
The command is run by a `testbot worker`,
which is continuously long polling `testbot farmer` via HTTP over TLS,
asking for test jobs.

Each `testbot worker` runs on Heroku with Go, Node, and Ruby buildpacks.
We run more instances to increase test parallelism.

In this example,
the tests for the Go, Node, and Ruby SDKs
will begin to run almost simultaneously
as different `testbot worker` processes pick them up.

Tests for `dashboard` and `server` will not run in this pull request
because no files in their directories were changed.

## Testing across services

To make it convenient to test across service boundaries,
we write a `with-serverd` script that:

- installs the `serverd` binary
- migrates the database
- creates a team and credential
- runs `serverd serve` without blocking programs
  passed as arguments to `with-serverd`

We place this script in `$ORG/bin` to make it available on our `$PATH`.

This script can be used in `Testfile`s:

```
$ cat $ORG/dashboard/Testfile
tests: cd $ORG/dashboard && with-serverd ./test.sh
$ cat $ORG/sdk/go/Testfile
tests: cd $ORG/sdk && with-serverd go test -cover ./...
```

Tests that depend on `with-serverd`
can avoid mocking on HTTP boundaries,
making actual requests to the backing service
and generating observable logs.

To broadly cover the product's surface area,
we move the Go SDK's `Testfile`
to the top of the file hierarchy
so its test jobs will run on all pull requests:

```
.
├── ...
├── Testfile
├── bin
│   └── with-serverd
└── sdk
    └── go
```

## Open source

Our SDKs should be open source on GitHub
and available on registries such as [NPM](https://www.npmjs.com/)
and [Rubygems](https://rubygems.org/)
but we like the monorepo as the canonical place
for our development workflow.

So, we mirror the SDK code to open source repos.
This provides community access.
Although the community patch workflow is a bit manual,
these are SDKs tightly coupled to our product
and we expect community patches to be rare.

To meet our design goals, we write a `mirrorbot` program:

```
.
├── ...
├── .mirrorbot.yml
└── cmd
    └── mirrorbot
        ├── Procfile
        ├── Testfile
        └── main.go
```

`mirrorbot` runs as a command-line program
that copies relevant changes from the `$ORG` monorepo
to one or more target repos.
It can be run on a laptop or a server.

On startup,
it clones the monorepo and each target repo.
Each copied commit includes an `upstream:[SHA]` line ([example][mirror]).
`mirrorbot` reads that SHA from the latest commit
on the destination branch of each target repo
to get its cursor.

[mirror]: https://github.com/sequence/sequence-sdk-ruby/commit/d00d15292808b7cf3c69c879ae9da781c819e658

Mirror then fast-forwards the local monorepo,
looking for new commits.

For each target repo,
a patch file is produced for each new commit.
The patch file is then filtered to remove anything not applicable to that repo.
If anything remains,
the patch is applied and committed.
Any new commits are then pushed to GitHub.

A `.mirrorbot.yml` file configures the program.
It maps monorepo branches to target repo branches
and monorepo directories and files to target repo directories and files.
For example:

```yaml
github.com/org/org-sdk-node:
  - branch:
      - main: main
  - mirror:
      - sdk/node: /
github.com/org/org-sdk-ruby:
  - branch:
      - main: main
  - mirror:
      - sdk/ruby: /
```

## Backwards compatibility

As our customers use the software,
we better understand their needs
and begin to design v2.

To do this well,
we add new interfaces and deprecate old interfaces
in the v1.x series of the SDKs.
The server continues to support the v1.x series now and
for a well-communicated time period past the release of v2.0
(such as 90 days).

As we release minor versions v1.1 and v1.2,
patch versions v1.2.1 and v1.2.2,
and eventually v2,
we want to carefully ensure backwards compatibility.

To meet our design goals,
we write `with-go-sdk`, `with-node-sdk`, and `with-ruby-sdk` scripts:

```
.
├── ...
└── bin
    ├── with-go-sdk
    ├── with-node-sdk
    └── with-ruby-sdk
```

We use these scripts to run processes using a given version of our SDKs:

```
with-go-sdk [version] [command]
with-node-sdk [version] [command]
with-ruby-sdk [version] [command]
```

The `version` can be a released version to a registry such as NPM or Rubygems
or a special value `monorepo` to use the current `$ORG` source code.

For example:

```bash
with-ruby-sdk 1.0 ruby -e "require 'org-sdk'; puts OrgSDK::VERSION"
```

We update the `Testfile` at the top of the file hierarchy
to test supported releases
and the current source code (pre-release)
on every pull request:

```
$ cat $ORG/Testfile
gohead: with-serverd with-go-sdk monorepo go test -cover ./...
gov1: with-serverd with-go-sdk 1.5 go test -cover ./...
gov2: with-serverd with-go-sdk 2 go test -cover ./...
# ... etc.
rubyv2: with-serverd with-ruby-sdk 2 $ORG/sdk/ruby/test.sh
```

The SDK scripts are composable with the `with-serverd` script.

In addition to running on CI,
they are useful for quickly testing bug reports
from customers on a particular version,
and hopefully providing a great customer experience for them.

## Deploy

We'll deploy `mirrorbot`, `serverd`, and `testbot` to
[Heroku](https://www.heroku.com/),
which expects a `Procfile` manifest for each program:

```
.
├── ...
├── cmd
│   ├── mirrorbot
│   │   ├── Procfile
│   │   └── main.go
│   ├── serverd
│   │   ├── Procfile
│   │   └── main.go
│   └── testbot
│       ├── Procfile
│       └── main.go
```

Heroku uses a stack of "buildpacks" to set up its build system.
We use a "monorepo" buildpack first in the stack,
which uses an `APP_BASE=subdir` config variable to tell Heroku
where to find the `Procfile` program.

```bash
heroku buildpacks:add -a <app> https://github.com/lstoll/heroku-buildpack-monorepo
heroku config:set APP_BASE=cmd/serverd -a <app>
```

We use the Go buildpack second in the stack,
which will install our Go program as a binary in `bin/program`.

```bash
heroku buildpacks:add -a <app> heroku/go
```

Our `Procfile` contains, for example:

```bash
web: bin/mirrorbot
```

## Conclusion

Working in a monorepo can encourage a feeling of "tight integration",
where service boundaries are well-defined and less likely to be mocked out.
Some tasks are a particular pleasure,
such as searching across projects
for callsites of a function or RPC.

For other tasks,
it can help to write custom tools.
Writing those custom tools offers an opportunity
to design an ideal experience for the engineering team.

The final directory structure looks like this:

```
.
├── .mirrorbot.yml
├── .prettierrc.yml
├── .rubocop.yml
├── .tool-versions
├── README.md
├── Testfile
├── bin
│   ├── setup-go
│   ├── setup-node
│   ├── setup-postgres
│   ├── setup-ruby
│   ├── with-go-sdk
│   ├── with-node-sdk
│   ├── with-ruby-sdk
│   └── with-serverd
├── cmd
│   ├── mirrorbot
│   │   ├── Procfile
│   │   ├── Testfile
│   │   ├── main.go
│   │   └── setup.sh
│   ├── serverd
│   │   ├── Procfile
│   │   └── main.go
│   └── testbot
│       ├── Procfile
│       └── main.go
├── dashboard
│   ├── Testfile
│   └── setup.sh
├── sdk
│   ├── go
│   │   └── setup.sh
│   ├── node
│   │   ├── Testfile
│   │   └── setup.sh
│   └── ruby
│       ├── Testfile
│       └── setup.sh
├── server
│   ├── Testfile
│   ├── http_handler.go
│   └── setup.sh
└── testbot
    ├── Testfile
    ├── farmer
    │   └── main.go
    ├── setup.sh
    └── worker
        └── main.go
```
