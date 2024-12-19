## Deposit

Bridge ERC 20 tokens from rootchain to childchain via deposit.

```mermaid
sequenceDiagram
	User->>Network: deposit
	Network->>RootERC20.sol: approve(RootERC20Predicate)
	Network->>RootERC20Predicate.sol: deposit()
	RootERC20Predicate.sol->>RootERC20Predicate.sol: mapToken()
	RootERC20Predicate.sol->>StateSender.sol: syncState(MAP_TOKEN_SIG), recv=ChildERC20Predicate
	RootERC20Predicate.sol-->>Network: TokenMapped Event
	StateSender.sol-->>Network: StateSynced Event to map tokens on child predicate
	RootERC20Predicate.sol->>StateSender.sol: syncState(DEPOSIT_SIG), recv=ChildERC20Predicate
	StateSender.sol-->>Network: StateSynced Event to deposit on child chain
	Network->>User: ok
	Network->>StateReceiver.sol:commit()
	StateReceiver.sol-->>Network: NewCommitment Event
	Network->>StateReceiver.sol:execute()
	StateReceiver.sol->>ChildERC20Predicate.sol:onStateReceive()
	ChildERC20Predicate.sol->>ChildERC20.sol: mint()
	StateReceiver.sol-->>Network:StateSyncResult Event
```

## Withdraw

Bridge ERC 20 tokens from childchain to rootchain via withdrawal.

```mermaid
sequenceDiagram
	User->>Network: withdraw
	Network->>ChildERC20Predicate.sol: withdrawTo()
	ChildERC20Predicate.sol->>ChildERC20: burn()
	ChildERC20Predicate.sol->>L2StateSender.sol: syncState(WITHDRAW_SIG), recv=RootERC20Predicate
	Network->>User: tx hash
	User->>Network: get tx receipt
	Network->>User: exit event id
	ChildERC20Predicate.sol-->>Network: L2ERC20Withdraw Event
	L2StateSender.sol-->>Network: StateSynced Event
	Network->>Network: Seal block
	Network->>CheckpointManager.sol: submit()
```
## Exit

Finalize withdrawal of ERC 20 tokens from childchain to rootchain.

```mermaid
sequenceDiagram
	User->>Network: exit, event id:X
	Network->>Network: bridge_generateExitProof()
	Network->>CheckpointManager.sol: getCheckpointBlock()
	CheckpointManager.sol->>Network: blockNum
	Network->>Network: getExitEventsForProof(epochNum, blockNum)
	Network->>Network: createExitTree(exitEvents)
	Network->>Network: generateProof()
	Network->>ExitHelper.sol: exit()
	ExitHelper.sol->>CheckpointManager.sol: getEventMembershipByBlockNumber()
	ExitHelper.sol->>RootERC20Predicate.sol:onL2StateReceive()
	RootERC20Predicate.sol->>RootERC20: transfer()
	Network->>User: ok
	RootERC20Predicate.sol-->>Network: ERC20Withdraw Event
	ExitHelper.sol-->>Network: ExitProcessed Event
```

