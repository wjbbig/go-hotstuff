#!/usr/bin/env bash

go build -o ../cmd/hotstuffgenkey/hotstuffgenkey ../cmd/hotstuffgenkey/hotstuffgenkey.go

../cmd/hotstuffgenkey/hotstuffgenkey -filepath /opt/hotstuff/keys/r1.key
../cmd/hotstuffgenkey/hotstuffgenkey -filepath /opt/hotstuff/keys/r2.key
../cmd/hotstuffgenkey/hotstuffgenkey -filepath /opt/hotstuff/keys/r3.key
../cmd/hotstuffgenkey/hotstuffgenkey -filepath /opt/hotstuff/keys/r4.key