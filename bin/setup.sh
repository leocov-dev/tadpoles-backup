#!/usr/bin/env bash

set -e

python -m venv .venv

source .venv/bin/activate

pip install -U -r requirements.txt
pip install -U pip setuptools