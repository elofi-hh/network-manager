// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package store

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
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

// StoreMetaData contains all meta data concerning the Store contract.
var StoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"initialKFactor\",\"type\":\"uint16\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"newKFactor\",\"type\":\"uint16\"}],\"name\":\"KFactorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"userId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"rank\",\"type\":\"uint16\"}],\"name\":\"RankingUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"K_FACTOR\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"userId\",\"type\":\"string\"}],\"name\":\"getRanking\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"rankings\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"newKFactor\",\"type\":\"uint16\"}],\"name\":\"setKFactor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"userId\",\"type\":\"string\"},{\"internalType\":\"uint16\",\"name\":\"rank\",\"type\":\"uint16\"}],\"name\":\"setRanking\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"userId\",\"type\":\"string\"},{\"internalType\":\"uint16\",\"name\":\"networkEloRating\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"isAbuser\",\"type\":\"bool\"}],\"name\":\"updateRanking\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StoreABI is the input ABI used to generate the binding from.
// Deprecated: Use StoreMetaData.ABI instead.
var StoreABI = StoreMetaData.ABI

// Store is an auto generated Go binding around an Ethereum contract.
type Store struct {
	StoreCaller     // Read-only binding to the contract
	StoreTransactor // Write-only binding to the contract
	StoreFilterer   // Log filterer for contract events
}

// StoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type StoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StoreSession struct {
	Contract     *Store            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StoreCallerSession struct {
	Contract *StoreCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StoreTransactorSession struct {
	Contract     *StoreTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type StoreRaw struct {
	Contract *Store // Generic contract binding to access the raw methods on
}

// StoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StoreCallerRaw struct {
	Contract *StoreCaller // Generic read-only contract binding to access the raw methods on
}

// StoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StoreTransactorRaw struct {
	Contract *StoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStore creates a new instance of Store, bound to a specific deployed contract.
func NewStore(address common.Address, backend bind.ContractBackend) (*Store, error) {
	contract, err := bindStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Store{StoreCaller: StoreCaller{contract: contract}, StoreTransactor: StoreTransactor{contract: contract}, StoreFilterer: StoreFilterer{contract: contract}}, nil
}

// NewStoreCaller creates a new read-only instance of Store, bound to a specific deployed contract.
func NewStoreCaller(address common.Address, caller bind.ContractCaller) (*StoreCaller, error) {
	contract, err := bindStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StoreCaller{contract: contract}, nil
}

// NewStoreTransactor creates a new write-only instance of Store, bound to a specific deployed contract.
func NewStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*StoreTransactor, error) {
	contract, err := bindStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StoreTransactor{contract: contract}, nil
}

// NewStoreFilterer creates a new log filterer instance of Store, bound to a specific deployed contract.
func NewStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*StoreFilterer, error) {
	contract, err := bindStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StoreFilterer{contract: contract}, nil
}

// bindStore binds a generic wrapper to an already deployed contract.
func bindStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *StoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Store.Contract.StoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *StoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.StoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *StoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.StoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *StoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Store.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *StoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *StoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.contract.Transact(opts, method, params...)
}

// KFACTOR is a free data retrieval call binding the contract method 0x77519b73.
//
// Solidity: function K_FACTOR() view returns(uint16)
func (_Store *StoreCaller) KFACTOR(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "K_FACTOR")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// KFACTOR is a free data retrieval call binding the contract method 0x77519b73.
//
// Solidity: function K_FACTOR() view returns(uint16)
func (_Store *StoreSession) KFACTOR() (uint16, error) {
	return _Store.Contract.KFACTOR(&_Store.CallOpts)
}

// KFACTOR is a free data retrieval call binding the contract method 0x77519b73.
//
// Solidity: function K_FACTOR() view returns(uint16)
func (_Store *StoreCallerSession) KFACTOR() (uint16, error) {
	return _Store.Contract.KFACTOR(&_Store.CallOpts)
}

// GetRanking is a free data retrieval call binding the contract method 0x46e7a1dc.
//
// Solidity: function getRanking(string userId) view returns(uint16)
func (_Store *StoreCaller) GetRanking(opts *bind.CallOpts, userId string) (uint16, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "getRanking", userId)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetRanking is a free data retrieval call binding the contract method 0x46e7a1dc.
//
// Solidity: function getRanking(string userId) view returns(uint16)
func (_Store *StoreSession) GetRanking(userId string) (uint16, error) {
	return _Store.Contract.GetRanking(&_Store.CallOpts, userId)
}

