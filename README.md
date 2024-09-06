# go-sweet-orm

## go的持久层框架 V2

### 基础支持
 - beego框架作为基础框架  (github.com/beego/beego/v2 v2.2.1)
 - 目前仅支持mysql和sqlite
### 使用方式

 - 1 引入包 
   ```text
    go get github.com/PurpleScorpion/go-sweet-orm/v2@latest
    
   ```
   ```text
    使用以下语句来引入包
    import "github.com/PurpleScorpion/go-sweet-orm/v2/mapper"
   ```
 - 2 注册数据库连接
     ```text
     mapper.Register(activeDB, connStr string, params ...int)  
         activeDB: 激活的数据库, 目前仅支持 mapper.MySQL 和 mapper.Sqlite
         connStr: 数据库连接字符串 (注意mysql是包含用户名密码的连接字符串, 但是Sqlite却是文件地址)
             mysql连接字符串: connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", username, password, host, port, dbName) 
         params: 可变参数, 但是目前仅前两个参数有效 , 可选项, 第一个为MaxIdleConns(默认50) , 第二个为MaxOpenConns(默认100) 
       
     使用示例:
       connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", username, password, host, port, dbName) 
       mapper.Register(mapper.MySQL, connStr)
     ```
 - 3 Wrapper介绍
   - 3.1 QueryWrapper
   ```text
    1. 创建查询器 
    var log []Logs
    qw := mapper.BuilderQueryWrapper(&log)
    注意, 此处必须使用指针参数 , 否则无法将结果映射到entity中
    2. 具体函数: 请查看 <5 查询器函数介绍>
   ```
   - 3.2 UpdateWrapper
   ```text
    1. 创建更新器(用于增删改)
        更新器大部分都可以使用匿名对象的形式进行传参, 但是执行Insert的时候必须传入指针对象, 否则若是自增ID则无法回填到entity中
        第二个参数其实是个可变参数 , 但是此处只接受1个参数, 该参数的作用是是否开启事务, 默认为false(不开启)
    2. 新增示例(构造器-无事务)
        var log Logs
        log.Title = "test"
        qw := mapper.BuilderUpdateWrapper(&log)
        qw.Insert()
    3. 新增示例(构造器-开启事务)
        var log Logs
        log.Title = "test"
        qw := mapper.BuilderUpdateWrapper(&log,true)
        count := qw.Insert()
        if (count == 0){
            qw.Rollback()
        }else {
            qw.Commit()
        }
    4. 新增示例(便捷插入-永无事务)
        var log Logs
        log.Title = "test"
        mapper.Insert(&log)
    5. 新增示例(构造器-批量插入)
        具体使用查看 Demo_test.go 下的 TestDemo1() 测试用例
    6. 删除示例(构造器-无事务)
        具体使用查看 Demo.go 下的 demo8() 测试用例
    7. 删除示例(构造器-开启事务)
        以此类推即可....
    8. 删除示例(根据ID删除)
        具体使用查看 Demo.go 下的 demo7() 测试用例
    9. 删除示例(根据ID批量删除)
        具体使用查看 Demo.go 下的 demo7_1() 测试用例
    10. 删除示例(根据ID批量删除-便捷使用)
        具体使用查看 Demo.go 下的 demo7_2() 测试用例
    11. 修改示例(构造器-无事务)
        以此类推即可....
    12. 修改示例(构造器-开启事务)
        以此类推即可....
   ```
 - 5 其他函数
   ```text
    SelectById(entity,id) 无返回值 , 查到的结果会直接映射到entity中
        根据id查询,传入的entity可以是单个对象,也可以是一个数组,但是建议使用数组
    QueryWrapper的使用请看2.1
    SelectList(QueryWrapper) 无返回值 , 查到的结果会直接映射到entity中
        根据QueryWrapper条件来查询数据列表
    SelectCount(QueryWrapper) 返回值为int64
        根据QueryWrapper条件来查询数据数量
    SelectList4SQL(entity,sql,values...) 无返回值 , 查到的结果会直接映射到entity中
        使用原生sql查询数据列表,values是可变参数列表,用来替换`?`占位符
        若不理解什么是`?`占位符,请先去学习java jdbc
    SelectCount4SQL(entity,sql,values...) 返回值为int64
        使用原生sql查询数据数量,values是可变参数列表,用来替换`?`占位符
    Update(QueryWrapper) 返回值为int64 (更新影响的行数)
        根据QueryWrapper条件来更新数据
        注意此时QueryWrapper必须使用Set函数进行设置待更新的字段,否则会抛出异常
    Update4SQL(sql,values...) 返回值为int64 (更新影响的行数)
        原生sql进行更新
        sql: 原生sql
        values... : 可变参数, 用于替换 `?` 占位符
    DeleteById(entity,id) 返回值为int64 (删除影响的行数)
        根据id删除,传入的entity必须是单个对象,不可以是数组
    Delete(QueryWrapper) 返回值为int64 (删除影响的行数)
        根据QueryWrapper条件来删除数据
    Delete4SQL(sql,values...) 返回值为int64 (删除影响的行数)
        原生sql进行删除
        sql: 原生sql
        values... : 可变参数, 用于替换 `?` 占位符
    Insert(entity,excludeField...) 返回值为int64 (新增影响的行数)
        根据传入实体类中的值进行新增 , 新增后该对象中会有新增后的id
        默认自增主键 , 空值排除(若为0或""等空值,则默认不向数据库中添加该值)
        excludeField为可变参数列表,传入的值为新增时忽视的字段
    InsertCustom(entity,autoId,excludeEmpty,excludeField...) 返回值为int64 (新增影响的行数)
        根据传入实体类中的值进行新增 若autoId为true 则新增后该对象中会有新增后的id
        autoId: 是否为自增主键 true:自增主键/false:自定义主键
        excludeEmpty: 是否进行空值排除 true:空值排除/false:空值仍然存储
        excludeField为可变参数列表,传入的值为新增时忽视的字段
    Insert4SQL(autoId,sql,values...) 返回值2个
        原生sql插入数据方式
        autoId: 是否是自增主键 true: 自增主键,返回值将返回插入后的自增主键值 / false: 不是自增主键,返回值中的自增主键值为0
        sql: 原生sql
        values... : 可变参数, 用于替换 `?` 占位符
        返回值1: 影响的行数
        返回值2: 插入后的自增主键值 , 受autoId参数影响
    Page(PageUtils) 分页工具 返回值为PageData对象
        若想使用分页工具,必须按以下流程进行编写代码
        1. 创建结果集 var resLog []Logs 必须: 必须是数组类型
        2. 创建查询器 qw := mapper.BuilderQueryWrapper(&resLog)
        3. 创建分页器 page := mapper.BuilderPageUtils(thisPage, pageSize, qw)
        4. 调用方法进行分页查询 pageData := mapper.Page(page) 必须: 必须接收返回值
            注: 可自行修改 PageData 类中的builder函数和Field来完成符合自己需求的分页结果集
   ```
