# #!/bin/bash

# DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# projects_array=("artifacts" "built-in-roles" "common" "patching" "deployment-producer" "executor" \
# "migration-producer" "test-infra" "verification-producer" \
# "ashandler" "digester" "infra-producer" "report-producer" \
# "tfhandler" "kunlun")

# for i in "${projects_array[@]}"
# do
#   go get -d github.com/kun-lun/$i
#   pushd $GOPATH/src/github.com/kun-lun/$i
#     echo "now checking out ${i}"
#     git remote remove origin
#     git remote add origin git@github.com:kun-lun/$i.git
#     # checkout the draft branch for now, remove this after v0.1 merged to the master.
#     git fetch
#     git checkout draft
#     git branch --set-upstream-to=origin/draft draft
#     git pull
#   popd
# done

# go get github.com/onsi/ginkgo/ginkgo
# go get github.com/onsi/gomega/...
