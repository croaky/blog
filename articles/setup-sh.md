# setup.sh

A contributor should be able to clone a [monorepo](monorepo),
change into the project directory,
and run one command to start contributing:

```
git clone git@github.com:org/monorepo.git
cd monorepo/project
./setup.sh
```

Example:

```
#!/bin/bash

# Set up Rails app after cloning codebase.

set -eou

# Install Ruby dependencies
gem install bundler --conservative
bundle check || bundle install

# Install JavaScript dependencies
bin/yarn

# Set up database
bin/rails db:setup
```
