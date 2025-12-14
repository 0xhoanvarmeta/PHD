// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {DeviceControl} from "../src/DeviceControl.sol";

/**
 * @title Verify Script
 * @notice Script to verify DeviceControl contract on block explorer
 *
 * Usage:
 * forge script script/Verify.s.sol --rpc-url <rpc_url> --verify --etherscan-api-key <api_key>
 *
 * Or use the helper command:
 * pnpm contract:verify <contract_address>
 */
contract VerifyScript is Script {
    function setUp() public {}

    function run() public view {
        address contractAddress = vm.envAddress("CONTRACT_ADDRESS");

        console.log("=== Contract Verification Info ===");
        console.log("Contract Address:", contractAddress);
        console.log("Contract Name: DeviceControl");
        console.log("Compiler Version: v0.8.28");
        console.log("Optimization: Enabled (200 runs)");
        console.log("");
        console.log("To verify manually, run:");
        console.log("forge verify-contract \\");
        console.log("  ", contractAddress, "\\");
        console.log("  src/DeviceControl.sol:DeviceControl \\");
        console.log("  --chain-id <chain_id> \\");
        console.log("  --verifier-url <verifier_url> \\");
        console.log("  --etherscan-api-key <api_key>");
        console.log("");
        console.log("For Hedera Testnet:");
        console.log("forge verify-contract \\");
        console.log("  ", contractAddress, "\\");
        console.log("  src/DeviceControl.sol:DeviceControl \\");
        console.log("  --chain-id 296 \\");
        console.log("  --verifier-url https://server-verify.hashscan.io \\");
        console.log("  --etherscan-api-key $ETHERSCAN_API_KEY");
    }
}
