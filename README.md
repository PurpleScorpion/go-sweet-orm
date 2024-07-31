# go-sweet-orm

## go的持久层框架

### 基础支持
 - beego框架作为基础框架  (github.com/beego/beego/v2 v2.2.1)
 - 目前仅支持mysql和sqlite
### 使用方式

 - 1 引入包 
   ```text
    go get github.com/PurpleScorpion/go-sweet-orm
    
   ```
   ```text
    使用以下语句来引入包
    import "github.com/PurpleScorpion/go-sweet-orm/mapper"
    使用以下函数注册驱动 , 注意 需要先在beego中完成数据库注册
    mapper.InitMapper(mapper.Sqlite, true)
    参数一: 
    mapper.Sqlite代表注册驱动为sqlite3
    mapper.MySQL代表注册驱动为mysql
    参数二:
    第二个参数为是否开启sql展示 true:开启/false:不开启
   ```
 - 2 包含的函数(自动事务)
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
   -  2.1 QueryWrapper的使用
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
 - 3 包含的函数(手动事务)   
   ```text
    若想使用手动事务, 请在增删改之前使用对应的Wapper构建器
   1. BuilderInsertWrapper Insert构建器 , 其中包含两个参数 , 第一个参数与2.1使用方式一致,若是在实际中使用不到第一个参数可填nil , 第二个参数为是否开启自动事务 true:开启/false:不开启
   注: 若开启自动事务 , 请勿忘记使用qw.Commit()进行提交
   1.1 用不到第一个参数的函数列表
      - Update4SQL
      - DeleteById
      - Delete4SQL
      - Insert
      - InsertCustom
      - Insert4SQL
   2. BuilderUpdateWrapper  Update构建器 使用方式同BuilderInsertWrapper
   3. BuilderDeleteWrapper  Delete构建器 使用方式同BuilderInsertWrapper
   
   4. Commit() 提交事务
   5. Rollback() 回滚事务
   
   使用示例:
   qw := mapper.BuilderInsertWrapper(nil,true)
   qw.Insert(&log)
   qw.Commit()
   // 若发生异常
   // qw.Rollback()
   
   ```
 - 4 具体函数使用
   ```text
    具体使用方式已在mapper.demo文件夹下
    其中Logs.go为实体类文件
    Demo.go是使用示例文件
   ```