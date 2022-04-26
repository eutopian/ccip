// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package offramp_router

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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
)

type CCIPMessage struct {
	SourceChainId  *big.Int
	SequenceNumber uint64
	Sender         common.Address
	Payload        CCIPMessagePayload
}

type CCIPMessagePayload struct {
	Tokens             []common.Address
	Amounts            []*big.Int
	DestinationChainId *big.Int
	Receiver           common.Address
	Executor           common.Address
	Data               []byte
	Options            []byte
}

var OffRampRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractOffRampInterface[]\",\"name\":\"offRamps\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"contractOffRampInterface\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"AlreadyConfigured\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"MessageFailure\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoOffRampsConfigured\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"contractOffRampInterface\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampNotConfigured\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractOffRampInterface\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractOffRampInterface\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"OffRampRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractOffRampInterface\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"addOffRamp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOffRamps\",\"outputs\":[{\"internalType\":\"contractOffRampInterface[]\",\"name\":\"offRamps\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractOffRampInterface\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"isOffRamp\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractOffRampInterface\",\"name\":\"offRamp\",\"type\":\"address\"}],\"name\":\"removeOffRamp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractCrossChainMessageReceiverInterface\",\"name\":\"receiver\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNumber\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"contractIERC20[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"options\",\"type\":\"bytes\"}],\"internalType\":\"structCCIP.MessagePayload\",\"name\":\"payload\",\"type\":\"tuple\"}],\"internalType\":\"structCCIP.Message\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"routeMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620014ff380380620014ff8339810160408190526200003491620002f9565b33806000816200008b5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000be57620000be816200019a565b50508151620000d69150600390602084019062000245565b5060005b815181101562000192576040518060400160405280826001600160601b0316815260200160011515815250600260008484815181106200011e576200011e620003cb565b6020908102919091018101516001600160a01b031682528181019290925260400160002082518154939092015115156c01000000000000000000000000026001600160681b03199093166001600160601b0390921691909117919091179055806200018981620003e1565b915050620000da565b505062000409565b336001600160a01b03821603620001f45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000082565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b8280548282559060005260206000209081019282156200029d579160200282015b828111156200029d57825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019062000266565b50620002ab929150620002af565b5090565b5b80821115620002ab5760008155600101620002b0565b634e487b7160e01b600052604160045260246000fd5b80516001600160a01b0381168114620002f457600080fd5b919050565b600060208083850312156200030d57600080fd5b82516001600160401b03808211156200032557600080fd5b818501915085601f8301126200033a57600080fd5b8151818111156200034f576200034f620002c6565b8060051b604051601f19603f83011681018181108582111715620003775762000377620002c6565b6040529182528482019250838101850191888311156200039657600080fd5b938501935b82851015620003bf57620003af85620002dc565b845293850193928501926200039b565b98975050505050505050565b634e487b7160e01b600052603260045260246000fd5b6000600182016200040257634e487b7160e01b600052601160045260246000fd5b5060010190565b6110e680620004196000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063991f65431161005b578063991f654314610130578063a40e69c714610143578063adb9f71b14610158578063f2fde38b1461016b57600080fd5b80631d7a74a01461008d5780633fa3c59a146100eb57806379ba5097146101005780638da5cb5b14610108575b600080fd5b6100d661009b366004610b4f565b73ffffffffffffffffffffffffffffffffffffffff166000908152600260205260409020546c01000000000000000000000000900460ff1690565b60405190151581526020015b60405180910390f35b6100fe6100f9366004610b73565b61017e565b005b6100fe6102ef565b60005460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100e2565b6100fe61013e366004610b4f565b6103ec565b61014b61076f565b6040516100e29190610bca565b6100fe610166366004610b4f565b6107de565b6100fe610179366004610b4f565b6109a1565b336000818152600260205260409020546c01000000000000000000000000900460ff166101f4576040517fd49ac04e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b6040517f49bc98f200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8416906349bc98f290610246908590600401610db6565b600060405180830381600087803b15801561026057600080fd5b505af1925050508015610271575060015b6102ea573d80801561029f576040519150601f19603f3d011682016040523d82523d6000602084013e6102a4565b606091505b506102b56040840160208501610f9d565b816040517fd487713c0000000000000000000000000000000000000000000000000000000081526004016101eb929190610fb8565b505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610370576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016101eb565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6103f46109b5565b6003546000819003610432576040517f22babb3200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82166000908152600260209081526040918290208251808401909352546bffffffffffffffffffffffff811683526c01000000000000000000000000900460ff1615159082018190526104de576040517fd49ac04e00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff841660048201526024016101eb565b600060036104ed60018561103d565b815481106104fd576104fd61107b565b60009182526020909120015482516003805473ffffffffffffffffffffffffffffffffffffffff9093169350916bffffffffffffffffffffffff9091169081106105495761054961107b565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16600361057860018661103d565b815481106105885761058861107b565b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600383600001516bffffffffffffffffffffffff16815481106105f6576105f661107b565b600091825260208083209190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff94851617905584519284168252600290526040902080547fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166bffffffffffffffffffffffff909216919091179055600380548061069d5761069d6110aa565b6000828152602080822083017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90810180547fffffffffffffffffffffffff000000000000000000000000000000000000000016905590920190925573ffffffffffffffffffffffffffffffffffffffff86168083526002909152604080832080547fffffffffffffffffffffffffffffffffffffff000000000000000000000000001690555190917fcf91daec21e3510e2f2aea4b09d08c235d5c6844980be709f282ef591dbf420c91a250505050565b606060038054806020026020016040519081016040528092919081815260200182805480156107d457602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116107a9575b5050505050905090565b6107e66109b5565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600260209081526040918290208251808401909352546bffffffffffffffffffffffff811683526c01000000000000000000000000900460ff1615801591830191909152610894576040517f3a4406b500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831660048201526024016101eb565b60016020808301828152600380546bffffffffffffffffffffffff908116865273ffffffffffffffffffffffffffffffffffffffff871660008181526002909552604080862088518154965115156c01000000000000000000000000027fffffffffffffffffffffffffffffffffffffff0000000000000000000000000090971694169390931794909417909155815494850182559083527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b90930180547fffffffffffffffffffffffff00000000000000000000000000000000000000001684179055517f78f53b26906785548b265fa08f4197f9f3fff73fe0d504d30400aacb527f4ce09190a25050565b6109a96109b5565b6109b281610a38565b50565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a36576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016101eb565b565b3373ffffffffffffffffffffffffffffffffffffffff821603610ab7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016101eb565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b73ffffffffffffffffffffffffffffffffffffffff811681146109b257600080fd5b600060208284031215610b6157600080fd5b8135610b6c81610b2d565b9392505050565b60008060408385031215610b8657600080fd5b8235610b9181610b2d565b9150602083013567ffffffffffffffff811115610bad57600080fd5b830160808186031215610bbf57600080fd5b809150509250929050565b6020808252825182820181905260009190848201906040850190845b81811015610c1857835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610be6565b50909695505050505050565b8035610c2f81610b2d565b919050565b803567ffffffffffffffff81168114610c2f57600080fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610c8157600080fd5b830160208101925035905067ffffffffffffffff811115610ca157600080fd5b8060051b3603831315610cb357600080fd5b9250929050565b81835260007f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115610cec57600080fd5b8260051b8083602087013760009401602001938452509192915050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112610d3e57600080fd5b830160208101925035905067ffffffffffffffff811115610d5e57600080fd5b803603831315610cb357600080fd5b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6000602080835283358184015267ffffffffffffffff610dd7828601610c34565b1660408401526040840135610deb81610b2d565b73ffffffffffffffffffffffffffffffffffffffff8082166060860152606086013591507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff21863603018212610e3f57600080fd5b608085810152908501906101808501610e588380610c4c565b60e060a089015291829052906000906101a088015b81831015610e96578335610e8081610b2d565b8516815292860192600192909201918601610e6d565b610ea287870187610c4c565b975094507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff609350838982030160c08a0152610ede818887610cba565b9650505050604083013560e0870152610ef960608401610c24565b73ffffffffffffffffffffffffffffffffffffffff81166101008801529150610f2460808401610c24565b73ffffffffffffffffffffffffffffffffffffffff81166101208801529150610f5060a0840184610d09565b92508187860301610140880152610f68858483610d6d565b945050610f7860c0840184610d09565b93509150808685030161016087015250610f93838383610d6d565b9695505050505050565b600060208284031215610faf57600080fd5b610b6c82610c34565b67ffffffffffffffff8316815260006020604081840152835180604085015260005b81811015610ff657858101830151858201606001528201610fda565b81811115611008576000606083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01692909201606001949350505050565b600082821015611076577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea164736f6c634300080d000a",
}

