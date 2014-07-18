#! /bin/bash
set -o errexit ; set -o nounset

error() { exec 1>&2 ; echo "$@" ; exit 1 ; }

# Check for required binaries in $PATH
protoc --help >/dev/null 2>&1 && rc=$? || rc=$?
[[ $rc -ne 0 ]] && error "Can't find 'protoc' binary in \$PATH, see README.md"

protoc-gen-go < /dev/null 2>&1 | fgrep -q 'protoc-gen-go: error:no files to generate' && rc=$? || rc=$?
[[ $rc -ne 0 ]] && error "Can't find 'protoc-gen-go' binary in \$PATH, see README.md"


# Check we are in the right place
[[ -f mesos.proto ]] || error "Can not find file 'mesos.proto' in current directory"
[[ -d ../proto ]] || error "Can not find '../proto' directory"


# Build all .proto files
for file in *.proto; do
    echo "Building ${file}..."
    base=$(basename $file .proto)

    mkdir -p "../proto/${base}.pb"
    protoc --go_out="../proto/${base}.pb" "${base}.proto"

    sed -i.bak -f sed-fix-mesos-pb-import.txt "../proto/${base}.pb/${base}.pb.go"
    rm -f "../proto/${base}.pb/${base}.pb.go.bak"
done

exit 0