// GetRanking is a free data retrieval call binding the contract method 0x46e7a1dc.
//
// Solidity: function getRanking(string userId) view returns(uint16)
func (_Store *StoreCallerSession) GetRanking(userId string) (uint16, error) {
	return _Store.Contract.GetRanking(&_Store.CallOpts, userId)
}

// Rankings is a free data retrieval call binding the contract method 0xfdf3f74d.
//
// Solidity: function rankings(string ) view returns(uint16)
func (_Store *StoreCaller) Rankings(opts *bind.CallOpts, arg0 string) (uint16, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "rankings", arg0)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// Rankings is a free data retrieval call binding the contract method 0xfdf3f74d.
//
// Solidity: function rankings(string ) view returns(uint16)
func (_Store *StoreSession) Rankings(arg0 string) (uint16, error) {
	return _Store.Contract.Rankings(&_Store.CallOpts, arg0)
}

// Rankings is a free data retrieval call binding the contract method 0xfdf3f74d.
//
// Solidity: function rankings(string ) view returns(uint16)
func (_Store *StoreCallerSession) Rankings(arg0 string) (uint16, error) {
	return _Store.Contract.Rankings(&_Store.CallOpts, arg0)
}

// SetKFactor is a paid mutator transaction binding the contract method 0x2a641fdb.
//
// Solidity: function setKFactor(uint16 newKFactor) returns()
func (_Store *StoreTransactor) SetKFactor(opts *bind.TransactOpts, newKFactor uint16) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "setKFactor", newKFactor)
}

// SetKFactor is a paid mutator transaction binding the contract method 0x2a641fdb.
//
// Solidity: function setKFactor(uint16 newKFactor) returns()
func (_Store *StoreSession) SetKFactor(newKFactor uint16) (*types.Transaction, error) {
	return _Store.Contract.SetKFactor(&_Store.TransactOpts, newKFactor)
}

// SetKFactor is a paid mutator transaction binding the contract method 0x2a641fdb.
//
// Solidity: function setKFactor(uint16 newKFactor) returns()
func (_Store *StoreTransactorSession) SetKFactor(newKFactor uint16) (*types.Transaction, error) {
	return _Store.Contract.SetKFactor(&_Store.TransactOpts, newKFactor)
}

// SetRanking is a paid mutator transaction binding the contract method 0x1f5bc0cf.
//
// Solidity: function setRanking(string userId, uint16 rank) returns()
func (_Store *StoreTransactor) SetRanking(opts *bind.TransactOpts, userId string, rank uint16) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "setRanking", userId, rank)
}

// SetRanking is a paid mutator transaction binding the contract method 0x1f5bc0cf.
//
// Solidity: function setRanking(string userId, uint16 rank) returns()
func (_Store *StoreSession) SetRanking(userId string, rank uint16) (*types.Transaction, error) {
	return _Store.Contract.SetRanking(&_Store.TransactOpts, userId, rank)
}

// SetRanking is a paid mutator transaction binding the contract method 0x1f5bc0cf.
//
// Solidity: function setRanking(string userId, uint16 rank) returns()
func (_Store *StoreTransactorSession) SetRanking(userId string, rank uint16) (*types.Transaction, error) {
	return _Store.Contract.SetRanking(&_Store.TransactOpts, userId, rank)
}

// UpdateRanking is a paid mutator transaction binding the contract method 0x950b409b.
//
// Solidity: function updateRanking(string userId, uint16 networkEloRating, bool isAbuser) returns()
func (_Store *StoreTransactor) UpdateRanking(opts *bind.TransactOpts, userId string, networkEloRating uint16, isAbuser bool) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "updateRanking", userId, networkEloRating, isAbuser)
}

// UpdateRanking is a paid mutator transaction binding the contract method 0x950b409b.
//
// Solidity: function updateRanking(string userId, uint16 networkEloRating, bool isAbuser) returns()
func (_Store *StoreSession) UpdateRanking(userId string, networkEloRating uint16, isAbuser bool) (*types.Transaction, error) {
	return _Store.Contract.UpdateRanking(&_Store.TransactOpts, userId, networkEloRating, isAbuser)
}

// UpdateRanking is a paid mutator transaction binding the contract method 0x950b409b.
//
// Solidity: function updateRanking(string userId, uint16 networkEloRating, bool isAbuser) returns()
func (_Store *StoreTransactorSession) UpdateRanking(userId string, networkEloRating uint16, isAbuser bool) (*types.Transaction, error) {
	return _Store.Contract.UpdateRanking(&_Store.TransactOpts, userId, networkEloRating, isAbuser)
}

