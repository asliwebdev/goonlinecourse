-- transaction - single unit of work

-- transaction  is something transferred. An exchange or transfer of goods, services or funds.

-- ACID(principle) - kislota

-- atomicity - all or nothing rule
-- consistent - data integrity
-- isolation - isolation of transactions
-- durability - data is not lost if system fails


-- SQL, PostgreSQL supports transactions fully

-- transaction commands
-- START TRANSACTION - begin transaction
-- COMMIT - commit transaction
-- ROLLBACK - rollback transaction
-- SAVEPOINT - savepoint
-- RELEASE SAVEPOINT - release savepoint
-- ROLLBACK TO SAVEPOINT - rollback to savepoint

start transaction;

update car set is_new = false where id = 5;

commit;