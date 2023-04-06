// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package wrapped_token_pool

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type IPoolRampUpdate struct {
	Ramp    common.Address
	Allowed bool
}

var WrappedTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"currentLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"ExceedsTokenLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NullAddressNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermissionsError\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"OffRampAllowanceSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"OnRampAllowanceSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"internalType\":\"structIPool.RampUpdate[]\",\"name\":\"onRamps\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"ramp\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"internalType\":\"structIPool.RampUpdate[]\",\"name\":\"offRamps\",\"type\":\"tuple[]\"}],\"name\":\"applyRampUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"onRamp\",\"type\":\"address\"}],\"name\":\"isOnRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lockOrBurn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"releaseOrMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b506040516200220f3803806200220f8339810160408190526200003491620002b2565b8282828282303380600081620000915760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c457620000c4816200013a565b50506001805460ff60a01b19169055506001600160a01b038116620000fc57604051634655efd160e11b815260040160405180910390fd5b6001600160a01b03166080526009620001168382620003c6565b50600a620001258282620003c6565b50505060ff1660a05250620004929350505050565b336001600160a01b03821603620001945760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000088565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200020d57600080fd5b81516001600160401b03808211156200022a576200022a620001e5565b604051601f8301601f19908116603f01168101908282118183101715620002555762000255620001e5565b816040528381526020925086838588010111156200027257600080fd5b600091505b8382101562000296578582018301518183018401529082019062000277565b83821115620002a85760008385830101525b9695505050505050565b600080600060608486031215620002c857600080fd5b83516001600160401b0380821115620002e057600080fd5b620002ee87838801620001fb565b945060208601519150808211156200030557600080fd5b506200031486828701620001fb565b925050604084015160ff811681146200032c57600080fd5b809150509250925092565b600181811c908216806200034c57607f821691505b6020821081036200036d57634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620003c157600081815260208120601f850160051c810160208610156200039c5750805b601f850160051c820191505b81811015620003bd57828155600101620003a8565b5050505b505050565b81516001600160401b03811115620003e257620003e2620001e5565b620003fa81620003f3845462000337565b8462000373565b602080601f831160018114620004325760008415620004195750858301515b600019600386901b1c1916600185901b178555620003bd565b600085815260208120601f198616915b82811015620004635788860151825594840194600190910190840162000442565b5085821015620004825787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60805160a051611d57620004b86000396000610259015260006101fa0152611d576000f3fe608060405234801561001057600080fd5b506004361061018d5760003560e01c806370a08231116100e3578063a9059cbb1161008c578063e2e59b3e11610066578063e2e59b3e146103c1578063ea6192a2146103d4578063f2fde38b146103e757600080fd5b8063a9059cbb14610355578063af51911214610368578063dd62ed3e1461037b57600080fd5b80638da5cb5b116100bd5780638da5cb5b1461031c57806395d89b411461033a578063a457c2d71461034257600080fd5b806370a08231146102d657806379ba50971461030c5780638456cb591461031457600080fd5b806323b872dd116101455780633f4ba83a1161011f5780633f4ba83a146102965780635c975abb146102a05780636f32b872146102c357600080fd5b806323b872dd1461023f578063313ce56714610252578063395093511461028357600080fd5b806318160ddd1161017657806318160ddd146101d35780631d7a74a0146101e557806321df0da7146101f857600080fd5b806306fdde0314610192578063095ea7b3146101b0575b600080fd5b61019a6103fa565b6040516101a791906118cb565b60405180910390f35b6101c36101be366004611967565b61048c565b60405190151581526020016101a7565b6008545b6040519081526020016101a7565b6101c36101f3366004611991565b6104a3565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101a7565b6101c361024d3660046119ac565b6104b0565b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101a7565b6101c3610291366004611967565b61059b565b61029e6105e4565b005b60015474010000000000000000000000000000000000000000900460ff166101c3565b6101c36102d1366004611991565b6105f6565b6101d76102e4366004611991565b73ffffffffffffffffffffffffffffffffffffffff1660009081526006602052604090205490565b61029e610603565b61029e610700565b60005473ffffffffffffffffffffffffffffffffffffffff1661021a565b61019a610710565b6101c3610350366004611967565b61071f565b6101c3610363366004611967565b6107f7565b61029e610376366004611b49565b610804565b6101d7610389366004611bad565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260076020908152604080832093909416825291909152205490565b61029e6103cf366004611be0565b610a10565b61029e6103e2366004611967565b610b17565b61029e6103f5366004611991565b610c37565b60606009805461040990611c03565b80601f016020809104026020016040519081016040528092919081815260200182805461043590611c03565b80156104825780601f1061045757610100808354040283529160200191610482565b820191906000526020600020905b81548152906001019060200180831161046557829003601f168201915b5050505050905090565b6000610499338484610c4b565b5060015b92915050565b600061049d600483610dfe565b60006104bd848484610e30565b73ffffffffffffffffffffffffffffffffffffffff8416600090815260076020908152604080832033845290915290205482811015610583576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206160448201527f6c6c6f77616e636500000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b6105908533858403610c4b565b506001949350505050565b33600081815260076020908152604080832073ffffffffffffffffffffffffffffffffffffffff8716845290915281205490916104999185906105df908690611c85565b610c4b565b6105ec6110e4565b6105f4611165565b565b600061049d600283610dfe565b60015473ffffffffffffffffffffffffffffffffffffffff163314610684576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161057a565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6107086110e4565b6105f461125e565b6060600a805461040990611c03565b33600090815260076020908152604080832073ffffffffffffffffffffffffffffffffffffffff86168452909152812054828110156107e0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760448201527f207a65726f000000000000000000000000000000000000000000000000000000606482015260840161057a565b6107ed3385858403610c4b565b5060019392505050565b6000610499338484610e30565b61080c6110e4565b60005b825181101561090b57600083828151811061082c5761082c611c9d565b6020026020010151905080602001511561085457805161084e9060029061134a565b50610864565b80516108629060029061136c565b505b7fbceff8f229c6dfcbf8bdcfb18726b84b0fd249b4803deb3948ff34d90401366284838151811061089757610897611c9d565b6020026020010151600001518584815181106108b5576108b5611c9d565b6020026020010151602001516040516108f292919073ffffffffffffffffffffffffffffffffffffffff9290921682521515602082015260400190565b60405180910390a15061090481611ccc565b905061080f565b5060005b8151811015610a0b57600082828151811061092c5761092c611c9d565b6020026020010151905080602001511561095457805161094e9060049061134a565b50610964565b80516109629060049061136c565b505b7fd8c3333ded377884ced3869cd0bcb9be54ea664076df1f5d39c468912031364883838151811061099757610997611c9d565b6020026020010151600001518484815181106109b5576109b5611c9d565b6020026020010151602001516040516109f292919073ffffffffffffffffffffffffffffffffffffffff9290921682521515602082015260400190565b60405180910390a150610a0481611ccc565b905061090f565b505050565b60015474010000000000000000000000000000000000000000900460ff1615610a95576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015260640161057a565b610a9e336105f6565b610ad4576040517f5307f5ab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610ade308361138e565b60405182815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a25050565b60015474010000000000000000000000000000000000000000900460ff1615610b9c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015260640161057a565b610ba5336104a3565b610bdb576040517f5307f5ab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610be5828261157b565b60405181815273ffffffffffffffffffffffffffffffffffffffff83169033907f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0906020015b60405180910390a35050565b610c3f6110e4565b610c4881611694565b50565b73ffffffffffffffffffffffffffffffffffffffff8316610ced576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f7265737300000000000000000000000000000000000000000000000000000000606482015260840161057a565b73ffffffffffffffffffffffffffffffffffffffff8216610d90576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f7373000000000000000000000000000000000000000000000000000000000000606482015260840161057a565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526007602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415155b9392505050565b73ffffffffffffffffffffffffffffffffffffffff8316610ed3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f20616460448201527f6472657373000000000000000000000000000000000000000000000000000000606482015260840161057a565b73ffffffffffffffffffffffffffffffffffffffff8216610f76576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201527f6573730000000000000000000000000000000000000000000000000000000000606482015260840161057a565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600660205260409020548181101561102c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206260448201527f616c616e63650000000000000000000000000000000000000000000000000000606482015260840161057a565b73ffffffffffffffffffffffffffffffffffffffff808516600090815260066020526040808220858503905591851681529081208054849290611070908490611c85565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516110d691815260200190565b60405180910390a350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146105f4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161057a565b60015474010000000000000000000000000000000000000000900460ff166111e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f5061757361626c653a206e6f7420706175736564000000000000000000000000604482015260640161057a565b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390a1565b60015474010000000000000000000000000000000000000000900460ff16156112e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f5061757361626c653a2070617573656400000000000000000000000000000000604482015260640161057a565b600180547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586112343390565b6000610e298373ffffffffffffffffffffffffffffffffffffffff8416611789565b6000610e298373ffffffffffffffffffffffffffffffffffffffff84166117d8565b73ffffffffffffffffffffffffffffffffffffffff8216611431576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360448201527f7300000000000000000000000000000000000000000000000000000000000000606482015260840161057a565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260066020526040902054818110156114e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60448201527f6365000000000000000000000000000000000000000000000000000000000000606482015260840161057a565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600660205260408120838303905560088054849290611523908490611d04565b909155505060405182815260009073ffffffffffffffffffffffffffffffffffffffff8516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b73ffffffffffffffffffffffffffffffffffffffff82166115f8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640161057a565b806008600082825461160a9190611c85565b909155505073ffffffffffffffffffffffffffffffffffffffff821660009081526006602052604081208054839290611644908490611c85565b909155505060405181815273ffffffffffffffffffffffffffffffffffffffff8316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610c2b565b3373ffffffffffffffffffffffffffffffffffffffff821603611713576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161057a565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60008181526001830160205260408120546117d05750815460018181018455600084815260208082209093018490558454848252828601909352604090209190915561049d565b50600061049d565b600081815260018301602052604081205480156118c15760006117fc600183611d04565b855490915060009061181090600190611d04565b905081811461187557600086600001828154811061183057611830611c9d565b906000526020600020015490508087600001848154811061185357611853611c9d565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061188657611886611d1b565b60019003818190600052602060002001600090559055856001016000868152602001908152602001600020600090556001935050505061049d565b600091505061049d565b600060208083528351808285015260005b818110156118f8578581018301518582016040015282016118dc565b8181111561190a576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461196257600080fd5b919050565b6000806040838503121561197a57600080fd5b6119838361193e565b946020939093013593505050565b6000602082840312156119a357600080fd5b610e298261193e565b6000806000606084860312156119c157600080fd5b6119ca8461193e565b92506119d86020850161193e565b9150604084013590509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715611a3a57611a3a6119e8565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611a8757611a876119e8565b604052919050565b600082601f830112611aa057600080fd5b8135602067ffffffffffffffff821115611abc57611abc6119e8565b611aca818360051b01611a40565b82815260069290921b84018101918181019086841115611ae957600080fd5b8286015b84811015611b3e5760408189031215611b065760008081fd5b611b0e611a17565b611b178261193e565b8152848201358015158114611b2c5760008081fd5b81860152835291830191604001611aed565b509695505050505050565b60008060408385031215611b5c57600080fd5b823567ffffffffffffffff80821115611b7457600080fd5b611b8086838701611a8f565b93506020850135915080821115611b9657600080fd5b50611ba385828601611a8f565b9150509250929050565b60008060408385031215611bc057600080fd5b611bc98361193e565b9150611bd76020840161193e565b90509250929050565b60008060408385031215611bf357600080fd5b82359150611bd76020840161193e565b600181811c90821680611c1757607f821691505b602082108103611c50577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008219821115611c9857611c98611c56565b500190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611cfd57611cfd611c56565b5060010190565b600082821015611d1657611d16611c56565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c634300080f000a",
}

