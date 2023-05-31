#!/bin/bash
set -o errexit

for config in $(find . -name ".mockdef.yml"); do
    for idx in $(yq -o tsv '.generate | keys' $config); do
        echo "Generating mocks for $(dirname $config)..."
        mockery \
            --dir $(dirname $config) \
            --inpackage \
            --with-expecter \
            --name $(yq ".generate.[$idx].interface" $config) \
            --structname $(yq ".generate.[$idx].mock" $config) \
            --filename $(yq ".generate.[$idx].target" $config)
    done
done

echo "Done!"