var OffRampRouterABI = OffRampRouterMetaData.ABI

var OffRampRouterBin = OffRampRouterMetaData.Bin

func DeployOffRampRouter(auth *bind.TransactOpts, backend bind.ContractBackend, offRamps []common.Address) (common.Address, *types.Transaction, *OffRampRouter, error) {
	parsed, err := OffRampRouterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OffRampRouterBin), backend, offRamps)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffRampRouter{OffRampRouterCaller: OffRampRouterCaller{contract: contract}, OffRampRouterTransactor: OffRampRouterTransactor{contract: contract}, OffRampRouterFilterer: OffRampRouterFilterer{contract: contract}}, nil
}

type OffRampRouter struct {
	address common.Address
	abi     abi.ABI
	OffRampRouterCaller
	OffRampRouterTransactor
	OffRampRouterFilterer
}

type OffRampRouterCaller struct {
	contract *bind.BoundContract
}

type OffRampRouterTransactor struct {
	contract *bind.BoundContract
}

type OffRampRouterFilterer struct {
	contract *bind.BoundContract
}

type OffRampRouterSession struct {
	Contract     *OffRampRouter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OffRampRouterCallerSession struct {
	Contract *OffRampRouterCaller
	CallOpts bind.CallOpts
}

type OffRampRouterTransactorSession struct {
	Contract     *OffRampRouterTransactor
	TransactOpts bind.TransactOpts
}

type OffRampRouterRaw struct {
	Contract *OffRampRouter
}

type OffRampRouterCallerRaw struct {
	Contract *OffRampRouterCaller
}

type OffRampRouterTransactorRaw struct {
	Contract *OffRampRouterTransactor
}

func NewOffRampRouter(address common.Address, backend bind.ContractBackend) (*OffRampRouter, error) {
	abi, err := abi.JSON(strings.NewReader(OffRampRouterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOffRampRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffRampRouter{address: address, abi: abi, OffRampRouterCaller: OffRampRouterCaller{contract: contract}, OffRampRouterTransactor: OffRampRouterTransactor{contract: contract}, OffRampRouterFilterer: OffRampRouterFilterer{contract: contract}}, nil
}

func NewOffRampRouterCaller(address common.Address, caller bind.ContractCaller) (*OffRampRouterCaller, error) {
	contract, err := bindOffRampRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampRouterCaller{contract: contract}, nil
}

func NewOffRampRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*OffRampRouterTransactor, error) {
	contract, err := bindOffRampRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffRampRouterTransactor{contract: contract}, nil
}

func NewOffRampRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*OffRampRouterFilterer, error) {
	contract, err := bindOffRampRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffRampRouterFilterer{contract: contract}, nil
}

func bindOffRampRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OffRampRouterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_OffRampRouter *OffRampRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRampRouter.Contract.OffRampRouterCaller.contract.Call(opts, result, method, params...)
}

