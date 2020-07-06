#!/usr/bin/env bash

go build -o ../server/hotstuff ../server/server.go

gnome-terminal -x bash -c "../server/hotstuff -id 1"
gnome-terminal -x bash -c "../server/hotstuff -id 2"
gnome-terminal -x bash -c "../server/hotstuff -id 3"
gnome-terminal -x bash -c "../server/hotstuff -id 4"