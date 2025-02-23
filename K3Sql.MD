# 1. 基本常用SQL

```sql
/*绕开KSQL, KSQL因为要兼容SQL SERVER和ORACLE，相较于标准SQL，功能大幅度缩水，很多复杂的语法，句式、函数都不支持。
实际开发过程中，尤其是能明确无需兼容多种数据库的场景下，用标准SQL太香了，有时难以抗拒。
在BOS底层实现上，如果发现SQL起始位置有dialect，KSQL就不会再翻译了，直接执行当前脚本。
dialect一定要放在SQL的首字符位置，它之前不要有任何其它字符，包括空格换行。*/
/*dialect*/
-------------------------------------数据库DBID-------------------------------------
-- 当前执行的数据库ID
SELECT  DB_ID()
-- 查询DBID
SELECT  DB_ID('DWQueue');
-------------------------------------字符串分割函数并是否存在重复-------------------------------------
--分割函数
ALTER FUNCTION [dbo].[f_splitStr_Tb]
(
@str NVARCHAR(500),--要分割的字符串
@splitStr NVARCHAR(20)--进行分割的符号
)RETURNS @table TABLE(id INT,val NVARCHAR(50))
AS
BEGIN 
    DECLARE @index INT,@startsplit INT,@id INT 
    SELECT @index=CHARINDEX(@splitStr,@str),@startsplit=1,@id=1
    WHILE @index>0
    BEGIN
        IF @id>1  
        BEGIN 
            SELECT @startsplit=@index+LEN(@splitStr) 
            SELECT @index=CHARINDEX(@splitStr,@str,@startsplit)
  
        END 
        IF @index>0  
            BEGIN 
                INSERT INTO @table VALUES (@id,SUBSTRING(@str,@startsplit,@index-@startsplit)) 
            END 
        ELSE 
        BEGIN  
            INSERT INTO @table VALUES (@id,SUBSTRING(@str,@startsplit,LEN(@str)-@startsplit+1))
        END 
        SELECT @id=@id+1
    END
    RETURN 
END
--分割中的字符是否存在重复
select t0.* from T_BD_MATERIALGROUP t0
where exists (select count(1)　from dbo.[f_ splitStr_Tb](t0.FFULLPARENTID,'.') group by val having count(1)>1);

-- 根据一个表来更新另外一个表，msSql和oracle都适应
MERGE INTO target_table
USING source_table
ON condition
WHEN MATCHED THEN 
operation
WHEN NOT MATCHED THEN 
operation;
```

# 2. 数据中心SQL

## 2.1 数据中心版本

```sql
-------------------------------------数据中心-------------------------------------
-- 注意：查询数据中心，必须连接管理中心数据库查询
SELECT b.FNAME AS [数据中心(账套名称)],a.FDATACENTERID AS [数据中心内码(DBID)],a.FVISION AS 产品版本号,a.* 
FROM T_BAS_DATACENTER a 
LEFT JOIN T_BAS_DATACENTER_L b ON a.FDATACENTERID=b.FDATACENTERID AND b.FLOCALEID=2052;

-- 重置管理中心的密码为6个8，注意选择执行的数据库为管理中心
update t_sec_user set fpassword='9EF0664B011C3DC6305C8B213378BACC2451FAEA',fmemo='8d2fea58-cda1-4035-9024-56cc6df8957c' where fuserid=16394

-- 管理中心版本
select FVersion as 版本,* from T_bos_installedpackage where FPKGID='K3Cloud_DC_VERSION'

-- 数据中心列表
select * from T_BAS_DATACENTER
```

## 2.2 性能分析

### 2.2.1 查等待

```sql
--查等待--
WITH [Waits] 
AS (SELECT wait_type, wait_time_ms/ 1000.0 AS [WaitS],
          (wait_time_ms - signal_wait_time_ms) / 1000.0 AS [ResourceS],
           signal_wait_time_ms / 1000.0 AS [SignalS],
           waiting_tasks_count AS [WaitCount],
           100.0 *  wait_time_ms / SUM (wait_time_ms) OVER() AS [Percentage],
           ROW_NUMBER() OVER(ORDER BY wait_time_ms DESC) AS [RowNum]
    FROM sys.dm_os_wait_stats WITH (NOLOCK)
    WHERE [wait_type] NOT IN (
        N'BROKER_EVENTHANDLER', N'BROKER_RECEIVE_WAITFOR', N'BROKER_TASK_STOP',
        N'BROKER_TO_FLUSH', N'BROKER_TRANSMITTER', N'CHECKPOINT_QUEUE',
        N'CHKPT', N'CLR_AUTO_EVENT', N'CLR_MANUAL_EVENT', N'CLR_SEMAPHORE',
        N'DBMIRROR_DBM_EVENT', N'DBMIRROR_EVENTS_QUEUE', N'DBMIRROR_WORKER_QUEUE',
        N'DBMIRRORING_CMD', N'DIRTY_PAGE_POLL', N'DISPATCHER_QUEUE_SEMAPHORE',
        N'EXECSYNC', N'FSAGENT', N'FT_IFTS_SCHEDULER_IDLE_WAIT', N'FT_IFTSHC_MUTEX',
        N'HADR_CLUSAPI_CALL', N'HADR_FILESTREAM_IOMGR_IOCOMPLETION', N'HADR_LOGCAPTURE_WAIT', 
        N'HADR_NOTIFICATION_DEQUEUE', N'HADR_TIMER_TASK', N'HADR_WORK_QUEUE',
        N'KSOURCE_WAKEUP', N'LAZYWRITER_SLEEP', N'LOGMGR_QUEUE', 
        N'MEMORY_ALLOCATION_EXT', N'ONDEMAND_TASK_QUEUE',
        N'PREEMPTIVE_OS_LIBRARYOPS', N'PREEMPTIVE_OS_COMOPS', N'PREEMPTIVE_OS_CRYPTOPS',
        N'PREEMPTIVE_OS_PIPEOPS', N'PREEMPTIVE_OS_AUTHENTICATIONOPS',
        N'PREEMPTIVE_OS_GENERICOPS', N'PREEMPTIVE_OS_VERIFYTRUST',
        N'PREEMPTIVE_OS_FILEOPS', N'PREEMPTIVE_OS_DEVICEOPS',
        N'PWAIT_ALL_COMPONENTS_INITIALIZED', N'QDS_PERSIST_TASK_MAIN_LOOP_SLEEP',
        N'QDS_ASYNC_QUEUE',
        N'QDS_CLEANUP_STALE_QUERIES_TASK_MAIN_LOOP_SLEEP', N'REQUEST_FOR_DEADLOCK_SEARCH',
        N'RESOURCE_QUEUE', N'SERVER_IDLE_CHECK', N'SLEEP_BPOOL_FLUSH', N'SLEEP_DBSTARTUP',
        N'SLEEP_DCOMSTARTUP', N'SLEEP_MASTERDBREADY', N'SLEEP_MASTERMDREADY',
        N'SLEEP_MASTERUPGRADED', N'SLEEP_MSDBSTARTUP', N'SLEEP_SYSTEMTASK', N'SLEEP_TASK',
        N'SLEEP_TEMPDBSTARTUP', N'SNI_HTTP_ACCEPT', N'SP_SERVER_DIAGNOSTICS_SLEEP',
        N'SQLTRACE_BUFFER_FLUSH', N'SQLTRACE_INCREMENTAL_FLUSH_SLEEP', N'SQLTRACE_WAIT_ENTRIES',
        N'WAIT_FOR_RESULTS', N'WAITFOR', N'WAITFOR_TASKSHUTDOWN', N'WAIT_XTP_HOST_WAIT',
        N'WAIT_XTP_OFFLINE_CKPT_NEW_LOG', N'WAIT_XTP_CKPT_CLOSE', N'XE_DISPATCHER_JOIN',
        N'XE_DISPATCHER_WAIT', N'XE_LIVE_TARGET_TVF', N'XE_TIMER_EVENT')
    AND waiting_tasks_count > 0)
SELECT
    MAX (W1.wait_type) AS [WaitType],
    CAST (MAX (W1.WaitS) AS DECIMAL (16,2)) AS [Wait_Sec],
    CAST (MAX (W1.ResourceS) AS DECIMAL (16,2)) AS [Resource_Sec],
    CAST (MAX (W1.SignalS) AS DECIMAL (16,2)) AS [Signal_Sec],
    MAX (W1.WaitCount) AS [Wait Count],
    CAST (MAX (W1.Percentage) AS DECIMAL (5,2)) AS [Wait Percentage],
    CAST ((MAX (W1.WaitS) / MAX (W1.WaitCount)) AS DECIMAL (16,4)) AS [AvgWait_Sec],
    CAST ((MAX (W1.ResourceS) / MAX (W1.WaitCount)) AS DECIMAL (16,4)) AS [AvgRes_Sec],
    CAST ((MAX (W1.SignalS) / MAX (W1.WaitCount)) AS DECIMAL (16,4)) AS [AvgSig_Sec]
FROM Waits AS W1
INNER JOIN Waits AS W2
ON W2.RowNum <= W1.RowNum
GROUP BY W1.RowNum
HAVING SUM (W2.Percentage) - MAX (W1.Percentage) < 99 -- percentage threshold
OPTION (RECOMPILE);
```

### 2.2.2 查阻塞

```sql
--查阻塞--
SELECT t1.resource_type AS [锁类型], DB_NAME(resource_database_id) AS [数据库名],
t1.resource_associated_entity_id AS [阻塞资源对象],t1.resource_description as [资源描述信息], t1.request_mode AS [请求的锁], 
t1.request_session_id AS [等待会话], t2.wait_duration_ms AS [等待时间],     
(SELECT [text] FROM sys.dm_exec_requests AS r WITH (NOLOCK)                    
CROSS APPLY sys.dm_exec_sql_text(r.[sql_handle]) 
    WHERE r.session_id = t1.request_session_id
) AS [等待会话执行的批SQL],
(SELECT SUBSTRING(qt.[text],r.statement_start_offset/2, 
(CASE WHEN r.statement_end_offset = -1 
THEN LEN(CONVERT(nvarchar(max), qt.[text])) * 2 
ELSE r.statement_end_offset END )/2) 
FROM sys.dm_exec_requests AS r WITH (NOLOCK)
CROSS APPLY sys.dm_exec_sql_text(r.[sql_handle]) AS qt
WHERE r.session_id = t1.request_session_id
) AS [等待会话执行的SQL],                  
t2.blocking_session_id AS [阻塞会话],                                      
(SELECT [text] FROM sys.sysprocesses AS p                                     
CROSS APPLY sys.dm_exec_sql_text(p.[sql_handle]) 
WHERE p.spid = t2.blocking_session_id
) AS [阻塞会话执行的批SQL]
FROM sys.dm_tran_locks AS t1 WITH (NOLOCK)
INNER JOIN sys.dm_os_waiting_tasks AS t2 WITH (NOLOCK)
ON t1.lock_owner_address = t2.resource_address OPTION (RECOMPILE);
```

### 2.2.3  查IO

```sql
--查IO--
SELECT TOP(20) DB_NAME(t.[dbid]) AS [Database Name], LEFT(t.[text], 50) AS [Short Query Text],
 qs.total_logical_reads AS [Total Logical Reads],
 qs.min_logical_reads AS [Min Logical Reads],
 qs.total_logical_reads/qs.execution_count AS [Avg Logical Reads],
 qs.max_logical_reads AS [Max Logical Reads],   
 qs.min_worker_time AS [Min Worker Time],
 qs.total_worker_time/qs.execution_count AS [Avg Worker Time], 
 qs.max_worker_time AS [Max Worker Time], 
 qs.min_elapsed_time AS [Min Elapsed Time], 
 qs.total_elapsed_time/qs.execution_count AS [Avg Elapsed Time], 
 qs.max_elapsed_time AS [Max Elapsed Time],
 qs.execution_count AS [Execution Count], qs.creation_time AS [Creation Time]
 ,t.[text] AS [Complete Query Text], qp.query_plan AS [Query Plan] -- uncomment out these columns if not copying results to Excel
 FROM sys.dm_exec_query_stats AS qs WITH (NOLOCK)
 CROSS APPLY sys.dm_exec_sql_text(plan_handle) AS t 
 CROSS APPLY sys.dm_exec_query_plan(plan_handle) AS qp 
 ORDER BY qs.total_logical_reads DESC OPTION (RECOMPILE);
```

### 2.2.4 查CPU

```sql
--查CPU--
SELECT TOP(20) DB_NAME(t.[dbid]) AS [Database Name], LEFT(t.[text], 50) AS [Short Query Text],  
qs.total_worker_time AS [Total Worker Time], qs.min_worker_time AS [Min Worker Time],
qs.total_worker_time/qs.execution_count AS [Avg Worker Time], 
qs.max_worker_time AS [Max Worker Time], 
qs.min_elapsed_time AS [Min Elapsed Time], 
qs.total_elapsed_time/qs.execution_count AS [Avg Elapsed Time], 
qs.max_elapsed_time AS [Max Elapsed Time],
qs.min_logical_reads AS [Min Logical Reads],
qs.total_logical_reads/qs.execution_count AS [Avg Logical Reads],
qs.max_logical_reads AS [Max Logical Reads], 
qs.execution_count AS [Execution Count], qs.creation_time AS [Creation Time]
 ,t.[text] AS [Query Text], qp.query_plan AS [Query Plan] -- uncomment out these columns if not copying results to Excel
FROM sys.dm_exec_query_stats AS qs WITH (NOLOCK)
CROSS APPLY sys.dm_exec_sql_text(plan_handle) AS t 
CROSS APPLY sys.dm_exec_query_plan(plan_handle) AS qp 
ORDER BY qs.total_worker_time DESC OPTION (RECOMPILE);
```

### 2.2.5 CPU高时查询最耗CPU的SQL

```sql
DECLARE @sessionsCPU TABLE 
(
 session_id INT NOT NULL,
 cpu_time BIGINT NULL ,
 initialCPUFlag BIT NOT NULL
)
  
INSERT INTO @sessionsCPU
SELECT 
 session_id
,sum(cpu_time) as cpu_time
,1 as initialCPUFlag
FROM sys.dm_exec_requests
WHERE session_id > 20
GROUP BY session_id
  
WAITFOR DELAY '0:00:10' -- 等待10秒
  
INSERT INTO @sessionsCPU
SELECT
 session_id
,sum(cpu_time) as cpu_time
,0 as initialCPUFlag
FROM sys.dm_exec_requests
WHERE session_id > 20
GROUP BY session_id;
  
WITH total AS
(
select
 initialCPUFlag
,sum(cpu_time) as cpu_time
FROM @sessionsCPU
GROUP BY initialCPUFlag
), 
delta AS
(
 SELECT 
s.cpu_time - f.cpu_time as total_cpu
 FROM total f
 CROSS JOIN total s
 WHERE f.initialCPUFlag = 1
 AND s.initialCPUFlag = 0
)  
SELECT
 i.session_id
,convert(numeric(5,2), (100. * (((i.cpu_time - l.cpu_time) * 1.) / (d.total_cpu * 1.)))) AS percentCPU
,convert(char(8), getdate() - r.start_time, 108) AS run_duration
,substring( st.text, 
(r.statement_start_offset/2)+1,
((CASE r.statement_end_offset
WHEN -1 THEN datalength(st.text)
ELSE r.statement_end_offset END - r.statement_start_offset)/2) + 1)
AS statement_text
,st.text AS full_query
,pt.query_plan
,s.login_name
,s.host_name
,s.program_name
FROM @sessionsCPU i
INNER JOIN @sessionsCPU l
ON i.session_id = l.session_id
AND i.initialCPUFlag = 0 AND l.initialCPUFlag = 1
INNER JOIN sys.dm_exec_requests r
ON i.session_id = r.session_id
INNER JOIN sys.dm_exec_sessions s
ON i.session_id = s.session_id
CROSS JOIN delta d
CROSS APPLY sys.dm_exec_sql_text(r.sql_handle) AS st
cross apply sys.dm_exec_query_plan(r.plan_handle) as pt
WHERE d.total_cpu <> 0
ORDER BY percentCPU DESC;
```

### 2.2.6 死锁

