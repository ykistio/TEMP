### 1. api

#### api.controller

```java
@ApiController(value = "应用", desc = "应用")
@ApiMapping("/api")
public class TestWebApiController {

    @ApiGetMapping("/testGet")
    public CustomApiResult<String> testGet(@ApiParam(value = "名称", required = true, message = "针对不能为空") String name) {
        if (StringUtils.equals(name, "kate")) {
            TestLogUtil.writeServeLog("testGet", "B", "传入的名称为：" + name, "", "名称不能为kate");
            return CustomApiResult.fail("500", "名称不能为kate");
        }
        String resp = "123" + name;
        TestLogUtil.writeServeLog("testGet", "A", "传入的名称为：" + name, resp, "");
        return CustomApiResult.success(resp);
    }


    @ApiPostMapping("/testPost")
    public CustomApiResult<@ApiResponseBody PersonModel> testPost(@ApiRequestBody(value = "传入参数", required = true, message = "测试确实不能为空") ParamModel paramModel) {
        String name = paramModel.getName();
        Integer age = paramModel.getAge();
        if (StringUtils.equals(name, "kate")) {
            TestLogUtil.writeServeLog("testPost", "B", "传入的JSON为：" + JSONObject.toJSONString(paramModel), "", "名称不能为kate");
            return CustomApiResult.fail("500", "名称不能为kate");
        }
        if (age == null) {
            TestLogUtil.writeServeLog("testPost", "B", "传入的JSON为：" + JSONObject.toJSONString(paramModel), "", "年龄不能为空");
            return CustomApiResult.fail("500", "年龄不能为空");
        }
        PersonModel personModel = new PersonModel();
        personModel.setName(name);
        personModel.setAge(age);
        personModel.setId(5);
        TestLogUtil.writeServeLog("testPost", "A", "请求报文为：" + JSONObject.toJSONString(paramModel), "响应报文为" + JSONObject.toJSONString(personModel), "");
        return CustomApiResult.success(personModel);
    }
}
```

#### api.log

```java
public class TestLogUtil {

    /**
     *  相应日志
     */
    public static final String SERVE_LOG = "shkd_resplog";
    public static final String CLIENT_LOG = "shkd_clientlog";

    /**
     * 写服务日志
     * @param opname 操作名称
     * @param isSuccess 是否成功 A 成功， B 失败
     * @param reqLog 请求日志
     * @param respLog 响应日志
     * @param errorLog 错误日志
     */
    public static void writeServeLog(String opname, String isSuccess, String reqLog, String respLog, String errorLog) {
        DynamicObject log = BusinessDataServiceHelper.newDynamicObject(SERVE_LOG);
        DynamicObject user = UserServiceHelper.getCurrentUser("name");
        log.set("username", user.getString("name"));
        log.set("opname", opname);
        log.set("opdesc", opname);
        log.set("opdate", DateUtil.getCurrentDateTime());
        log.set("shkd_reqlog_tag", reqLog);
        log.set("shkd_responselog_tag", respLog);
        errorLog = (errorLog.length() > 2000) ? errorLog.substring(0, 2000) : errorLog;
        log.set("shkd_errorlog", errorLog);
        log.set("shkd_issuccess", isSuccess);
        LogORM logOrm = LogORM.create();
        long[] insert = logOrm.insert(Collections.singletonList(log));
        if (insert.length == 0) {
            throw new RuntimeException("日志写入失败");
        }
    }

    /**
     * 写服务日志
     * @param opname 操作名称
     * @param isSuccess 是否成功 A 成功， B 失败
     * @param reqLog 请求日志
     * @param respLog 响应日志
     * @param errorLog 错误日志
     */
    public static void writeClientLog(String opname, String isSuccess, String reqLog, String respLog, String errorLog) {
        DynamicObject log = BusinessDataServiceHelper.newDynamicObject(CLIENT_LOG);
        DynamicObject user = UserServiceHelper.getCurrentUser("name");
        log.set("username", user.getString("name"));
        log.set("opname", opname);
        log.set("opdesc", opname);
        log.set("opdate", DateUtil.getCurrentDateTime());
        log.set("shkd_issuccess", isSuccess);
        log.set("shkd_requestlog_tag", reqLog);
        log.set("shkd_responselog_tag", respLog);
        errorLog = (errorLog.length() > 2000) ? errorLog.substring(0, 2000) : errorLog;
        log.set("shkd_errorlog", errorLog);
        LogORM logOrm = LogORM.create();
        long[] insert = logOrm.insert(Collections.singletonList(log));
        if (insert.length == 0) {
            throw new RuntimeException("日志写入失败");
        }
    }

}
```

