# microlisp TODO

## 未实现的功能

### Reader / Readtable
- [ ] `set-dispatch-macro-character` 注册的 dispatch 函数未被 parser 调用（# 分发完全硬编码）
- [ ] `read-delimited-list` 未实现

### Destructuring
- [ ] `destructuring-bind` 不支持 `&key` 参数

### Setf 扩展
- [ ] `(defun (setf foo) ...)` 不支持复合函数名
- [ ] `(setf (values ...))` 不支持
- [ ] `(setf (macro-function ...))` 不支持
- [ ] `defsetf` 不支持 `&environment` 参数

### CLOS / 对象系统
- [ ] CLOS method combinations 未实现
- [ ] `find-method` / `remove-method` 未实现
- [ ] `call-next-method` 未实现

### Loop
- [ ] `loop ... from ... downto ... by ...` 语法不支持
- [ ] `loop` hash-key 迭代挂起
- [ ] `loop` destructuring 不完整

### 其他
- [ ] `#p""` pathname 字面量语法未支持
- [ ] `*posix-argv*` 未实现（sbcl 扩展）
- [ ] `*random-state*` 未定义为特殊变量
- [ ] `|...|` 转义符号读取不支持
- [ ] `sb-int:constant-form-value` 不适用（sbcl 特有）
- [ ] `checked-compile` 不适用（sbcl 特有）
- [ ] `char-code-limit` 未定义为常量
- [ ] `#+sb-unicode` 特性不存在

## 已修复的 Bug（来自 sbcl-tests 测试）

1. `macrolet` / `lambda` 不处理 `&rest`/`&body` 参数
2. `assert` 宏逻辑反转
3. `(function name)` 只查全局环境，找不到 `flet`/`labels` 绑定
4. `get-macro-character` 对标准宏字符返回 nil
5. `set-macro-character` 不接受 nil 参数
6. `car`/`cdr` 对 nil 报错而非返回 nil
7. `division-by-zero` 未通过条件系统发出
8. `ecase`/`etypecase`/`ctypecase` 错误信息格式错误
9. `fmt.Errorf` 使用 Lisp 格式符 `~D`/`~A` 而非 Go 的 `%d`/`%v`
10. `defun` 中 `&optional`/`&key`/`&aux` 后的参数未收集
11. compound cons 类型说明符 `(cons integer string)` 未递归验证 car/cdr
12. `null` 未作为 `nil` 类型别名的 type specifier
13. `length` 不支持 VArray 类型
14. EQL specializer 分发完全失效（methodApplicable 只检查 VInstance，specializer 被静默丢弃）
15. `export` 不接受符号列表（只支持单个符号）
16. 缺少 `CL-USER` 包（USER 包未设置 CL-USER 昵称）
17. `export` 不接受字符串参数（`defpackage` 的 `:export` 选项传递字符串时会失败）
18. `cl:NAME` 包限定符号无法解析（CL 包没有导出符号）
19. `nil` 类型说明符被错误当作 `null`（ANSI CL 中 nil 是空类型，null 才匹配 nil）
