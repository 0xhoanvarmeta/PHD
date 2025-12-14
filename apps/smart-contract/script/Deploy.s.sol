// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {DeviceControl} from "../src/DeviceControl.sol";

contract DeployScript is Script {
    function setUp() public {}

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerPrivateKey);

        DeviceControl deviceControl = new DeviceControl();

        console.log("DeviceControl deployed at:", address(deviceControl));
        console.log("Admin address:", deviceControl.admin());

        vm.stopBroadcast();
    }
}
