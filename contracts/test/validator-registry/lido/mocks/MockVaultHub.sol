// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

contract MockVaultHub {
    mapping(address => bool) public connected;
    mapping(address => uint256) public _totalValue;

    function setConnected(address _vault, bool _connected) external {
        connected[_vault] = _connected;
    }

    function setTotalValue(address _vault, uint256 _amount) external {
        _totalValue[_vault] = _amount;
    }

    // ---- IVaultHubMinimal ----
    function isVaultConnected(address _vault) external view returns (bool) {
        return connected[_vault];
    }

    function totalValue(address _vault) external view returns (uint256) {
        return _totalValue[_vault];
    }
}