#### api.reqmodel

```java
import kd.bos.openapi.common.custom.annotation.ApiModel;
import kd.bos.openapi.common.custom.annotation.ApiParam;

import javax.validation.constraints.NotNull;
import java.io.Serializable;

/**
 * @Description: 参数模型
 * @Author: hyk
 * @Date: 2024/11/30
 **/
@ApiModel
public class ParamModel implements Serializable {

    @ApiParam(value = "名称", example = "张三")
    private String name;

    @ApiParam(value = "年龄", example = "20")
    @NotNull(message = "年龄不能为空")
    private Integer age;
}
```

#### api.respmodel

```java
import kd.bos.openapi.common.custom.annotation.ApiModel;

import java.io.Serializable;

/**
 * @Description: TODO
 * @Author: hyk
 * @Date: 2024/11/30
 **/
@ApiModel
public class PersonModel implements Serializable {

    private Integer id;

    private String name;

    private Integer age;
}
```

### 2. ctx

```java
RequestContext requestContext = RequestContext.get();
// 租户ID
String tenantId = requestContext.getTenantId();
// 租户别名
String tenantCode = requestContext.getTenantCode();
// 账套ID
String accountId = requestContext.getAccountId();
// 用户ID
String uid = requestContext.getUid();
// 用户名称
String userName = requestContext.getUserName();
// 用户类型
String userType = requestContext.getUserType();
// 获取用户登录的组织内码
Long orgId = requestContext.getOrgId();
// 获取第三方单点切换组织代码
String loginOrg = requestContext.getLoginOrg();
// 云之家用户在IERP工作圈内ID
String userOpenId = requestContext.getUserOpenId();
// 获取语言
Lang lang = requestContext.getLang();
// 获取登录IP地址
String loginIP = requestContext.getLoginIP();
// 获取登录时间
String loginTime = requestContext.getLoginTime();
// 获取云之家应用登录的AppId
String yzjAppId = requestContext.getYzjAppId();
// 获取云之家应用登录的Ticket
String yzjAppTicket = requestContext.getYzjAppTicket();
// 获取用户登录客户端
String client = requestContext.getClient();
```

### 3. db

#### db.exec

```java
try {
    String sql = "UPDATE t_xxx SET field=? WHERE fid=?;";
    boolean execute = DB.execute(DBRoute.basedata, sql, new Object[]{"xxx", "xxx"});
    if (!execute) {
        throw new KDBizException("写数据库失败");
    }
} catch (Exception e) {
    throw new KDBizException(e.getMessage());
}
```

#### db.execbatch

```java
try {
    String sql = "UPDATE t_xxx SET field=? WHERE fid=?;";
    List<Object[]> paramsList = new ArrayList();
    Object[] params = new Object[]{"xxx", "xxx"};
    paramsList.add(params);
    int[] executeBatch = DB.executeBatch(DBRoute.basedata, sql, paramsList);
    if (executeBatch.length == 0) {
        throw new KDBizException("更新失败");
    }
} catch (Exception e) {
    throw new KDBizException(e.getMessage());
}
```

#### db.query

```java
String sql = "select fid, ftruename from t_SEC_User where fid = ?";
List<String> fids = new ArrayList<>();
List<String> query = DB.query(DBRoute.of("sys"), sql, new Object[] {2054602484386629632L},  new ResultSetHandler<List<String>>() {
    @Override
    public List<String> handle(ResultSet rs) {
        try {
            while (rs.next()) {
                fids.add(rs.getString("ftruename"));
            }
        } catch (Exception e) {
        }
        return fids;
    }
});
this.getView().showTipNotification("结果:" + query);
```

