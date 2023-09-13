#!/bin/bash

go build -o /tmp/deploy ./cmd/system-deploy 

DOCS=./docs/docs/actions/

function gendoc() {
    cat > ${DOCS}/$1.md <<EOT
---
layout: default
parent: Actions
title: $1
nav_order: 1
---
EOT
    /tmp/deploy describe $1 --markdown >> ${DOCS}/$1.md ;
}

gendoc InstallPackages
gendoc Platform
gendoc Systemd
gendoc Copy
gendoc Exec
gendoc OnChange
gendoc EditFile

cat > ./docs/docs/concepts/task-props.md <<EOT
---
layout: default
parent: Documentation
title: Task Properties
nav_order: 3
---
EOT

/tmp/deploy describe task --markdown >> ./docs/docs/concepts/task-props.md

# finally, update the search index as well
(cd docs && bundle exec just-the-docs rake search:init)
