- this project is based on udemy course: https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/
- this project uses ``sqlc`` for db communication & ``migrate`` for <br>
db migrations

### psql locks
- documentation: https://www.postgresql.org/docs/current/explicit-locking.html
- the link for the queries: https://wiki.postgresql.org/wiki/Lock_Monitoring
- to show blocked queries:
```sql
SELECT blocked_locks.pid     AS blocked_pid,
         blocked_activity.usename  AS blocked_user,
         blocking_locks.pid     AS blocking_pid,
         blocking_activity.usename AS blocking_user,
         blocked_activity.query    AS blocked_statement,
         blocking_activity.query   AS current_statement_in_blocking_process
   FROM  pg_catalog.pg_locks         blocked_locks
    JOIN pg_catalog.pg_stat_activity blocked_activity  ON blocked_activity.pid = blocked_locks.pid
    JOIN pg_catalog.pg_locks         blocking_locks 
        ON blocking_locks.locktype = blocked_locks.locktype
        AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
        AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
        AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
        AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
        AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
        AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
        AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
        AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
        AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
        AND blocking_locks.pid != blocked_locks.pid

    JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
   WHERE NOT blocked_locks.granted;
```
- get the list of locks in db:
```sql
SELECT 
--     a.datname, //shows the database name
    a.application_name, -- from what the application lock comes from
    l.relation::regclass, -- name of the table
    l.transactionid,
    l.mode,
    l.locktype,
    l.GRANTED,
    a.usename,
    a.query,
a.pid
FROM pg_stat_activity a
JOIN pg_locks l ON l.pid = a.pid
where a.application_name = 'psql' -- your app name can be different
ORDER BY a.pid;
```