#### db.update

```java
try {
    // 修改数据示例
    String sql = "update t_SEC_User set ftruename = ? from  where fid = ?";
    int sys = DB.update(DBRoute.of("sys"), sql, new Object[]{"AA", 2012695299436315648L});
    if (sys > 0) {
        this.getView().showTipNotification("结果: 更新成功");
    }
} catch (Exception e) {
    throw new KDBizException(e.getMessage());
}
```

### 4. form

#### form.abd

```java
/**
 * 界面数据包构建完毕，生成指令，刷新前端字段值，控件状态之后
 */
@Override
public void afterBindData(EventObject e) {
    super.afterBindData(e);
}
```

#### form.acnd

```java
/**
 * 界面初始化或者刷新，创建数据包成功，字段填写默认值之后
 */
@Override
public void afterCreateNewData(EventObject e) {
    super.afterCreateNewData(e);
}
```

#### form.ado

```java
/**
 * 用户点击按钮、菜单，执行完绑定的操作后，不论成功与否，均会触发此事件
 */
@Override
public void afterDoOperation(AfterDoOperationEventArgs afterDoOperationEventArgs) {
    super.afterDoOperation(afterDoOperationEventArgs);
    FormOperate operate = (FormOperate) afterDoOperationEventArgs.getSource();
    String operateKey = operate.getOperateKey();
    if (StringUtils.equals(operateKey, "")) {

    }
}
```

#### form.adr

```java
/**
 * 插件可以在此事件，进行删除行后的数据同步处理，如刷新前端表格控件。
 */
@Override
public void afterDeleteRow(AfterDeleteRowEventArgs e) {
    EntryProp entryProp = e.getEntryProp();
    String entryPropName = entryProp.getName();
    if (StringUtils.equals(entryPropName, "")) {
      
    }
}
```

#### form.ald

```java
/**
 * 单据界面加载完毕，会根据传入的单据主键值，到数据库加载单据数据包。
 */
@Override
public void afterLoadData(EventObject e) {
    super.afterCreateNewData(e);
}
```

#### form.bbd

```java
/**
 * 界面数据包构建完毕，开始生成指令，刷新前端字段值，控件状态之前。
 */
@Override
public void beforeBindData(EventObject e) {
    super.beforeBindData(e);
}
```

#### form.bc

```java
/**
 * 界面关闭之前触发此事件
 */
@Override
public void beforeClosed(BeforeClosedEvent e) {
    super.beforeClosed(e);
}
```

#### form.bclick

```java
/**
 * 用户点击按钮, 标签，可以注册点击的事件
 */
@Override
public void beforeClick(BeforeClickEvent evt) {
    super.beforeClick(evt);
    Control source = (Control) evt.getSource();
    String key = source.getKey();
}
```

#### form.bdo

```java
/**
 * 用户点击按钮、菜单，执行绑定的操作逻辑前，触发此事件
 */
@Override
public void beforeDoOperation(BeforeDoOperationEventArgs args) {
    super.beforeDoOperation(args);
    FormOperate operate = (FormOperate) args.getSource();
    String operateKey = operate.getOperateKey();
    if (StringUtils.equals(operateKey, "")) {

    }
}
```

#### form.bic

```java
/**
 * 用户点击菜单按钮，在执行按钮绑定的操作前，触发此事件
 */
@Override
public void beforeItemClick(BeforeItemClickEvent evt) {
    String itemKey = evt.getItemKey();
    String operationKey = evt.getOperationKey();
    if (StringUtils.equals(operationKey, "")) {
    }
}
```

#### form.bid