```sql
--获取死锁信息 sql server 2008 r2
SELECT  XEventData.XEvent.value('@timestamp', 'datetime2(3)'),
cast(XEventData.XEvent.value('(data/value)[1]', 'varchar(max)') as xml)
FROM (SELECT CAST (target_data AS XML) AS TargetData
    FROM sys.dm_xe_session_targets st
    JOIN sys.dm_xe_sessions s ON s.address = st.event_session_address
    WHERE [name] = 'system_health' AND st.target_name = 'ring_buffer') AS Data
CROSS APPLY TargetData.nodes ('//RingBufferTarget/event') AS XEventData (XEvent)
    WHERE XEventData.XEvent.value('@name', 'varchar(4000)') = 'xml_deadlock_report';
--or
SELECT event.c.value('@timestamp', 'datetime2(3)'), 
       datavalue.c.value('(./text())[1]', 'nvarchar(max)') AS DeadlockGraph
FROM
(
SELECT CAST(target_data AS XML) AS TargetData
FROM sys.dm_xe_session_targets st
JOIN sys.dm_xe_sessions s ON s.address = st.event_session_address
WHERE name = 'system_health'
) AS DATA
CROSS APPLY TargetData.nodes('/RingBufferTarget/event') AS event(c)
CROSS APPLY event.c.nodes ('data/value') AS datavalue(c)
WHERE event.c.value('@name', 'varchar(4000)') = 'xml_deadlock_report'
--or
SELECT
xed.value('@timestamp', 'datetime2(3)') as CreationDate,
xed.query('.') AS XEvent
FROM
(
SELECT CAST([target_data] AS XML) AS TargetData
FROM sys.dm_xe_session_targets AS st
INNER JOIN sys.dm_xe_sessions AS s
ON s.address = st.event_session_address
WHERE s.name = N'system_health'
AND st.target_name = N'ring_buffer'
) AS Data
CROSS APPLY TargetData.nodes('RingBufferTarget/event[@name="xml_deadlock_report"]') AS XEventData (xed)
ORDER BY CreationDate DESC
SELECT  CAST(event_data.value('(event/data/value)[1]',
                               'varchar(max)') AS XML) AS DeadlockGraph
FROM    ( SELECT    XEvent.query('.') AS event_data
          FROM      (    -- Cast the target_data to XML 
                      SELECT    CAST(target_data AS XML) AS TargetData
                      FROM      sys.dm_xe_session_targets st
                                JOIN sys.dm_xe_sessions s
                                 ON s.address = st.event_session_address
                      WHERE     name = 'system_health'
                                AND target_name = 'ring_buffer'
                    ) AS Data -- Split out the Event Nodes 
                    CROSS APPLY TargetData.nodes('RingBufferTarget/
                                     event[@name="xml_deadlock_report"]')
                    AS XEventData ( XEvent )
        ) AS tab ( event_data )
-- Retrieve from Extended Events in 2012
SELECT  XEvent.query('(event/data/value/deadlock)[1]') AS DeadlockGraph
FROM    ( SELECT    XEvent.query('.') AS XEvent
          FROM      ( SELECT    CAST(target_data AS XML) AS TargetData
                      FROM      sys.dm_xe_session_targets st
                                JOIN sys.dm_xe_sessions s
                                 ON s.address = st.event_session_address
                      WHERE     s.name = 'system_health'
                                AND st.target_name = 'ring_buffer'
                    ) AS Data
                    CROSS APPLY TargetData.nodes
                  ('RingBufferTarget/event[@name="xml_deadlock_report"]')
                    AS XEventData ( XEvent )
        ) AS src;
SELECT XEvent.query('(event/data/value/deadlock)[1]') AS DeadlockGraph
, XEvent.value('(/event/@timestamp)[1]','DATETIMEOFFSET') AS BeginTime
    FROM ( SELECT XEvent.query('.') AS XEvent
        FROM ( SELECT CAST(target_data AS XML) AS TargetData
            FROM sys.dm_xe_session_targets st
                JOIN sys.dm_xe_sessions s
                ON s.address = st.event_session_address
            WHERE s.name = 'system_health'
                AND st.target_name = 'ring_buffer' 
              ) AS Data 
              CROSS APPLY 
                 TargetData.nodes 
                    ('RingBufferTarget/event[@name="xml_deadlock_report"]')
              AS XEventData ( XEvent ) 
      ) AS src
;
-------------------------------------------------
DECLARE @SessionName SysName 
SELECT @SessionName = 'system_health'
/* 
SELECT  Session_Name = s.name, s.blocked_event_fire_time, s.dropped_buffer_count, s.dropped_event_count, s.pending_buffers
FROM sys.dm_xe_session_targets t
INNER JOIN sys.dm_xe_sessions s ON s.address = t.event_session_address
WHERE target_name = 'event_file'
--*/
IF OBJECT_ID('tempdb..#Events') IS NOT NULL BEGIN
DROP TABLE #Events
END
DECLARE @Target_File NVarChar(1000)
, @Target_Dir NVarChar(1000)
, @Target_File_WildCard NVarChar(1000)
SELECT @Target_File = CAST(t.target_data as XML).value('EventFileTarget[1]/File[1]/@name', 'NVARCHAR(256)')
FROM sys.dm_xe_session_targets t
INNER JOIN sys.dm_xe_sessions s ON s.address = t.event_session_address
WHERE s.name = @SessionName
AND t.target_name = 'event_file'
SELECT @Target_Dir = LEFT(@Target_File, Len(@Target_File) - CHARINDEX('\', REVERSE(@Target_File))) 
SELECT @Target_File_WildCard = @Target_Dir + '\'  + @SessionName + '_*.xel'
--Keep this as a separate table because it's called twice in the next query.  You don't want this running twice.
SELECT DeadlockGraph = CAST(event_data AS XML)
, DeadlockID = Row_Number() OVER(ORDER BY file_name, file_offset)
INTO #Events
FROM sys.fn_xe_file_target_read_file(@Target_File_WildCard, null, null, null) AS F
WHERE event_data like '<event name="xml_deadlock_report%'
;WITH Victims AS
(
SELECT VictimID = Deadlock.Victims.value('@id', 'varchar(50)')
, e.DeadlockID 
FROM #Events e
CROSS APPLY e.DeadlockGraph.nodes('/event/data/value/deadlock/victim-list/victimProcess') as Deadlock(Victims)
)
, DeadlockObjects AS
(
SELECT DISTINCT e.DeadlockID
, ObjectName = Deadlock.Resources.value('@objectname', 'nvarchar(256)')
FROM #Events e
CROSS APPLY e.DeadlockGraph.nodes('/event/data/value/deadlock/resource-list/*') as Deadlock(Resources)
)
SELECT *
FROM
(
SELECT e.DeadlockID
, TransactionTime = Deadlock.Process.value('@lasttranstarted', 'datetime')
, DeadlockGraph
, DeadlockObjects = substring((SELECT (', ' + o.ObjectName)
FROM DeadlockObjects o
WHERE o.DeadlockID = e.DeadlockID
ORDER BY o.ObjectName
FOR XML PATH ('')
), 3, 4000)
, Victim = CASE WHEN v.VictimID IS NOT NULL 
THEN 1 
ELSE 0 
END
, SPID = Deadlock.Process.value('@spid', 'int')
, ProcedureName = Deadlock.Process.value('executionStack[1]/frame[1]/@procname[1]', 'varchar(200)')
, LockMode = Deadlock.Process.value('@lockMode', 'char(1)')
, Code = Deadlock.Process.value('executionStack[1]/frame[1]', 'varchar(1000)')
, ClientApp = CASE LEFT(Deadlock.Process.value('@clientapp', 'varchar(100)'), 29)
WHEN 'SQLAgent - TSQL JobStep (Job '
THEN 'SQLAgent Job: ' + (SELECT name FROM msdb..sysjobs sj WHERE substring(Deadlock.Process.value('@clientapp', 'varchar(100)'),32,32)=(substring(sys.fn_varbintohexstr(sj.job_id),3,100))) + ' - ' + SUBSTRING(Deadlock.Process.value('@clientapp', 'varchar(100)'), 67, len(Deadlock.Process.value('@clientapp', 'varchar(100)'))-67)
ELSE Deadlock.Process.value('@clientapp', 'varchar(100)')
END 
, HostName = Deadlock.Process.value('@hostname', 'varchar(20)')
, LoginName = Deadlock.Process.value('@loginname', 'varchar(20)')
, InputBuffer = Deadlock.Process.value('inputbuf[1]', 'varchar(1000)')
FROM #Events e
CROSS APPLY e.DeadlockGraph.nodes('/event/data/value/deadlock/process-list/process') as Deadlock(Process)
LEFT JOIN Victims v ON v.DeadlockID = e.DeadlockID AND v.VictimID = Deadlock.Process.value('@id', 'varchar(50)')
) X --In a subquery to make filtering easier (use column names, not XML parsing), no other reason
ORDER BY DeadlockID DESC
CREATE TABLE #errorlog (
            LogDate DATETIME 
            , ProcessInfo VARCHAR(100)
            , [Text] VARCHAR(MAX)
            );
DECLARE @tag VARCHAR (MAX) , @path VARCHAR(MAX);
INSERT INTO #errorlog EXEC sp_readerrorlog;
SELECT @tag = text
FROM #errorlog 
WHERE [Text] LIKE 'Logging%MSSQL\Log%';
DROP TABLE #errorlog;
SET @path = SUBSTRING(@tag, 38, CHARINDEX('MSSQL\Log', @tag) - 29);
SELECT 
  CONVERT(xml, event_data).query('/event/data/value/child::*') AS DeadlockReport,
  switchoffset(CONVERT(xml, event_data).value('(event[@name="xml_deadlock_report"]/@timestamp)[1]', 'datetime2'),'+08:00')
  AS Execution_Time
FROM sys.fn_xe_file_target_read_file(@path + '\system_health*.xel', NULL, NULL, NULL)
WHERE OBJECT_NAME like 'xml_deadlock_report';
--死锁分析--
DECLARE @SessionName SysName 
SELECT @SessionName = 'system_health'
/* 
SELECT  Session_Name = s.name, s.blocked_event_fire_time, s.dropped_buffer_count, s.dropped_event_count, s.pending_buffers
FROM sys.dm_xe_session_targets t
INNER JOIN sys.dm_xe_sessions s ON s.address = t.event_session_address
WHERE target_name = 'event_file'
--*/
IF OBJECT_ID('tempdb..#Events') IS NOT NULL BEGIN
DROP TABLE #Events
END
DECLARE @Target_File NVarChar(1000)
, @Target_Dir NVarChar(1000)
, @Target_File_WildCard NVarChar(1000)
SELECT @Target_File = CAST(t.target_data as XML).value('EventFileTarget[1]/File[1]/@name', 'NVARCHAR(256)')
FROM sys.dm_xe_session_targets t
INNER JOIN sys.dm_xe_sessions s ON s.address = t.event_session_address
WHERE s.name = @SessionName
AND t.target_name = 'event_file'
SELECT @Target_Dir = LEFT(@Target_File, Len(@Target_File) - CHARINDEX('\', REVERSE(@Target_File))) 
SELECT @Target_File_WildCard = @Target_Dir + '\'  + @SessionName + '_*.xel'
--Keep this as a separate table because it's called twice in the next query.  You don't want this running twice.
SELECT DeadlockGraph = CAST(event_data AS XML)
, DeadlockID = Row_Number() OVER(ORDER BY file_name, file_offset)
INTO #Events
FROM sys.fn_xe_file_target_read_file(@Target_File_WildCard, null, null, null) AS F
WHERE event_data like '<event name="xml_deadlock_report%'
;WITH Victims AS
(
SELECT VictimID = Deadlock.Victims.value('@id', 'varchar(50)')
, e.DeadlockID 
FROM #Events e
CROSS APPLY e.DeadlockGraph.nodes('/event/data/value/deadlock/victim-list/victimProcess') as Deadlock(Victims)
)
, DeadlockObjects AS
(
SELECT DISTINCT e.DeadlockID
, ObjectName = Deadlock.Resources.value('@objectname', 'nvarchar(256)')
FROM #Events e
CROSS APPLY e.DeadlockGraph.nodes('/event/data/value/deadlock/resource-list/*') as Deadlock(Resources)
)
SELECT *
FROM
(
SELECT e.DeadlockID
, TransactionTime = Deadlock.Process.value('@lasttranstarted', 'datetime')
, DeadlockGraph
, DeadlockObjects = substring((SELECT (', ' + o.ObjectName)
FROM DeadlockObjects o
WHERE o.DeadlockID = e.DeadlockID
ORDER BY o.ObjectName
FOR XML PATH ('')
), 3, 4000)
, Victim = CASE WHEN v.VictimID IS NOT NULL 
THEN 1 
ELSE 0 
END
, SPID = Deadlock.Process.value('@spid', 'int')
, ProcedureName = Deadlock.Process.value('executionStack[1]/frame[1]/@procname[1]', 'varchar(200)')
, LockMode = Deadlock.Process.value('@lockMode', 'char(1)')
, Code = Deadlock.Process.value('executionStack[1]/frame[1]', 'varchar(1000)')
, ClientApp = CASE LEFT(Deadlock.Process.value('@clientapp', 'varchar(100)'), 29)
WHEN 'SQLAgent - TSQL JobStep (Job '
THEN 'SQLAgent Job: ' + (SELECT name FROM msdb..sysjobs sj WHERE substring(Deadlock.Process.value('@clientapp', 'varchar(100)'),32,32)=(substring(sys.fn_varbintohexstr(sj.job_id),3,100))) + ' - ' + SUBSTRING(Deadlock.Process.value('@clientapp', 'varchar(100)'), 67, len(Deadlock.Process.value('@clientapp', 'varchar(100)'))-67)
ELSE Deadlock.Process.value('@clientapp', 'varchar(100)')
END 
, HostName = Deadlock.Process.value('@hostname', 'varchar(20)')
, LoginName = Deadlock.Process.value('@loginname', 'varchar(20)')
, InputBuffer = Deadlock.Process.value('inputbuf[1]', 'varchar(1000)')
FROM #Events e
CROSS APPLY e.DeadlockGraph.nodes('/event/data/value/deadlock/process-list/process') as Deadlock(Process)
LEFT JOIN Victims v ON v.DeadlockID = e.DeadlockID AND v.VictimID = Deadlock.Process.value('@id', 'varchar(50)')
) X --In a subquery to make filtering easier (use column names, not XML parsing), no other reason
ORDER BY DeadlockID DESC
CREATE TABLE #errorlog (
            LogDate DATETIME 
            , ProcessInfo VARCHAR(100)
            , [Text] VARCHAR(MAX)
            );
DECLARE @tag VARCHAR (MAX) , @path VARCHAR(MAX);
INSERT INTO #errorlog EXEC sp_readerrorlog;
SELECT @tag = text
FROM #errorlog 
WHERE [Text] LIKE 'Logging%MSSQL\Log%';
DROP TABLE #errorlog;
SET @path = SUBSTRING(@tag, 38, CHARINDEX('MSSQL\Log', @tag) - 29);
SELECT 
  CONVERT(xml, event_data).query('/event/data/value/child::*') AS DeadlockReport,
  switchoffset(CONVERT(xml, event_data).value('(event[@name="xml_deadlock_report"]/@timestamp)[1]', 'datetime2'),'+08:00')
  AS Execution_Time
FROM sys.fn_xe_file_target_read_file(@path + '\system_health*.xel', NULL, NULL, NULL)
WHERE OBJECT_NAME like 'xml_deadlock_report';
```

# 3. 元数据

## 3.1 常用业务对象查询SQL

```sql
-------------------------------------菜单-------------------------------------
--业务领域  基础管理 60
 select * from T_META_TOPCLASS_L where FLOCALEID=2052 and FNAME='基础管理'
--子系统 基础资料  6002
select t1.FNAME,* from T_META_SUBSYSTEM t0 
inner join T_META_SUBSYSTEM_L t1 on t1.FID = t0.FID
where  t1.FLOCALEID=2052 and t0.FTOPCLASSID='60';
---得到所有的菜单
select t1.FNAME, * from T_META_OBJECTTYPE t0
inner join T_META_OBJECTTYPE_L t1 on t1.FID = t0.FID
where   (FModelTypeID = 100 AND t0.FID IN 
(SELECT FOBJECTID FROM T_META_CONSOLEDETAIL WHERE FSTATUS = '1') or 
(FModelTypeID = 400 AND FDEVTYPE<>2))
 and t0.FSUBSYSID='6002' and t1.FLOCALEID=2052;
-------------------------------------查询单据关联的并发布业务流程图-------------------------------------
SELECT DISTINCT A.FID, A.FObjectTypeId AS FThisForm, B.FObjectTypeId AS FEntryForm, C.FBeginDate, C.FEndDate 
FROM T_BF_DEFVERSIONLOOKUP A 
INNER JOIN T_BF_DEFVERSION B ON (A.FID = B.FID) 
LEFT JOIN T_BF_PUBLISHENTRY C ON (A.FID = C.FBFVERSIONID) WHERE (A.FObjectTypeId = 'PAEZ_BillC');
-------------------------------------子系统-------------------------------------
SELECT b.FNAME AS 子系统名称,a.FID AS 子系统内码,a.* 
FROM T_META_SUBSYSTEM a 
LEFT JOIN T_META_SUBSYSTEM_L b ON a.FID=b.FID AND b.FLOCALEID=2052
-------------------------------------领域模型-------------------------------------
SELECT b.FNAME AS 领域模型名称,a.FID AS 领域模型内码,a.* 
FROM T_MDL_DOMAINMODETYPE a 
LEFT JOIN T_MDL_DOMAINMODETYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052
-------------------------------------业务对象-------------------------------------
SELECT b.FNAME AS 业务对象名称,a.FID AS [业务对象标识(FormId)],* 
FROM T_META_OBJECTTYPE a 
LEFT JOIN T_META_OBJECTTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
WHERE a.FID='SAL_OUTSTOCK'
-------------------------------------单据视图-------------------------------------
SELECT c.FNAME AS 业务对象名称,a.FDEPENDENCYOBJECTID AS [业务对象内码 (FormId)],b.FNAME AS 视图名称,a.FID AS [视图内码(layoutid)],* 
FROM T_META_OBJECTTYPEVIEW a 
LEFT JOIN T_META_OBJECTTYPEVIEW_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
LEFT JOIN T_META_OBJECTTYPE_L c ON a.FDEPENDENCYOBJECTID=c.FID AND  c.FLOCALEID=2052 
WHERE a.FDEPENDENCYOBJECTID IN ('SAL_OUTSTOCK','DE_CD_OutStock','00f0614c-d928-46b6-8e79-d1ac94ae867a')
-------------------------------------表单操作-------------------------------------
SELECT b.FNAME AS 操作名称,a.FID AS 操作内码,a.* 
FROM T_MDL_FORMOPERATIONTYPE a 
LEFT JOIN T_MDL_FORMOPERATIONTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
ORDER BY a.FID
-------------------------------------表单服务-------------------------------------
SELECT b.FDESC AS 服务名称,a.FACTIONID AS 服务内码,a.* 
FROM T_MDL_FORMBUSINESS a 
LEFT JOIN T_MDL_FORMBUSINESS_L b ON a.FACTIONID=b.FACTIONID AND b.FLOCALEID=2052
-------------------------------------操作内置校验器-------------------------------------
SELECT b.FNAME AS 校验器名称,a.* FROM T_MDL_FORMVALIDATIONTYPE a 
LEFT JOIN T_MDL_FORMVALIDATIONTYPE_L b ON a.FTYPEID=b.FTYPEID AND b.FLOCALEID=2052
-------------------------------------过滤方案-------------------------------------
-- 过滤方案
SELECT b.FDESCRIPTION AS 过滤方案名称,a.FSCHEMEID AS 过滤方案内码,c.FNAME AS 业务对象名称,d.FNAME AS 用户名,* 
FROM T_BAS_FILTERSCHEME a 
LEFT JOIN T_BAS_FILTERSCHEME_L b ON a.FSCHEMEID=b.FSCHEMEID AND b.FLOCALEID=2052 
LEFT JOIN T_META_OBJECTTYPE_L c ON a.FFORMID=c.FID AND c.FLOCALEID=2052 
LEFT JOIN T_SEC_USER d ON a.FUSERID=d.FUSERID 
WHERE a.FFORMID='PUR_PurchaseOrder'
-- 过滤符
SELECT b.FNAME AS 过滤符名称,a.FCONDITIONTYPEID AS 符号集编码,a.* FROM T_MDL_FILTERCONDITION a 
LEFT JOIN T_MDL_FILTERCONDITION_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
ORDER BY a.FCONDITIONTYPEID
-------------------------------------语言-------------------------------------
SELECT a.FMSNUMBER AS [语言ID(LCID)],b.FNAME AS 语言名称,a.* 
FROM T_BD_LANGUAGE a 
LEFT JOIN T_BD_LANGUAGE_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
ORDER BY a.FMSNUMBER
-------------------------------------编码规则-------------------------------------
-- 编码规则
SELECT b.FNAME AS [编码规则名称],a.FRULEID AS 编码规则内码,c.FNAME AS 业务对象名称,a.* 
FROM T_BAS_BILLCODERULE a 
LEFT JOIN T_BAS_BILLCODERULE_L b ON a.FRULEID=b.FRULEID AND b.FLOCALEID=2052 
LEFT JOIN T_META_OBJECTTYPE_L c ON a.FBILLFORMID=c.FID AND c.FLOCALEID=2052 
WHERE a.FBILLFORMID='BD_Supplier'
-- 编码规则编码模板
SELECT FRULEID AS 编码规则内码,* 
FROM T_BAS_BILLCODERULEENTRY 
WHERE FRULEID='cbda6d5a48e0438ea9729e9279a1e59f'
-- 编码规则流水号
SELECT FRULEID AS 编码规则内码,FNUMMAX AS 当前编码规则最大流水号,* 
FROM T_BAS_BILLCODES 
WHERE FRULEID='cbda6d5a48e0438ea9729e9279a1e59f'
-- 编码规则启用补号 （银行转账单，启用补号功能。）
UPDATE
    T_BAS_BILLCODERULE
SET
    FISREPAIRBILLNO = 1,
    --启用手动补号
    FISAUTOBILLNO = 1 --启用自动补号
WHERE
    FRULEID IN (
        SELECT
            FRULEID
        FROM
            T_BAS_BILLCODERULE_L
        WHERE
            FLOCALEID = 2052
            AND FNAME = N'银行转账单标准编码规则'
    )
-------------------------------------辅助资料-------------------------------------
-- 辅助资料类型
SELECT a.FID AS 辅助资料类型内码,b.FNAME AS 辅助资料类型名称,a.* 
FROM T_BAS_ASSISTANTDATA a 
LEFT JOIN T_BAS_ASSISTANTDATA_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
WHERE b.FNAME='等级'
-- 辅助资料
SELECT a.FID AS 辅助资料类型内码,a.FNUMBER AS 辅助资料编码,b.FDATAVALUE AS 辅助资料名称,a.FENTRYID AS 辅助资料内码,a.* 
FROM T_BAS_ASSISTANTDATAENTRY a 
LEFT JOIN T_BAS_ASSISTANTDATAENTRY_L b ON a.FENTRYID=b.FENTRYID AND b.FLOCALEID=2052 
WHERE a.FID='005056a3024b911211e36619ce35ebd3' 
ORDER BY a.FSEQ
-------------------------------------枚举-------------------------------------
-- 枚举类型
SELECT a.FID AS 枚举类型内码,b.FNAME AS 枚举类型名称,* 
FROM T_META_FORMENUM a 
LEFT JOIN T_META_FORMENUM_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
WHERE b.FNAME=N'发货状态'
-- 枚举项
SELECT a.FID AS 枚举类型内码,a.FENUMID AS 枚举项内码,b.FCAPTION AS 枚举项名称,a.FVALUE AS 枚举项值,* 
FROM T_META_FORMENUMITEM a 
LEFT JOIN T_META_FORMENUMITEM_L b ON a.FENUMID=b.FENUMID AND b.FLOCALEID=2052 
WHERE a.FID='0113bd8e-03b1-4c47-8fb4-9ed5b9074680' 
ORDER BY a.FSEQ
-------------------------------------通过表名查询业务对象-------------------------------------
-- 查询使用了某个表的业务对象
SELECT b.FNAME,a.* FROM T_META_OBJECTTYPE a 
LEFT JOIN T_META_OBJECTTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
WHERE a.FMODELTYPEID IN (100,400) AND CAST(a.FKERNELXML AS NVARCHAR(MAX)) LIKE '%<TableName>T_PUR_POORDER<%'

-- 查询使用了某个表的业务对象，注意：XML查询区分大小写
SELECT b.FNAME,a.* FROM T_META_OBJECTTYPE a 
LEFT JOIN T_META_OBJECTTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
WHERE a.FMODELTYPEID IN (100,400) AND a.FKERNELXML.exist('//TableName[.="t_PUR_POOrder"]')=1

-------------------------------------查询业务对象发布状态-------------------------------------
SELECT a.FID AS 业务对象标识,b.FNAME AS 业务对象名称
,(CASE WHEN d.FDETAILFUNCID IS NULL THEN '否' ELSE '是' END) AS 是否已创建菜单
,e.FNAME AS 菜单名称
,(CASE d.FSTATUS WHEN '1' THEN '已发布' WHEN '0' THEN '未发布' ELSE '' END) AS 菜单发布状态
,g.FNAME AS 子功能名称
,i.FNAME AS 所属子系统
,j.FNAME AS 所属业务领域
FROM T_META_OBJECTTYPE a
JOIN T_META_OBJECTTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052
LEFT JOIN T_META_CONSOLEDETAIL d ON a.FID=d.FOBJECTID
LEFT JOIN T_META_CONSOLEDETAIL_L e ON d.FDETAILFUNCID=e.FDETAILFUNCID AND e.FLOCALEID=2052
LEFT JOIN T_META_CONSOLESUBFUNC f ON d.FSUBFUNCID=f.FSUBFUNCID
LEFT JOIN T_META_CONSOLESUBFUNC_L g ON f.FSUBFUNCID=g.FSUBFUNCID AND g.FLOCALEID=2052
LEFT JOIN T_META_SUBSYSTEM h ON f.FSUBSYSTEMID=h.FID
LEFT JOIN T_META_SUBSYSTEM_L i ON h.FID=i.FID AND i.FLOCALEID=2052
LEFT JOIN T_META_TOPCLASS_L j ON h.FTOPCLASSID=j.FTOPCLASSID AND j.FLOCALEID=2052
WHERE a.FDEVTYPE IN (0,1)
```

