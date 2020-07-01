#!/usr/bin/env bash

go build -o ../cmd/hotstuffgenkey/hotstuffgenkey ../cmd/hotstuffgenkey/main.go

../cmd/hotstuffgenkey/hotstuffgenkey -p /opt/hotstuff/keys -k 3 -l 4