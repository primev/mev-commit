// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

// Adjust these paths to your repo layout:
import {ValidatorOptInHub} from "../../contracts/validator-registry/ValidatorOptInHub.sol";
import {IValidatorOptInHub} from "../../contracts/interfaces/IValidatorOptInHub.sol";
import {IRegistry} from "../../contracts/interfaces/IRegistry.sol";
import {Errors} from "../../contracts/utils/Errors.sol";

contract ValidatorOptInHubTest is Test {
    ValidatorOptInHub hubImplementation;
    ValidatorOptInHub hub;

    MockRegistry registryA;
    MockRegistry registryB;
    MockRegistry registryC;
    MockRegistry registryD; // used when testing addRegistry on top of 3 initial ones

    address badRegistry = address(0xBAD);
    address contractOwner = address(this);

    // Sample validator public keys
    bytes validatorPublicKeyOne   = hex"01";
    bytes validatorPublicKeyTwo   = hex"02";
    bytes validatorPublicKeyThree = hex"03";

    function setUp() public {
        // --- Deploy registry mocks ---
        registryA = new MockRegistry();
        registryB = new MockRegistry();
        registryC = new MockRegistry();
        registryD = new MockRegistry();

        // Seed sample statuses:
        // - registryA: pk1=true, pk2=false
        // - registryB: pk1=false, pk2=true
        // - registryC/D: defaults to false
        registryA.setOptIn(validatorPublicKeyOne, true);
        registryA.setOptIn(validatorPublicKeyTwo, false);
        registryB.setOptIn(validatorPublicKeyOne, false);
        registryB.setOptIn(validatorPublicKeyTwo, true);

        // --- Deploy hub implementation (logic) ---
        hubImplementation = new ValidatorOptInHub();

        // Prepare initial registries for initializer (NOW 3)
        address[] memory initialRegistries = new address[](3);
        initialRegistries[0] = address(registryA);
        initialRegistries[1] = address(registryB);
        initialRegistries[2] = address(registryC);

        // --- Deploy UUPS proxy with initializer call ---
        bytes memory initializerCalldata =
            abi.encodeCall(ValidatorOptInHub.initialize, (initialRegistries, contractOwner));
        ERC1967Proxy proxy = new ERC1967Proxy(address(hubImplementation), initializerCalldata);

        // Contract has payable receive/fallback → cast via payable(address)
        hub = ValidatorOptInHub(payable(address(proxy)));
    }

    // ---------------------------
    // Internal assertion helpers
    // ---------------------------

    function _expectCustomError(bytes4 selector) internal {
        vm.expectRevert(abi.encodeWithSelector(selector));
    }

    function _assertEqualBoolArray(bool[] memory actual, bool[] memory expected) internal {
        assertEq(actual.length, expected.length, "length mismatch");
        for (uint256 index = 0; index < actual.length; ++index) {
            assertEq(actual[index], expected[index], "element mismatch");
        }
    }

    function _assertEqualBoolMatrix(bool[][] memory actual, bool[][] memory expected) internal {
        assertEq(actual.length, expected.length, "row count mismatch");
        for (uint256 row = 0; row < actual.length; ++row) {
            _assertEqualBoolArray(actual[row], expected[row]);
        }
    }

    // ---------------
    // Core behaviours
    // ---------------

    function test_initialize_and_queryMatrix() public {
        // Query two validator keys
        bytes[] memory validatorPublicKeys = new bytes[](2);
        validatorPublicKeys[0] = validatorPublicKeyOne;
        validatorPublicKeys[1] = validatorPublicKeyTwo;

        bool[][] memory actualResultsMatrix = hub.areValidatorsOptedInList(validatorPublicKeys);

        // Expected columns NOW: [registryA, registryB, registryC]
        bool[][] memory expectedResultsMatrix = new bool[][](2);
        expectedResultsMatrix;
        expectedResultsMatrix[0] = new bool[](3);
        expectedResultsMatrix[1] = new bool[](3);

        // pk1: A=true,  B=false, C=false
        expectedResultsMatrix[0][0] = true;
        expectedResultsMatrix[0][1] = false;
        expectedResultsMatrix[0][2] = false;

        // pk2: A=false, B=true,  C=false
        expectedResultsMatrix[1][0] = false;
        expectedResultsMatrix[1][1] = true;
        expectedResultsMatrix[1][2] = false;

        _assertEqualBoolMatrix(actualResultsMatrix, expectedResultsMatrix);

        // “Any” aggregate checks
        assertTrue(hub.isValidatorOptedIn(validatorPublicKeyOne));
        assertTrue(hub.isValidatorOptedIn(validatorPublicKeyTwo));
        assertFalse(hub.isValidatorOptedIn(validatorPublicKeyThree));
    }

    function test_addRegistry_appends_and_affectsResults() public {
        // Add registryD (defaults false for all keys) on top of 3 initial ones
        hub.addRegistry(address(registryD));

        // pk1 remains true (registryA says true)
        assertTrue(hub.isValidatorOptedIn(validatorPublicKeyOne));

        // pk3 remains false (all registries false)
        assertFalse(hub.isValidatorOptedIn(validatorPublicKeyThree));

        // Matrix now has 4 columns
        bytes[] memory validatorPublicKeys = new bytes[](1);
        validatorPublicKeys[0] = validatorPublicKeyOne;

        bool[][] memory actualResultsMatrix = hub.areValidatorsOptedInList(validatorPublicKeys);
        assertEq(actualResultsMatrix[0].length, 4, "expected 4 registries");
        assertEq(actualResultsMatrix[0][0], true);  // A
        assertEq(actualResultsMatrix[0][1], false); // B
        assertEq(actualResultsMatrix[0][2], false); // C
        assertEq(actualResultsMatrix[0][3], false); // D (added)
    }

    function test_addRegistry_zeroAddress_reverts() public {
        _expectCustomError(IValidatorOptInHub.ZeroAddress.selector);
        hub.addRegistry(address(0));
    }

    function test_addRegistry_invalidRegistry_reverts_onProbe() public {
        // EOA / non-contract → probe should fail decoding and revert
        _expectCustomError(IValidatorOptInHub.InvalidRegistry.selector);
        hub.addRegistry(badRegistry);
    }

    function test_addRegistry_invalidRegistry_reverts_onProbe_wrongSelector() public {
        // Contract exists but does NOT implement isValidatorOptedIn(bytes) → staticcall ok with empty/invalid returndata
        // → ret.length != 32 → InvalidRegistry
        MockWrongSigRegistry wrong = new MockWrongSigRegistry();
        _expectCustomError(IValidatorOptInHub.InvalidRegistry.selector);
        hub.addRegistry(address(wrong));
    }

    function test_addRegistry_invalidRegistry_reverts_onProbe_mutating() public {
        // Contract implements the function but writes state → STATICCALL will fail (ok == false) → InvalidRegistry
        MockMutatingRegistry bad = new MockMutatingRegistry();
        _expectCustomError(IValidatorOptInHub.InvalidRegistry.selector);
        hub.addRegistry(address(bad));
    }

    function test_updateRegistry_replacesIndex_stably() public {
        // We start at [A, B, C]; replace index 1 (B) with C → [A, C, C]
        hub.updateRegistry(1, address(registryB), address(registryC));

        // pk2 was true on B, false on C → col1 should now be false
        bytes[] memory validatorPublicKeys = new bytes[](1);
        validatorPublicKeys[0] = validatorPublicKeyTwo;

        bool[][] memory actualResultsMatrix = hub.areValidatorsOptedInList(validatorPublicKeys);
        assertEq(actualResultsMatrix[0].length, 3);
        assertEq(actualResultsMatrix[0][0], false); // A says false for pk2
        assertEq(actualResultsMatrix[0][1], false); // replaced with C
        assertEq(actualResultsMatrix[0][2], false); // original C
    }

    function test_updateRegistry_invalidIndex_reverts() public {
        _expectCustomError(IValidatorOptInHub.InvalidIndex.selector);
        hub.updateRegistry(99, address(registryA), address(registryC));
    }

    function test_updateRegistry_zeroAddress_reverts() public {
        _expectCustomError(IValidatorOptInHub.ZeroAddress.selector);
        hub.updateRegistry(0, address(0), address(registryC));
    }

    function test_updateRegistry_indexRegistryMismatch_reverts() public {
        // Index 0 holds registryA; pass the wrong oldRegistry (registryB)
        _expectCustomError(IValidatorOptInHub.IndexRegistryMismatch.selector);
        hub.updateRegistry(0, address(registryB), address(registryC));
    }

    function test_removeRegistry_zeroesSlot_preservingIndexing() public {
        // “Logical delete” index 0 (registryA) to preserve column indices
        hub.removeRegistry(0, address(registryA));

        // pk1 used to be true at column 0; now zeroed → false
        bytes[] memory validatorPublicKeys = new bytes[](1);
        validatorPublicKeys[0] = validatorPublicKeyOne;

        bool[][] memory actualResultsMatrix = hub.areValidatorsOptedInList(validatorPublicKeys);
        assertEq(actualResultsMatrix[0][0], false);
        // Column 1 remains registryB and still reports false for pk1
        assertEq(actualResultsMatrix[0][1], false);
        // Column 2 remains registryC and reports false for pk1
        assertEq(actualResultsMatrix[0][2], false);
    }

    function test_removeRegistry_invalidIndex_reverts() public {
        _expectCustomError(IValidatorOptInHub.InvalidIndex.selector);
        hub.removeRegistry(42, address(registryA));
    }

    function test_removeRegistry_zeroAddress_reverts() public {
        _expectCustomError(IValidatorOptInHub.ZeroAddress.selector);
        hub.removeRegistry(0, address(0));
    }

    function test_removeRegistry_indexRegistryMismatch_reverts() public {
        _expectCustomError(IValidatorOptInHub.IndexRegistryMismatch.selector);
        hub.removeRegistry(1, address(registryA)); // slot 1 is registryB at setup
    }
    

    function test_receive_reverts() public {
        (bool ok, bytes memory r) = address(hub).call{value: 1 ether}("");
        assertFalse(ok, "receive should revert");
        assertGe(r.length, 4);
        assertEq(bytes4(r), Errors.InvalidReceive.selector);
    }

    function test_fallback_reverts() public {
        (bool ok, bytes memory r) = address(hub).call(hex"deadbeef");
        assertFalse(ok, "fallback should revert");
        assertGe(r.length, 4);
        assertEq(bytes4(r), Errors.InvalidFallback.selector);
    }

    function test_isValidatorOptedInList_returnsPerRegistryFlags_forEachValidator_onInitialSetup() public {
        // validatorPublicKeyOne is true on registryA, false on registryB/C
        bool[] memory perRegistryFlagsForValidatorOne = hub.isValidatorOptedInList(validatorPublicKeyOne);
        bool[] memory expectedFlagsForValidatorOne = new bool[](3);
        expectedFlagsForValidatorOne[0] = true;  // A
        expectedFlagsForValidatorOne[1] = false; // B
        expectedFlagsForValidatorOne[2] = false; // C
        _assertEqualBoolArray(perRegistryFlagsForValidatorOne, expectedFlagsForValidatorOne);

        // validatorPublicKeyTwo is false on registryA, true on registryB, false on registryC
        bool[] memory perRegistryFlagsForValidatorTwo = hub.isValidatorOptedInList(validatorPublicKeyTwo);
        bool[] memory expectedFlagsForValidatorTwo = new bool[](3);
        expectedFlagsForValidatorTwo[0] = false; // A
        expectedFlagsForValidatorTwo[1] = true;  // B
        expectedFlagsForValidatorTwo[2] = false; // C
        _assertEqualBoolArray(perRegistryFlagsForValidatorTwo, expectedFlagsForValidatorTwo);

        // validatorPublicKeyThree is false everywhere
        bool[] memory perRegistryFlagsForValidatorThree = hub.isValidatorOptedInList(validatorPublicKeyThree);
        bool[] memory expectedFlagsForValidatorThree = new bool[](3);
        expectedFlagsForValidatorThree[0] = false;
        expectedFlagsForValidatorThree[1] = false;
        expectedFlagsForValidatorThree[2] = false;
        _assertEqualBoolArray(perRegistryFlagsForValidatorThree, expectedFlagsForValidatorThree);
    }

    function test_areValidatorsOptedIn_returnsAnyAggregation_forBatch_onInitialSetup() public {
        bytes[] memory batchOfValidatorPublicKeys = new bytes[](3);
        batchOfValidatorPublicKeys[0] = validatorPublicKeyOne;   // true on A
        batchOfValidatorPublicKeys[1] = validatorPublicKeyTwo;   // true on B
        batchOfValidatorPublicKeys[2] = validatorPublicKeyThree; // false everywhere

        bool[] memory anyAggregationStatuses = hub.areValidatorsOptedIn(batchOfValidatorPublicKeys);

        bool[] memory expectedAnyAggregationStatuses = new bool[](3);
        expectedAnyAggregationStatuses[0] = true;  // pk1 → true (A)
        expectedAnyAggregationStatuses[1] = true;  // pk2 → true (B)
        expectedAnyAggregationStatuses[2] = false; // pk3 → false
        _assertEqualBoolArray(anyAggregationStatuses, expectedAnyAggregationStatuses);
    }

    function test_areValidatorsOptedIn_reflectsChanges_afterUpdateRegistryReplacement() public {
        // Replace index 1 (registryB) with registryC; pk2 was true on B, but false on C
        hub.updateRegistry(1, address(registryB), address(registryC));

        bytes[] memory batchOfValidatorPublicKeys = new bytes[](2);
        batchOfValidatorPublicKeys[0] = validatorPublicKeyOne; // still true on A
        batchOfValidatorPublicKeys[1] = validatorPublicKeyTwo; // should become false (B replaced)

        bool[] memory anyAggregationStatuses = hub.areValidatorsOptedIn(batchOfValidatorPublicKeys);

        bool[] memory expectedAnyAggregationStatuses = new bool[](2);
        expectedAnyAggregationStatuses[0] = true;  // pk1 true via A
        expectedAnyAggregationStatuses[1] = false; // pk2 now false (B→C)
        _assertEqualBoolArray(anyAggregationStatuses, expectedAnyAggregationStatuses);
    }

    function test_isValidatorOptedInList_expandsLength_afterAddingAdditionalRegistry() public {
        // Add a fourth registry (registryD) that returns false for all keys by default
        hub.addRegistry(address(registryD));

        // Query per-registry flags for validatorPublicKeyOne; should now have length 4
        bool[] memory perRegistryFlagsForValidatorOne = hub.isValidatorOptedInList(validatorPublicKeyOne);
        assertEq(perRegistryFlagsForValidatorOne.length, 4, "expected per-registry list to include the newly added registry");
        assertEq(perRegistryFlagsForValidatorOne[0], true);   // A
        assertEq(perRegistryFlagsForValidatorOne[1], false);  // B
        assertEq(perRegistryFlagsForValidatorOne[2], false);  // C
        assertEq(perRegistryFlagsForValidatorOne[3], false);  // D (newly added defaults to false)
    }

    function test_areValidatorsOptedIn_returnsEmpty_whenInputArrayIsEmpty() public {
        bytes[] memory emptyBatchOfValidatorPublicKeys = new bytes[](0);
        bool[] memory anyAggregationStatuses = hub.areValidatorsOptedIn(emptyBatchOfValidatorPublicKeys);
        assertEq(anyAggregationStatuses.length, 0, "empty input should produce empty result");
    }

    function test_isValidatorOptedInList_allFalse_forUnknownValidatorKey() public {
        bool[] memory perRegistryFlagsForUnknownKey = hub.isValidatorOptedInList(validatorPublicKeyThree);
        bool[] memory expectedFlagsForUnknownKey = new bool[](3);
        expectedFlagsForUnknownKey[0] = false;
        expectedFlagsForUnknownKey[1] = false;
        expectedFlagsForUnknownKey[2] = false;
        _assertEqualBoolArray(perRegistryFlagsForUnknownKey, expectedFlagsForUnknownKey);
    }

}


// ---------------------
// Local mock registry
// ---------------------
contract MockRegistry is IRegistry {
    mapping(bytes32 => bool) private optedInByKeyHash;

    function setOptIn(bytes memory validatorPublicKey, bool isOptedIn) external {
        optedInByKeyHash[keccak256(validatorPublicKey)] = isOptedIn;
    }

    function isValidatorOptedIn(bytes calldata validatorPublicKey)
        external
        view
        returns (bool)
    {
        return optedInByKeyHash[keccak256(validatorPublicKey)];
    }
}

contract MockWrongSigRegistry {
    // Wrong arg type → selector mismatch with IRegistry.isValidatorOptedIn(bytes)
    function isValidatorOptedIn(bytes32 /*validatorPublicKey*/) external pure returns (bool) {
        return true;
    }
}

contract MockMutatingRegistry {
    uint256 internal writes;
    // Not marked view and performs SSTORE → will revert under STATICCALL during validation probe
    function isValidatorOptedIn(bytes calldata /*validatorPublicKey*/) external returns (bool) {
        writes += 1;
        return true;
    }
}