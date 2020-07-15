#!/usr/bin/env bash

go build -o ../server/hotstuff ../server/server.go

gnome-terminal -x bash -c "../server/hotstuff -id 1 -type event-driven"
gnome-terminal -x bash -c "../server/hotstuff -id 2 -type event-driven"
gnome-terminal -x bash -c "../server/hotstuff -id 3 -type event-driven"
gnome-terminal -x bash -c "../server/hotstuff -id 4 -type event-driven"