```java
/**
 * 引入数据填写之前触发的事件，批量处理时，每张单据触发一次。可以对本次引入的源单数据包，进行修订，或者取消本单引入。
 * 例如：判断单据状态不为暂存取消本单引入:
 * e.setCancel(true);
 * e.setCancelMessage();
 * @param e 引入时上下问参数
 */
@Override
public void beforeImportData(BeforeImportDataEventArgs e) {
    JSONObject bizJsonObject = (JSONObject) e.getSourceData();

}
```

#### form.c

```java
/**
 * 用户点击按钮, 标签，可以注册点击的事件
 */
@Override
public void click(EventObject evt) {
    super.click(evt);
    Control source = (Control) evt.getSource();
    String key = source.getKey();
    if (StringUtils.equals(key, "")) {

    }
}
```

#### form.cc

```java
/**
 * 用户确认了交互提示信息后，触发此事件，通知插件进行后续处理
 */
@Override
public void confirmCallBack(MessageBoxClosedEvent evt) {
    // 回调ID
    String callBackId = evt.getCallBackId();
    // 消息选择结果
    MessageBoxResult messageboxResult = evt.getResult();
    int resultInt = messageboxResult.getValue();
    String resultVal = evt.getResultValue();
}
```

#### form.ccb

```java
/**
 * 子界面关闭时，触发父界面的closedCallBack事件
 */
@Override
public void closedCallBack(ClosedCallBackEvent closedCallBackEvent) {
    super.closedCallBack(closedCallBackEvent);
}
```

#### form.ce

```java
/**
 * 用于触发自定义控件的定制事件
 */
@Override
public void customEvent(CustomEventArgs e) {
    super.customEvent(e);
}
```

#### form.cnd

```java
/**
 * 界面初始化或者刷新，创建数据包时
 */
@Override
public void createNewData(BizDataEventArgs e) {
    super.createNewData(e);
}
```

#### form.erc

```java
/**
 * 分录行点击事件
 * @param evt 分录行点击事件
 */
@Override
public void entryRowClick(RowClickEvent evt) {
    IPageCache pageCache = this.getPageCache();
    AbstractFormDataModel model = (AbstractFormDataModel)this.getModel();
    EntryGrid grid = (EntryGrid)evt.getSource();
    String entryKey = grid.getKey();
    if (StringUtils.equals(entryKey, "")) {

    }
}
```

#### form.fbc

```java
/**
 * 弹性域维护界面关闭时，触发父界面此事件
 */
@Override
public void flexBeforeClosed(FlexBeforeClosedEvent e) {
    super.flexBeforeClosed(e);
}
```

#### form.get

```java
/**
 * 创建数据包之前
 */
@Override
public void getEntityType(GetEntityTypeEventArgs e) {
    super.getEntityType(e);
}
```

#### form.i

```java
/**
 * 视图初始化，创建插件后
 */
@Override
public void initialize() {
    super.initialize();
}
```

#### form.ic

```java
/**
 * 用户点击菜单项，触发此事件
 */
@Override
public void itemClick(ItemClickEvent evt) {
    String itemKey = evt.getItemKey();
    String operationKey = evt.getOperationKey();
    if (StringUtils.equals(operationKey, "")) {
    }
}
```

#### form.lccm

```java
/**
 * 界面显示前，构建界面显示参数时
 */
@Override
public void loadCustomControlMetas(LoadCustomControlMetasArgs e) {
    super.loadCustomControlMetas(e);
}
```

#### form.mcb

```java
/**
 * 确认回调接口触发
 */
@Override
public void confirmCallBack(MessageBoxClosedEvent messageBoxClosedEvent) {
    String callBackId = messageBoxClosedEvent.getCallBackId();
    String resultValue = messageBoxClosedEvent.getResultValue();
    if (StringUtils.equals(callBackId, "") && StringUtils.equals(resultValue, "OK")) {
      
    }
}
```

#### form.ogc

```java
/**
 * 在有代码尝试获取控件的编程模型时，触发此事件
 */
@Override
public void onGetControl(OnGetControlArgs e) {
    super.onGetControl(e);
}
```

#### form.pc

