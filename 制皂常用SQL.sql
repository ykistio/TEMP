use AIS20241024092116;

-- 查询类名是否已经注册到那个业务对象上
select a.FID,a.FKERNELXML,b.FNAME,a.Dome  
from  
(select FKERNELXML.query('//ClassName') 'Dome', * from T_META_OBJECTTYPE)  a   
left join T_META_OBJECTTYPE_L b on b.FID = a.FID
 --插件类名
where convert(varchar(max),Dome) like '%QuotationClose%';

-- 查询单位
SELECT * FROM T_BD_UNIT;
-- 查询单位不一致的数据
SELECT * FROM T_BD_MATERIAL WHERTE  FBASEUNITID <> FSALEURNUM;
-- 查询销售信息
select FSALEURNUM from t_BD_MaterialSale;
-- 查询组织信息
select * from T_ORG_Organizations where  fnumber = '500';
-- 查询物料基础信息
select FSAPNUMBER,FNUMBER from T_BD_MATERIAL where FSAPNUMBER in('80666867', '80666869', '80666870', '80667641') and FUSEORGID = '406423';
-- 查询SAP编码分组，sap编码必须唯一
select FSAPNUMBER,count(*) from T_BD_MATERIAL where FSAPNUMBER != '' and FUSEORGID = '406423' group by FSAPNUMBER  ;
-- 根据物料编码查询物料
select * from T_BD_MATERIAL where fnumber = '80092651';

-- 查询月度要货计划
select * from BBG_t_Cust100010 where FREQUESTORGID = '406423' order by FCREATEDATE;

-- 查询信用
select * from BBG_t_CustomerCredit;

-- 查询物料编码
select * from t_bd_material where fnumber = '21011062L' and FUSEORGID = '406423' 

-- 更新物料，SAP编码，方便可以添加月度要货计划
select * from T_BD_MATERIAL;
update T_BD_MATERIAL set FSAPNUMBER = '8383iekydk' WHERE FMATERIALID = '495806';

-- 查询要货计划汇总单
select tma.fmaterialid, tma.FNUMBER, t2.F_BBG_Qty, t2.F_BBG_SALEQTY, FENTRYID
from BBG_t_Cust100010 t1
inner join BBG_t_Cust_Entry100026 t2
inner join T_BD_MATERIAL tma on t2.FMATERIALID = tma.FMATERIALID
on t1.FID = t2.FID
where FREQUESTORGID = 406423;

-- 更新分录中的数据
update BBG_t_Cust_Entry100026 
SET F_BBG_Qty = 2, 
F_BBG_SALEQTY = 2
where FENTRYID = 100222;

update BBG_t_Cust_Entry100026 
SET F_BBG_Qty = 2, 
F_BBG_SALEQTY = 2
where FENTRYID = 100223;

update BBG_t_Cust_Entry100026 
SET F_BBG_Qty = 8, 
F_BBG_SALEQTY = 8
where FENTRYID = 100217;
-- where FENTRYID in (100217, 100219, 100220, 100221, 100222, 100223, 100224);

-- 销售报价单校验， 暂存，更新SAP订单号，提交
select * from T_SAL_QUOTATION where fdocumentstatus = 'Z' order by FDate desc;
update T_SAL_QUOTATION set F_BBG_CONTRACTNO = '12', F_ora_Mailtxt = '2059416370@qq.com'  where fid = 153281;

-- 修改客户信用
select * from BBG_t_CustomerCredit;
update BBG_t_CustomerCredit set F_BBG_SALECHANNEL = 18000000 where fnumber = '30080003';

-- 销售报价单关闭
update T_SAL_QUOTATION set FBillStatus = 'A' where FBillNo in ('BJ550118002503001', 'BJ310728712503002');

-- 销售订单关闭
select * from T_SAL_ORDER where FID in ('195524', '195598');
select * from T_SAL_ORDERENTRY where FID = 195598;

-- 更新SAP发货数量
update T_SAL_ORDERENTRY set  F_BBG_SAPSTOCKOUTQTY = 28
 where FID = 195524;

 -- 移除其他不在库存单据中的记录，方便测试
delete from T_SAL_ORDERENTRY where fentryid in (373240, 373241, 373242, 373243, 373244);
