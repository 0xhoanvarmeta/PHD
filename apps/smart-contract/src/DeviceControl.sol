// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

/**
 * @title DeviceControl
 * @notice Smart contract for managing device control commands
 * @dev Emits events when admin triggers commands, clients query to get commands
 */
contract DeviceControl {
    // Events
    event CommandTriggered(
        uint256 indexed commandId,
        uint256 timestamp,
        CommandType commandType,
        string backendCommandId
    );

    event AdminUpdated(address indexed oldAdmin, address indexed newAdmin);

    // Enums
    enum CommandType {
        SCRIPT,      // Execute a script directly
        URL          // Fetch from URL and execute
    }

    // Structs
    struct Command {
        uint256 id;
        CommandType commandType;
        string data;           // Script content or URL
        uint256 timestamp;
        address triggeredBy;
        string backendCommandId; // Backend command identifier
    }

    // State variables
    address public admin;
    uint256 public currentCommandId;
    mapping(uint256 => Command) public commands;

    // Modifiers
    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can call this function");
        _;
    }

    // Constructor
    constructor() {
        admin = msg.sender;
        currentCommandId = 0;
    }

    /**
     * @notice Trigger a new command (emit event for clients to listen)
     * @param commandType Type of command (SCRIPT or URL)
     * @param data Script content or URL
     * @param backendCommandId Backend command identifier for tracking
     */
    function Trigger(CommandType commandType, string calldata data, string calldata backendCommandId)
        external
        onlyAdmin
    {
        currentCommandId++;

        commands[currentCommandId] = Command({
            id: currentCommandId,
            commandType: commandType,
            data: data,
            timestamp: block.timestamp,
            triggeredBy: msg.sender,
            backendCommandId: backendCommandId
        });

        emit CommandTriggered(
            currentCommandId,
            block.timestamp,
            commandType,
            backendCommandId
        );
    }

    /**
     * @notice Get the current command details
     * @return id Command ID
     * @return commandType Type of command
     * @return data Script or URL
     * @return timestamp When command was triggered
     * @return backendCommandId Backend command identifier
     */
    function GetFunction()
        external
        view
        returns (
            uint256 id,
            CommandType commandType,
            string memory data,
            uint256 timestamp,
            string memory backendCommandId
        )
    {
        require(currentCommandId > 0, "No command available");

        Command memory cmd = commands[currentCommandId];
        return (cmd.id, cmd.commandType, cmd.data, cmd.timestamp, cmd.backendCommandId);
    }

    /**
     * @notice Get a specific command by ID
     * @param commandId The ID of the command to retrieve
     */
    function getCommand(uint256 commandId)
        external
        view
        returns (
            uint256 id,
            CommandType commandType,
            string memory data,
            uint256 timestamp,
            address triggeredBy,
            string memory backendCommandId
        )
    {
        require(commandId > 0 && commandId <= currentCommandId, "Invalid command ID");

        Command memory cmd = commands[commandId];
        return (
            cmd.id,
            cmd.commandType,
            cmd.data,
            cmd.timestamp,
            cmd.triggeredBy,
            cmd.backendCommandId
        );
    }

    /**
     * @notice Transfer admin role to a new address
     * @param newAdmin Address of the new admin
     */
    function transferAdmin(address newAdmin) external onlyAdmin {
        require(newAdmin != address(0), "Invalid address");
        address oldAdmin = admin;
        admin = newAdmin;
        emit AdminUpdated(oldAdmin, newAdmin);
    }

    /**
     * @notice Get the latest command ID
     */
    function getLatestCommandId() external view returns (uint256) {
        return currentCommandId;
    }
}