func (_OffRampRouter *OffRampRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampRouter.Contract.OffRampRouterTransactor.contract.Transfer(opts)
}

func (_OffRampRouter *OffRampRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRampRouter.Contract.OffRampRouterTransactor.contract.Transact(opts, method, params...)
}

func (_OffRampRouter *OffRampRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffRampRouter.Contract.contract.Call(opts, result, method, params...)
}

func (_OffRampRouter *OffRampRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampRouter.Contract.contract.Transfer(opts)
}

func (_OffRampRouter *OffRampRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffRampRouter.Contract.contract.Transact(opts, method, params...)
}

func (_OffRampRouter *OffRampRouterCaller) GetOffRamps(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _OffRampRouter.contract.Call(opts, &out, "getOffRamps")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_OffRampRouter *OffRampRouterSession) GetOffRamps() ([]common.Address, error) {
	return _OffRampRouter.Contract.GetOffRamps(&_OffRampRouter.CallOpts)
}

func (_OffRampRouter *OffRampRouterCallerSession) GetOffRamps() ([]common.Address, error) {
	return _OffRampRouter.Contract.GetOffRamps(&_OffRampRouter.CallOpts)
}

func (_OffRampRouter *OffRampRouterCaller) IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error) {
	var out []interface{}
	err := _OffRampRouter.contract.Call(opts, &out, "isOffRamp", offRamp)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_OffRampRouter *OffRampRouterSession) IsOffRamp(offRamp common.Address) (bool, error) {
	return _OffRampRouter.Contract.IsOffRamp(&_OffRampRouter.CallOpts, offRamp)
}