## 3.2 修改数据库表字段

```sql
-- 存储过程名称：p_AlterColumn
-- 功能说明：修改列定义，包括：数据类型，精度，是否空，默认值。
-- 参数说明：
@TableName VARCHAR(50), -- 表名，例如：采购订单，T_PUR_POORDER
@ColumnName VARCHAR(50), -- 列名，例如：单据编号，FBILLNO
@ColumnType VARCHAR(50), -- 列类型，，例如：NVARCHAR(100)
@ColumnIsNull VARCHAR(10), -- 列是否为空，例如：NOT NULL
@ModifyFlag VARCHAR(10), -- 修改类型，固定4位长度，分别是数据类型，精度，为空，缺省值（各占1位）。使用示例：改精度和缺省值：'0101'
@DefaultValue VARCHAR(255) -- 缺省值，例如，缺省值空：''''''，字符串A：'''A'''

-- 辅助资料的备注字段默认长度255，将该字段长度调整为4000。 BOS设计器对应字段长度范围清空，移除控制； 修改数据库
EXEC p_AlterColumn 'T_BAS_ASSISTANTDATAENTRY_L',-- 表名
'FDESCRIPTION',-- 列名
'NVARCHAR(4000)',-- 字段类型 ModifyFlag='1000' ColumnType = 'NULL' 表示删除
'NOT NULL',-- 是否为空
'0101',-- 功能控制标志位：共4位，分别是数据类型，精度，为空，缺省值。
''' '''-- 默认值

-- 执行上面的存储过程，实际执行的是下面这3条语句
ALTER TABLE T_BAS_ASSISTANTDATAENTRY_L DROP CONSTRAINT DF__T_BAS_ASS__FDESC__24485945
ALTER TABLE T_BAS_ASSISTANTDATAENTRY_L ALTER COLUMN FDESCRIPTION NVARCHAR(4000) NOT NULL 
ALTER TABLE T_BAS_ASSISTANTDATAENTRY_L ADD  CONSTRAINT DF__T_BAS_ASS__FDESC__24485945 DEFAULT (' ') FOR FDESCRIPTION

-- 查看表结构
sp_help '表名' 


-- BOS设计器能支持的文本字段的最大长度是2000，如果需要更大长度的字文本段，需要提前用SQL预生成列，或者后期用SQL修改数据库中的列定义。
-- 调整文本字段：NVARCHAR(4000) NOT NULL DEFAULT('')
EXEC p_AlterColumn 'T_PUR_POORDER', 'F_JAC_TEXT', 'NVARCHAR(4000)', 'NOT NULL', '1101', '''''';

-- 调整文本字段：NVARCHAR(MAX) NOT NULL DEFAULT('')
EXEC p_AlterColumn 'T_PUR_POORDER', 'F_JAC_TEXT', 'NVARCHAR(MAX)', 'NOT NULL', '1101', '''''';

-- 调整文本字段：NTEXT NOT NULL DEFAULT('')
EXEC p_AlterColumn 'T_PUR_POORDER', 'F_JAC_TEXT', 'NTEXT', 'NOT NULL', '1101', '''''';

-- 调整整数字段类型：bigint NOT NULL DEFAULT(0)
EXEC p_AlterColumn 'T_PUR_POORDER', 'F_JAC_INTEGER', 'BIGINT', 'NOT NULL', '1101', '0';

```

## 3.3 查询领域下元素中属性默认值

```sql
-- 比如得到移动领域下的下拉列表字段中的必录显示空白行属性的默认值
--得到移动领域FID 1400
select * from T_MDL_DOMAINMODETYPE_L where FName like '%移动业务%' and FLocaleId=2052
--得到下拉列表元素FID 9
select * from T_MDL_ELEMENTTYPE_l where Fname like '%下拉列表%' and FLocaleId=2052
--得到领域和元素的关系FEntryId 34121
select * from T_MDL_DOMAINMODELELEMENTMAP where Fid=1400 and FElementTypeId=9 
--得到必录显示空白行属性FID 1946
select * from T_MDL_ELEMENTPROPERTY_L where Fname like '%必录显示空白行%' and FLocaleId=2052
--得到必录显示空白行在移动领域下拉列表中的默认值
select * from T_MDL_DOMAINMODELELPROPMAP where FEntryId=34121 and FElementPropertyId=1946
```

## 3.4 查询业务对象，子系统，子系统分组

```sql
SELECT t0.FID,t1.FNAME,  t3.FNAME FTOPCLASSNAME,t4.FNAME FSUBSYSTEMNAME
FROM
T_META_OBJECTTYPE t0
LEFT JOIN V_META_OBJECTTYPE_L t1 ON t0.FID= t1.FID AND t1.FLOCALEID = 2052 --业务对象
INNER JOIN T_META_SUBSYSTEM t2 ON t0.FSUBSYSID = t2.FID --子系统
LEFT JOIN T_META_SUBSYSTEM_L t3 ON t0.FSUBSYSID = t3.FID AND t3.FLOCALEID =2052
LEFT JOIN T_META_TOPCLASS_L t4 ON T2.FTOPCLASSID=  t4.FTOPCLASSID AND t4.FLOCALEID = 2052 --子系统分组
WHERE  t0.FDEVTYPE<>2
```

## 3.5 查询，删除单据转换规则

### 3.5.1 查询

```sql
-- 查询单据转换规则信息
select tobjSource.FNAME as [源单]
     , tobjTarget.FNAME as [目标单]
     , tul.FNAME        as [规则名称]
     , tu.FID           -- 单据转换规则ID
     , tu.FSOURCEFORMID -- 源单业务对象内码
     , tu.FTARGETFORMID -- 目标单业务对象内码
     , tu.FSTATUS       -- 规则状态|0，停用|1，启用
     , tu.FISDEFAULT    -- 是否默认规则|1 是|0 否
      --, tu.FKERNELXML	-- 通用版本XML
from T_META_CONVERTRULE tu -- 单据转换规则
         left join T_META_CONVERTRULE_l tul -- 单据转换规则 多语言
                   on tul.FID = tu.FID
                       and tul.FLOCALEID = '2052'
         left join T_META_OBJECTTYPE_L tobjSource -- 业务对象多语言
                   on tobjSource.FLOCALEID = '2052'
                       and tobjSource.FID = tu.FSOURCEFORMID
         left join T_META_OBJECTTYPE_L tobjTarget -- 业务对象多语言
                   on tobjTarget.FLOCALEID = '2052'
                       and tobjTarget.FID = tu.FTARGETFORMID
where tobjTarget.FNAME like '应付单%' -- 目标单为
-- OR tobjSource.FNAME like  '发货%' -- 源单为
order by tobjTarget.FNAME

-- 查询采购申请单-采购订单的转换规则
SELECT b.FNAME AS 单据转换规则名称,a.FID AS 单据转换规则内码,a.FKERNELXML AS 转换规则元数据,c.FNAME AS 源单名称,d.FNAME AS 目标单名称,a.* FROM T_META_CONVERTRULE a
LEFT JOIN T_META_CONVERTRULE_L b ON a.FID=b.FID AND b.FLOCALEID=2052
LEFT JOIN T_META_OBJECTTYPE_L c ON a.FSOURCEFORMID=c.FID AND c.FLOCALEID=2052
LEFT JOIN T_META_OBJECTTYPE_L d ON a.FTARGETFORMID=d.FID AND d.FLOCALEID=2052
WHERE a.FSOURCEFORMID='PUR_Requisition' AND a.FTARGETFORMID='PUR_PurchaseOrder'

```

### 3.5.2 删除

```sql
-- 删除转换规则，FID是转换规则的内码
DELETE FROM T_META_CONVERTRULE WHERE FID='005056af-13b3-bd30-11e4-e747e5b3cdfe';
DELETE FROM T_META_CONVERTRULE_L WHERE FID='005056af-13b3-bd30-11e4-e747e5b3cdfe';
```

## 3.6 下推关系查询

```sql
-- 下推关系查询
select T2H.FID      AS [下游单据主表ID]
     , T2E.FENTRYID AS [下游单据明细ID]
     , lk.FENTRYID  as [下游单据关联表FENTRYID对应下游单据明细ID]
     , lk.FSBILLID  AS [下游单据关联表FSBILLID对应上游单据主表ID]
     , T1H.FID      AS [上游单据主表ID]
     , lk.FSID      AS [下游单据关联表FSID对应上游单据明细ID]
     , T1E.FENTRYID AS [上游单据明细ID]
     , lk.FSTABLENAME  [下游单据关联表FSTABLENAME]
from T_SAL_OUTSTOCK T2H -- 销售出库单
         left join T_SAL_OUTSTOCKENTRY T2E -- 销售订单明细记录
                   on T2H.fid = T2E.fid
         left join T_SAL_OUTSTOCKENTRY_LK lk -- 销售订单明细下推关系表
                   on T2E.FENTRYID = lk.FENTRYID
         LEFT JOIN T_SAL_ORDERENTRY T1E -- 销售订单明细表
                   ON LK.FSID = T1E.FENTRYID
         LEFT JOIN T_SAL_ORDER T1H -- 销售订单主表
                   on lk.FSBILLID = T1H.FID
where T2H.FDOCUMENTSTATUS = 'C'
  and T2H.fbillno = 'XSCKD145137'
```

## 3.7 查询元数据注册的插件信息

```sql
-- 从元数据中查找插件注册信息
SELECT FKERNELXML.query('//PlugIn') AS PlugIn, * FROM T_META_OBJECTTYPE
WHERE FID='PUR_PurchaseOrder'

SELECT
    FID,
    FKERNELXML.query('//PlugIn//ClassName')
FROM
    T_META_OBJECTTYPE
WHERE
    FMODELTYPEID = 100
    AND (
        FID = 'PUR_PurchaseOrder'
        OR FINHERITPATH LIKE ',%PUR_PurchaseOrder%,'
    )


-- 查询元数据中的Python插件代码。
SELECT
    TOP 10 FID,
    FKERNELXML.query('//PlugIn/PyScript') AS Python脚本
FROM
    T_META_OBJECTTYPE
WHERE
    FKERNELXML.exist('//PlugIn/PlugInType[.="1"]') = 1;
```

## 3.8 根据枚举类型查询引用该类型的单据

```sql
/*查询引用了枚举类型发货状态的单据，发货状态的内码：0113bd8e-03b1-4c47-8fb4-9ed5b9074680*/
SELECT top 100 b.FNAME,a.*
FROM T_META_OBJECTTYPE a
LEFT JOIN T_META_OBJECTTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052
WHERE CAST(a.FKERNELXML AS NVARCHAR(MAX)) LIKE '%<EnumType>0113bd8e-03b1-4c47-8fb4-9ed5b9074680</EnumType>%'

/*查询枚举类型*/
SELECT b.FNAME,a.*
FROM T_META_FORMENUM a
LEFT JOIN T_META_FORMENUM_L b ON a.FID=b.FID AND b.FLOCALEID=2052
WHERE b.FNAME='发货状态'
```

## 3.9 查询所有下拉列表枚举值

```sql
-- 查询查询所有枚举值信息
SELECT b.FNAME "枚举类型名称",c.FSEQ "序号",c.FVALUE "枚举值",d.FCAPTION "枚举项名称",c.FINVALID "是否禁用"  FROM T_META_FORMENUM a
INNER JOIN T_META_FORMENUM_L b ON a.FID =b.FID AND b.FLOCALEID = 2052
INNER JOIN T_META_FORMENUMITEM c ON c.fid = a.FID
INNER JOIN T_META_FORMENUMITEM_L d ON d.FENUMID = c.FENUMID;
```

## 

# 4. 业务对象类

## 4.1 参数

```sql
-------------------------------------参数-------------------------------------
-- 查询系统参数
SELECT * FROM T_BAS_SYSPARAMETER;
-- 查询单据参数
SELECT * FROM T_BAS_BILLGLOBALPARAMETER;
SELECT * FROM T_BAS_BILLGLOBALPARAMFIELD;
SELECT * FROM T_BAS_BILLGLOBALPARAMSEARCH;
-- 查询单据类型参数
SELECT a.FBILLTYPEID AS 单据类型内码, a.FNUMBER AS 单据类型编码, b.FNAME AS 单据类型名称, FPARASETTINGXML AS 单据类型参数, a.* 
FROM T_BAS_BILLTYPE a 
LEFT JOIN T_BAS_BILLTYPE_L b ON a.FBILLTYPEID = b.FBILLTYPEID and b.FLOCALEID = 2052;
-- 查询用户参数
SELECT * FROM T_BAS_USERPARAMETER;
```

## 4.2 用户

### 4.2.1 登录用户

```sql
-------------------------------------登录用户-------------------------------------
SELECT FUSERID AS 用户内码,FNAME AS 用户名称,FUSERACCOUNT AS 用户账号,* 
FROM T_SEC_USER;
```

### 4.2.2 用户和员工的关系

```sql
-- 用户：SEC_User\T_SEC_USER
-- 员工：BD_Empinfo\T_HR_EMPINFO
-- 员工->岗位信息单据体：T_BD_STAFFTEMP
-- 人员详细信息(公共)：DB_Person\T_BD_PERSON （V_BD_CONTACTOBJECT）
-- 员工任岗信息：BD_NEWSTAFF\T_BD_STAFF
-- 员工任岗信息->任职信息： T_BD_STAFFPOSTINFO
-- 员工汇报体系（查看用）
-- 部门：BD_Department\T_BD_DEPARTMENT
-- 岗位信息：HR_ORG_HRPOST\T_ORG_POST
-- 岗位信息->上级汇报关系： T_ORG_POSTREPORTLINE   
-- 岗位信息->岗位角色：T_ORG_POSTROLE
-- 岗位信息->职位子单据头：T_ORG_HRPOST
-- 组织机构：ORG_Organizations\T_ORG_Organizations

-- 通过用户Id找到人员详细信息(公共)  和岗位相关信息
select t0.FUSERID,t0.FNAME,t1.FPERSONID,t2.FNAME,t3.FDEPTID '部门Id',t3.FPOSTID '岗位Id',
ISNULL(t4.FISFIRSTPOST,'0') AS '是否为主任岗位' 
from T_SEC_USER t0
join T_BD_PERSON t1 on t1.FPERSONID = t0.FLINKOBJECT
join T_BD_PERSON_L t2 on t2.FPERSONID = t1.FPERSONID and t2.FLOCALEID=2052
left join T_BD_STAFF t3 on t3.FPERSONID = t1.FPERSONID
left join T_BD_STAFFPOSTINFO t4 on t4.FSTAFFID = t3.FSTAFFID
where t0.FTYPE=1 and t0.FUSERID =624152

-- 用户-人员-员工任岗关系
SELECT
    a.FUSERID AS 用户内码,
    a.FNAME 用户名,
    b.FPERSONID AS 人员内码,
    bl.FNAME AS 人员姓名,
    dl.FNAME AS 任岗部门名称,
    el.FNAME AS 任岗岗位名称
FROM
    T_SEC_USER a -- 用户表
    INNER JOIN T_BD_PERSON b ON a.FLINKOBJECT = b.FPERSONID -- 人员表
    LEFT JOIN T_BD_PERSON_L bl ON b.FPERSONID = bl.FPERSONID
    AND bl.FLOCALEID = 2052
    INNER JOIN T_BD_STAFF c ON b.FPERSONID = c.FPERSONID -- 员工任岗表
    LEFT JOIN T_BD_DEPARTMENT d ON c.FDEPTID = d.FDEPTID -- 部门表
    LEFT JOIN T_BD_DEPARTMENT_L dl ON d.FDEPTID = dl.FDEPTID
    AND dl.FLOCALEID = 2052
    LEFT JOIN T_ORG_POST e ON c.FPOSTID = e.FPOSTID -- 岗位表
    LEFT JOIN T_ORG_POST_L el ON e.FPOSTID = el.FPOSTID
    AND el.FLOCALEID = 2052
WHERE
    a.FNAME = 'demo'
```

### 4.2.3 查询用户组织和角色的关系

```sql
select t0.FUSERID,t0.FNAME userName,t1.FORGID,t12.FNAME orgName,t2.FROLEID,t22.FNAME roleName from T_SEC_USER t0
inner join T_SEC_USERORG  t1 on t1.FUSERID = t0.FUSERID
left join T_ORG_ORGANIZATIONS_L t12 on t12.FORGID = t1.FORGID and t12.FLOCALEID=2052
inner join T_SEC_USERROLEMAP t2 on t2.FENTITYID = t1.FENTITYID
left join T_SEC_ROLE_L t22 on t22.FROLEID = t2.FROLEID and t22.FLOCALEID=2052
where t0.FNAME='Demo'
```

### 4.2.4 查询用户业务对象字段权限

```sql
SELECT t0.FUSERID,t0.FName,t1.FORGID,t2.FNAME orgName ,t3.FROLEID,t4.FNAME roleName,t5.FOBJECTTYPEID,t7.Fname objectName,
case WHEN t6.FVIEW is null then '有权'
when t6.FVIEW =1 then '无权' else '禁用'  end 'ViewPermission',
case WHEN t6.FEDIT is null then '有权'
when t6.FEDIT =1 then '无权' else '禁用'  end 'EditPermission'
,t6.FFIELDKEY
FROM  T_SEC_USER t0
INNER JOIN T_SEC_USERORG t1 ON  t1.FUSERID = t0.FUSERID
left join T_ORG_ORGANIZATIONS_L t2 on t2.FORGID = t1.FORGID and t2.FLOCALEID=2052
INNER JOIN T_SEC_USERROLEMAP t3 ON t3.FENTITYID = t1.FENTITYID
left join T_SEC_ROLE_L t4 on t4.FROLEID = t3.FROLEID and t4.FLOCALEID=2052
LEFT JOIN T_SEC_FIELDPERMISSION t5 ON t5.FROLEID = t3.FROLEID
LEFT JOIN T_SEC_FIELDPERMISSIONENTRY t6 ON t6.FITEMID = t5.FITEMID
left join V_META_OBJECTTYPE_L t7 on t7.FID = t5.FOBJECTTYPEID and t7.FLOCALEID=2052
WHERE t0.FName = 'demo' AND t5.FOBJECTTYPEID ='BOS_BillEE'
```

### 4.2.5 查询指定用户最近执行的SQL

