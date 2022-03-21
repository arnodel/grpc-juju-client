#!/bin/bash

# This is just so that Python can find the generated code
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
export PYTHONPATH=$(realpath $SCRIPT_DIR/../../gen/proto/python)

python3 client.py "$@"
