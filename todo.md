# microlisp TODO

## 未实现的功能

### Reader / Readtable
- [ ] `set-dispatch-macro-character` 注册的 dispatch 函数未被 parser 调用（# 分发完全硬编码）
- [x] `read-delimited-list` 未实现 — 已实现 `builtinReadDelimitedList`
- [x] `readtable-case` 的 `:preserve` 和 `:invert` 模式未实现（lexer 总是 uppercase 符号名）— 已修复

### Destructuring
- [x] `destructuring-bind` 不支持 `&key` 参数 — 已修复

### Setf 扩展
- [x] `(defun (setf foo) ...)` 不支持复合函数名 — 已修复
- [x] `(setf (values ...))` 不支持 — 已修复
- [x] `(setf (macro-function ...))` 不支持 — 已修复（实现 `builtinMacroFunction` 和 `builtinMacroFunctionSetf`，`expandMacro` 增加 VFunc/VPrim 直接调用分支，`&whole` 机制传递完整宏调用表单）
- [ ] `defsetf` 不支持 `&environment` 参数

### CLOS / 对象系统
- [ ] CLOS method combinations 未实现
- [ ] `find-method` / `remove-method` 未实现
- [ ] `call-next-method` 未实现

### Loop
- [x] `loop ... from ... downto ... by ...` 语法不支持 — 已验证实现（downto/by/above/below 均已支持）
- [ ] `loop` hash-key 迭代挂起
- [ ] `loop` destructuring 不完整

### 其他
- [x] `#p""` pathname 字面量语法未支持 — 已修复（lexer 和 parser 均已实现）
- [x] `*posix-argv*` 未实现（sbcl 扩展）— 已实现（Go 初始化时使用 `os.Args` 填充列表）
- [x] `*random-state*` 未定义为特殊变量 — 已修复（`builtinRandom` 现在会从 `globalEnv` 查找 `*random-state*` 作为默认 rng）
- [x] `|...|` 转义符号读取不支持 — 已修复（添加 lexBarSym，保留大小写，支持 \\ 和 \| 转义）
- [ ] `sb-int:constant-form-value` 不适用（sbcl 特有）
- [ ] `checked-compile` 不适用（sbcl 特有）
- [x] `char-code-limit` 未定义为常量 — 已修复，定义为 1114112（Unicode 码点上限）
- [x] `#+sb-unicode` 特性不存在 — 已修复（添加 `:sb-unicode` 到 features 列表，并修复 feature lookup 大小写不敏感问题）

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
20. `subtypep` 不支持复合 CONS 类型说明符（如 `(cons integer *)`）
21. `compile` 返回值被包装成 VPair 而非 VMultiVal
22. `make-sequence` 不接受 `:initial-element` 关键字
23. `subseq` 对字符串返回 nil 而非子字符串
24. `nreverse` 对列表原地修改但返回错误结果
25. `defmacro` 的 `&optional`/`&key`/`&aux` 参数默认值未求值（参数名被提取，但无参数时使用 nil 而非默认值）
26. `copy-seq` 在向量上调用 Lisp 定义覆盖了 Go 内置函数
27. `assoc` 找不到时返回 `#f` 而非 `nil`
28. `loop` 的 `for x on ...` 无限循环
29. `loop` 的 `with x = value` 子句解析错误
30. `mapcon` 返回错误结果
31. `typep` 缺少 `'vector` 类型检查
32. `typep` 缺少 `'atom` 类型检查
33. `subtypep` 返回 list 而非 VMultiVal（导致 `not` 接收整个列表）
34. 浮点数指数标记（d/D/f/F/s/S/l/L）不被 `parseFloatStr` 支持
35. `ignore-errors` 出错时返回 `(nil . condition)` 而非 `nil`
36. `destructuring-bind` 不支持 `&rest`/`&body`/`&optional`/`&key`（`&rest` 被当作普通变量绑定到错误值）