```sql
-- 查询连接到当前数据库的进程信息
SELECT * FROM sys.sysprocesses WHERE dbid=db_id()

-- 查询sa用户最后执行的sql语句
SELECT * FROM sys.sysprocesses a
CROSS APPLY sys.dm_exec_sql_text (a.sql_handle) b
WHERE a.dbid=db_id() and a.loginame='sa'

--查询数据库最近执行的前100条sql语句
SELECT TOP 100 QS.creation_time AS '执行时间',ST.text AS '执行的SQL语句',*
FROM sys.dm_exec_query_stats QS
CROSS APPLY sys.dm_exec_sql_text (QS.sql_handle) ST
ORDER BY QS.creation_time DESC
```

## 4.3 网控

```sql
-------------------------------------网控-------------------------------------
-- 网控类型
SELECT a.FID AS 网控类型内码, b.FNAME AS 网控类型名称, a.* 
FROM T_BAS_NETWORKCTRLTYPE a 
LEFT JOIN T_BAS_NETWORKCTRLTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052
-- 网控对象
SELECT a.FID AS 网控对象内码,b.FNAME AS 网控对象名称,t.FNAME AS 网控类型名称, c.FNAME AS 业务对象名称, e.FNAME AS 操作名称, a.* 
FROM T_BAS_NETWORKCTRLOBJECT a 
LEFT JOIN T_BAS_NETWORKCTRLOBJECT_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
LEFT JOIN T_BAS_NETWORKCTRLTYPE_L t ON a.FTYPE=t.FID AND t.FLOCALEID=2052 
LEFT JOIN T_META_OBJECTTYPE_L c ON a.FMETAOBJECTID=c.FID AND c.FLOCALEID=2052 
LEFT JOIN T_MDL_FORMOPERATIONTYPE d ON a.FNUMBER=d.FOPERATION 
LEFT JOIN T_MDL_FORMOPERATIONTYPE_L e ON d.FID=e.FID AND e.FLOCALEID=2052 
--WHERE a.FID='903eafb2784042b5948fe9a245b5dddf'
-- 网控互斥对象
SELECT b.FNAME AS 网控对象,c.FNAME AS 互斥对象, a.* 
FROM T_BAS_NETWORKCTRLMUTEX a 
LEFT JOIN T_BAS_NETWORKCTRLOBJECT_L b ON a.FNETCTRLID=b.FID AND b.FLOCALEID=2052 
LEFT JOIN T_BAS_NETWORKCTRLOBJECT_L c ON a.FMUTEXOBJID=c.FID AND c.FLOCALEID=2052 
--WHERE a.FNETCTRLID='903eafb2784042b5948fe9a245b5dddf'
-- 网控数据
SELECT * FROM T_BAS_NETWORKCTRLRECORDS;
```

## 4.4 供应链

### 4.4.1 基础资料

```sql
-- 物料
SELECT *
FROM (SELECT 
  a.FMATERIALID '物料ID',
  a.FNUMBER '物料编码',
  b.FNAME '物料名称',
  b.FSPECIFICATION '规格型号',
  c.FNUMBER '分组编码',
  d.FNAME '分组名称'
FROM  dbo.T_BD_MATERIAL a
INNER JOIN dbo.T_BD_MATERIAL_L b ON b.FMATERIALID = a.FMATERIALID
INNER JOIN T_BD_MATERIALGROUP c ON a.FMATERIALGROUP = c.FID
INNER JOIN dbo.T_BD_MATERIALGROUP_L d ON c.FID = d.FID) MAT;

-- 客户
SELECT 
a.FCUSTID '客户ID',
a.FNUMBER '客户编码',
b.FNAME '客户名称'
FROM dbo.T_BD_CUSTOMER a
INNER JOIN dbo.T_BD_CUSTOMER_L b ON b.FCUSTID = a.FCUSTID;

-- 客户分组级次
SELECT A.FCUSTID '客户ID',
           D.FLevel '级次',
        GROUP1.分组名称 '一级分组名称',
        GROUP2.分组名称 '二级分组名称',
        GROUP3.分组名称 '三级分组名称',
        GROUP4.分组名称 '四级分组名称',
        GROUP5.分组名称 '五级分组名称',
        GROUP6.分组名称 '六级分组名称',
        GROUP7.分组名称 '七级分组名称'
    FROM dbo.T_BD_CUSTOMER A
        INNER JOIN T_BD_CUSTOMERGROUP B ON A.FPRIMARYGROUP = B.FID
        INNER JOIN dbo.T_BD_CUSTOMERGROUP_L C ON B.FID = C.FID
        INNER JOIN T_BD_CUSTOMERGROUPLEVEL D ON B.FID = D.FID
  
     LEFT JOIN
     (SELECT A.FID '分组ID',B.FNAME '分组名称' FROM T_BD_CUSTOMERGROUP A
        INNER JOIN dbo.T_BD_CUSTOMERGROUP_L B
            ON A.FID = B.FID) GROUP1 ON D.FGROUP1 = GROUP1.分组ID

     LEFT JOIN
     (SELECT A.FID '分组ID',B.FNAME '分组名称' FROM T_BD_CUSTOMERGROUP A
        INNER JOIN dbo.T_BD_CUSTOMERGROUP_L B
            ON A.FID = B.FID) GROUP2 ON D.FGROUP2 = GROUP2.分组ID

     LEFT JOIN
     (SELECT A.FID '分组ID',B.FNAME '分组名称' FROM T_BD_CUSTOMERGROUP A
        INNER JOIN dbo.T_BD_CUSTOMERGROUP_L B
            ON A.FID = B.FID) GROUP3 ON D.FGROUP3 = GROUP3.分组ID

     LEFT JOIN
     (SELECT A.FID '分组ID',B.FNAME '分组名称' FROM T_BD_CUSTOMERGROUP A
        INNER JOIN dbo.T_BD_CUSTOMERGROUP_L B
            ON A.FID = B.FID) GROUP4 ON D.FGROUP4 = GROUP4.分组ID

     LEFT JOIN
     (SELECT A.FID '分组ID',B.FNAME '分组名称' FROM T_BD_CUSTOMERGROUP A
        INNER JOIN dbo.T_BD_CUSTOMERGROUP_L B
            ON A.FID = B.FID) GROUP5 ON D.FGROUP5 = GROUP5.分组ID

     LEFT JOIN
     (SELECT A.FID '分组ID',B.FNAME '分组名称' FROM T_BD_CUSTOMERGROUP A
        INNER JOIN dbo.T_BD_CUSTOMERGROUP_L B
            ON A.FID = B.FID) GROUP6 ON D.FGROUP6 = GROUP6.分组ID

     LEFT JOIN
     (SELECT A.FID '分组ID',B.FNAME '分组名称' FROM T_BD_CUSTOMERGROUP A
        INNER JOIN dbo.T_BD_CUSTOMERGROUP_L B
            ON A.FID = B.FID) GROUP7 ON D.FGROUP7 = GROUP7.分组ID
    WHERE A.FNUMBER = '01.01.01.12.01';

-- 供应商
SELECT  
  a.FSUPPLIERID '供应商ID',
  a.FNUMBER '供应商编码',
  b.FNAME '供应商名称'
FROM dbo.T_BD_SUPPLIER a
INNER JOIN dbo.T_BD_SUPPLIER_L b ON b.FSUPPLIERID = a.FSUPPLIERID;

-- 组织
SELECT  
a.FORGID '组织ID',
a.FNUMBER '组织编码',
b.FNAME '组织名称'
FROM dbo.T_ORG_ORGANIZATIONS a
INNER JOIN T_ORG_ORGANIZATIONS_L b ON b.FORGID = a.FORGID;

-- 部门
SELECT  
a.FDEPTID '部门ID',
a.FNUMBER '部门编码',
b.FNAME '部门名称'
FROM dbo.T_BD_DEPARTMENT a
INNER JOIN dbo.T_BD_DEPARTMENT_L b ON a.FDEPTID = b.FDEPTID;

-- 仓库
SELECT 
a.FSTOCKID '仓库ID',
a.FNUMBER '仓库编码',
b.FNAME '仓库名称'
FROM dbo.T_BD_STOCK a INNER JOIN dbo.T_BD_STOCK_L b ON b.FSTOCKID = a.FSTOCKID;
```

### 4.4.2 业务单据

```sql
-- 销售订单
SELECT
    *
FROM
    (
        SELECT
            *
        FROM
            dbo.T_SAL_ORDER A
            INNER JOIN dbo.T_SAL_ORDERENTRY B ON B.FID = A.FID
            INNER JOIN dbo.T_SAL_ORDERENTRY_F C ON C.FENTRYID = B.FENTRYID
            INNER JOIN dbo.T_SAL_ORDERENTRY_R D ON D.FENTRYID = B.FENTRYID
        WHERE
            A.FDOCUMENTSTATUS = 'C'
    ) SAL;

-- 销售出库单
SELECT
    *
FROM
    dbo.T_SAL_OUTSTOCK A
    INNER JOIN dbo.T_SAL_OUTSTOCKENTRY B ON B.FID = A.FID
```

### 4.4.3 仓库，仓位

```sql
-- 仓库和仓位取值关系
SELECT
B.FSTOCKLOCID,
B.*
FROM dbo.T_SAL_OUTSTOCK A INNER JOIN dbo.T_SAL_OUTSTOCKENTRY B ON B.FID = A.FID
WHERE A.FBILLNO = 'XSCKD000054'; --102345

--根据仓库编码查询仓库信息
SELECT 
FSTOCKID
FROM dbo.T_BD_STOCK WHERE FNUMBER = 'CK01LQ'; 

-- 仓库和仓位值集表，通过仓库ID查询出仓位值的FENTRYID
SELECT 
FENTRYID
FROM dbo.T_BD_STOCKFLEXITEM
WHERE FSTOCKID = '128344'; 

--根据上表查询出来的FENTRYID查询出仓位FFLEXENTRYID这个值
SELECT *
FROM T_BD_STOCKFLEXDETAIL
WHERE FENTRYID = '100013';

--根据上表的FFLEXENTRYID值查询出仓位的ID，单据上的仓位字段存的就是这个值
SELECT
loc.FID AS '仓位ID'
FROM  T_BAS_FLEXVALUESDETAIL loc
LEFT JOIN T_BAS_FLEXVALUESENTRY locCol1 ON loc.FF100015 = locCol1.FENTRYID
LEFT JOIN T_BAS_FLEXVALUESENTRY_L locCol1L ON locCol1.FENTRYID = locCol1L.FENTRYID
WHERE locCol1.FENTRYID = '100146';

-- 【总结】通过仓位ID查询仓库SQL
SELECT  
a.仓库ID,
a.仓库名称,
d.仓位ID,
d.仓位名称,
c.无值内码,
c.内码
FROM (   
  SELECT  
  a.FSTOCKID '仓库ID',
  b.FNAME '仓库名称'
  FROM      dbo.T_BD_STOCK a
  INNER JOIN dbo.T_BD_STOCK_L b ON a.FSTOCKID = b.FSTOCKID
) a
LEFT JOIN (   
  SELECT 
  FENTRYID '仓位值表FENTRYID',
  FSTOCKID '仓库ID'
  FROM dbo.T_BD_STOCKFLEXITEM
) b ON b.仓库ID = a.仓库ID
LEFT JOIN (   
  SELECT 
  FDETAILID '内码',
  FENTRYID '仓位值表FENTRYID',
  FFLEXENTRYID '仓位表FENTRYID',
  F_T_IIDE '无值内码'
  FROM T_BD_STOCKFLEXDETAIL
) c ON c.仓位值表FENTRYID = b.仓位值表FENTRYID
LEFT JOIN (   
  SELECT  
  loc.FID AS '仓位ID',
  locCol1L.FNAME '仓位名称',
  locCol1.FENTRYID '仓位表FENTRYID'
  FROM T_BAS_FLEXVALUESDETAIL loc
  LEFT JOIN T_BAS_FLEXVALUESENTRY locCol1 ON loc.FF100015 = locCol1.FENTRYID
  LEFT JOIN T_BAS_FLEXVALUESENTRY_L locCol1L ON locCol1.FENTRYID = locCol1L.FENTRYID
) d ON d.仓位表FENTRYID  = c.仓位表FENTRYID
WHERE d.仓位ID = '102345';

```

### 4.4.4 即时库存

#### 4.4.4.1 查询即时库存

```sql
SELECT
    TI.FSTOCKORGID,
    OL.FNAME fstockorgname,
    M.FNUMBER '物料编码', --物料编码   
    ML.FNAME '物料名称', --物料名称   
    TSL.FNAME '仓库', --仓库   
    TI.FSTOCKLOCID '仓位ID', --仓位ID   
    TUL0.FNAME '库存主单位', --库存主单位   
    TI.FQTY '主单位库存量',
    TUL1.FNAME '基本单位', --基本单位  
    TI.FBASEQTY '基本单位库存量',
    TUL2.FNAME '库存辅单位', --库存辅单位   
    TI.FSECQTY '库存辅单位库存量',
    TL.FNUMBER flotnumber,
    TI.FAUXPROPID,
    CASE
        WHEN TMS.FISEXPPARTOFLOT = '1' THEN CONVERT (
            CHAR (10),
            TL.FPRODUCEDATE,
            20
        )
        ELSE CONVERT (
            CHAR (10),
            TI.FPRODUCEDATE,
            20
        )
    END fproducedate,
    CASE
        WHEN TMS.FISEXPPARTOFLOT = '1' THEN CONVERT (CHAR(10), TL.FEXPIRYDATE, 20)
        ELSE CONVERT (CHAR(10), TI.FEXPIRYDATE, 20)
    END fexpirydate,
    TB.FNUMBER fbomnumber, --BOM编号  
    TSSL.FNAME fstockstatus,
    TI.FOWNERTYPEID,
    VO_L.FNAME fownername,
    TI.FKEEPERTYPEID, --保管者类型  
    VK_L.FNAME fkeepername,
    TMS.FSTOREURNUM,
    TMS.FSTOREURNOM,
    TMS.FISSNMANAGE,
    TSK.FALLOWMINUSQTY,
    TUS.FPRECISION fstkprecision,
    TUS.FROUNDTYPE froundtype,
    TUE.FPRECISION fsecprecision,
    TI.FMTONO,
    TI.FPROJECTNO,
    TSUB.FBASELOCKQTY fbaselockqty,
    TSUB.FSECLOCKQTY fseclockqty,
    '' fstocklocname
FROM
    dbo.T_STK_INVENTORY TI
    LEFT OUTER JOIN T_ORG_ORGANIZATIONS_L OL ON (
        TI.FSTOCKORGID = OL.FORGID
        AND OL.FLOCALEID = 2052
    ) --机构名称 
    INNER JOIN T_BD_MATERIAL M ON (
        M.FMASTERID = TI.FMATERIALID
        AND (
            M.FUSEORGID = TI.FSTOCKORGID
            OR EXISTS (
                SELECT
                    1
                FROM
                    T_META_BASEDATATYPE BT
                WHERE
                    (
                        BT.FBASEDATATYPEID = 'BD_MATERIAL'
                        AND BT.FSTRATEGYTYPE = 1
                    )
            )
        )
    ) --物料 
    LEFT OUTER JOIN T_BD_MATERIAL_L ML ON (
        M.FMATERIALID = ML.FMATERIALID
        AND ML.FLOCALEID = 2052
    ) --物料名称、规格型号 
    INNER JOIN T_BD_MATERIALSTOCK TMS ON M.FMATERIALID = TMS.FMATERIALID
    LEFT OUTER JOIN T_ENG_BOM TB ON (
        TB.FMASTERID = TI.FBOMID
        AND (
            TB.FUSEORGID = TI.FSTOCKORGID
            OR EXISTS (
                SELECT
                    1
                FROM
                    T_META_BASEDATATYPE BT
                WHERE
                    (
                        BT.FBASEDATATYPEID = 'ENG_BOM'
                        AND BT.FSTRATEGYTYPE = 1
                    )
            )
        )
    )
    LEFT OUTER JOIN T_BD_STOCKSTATUS_L TSSL ON (
        TI.FSTOCKSTATUSID = TSSL.FSTOCKSTATUSID
        AND TSSL.FLOCALEID = 2052
    )
    LEFT OUTER JOIN T_BD_LOTMASTER TL ON (
        (
            TI.FLOT = TL.FMASTERID
            AND TI.FSTOCKORGID = TL.FUSEORGID
        )
        AND TL.FBIZTYPE = '1'
    )
    LEFT OUTER JOIN T_BD_UNIT TUS ON TMS.FSTOREUNITID = TUS.FUNITID
    LEFT OUTER JOIN T_BD_UNIT_L TUL0 ON (
        TMS.FSTOREUNITID = TUL0.FUNITID
        AND TUL0.FLOCALEID = 2052
    )
    INNER JOIN T_BD_STOCK TSK ON (
        TSK.FMASTERID = TI.FSTOCKID
        AND (
            TSK.FUSEORGID = TI.FSTOCKORGID
            OR EXISTS (
                SELECT
                    1
                FROM
                    T_META_BASEDATATYPE BT
                WHERE
                    (
                        BT.FBASEDATATYPEID = 'BD_STOCK'
                        AND BT.FSTRATEGYTYPE = 1
                    )
            )
        )
    )
    LEFT OUTER JOIN T_BD_STOCK_L TSL ON (
        TSL.FSTOCKID = TI.FSTOCKID
        AND TSL.FLOCALEID = 2052
    )
    LEFT OUTER JOIN T_BD_UNIT_L TUL1 ON (
        TI.FBASEUNITID = TUL1.FUNITID
        AND TUL1.FLOCALEID = 2052
    )
    LEFT OUTER JOIN T_BD_UNIT TUE ON TMS.FAUXUNITID = TUE.FUNITID
    LEFT OUTER JOIN T_BD_UNIT_L TUL2 ON (
        TI.FSECUNITID = TUL2.FUNITID
        AND TUL2.FLOCALEID = 2052
    )
    LEFT OUTER JOIN V_ITEMCLASS_OWNER VO ON (
        (
            VO.FMASTERID = TI.FOWNERID
            AND VO.FFORMID = TI.FOWNERTYPEID
        )
        AND (
            VO.FUSEORGID = TI.FSTOCKORGID
            OR VO.FUSEORGID = 0
            OR EXISTS (
                SELECT
                    1
                FROM
                    T_META_BASEDATATYPE BT
                WHERE
                    (
                        BT.FBASEDATATYPEID = VO.FFORMID
                        AND BT.FSTRATEGYTYPE = 1
                    )
            )
        )
    )
    LEFT OUTER JOIN V_ITEMCLASS_OWNER_L VO_L ON (
        VO.FITEMID = VO_L.FITEMID
        AND VO_L.FLOCALEID = 2052
    )
    LEFT OUTER JOIN V_ITEMCLASS_KEEPER VK ON (
        (
            VK.FMASTERID = TI.FKEEPERID
            AND VK.FFORMID = TI.FKEEPERTYPEID
        )
        AND (
            VK.FUSEORGID = TI.FSTOCKORGID
            OR VK.FUSEORGID = 0
            OR EXISTS (
                SELECT
                    1
                FROM
                    T_META_BASEDATATYPE BT
                WHERE
                    (
                        BT.FBASEDATATYPEID = VK.FFORMID
                        AND BT.FSTRATEGYTYPE = 1
                    )
            )
        )
    )
    LEFT OUTER JOIN V_ITEMCLASS_KEEPER_L VK_L ON (
        VK.FITEMID = VK_L.FITEMID
        AND VK_L.FLOCALEID = 2052
    ) --以下仓位 
    LEFT OUTER JOIN T_BAS_FLEXVALUESDETAIL FVD ON TI.FSTOCKLOCID = FVD.FID
    LEFT OUTER JOIN T_BAS_FLEXVALUESENTRY_L VFF100001 ON (
        FVD.FOPCODE = VFF100001.FENTRYID
        AND VFF100001.FLOCALEID = 2052
    )
    LEFT OUTER JOIN T_BAS_FLEXVALUESENTRY_L VFF100002 ON (
        FVD.FOPCODE = VFF100002.FENTRYID
        AND VFF100002.FLOCALEID = 2052
    )
    LEFT OUTER JOIN T_BAS_FLEXVALUESENTRY_L VFF100003 ON (
        FVD.FOPCODE = VFF100003.FENTRYID
        AND VFF100003.FLOCALEID = 2052
    ) --以下锁库和预留 
    LEFT OUTER JOIN (
        SELECT
            TLKE.FSUPPLYINTERID finventryid,
            SUM (TLKE.FBASEQTY) fbaselockqty,
            SUM (TLKE.FSECQTY) fseclockqty
        FROM
            T_PLN_RESERVELINKENTRY TLKE
            INNER JOIN T_PLN_RESERVELINK TLKH ON TLKE.FID = TLKH.FID
        WHERE
            (
                TLKE.FSUPPLYFORMID = 'STK_Inventory'
                AND TLKE.FLINKTYPE = '4'
            )
        GROUP BY
            TLKE.FSUPPLYINTERID
    ) tsub ON TI.FID = TSUB.FINVENTRYID
WHERE
    TI.FBASEQTY > 0
```

