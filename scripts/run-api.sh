#!/bin/sh

set -eux

docker-compose run --rm --service-ports api
