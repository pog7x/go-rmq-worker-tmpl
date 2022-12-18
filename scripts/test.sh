#!/usr/bin/env bash
set -xue

COMPOSE_FILE="docker-compose.test.yml"
CMD="docker-compose -f ${COMPOSE_FILE}"

function docker_compose_down() {
    # Print containers logs in case of non-zero exit code.
    if [[ $? -ne 0 ]]; then
        ${CMD} logs app
    fi
    ${CMD} down
    ${CMD} rm -f
}

docker_compose_down

# The following command registers deferred function which is called at script exit
# even in case of the error.
trap docker_compose_down EXIT

${CMD} pull
${CMD} build --no-cache --pull
${CMD} up -d --no-color

${CMD} run --rm -T app make test
