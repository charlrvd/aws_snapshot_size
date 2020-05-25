#!/bin/bash

# comment out GOOS=darwin on the Dockerfile if not building the binary to run on macos

set -e

substitute=$1
sed -i '' -e "s!PLACEHOLDER!${substitute}!" Dockerfile-build


docker build -t artefact_container:latest -f Dockerfile-build .
docker run --rm --name temp -v "$(PWD)/artefacts":"/tmp/artefacts" artefact_container
ls -l ./artefacts
docker rmi $(docker images -f 'dangling=true' -q)

sed -i '' -e "s!${substitute}!PLACEHOLDER!" Dockerfile-build
