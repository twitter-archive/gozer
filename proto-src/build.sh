#! /bin/bash
set -o errexit ; set -o nounset

MESOS_TAG="0.20.1"
MESOS_REPO="http://git-wip-us.apache.org/repos/asf/mesos.git"


error() { exec 1>&2 ; echo "$@" ; exit 1 ; }

# Check for required binaries in $PATH
protoc --help >/dev/null 2>&1 && rc=$? || rc=$?
[[ $rc -ne 0 ]] && error "Can't find 'protoc' binary in \$PATH, see README.md"

protoc-gen-go < /dev/null 2>&1 | fgrep -q 'protoc-gen-go: error:no files to generate' && rc=$? || rc=$?
[[ $rc -ne 0 ]] && error "Can't find 'protoc-gen-go' binary in \$PATH, see README.md"

# Check we are in the right place
[[ -d ../proto ]] || error "Can not find '../proto' directory"


# Grab the Mesos repo
[[ -d mesos-repo ]] || git clone "$MESOS_REPO" mesos-repo
info=$(cd mesos-repo && git checkout -q "$MESOS_TAG" && git log -1 --format='format:%H,%at')
mesos_sha=$(echo "${info}" | cut -d, -f1)
mesos_ts=$(echo "${info}" | cut -d, -f2-)

# These should be used to auto-generate part of the golang packages, such that they can be
# queried as to the version of the Mesos protobuf files they are generated/built against.
echo "Mesos SHA for tag ${MESOS_TAG} is ${mesos_sha}"
echo "Commit timestamp = ${mesos_ts}"

# Copy all .proto files local
for file in $(find mesos-repo -name \*.proto | fgrep -v 3rdparty/libprocess/3rdparty/stout/tests); do
    echo "Copying ${file}"
    cp "${file}" .
done

rm -rf mesos
mkdir mesos && mv mesos.proto mesos

# Build all the .proto files
for file in *.proto; do
    echo "Building ${file}..."
    base=$(basename "${file}" .proto)

    mkdir -p "../proto/${base}.pb"
    protoc --go_out="../proto/${base}.pb" "${file}"

    sed -i.bak -f sed-fix-mesos-pb-import.txt "../proto/${base}.pb/${base}.pb.go"
    rm -f "../proto/${base}.pb/${base}.pb.go.bak"
done

# Generate proto.go file so SHA/DateTime can be queried
sed -e "s/@GIT_SHA@/${mesos_sha}/" \
    -e "s/@GIT_TS@/${mesos_ts}/" \
    -e "s/@GIT_TAG@/${MESOS_TAG}/" \
    < proto-template.go > ../proto/proto.go

# Cleanup
rm -rf mesos-repo mesos
rm -f *.proto

exit 0
