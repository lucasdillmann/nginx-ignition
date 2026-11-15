#!/bin/sh
# Verifies that every installer translation defines exactly the same strings as
# the canonical English file.
#
# Usage: dist/windows/lang/check.sh

set -eu

cd "$(dirname "$0")"

reference="en.nsh"
status=0
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

keys() {
    sed -n 's/^LangString[[:space:]]\{1,\}\([A-Za-z0-9_]\{1,\}\).*/\1/p' "$1" | sort
}

keys "$reference" > "$tmp/expected"
expected_count=$(wc -l < "$tmp/expected" | tr -d ' ')

for file in *.nsh; do
    [ "$file" = "$reference" ] && continue

    keys "$file" > "$tmp/actual"

    comm -23 "$tmp/expected" "$tmp/actual" > "$tmp/missing"
    comm -13 "$tmp/expected" "$tmp/actual" > "$tmp/extra"
    uniq -d "$tmp/actual" > "$tmp/duplicate"

    if [ -s "$tmp/missing" ]; then
        echo "$file: missing strings:"
        sed 's/^/  /' "$tmp/missing"
        status=1
    fi

    if [ -s "$tmp/extra" ]; then
        echo "$file: unexpected strings:"
        sed 's/^/  /' "$tmp/extra"
        status=1
    fi

    if [ -s "$tmp/duplicate" ]; then
        echo "$file: duplicate strings:"
        sed 's/^/  /' "$tmp/duplicate"
        status=1
    fi

    actual_count=$(wc -l < "$tmp/actual" | tr -d ' ')
    if [ "$actual_count" -ne "$expected_count" ]; then
        echo "$file: defines $actual_count strings, $reference defines $expected_count"
        status=1
    fi
done

if [ "$status" -eq 0 ]; then
    echo "All $expected_count installer strings are present in every language."
fi

exit "$status"