// StoreKFactorUpdatedIterator is returned from FilterKFactorUpdated and is used to iterate over the raw logs and unpacked data for KFactorUpdated events raised by the Store contract.
type StoreKFactorUpdatedIterator struct {
	Event *StoreKFactorUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StoreKFactorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StoreKFactorUpdated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StoreKFactorUpdated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StoreKFactorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StoreKFactorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StoreKFactorUpdated represents a KFactorUpdated event raised by the Store contract.
type StoreKFactorUpdated struct {
	NewKFactor uint16
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterKFactorUpdated is a free log retrieval operation binding the contract event 0x02c29147f9dbc311a4ccf9949488a20725c70e7f7ef01e2f436c86cdabc2d2d8.
//
// Solidity: event KFactorUpdated(uint16 newKFactor)
func (_Store *StoreFilterer) FilterKFactorUpdated(opts *bind.FilterOpts) (*StoreKFactorUpdatedIterator, error) {

	logs, sub, err := _Store.contract.FilterLogs(opts, "KFactorUpdated")
	if err != nil {
		return nil, err
	}
	return &StoreKFactorUpdatedIterator{contract: _Store.contract, event: "KFactorUpdated", logs: logs, sub: sub}, nil
}

// WatchKFactorUpdated is a free log subscription operation binding the contract event 0x02c29147f9dbc311a4ccf9949488a20725c70e7f7ef01e2f436c86cdabc2d2d8.
//
// Solidity: event KFactorUpdated(uint16 newKFactor)
func (_Store *StoreFilterer) WatchKFactorUpdated(opts *bind.WatchOpts, sink chan<- *StoreKFactorUpdated) (event.Subscription, error) {

	logs, sub, err := _Store.contract.WatchLogs(opts, "KFactorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StoreKFactorUpdated)
				if err := _Store.contract.UnpackLog(event, "KFactorUpdated", log); err != nil {
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

// ParseKFactorUpdated is a log parse operation binding the contract event 0x02c29147f9dbc311a4ccf9949488a20725c70e7f7ef01e2f436c86cdabc2d2d8.
//
// Solidity: event KFactorUpdated(uint16 newKFactor)
func (_Store *StoreFilterer) ParseKFactorUpdated(log types.Log) (*StoreKFactorUpdated, error) {
	event := new(StoreKFactorUpdated)
	if err := _Store.contract.UnpackLog(event, "KFactorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StoreRankingUpdatedIterator is returned from FilterRankingUpdated and is used to iterate over the raw logs and unpacked data for RankingUpdated events raised by the Store contract.
type StoreRankingUpdatedIterator struct {
	Event *StoreRankingUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StoreRankingUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StoreRankingUpdated)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StoreRankingUpdated)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StoreRankingUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StoreRankingUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StoreRankingUpdated represents a RankingUpdated event raised by the Store contract.
type StoreRankingUpdated struct {
	UserId string
	Rank   uint16
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRankingUpdated is a free log retrieval operation binding the contract event 0x60823246bd9724a87f737174a65e2a1c68a909f877c4005413c10764f781b8f1.
//
// Solidity: event RankingUpdated(string userId, uint16 rank)
func (_Store *StoreFilterer) FilterRankingUpdated(opts *bind.FilterOpts) (*StoreRankingUpdatedIterator, error) {

	logs, sub, err := _Store.contract.FilterLogs(opts, "RankingUpdated")
	if err != nil {
		return nil, err
	}
	return &StoreRankingUpdatedIterator{contract: _Store.contract, event: "RankingUpdated", logs: logs, sub: sub}, nil
}

// WatchRankingUpdated is a free log subscription operation binding the contract event 0x60823246bd9724a87f737174a65e2a1c68a909f877c4005413c10764f781b8f1.
//
// Solidity: event RankingUpdated(string userId, uint16 rank)
func (_Store *StoreFilterer) WatchRankingUpdated(opts *bind.WatchOpts, sink chan<- *StoreRankingUpdated) (event.Subscription, error) {

	logs, sub, err := _Store.contract.WatchLogs(opts, "RankingUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StoreRankingUpdated)
				if err := _Store.contract.UnpackLog(event, "RankingUpdated", log); err != nil {
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

// ParseRankingUpdated is a log parse operation binding the contract event 0x60823246bd9724a87f737174a65e2a1c68a909f877c4005413c10764f781b8f1.
//
// Solidity: event RankingUpdated(string userId, uint16 rank)
func (_Store *StoreFilterer) ParseRankingUpdated(log types.Log) (*StoreRankingUpdated, error) {
	event := new(StoreRankingUpdated)
	if err := _Store.contract.UnpackLog(event, "RankingUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

