# go-sweet-orm

## go的持久层框架 V3

### 基础支持
 - gorm作为基础框架  (gorm.io/gorm v1.31.1)
 - 目前仅支持mysql和sqlite
 - go版本1.25.6及以上
### 使用方式

 - 1 引入包 
   ```text
    go get github.com/PurpleScorpion/go-sweet-orm/v3@latest
    
   ```
   ```text
    使用以下语句来引入包
    import "github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
   ```
 - 2 注册数据库连接
     ```text
     mapper.SetMySqlConf(mapper.MySQLConf{
        UserName: "root",
        Password: "root",
        DbName:   "demo",
        Port:     3308,
        Host:     "localhost",
    })
    mapper.RegisterMySql()
    
    配置项如下所示
    type MySQLConf struct {
        UserName     string // 用户名
        Password     string // 密码
        DbName       string // 数据库名称
        Port         int   // 端口 不填默认3306
        Host         string // 数据库地址
        Charset      string // 字符集 不填默认 utf8mb4
        Loc          string // 时区 不填默认LOCAL
        MaxIdleConn  int // 最大空闲连接数 不填默认10
        MaxOpenConn  int // 最大连接数 不填默认100
        TlsCertPool  *x509.CertPool // 若MYSQL需要证书,则需要配置根证书池
    } 
    若你的MySQL需要证书,则需要配置根证书池
    rootCertPool := x509.NewCertPool()
	pem, _ := ioutil.ReadFile("你的根证书路径")
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		panic("Failed to append PEM.")
	}
   
     ```
 - 3 Wrapper介绍
   - 3.1 QueryWrapper
   ```text
    1. 创建查询器 
    qw := mapper.BuilderQueryWrapper()
    2. 具体函数: 请查看 <5 查询器函数介绍>
   ```
   - 3.2 UpdateWrapper
   ```text
    1. 创建更新器(用于增删改)
        详细使用请参考 
            [DeleteDemo_test.go](demo/DeleteDemo_test.go)
            [InsertDemo_test.go](demo/InsertDemo_test.go)
            [UpdateDemo_test.go](demo/UpdateDemo_test.go)
        实际使用方式与mybatis-plus类似
    2. 构造器-无事务
        wrapper := mapper.BuilderUpdateWrapper(false)
    3. 构造器-开启事务
        wrapper := mapper.BuilderUpdateWrapper(true)
        
    
   ```
 - 5 其他函数
   ```text
    SelectById[泛型](id) []泛型 , 
        根据id查询,返回泛型数组(为了避免差不到报错的问题 , 直接返回一个数组)
   
    SelectList[泛型](QueryWrapper) []泛型 ,
        根据QueryWrapper条件来查询数据列表 , 为了避免遍历表中的所有数据造成事故 , 这里限制了QueryWrapper不能为nil
   
    SelectCount(QueryWrapper) 返回值为int
        根据QueryWrapper条件来查询数据数量,若想查表中数量 , 参数可以传递nil
   
    SelectList4SQL[泛型](,sql,values...) []泛型
        使用原生sql查询数据列表,values是可变参数列表,用来替换`?`占位符
        若不理解什么是`?`占位符,请先去学习java jdbc
   
    SelectCount4SQL(sql,values...) 返回值为int
        使用原生sql查询数据数量,values是可变参数列表,用来替换`?`占位符
   
    Update[泛型](UpdateWrapper) 返回值为int64 (更新影响的行数)
        根据UpdateWrapper条件来更新数据
        注意此时UpdateWrapper必须使用Set函数进行设置待更新的字段,否则会抛出异常
   
    Update4SQL(UpdateWrapper) 返回值为int64 (更新影响的行数)
        需调用函数 wrapper.SQL(sql,values...)
        原生sql进行更新
        sql: 原生sql
        values... : 可变参数, 用于替换 `?` 占位符
   
    DeleteById[泛型](id,UpdateWrapper) 返回值为int64 (删除影响的行数)
        根据id删除,UpdateWrapper可传nil , 有UpdateWrapper的目的是为了开启事务
   
    Delete[泛型](UpdateWrapper) 返回值为int64 (删除影响的行数)
        根据UpdateWrapper条件来删除数据
   
    Delete4SQL(UpdateWrapper) 返回值为int64 (删除影响的行数)
        需调用函数 wrapper.SQL(sql,values...)
        原生sql进行更新
        sql: 原生sql
        values... : 可变参数, 用于替换 `?` 占位符
   
    Insert[泛型](&entity,UpdateWrapper) 返回值为int64 (新增影响的行数)
        UpdateWrapper可为nil
        wrapper.SetExcludeEmpty(true/false) 用来设置是否排除空值, true:排除 (int: 0,string: "")这种空值 , 默认false
        wrapper.SetExcludeField(excludeField ...string) 用来设置是否的字段 , 可多次调用
        其中默认自增,若不想使用自增, 请使用wrapper.CloseAutoId()
        若嫌太麻烦 , 则可以使用mapper.CloseGlobalAutoId() , 这样全局都不会出现自增ID
        当然 , 若有个别的表是AutoId, 则可以使用wrapper.OpenAutoId() , 对此次操作开启
   
    InsertAll[泛型](&entityList,UpdateWrapper) 返回值为int64 (新增影响的行数) 
        UpdateWrapper可为nil
   
    Insert4SQL(UpdateWrapper) 返回值2个
        原生sql插入数据方式
        需调用函数 wrapper.SQL(sql,values...)
        sql: 原生sql
        values... : 可变参数, 用于替换 `?` 占位符
        其中autoId受UpdateWrapper影响
   
        返回值1: 影响的行数
        返回值2: 插入后的自增主键值 , 受autoId参数影响
   
   
    Page[泛型](PageUtils) 分页工具 返回值为PageData对象
        两种使用方式 , 为了美观 , 可以选择被注释的方式
        //pageUtils := mapper.BuilderPageUtils(1, 2, mapper.BuilderQueryWrapper())
	    //page := mapper.Page[user](pageUtils)
   
	    page := mapper.Page[user](
		    mapper.BuilderPageUtils(1, 2, 
			    mapper.BuilderQueryWrapper(),
		    ),
	    )

	    logger.Info("当前页: {}", page.Current)
	    logger.Info("总页数: {}", page.TotalPage)
	    logger.Info("总记录数: {}", page.TotalCount)

	    list := page.List
	    for _, u := range list {
		    logger.Info("姓名: {} , 年龄: {}", u.UserName, u.Age)
	    }    
    ConvertPageData[泛型1, 泛型2](data PageData[泛型1], list []泛型2) PageData[泛型2]
        其作用为有时需将查询结果进行转换,比如查询结果为PageData对象,而你希望将其转换成PageData[T]对象,那么就可以使用此函数
        当然 , 你自己手动创建PageData对象 , 然后手动赋值给新的PageData对象也是可以的
   
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
    具体使用方式已在demo文件夹下
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
 - 9 Add组与Or组使用
  ```text

        使用示例
    func TestNestedQuery(t *testing.T) {
        // 构建复杂嵌套查询
        wrapper := mapper.BuilderQueryWrapper()
        result := wrapper.Eq(true, "name", "jinzhu").
            Eq(true, "sex", 1).
            And(
                mapper.NewAndGroup().
                    Eq(true, "name", "jinzhu 2").
                    Eq(true, "age", 18),
            ).
            And(
                mapper.NewOrGroup().
                    Eq(true, "word1", "bbbb").
                    Eq(true, "word2", "aaa"),
            ).
            Or(
                mapper.NewAndGroup().
                    Eq(true, "key1", "123").
                    Eq(true, "key2", 456),
            ).
            Or(
                mapper.NewOrGroup().
                    Eq(true, "key3", "789").
                    Eq(true, "key4", 666),
            )
        
        生成的sql语句如下所示
        Generated SQL:  and name = ?  and sex = ?  AND (name = ? AND age = ?) AND (word1 = ? OR word2 = ?) OR (key1 = ? AND key2 = ?) OR (key3 = ? OR key4 = ?)
    }
```