```java
/**
 *  界面数据初始时，不触发此事件，即在afterCreateNewData事件中，修改字段值，不会触发此事件；
 *  此事件在实体服务规则执行之后发生：用户手工录入字段值时，先执行与此字段相关的实体服务规则，然后再触发插件propertyChanged事件。
 */
@Override
public void propertyChanged(PropertyChangedArgs e) {
    IDataEntityProperty property = e.getProperty();
    String name = property.getName();
    if (StringUtils.equals(name, "")) {
        ChangeData[] changeSet = e.getChangeSet();
        int rowIndex = changeSet[0].getRowIndex();
        Object newValue = changeSet[0].getNewValue();
        if (newValue != null) {

        }
    }
}
```

#### form.pof

```java
/**
 * 接收到showForm指令后触发, 可以修改界面显示参数或取消显示
 */
@Override
public void preOpenForm(PreOpenFormEventArgs e) {
    super.preOpenForm(e);
}
```

#### form.pr

```java
/**
 * 界面关闭后，系统开始释放表单所占的资源时，触发此事件
 */
@Override
public void pageRelease(EventObject e) {
    super.pageRelease(e);
}
```

#### form.rl

```java
/**
 * 用户交互时触发
 */
@Override
public void registerListener(EventObject e) {
    super.registerListener(e);
}
```

#### form.sv

```java
/**
 * 视图初始化，创建插件时
 */
@Override
public void setView(IFormView formView) {
    super.setView(formView);
}
```

#### form.te

```java
/**
 * 前端定时发送请求，触发此界面事件
 */
@Override
public void TimerElapsed(TimerElapsedArgs e) {
    super.TimerElapsed(e);
}
```

### 5. list

#### list.bcc

```java
/**
 * 列表打开的单据界面关闭并返回到列表时，触发此事件
 */
@Override
public void billClosedCallBack(BillClosedCallBackEvent e) {
    CloseCallBack closeCallBack = e.getCloseCallBack();
    String actionId = closeCallBack.getActionId();
    Object pkId = e.getPkId();

}
```

#### list.bclc

```java
/**
 * 刷新单据列表，构建单据列表显示的列时，触发此事件
 */
@Override
public void beforeCreateListColumns(BeforeCreateListColumnsArgs args) {
    super.beforeCreateListColumns(args);
}
```

#### list.bcldp

```java
/**
 * 构建列表取数器之前，触发此事件
 */
@Override
public void beforeCreateListDataProvider(BeforeCreateListDataProviderArgs args) {
    super.beforeCreateListDataProvider(args);
}
```

#### list.bdcdfs

```java
/**
 * 设置常用过滤的基础资料依赖字段
 */
@Override
public void baseDataColumnDependFieldSet(BaseDataColumnDependFieldSetEvent args) {
    super.baseDataColumnDependFieldSet(args);
}
```

#### list.bic

```java
/**
 * 用户点击单据列表主菜单工具栏上的按钮时，触发此事件
 */
@Override
public void beforeItemClick(BeforeItemClickEvent evt) {
    super.beforeItemClick(evt);
}
```

#### list.blhlc

```java
/**
 * 超链接的单元格，用户点击时，在打开单据界面之前，触发此事件
 */
@Override
public void billListHyperLinkClick(HyperLinkClickArgs args) {
    super.billListHyperLinkClick(args);
}
```

#### list.bsb

```java
/**
 * 在执行操作时，打开子界面之前，触发此事件
 */
@Override
public void beforeShowBill(BeforeShowBillFormEvent e) {
    super.beforeShowBill(e);
}
```

#### list.fcasc

```java
/**
 * 过滤容器搜索点击后的处理方法,此事件发生在过滤条件解析后，主要用于点击过滤条件时联动修改其他过滤字段控件。
 */
@Override
public void filterContainerAfterSearchClick(FilterContainerSearchClickArgs args) {
    super.filterContainerAfterSearchClick(args);
}
```

#### list.fcbf7s