var WrappedTokenPoolABI = WrappedTokenPoolMetaData.ABI

var WrappedTokenPoolBin = WrappedTokenPoolMetaData.Bin

func DeployWrappedTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals uint8) (common.Address, *types.Transaction, *WrappedTokenPool, error) {
	parsed, err := WrappedTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WrappedTokenPoolBin), backend, name, symbol, decimals)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WrappedTokenPool{WrappedTokenPoolCaller: WrappedTokenPoolCaller{contract: contract}, WrappedTokenPoolTransactor: WrappedTokenPoolTransactor{contract: contract}, WrappedTokenPoolFilterer: WrappedTokenPoolFilterer{contract: contract}}, nil
}

type WrappedTokenPool struct {
	address common.Address
	abi     abi.ABI
	WrappedTokenPoolCaller
	WrappedTokenPoolTransactor
	WrappedTokenPoolFilterer
}

type WrappedTokenPoolCaller struct {
	contract *bind.BoundContract
}

type WrappedTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type WrappedTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type WrappedTokenPoolSession struct {
	Contract     *WrappedTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type WrappedTokenPoolCallerSession struct {
	Contract *WrappedTokenPoolCaller
	CallOpts bind.CallOpts
}

type WrappedTokenPoolTransactorSession struct {
	Contract     *WrappedTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type WrappedTokenPoolRaw struct {
	Contract *WrappedTokenPool
}

type WrappedTokenPoolCallerRaw struct {
	Contract *WrappedTokenPoolCaller
}

type WrappedTokenPoolTransactorRaw struct {
	Contract *WrappedTokenPoolTransactor
}

func NewWrappedTokenPool(address common.Address, backend bind.ContractBackend) (*WrappedTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(WrappedTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindWrappedTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPool{address: address, abi: abi, WrappedTokenPoolCaller: WrappedTokenPoolCaller{contract: contract}, WrappedTokenPoolTransactor: WrappedTokenPoolTransactor{contract: contract}, WrappedTokenPoolFilterer: WrappedTokenPoolFilterer{contract: contract}}, nil
}

func NewWrappedTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*WrappedTokenPoolCaller, error) {
	contract, err := bindWrappedTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolCaller{contract: contract}, nil
}

func NewWrappedTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*WrappedTokenPoolTransactor, error) {
	contract, err := bindWrappedTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolTransactor{contract: contract}, nil
}

func NewWrappedTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*WrappedTokenPoolFilterer, error) {
	contract, err := bindWrappedTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolFilterer{contract: contract}, nil
}

func bindWrappedTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WrappedTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_WrappedTokenPool *WrappedTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WrappedTokenPool.Contract.WrappedTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_WrappedTokenPool *WrappedTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.WrappedTokenPoolTransactor.contract.Transfer(opts)
}

func (_WrappedTokenPool *WrappedTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.WrappedTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WrappedTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.contract.Transfer(opts)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WrappedTokenPool.Contract.Allowance(&_WrappedTokenPool.CallOpts, owner, spender)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WrappedTokenPool.Contract.Allowance(&_WrappedTokenPool.CallOpts, owner, spender)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WrappedTokenPool.Contract.BalanceOf(&_WrappedTokenPool.CallOpts, account)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WrappedTokenPool.Contract.BalanceOf(&_WrappedTokenPool.CallOpts, account)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) Decimals() (uint8, error) {
	return _WrappedTokenPool.Contract.Decimals(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) Decimals() (uint8, error) {
	return _WrappedTokenPool.Contract.Decimals(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) GetToken() (common.Address, error) {
	return _WrappedTokenPool.Contract.GetToken(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _WrappedTokenPool.Contract.GetToken(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "isOffRamp", offRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) IsOffRamp(offRamp common.Address) (bool, error) {
	return _WrappedTokenPool.Contract.IsOffRamp(&_WrappedTokenPool.CallOpts, offRamp)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) IsOffRamp(offRamp common.Address) (bool, error) {
	return _WrappedTokenPool.Contract.IsOffRamp(&_WrappedTokenPool.CallOpts, offRamp)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) IsOnRamp(opts *bind.CallOpts, onRamp common.Address) (bool, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "isOnRamp", onRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) IsOnRamp(onRamp common.Address) (bool, error) {
	return _WrappedTokenPool.Contract.IsOnRamp(&_WrappedTokenPool.CallOpts, onRamp)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) IsOnRamp(onRamp common.Address) (bool, error) {
	return _WrappedTokenPool.Contract.IsOnRamp(&_WrappedTokenPool.CallOpts, onRamp)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) Name() (string, error) {
	return _WrappedTokenPool.Contract.Name(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) Name() (string, error) {
	return _WrappedTokenPool.Contract.Name(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) Owner() (common.Address, error) {
	return _WrappedTokenPool.Contract.Owner(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) Owner() (common.Address, error) {
	return _WrappedTokenPool.Contract.Owner(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) Paused() (bool, error) {
	return _WrappedTokenPool.Contract.Paused(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) Paused() (bool, error) {
	return _WrappedTokenPool.Contract.Paused(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) Symbol() (string, error) {
	return _WrappedTokenPool.Contract.Symbol(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) Symbol() (string, error) {
	return _WrappedTokenPool.Contract.Symbol(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WrappedTokenPool.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_WrappedTokenPool *WrappedTokenPoolSession) TotalSupply() (*big.Int, error) {
	return _WrappedTokenPool.Contract.TotalSupply(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolCallerSession) TotalSupply() (*big.Int, error) {
	return _WrappedTokenPool.Contract.TotalSupply(&_WrappedTokenPool.CallOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_WrappedTokenPool *WrappedTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.AcceptOwnership(&_WrappedTokenPool.TransactOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.AcceptOwnership(&_WrappedTokenPool.TransactOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) ApplyRampUpdates(opts *bind.TransactOpts, onRamps []IPoolRampUpdate, offRamps []IPoolRampUpdate) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "applyRampUpdates", onRamps, offRamps)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) ApplyRampUpdates(onRamps []IPoolRampUpdate, offRamps []IPoolRampUpdate) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.ApplyRampUpdates(&_WrappedTokenPool.TransactOpts, onRamps, offRamps)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) ApplyRampUpdates(onRamps []IPoolRampUpdate, offRamps []IPoolRampUpdate) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.ApplyRampUpdates(&_WrappedTokenPool.TransactOpts, onRamps, offRamps)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "approve", spender, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.Approve(&_WrappedTokenPool.TransactOpts, spender, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.Approve(&_WrappedTokenPool.TransactOpts, spender, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.DecreaseAllowance(&_WrappedTokenPool.TransactOpts, spender, subtractedValue)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.DecreaseAllowance(&_WrappedTokenPool.TransactOpts, spender, subtractedValue)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.IncreaseAllowance(&_WrappedTokenPool.TransactOpts, spender, addedValue)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.IncreaseAllowance(&_WrappedTokenPool.TransactOpts, spender, addedValue)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, amount *big.Int, arg1 common.Address) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "lockOrBurn", amount, arg1)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) LockOrBurn(amount *big.Int, arg1 common.Address) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.LockOrBurn(&_WrappedTokenPool.TransactOpts, amount, arg1)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) LockOrBurn(amount *big.Int, arg1 common.Address) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.LockOrBurn(&_WrappedTokenPool.TransactOpts, amount, arg1)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "pause")
}

func (_WrappedTokenPool *WrappedTokenPoolSession) Pause() (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.Pause(&_WrappedTokenPool.TransactOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) Pause() (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.Pause(&_WrappedTokenPool.TransactOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "releaseOrMint", recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) ReleaseOrMint(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.ReleaseOrMint(&_WrappedTokenPool.TransactOpts, recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) ReleaseOrMint(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.ReleaseOrMint(&_WrappedTokenPool.TransactOpts, recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "transfer", recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.Transfer(&_WrappedTokenPool.TransactOpts, recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.Transfer(&_WrappedTokenPool.TransactOpts, recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.TransferFrom(&_WrappedTokenPool.TransactOpts, sender, recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.TransferFrom(&_WrappedTokenPool.TransactOpts, sender, recipient, amount)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_WrappedTokenPool *WrappedTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.TransferOwnership(&_WrappedTokenPool.TransactOpts, to)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.TransferOwnership(&_WrappedTokenPool.TransactOpts, to)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappedTokenPool.contract.Transact(opts, "unpause")
}

func (_WrappedTokenPool *WrappedTokenPoolSession) Unpause() (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.Unpause(&_WrappedTokenPool.TransactOpts)
}

func (_WrappedTokenPool *WrappedTokenPoolTransactorSession) Unpause() (*types.Transaction, error) {
	return _WrappedTokenPool.Contract.Unpause(&_WrappedTokenPool.TransactOpts)
}

type WrappedTokenPoolApprovalIterator struct {
	Event *WrappedTokenPoolApproval

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolApprovalIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolApprovalIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WrappedTokenPoolApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolApprovalIterator{contract: _WrappedTokenPool.contract, event: "Approval", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolApproval)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseApproval(log types.Log) (*WrappedTokenPoolApproval, error) {
	event := new(WrappedTokenPoolApproval)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolBurnedIterator struct {
	Event *WrappedTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolBurned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolBurned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*WrappedTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolBurnedIterator{contract: _WrappedTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolBurned)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseBurned(log types.Log) (*WrappedTokenPoolBurned, error) {
	event := new(WrappedTokenPoolBurned)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolLockedIterator struct {
	Event *WrappedTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolLocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolLocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*WrappedTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolLockedIterator{contract: _WrappedTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolLocked)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseLocked(log types.Log) (*WrappedTokenPoolLocked, error) {
	event := new(WrappedTokenPoolLocked)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolMintedIterator struct {
	Event *WrappedTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*WrappedTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolMintedIterator{contract: _WrappedTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolMinted)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseMinted(log types.Log) (*WrappedTokenPoolMinted, error) {
	event := new(WrappedTokenPoolMinted)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolOffRampAllowanceSetIterator struct {
	Event *WrappedTokenPoolOffRampAllowanceSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolOffRampAllowanceSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolOffRampAllowanceSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolOffRampAllowanceSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolOffRampAllowanceSetIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolOffRampAllowanceSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolOffRampAllowanceSet struct {
	OnRamp  common.Address
	Allowed bool
	Raw     types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterOffRampAllowanceSet(opts *bind.FilterOpts) (*WrappedTokenPoolOffRampAllowanceSetIterator, error) {

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "OffRampAllowanceSet")
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolOffRampAllowanceSetIterator{contract: _WrappedTokenPool.contract, event: "OffRampAllowanceSet", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchOffRampAllowanceSet(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolOffRampAllowanceSet) (event.Subscription, error) {

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "OffRampAllowanceSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolOffRampAllowanceSet)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "OffRampAllowanceSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseOffRampAllowanceSet(log types.Log) (*WrappedTokenPoolOffRampAllowanceSet, error) {
	event := new(WrappedTokenPoolOffRampAllowanceSet)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "OffRampAllowanceSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolOnRampAllowanceSetIterator struct {
	Event *WrappedTokenPoolOnRampAllowanceSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolOnRampAllowanceSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolOnRampAllowanceSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolOnRampAllowanceSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolOnRampAllowanceSetIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolOnRampAllowanceSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolOnRampAllowanceSet struct {
	OnRamp  common.Address
	Allowed bool
	Raw     types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterOnRampAllowanceSet(opts *bind.FilterOpts) (*WrappedTokenPoolOnRampAllowanceSetIterator, error) {

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "OnRampAllowanceSet")
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolOnRampAllowanceSetIterator{contract: _WrappedTokenPool.contract, event: "OnRampAllowanceSet", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchOnRampAllowanceSet(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolOnRampAllowanceSet) (event.Subscription, error) {

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "OnRampAllowanceSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolOnRampAllowanceSet)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "OnRampAllowanceSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseOnRampAllowanceSet(log types.Log) (*WrappedTokenPoolOnRampAllowanceSet, error) {
	event := new(WrappedTokenPoolOnRampAllowanceSet)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "OnRampAllowanceSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolOwnershipTransferRequestedIterator struct {
	Event *WrappedTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolOwnershipTransferRequestedIterator{contract: _WrappedTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolOwnershipTransferRequested)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*WrappedTokenPoolOwnershipTransferRequested, error) {
	event := new(WrappedTokenPoolOwnershipTransferRequested)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolOwnershipTransferredIterator struct {
	Event *WrappedTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolOwnershipTransferredIterator{contract: _WrappedTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolOwnershipTransferred)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*WrappedTokenPoolOwnershipTransferred, error) {
	event := new(WrappedTokenPoolOwnershipTransferred)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolPausedIterator struct {
	Event *WrappedTokenPoolPaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolPausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolPausedIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolPaused struct {
	Account common.Address
	Raw     types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterPaused(opts *bind.FilterOpts) (*WrappedTokenPoolPausedIterator, error) {

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolPausedIterator{contract: _WrappedTokenPool.contract, event: "Paused", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolPaused) (event.Subscription, error) {

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolPaused)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParsePaused(log types.Log) (*WrappedTokenPoolPaused, error) {
	event := new(WrappedTokenPoolPaused)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolReleasedIterator struct {
	Event *WrappedTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolReleased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolReleased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*WrappedTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolReleasedIterator{contract: _WrappedTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolReleased)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseReleased(log types.Log) (*WrappedTokenPoolReleased, error) {
	event := new(WrappedTokenPoolReleased)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolTransferIterator struct {
	Event *WrappedTokenPoolTransfer

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolTransferIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolTransferIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedTokenPoolTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolTransferIterator{contract: _WrappedTokenPool.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolTransfer)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseTransfer(log types.Log) (*WrappedTokenPoolTransfer, error) {
	event := new(WrappedTokenPoolTransfer)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type WrappedTokenPoolUnpausedIterator struct {
	Event *WrappedTokenPoolUnpaused

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *WrappedTokenPoolUnpausedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappedTokenPoolUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(WrappedTokenPoolUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *WrappedTokenPoolUnpausedIterator) Error() error {
	return it.fail
}

func (it *WrappedTokenPoolUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type WrappedTokenPoolUnpaused struct {
	Account common.Address
	Raw     types.Log
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) FilterUnpaused(opts *bind.FilterOpts) (*WrappedTokenPoolUnpausedIterator, error) {

	logs, sub, err := _WrappedTokenPool.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &WrappedTokenPoolUnpausedIterator{contract: _WrappedTokenPool.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolUnpaused) (event.Subscription, error) {

	logs, sub, err := _WrappedTokenPool.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(WrappedTokenPoolUnpaused)
				if err := _WrappedTokenPool.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_WrappedTokenPool *WrappedTokenPoolFilterer) ParseUnpaused(log types.Log) (*WrappedTokenPoolUnpaused, error) {
	event := new(WrappedTokenPoolUnpaused)
	if err := _WrappedTokenPool.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_WrappedTokenPool *WrappedTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _WrappedTokenPool.abi.Events["Approval"].ID:
		return _WrappedTokenPool.ParseApproval(log)
	case _WrappedTokenPool.abi.Events["Burned"].ID:
		return _WrappedTokenPool.ParseBurned(log)
	case _WrappedTokenPool.abi.Events["Locked"].ID:
		return _WrappedTokenPool.ParseLocked(log)
	case _WrappedTokenPool.abi.Events["Minted"].ID:
		return _WrappedTokenPool.ParseMinted(log)
	case _WrappedTokenPool.abi.Events["OffRampAllowanceSet"].ID:
		return _WrappedTokenPool.ParseOffRampAllowanceSet(log)
	case _WrappedTokenPool.abi.Events["OnRampAllowanceSet"].ID:
		return _WrappedTokenPool.ParseOnRampAllowanceSet(log)
	case _WrappedTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _WrappedTokenPool.ParseOwnershipTransferRequested(log)
	case _WrappedTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _WrappedTokenPool.ParseOwnershipTransferred(log)
	case _WrappedTokenPool.abi.Events["Paused"].ID:
		return _WrappedTokenPool.ParsePaused(log)
	case _WrappedTokenPool.abi.Events["Released"].ID:
		return _WrappedTokenPool.ParseReleased(log)
	case _WrappedTokenPool.abi.Events["Transfer"].ID:
		return _WrappedTokenPool.ParseTransfer(log)
	case _WrappedTokenPool.abi.Events["Unpaused"].ID:
		return _WrappedTokenPool.ParseUnpaused(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (WrappedTokenPoolApproval) Topic() common.Hash {
	return common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")
}

func (WrappedTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (WrappedTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (WrappedTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (WrappedTokenPoolOffRampAllowanceSet) Topic() common.Hash {
	return common.HexToHash("0xd8c3333ded377884ced3869cd0bcb9be54ea664076df1f5d39c4689120313648")
}

func (WrappedTokenPoolOnRampAllowanceSet) Topic() common.Hash {
	return common.HexToHash("0xbceff8f229c6dfcbf8bdcfb18726b84b0fd249b4803deb3948ff34d904013662")
}

func (WrappedTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (WrappedTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (WrappedTokenPoolPaused) Topic() common.Hash {
	return common.HexToHash("0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258")
}

func (WrappedTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (WrappedTokenPoolTransfer) Topic() common.Hash {
	return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
}

func (WrappedTokenPoolUnpaused) Topic() common.Hash {
	return common.HexToHash("0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa")
}

func (_WrappedTokenPool *WrappedTokenPool) Address() common.Address {
	return _WrappedTokenPool.address
}

type WrappedTokenPoolInterface interface {
	Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error)

	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)

	Decimals(opts *bind.CallOpts) (uint8, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error)

	IsOnRamp(opts *bind.CallOpts, onRamp common.Address) (bool, error)

	Name(opts *bind.CallOpts) (string, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	Paused(opts *bind.CallOpts) (bool, error)

	Symbol(opts *bind.CallOpts) (string, error)

	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ApplyRampUpdates(opts *bind.TransactOpts, onRamps []IPoolRampUpdate, offRamps []IPoolRampUpdate) (*types.Transaction, error)

	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)

	DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error)

	IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, amount *big.Int, arg1 common.Address) (*types.Transaction, error)

	Pause(opts *bind.TransactOpts) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	Unpause(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WrappedTokenPoolApprovalIterator, error)

	WatchApproval(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolApproval, owner []common.Address, spender []common.Address) (event.Subscription, error)

	ParseApproval(log types.Log) (*WrappedTokenPoolApproval, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*WrappedTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*WrappedTokenPoolBurned, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*WrappedTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*WrappedTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*WrappedTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*WrappedTokenPoolMinted, error)

	FilterOffRampAllowanceSet(opts *bind.FilterOpts) (*WrappedTokenPoolOffRampAllowanceSetIterator, error)

	WatchOffRampAllowanceSet(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolOffRampAllowanceSet) (event.Subscription, error)

	ParseOffRampAllowanceSet(log types.Log) (*WrappedTokenPoolOffRampAllowanceSet, error)

	FilterOnRampAllowanceSet(opts *bind.FilterOpts) (*WrappedTokenPoolOnRampAllowanceSetIterator, error)

	WatchOnRampAllowanceSet(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolOnRampAllowanceSet) (event.Subscription, error)

	ParseOnRampAllowanceSet(log types.Log) (*WrappedTokenPoolOnRampAllowanceSet, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*WrappedTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*WrappedTokenPoolOwnershipTransferred, error)

	FilterPaused(opts *bind.FilterOpts) (*WrappedTokenPoolPausedIterator, error)

	WatchPaused(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolPaused) (event.Subscription, error)

	ParsePaused(log types.Log) (*WrappedTokenPoolPaused, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*WrappedTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*WrappedTokenPoolReleased, error)

	FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WrappedTokenPoolTransferIterator, error)

	WatchTransfer(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolTransfer, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseTransfer(log types.Log) (*WrappedTokenPoolTransfer, error)

	FilterUnpaused(opts *bind.FilterOpts) (*WrappedTokenPoolUnpausedIterator, error)

	WatchUnpaused(opts *bind.WatchOpts, sink chan<- *WrappedTokenPoolUnpaused) (event.Subscription, error)

	ParseUnpaused(log types.Log) (*WrappedTokenPoolUnpaused, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
