// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Test, console} from "forge-std/Test.sol";
import {DeviceControl} from "../src/DeviceControl.sol";

contract DeviceControlTest is Test {
    DeviceControl public deviceControl;
    address public admin;
    address public user;

    event CommandTriggered(
        uint256 indexed commandId,
        uint256 timestamp,
        DeviceControl.CommandType commandType
    );

    function setUp() public {
        admin = address(this);
        user = address(0x1);
        deviceControl = new DeviceControl();
    }

    function test_InitialState() public view {
        assertEq(deviceControl.admin(), admin);
        assertEq(deviceControl.currentCommandId(), 0);
    }

    function test_TriggerScript() public {
        string memory script = "set-wallpaper https://example.com/image.jpg";

        vm.expectEmit(true, false, false, true);
        emit CommandTriggered(1, block.timestamp, DeviceControl.CommandType.SCRIPT);

        deviceControl.Trigger(DeviceControl.CommandType.SCRIPT, script);

        assertEq(deviceControl.currentCommandId(), 1);
    }

    function test_TriggerURL() public {
        string memory url = "https://api.example.com/commands/wallpaper";

        vm.expectEmit(true, false, false, true);
        emit CommandTriggered(1, block.timestamp, DeviceControl.CommandType.URL);

        deviceControl.Trigger(DeviceControl.CommandType.URL, url);

        assertEq(deviceControl.currentCommandId(), 1);
    }

    function test_GetFunction() public {
        string memory testData = "test command";
        deviceControl.Trigger(DeviceControl.CommandType.SCRIPT, testData);

        (
            uint256 id,
            DeviceControl.CommandType commandType,
            string memory data,
            uint256 timestamp
        ) = deviceControl.GetFunction();

        assertEq(id, 1);
        assertEq(uint256(commandType), uint256(DeviceControl.CommandType.SCRIPT));
        assertEq(data, testData);
        assertEq(timestamp, block.timestamp);
    }

    function test_GetFunctionRevertsWhenNoCommand() public {
        vm.expectRevert("No command available");
        deviceControl.GetFunction();
    }

    function test_GetCommand() public {
        string memory testData = "test command";
        deviceControl.Trigger(DeviceControl.CommandType.URL, testData);

        (
            uint256 id,
            DeviceControl.CommandType commandType,
            string memory data,
            uint256 timestamp,
            address triggeredBy
        ) = deviceControl.getCommand(1);

        assertEq(id, 1);
        assertEq(uint256(commandType), uint256(DeviceControl.CommandType.URL));
        assertEq(data, testData);
        assertEq(timestamp, block.timestamp);
        assertEq(triggeredBy, admin);
    }

    function test_OnlyAdminCanTrigger() public {
        vm.prank(user);
        vm.expectRevert("Only admin can call this function");
        deviceControl.Trigger(DeviceControl.CommandType.SCRIPT, "test");
    }

    function test_TransferAdmin() public {
        address newAdmin = address(0x2);

        deviceControl.transferAdmin(newAdmin);

        assertEq(deviceControl.admin(), newAdmin);
    }

    function test_TransferAdminOnlyByAdmin() public {
        vm.prank(user);
        vm.expectRevert("Only admin can call this function");
        deviceControl.transferAdmin(user);
    }

    function test_TransferAdminRevertsZeroAddress() public {
        vm.expectRevert("Invalid address");
        deviceControl.transferAdmin(address(0));
    }

    function test_MultipleCommands() public {
        deviceControl.Trigger(DeviceControl.CommandType.SCRIPT, "command 1");
        deviceControl.Trigger(DeviceControl.CommandType.URL, "command 2");
        deviceControl.Trigger(DeviceControl.CommandType.SCRIPT, "command 3");

        assertEq(deviceControl.currentCommandId(), 3);

        (, , string memory data, ) = deviceControl.GetFunction();
        assertEq(data, "command 3"); // Should return latest command
    }

    function test_GetLatestCommandId() public {
        assertEq(deviceControl.getLatestCommandId(), 0);

        deviceControl.Trigger(DeviceControl.CommandType.SCRIPT, "test");
        assertEq(deviceControl.getLatestCommandId(), 1);

        deviceControl.Trigger(DeviceControl.CommandType.URL, "test2");
        assertEq(deviceControl.getLatestCommandId(), 2);
    }
}