#### 4.4.4.2 查询即时库存可用量

```sql
SELECT
    TI.FSTOCKORGID,
    TI.FKEEPERTYPEID,
    TI.FKEEPERID,
    TI.FOWNERTYPEID,
    TI.FOWNERID,
    TI.FMATERIALID,
    TM.FMATERIALID AS MaterialId,
    TM.FNUMBER AS MaterialNum,
    TI.FSTOCKID,
    TI.FSTOCKLOCID,
    TI.FAUXPROPID,
    TI.FSTOCKSTATUSID,
    TI.FLOT,
    TI.FBOMID,
    TI.FMTONO,
    TI.FPROJECTNO,
    (
        CASE
            WHEN TMS.FISEXPPARTOFLOT = '1' THEN TL.FPRODUCEDATE
            ELSE TI.FPRODUCEDATE
        END
    ) AS FPRODUCEDATE,
    (
        CASE
            WHEN TMS.FISEXPPARTOFLOT = '1' THEN TL.FEXPIRYDATE
            ELSE TI.FEXPIRYDATE
        END
    ) AS FEXPIRYDATE,
    TI.FQTY,
    --基本单位库存量、辅助单位库存量
    TI.FBASEUNITID,
    TI.FBASEQTY,
    TI.FSECUNITID,
    TI.FSECQTY --基本单位锁库量、辅助单位锁库量
,
    TSUB.FBASELOCKQTY AS FBASELOCKQTY,
    (
        CASE
            WHEN TMS.FSTOREURNUM = 0 THEN (TI.FBASEQTY - ISNULL(TSUB.FBASELOCKQTY, 0))
            ELSE (TI.FBASEQTY - ISNULL(TSUB.FBASELOCKQTY, 0)) * TMS.FSTOREURNOM / TMS.FSTOREURNUM
        END
    ) FAVBQTY,
    (TI.FBASEQTY - ISNULL(TSUB.FBASELOCKQTY, 0)) FBASEAVBQTY,
    TSUB.FSECLOCKQTY AS FSECLOCKQTY,
    --库存单位数量
    TMS.FSTOREURNUM fstoreurnum,
    TMS.FSTOREURNOM fstoreurnom,
    TUS.FPRECISION as fprecision,
    Convert(decimal(23, 10), TI.FBASEQTY / TMS.FSTOREURNUM) fstockunitqty
FROM
    T_STK_INVENTORY TI
    INNER JOIN T_BD_MATERIAL TM ON TI.FMATERIALID = TM.FMASTERID
    AND TI.FSTOCKORGID = TM.FUSEORGID
    INNER JOIN T_BD_MATERIALSTOCK TMS ON TM.FMATERIALID = TMS.FMATERIALID
    LEFT JOIN T_BD_LOTMASTER TL ON TI.FLOT = TL.FMASTERID
    AND TI.FSTOCKORGID = TL.FUSEORGID
    AND TL.FBIZTYPE = '1'
    LEFT OUTER JOIN T_BD_UNIT TUS ON TMS.FSTOREUNITID = TUS.FUNITID
    LEFT OUTER JOIN T_BD_UNIT_L TUL0 ON (
        TMS.FSTOREUNITID = TUL0.FUNITID
        AND TUL0.FLOCALEID = 2052
    ) --以下锁库和预留
    LEFT JOIN (
        SELECT
            TLKE.FSUPPLYINTERID AS FINVENTRYID,
            SUM(TLKE.FBASEQTY) AS FBASELOCKQTY,
            SUM(TLKE.FSECQTY) AS FSECLOCKQTY
        FROM
            T_PLN_RESERVELINKENTRY TLKE
            INNER JOIN T_PLN_RESERVELINK TLKH ON TLKE.FID = TLKH.FID
        WHERE
            TLKE.FSUPPLYFORMID = 'STK_Inventory'
            AND TLKE.FLINKTYPE = '4'
        GROUP BY
            TLKE.FSUPPLYINTERID
    ) TSUB ON TI.FID = TSUB.FINVENTRYID
WHERE
    TI.FISEFFECTIVED = '1' --AND TM.FMATERIALID='328233'  ---测试查询的物料内码
```

### 4.4.5 全流程跟踪

```sql
-- 全流程跟踪 - 下查
-- 【按分录查询】
-- FTId：分录数据行ID
-- FTTableName ：分录表名
-- FTTableId：分录表ID
-- 【按整单查询】
-- FTId：单据ID
-- FTTableName ：单据头表名
-- FTTableId：单据头表ID
-- 其中，FTTableName 和FTTableId通过以下SQL获取：SELECT * FROM t_BF_TableDefine WHERE FFORMID='STK_INSTOCK'

-- 全流程跟踪查询路线
-- 1先查流程表，找到FMasterId
-- 2根据FMasterId，再去流程表t_BF_Instance查FInstanceId
-- 3查询业务流程实例
SELECT *
FROM t_BF_Instance TInst 
INNER JOIN t_BF_InstanceEntry TEntry ON (TInst.FInstanceId = TEntry.FInstanceId)  AND TEntry.FTId=100494
WHERE TEntry.FTTableName = 'T_STK_INSTOCKENTRY' 

-- 4查询业务流程历史表实例
-- FTTableId来自t_BF_TableDefine.FSeq：SELECT * FROM t_BF_TableDefine WHERE FFORMID='STK_INSTOCK'
SELECT *
FROM t_BF_InstanceHis TInst 
INNER JOIN t_BF_InstanceEntryHis TEntry ON (TInst.FInstanceId = TEntry.FInstanceId) AND TEntry.FTId=100494
WHERE TEntry.FTTableId = 5100 

-- 5查询业务流程备份表实例
SELECT *
FROM t_BF_InstBackUp TInst 
INNER JOIN t_BF_InstEntryBackUp TEntry ON (TInst.FInstanceId = TEntry.FInstanceId) AND TEntry.FTId=100494
WHERE TEntry.FTTableId = 5100 

-- 6根据FMasterId，再去流程表t_BF_Instance查FInstanceId
SELECT Distinct t0.FInstanceId, t0.FFlowId, t0.FSourceId, t0.FMasterId, t0.FStatus, t0.FStartTime 
FROM  t_BF_Instance t0 where t0.FMasterId ='5a1cbdcce64e72'
```

## 4.5 附件及文件服务

```sql
-- 附件信息表
SELECT * FROM T_BAS_ATTACHMENT 
WHERE FBILLTYPE='PUR_PurchaseOrder'
-- 枚举类型：文件存储位置
SELECT b.FNAME AS 枚举类型名称,* 
FROM T_META_FORMENUM a 
JOIN T_META_FORMENUM_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
WHERE b.FNAME=N'文件存储位置'
-- 枚举项：文件存储位置
SELECT b.FCAPTION AS 枚举项名称,a.FVALUE AS  枚举项值,* 
FROM T_META_FORMENUMITEM a 
JOIN T_META_FORMENUMITEM_L b ON a.FENUMID=b.FENUMID AND b.FLOCALEID=2052 
WHERE a.FID='5f910924-50c4-4c1b-8fd6-d006ef9d2aac' 
ORDER BY a.FVALUE
-- 文件信息表，查询文件信息
SELECT * FROM T_BAS_FILESERVERFILEINFO ORDER BY FCREATEDATE DESC
-- 查询文件存储目录信息（对应存储类型：文件服务器）
SELECT * FROM T_BAS_FILESERVERDIRINFO
-- 文件服务器信息表（对应存储类型：文件服务器）*/ 
SELECT * FROM T_BAS_FILESERVERINFO
```

## 4.6 辅助资料

### 4.6.1 辅助资料查询

```sql
-- 查询辅助资料数据
SELECT a.FID AS 操作内码,b.FNAME AS 操作名称,a.FOPERATION AS 操作编码
,a.FSERVICECLASS AS 运行时类全名,a.FDESIGNERCLASS AS 设计时类全名
,a.* FROM T_MDL_FORMOPERATIONTYPE a 
JOIN T_MDL_FORMOPERATIONTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052
-- 查询辅助资料类型
SELECT  a.FID AS 辅助资料类型内码 ,
        a.FNUMBER AS 辅助资料类型编码 ,
        b.FNAME AS 辅助资料类型名称 ,
        c.FNAME AS 所属业务领域
FROM    T_BAS_ASSISTANTDATA a
        LEFT JOIN T_BAS_ASSISTANTDATA_L b ON a.FID = b.FID
                                             AND b.FLOCALEID = 2052
        LEFT JOIN T_META_TOPCLASS_L c ON a.FTOPCLASSID = c.FTOPCLASSID
                                             AND c.FLOCALEID = 2052
WHERE   b.FNAME = '国家';
-- 查询辅助资料
SELECT  a.FID AS 辅助资料类型内码 ,
        c.FNAME AS 辅助资料类型名称 ,
        a.FENTRYID AS 辅助资料内码 ,
        a.FNUMBER AS 辅助资料编码 ,
        b.FDATAVALUE AS 辅助资料名称
FROM    T_BAS_ASSISTANTDATAENTRY a
        LEFT JOIN T_BAS_ASSISTANTDATAENTRY_L b ON a.FENTRYID = b.FENTRYID
                                                  AND b.FLOCALEID = 2052
        LEFT JOIN T_BAS_ASSISTANTDATA_L c ON a.FID = c.FID
                                             AND c.FLOCALEID = 2052
WHERE   c.FNAME = '国家';
```

## 4.7 弹性域

### 4.7.1 仓位

```sql
-- 仓库
SELECT b.FNAME AS 仓库名称,a.FNUMBER AS 仓库编码,a.FSTOCKID AS 仓库内码,* 
FROM T_BD_STOCK a
LEFT JOIN T_BD_STOCK_L b ON a.FSTOCKID=b.FSTOCKID AND b.FLOCALEID=2052

-- 仓库已启用仓位值集
SELECT b.FNAME AS 仓库名称,d.FNAME AS 仓位值集名称,* 
FROM T_BD_STOCK a
LEFT JOIN T_BD_STOCK_L b ON a.FSTOCKID=b.FSTOCKID AND b.FLOCALEID=2052
LEFT JOIN T_BD_STOCKFLEXITEM c ON a.FSTOCKID=c.FSTOCKID
LEFT JOIN T_BAS_FLEXVALUES_L d ON c.FFLEXID=d.FID AND d.FLOCALEID=2052
WHERE b.FNAME='菲律宾1库'

-- 仓位值集类型
SELECT b.FNAME AS 仓位值集名称,a.FNUMBER AS 仓位值集编码,a.FID AS 仓位值集内码,* 
FROM T_BAS_FLEXVALUES a
LEFT JOIN T_BAS_FLEXVALUES_L b ON a.FID=b.FID AND b.FLOCALEID=2052

-- 仓位值集维度数据
SELECT * FROM T_BAS_FLEXVALUESDETAIL
```

# 5. 工作流程类

## 5.1 待办，已办，转发，加签，委托

```sql
DECLARE @FUSERID INT
SELECT @FUSERID=FUSERID FROM dbo.T_SEC_USER WHERE FNAME='demo'

-- 我的待办
SELECT  t0.FSENDERID AS FSENDERID_Id ,
        t0.FTITLE AS FTITLE ,
        t0.FCREATETIME AS FCREATETIME ,
        t0_L.fpreresultname AS FPRERESULTNAME ,
        t0.fpredisposition AS FPREDISPOSITION ,
        t0.fbillnumber AS FBILLNUMBER ,
        t0.FPROCDEFID AS FProcDefId_Id ,
        t0.FCONTENT AS FCONTENT ,
        t0.FNUMBER AS FNUMBER ,
        t0.fdiscusscount AS FDISCUSSCOUNT ,
        t0.freadstatus AS FREADSTATUS ,
        t0.FASSIGNID AS FASSIGNID ,
        ROW_NUMBER() OVER ( ORDER BY t0.FASSIGNID DESC ) fidentityid
FROM    V_WF_ASSIGN t0
        LEFT JOIN V_WF_ASSIGN_L t0_L ON ( t0.FASSIGNID = t0_L.FASSIGNID
                                          AND t0_L.FLOCALEID = 2052
                                        )
WHERE   ( (t0.FRECEIVERID = @FUSERID
          AND t0.fstatus = '0')
        );
    
-- 我的已办
SELECT  t0.FSENDERID AS FSENDERID_Id ,
        t0.FTITLE AS FTITLE ,
        t0.FCREATETIME AS FCREATETIME ,
        t0_L.fpreresultname AS FPRERESULTNAME ,
        t0.FDISPOSITION AS FPREDISPOSITION ,
        t0.FPROCDEFID AS FProcDefId_Id ,
        t0.FBILLNO AS FBILLNO ,
        t0.FNUMBER AS FNUMBER ,
        t0.FCOMPLETEDTIME AS FCOMPLETEDTIME ,
        t0.fdiscusscount AS FDISCUSSCOUNT ,
        t0.FREADSTATUS AS FREADSTATUS ,
        t0.FASSIGNID AS FASSIGNID ,
        ROW_NUMBER() OVER ( ORDER BY t0.FCOMPLETEDTIME DESC ) fidentityid
FROM    V_WF_ASSIGNEND t0
        LEFT JOIN V_WF_ASSIGNEnd_L t0_L ON ( t0.FASSIGNID = t0_L.FAssignId
                                             AND t0_L.FLocaleId = 2052
                                           )
WHERE   ( (t0.FRECEIVERID = @FUSERID
          AND ( t0.FSTATUS > '0' ))
        );
 

-- 我的转发
SELECT  t0.FSenderId AS FSENDERID_Id ,
        t0.FCreateTime AS FCREATETIME ,
        t0.FFORWARDRECEIVERID AS FFORWARDRECEIVERID_Id ,
        t0.FTitle AS FTITLE ,
        t0.FCompletedTime AS FCOMPLETEDTIME ,
        t0.FProcDefId AS FProcDefId_Id ,
        t0.FDisposition AS FDISPOSITION ,
        t0.FNumber AS FNUMBER ,
        t0.fforwardstatus AS FFORWARDSTATUS ,
        t0.FReadStatus AS FREADSTATUS ,
        t0.FAssignId AS FASSIGNID ,
        ROW_NUMBER() OVER ( ORDER BY t0.FCreateTime DESC ) fidentityid
FROM    V_WF_FORWARD t0
WHERE   ( (FReceiverId = @FUSERID
          AND ( t0.FFORWARDRECEIVERID > 0 ))
        );

-- 我的加签
SELECT  t0.FSenderId AS FSENDERID_Id ,
        t0.FCreateTime AS FCREATETIME ,
        t0.FTitle AS FTITLE ,
        t0.FCompletedTime AS FCOMPLETEDTIME ,
        t0.FProcDefId AS FProcDefId_Id ,
        t0.FNumber AS FNUMBER ,
        t0.FReadStatus AS FREADSTATUS ,
        t0.FAssignId AS FASSIGNID ,
        ROW_NUMBER() OVER ( ORDER BY t0.FCreateTime DESC ) fidentityid
FROM    V_WF_ADDSIGN t0
WHERE   ( t0.FReceiverId = 100008 );
 

-- 我的委托
 SELECT t0.FCREATETIME AS FCREATETIME ,
        t0.fentrustreceiverid AS FReceiverId_Id ,
        t0.FTITLE AS FTITLE ,
        t0.FCONTENT AS FCONTENT ,
        t0_L.fpreresultname AS FPRERESULTNAME ,
        t0.FDISPOSITION AS FPREDISPOSITION ,
        t0.FNUMBER AS FNUMBER ,
        t0.fstatus AS FSTATUS ,
        t0.FREADSTATUS AS FREADSTATUS ,
        t0.fpkid AS FPKID ,
        ROW_NUMBER() OVER ( ORDER BY t0.FCREATETIME DESC ) fidentityid
 FROM   V_WF_Consign t0
        LEFT JOIN V_WF_Consign_L t0_L ON ( t0.fpkid = t0_L.fpkid
                                           AND t0_L.FLocaleId = 2052
                                         )
 WHERE  ( (( fentrustconsignor = @FUSERID
             AND fentrustconsignstatus = 1--0非委托1委托2已撤销
           )
          AND ( t0.fstatus <> '3' ))--0未处理1已处理3已转发
        );

-- 更新待办任务状态 待办任务状态：0未处理，1已处理
UPDATE T_WF_ASSIGN SET FSTATUS='1' WHERE FASSIGNID='636df1cc33b34c'

-- 流程实例状态：1完成,2运行中,3挂起,4终止,5暂停
UPDATE T_WF_PROCINST SET FSTATUS='1' WHERE FPROCINSTID='636df15289e012e65e'
```

## 5.2 查询已超时的代办任务

```sql
-- 查询已超时的待办任务（对应执行计划中的任务【工作流超时处理服务】）
-- 1. 未启用工作流日历
SELECT  A.FAPPROVALASSIGNID ,
        A.FASSIGNID ,
        A.FOBJECTTYPEID ,
        A.FKEYVALUE ,
        A.FPOLICYXML ,
        A.FAPPROVALACTION ,
        A.FVARIABLENAME ,
        A.FDEFAULTRESULT ,
        A.FRESULT ,
        A.FISLOCKED ,
        A.FLOCKEDUSERID ,
        A.FISALLOWLOCK ,
        A.FISALLOWFORWARD ,
        A.FISALLOWBATCHDEAL ,
        A.FCANSPECIFYNEXT ,
        A.FTIMEOUTDEAL ,
        A.FTIMEOUTHOURS ,
        A.FTIMEOUTSTRATEGY ,
        A.FAUTOAPPROVAL ,
        A.FTIMEOUTDEALWITH
FROM    T_WF_APPROVALASSIGN A
        INNER JOIN ( SELECT T.FAPPROVALASSIGNID
                     FROM   T_WF_APPROVALASSIGN T
                            INNER JOIN T_WF_ASSIGN B ON B.FASSIGNID = T.FASSIGNID
                     WHERE  B.FSTATUS = '0'
                            AND T.FTIMEOUTDEAL = '1'
                            AND T.FDEALTIME < GETDATE()
                            - ( 1.0 * T.FTIMEOUTHOURS / 24 + 1.0
                                * T.FTIMEOUTMINUTE / 60 / 24 )
                            AND EXISTS ( SELECT 1
                                         FROM   T_WF_PROCINST C
                                         WHERE  B.FPROCINSTID = C.FPROCINSTID
                                                AND C.FSTATUS IN ( '0', '2' ) )
                   ) U ON U.FAPPROVALASSIGNID = A.FAPPROVALASSIGNID

-- 启用了工作流日历
SELECT  A.FAPPROVALASSIGNID ,
        A.FASSIGNID ,
        A.FOBJECTTYPEID ,
        A.FKEYVALUE ,
        A.FPOLICYXML ,
        A.FAPPROVALACTION ,
        A.FVARIABLENAME ,
        A.FDEFAULTRESULT ,
        A.FRESULT ,
        A.FISLOCKED ,
        A.FLOCKEDUSERID ,
        A.FISALLOWLOCK ,
        A.FISALLOWFORWARD ,
        A.FISALLOWBATCHDEAL ,
        A.FCANSPECIFYNEXT ,
        A.FTIMEOUTDEAL ,
        A.FTIMEOUTHOURS ,
        A.FTIMEOUTSTRATEGY ,
        A.FAUTOAPPROVAL ,
        A.FTIMEOUTDEALWITH
FROM    T_WF_APPROVALASSIGN A
        INNER JOIN ( SELECT T.FAPPROVALASSIGNID
                     FROM   T_WF_APPROVALASSIGN T
                            INNER JOIN T_WF_ASSIGN B ON B.FASSIGNID = T.FASSIGNID
                     WHERE  B.FSTATUS = '0'
                            AND T.FTIMEOUTDEAL = '1'
                            AND ( T.FCALENDARID > 0
                                  AND EXISTS ( SELECT   1
                                               FROM     T_WF_CALENDAR tc
                                               WHERE    tc.FID = T.FCALENDARID
                                                        AND tc.FFORBIDSTATUS = 'A' )
                                )
                            AND EXISTS ( SELECT 1
                                         FROM   T_WF_PROCINST C
                                         WHERE  B.FPROCINSTID = C.FPROCINSTID
                                                AND C.FSTATUS IN ( '0', '2' ) )
                            AND ( ISNULL(FDEALDATEHOURS, 0) --当天提交单据时剩余工作小时数（比如当天工作8小时，8点上班，9点提交单据，那么剩余工作小时是7小时）
                                  + ( SELECT    ISNULL(SUM(FTEMPCOSTHOURS), 0)
                                      FROM      T_WF_CALENDAR
                                      WHERE     FID = T.FCALENDARID
                                    )     -- 截止当前当天已耗工时（比如当天工作8小时，8点上班，10点的时候去查询是否超时，那么当前已使用工作小时数是2小时）
                                  - ( CASE WHEN DATEDIFF(DD, FDEALTIME,
                                                         GETDATE()) <= 0
                                           THEN ( SELECT    ISNULL(SUM(FWORKHOURS),
                                                              0)
                                                  FROM      T_WF_CALENDARDATA
                                                  WHERE     FID = T.FCALENDARID
                                                            AND FISWORKTIME = '1'
                                                            AND DATEDIFF(DD,
                                                              FDAY, GETDATE()) = 0
                                                )
                                           ELSE 0
                                      END-- 如果提交日期和当前查询日期是同一天，还需要减去当天工作小时数（按上述例子，当天工作8小时，7小时+2小时-8小时=1小时，即得出已耗工作时间）
                                           )
                                  + ( SELECT    ISNULL(SUM(FWORKHOURS), 0)
                                      FROM      T_WF_CALENDARDATA
                                      WHERE     FID = T.FCALENDARID
                                                AND FISWORKTIME = '1'
                                                AND FDAY > FDEALTIME
                                                AND FDAY <= GETDATE() - 1
                                    )-- 如果提交日期和当前查询日期不是同一天，还需要加上期间所有工作小时数，即得出已耗工作时间
                                  ) * 1.0 / 24 > ( 1.0 * T.FTIMEOUTHOURS / 24
                                                   + 1.0 * T.FTIMEOUTMINUTE
                                                   / 60 / 24 )
                     UNION
                     SELECT T.FAPPROVALASSIGNID
                     FROM   T_WF_APPROVALASSIGN T
                            INNER JOIN T_WF_ASSIGN B ON B.FASSIGNID = T.FASSIGNID
                     WHERE  B.FSTATUS = '0'
                            AND T.FTIMEOUTDEAL = '1'
                            AND T.FDEALTIME < GETDATE() - ( 1.0
                                                            * T.FTIMEOUTHOURS
                                                            / 24 + 1.0
                                                            * T.FTIMEOUTMINUTE
                                                            / 60 / 24 )
                            AND EXISTS ( SELECT 1
                                         FROM   T_WF_PROCINST C
                                         WHERE  B.FPROCINSTID = C.FPROCINSTID
                                                AND C.FSTATUS IN ( '0', '2' ) )
                            AND ( T.FCALENDARID <= 0
                                  OR EXISTS ( SELECT    1
                                              FROM      T_WF_CALENDAR tc
                                              WHERE     tc.FID = T.FCALENDARID
                                                        AND tc.FFORBIDSTATUS = 'B' )
                                )
                   ) U ON U.FAPPROVALASSIGNID = A.FAPPROVALASSIGNID
```

