// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {MevCommitMiddleware} from "../../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";

contract MevCommitMiddlewareTest is Test {

    MevCommitMiddleware public mevCommitMiddleware;

    function setUp() public {
        mevCommitMiddleware = new MevCommitMiddleware();
    }
}
