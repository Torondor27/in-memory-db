# in-memory-db

The storage for the in-memory database is a simple map, chosen for its speed and simplicity in managing key-value storage.

For the transaction implementation, I've opted to store transactions as a stack containing operation stacks. 
This method simplifies handling nested transactions - retrieving the last transaction becomes straightforward.
Each transaction is represented as a stack containing all the operations performed within it. Consequently, when it 
comes to committing, everything is already completed, and we simply remove the transaction from the stack. Regarding rollback,
we get all operations from the operation stack and perform the opposite (delete - create, update - update back to old value).