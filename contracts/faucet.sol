// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Faucet {
    event Withdrawal(address indexed to, uint256 amount);
    event Deposit(address indexed from, uint256 amount);

    function withdraw(uint256 withdraw_amount) public {
        require(withdraw_amount <= 1 ether, "You can withdraw up to 1 ETH");

        require(address(this).balance >= withdraw_amount, "Insufficient balance in faucet");

        payable(msg.sender).transfer(withdraw_amount);
        emit Withdrawal(msg.sender, withdraw_amount);
    }

    receive() external payable {
        emit Deposit(msg.sender, msg.value);
    }
}