func (_OffRampRouter *OffRampRouterCallerSession) IsOffRamp(offRamp common.Address) (bool, error) {
	return _OffRampRouter.Contract.IsOffRamp(&_OffRampRouter.CallOpts, offRamp)
}

func (_OffRampRouter *OffRampRouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OffRampRouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_OffRampRouter *OffRampRouterSession) Owner() (common.Address, error) {
	return _OffRampRouter.Contract.Owner(&_OffRampRouter.CallOpts)
}

func (_OffRampRouter *OffRampRouterCallerSession) Owner() (common.Address, error) {
	return _OffRampRouter.Contract.Owner(&_OffRampRouter.CallOpts)
}

func (_OffRampRouter *OffRampRouterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffRampRouter.contract.Transact(opts, "acceptOwnership")
}

func (_OffRampRouter *OffRampRouterSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRampRouter.Contract.AcceptOwnership(&_OffRampRouter.TransactOpts)
}

func (_OffRampRouter *OffRampRouterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _OffRampRouter.Contract.AcceptOwnership(&_OffRampRouter.TransactOpts)
}

func (_OffRampRouter *OffRampRouterTransactor) AddOffRamp(opts *bind.TransactOpts, offRamp common.Address) (*types.Transaction, error) {
	return _OffRampRouter.contract.Transact(opts, "addOffRamp", offRamp)
}

func (_OffRampRouter *OffRampRouterSession) AddOffRamp(offRamp common.Address) (*types.Transaction, error) {
	return _OffRampRouter.Contract.AddOffRamp(&_OffRampRouter.TransactOpts, offRamp)
}

func (_OffRampRouter *OffRampRouterTransactorSession) AddOffRamp(offRamp common.Address) (*types.Transaction, error) {
	return _OffRampRouter.Contract.AddOffRamp(&_OffRampRouter.TransactOpts, offRamp)
}

func (_OffRampRouter *OffRampRouterTransactor) RemoveOffRamp(opts *bind.TransactOpts, offRamp common.Address) (*types.Transaction, error) {
	return _OffRampRouter.contract.Transact(opts, "removeOffRamp", offRamp)
}

func (_OffRampRouter *OffRampRouterSession) RemoveOffRamp(offRamp common.Address) (*types.Transaction, error) {
	return _OffRampRouter.Contract.RemoveOffRamp(&_OffRampRouter.TransactOpts, offRamp)
}

func (_OffRampRouter *OffRampRouterTransactorSession) RemoveOffRamp(offRamp common.Address) (*types.Transaction, error) {
	return _OffRampRouter.Contract.RemoveOffRamp(&_OffRampRouter.TransactOpts, offRamp)
}

func (_OffRampRouter *OffRampRouterTransactor) RouteMessage(opts *bind.TransactOpts, receiver common.Address, message CCIPMessage) (*types.Transaction, error) {
	return _OffRampRouter.contract.Transact(opts, "routeMessage", receiver, message)
}

func (_OffRampRouter *OffRampRouterSession) RouteMessage(receiver common.Address, message CCIPMessage) (*types.Transaction, error) {
	return _OffRampRouter.Contract.RouteMessage(&_OffRampRouter.TransactOpts, receiver, message)
}

func (_OffRampRouter *OffRampRouterTransactorSession) RouteMessage(receiver common.Address, message CCIPMessage) (*types.Transaction, error) {
	return _OffRampRouter.Contract.RouteMessage(&_OffRampRouter.TransactOpts, receiver, message)
}

func (_OffRampRouter *OffRampRouterTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _OffRampRouter.contract.Transact(opts, "transferOwnership", to)
}

func (_OffRampRouter *OffRampRouterSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRampRouter.Contract.TransferOwnership(&_OffRampRouter.TransactOpts, to)
}

