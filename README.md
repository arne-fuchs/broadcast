# broadcast
A plugin for IOTA's GoShimmer Node to broadcast every message on the message layer an write it to active tcp connections over port 5050

## Installation
Move the project's folder into your goshimmer/plugins/ folder.

In goshimmer/plugins/research.go add the following line:
```
broadcast.Plugin(),
```
in the node.Plugins(...) list and add this in the import statement:
```
"github.com/iotaledger/goshimmer/plugins/broadcast"
```
You may need to recompile the goshimmer software.

In the config.json you need to add "broadcast" to the "node" sections as followed:
```
"node": {
    "disablePlugins": [],
    "enablePlugins": ["broadcast"]
},
```
You also want to configure the plugin in the same file. You need to paste this config between the other plugins:
```
"broadcast": {
    "bindAddress": "127.0.0.1:5050"
    },
```
Please notice, that you can only connect locally with this ip. If you want it to be reachable from the outside you need to use 0.0.0.0.

If you get a compile error while buidling your goshimmer software, it is mostly because of missing libraries.
If so, do the following:
Go to
```
~/go/src/github.com/iotaledger/goshimmer/plugins
```
and paste in the broadcast folder so go can access the missing files.

## Usage
Just connect to the plugin's port 5050 and you get the messages in real time as long as you are connected.
A maximum of 256 Connections are possible before it throws errors.

## Donations
If you want to keep me motivated to do more open source stuff you can donate me some IOTA's. Even very small amounts makes me happy:

```
iota1qqvrqjfscx5ax7vnt8mmtmzj30af3xf7zfm8t7lnaxyrt73awgqckz02upv
```
<p align="center">
  <img src="https://paesserver.de/img/qr.png?raw=true" width="250" title="logo">
</p>

