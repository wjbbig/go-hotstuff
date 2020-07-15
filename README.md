# go-hotstuff

`go-hotstuff` is a simple implementation of hotstuff consensus protocol.

## Run the example

First, run `scripts/generate_keys.sh` to generate threshold keys, then run `scripts/run_server.sh` to start four hotstuff servers. There are three kinds of networks you can choose to use, they are `basic`, `chained` and `event-driven`, the default network type is `event-driven`. There is a simple client where you can find in `cmd/hotstuffclient`, or you can run `scripts/run_client.sh` to start it.

## TODO

1. Fix bugs, especially chained hotstuff.
2. Complete function of syncing blocks.

## Reference

M. Yin, D. Malkhi, M. K. Reiter, G. Golan Gueta, and I. Abraham, “HotStuff: BFT Consensus in the Lens of Blockchain,” Mar 2018.