## 5.3 终止流程实例

```sql
-- 单据启动工作流后，单据数据被非法删除，导致该单据的流程实例没法终止了。
-- 查询采购订单，单据已被删，但是待办任务还未终止的流程实例编码
SELECT  a.FNUMBER AS 流程实例编码
FROM    T_WF_PROCINST a
        JOIN T_WF_PIBIMAP b ON a.FPROCINSTID = b.FPROCINSTID AND b.FOBJECTTYPEID = 'PUR_PurchaseOrder'
        LEFT JOIN t_PUR_POOrder c ON b.FKEYVALUE = c.FID
WHERE   a.FCREATETIME > '2000-9-1'
        AND c.FID IS NULL
-- 将异常的流程实例设置为终止（流程实例状态：1完成2运行中3挂起4终止5暂停）   
UPDATE  T_WF_PROCINST
SET     FSTATUS = 4--终止流程实例
WHERE   FNUMBER IN (
        SELECT  a.FNUMBER AS 流程实例编码
        FROM    T_WF_PROCINST a
                JOIN T_WF_PIBIMAP b ON a.FPROCINSTID = b.FPROCINSTID
                                       AND b.FOBJECTTYPEID = 'PUR_PurchaseOrder'
                LEFT JOIN t_PUR_POOrder c ON b.FKEYVALUE = c.FID
        WHERE   a.FCREATETIME > '2000-9-1'
                AND c.FID IS NULL )
        AND FSTATUS = 2

-- 终止付款申请单的异常流程实例：
-- 查询付款申请单，单据已被删，但是待办任务还未终止的流程实例编码
SELECT  a.FNUMBER AS 流程实例编码
FROM    T_WF_PROCINST a
        JOIN T_WF_PIBIMAP b ON a.FPROCINSTID = b.FPROCINSTID AND b.FOBJECTTYPEID = 'CN_PAYAPPLY'
        LEFT JOIN T_CN_PAYAPPLY c ON b.FKEYVALUE = c.FID
WHERE   a.FCREATETIME > '2022-9-1'
        AND c.FID IS NULL
      
-- 将异常的流程实例设置为终止（流程实例状态：1完成2运行中3挂起4终止5暂停）   
UPDATE  T_WF_PROCINST
SET     FSTATUS = 4--终止流程实例
WHERE   FNUMBER IN (
        SELECT  a.FNUMBER AS 流程实例编码
        FROM    T_WF_PROCINST a
                JOIN T_WF_PIBIMAP b ON a.FPROCINSTID = b.FPROCINSTID
                                       AND b.FOBJECTTYPEID = 'CN_PAYAPPLY'
                LEFT JOIN T_CN_PAYAPPLY c ON b.FKEYVALUE = c.FID
        WHERE   a.FCREATETIME > '2022-9-1'
                AND c.FID IS NULL )
        AND FSTATUS = 2
```

## 5.4 查询节点类型

```sql
-- <1>通过待办任务内码，查询待办任务的节点信息
SELECT
    a.FASSIGNID AS 待办任务内码,
    a.FACTINSTID AS 节点内码,
    b.FACTIVITYID AS 节点ID,
    c.FACTNAME AS 节点名称
FROM
    T_WF_ASSIGN a
    JOIN T_WF_ACTINST b ON a.FACTINSTID = b.FACTINSTID
    LEFT JOIN T_WF_ACTINST_L c ON b.FACTINSTID = c.FACTINSTID
    AND c.FLOCALEID = 2052
WHERE
    a.FASSIGNID = '005056af13b388b411e4d685714b81f4';

-- <2>通过待办任务内码和节点ID，查询节点XML，通过节点名称判断当前节点是什么类型的节点                                
SELECT
    CASE
        e.FPROCESSTYPE
        WHEN 300 THEN '工作流'
        ELSE '审批流'
    END AS 流程类型,
    d.FDEFINITIONXML.query('//*/Activities/*/ActivityId[.="22"]/parent::*') AS 节点XML -- 这里的22是节点ID，需手动替换
FROM
    T_WF_ASSIGN a
    JOIN T_WF_ACTINST b ON a.FACTINSTID = b.FACTINSTID
    JOIN T_WF_PROCINST c ON a.FPROCINSTID = c.FPROCINSTID
    JOIN T_WF_DEFVERSION d ON c.FVERSIONID = d.FVERSIONID
    JOIN T_WF_PROCDEF e ON d.FPROCDEFID = e.FPROCDEFID
WHERE
    a.FASSIGNID = '005056af13b388b411e4d685714b81f4';

-- 工作流节点类型：
--StartActivity：开始节点
--NativeActivity：普通节点
-- 审批流节点类型：
--ApprovalStartActivity：开始节点
--ApprovalEndActivity：结束节点
--ApprovalActivity：审批节点
--FinalApprovalActivity：终审节点
--ApprovalMessageActivity：消息节点
--InvokeOperationActivity：自动服务节点
--ManualOperationActivity：人工服务节点
--ApprovalShareServiceActivity：共享服务节点
--ApprovalImageInteractActivity：影像节点
```

## 5.5 查询流程模板

```sql
SELECT
    c.FNAME AS 业务对象名称,
    a.FNUMBER AS 流程模板编码,
    b.FNAME AS 流程模板名称,
    e.FNUMBER AS 流程编码,
    f.FDISPLAYNAME AS 流程名称,
    g.FVERSIONID AS 流程版本内码,
    g.FVERSION AS 流程版本号,
    i.FACTIVITYID AS 流程节点标识,
    j.FNAME AS 流程节点绑定的视图名称,
    i.FACTIVITYCONFIG AS 节点界面布局XML
FROM
    T_WF_TEMPLATE a --流程模板
    LEFT JOIN T_WF_TEMPLATE_L b ON a.FTMPID = b.FTMPID
    AND b.FLOCALEID = 2052
    LEFT JOIN T_META_OBJECTTYPE_L c ON a.FFORMID = c.FID
    AND c.FLOCALEID = 2052
    LEFT JOIN T_WF_PRCLAUNCHCONFIG d ON a.FTMPID = d.FTMPID --流程模板明细
    LEFT JOIN T_WF_PROCDEF e ON d.FPROCDEFID = e.FPROCDEFID --流程定义
    LEFT JOIN T_WF_PROCDEF_L f ON e.FPROCDEFID = f.FPROCDEFID
    AND f.FLOCALEID = 2052
    LEFT JOIN T_WF_DEFVERSION g ON d.FVERSIONID = g.FVERSIONID --流程版本
    LEFT JOIN T_WF_PRCCUSTOMIZEDCONFIG h ON d.FTMPDETAILID = h.FTMPDETAILID --流程节点界面配置
    LEFT JOIN T_WF_ACTIVITYCONFIG i ON h.FTMPVERDETAILID = i.FTMPVERDETAILID --流程节点界面配置明细
    LEFT JOIN T_META_OBJECTTYPEVIEW_L j ON CAST(
        i.FACTIVITYCONFIG.query('//ElementControlCollection/LayoutId/text()') AS VARCHAR(100)
    ) = j.FID
    AND j.FLOCALEID = 2052
WHERE
    a.FFORMID = 'PUR_PurchaseOrder';
```

## 5.6 通过代办任务内码查询流程发起人

```sql
SELECT
    b.FORIGINATORID AS 发起人内码,
    c.FNAME AS 发起人名称
FROM
    T_WF_ASSIGN a
    JOIN T_WF_PROCINST b ON a.FPROCINSTID = b.FPROCINSTID
    JOIN T_SEC_USER c ON b.FORIGINATORID = c.FUSERID
WHERE
    a.FASSIGNID = 'b_60af3f26e0ca06';
```

## 5.7 查询多级审核单据体表名

```sql
SELECT a.FFORMID  AS 业务对象标识,b.FNAME 业务对象名称,a.FMATABNAME AS 多级审核单据体表名 
FROM T_WF_MULTIAPPROVALDESC a
LEFT JOIN T_META_OBJECTTYPE_L b ON a.FFORMID=b.FID AND b.FLOCALEID=2052
WHERE b.FNAME='采购订单'
```

## 5.8 查询流程信息

