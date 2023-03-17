package chain

/*
# ERC20
Triggers when tokens are transferred, including zero value transfers. A token contract which creates new tokens SHOULD trigger a Transfer event with the _from address set to 0x0 when tokens are created.

event Transfer(address indexed _from, address indexed _to, uint256 _value)
topic[0] "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
Total topics: 3 (each indexed field is a topic + topic[0] is the event signature, which doesn't including the indexed keyword when calculating keccak)
Data field holds: uint256 value

# ERC721
This emits when ownership of any NFT changes by any mechanism. This event emits when NFTs are created (`from` == 0) and destroyed (`to` == 0). Exception: during contract creation, any number of NFTs may be created and assigned without emitting Transfer.

event Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId);
topic[0] is the same as ERC20, "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
Total topics: 4
Data field holds: nothing, it is empty

# ERC1155
Either `TransferSingle` or `TransferBatch` MUST emit when tokens are transferred, including zero value transfers as well as minting or burning.
The `_operator` argument MUST be the address of an account/contract that is approved to make the transfer (SHOULD be msg.sender).
The `_from` argument MUST be the address of the holder whose balance is decreased.
The `_to` argument MUST be the address of the recipient whose balance is increased.
The `_id` argument MUST be the token type being transferred.
The `_value` argument MUST be the number of tokens the holder balance is decreased by and match what the recipient balance is increased by.
When minting/creating tokens, the `_from` argument MUST be set to `0x0` (i.e. zero address).
When burning/destroying tokens, the `_to` argument MUST be set to `0x0` (i.e. zero address).

## Single
event TransferSingle(address indexed _operator, address indexed _from, address indexed _to, uint256 _id, uint256 _value);
topic[0] for single is "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
Total topics: 4
Data field holds:  uint256 _id, uint256 _value

## Batch
event TransferBatch(address indexed _operator, address indexed _from, address indexed _to, uint256[] _ids, uint256[] _values);
topic[0] for batch is  "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb"
Total topics: 4
Data field holds: uint256[] _ids, uint256[] _values

*/

import (
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

var zeroAddress = common.Address{}

// transferEventKeccakTokens is the keccak256 hash that corresponds to the Transfer event signature
// of an ERC20 or ERC721 token
const transferEventKeccakTokens string = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// transferEventKeccakHybridSingle is the keccak256 hash that corresponds to the Transfer event signature
// of a single transfer of an ERC1155 token
const transferEventKeccakHybridSingle string = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"

// transferEventKeccakHybridBatch is the keccak256 hash that corresponds to the Transfer event signature
// of a batch transfer of an ERC1155 token
const transferEventKeccakHybridBatch string = "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb"

var transferSingleAbi abi.ABI
var transferBatchAbi abi.ABI

func init() {
	transferSingleAbiJSON := "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"}]"
	transferBatchAbiJSON := "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"}]"

	var err error
	transferSingleAbi, err = abi.JSON(strings.NewReader(transferSingleAbiJSON))
	if err != nil {
		logr.Error.Panicf("failed to parse transferSingleAbi: %s", err)
	}
	transferBatchAbi, err = abi.JSON(strings.NewReader(transferBatchAbiJSON))
	if err != nil {
		logr.Error.Panicf("failed to parse transferBatchAbi: %s", err)
	}
}

// DecodeDataForErc1155Single decodes the log.Data field from a ERC1155 "single transfer" event log
// and returns zeroes if any troubles
func DecodeDataForErc1155Single(logData []byte) (id *big.Int, value *big.Int) {

	var transferSingleEvent struct {
		Operator common.Address
		From     common.Address
		To       common.Address
		Id       *big.Int
		Value    *big.Int
	}

	err := transferSingleAbi.UnpackIntoInterface(&transferSingleEvent, "TransferSingle", logData)
	if err != nil {
		logr.Warning.Printf("Failed to unpack transferSingleEvent: %s\n", err)
		return big.NewInt(0), big.NewInt(0)
	}
	return transferSingleEvent.Id, transferSingleEvent.Value
}

// DecodeDataForErc1155Batch decodes the log.Data field from a ERC1155 "batch transfer" event log and returns
// empty slices if any troubles. []ids and []values are guaranteed to have same number of entries, which could
// be zero.
func DecodeDataForErc1155Batch(logData []byte) (ids []*big.Int, values []*big.Int) {

	var transferBatchEvent struct {
		Operator common.Address
		From     common.Address
		To       common.Address
		Ids      []*big.Int
		Values   []*big.Int
	}

	err := transferBatchAbi.UnpackIntoInterface(&transferBatchEvent, "TransferBatch", logData)
	if err != nil {
		logr.Warning.Printf("Failed to unpack transferBatchEvent: %s\n", err)
		return make([]*big.Int, 0), make([]*big.Int, 0)
	}

	for _, id := range transferBatchEvent.Ids {
		ids = append(ids, big.NewInt(0).Set(id))
	}
	for _, value := range transferBatchEvent.Values {
		values = append(values, big.NewInt(0).Set(value))
	}
	if len(ids) != len(values) {
		logr.Warning.Println("Unpacking transferBatchEvent had different numbers of ids and values")
		return make([]*big.Int, 0), make([]*big.Int, 0)
	}

	return ids, values
}
