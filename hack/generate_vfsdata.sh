#!/bin/sh

rm -f vendor/github.com/openshift/installer/data/assets_vfsdata.go
env GO111MODULE=off OPENSHIFT_INSTALL_DATA='' go generate ./vendor/github.com/openshift/installer/data
sed -i '1,4d' vendor/github.com/openshift/installer/data/assets_vfsdata.go
rm -f vendor/github.com/openshift/installer/data/{assets.go,assets_generate.go}