func (_OffRampRouter *OffRampRouterTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _OffRampRouter.Contract.TransferOwnership(&_OffRampRouter.TransactOpts, to)
}

type OffRampRouterOffRampAddedIterator struct {
	Event *OffRampRouterOffRampAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampRouterOffRampAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampRouterOffRampAdded)
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
		it.Event = new(OffRampRouterOffRampAdded)
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

func (it *OffRampRouterOffRampAddedIterator) Error() error {
	return it.fail
}

func (it *OffRampRouterOffRampAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampRouterOffRampAdded struct {
	OffRamp common.Address
	Raw     types.Log
}

func (_OffRampRouter *OffRampRouterFilterer) FilterOffRampAdded(opts *bind.FilterOpts, offRamp []common.Address) (*OffRampRouterOffRampAddedIterator, error) {

	var offRampRule []interface{}
	for _, offRampItem := range offRamp {
		offRampRule = append(offRampRule, offRampItem)
	}

	logs, sub, err := _OffRampRouter.contract.FilterLogs(opts, "OffRampAdded", offRampRule)
	if err != nil {
		return nil, err
	}
	return &OffRampRouterOffRampAddedIterator{contract: _OffRampRouter.contract, event: "OffRampAdded", logs: logs, sub: sub}, nil
}

func (_OffRampRouter *OffRampRouterFilterer) WatchOffRampAdded(opts *bind.WatchOpts, sink chan<- *OffRampRouterOffRampAdded, offRamp []common.Address) (event.Subscription, error) {

	var offRampRule []interface{}
	for _, offRampItem := range offRamp {
		offRampRule = append(offRampRule, offRampItem)
	}

	logs, sub, err := _OffRampRouter.contract.WatchLogs(opts, "OffRampAdded", offRampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampRouterOffRampAdded)
				if err := _OffRampRouter.contract.UnpackLog(event, "OffRampAdded", log); err != nil {
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

func (_OffRampRouter *OffRampRouterFilterer) ParseOffRampAdded(log types.Log) (*OffRampRouterOffRampAdded, error) {
	event := new(OffRampRouterOffRampAdded)
	if err := _OffRampRouter.contract.UnpackLog(event, "OffRampAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampRouterOffRampRemovedIterator struct {
	Event *OffRampRouterOffRampRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampRouterOffRampRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampRouterOffRampRemoved)
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
		it.Event = new(OffRampRouterOffRampRemoved)
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

func (it *OffRampRouterOffRampRemovedIterator) Error() error {
	return it.fail
}

func (it *OffRampRouterOffRampRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampRouterOffRampRemoved struct {
	OffRamp common.Address
	Raw     types.Log
}

func (_OffRampRouter *OffRampRouterFilterer) FilterOffRampRemoved(opts *bind.FilterOpts, offRamp []common.Address) (*OffRampRouterOffRampRemovedIterator, error) {

	var offRampRule []interface{}
	for _, offRampItem := range offRamp {
		offRampRule = append(offRampRule, offRampItem)
	}

	logs, sub, err := _OffRampRouter.contract.FilterLogs(opts, "OffRampRemoved", offRampRule)
	if err != nil {
		return nil, err
	}
	return &OffRampRouterOffRampRemovedIterator{contract: _OffRampRouter.contract, event: "OffRampRemoved", logs: logs, sub: sub}, nil
}

func (_OffRampRouter *OffRampRouterFilterer) WatchOffRampRemoved(opts *bind.WatchOpts, sink chan<- *OffRampRouterOffRampRemoved, offRamp []common.Address) (event.Subscription, error) {

	var offRampRule []interface{}
	for _, offRampItem := range offRamp {
		offRampRule = append(offRampRule, offRampItem)
	}

	logs, sub, err := _OffRampRouter.contract.WatchLogs(opts, "OffRampRemoved", offRampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampRouterOffRampRemoved)
				if err := _OffRampRouter.contract.UnpackLog(event, "OffRampRemoved", log); err != nil {
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

func (_OffRampRouter *OffRampRouterFilterer) ParseOffRampRemoved(log types.Log) (*OffRampRouterOffRampRemoved, error) {
	event := new(OffRampRouterOffRampRemoved)
	if err := _OffRampRouter.contract.UnpackLog(event, "OffRampRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampRouterOwnershipTransferRequestedIterator struct {
	Event *OffRampRouterOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampRouterOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampRouterOwnershipTransferRequested)
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
		it.Event = new(OffRampRouterOwnershipTransferRequested)
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

func (it *OffRampRouterOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OffRampRouterOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampRouterOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRampRouter *OffRampRouterFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampRouterOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampRouter.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampRouterOwnershipTransferRequestedIterator{contract: _OffRampRouter.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OffRampRouter *OffRampRouterFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampRouter.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampRouterOwnershipTransferRequested)
				if err := _OffRampRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OffRampRouter *OffRampRouterFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffRampRouterOwnershipTransferRequested, error) {
	event := new(OffRampRouterOwnershipTransferRequested)
	if err := _OffRampRouter.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffRampRouterOwnershipTransferredIterator struct {
	Event *OffRampRouterOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffRampRouterOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffRampRouterOwnershipTransferred)
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
		it.Event = new(OffRampRouterOwnershipTransferred)
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

func (it *OffRampRouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OffRampRouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffRampRouterOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffRampRouter *OffRampRouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampRouterOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampRouter.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffRampRouterOwnershipTransferredIterator{contract: _OffRampRouter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OffRampRouter *OffRampRouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffRampRouter.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffRampRouterOwnershipTransferred)
				if err := _OffRampRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OffRampRouter *OffRampRouterFilterer) ParseOwnershipTransferred(log types.Log) (*OffRampRouterOwnershipTransferred, error) {
	event := new(OffRampRouterOwnershipTransferred)
	if err := _OffRampRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OffRampRouter *OffRampRouter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OffRampRouter.abi.Events["OffRampAdded"].ID:
		return _OffRampRouter.ParseOffRampAdded(log)
	case _OffRampRouter.abi.Events["OffRampRemoved"].ID:
		return _OffRampRouter.ParseOffRampRemoved(log)
	case _OffRampRouter.abi.Events["OwnershipTransferRequested"].ID:
		return _OffRampRouter.ParseOwnershipTransferRequested(log)
	case _OffRampRouter.abi.Events["OwnershipTransferred"].ID:
		return _OffRampRouter.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OffRampRouterOffRampAdded) Topic() common.Hash {
	return common.HexToHash("0x78f53b26906785548b265fa08f4197f9f3fff73fe0d504d30400aacb527f4ce0")
}

func (OffRampRouterOffRampRemoved) Topic() common.Hash {
	return common.HexToHash("0xcf91daec21e3510e2f2aea4b09d08c235d5c6844980be709f282ef591dbf420c")
}

func (OffRampRouterOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OffRampRouterOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_OffRampRouter *OffRampRouter) Address() common.Address {
	return _OffRampRouter.address
}

type OffRampRouterInterface interface {
	GetOffRamps(opts *bind.CallOpts) ([]common.Address, error)

	IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddOffRamp(opts *bind.TransactOpts, offRamp common.Address) (*types.Transaction, error)

	RemoveOffRamp(opts *bind.TransactOpts, offRamp common.Address) (*types.Transaction, error)

	RouteMessage(opts *bind.TransactOpts, receiver common.Address, message CCIPMessage) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOffRampAdded(opts *bind.FilterOpts, offRamp []common.Address) (*OffRampRouterOffRampAddedIterator, error)

	WatchOffRampAdded(opts *bind.WatchOpts, sink chan<- *OffRampRouterOffRampAdded, offRamp []common.Address) (event.Subscription, error)

	ParseOffRampAdded(log types.Log) (*OffRampRouterOffRampAdded, error)

	FilterOffRampRemoved(opts *bind.FilterOpts, offRamp []common.Address) (*OffRampRouterOffRampRemovedIterator, error)

	WatchOffRampRemoved(opts *bind.WatchOpts, sink chan<- *OffRampRouterOffRampRemoved, offRamp []common.Address) (event.Subscription, error)

	ParseOffRampRemoved(log types.Log) (*OffRampRouterOffRampRemoved, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampRouterOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffRampRouterOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OffRampRouterOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffRampRouterOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffRampRouterOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OffRampRouterOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