```sql
DECLARE @FormId VARCHAR(50)
DECLARE @FBillNo NVARCHAR(50)
SET @FormId='PUR_PurchaseOrder'
SET @FBillNo='CGDD002129'

---- 查询元数据信息
--SELECT a.FID AS 业务对象标识,b.FNAME AS 业务对象名称,CAST(FKERNELXML.query('//HeadEntity/TableName/text()') AS VARCHAR) AS 主表名,a.* 
--FROM T_META_OBJECTTYPE a
--LEFT JOIN T_META_OBJECTTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052
--WHERE a.FMODELTYPEID<=400 AND b.FNAME LIKE '%采购订单%'
---- AND (a.FID LIKE '%PUR_PurchaseOrder%' OR FINHERITPATH LIKE '%PUR_PurchaseOrder%')

-- 查询单据信息
DECLARE @FID INT
SELECT @FID=FID FROM T_PUR_POORDER WHERE FBILLNO=@FBillNo
SELECT '单据信息'
,a.FBILLNO AS 单据编号
,(CASE a.FDOCUMENTSTATUS WHEN 'A' THEN '创建' WHEN 'B' THEN '审核中' WHEN 'C' THEN '已审核' WHEN 'D' THEN '重新审核' WHEN 'Z' THEN '暂存' ELSE '' END) AS 单据状态
,b.FNAME AS 最后审批人,FAPPROVEDATE AS 最后审批时间,c.FNAME AS 最后修改人,a.FMODIFYDATE AS 最后修改时间,a.* 
FROM T_PUR_POORDER a
LEFT JOIN T_SEC_USER b ON a.FAPPROVERID=b.FUSERID
LEFT JOIN T_SEC_USER c ON a.FMODIFIERID=c.FUSERID
WHERE a.FBILLNO=@FBillNo

-- 查询流程实例信息
DECLARE @FPROCINSTID VARCHAR(50)
SELECT TOP 1 @FPROCINSTID=FPROCINSTID FROM T_WF_PROCINST WHERE FNUMBER LIKE @FBillNo+'%' --AND FSTATUS='4'
ORDER BY FPROCINSTID DESC
SELECT '流程实例信息(T_WF_PROCINST)'
,(CASE a.FSTATUS WHEN '1' THEN '完成' WHEN '2' THEN '运行中' WHEN '3' THEN '挂起' WHEN '4' THEN '终止' WHEN '5' THEN '暂停' ELSE '' END)  AS 流程状态
,a.FNUMBER AS 流程实例编码,a.FPROCINSTID AS 流程实例内码,b.FNUMBER AS 流程编码,c.FDISPLAYNAME AS 流程名称
,d.FVERSION AS 流程版本,d.FDEFINITIONXML AS 流程定义XML,d.FDESIGNERXML AS 流程设计XML
,a.FFAILUREOUTINFO AS 错误信息,a.FFAILUREMESSAGE AS 错误堆栈
,a.* 
FROM T_WF_PROCINST a 
LEFT JOIN T_WF_PROCDEF b ON a.FPROCDEFID=b.FPROCDEFID
LEFT JOIN T_WF_PROCDEF_L c ON b.FPROCDEFID=c.FPROCDEFID AND c.FLOCALEID=2052
LEFT JOIN T_WF_DEFVERSION d ON a.FVERSIONID=d.FVERSIONID
WHERE a.FNUMBER LIKE @FBillNo+'%'
ORDER BY a.FPROCINSTID DESC

-- 查询流程中转消息（此表有数据时需要高度关注）
SELECT * FROM T_WF_MAPSTATEMQSTORE

-- 查询流程节点信息
SELECT '流程节点信息(T_WF_ACTINST)' 
,b.FACTNAME AS 节点名称
,CASE a.FSTATUS WHEN '0' THEN '准备(未处理)' WHEN '1' THEN '完成(已处理)' WHEN '2' THEN '空闲/挂起(运行中)' WHEN '3' THEN '被舍弃' WHEN '4' THEN '失败' WHEN '5' THEN '终止' WHEN '6' THEN '驳回' ELSE '' END AS 节点状态
,a.FACTINSTID AS 节点内码,a.* 
FROM T_WF_ACTINST a
LEFT JOIN T_WF_ACTINST_L b ON a.FACTINSTID=b.FACTINSTID AND b.FLOCALEID=2052
WHERE a.FPROCINSTID=@FPROCINSTID
ORDER BY a.FACTINSTID

-- 查询流程节点实例消息
SELECT '流程节点实例消息(T_WF_ACTINSTMSG)'
,d.FACTNAME AS 节点名称
,e.FNAME AS 接收人
,a.FCREATETIME AS 接收时间
,a.*,c.FACTINSTID
FROM T_WF_ACTINSTMSG a
LEFT JOIN T_WF_PROCINST b ON a.FPROCINSTID=b.FPROCINSTID
LEFT JOIN T_WF_ACTINST c ON b.FPROCINSTID=c.FPROCINSTID AND a.FACTIVITYID=c.FACTIVITYID
LEFT JOIN T_WF_ACTINST_L d ON c.FACTINSTID=d.FACTINSTID AND d.FLOCALEID=2052
LEFT JOIN T_SEC_USER e ON a.FRECEIVERID=e.FUSERID
WHERE a.FPROCINSTID=@FPROCINSTID
ORDER BY a.FACTINSTMSGID

-- 查询节点动作执行日志信息
SELECT '节点动作执行日志信息(T_WF_ACTIONEVENTLOG)'
,c.FACTNAME AS 节点名称
,a.FEXECUTEDTIME AS 执行时间
,a.FEXECUTELOG AS 执行日志
,a.* 
FROM T_WF_ACTIONEVENTLOG a
LEFT JOIN  T_WF_ACTINST b ON a.FACTINSTID=b.FACTINSTID
LEFT JOIN T_WF_ACTINST_L c ON b.FACTINSTID=c.FACTINSTID AND c.FLOCALEID=2052
WHERE b.FPROCINSTID=@FPROCINSTID
ORDER BY a.FACTIONLOGID

-- 查询待办任务信息
SELECT '待办任务信息(T_WF_ASSIGN)'
,c.FACTNAME AS 节点名称
,CASE a.FSTATUS WHEN '0' THEN '未处理' WHEN '1' THEN '已处理' ELSE '' END AS 待办任务状态
,b.FNAME AS 发送人
,a.FRECEIVERNAMES AS 接收人_未处理用户
,a.FASSIGNID AS 待办任务内码,a.* 
FROM T_WF_ASSIGN a
LEFT JOIN T_SEC_USER b ON a.FSENDERID=b.FUSERID
LEFT JOIN T_WF_ACTINST_L c ON a.FACTINSTID=c.FACTINSTID AND c.FLOCALEID=2052
WHERE a.FPROCINSTID=@FPROCINSTID
ORDER BY a.FASSIGNID

-- 查询待办任务接收人信息
SELECT '待办任务接收人信息(T_WF_RECEIVER)'
,c.FACTNAME AS 节点名称
,d.FNAME AS 接收人
,CASE a.FADDSIGNSTATUS WHEN '0' THEN '不涉及加签，计入审批票数' WHEN '1' THEN '作为加签发起人，其发起的加签任务组全部已完成，计入审批' WHEN '2' THEN '作为加签发起人，其发起的加签任务组未全部完成，不计入审批' WHEN '4' THEN '作为加签接收人，不计入审批票数' ELSE '' END AS 加签状态
,a.* 
FROM T_WF_RECEIVER a 
LEFT JOIN T_WF_ASSIGN b ON a.FASSIGNID=b.FASSIGNID
LEFT JOIN T_WF_ACTINST_L c ON b.FACTINSTID=c.FACTINSTID AND c.FLOCALEID=2052
LEFT JOIN T_SEC_USER d ON a.FRECEIVERID=d.FUSERID
WHERE b.FPROCINSTID=@FPROCINSTID 
ORDER BY a.FID

-- 查询审批任务信息
SELECT '审批任务信息(T_WF_APPROVALASSIGN)'
,c.FACTNAME AS 节点名称
,a.* 
FROM T_WF_APPROVALASSIGN a
LEFT JOIN T_WF_ASSIGN b ON a.FASSIGNID=b.FASSIGNID
LEFT JOIN T_WF_ACTINST_L c ON b.FACTINSTID=c.FACTINSTID AND c.FLOCALEID=2052
WHERE b.FPROCINSTID=@FPROCINSTID
ORDER BY a.FAPPROVALASSIGNID

--  查询审批任务处理明细信息
SELECT '审批任务处理明细信息(T_WF_APPROVALITEM)' AS 审批任务处理明细信息
,e.FACTNAME AS 节点名称
,d.FNAME AS 接收人
,CASE a.FSTATUS WHEN '0' THEN '未处理' WHEN '1' THEN '已处理' WHEN '3' THEN '转发' WHEN '4' THEN '准备中(顺序会签)' ELSE '' END AS 处理状态
,a.* 
FROM T_WF_APPROVALITEM a 
LEFT JOIN T_WF_APPROVALASSIGN b ON a.FAPPROVALASSIGNID=b.FAPPROVALASSIGNID
LEFT JOIN T_WF_ASSIGN c ON b.FASSIGNID=c.FASSIGNID
LEFT JOIN T_SEC_USER d ON a.FRECEIVERID=d.FUSERID
LEFT JOIN T_WF_ACTINST_L e ON c.FACTINSTID=e.FACTINSTID AND e.FLOCALEID=2052
WHERE c.FPROCINSTID=@FPROCINSTID
ORDER BY a.FAPPROVALITEMID

SELECT '待办任务加签信息(T_WF_ADDSIGNASSIGN)'
,a.FADDSIGNASSIGNID AS 加签任务内码
,CASE a.FSTATUS WHEN '0' THEN '未处理' WHEN '1' THEN '已处理' ELSE '' END AS 加签任务状态
,e.FACTNAME AS 节点名称
,d.FNAME AS 待办发送人
,c.FRECEIVERNAMES AS 待办接收人
,f.FNAME AS 加签任务创建人
,CASE a.FADDSIGNTYPE WHEN 0 THEN '前加签' WHEN 1 THEN '后加签' ELSE '' END AS 加签类型
,g.FNAME AS 加签任务接收人
,CASE b.FSTATUS WHEN '0' THEN '未处理' WHEN '1' THEN '已处理' ELSE '' END AS 加签处理状态
,b.FACTIONRESULT AS 加签审批结果
,a.* 
FROM T_WF_ADDSIGNASSIGN a 
LEFT JOIN T_WF_ADDSIGNRECEIVER b ON a.FADDSIGNASSIGNID=b.FADDSIGNASSIGNID
LEFT JOIN T_WF_ASSIGN c ON a.FASSIGNID=c.FASSIGNID
LEFT JOIN T_SEC_USER d ON c.FSENDERID=d.FUSERID
LEFT JOIN T_WF_ACTINST_L e ON c.FACTINSTID=e.FACTINSTID AND e.FLOCALEID=2052
LEFT JOIN T_SEC_USER f ON a.FCREATORID=f.FUSERID
LEFT JOIN T_SEC_USER g ON b.FRECEIVERID=g.FUSERID
WHERE c.FPROCINSTID=@FPROCINSTID AND ((a.FADDSIGNTYPE=0 AND a.FISORIGIN=1) OR (a.FADDSIGNTYPE=1 AND a.FISORIGIN=0))
ORDER BY a.FASSIGNID,a.FADDSIGNASSIGNID

-- 查询消息
SELECT '消息',cl.FACTNAME AS 节点名称,d.FNAME AS 接收人,e.FNAME AS 发送人,a.* FROM T_WF_MESSAGE a
LEFT JOIN T_WF_PROCINST b ON a.FPROCINSTID=b.FPROCINSTID
LEFT JOIN T_WF_ACTINST c ON a.FPROCINSTID=c.FPROCINSTID AND a.FACTIVITYID=c.FACTIVITYID
LEFT JOIN T_WF_ACTINST_L cl ON c.FACTINSTID=cl.FACTINSTID AND cl.FLOCALEID=2052
LEFT JOIN T_SEC_USER d ON a.FRECEIVERID=d.FUSERID
LEFT JOIN T_SEC_USER e ON a.FSENDERID=e.FUSERID
WHERE a.FPROCINSTID=@FPROCINSTID

-- 查询工作流日志信息
SELECT '工作流日志信息(T_WF_LOG)'
,c.FACTNAME AS 节点名称
,b.FNAME AS 操作人
,a.* 
FROM T_WF_LOG a 
LEFT JOIN T_SEC_USER b ON a.FCREATORID=b.FUSERID
LEFT JOIN T_WF_ACTINST_L c ON a.FACTINSTID=c.FACTINSTID AND c.FLOCALEID=2052
WHERE a.FPROCINSTID=@FPROCINSTID
ORDER BY a.FID

-- 查询上机操作日志
SELECT '上机操作日志(单据编号)' AS [上机操作日志(单据编号)]
,a.FCLIENTTYPE AS 客户端类型,b.FNAME AS 用户名,a.FDATETIME AS 操作时间,a.FDESCRIPTION AS 操作描述,a.* 
FROM T_BAS_OPERATELOG a 
LEFT JOIN T_SEC_USER b ON a.FUSERID=b.FUSERID
WHERE a.FOBJECTTYPEID=@FormId AND a.FDESCRIPTION LIKE '%'+@FBillNo+'%'
ORDER BY a.FDATETIME

-- 查询已归档上机操作日志
SELECT '已归档上机操作日志(单据编号)' AS [已归档上机操作日志(单据编号)]
,a.FCLIENTTYPE AS 客户端类型,b.FNAME AS 用户名,a.FDATETIME AS 操作时间,a.FDESCRIPTION AS 操作描述,a.* 
FROM T_BAS_OPERATELOGBK a 
LEFT JOIN T_SEC_USER b ON a.FUSERID=b.FUSERID
WHERE a.FOBJECTTYPEID=@FormId AND a.FDESCRIPTION LIKE '%'+@FBillNo+'%'
ORDER BY a.FDATETIME

---- 查询上机操作日志
--SELECT '上机操作日志(单据内码)' AS [上机操作日志(单据内码)]
--,a.FCLIENTTYPE AS 客户端类型,b.FNAME AS 用户名,a.FDATETIME AS 操作时间,a.FDESCRIPTION AS 操作描述,a.* 
--FROM T_BAS_OPERATELOG a 
--LEFT JOIN T_SEC_USER b ON a.FUSERID=b.FUSERID
--WHERE a.FOBJECTTYPEID=@FormId AND a.FINTERID=CAST(@FID AS VARCHAR(50))
--ORDER BY FDATETIME

---- 查询已归档上机操作日志
--SELECT '已归档上机操作日志(单据内码)' AS [已归档上机操作日志(单据内码)]
--,a.FCLIENTTYPE AS 客户端类型,b.FNAME AS 用户名,a.FDATETIME AS 操作时间,a.FDESCRIPTION AS 操作描述,a.* 
--FROM T_BAS_OPERATELOGBK a 
--LEFT JOIN T_SEC_USER b ON a.FUSERID=b.FUSERID
--WHERE a.FOBJECTTYPEID=@FormId AND a.FINTERID=CAST(@FID AS VARCHAR(50))
--ORDER BY FDATETIME

-- 查询元数据信息
SELECT '元数据信息' AS 元数据信息
,a.FID AS 业务对象标识,b.FNAME AS 业务对象名称,CAST(FKERNELXML.query('//HeadEntity/TableName/text()') AS VARCHAR) AS 主表名,a.* 
FROM T_META_OBJECTTYPE a
LEFT JOIN T_META_OBJECTTYPE_L b ON a.FID=b.FID AND b.FLOCALEID=2052
WHERE a.FMODELTYPEID<=400 AND (a.FID LIKE '%'+@FormId+'%' OR FINHERITPATH LIKE '%'+@FormId+'%')

---- 查询单据已审核流程未结束的异常单据
--SELECT '单据已审核流程未结束的异常单据' AS 单据已审核流程未结束的异常单据
--,c.FNUMBER AS 流程实例编码,c.FPROCINSTID AS 流程实例内码,a.FBILLNO AS 单号,a.* 
--FROM T_PUR_POORDER a
--JOIN T_WF_PIBIMAP b ON b.FOBJECTTYPEID=@FormId AND a.FID=b.FKEYVALUE
--JOIN T_WF_PROCINST c ON b.FPROCINSTID=c.FPROCINSTID
--WHERE a.FDOCUMENTSTATUS = 'C' AND c.FSTATUS='2'

---- 流程已结束单据未审核的异常单据
--SELECT '流程已结束单据未审核的异常单据' AS 流程已结束单据未审核的异常单据
--,c.FNUMBER AS 流程实例编码,c.FPROCINSTID AS 流程实例内码,a.FBILLNO AS 单号,a.* 
--FROM T_PUR_POORDER a
--JOIN T_WF_PIBIMAP b ON b.FOBJECTTYPEID=@FormId AND a.FID=b.FKEYVALUE
--JOIN T_WF_PROCINST c ON b.FPROCINSTID=c.FPROCINSTID
--WHERE a.FDOCUMENTSTATUS <> 'C' AND c.FSTATUS<>'2'

-- 查询流程配置中心的节点界面配置
SELECT
    ml.FNAME AS 业务对象名称,
    ol.FNAME AS 组织名称,
    a.FTMPID AS 流程配置内码,
    a.FNUMBER AS 流程配置编码,
    al.FNAME AS 流程配置名称,
    f.FPROCDEFID AS 流程内码,
    f.FNUMBER AS 流程编码,
    fl.FDISPLAYNAME AS 流程名称,
    b.FVERSIONID AS 流程版本内码,
    e.FVERSION AS 流程版本,
    d.FACTIVITYID AS 节点ID,
    d.FACTIVITYCONFIG AS 节点配置XML,
    *
FROM
    T_WF_TEMPLATE a
    JOIN T_WF_PRCLAUNCHCONFIG b ON a.FTMPID = b.FTMPID
    LEFT JOIN T_WF_PRCCUSTOMIZEDCONFIG c ON b.FTMPDETAILID = c.FTMPDETAILID
    AND c.FVERSIONID IN (
        SELECT
            y.FVERSIONID
        FROM
            T_WF_DEFVERSION x
            INNER JOIN T_WF_DEFVERSION y ON (x.FPROCDEFID = y.FPROCDEFID)
        WHERE
            x.FVERSIONID = b.FVERSIONID
    )
    LEFT JOIN T_WF_ACTIVITYCONFIG d ON c.FTMPVERDETAILID = d.FTMPVERDETAILID
    LEFT JOIN T_WF_DEFVERSION e ON b.FVERSIONID = e.FVERSIONID
    LEFT JOIN T_WF_PROCDEF f ON b.FPROCDEFID = f.FPROCDEFID
    LEFT JOIN T_WF_TEMPLATE_L al ON a.FTMPID = al.FTMPID
    AND al.FLOCALEID = 2052
    LEFT JOIN T_WF_PROCDEF_L fl ON f.FPROCDEFID = fl.FPROCDEFID
    AND fl.FLOCALEID = 2052
    LEFT JOIN T_META_OBJECTTYPE_L ml ON a.FFORMID = ml.FID
    AND ml.FLOCALEID = 2052
    LEFT JOIN T_ORG_ORGANIZATIONS_L ol ON a.FORGID = ol.FORGID
WHERE
    a.FNUMBER = '币别工作流';
```

## 5.9 查询单据的运行中的流程实例的当前节点信息

```sql
SELECT
    c.FACTINSTID AS 节点实例内码,
    c.FACTIVITYID AS 节点ID,
    cl.FACTNAME AS 节点名称,
    c.FSTATUS AS 节点状态
FROM
    T_WF_PIBIMAP a
    JOIN T_WF_PROCINST b ON a.FPROCINSTID = b.FPROCINSTID
    JOIN T_WF_ACTINST c ON b.FPROCINSTID = c.FPROCINSTID
    JOIN T_WF_ACTINST_L cl ON c.FACTINSTID = cl.FACTINSTID
    AND cl.FLOCALEID = 2052
WHERE
    a.FOBJECTTYPEID = 'BD_Currency' --业务对象标识
    AND a.FKEYVALUE = '702525' -- 单据内码
    AND b.FSTATUS IN ('2', '3', '5') -- 流程状态[1完成2运行中3挂起4终止5暂停]
    AND c.FSTATUS IN ('0', '2'); --节点状态[0准备(未处理)1完成(已处理)2空闲/挂起(运行中)3被舍弃4失败5终止6驳回]
```

## 5.10 查询单据状态已审核待办任务未处理的业务单据数据

```sql
-- 查询采购订单和销售订单的单据状态已审核待办任务未处理的业务单据数据。
--<1>所有待办任务涉及的业务对象
SELECT
    DISTINCT FOBJECTTYPEID
FROM
    V_WF_ASSIGN
WHERE
    FSTATUS = 0 
--<2>基于步骤1查询到的业务对象，再次查询单据状态已审核，待办任务未处理的业务单据数据
SELECT
    b.FNAME AS 单据名称,
    a.FKEYVALUE 单据内码,
    a.FNUMBER AS 流程实例编码,
    c.FNAME AS 待办任务处理人
FROM
    V_WF_ASSIGN a
    LEFT JOIN T_META_OBJECTTYPE_L b ON a.FOBJECTTYPEID = b.FID
    AND b.FLOCALEID = 2052
    LEFT JOIN T_SEC_USER c ON a.FRECEIVERID = c.FUSERID
WHERE
    a.fstatus = 0
    AND (
        a.FOBJECTTYPEID = 'PUR_PurchaseOrder'
        AND EXISTS (
            SELECT
                *
            FROM
                T_PUR_POORDER
            WHERE
                FID = a.FKEYVALUE
                AND FDOCUMENTSTATUS = 'C'
        )
    ) --采购订单
    OR (
        a.FOBJECTTYPEID = 'SAL_SaleOrder'
        AND EXISTS (
            SELECT
                *
            FROM
                T_SAL_ORDER
            WHERE
                FID = a.FKEYVALUE
                AND FDOCUMENTSTATUS = 'C'
        )
    ); --销售订单

-- 查询某个业务对象下的单据是否存在未完成的流程实例
SELECT  *
FROM    T_WF_PIBIMAP a
        JOIN T_WF_PROCINST b ON a.FPROCINSTID = b.FPROCINSTID
WHERE   a.FOBJECTTYPEID = 'AP_Payable'--业务对象标识
        AND a.FKEYVALUE = '123456'--单据的主键
        AND b.FSTATUS NOT IN ( '1', '4' );--流程实例的状态不等于完成和终止
```

## 5.11 查询某张采购订单的所有流程实例内码。

```sql
-- 查询单据内码
DECLARE @FKEYVALUE INT;

SELECT
    @FKEYVALUE = FID
FROM
    T_PUR_POORDER
WHERE
    FBILLNO = 'CGDD002058';

-- 业务对象标识
DECLARE @FOBJECTTYPEID VARCHAR(36);

SET
    @FOBJECTTYPEID = 'PUR_PurchaseOrder';

-- 查询指定单据的所有流程实例
SELECT
    PIBI.FPROCINSTID,
    PROCINST.FSTATUS
FROM
    V_WF_PIBIMAPAll PIBI
    INNER JOIN V_WF_PROCINSTMG PROCINST ON PIBI.FPROCINSTID = PROCINST.FPROCINSTID
WHERE
    FOBJECTTYPEID = @FOBJECTTYPEID
    AND FKEYVALUE = @FKEYVALUE
ORDER BY
    PROCINST.FCREATETIME DESC;
```

## 5.12 查询单据在指定用户下的代办任务

```sql
/* 查询某张采购订单在demo用户下的待办任务信息，用于判断demo用户是否有权限进行审批操作，判断是用工作流审批还是传统审批。
 <1>查询当前单据的待办任务。
 如果能查到待办任务，表示接下来的审批操作将走工作流审批。
 如果查不到待办任务，需进一步分析该单是个什么情况。*/
-- 查询用户内码
DECLARE @FUSERID INT;

SELECT
    @FUSERID = FUSERID
FROM
    T_SEC_USER
WHERE
    FNAME = 'demo';

-- 查询单据内码
DECLARE @FKEYVALUE INT;

SELECT
    @FKEYVALUE = FID
FROM
    T_PUR_POORDER
WHERE
    FBILLNO = 'CGDD002058';

-- 业务对象标识
DECLARE @FOBJECTTYPEID VARCHAR(36);

SET
    @FOBJECTTYPEID = 'PUR_PurchaseOrder';

-- 查询指定单据是否存在指定用户的待办任务
-- 【注意】如果某个单据，既找不到某个用户的待办任务，又找不到某个用户不能处理的原因，该单将走传统审批。
SELECT
    APPROVAL.FKEYVALUE,
    ASSIGN.FASSIGNID
FROM
    T_WF_PROCINST PROCINST
    INNER JOIN T_WF_ASSIGN ASSIGN ON ASSIGN.FPROCINSTID = PROCINST.FPROCINSTID
    INNER JOIN T_WF_ACTINST ACTINST ON ASSIGN.FACTINSTID = ACTINST.FACTINSTID
    INNER JOIN T_WF_APPROVALASSIGN APPROVAL ON APPROVAL.FASSIGNID = ASSIGN.FASSIGNID
    INNER JOIN T_WF_RECEIVER RECIVER ON RECIVER.FASSIGNID = ASSIGN.FASSIGNID
    LEFT JOIN T_WF_APPROVALITEM ITEM ON ITEM.FAPPROVALASSIGNID = APPROVAL.FAPPROVALASSIGNID
    AND ITEM.FRECEIVERID = RECIVER.FRECEIVERID
WHERE
    --流程实例处于IDLE状态
    PROCINST.FSTATUS = '2'
    AND --节点处于IDLE状态
    ACTINST.FSTATUS = '2'
    AND --当前任务为未完成状态
    ASSIGN.FSTATUS = '0'
    AND --当前任务没有锁定或者锁定人是当前用户
    (
        APPROVAL.FISLOCKED = '0'
        OR (
            APPROVAL.FISLOCKED = '1'
            AND APPROVAL.FLOCKEDUSERID = @FUSERID
        )
    )
    AND (
        RECIVER.FRECEIVERID = @FUSERID
        OR EXISTS (
            SELECT
                1
            FROM
                T_WF_ENTRUSTRECEIVER C_ENTRUSTRECEIVER
            WHERE
                RECIVER.FID = C_ENTRUSTRECEIVER.FID
                AND FCONSIGNOR = @FUSERID
        )
    )
    AND (
        ITEM.FSTATUS IS NULL
        OR ITEM.FSTATUS = '0'
    )
    AND APPROVAL.FOBJECTTYPEID = @FOBJECTTYPEID
    AND APPROVAL.FKEYVALUE = @FKEYVALUE;

/* 查询当前单据的待办任务的当前状态和处理权限。
 可能的情况有：
 （1）流程挂起了。
 （2）流程正在分配处理人。
 （3）流程有其它用户正在处理中。
 （4）当前用户不是当前单据在流程的当前节点的处理人，无权审批。
 （5）查不到待办任务。*/
-- 查询用户内码
DECLARE @FUSERID INT;

SELECT
    @FUSERID = FUSERID
FROM
    T_SEC_USER
WHERE
    FNAME = 'yh';

-- 查询单据内码
DECLARE @FKEYVALUE INT;

SELECT
    @FKEYVALUE = FID
FROM
    T_PUR_POORDER
WHERE
    FBILLNO = 'CGDD002058';

-- 业务对象标识
DECLARE @FOBJECTTYPEID VARCHAR(36);

SET
    @FOBJECTTYPEID = 'PUR_PurchaseOrder';

-- 查询指定单据的待办任务状态，以及指定用户是否是当前待办任务的处理人
-- 【注意】如果某个单据，既找不到某个用户的待办任务，又找不到某个用户不能处理的原因，该单将走传统审批。
SELECT
    FKEYVALUE,
    PROCINST.FSTATUS [流程状态(1完成2运行中3挂起4终止5暂停)],
    CASE
        WHEN PROCINST.FSTATUS = 0
        OR ASSIGN.FSTATUS = 1 THEN 1
        ELSE 0
    END AS 流程正在分配处理人,
    APPROVAL.FISLOCKED 流程已有其他用户在处理,
    RECIVER.FNOPERMISSION 您不是流程当前节点的处理人
FROM
    T_WF_PROCINST PROCINST
    INNER JOIN T_WF_PIBIMAP PIBIMAP ON PROCINST.FPROCINSTID = PIBIMAP.FPROCINSTID
    LEFT JOIN T_WF_ASSIGN ASSIGN ON PROCINST.FPROCINSTID = ASSIGN.FPROCINSTID
    AND ASSIGN.FSTATUS = 1
    AND PROCINST.FSTATUS = 2
    AND ASSIGN.FCOMPLETEDTIME IS NOT NULL
    AND ASSIGN.FCOMPLETEDTIME > PROCINST.FMODIFYTIME
    LEFT JOIN (
        SELECT
            C_ASSIGN.FPROCINSTID,
            1 FISLOCKED
        FROM
            T_WF_APPROVALASSIGN C_APPROVAL
            INNER JOIN T_WF_ASSIGN C_ASSIGN ON C_APPROVAL.FASSIGNID = C_ASSIGN.FASSIGNID
            INNER JOIN T_WF_RECEIVER C_RECIVER ON C_RECIVER.FASSIGNID = C_ASSIGN.FASSIGNID
        WHERE
            C_ASSIGN.FSTATUS = 0
            AND (
                C_RECIVER.FRECEIVERID = @FUSERID
                OR EXISTS (
                    SELECT
                        1
                    FROM
                        T_WF_ENTRUSTRECEIVER C_ENTRUSTRECEIVER
                    WHERE
                        C_RECIVER.FID = C_ENTRUSTRECEIVER.FID
                        AND FCONSIGNOR = @FUSERID
                )
            )
            AND C_APPROVAL.FISLOCKED = '1'
            AND C_APPROVAL.FLOCKEDUSERID <> @FUSERID
    ) APPROVAL ON PROCINST.FPROCINSTID = APPROVAL.FPROCINSTID
    LEFT JOIN (
        SELECT
            C_ASSIGN.FPROCINSTID,
            1 FNOPERMISSION
        FROM
            T_WF_ASSIGN C_ASSIGN
        WHERE
            C_ASSIGN.FSTATUS = 0
            AND NOT EXISTS (
                SELECT
                    1
                FROM
                    T_WF_RECEIVER C_RECIVER
                WHERE
                    C_RECIVER.FASSIGNID = C_ASSIGN.FASSIGNID
                    AND (
                        C_RECIVER.FRECEIVERID = @FUSERID
                        OR EXISTS (
                            SELECT
                                1
                            FROM
                                T_WF_ENTRUSTRECEIVER C_ENTRUSTRECEIVER
                            WHERE
                                C_RECIVER.FID = C_ENTRUSTRECEIVER.FID
                                AND FCONSIGNOR = @FUSERID
                        )
                    )
            )
    ) RECIVER ON PROCINST.FPROCINSTID = RECIVER.FPROCINSTID
WHERE
    (
        PROCINST.FSTATUS = 0
        OR PROCINST.FSTATUS = 2
        OR PROCINST.FSTATUS = 3
    )
    AND FOBJECTTYPEID = @FOBJECTTYPEID
    AND FKEYVALUE = @FKEYVALUE;

/* <3>查询当前单据是否存在运行中的流程实例。
 只要单据有运行中的流程实例，就必须走工作流审批。*/
-- 查询单据内码
DECLARE @FKEYVALUE INT;

SELECT
    @FKEYVALUE = FID
FROM
    T_PUR_POORDER
WHERE
    FBILLNO = 'CGDD002058';

-- 业务对象标识
DECLARE @FOBJECTTYPEID VARCHAR(36);

SET
    @FOBJECTTYPEID = 'PUR_PurchaseOrder';

-- 查询指定单据是否存在运行中的流程实例
SELECT
    DISTINCT A.FKEYVALUE
FROM
    T_WF_PIBIMAP A
    LEFT JOIN T_WF_PROCINST B ON A.FPROCINSTID = B.FPROCINSTID
WHERE
    A.FKEYVALUE = @FKEYVALUE
    AND A.FOBJECTTYPEID = @FOBJECTTYPEID
    AND (ISNULL(B.FSTATUS, '0') NOT IN ('1', '4'));
```