```java
/**
 * 过滤容器内F7弹出前的处理方法
 */
@Override
public void filterContainerBeforeF7Select(BeforeFilterF7SelectEvent args) {
    String fieldName = args.getFieldName();
    if (StringUtils.equals(fieldName, "")) {
        // 方案过滤视图, 在物料为条件配置过滤方案时, 只能选择名称包含"主机"的物料数据
        args.addCustomQFilter(new QFilter("name", QFilter.like, "%主机%"));
    }
}
```

#### list.fci

```java
/**
 * 列表界面，初始化过滤面板时，触发此事件
 * 用户在过滤面板点击搜索时，也会重新初始化过滤面板，触发此事件
 */
@Override
public void filterContainerInit(FilterContainerInitArgs args) {
    super.filterContainerInit(args);
}
```

#### list.fcsc

```java
/**
 * 用户在过滤条件面版上改变过滤条件，触发此事件
 */
@Override
public void filterContainerSearchClick(FilterContainerSearchClickArgs args) {
    super.filterContainerSearchClick(args);
}
```

#### list.fcsf

```java
/**
 * 过滤字段上的基础资料字段过滤条件调整事件
 */
@Override
public void filterColumnSetFilter(SetFilterEvent args) {
    super.filterColumnSetFilter(args);
    String fieldName = args.getFieldName();
    if (StringUtils.equals(fieldName, "materielfield.name")) {
        // 常用过滤视图下的物料字段只展示名称包含"可乐"的数据
        args.addCustomQFilter(new QFilter("name", QFilter.like, "%可乐%"));
    }
}
```

#### list.ic

```java
/**
 * 用户点击单据列表主菜单工具栏上的按钮时，触发此事件
 */
@Override
public void itemClick(ItemClickEvent evt) {
    String itemKey = evt.getItemKey();
    if (StringUtils.equals(itemKey, "")) {

    }
}
```

#### list.lrc

```java
/**
 * 用户点击单据列表行时，触发此事件。
 */
@Override
public void listRowClick(ListRowClickEvent evt) {
    super.listRowClick(evt);
}
```

#### list.lrdc

```java
/**
 * 用户双击单据列表行时，触发此事件。
 */
@Override
public void listRowDoubleClick(ListRowClickEvent evt) {
    super.listRowDoubleClick(evt);
}
```

#### list.scfv

```java
/**
 * 设置单元格指令
 */
@Override
public void setCellFieldValue(SetCellFieldValueArgs args) {
    super.setCellFieldValue(args);
}
```

#### list.sf

```java
/**
 * 单据列表控件，在构建好取数条件，准备重新取数之前，触发此事件
 */
@Override
public void setFilter(SetFilterEvent e) {
    super.setFilter(e);
}
```

### 6. op

#### op.aeot

```java
/**
 * 操作执行完毕，事务提交之后，触发此事件
 */
@Override
public void afterExecuteOperationTransaction(AfterOperationArgs e) {
    String operationKey = e.getOperationKey();
    if (StringUtils.equals(operationKey, "")) {
        DynamicObject[] dataEntityArray = e.getDataEntities();
        for (DynamicObject dataEntity : dataEntityArray) {

        }
    }
}
```

#### op.beot

```java
/**
 * 操作校验通过之后，开启事务存储数据之前，触发此事件
 */
@Override
public void beforeExecuteOperationTransaction(BeforeOperationArgs e) {
    super.beforeExecuteOperationTransaction(e);
}
```

#### op.bot

```java
/**
 * 操作校验通过，开启了事务之后，还没有把数据提交到数据库之前触发此事件
 */
@Override
public void beginOperationTransaction(BeginOperationTransactionArgs e) {
    super.beginOperationTransaction(e);
}
```

#### op.eot

```java
/**
 * 单据数据已经提交到数据库之后，事务未提交之前，触发此事件
 */
@Override
public void endOperationTransaction(EndOperationTransactionArgs e) {
    String operationKey = e.getOperationKey();
    if (StringUtils.equals(operationKey, "")) {
        DynamicObject[] dataEntityArray = e.getDataEntities();
        for (DynamicObject dataEntity : dataEntityArray) {

        }
    }
}
```