37. Go 词法分析器对超出 float64 尾数精度的大整数（>2^53）丢失精度
38. setf 对未绑定变量报错而非创建全局绑定（ANSI CL 语义）
39. `destructuring-bind` 的 `&key` 使用位置绑定而非关键字匹配
40. `butlast` 对 n<=0 返回原列表而非副本
41. `block`/`return-from` 不接受 nil 作为块名
42. `eq`/`equal` 不将 nil 符号和 VNil（空列表）视为相等
43. 双反引号嵌套求值错误（`(quasiquote (quasiquote X))` 未正确解包）
44. `unquote`/`unquote-splicing` 在 depth>0 时未递归处理
45. `loop` 的 `for x = expr` 子句在 expr 中引用其他循环变量时报 undefined
46. `load` 不支持 `:if-does-not-exist nil` 关键字参数
47. `stringp`/`numberp` 谓词函数未实现
48. `loop` 不支持 `being each present-symbol/external-symbol of package` 子句
49. `random` 函数接受浮点数参数时总是返回 0（截断为整数导致 rand.Intn(0/1)）
50. `macroexpand` 不展开 quasiquote 形式（返回原始形式不变）
51. `loop` 不支持解构模式如 `(for (a b) in list)`
52. `functionp` 谓词函数未实现
53. ✅ `defun` 接受 `(setf name)` 作为函数名 — 已修复
54. `ignore-errors` 错误时未返回 `(values nil condition)`
55. `nth-value` 无法从 VMultiVal 正确提取第 n 个值
56. `delete-if`/`nsubstitute-if` 谓词函数调用方式错误（eval 而非 callFnOnSeq）
57. `delete-duplicates` 使用指针相等而非值相等判断重复
58. `*random-state*` 未初始化 — 已修复
59. `coerce` 不支持 `'vector` 和 `'array` 结果类型
60. `typep` 不处理复合 `vector` 类型说明符如 `(vector *)`
61. `logand`/`logior`/`logxor` 对非整数参数静默转为0而非报type-error
62. `(setf (values ...) ...)` 未实现
63. `char-name` 对 C1 控制字符（128-159）返回 nil
64. `type-of` 返回 `"unknown"` 对于 `pathname`、`random-state`、`array`、`integer`（大整数）类型
65. `typep`/`subtypep` 类型比较大小写不敏感问题（符号名大写后比较失败）
66. `destructuring-bind` 的 Go 实现中 lambda-list 关键字大小写不匹配（`&rest` vs `&REST`）
67. `set!`/`setq` 不更新 globalEnv 中已定义的全局变量
68. `isNil()` 不识别 VSym "NIL"（导致 length/butlast 等函数对 nil 报错）
69. `find-all-symbols` 函数未实现
70. `coerce` 类型说明符大小写敏感（`'STRING` vs `'string`）
71. 关键字参数大小写不匹配（reader 大写化后 Go 侧用小写匹配）
72. `checked-compile` 宏引用 bug（`eval` 未正确展开变量）
73. `destructuring-bind` 点对模式匹配 nil 值时 Go nil 指针崩溃
74. `string-upcase/downcase/capitalize` 不接受 string designators（符号/字符）
75. `string-capitalize` 不支持 `:start`/`:end` 关键字参数
76. `nstring-upcase/downcase/capitalize` 不支持 VArray 和 fill-pointer
77. `(setf fill-pointer)` 未实现
78. `butlast` 对 dotted list 处理错误
79. `floor`/`ceiling`/`truncate`/`round` 返回 list 而非 VMultiVal（多值应使用 multiVal 而非 list）— 已修复
80. `=`/`/=` 等数值比较对复数只比较实部（`compareNumeric` 忽略虚部）— 已修复
81. `coerce` 不支持 `standard-char`/`base-char` 作为结果类型 — 已修复

## 新发现并修复的 Bug

82. `coerce` 到 `(complex float)` / `(complex single-float)` / `(complex double-float)` 类型说明符时返回实数而非复数（vcomplex 在虚部为0时返回 vnum）— 已修复（新增 vcomplexAlways 函数，coerce 中对 compound type specifier 使用 vcomplexAlways）
83. `coerce` 将列表如 `'(a b c)` 转为字符串时返回空字符串（list-to-string 未处理 VSym/VStr 元素）— 已修复
84. `substitute` 对字符串输入返回列表而非字符串 — 已修复
85. `substitute-if` 对字符串输入返回列表而非字符串 — 已修复
86. `remove` 对字符串输入返回列表而非字符串 — 已修复
87. `remove-if` 对字符串输入返回列表而非字符串 — 已修复
88. `delete` 不支持字符串输入（只接受 VPair）— 已修复（对字符串委托给 remove）
89. `delete-if` 不支持字符串输入（只接受 VPair）— 已修复（对字符串委托给 remove-if）
90. `round` 的 two-argument 形式（`round x d`）未使用 round-half-to-even 规则 — 已修复
91. `lambda`/`defun` 中 `&optional` 和 `&key` 参数的默认值不生效（parseParams 中 elem 变量作用域问题）— 已修复（重构 parseParams 和 apply 函数，正确处理可选/关键字参数的默认值和绑定）
