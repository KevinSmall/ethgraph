# EthGraph
EthGraph is a tool to generate GraphML for token movements for ERC20, ERC721 and ERC1155 (or their equivalents) on any EVM-compatible blockchain. [GraphML](https://en.wikipedia.org/wiki/GraphML) is an XML file format for graph data which is commonly used by graphing tools. We can then use those graphing tools to analyse token movements.

## How can this be used?
The token movements form a `graph` made up of `nodes` (addresses and movements) and `edges` (links between addresses and movement). A `movement` is a transfer of a single token type or single NFT between one address and another. ERC1155 transactions, which can move many NFTs in a single transaction, are decomposed into their individual movements.

When a GraphML file of token movements is opened in a graphing tool like [Gephi](https://gephi.org/) (free, open-source and cross-platform) we can easily visualise the token movements:

![Movements as graph](./docs/movements_as_graph.png "Token Movements as Graph (Gephi screenshot)")

In the above screenshot, taken from Gephi, the large green circles are addresses (full address hashes are available in Gephi and can be copied to the clipboard, but in visuals they clutter up the display). The smaller circles are movements of a particular token or NFT, color coded by token type. USDC is light green, USDT is orange. 

The above sample is part of a larger dataset of some 60k nodes and 100k edges. Using the exact same dataset in Gephi, we can zoom out, switch off labels and adjust the coloring. When zooming out, some common patterns start to appear:

![Far graph](./docs/movements_as_graph_far.png "Token Movements as Graph from Afar (Gephi screenshot")

In the above screenshot, the dense "starburst" area at the bottom left is the zero address. The zero address is involved in many minting and burning movements. Towards the centre of the image, other dense areas form that are the crypto exchange addresses. To the right of the image is the "dust" of addresses that have transacted little and so are not very connected.

This is kind-of pretty, but it's hard to make any sense of very dense graphs. In graph analysis very dense graphs are known as "hairballs". Graphing tools also allow statistical analysis, and Gephi offers many standard measures of connectedness and identification of clusters. Gephi also allows animated views showing how graphs grow over time. More samples can be found in the project [Wiki](https://github.com/KevinSmall/ethgraph/wiki).

This dataset was formed from about 200 blocks of Ethereum mainnet, about 40 minutes worth of movements. Since `ethgraph` is compatible with any EVM-based chain that uses similar ERC token standards, we can easily collect data from other chains. Avalanche has a similar appearance, with an exchange and the zero address dominating:

![Avalanche graph](./docs/movements_avalanche.png "Avalanche Token Movements as Graph (Gephi screenshot")

## How to install

### 1. Install Go
Official instructions [here](https://go.dev/doc/install), or for Ubuntu:
```
$ sudo apt update && sudo apt upgrade
$ sudo apt install golang-go
```
To check it installed ok:
```
$ go version
go version go1.18.1 linux/amd64
```
### 2. Install EthGraph
Download and build:
```
$ git clone https://github.com/KevinSmall/ethgraph.git 
$ cd ethgraph
$ go build .
```
To check it built ok, execute this to see version information:
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
$ ./ethgraph byblock  "https://mainnet.infura.io/v3/<your API key>>"  -f 16_835_977 -t 16_835_986
```
Or for Avalanche on Infura:
```
$ ./ethgraph byblock "https://avalanche-mainnet.infura.io/v3/<your API key>" -f 27_486_035 -t 27_486_094
```

The file created is called `<chainname>.graphml`. This can then be opened in Gephi or other graph tools, see [Wiki](https://github.com/KevinSmall/ethgraph/wiki) for more detailed usage.