#### op.ior

```java
/**
 * 初始化操作结果集
 */
@Override
public void initializeOperationResult(OperationResult result) {
    super.initializeOperationResult(result);
}
```

#### op.oav

```java
/**
 * 构建好操作校验器之后，执行操作校验之前，触发此事件
 */
@Override
public void onAddValidators(AddValidatorsEventArgs e) {
    super.onAddValidators(e);
}
```

#### op.pp

```java
/**
 * 操作执行，加载单据数据包之前，触发此事件
 */
@Override
public void onPreparePropertys(PreparePropertysEventArgs e) {
    super.onPreparePropertys(e);
}
```

#### op.ro

```java
/**
 * 操作事务提交失败，事务回滚之后触发此事件
 */
@Override
public void rollbackOperation(RollbackOperationArgs e) {
    super.rollbackOperation(e);
}
```

#### op.sc

```java
/**
 * 设置上下文
 */
@Override
public void setContext(MainEntityType billEntityType, Map<String, Object> operateMeta, OperateOption option) {
    super.setContext(billEntityType, operateMeta, option);
}
```

#### op.v

```java
@Override
public void validate() {
    ExtendedDataEntity[] dataEntities = this.getDataEntities();
    for (ExtendedDataEntity extendedDataEntity : dataEntities) {
        DynamicObject dataEntity = extendedDataEntity.getDataEntity();

    }
}
```

### 7. tx

```java
try (TXHandle handler = TX.requiresNew("fi_voucher_submit")) {
    try {
        // 在这里执行的，有事务。
    } catch (Throwable e) {
        handler.markRollback();  
        throw e;
    }
}
```

### 8. mq

mq.consumer

```java
public class DemoConsumer implements MessageConsumer{
    public void onMessage(Object body, String messageId, boolean resend, MessageAcker acker) {
        try {
            String str = (String)body; //new String(bytes, "UTF-8");
            logger.info(str);
            RequestContext rc = RequestContext.get(); //获取上下文
            acker.ack(messageId);
        } catch (Throwable e) {
            boolean discard = e instanceof IOException; //是否废弃这条消息，根据具体场景判断
            if (discard){
                acker.discard(messageId);//废弃
                // 记录废弃原因，并写业务日志
            } else{
                acker.deny(messageId);//告诉mq重发这条消息
                // 记录异常原因，并写业务日志
            }
        }
    }
    //当事务敏感，发送方和消费方不在同一个库时，需要重写设置消费方的routekey
    String getRouteKey() {
        return  "consumerDbRouteKey”;
    }
}
```

mq.send

```java
MessagePublisher pub = MQFactory.get().createSimplePublisher("demo", "demo_queue");
try {
    pub.publish("hello world");
} finally{
    //用完关闭
    pub.close();
}
```

mq.sendwithTX

```java
// Trascation begin
MessagePublisher pub = MQFactory.get().createSimplePublisher("demo", "demo_queue");
try {
    //第一个参数现在不用了，不需要传
    pub.publishInDbTranscation(null, "hello world");
} finally {
    //用完关闭
    pub.close();
}
// Transcation end
```

mq.config

```java
<!-- 项目路径：resources/materialcreatemqconfig.xml， MC环境公共配置项-->
<!--region为云的标识，queue name可自定义，appid为应用标识，class为消费者的全路径-->
<root>
    <region name="scmc" appid="sm">
        <queue name="scmc.sm.sm_salepricelist" sequential="true">
            <consumer class="shkd.yaduo.bamp.basedata.plugin.materialmanager.consumer.SalPriceListConsumer"></consumer>
        </queue>
    </region>
    <region name="scmc" appid="pm">
        <queue name="scmc.pm.pm_purpricelist" sequential="true">
            <consumer class="shkd.yaduo.bamp.basedata.plugin.materialmanager.consumer.PurPriceListConsumer"></consumer>
        </queue>
    </region>
</root>
```
