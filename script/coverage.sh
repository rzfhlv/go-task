#!/bin/bash
set -e

EXCLUDE_DIRS="cmd/ docs/ config/ script/ /mocks/ /infrastructure/ /presenter/"

cp coverage.out coverage_filtered.out
for d in $EXCLUDE_DIRS; do
    grep -v "$d" coverage_filtered.out > tmp.out && mv tmp.out coverage_filtered.out
done

go tool cover -func=coverage_filtered.out