- 6 查询器函数介绍
     - 如果你熟悉Java中的MyBatis的QueryWrapper,那本教程你将会很轻松
       - QueryWrapper的创建: `mapper.BuilderQueryWrapper(&log)`
           - 创建结果集接收对象 `var log []Logs` `var log Logs`
           - 传入的参数`&log`是用来接收结果集的对象
       - QueryWrapper的操作:
         - Eq(flag,column,value) `释义: 等于`
           - 参数1: flag , 是否执行此比较操作,若值为false,即使填充了内容,查询器也会无视
           - 参数2: column, 数据库列名 , 注意是数据库的
           - 参数3: value, 需要比较的值
         - Ne(flag,column,value) `释义: 不等于`
           - 参数列表: 与Eq一致
         - Gt(flag,column,value) `释义: 大于`
           - 参数列表: 与Eq一致
         - Ge(flag,column,value) `释义: 大于等于`
             - 参数列表: 与Eq一致
         - Lt(flag,column,value) `释义: 小于`
           - 参数列表: 与Eq一致
         - Le(flag,column,value) `释义: 小于等于`
           - 参数列表: 与Eq一致
         - Like(flag,column,value) `释义: 模糊查询`
           - 参数列表: 与Eq一致 注意,查询器会自动添加%, 例如: `%张三%`
         - LikeLeft(flag,column,value) `释义: 模糊查询`
           - 参数列表: 与Eq一致 注意,查询器会自动添加%, 例如: `%张三`
         - LikeRight(flag,column,value) `释义: 模糊查询`
           - 参数列表: 与Eq一致 注意,查询器会自动添加%, 例如: `张三%`
         - NotLike(flag,column,value) `释义: 模糊查询`
           - 参数列表: 与Eq一致 注意,查询器会自动添加%, 例如: `%张三%`
         - IsNull(flag,column) `释义: 空值判断`
           - 参数列表: 与Eq一致,但是没有value参数
         - IsNotNull(flag,column) `释义: 非空判断`
           - 参数列表: 与Eq一致,但是没有value参数
         - In(flag,column,value...) `释义: mysql中in函数`
           - 参数列表: 与Eq一致,但是该value为可变参数,且要求参数类型必须全部一致
         - NotIn(flag,column,value...) `释义: mysql中 not in 函数`
           - 参数列表: 与Eq一致,但是该value为可变参数,且要求参数类型必须全部一致
         - Between(flag,column,from,to) `释义: mysql中 between and 函数`
           - 参数列表: 前两个参数与Eq一致
           - 参数3: string类型
           - 参数4: string类型
         - OrderByAsc(flag,column) `释义: 正序排序`
           - 参数1: flag , 是否执行此比较操作,若值为false,即使填充了内容,查询器也会无视
           - 参数2: column, 数据库列名 , 注意是数据库的
         - OrderByDesc(flag,column) `释义: 倒序排序`
           - 参数列表: 与OrderByAsc一致
         - OrderByTimeAsc(flag,column)  `释义: 仅用于处理时间格式的排序`
           - 参数列表: 与OrderByAsc一致
         - OrderByTimeDesc(flag,column)  `释义: 仅用于处理时间格式的排序`
           - 参数列表: 与OrderByAsc一致
         - LastSql(sql) `释义: 永远拼接在最后位置的sql文本`
           - 参数列表: string类型的自定义sql
         - Set(flag,column,value) `释义: 仅使用更新函数时该方法生效`
           - 参数1: flag , 是否执行此字段的更新
           - 参数2: column, 数据库列名 , 注意是数据库的
           - 参数3: value, 更新的值
 - 7 具体函数使用
   ```text
    具体使用方式已在mapper.demo文件夹下
    其中Logs.go为实体类文件
    Demo.go是使用示例文件
   ```
 - 8 其他函数的使用
   - OpenLog() 开启日志
   - OpenTransaction() 开启全局事务. 这样的话BuilderUpdateWrapper第二个参数不填也是true, 但若是填了,则以你填的为准
   - GetOrm() 获取orm对象(禁止自己去NewOrm), 若使用,则直接通过该函数获取即可
   - SpreadTransaction() 事务传播, 具体可看demo14()
   - SetTransaction() 设置事务
   - GetTransaction() 获取事务
   - BenginTransaction() 开启事务
   - Commit() 提交事务
   - Rollback() 回滚事务