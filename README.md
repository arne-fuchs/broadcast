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

## Usage
Just connect to the plugin's port 5050 and you get the messages in real time as long as you are connected.
A maximum of 256 Connections are possible before it throws errors.

## Donations
If you want to keep me motivated to do more open source stuff you can donate me some IOTA's. Even very small amounts makes me happy:

```
iota1qqvrqjfscx5ax7vnt8mmtmzj30af3xf7zfm8t7lnaxyrt73awgqckz02upv
```
