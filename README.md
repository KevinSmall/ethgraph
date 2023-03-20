# EthGraph
EthGraph is a tool written in Go to generate GraphML for token movements for ERC20, ERC721 and ERC1155 (or their equivalents) on any EVM-compatible blockchain. [GraphML](https://en.wikipedia.org/wiki/GraphML) is an XML file format for graph data which is commonly used by graphing tools. We can then use those graphing tools to analyse token movements.

## How can this be used?
The token movements form a `graph` made up of `nodes` (addresses and movements) and `edges` (links between addresses and movement). A `movement` is a transfer of a single token type or single NFT between one address and another. ERC1155 transactions, which can move many NFTs in a single transaction, are decomposed into their individual movements.

When a GraphML file of token movements is opened in a graphing tool like [Gephi](https://gephi.org/) (free, open-source and cross-platform) we can easily visualise the token movements:

![Movements as graph](./docs/movements_as_graph.png "Token Movements as Graph (Gephi screenshot)")

In the above screenshot, taken from Gephi, the large green circles are addresses. Full address hashes are available in Gephi and can be copied to the clipboard, but in visuals they clutter up the display. The smaller circles are movements of a particular token or NFT, color coded by token type. USDC is light green, USDT is orange. 

The above sample is part of a larger dataset. Using the _exact same dataset_ in Gephi, we can zoom out, switch off labels and adjust the coloring. When zooming out, some common patterns start to appear:

![Far graph](./docs/movements_as_graph_far.png "Token Movements as Graph from afar (Gephi screenshot)")

In the above screenshot, the dense "starburst" area at the bottom left is the zero address. The zero address is involved in many movements, very often minting happens from the zero address. Towards the centre of the image, other dense areas form that are the crypto exchange addresses. To the right of the image is the "dust" of addresses that have transacted little (in the block ranges used for data selection) and so are not very connected.

This is all kind-of pretty, but it's hard to make any sense of very dense graphs. In graph analysis very dense graphs are known as "hairballs". Graphing tools also allow statistical analysis, and Gephi offers many standard measures of connectedness and identification of clusters. Gephi also allows animated views showing how graphs grow over time. More samples can be found in the project [Wiki](https://github.com/KevinSmall/ethgraph/wiki).

Since `ethgraph` is compatible with any EVM-based chain that uses similar ERC token standards, we can easily collect data from other chains. Avalanche has a similar appearance, although many fewer transactions. Here an exchange and the zero address are dominating:

![Avalanche graph](./docs/movements_avalanche.png "Avalanche Token Movements as Graph (Gephi screenshot)")

## Performance and Limitations
The above Ethereum dataset was formed from about 200 blocks (~40 minutes) worth of transactions from Ethereum mainnet. This produces some 60k nodes and 90k edges. As a rule of thumb, graphs of 50k to 100k nodes become cumbersome to use with Gephi or similar tools.

`ethgraph` is designed to perform well. Processing 200 blocks of mainnet, including master data retrieval for thousands of tokens, takes ~7 seconds on a reasonable laptop. This produces a file that starts to reach the limits of Gephi. Smaller extracts are much easier to manage. When experimenting, start with just a few blocks and work up.

## How to install

### 1. Install Go
Check if Go is already installed:
```
$ go version
go version go1.18.1 linux/amd64
```
If Go is not installed yet, see the [official instructions](https://go.dev/doc/install), or for Ubuntu:
```
$ sudo apt update && sudo apt upgrade
$ sudo apt install golang-go
```
To check that Go installed ok:
```
$ go version
go version go1.18.1 linux/amd64
```
Builds have been tested with `go1.18.1` and `go1.19.5`.

### 2. Install EthGraph
Download and build:
```
$ git clone https://github.com/KevinSmall/ethgraph.git 
$ cd ethgraph
$ go build .
```
To ensure that `ethgraph` built ok, check version information:
```
$ ./ethgraph -v
ethgraph version 1.0.0
```

## How to use
Sample usage, here selecting block numbers from Ethereum mainnet:
```
$ ./ethgraph byblock "https://<RPC endpoint>"  -f 16_835_977 -t 16_835_978
```
If you were using Infura it might be:
```
$ ./ethgraph byblock  "https://mainnet.infura.io/v3/<your API key>"  -f 16_835_977 -t 16_835_986
```
Or for Avalanche on Infura:
```
$ ./ethgraph byblock "https://avalanche-mainnet.infura.io/v3/<your API key>" -f 27_486_035 -t 27_486_094
```

The file created is called `<chainname>.graphml`. It will overwrite any existing file with the same name. This file can then be opened in Gephi or other graph tools, see [Wiki](https://github.com/KevinSmall/ethgraph/wiki) for more detailed usage.