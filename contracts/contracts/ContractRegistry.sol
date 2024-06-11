// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

contract ContractRegistry {
    struct ContractInfo {
        string name;
        address addr;
    }

    ContractInfo[] public contracts;
    mapping(string => uint256) private nameToIndex;
    mapping(string => bool) private nameExists;

    event ContractAdded(string name, address addr);
    event ContractUpdated(string name, address addr);

    function addContract(string memory _name, address _addr) public {
        require(!nameExists[_name], "Contract name already exists");
        contracts.push(ContractInfo({name: _name, addr: _addr}));
        nameToIndex[_name] = contracts.length - 1;
        nameExists[_name] = true;
        emit ContractAdded(_name, _addr);
    }

    function updateContract(string memory _name, address _addr) public {
        require(nameExists[_name], "Contract name does not exist");
        uint256 index = nameToIndex[_name];
        contracts[index].addr = _addr;
        emit ContractUpdated(_name, _addr);
    }

    function getContractAddress(string memory _name) public view returns (address) {
        require(nameExists[_name], "Contract name does not exist");
        uint256 index = nameToIndex[_name];
        return contracts[index].addr;
    }

    function getAllContracts() public view returns (ContractInfo[] memory) {
        return contracts;
    }
}