# 6. 业务流程类

## 6.1 业务流程清理脏数据

```sql
--执行计划归档报插入重复性错误，直接删除历史表数据即可
-- 删除存在当前表中的历史表反写值
delete from T_BF_INSTANCEAMOUNTHIS where FRouteId in(select FRouteId from T_BF_INSTANCEAMOUNT )
--删除存在当前表的节点值
delete from T_BF_INSTANCEEntryHis where FRouteId in(select FRouteId from t_BF_InstanceEntry )
--删除存在当前表的实例值
delete from T_BF_INSTANCEHis where FInstanceId in(select FInstanceId from t_BF_Instance )

--在操作时，如果堆栈中有报对象引用错误，即get_Instance()这样的情况，则存在没有实例的节点，需要把这样的节点删除或者创建实例
--清理没有实例的节点
delete from t_bf_instanceentry where finstanceid not in (select finstanceid from t _bf_instance)

-- 清理表单被删除的反写规则
--清除多语言
delete  from  T_BF_WRITEBACKRULE_L where FID in(
select FID  from T_BF_WRITEBACKRULE where FSOURCEFORMID not in(select FID from T_META_OBJECTTYPE))
delete  from  T_BF_WRITEBACKRULE_L where FID in(
select FID  from T_BF_WRITEBACKRULE where FTARGETFORMID not in(select FID from T_META_OBJECTTYPE))
--清除反写规则
delete  from T_BF_WRITEBACKRULE where FSOURCEFORMID not in(select FID from T_META_OBJECTTYPE)
delete  from T_BF_WRITEBACKRULE where FTARGETFORMID not in(select FID from T_META_OBJECTTYPE)

-- 清理不存在实际单据数据的节点
  --首先根据表单找出表名，假设为销售订单，T_SAL_ORDERENTRY
  select t0.FTABLENUMBER from T_BF_TABLEDEFINE t0
  inner join T_META_OBJECTTYPE_L t1 on t0.FFORMID=t1.FID 
  where t1.FNAME='销售订单' and FLOCALEID=2052
  --根据表名找到不存在实际数据的节点
  select * from t_bf_instanceEntry where FTTABLENAME='T_SAL_ORDERENTRY'
  and FTID not in(select FENTRYID from T_SAL_ORDERENTRY)
  --删除节点数据
  delete from t_bf_instanceEntry where FTTABLENAME='T_SAL_ORDERENTRY'
  and FTID not in(select FENTRYID from T_SAL_ORDERENTRY)

--  清理不存在实际单据的转换规则
delete from  T_meta_ConvertRule where FSOURCEFORMID not in(select FID from T_META_OBJECTTYPE) or FTARGETFORMID not in(select FID from T_META_OBJECTTYPE) 
delete from  T_meta_ConvertRule_L where fid not in(select FID from T_meta_ConvertRule)

-- 查询是否有归档条件的业务流程实例
SELECT TOP 10000 ROW_NUMBER() OVER(ORDER BY inst0.FInstanceId) AS FId, inst0.FInstanceId AS FInstanceId 
             FROM t_BF_Instance inst0 
            INNER JOIN t_BF_Instance inst1 ON (inst0.FInstanceId = inst1.FMasterId) 
          INNER JOIN t_BF_InstanceEntry instEntry ON (inst1.FInstanceId = instEntry.FInstanceId) 
           GROUP BY inst0.FInstanceId
         HAVING max(instEntry.FCreateTime) < '当前时间减去设置值'
```

# 7. 日志类

## 7.1 上机操作日志

```sql
-------------------------------------上机操作日志-------------------------------------
-- 上级操作日志
SELECT b.FNAME AS 用户名,c.FNAME AS 子系统,d.FNAME AS 业务对象名称,a.FOPERATENAME  AS 操作名称,a.FDESCRIPTION AS 操作描述,a.* 
FROM T_BAS_OPERATELOG a 
LEFT JOIN T_SEC_USER b ON a.FUSERID=b.FUSERID 
LEFT JOIN T_META_SUBSYSTEM_L c ON a.FSUBSYSTEMID=c.FID AND c.FLOCALEID=2052 
LEFT JOIN T_META_OBJECTTYPE_L d ON a.FOBJECTTYPEID=d.FID AND d.FLOCALEID=2052 
WHERE a.FUSERID=100008
-- 上级操作日志归档表
SELECT b.FNAME AS 用户名,c.FNAME AS 子系统,d.FNAME AS 业务对象名称,a.FOPERATENAME AS 操作名称,a.FDESCRIPTION AS 操作描述,a.* 
FROM T_BAS_OPERATELOGBK a 
LEFT JOIN T_SEC_USER b ON a.FUSERID=b.FUSERID 
LEFT JOIN T_META_SUBSYSTEM_L c ON a.FSUBSYSTEMID=c.FID AND c.FLOCALEID=2052 
LEFT JOIN T_META_OBJECTTYPE_L d ON a.FOBJECTTYPEID=d.FID AND d.FLOCALEID=2052 
WHERE a.FUSERID=100008
```

## 7.2 元数据修改日志

```sql
-------------------------------------元数据修改日志-------------------------------------
SELECT TOP 10 c.FNAME AS 业务对象名称,d.FNAME AS 修改人,a.*,b.FKERNELXMLLANG 
FROM T_BAS_METADATALOG a 
LEFT JOIN T_BAS_METADATALOG_L b ON a.FID=b.FID AND b.FLOCALEID=2052 
LEFT JOIN T_META_OBJECTTYPE_L c ON a.FOBJECTID=c.FID AND c.FLOCALEID=2052 
LEFT JOIN T_SEC_USER d ON a.FUSERID=d.FUSERID 
ORDER BY a.FCREATETIME DESC
```

# 8. 后台任务

## 8.1 调度计划

### 8.1.1 执行计划

```sql
-- 执行计划
SELECT b.FNAME AS 执行计划名称,a.FSCHEDULETYPEID AS 执行计划内码,a.FSCHEDULECLASS AS 执行计划服务插件 
,CASE c.FSTATUS WHEN 0 THEN N'准备' WHEN 1 THEN N'停止' WHEN 2 THEN N'运行' WHEN 5 THEN N'异常' ELSE '' END AS 任务状态,* 
FROM T_BAS_SCHEDULETYPE a 
LEFT JOIN T_BAS_SCHEDULETYPE_L b ON a.FSCHEDULETYPEID=b.FSCHEDULETYPEID AND b.FLOCALEID=2052 
LEFT JOIN T_BAS_SCHEDULEINFO c ON a.FSCHEDULETYPEID=c.FSCHEDULETYPEID

-- 执行计划日志
-- 查询日志
SELECT * FROM T_BAS_SCHEDULEINFOLOG WHERE FSCHEDULETYPEID IN 
( 
SELECT FSCHEDULETYPEID FROM T_BAS_SCHEDULETYPE_L WHERE FNAME LIKE '%智能补货调度服务%'
)

-- 删除日志
DELETE FROM T_BAS_SCHEDULEINFOLOG WHERE FSCHEDULETYPEID IN 
( 
SELECT FSCHEDULETYPEID FROM T_BAS_SCHEDULETYPE_L WHERE FNAME LIKE '%智能补货调度服务%'
)
```


# 9. PLM

## 9.1 创建BOM展开视图

```sql
Create VIEW  V_K3_BOM
AS 
select distinct
BOM.FID BOMID, BOM.FMATERIALID PFMATERIALID,MATERIAL_L.FNAME PFMATERIALFNAME,MATERIAL.fnumber PFMATERIALFNUMBER,MATERIALBase.FERPCLSID PFERPCLSID,
BOMCHILD.FMATERIALID subFMATERIALID,subMATERIAL.FNUMBER subMATERIALFNUMBER,MATERIAL_L2.FNAME subMATERIALFNAME,subMATERIALBase.FERPCLSID subFERPCLSID,
BOM.FNUMBER BOMFNUMBER,BOM.FUSEORGID FUSEORGID,ORGANIZATIONS.FNUMBER ORGFNUMBER,ORGANIZATIONS_L.FNAME ORGFNAME,BOM.FDOCUMENTSTATUS FDOCUMENTSTATUS
,BOMCHILD.FUQTY,BOM.F_KD_MAINPRODUCTNO,BOM.F_KD_INTEGER
from T_ENG_BOM BOM
inner join T_ENG_BOMCHILD BOMCHILD  
on BOMCHILD.FID=BOM.FID
inner join T_BD_MATERIAL  MATERIAL 
on MATERIAL.FMATERIALID=BOM.FMATERIALID
inner join t_BD_MaterialBase MATERIALBase  
on MATERIALBase.FMATERIALID=BOM.FMATERIALID
inner join t_BD_MaterialBase subMATERIALBase  
on subMATERIALBase.FMATERIALID=BOMCHILD.FMATERIALID
inner join T_BD_MATERIAL  subMATERIAL
on subMATERIAL.FMATERIALID=BOMCHILD.FMATERIALID
inner join T_BD_MATERIAL_L MATERIAL_L
on (MATERIAL_L.FMATERIALID=BOM.FMATERIALID and MATERIAL_L.FLOCALEID='2052')
inner join T_BD_MATERIAL_L MATERIAL_L2
on (MATERIAL_L2.FMATERIALID=BOMCHILD.FMATERIALID and MATERIAL_L2.FLOCALEID='2052')
inner join T_ORG_ORGANIZATIONS_L ORGANIZATIONS_L
on (ORGANIZATIONS_L.FORGID=BOM.FUSEORGID and ORGANIZATIONS_L.FLOCALEID='2052')
inner join T_ORG_ORGANIZATIONS ORGANIZATIONS 
on ORGANIZATIONS.FORGID=BOM.FUSEORGID;


-- 存储过程传入字符串后需要拆分 例如‘a,b,c’
-- 1)拆分： 函数一
create  function  SplitIn(@c   varchar (2000),@split   varchar (2))  
returns   @t   table (col   varchar (20))  
as  
begin   
   while(charindex(@split,@c)<>0)  
     begin  
       insert   @t(col)   values   ( substring (@c,1,charindex(@split,@c)-1))  
       set   @c   =   stuff(@c,1,charindex(@split,@c), '' )  
     end  
   insert   @t(col)   values   (@c)  
   return  
end
-- 2)拆分：函数二：
create function [dbo].[f_split](@SourceSql varchar(max),@StrSeprate varchar(10))
    returns @temp table(Rowvalue varchar(1000))
as
    begin
        declare @i int
        set @SourceSql=rtrim(ltrim(@SourceSql))
        set @i=charindex(@StrSeprate,@SourceSql)
        while @i>=1
        begin
            insert @temp values(left(@SourceSql,@i-1))
            set @SourceSql=substring(@SourceSql,@i+1,len(@SourceSql)-@i)
            set @i=charindex(@StrSeprate,@SourceSql)
        end
        if @SourceSql<>''
        insert @temp values(@SourceSql)
     return
    end
```

# 10 问题清单

> 单据体先有数据，后启用序号列，新的数据行会自动刷新序号，而历史数据是不会自动刷新序号的，需要手动刷新。
>
> 下面以流程配置中心为例，演示如何使用SQL脚本刷新序号列的数据。

```sql
-- 更新单行分录的单据体的序号
UPDATE
    T_WF_PRCLAUNCHCONFIG
SET
    FSEQ = 1
WHERE
    FTMPID IN (
        SELECT
            FTMPID
        FROM
            T_WF_PRCLAUNCHCONFIG
        GROUP BY
            FTMPID
        HAVING
            (COUNT(*) = 1)
    );
  
-- 更新多行分录的单据体的序号（拷贝此查询语句生成的SQL并执行，即可实现单据体序号的刷新）
-- SET     a.FSEQ = b.ROWID -- 完全覆盖
-- SET     a.FSEQ = CASE WHEN a.FSEQ=0 THEN b.ROWID ELSE a.FSEQ END -- 仅更新值等于0的序号
SELECT
    'UPDATE  a
SET     a.FSEQ = b.ROWID
FROM    T_WF_PRCLAUNCHCONFIG a
        JOIN ( SELECT   ROW_NUMBER() OVER ( ORDER BY FTMPDETAILID ) AS ROWID ,
                        *
               FROM     T_WF_PRCLAUNCHCONFIG
               WHERE    FTMPID = ' + CAST(FTMPID AS VARCHAR(50)) + '
             ) b ON a.FTMPDETAILID = b.FTMPDETAILID'
FROM
    T_WF_PRCLAUNCHCONFIG
GROUP BY
    FTMPID
HAVING
    (COUNT(*) > 1);
  
/*
    -- 更新序号示例SQL：
    UPDATE  a
    SET     a.FSEQ = b.ROWID
    FROM    T_WF_PRCLAUNCHCONFIG a
    JOIN ( SELECT   ROW_NUMBER() OVER ( ORDER BY FTMPDETAILID ) AS ROWID ,
    *
    FROM     T_WF_PRCLAUNCHCONFIG
    WHERE    FTMPID = 100002
    ) b ON a.FTMPDETAILID = b.FTMPDETAILID
    */
```

> 在数据库直接改了枚举类型定义后，使用了该枚举项的单据上的下拉列表的下拉项并不会马上生效？

```sql
/*查询采购订单的元数据版本号*/
SELECT FMAINVERSION,* FROM T_META_OBJECTTYPE WHERE FID='PUR_PurchaseOrder' OR FINHERITPATH LIKE '%,PUR_PurchaseOrder,%'
/*修改采购订单的元数据版本号，版本号变化后，元数据缓存会自动缓存*/
UPDATE T_META_OBJECTTYPE SET FMAINVERSION=CAST(CAST(FMAINVERSION AS BIGINT)+1 AS varchar(100)) WHERE FID='PUR_PurchaseOrder' OR FINHERITPATH LIKE '%,PUR_PurchaseOrder,%'
```

> 在业务并发量大的情况下,数据库CPU达到100%,如何排查可能导致数据库服务器CPU异常的SQL?

```sql
SELECT
    s.session_id,
    r.start_time,
    r.total_elapsed_time / 1000 AS elapsed_seconds,
    r.status AS RequestStatus,
    r.database_id,
    DB_NAME(r.database_id) AS DatabaseName,
    r.reads,
    r.writes,
    r.logical_reads,
    r.cpu_time,
    r.total_elapsed_time,
    s.status AS SessionStatus,
    qt.event_info as text,
    getdate() createtime
FROM
    sys.dm_exec_requests(nolock) r
    INNER JOIN sys.dm_exec_sessions(nolock) s ON r.session_id = s.session_id
    LEFT JOIN sys.dm_exec_connections(nolock) c ON c.session_id = s.session_id
    cross apply sys.dm_exec_input_buffer (s.session_id, 0) as qt
WHERE
    DB_NAME(r.database_id) <> 'master'
    and qt.event_info is not null
order by
    r.cpu_time desc
```

> 成本计算单打开慢，优化成本表相关索引？

```sql
alter index all on T_CB_PROORDERINFO_H rebuild;

UPDATE
    STATISTICS T_CB_PROORDERINFO_H;

ALTER TABLE
    T_CB_PROORDERINFO_H rebuild WITH (DATA_COMPRESSION = ROW);

alter index all on T_CB_PROORDERDIME rebuild;

UPDATE
    STATISTICS T_CB_PROORDERDIME;

ALTER TABLE
    T_CB_PROORDERDIME rebuild WITH (DATA_COMPRESSION = ROW);

alter index all on T_CB_PROORDERTYPE rebuild;

UPDATE
    STATISTICS T_CB_PROORDERTYPE;

ALTER TABLE
    T_CB_PROORDERTYPE rebuild WITH (DATA_COMPRESSION = ROW);

alter index all on T_CB_CostCalExpense_H rebuild;

UPDATE
    STATISTICS T_CB_CostCalExpense_H;

ALTER TABLE
    T_CB_CostCalExpense_H rebuild WITH (DATA_COMPRESSION = ROW);

alter index all on T_CB_COSTCENTER rebuild;

UPDATE
    STATISTICS T_CB_COSTCENTER;

ALTER TABLE
    T_CB_COSTCENTER rebuild WITH (DATA_COMPRESSION = ROW);


-- 成本计算单(横向）打开慢
create index IDX_CB_PROORDERTYPE_INDEX0523 on T_CB_PROORDERTYPE(
    FPROORDERTYPE,
    FPROORDERENTRYID,
    FNUMBER,
    FPRODUCTID,
    FCOSTCENTERID
);
```

> 工作流，参与人变量设置为采购员，解析异常时，可通过SQL分析异常原因?

```sql
-- <1>查询采购员内码
SELECT
    FPURCHASERID
FROM
    T_PUR_POORDER
WHERE
    FBILLNO = 'CGDD001938';

-- <2>查询参与人，（用采购员内码查询，模拟工作流解析：参与人变量.采购员）
SELECT
    U.FUSERID,
    U.FNAME,
    STAFF.FPOSTID
FROM
    T_BD_OPERATORENTRY OPERA
    INNER JOIN T_BD_STAFF STAFF ON OPERA.FSTAFFID = STAFF.FSTAFFID
    AND STAFF.FDOCUMENTSTATUS = 'C'
    AND STAFF.FFORBIDSTATUS = 'A'
    INNER JOIN V_BD_CONTACTOBJECT CONTACT ON STAFF.FPERSONID = CONTACT.fid
    AND CONTACT.FDOCUMENTSTATUS = 'C'
    AND CONTACT.FFORBIDSTATUS = 'A'
    INNER JOIN T_SEC_USER U ON U.FLINKOBJECT = CONTACT.fid
    AND U.FFORBIDSTATUS = 'A'
WHERE
    OPERA.FISUSE = '1'
    AND OPERA.FENTRYID = 300752;
```
