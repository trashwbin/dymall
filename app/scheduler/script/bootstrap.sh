#! /usr/bin/env bash
CURDIR=$(cd $(dirname $0); pwd)
echo "$CURDIR/bin/scheduler"
exec "$CURDIR/bin/scheduler"
