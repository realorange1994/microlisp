# microlisp TODO

## 未实现的功能

### Reader / Readtable
- [x] `set-dispatch-macro-character` 注册的 dispatch 函数未被 parser 调用（# 分发完全硬编码）— 已修复（添加 TSharpMacro token 类型，lexer 检查 currentReadtable.dispFns，parser 调用 invokeDispatchMacro）
- [x] `read-delimited-list` 未实现 — 已实现 `builtinReadDelimitedList`
- [x] `readtable-case` 的 `:preserve` 和 `:invert` 模式未实现（lexer 总是 uppercase 符号名）— 已修复

### Destructuring
- [x] `destructuring-bind` 不支持 `&key` 参数 — 已修复
- [x] `destructuring-bind` 不支持 `&key` supplied-p 变量 — 已修复（Bug #98）
- [x] `destructuring-bind` `&key` 默认值不生效 — 已修复（Bug #99）

### Setf 扩展
- [x] `(defun (setf foo) ...)` 不支持复合函数名 — 已修复
- [x] `(setf (values ...))` 不支持 — 已修复
- [x] `(setf (macro-function ...))` 不支持 — 已修复（实现 `builtinMacroFunction` 和 `builtinMacroFunctionSetf`，`expandMacro` 增加 VFunc/VPrim 直接调用分支，`&whole` 机制传递完整宏调用表单）
- [x] `defsetf` 不支持 `&environment` 参数 — 已修复（添加 `remove-env` 辅助函数过滤 `&environment`，修复 `-SETF` 后缀大小写匹配问题）

### CLOS / 对象系统
- [x] CLOS method combinations 未实现 — 已修复（实现 standard/progn/and/or/list/append/nconc/min/max/+ 方法组合，添加 methodCombo 字段到 VGeneric，set-method-combination 内置函数，defgeneric 支持 :method-combination 选项，isTypeSpecializerMatch 和 methodSpecificity 使用大写类型名匹配）
- [x] `find-method` / `remove-method` 未实现 — 已修复（添加 VMethod 类型，实现 builtinFindMethod 和 builtinRemoveMethod，支持 qualifier/specializer-list/errorp 参数，特化器匹配处理大小写和 t="" 等价）
- [ ] `call-next-method` 未实现 — 已确认可用（defmethod + CLOS dispatch 中已有 call-next-method/next-method-p 绑定实现，支持方法链调用）

### Loop
- [x] `loop ... from ... downto ... by ...` 语法不支持 — 已验证实现（downto/by/above/below 均已支持）
- [x] `loop` hash-key 迭代挂起 — 已修复（实现 `hash-table-keys`/`hash-table-values` 函数，loop 宏的 `being` 子句支持 `hash-keys`/`hash-key`/`hash-values`/`hash-value`，支持 `using (hash-value v)` 并行绑定）
- [ ] `loop` destructuring 不完整 — 注：`for (a b) on list` 的析构模式在 loop 宏中未实现（loop 宏已被禁止修改）
- [x] `loop for-across` 未实现（遍历数组）— 已修复（转换为 idx from 0 below length + body set var to aref）

### 其他
- [x] `#p""` pathname 字面量语法未支持 — 已修复（lexer 和 parser 均已实现）
- [x] `*posix-argv*` 未实现（sbcl 扩展）— 已实现（Go 初始化时使用 `os.Args` 填充列表）
- [x] `*random-state*` 未定义为特殊变量 — 已修复（`builtinRandom` 现在会从 `globalEnv` 查找 `*random-state*` 作为默认 rng）
- [x] `|...|` 转义符号读取不支持 — 已修复（添加 lexBarSym，保留大小写，支持 \\ 和 \| 转义）
- [ ] `sb-int:constant-form-value` 不适用（sbcl 特有）
- [ ] `checked-compile` 不适用（sbcl 特有）
- [x] `char-code-limit` 未定义为常量 — 已修复，定义为 1114112（Unicode 码点上限）
- [x] `#+sb-unicode` 特性不存在 — 已修复（添加 `:sb-unicode` 到 features 列表，并修复 feature lookup 大小写不敏感问题）
- [x] `with-standard-io-syntax` 宏未实现 — 已修复（Go 初始化代码中添加 *print-radix*、*print-readably*、*print-gensym*、*print-array*、*read-base*、*read-default-float-format*、*read-eval*、*read-suppress* 等 IO 变量，defmacro 用 unwind-protect 实现 save/restore 模式）
- [x] `with-hash-table-iterator` 宏未实现 — 已修复（实现为 macrolet 包装，内部绑定 keys 列表和 index 计数器，每次调用返回 (values key value present-p) 三值）

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
21. `compile` 返回值被包装成 VPair 而非 VMultiVal — 已修复（compile 返回 multiVal(fn, vnil(), vnil())）
22. `make-sequence` 不接受 `:initial-element` 关键字 — 已修复（实现 `builtinMakeSequence`，支持 list/vector/string/bit-vector 类型及 `:initial-element` 关键字）
23. `subseq` 对字符串返回 nil 而非子字符串
24. `nreverse` 对列表原地修改但返回错误结果
25. `defmacro` 的 `&optional`/`&key`/`&aux` 参数默认值未求值（参数名被提取，但无参数时使用 nil 而非默认值）— 已修复（parseMacroParams 返回 optDefaults/keySpecs，expandMacro 绑定可选/关键字参数时评估默认值表达式，defmacro/define-macro/macrolet 均存储 optDefaults 和 keySpecs）— 已修复（parseMacroParams 返回 optDefaults/keySpecs，expandMacro 绑定可选/关键字参数时评估默认值表达式）
26. `copy-seq` 在向量上调用 Lisp 定义覆盖了 Go 内置函数
27. `assoc` 找不到时返回 `#f` 而非 `nil`
28. `loop` 的 `for x on ...` 无限循环
29. `loop` 的 `with x = value` 子句解析错误
30. `mapcon` 返回错误结果
31. `typep` 缺少 `'vector` 类型检查
32. `typep` 缺少 `'atom` 类型检查
33. `subtypep` 返回 list 而非 VMultiVal（导致 `not` 接收整个列表）
34. 浮点数指数标记（d/D/f/F/s/S/l/L）不被 `parseFloatStr` 支持
35. `ignore-errors` 出错时返回 `(nil . condition)` 而非 `nil` — 已修复（错误时返回 multiVal(vnil(), condition)，成功时返回 multiVal(result, vnil())）
36. `destructuring-bind` 不支持 `&rest`/`&body`/`&optional`/`&key`（`&rest` 被当作普通变量绑定到错误值）— 已修复（在 &optional 内部循环中添加 lambda-list 关键字检测，遇到 &rest/&key/&aux 时 break；将 &optional 循环后的 return nil 替换为 continue 让外层循环处理 &rest）

37. Go 词法分析器对超出 float64 尾数精度的大整数（>2^53）丢失精度 — 已修复（`compareNumeric` 添加 `toBigIntExact` 和 `toBigRat` 辅助函数，使用 `big.Int.Cmp` 和 `big.Rat.Cmp` 进行精确比较）
38. setf 对未绑定变量报错而非创建全局绑定（ANSI CL 语义）— 已修复（set! 和 eval 中的 SETQ 分支在 globalEnv 未绑定时自动创建全局绑定）
39. `destructuring-bind` 的 `&key` 使用位置绑定而非关键字匹配 — 已修复（bindPatternRec 中 &key 分支使用 keyword matching 而非 positional）
40. `butlast` 对 n<=0 返回原列表而非副本 — 已修复（Go 实现已正确返回副本；移除覆盖的 Lisp 定义）
41. `block`/`return-from` 不接受 nil 作为块名 — 已修复（block 和 return-from 均处理 VNil 作为块名，转换为 "NIL"）
42. `eq`/`equal` 不将 nil 符号和 VNil（空列表）视为相等 — 已修复（eq 和 equal 均处理 VSym "NIL" 与 VNil 的等价）
43. 双反引号嵌套求值错误（`(quasiquote (quasiquote X))` 未正确解包）— 已修复（evalQuasiquote 重写：QUASIQUOTE 递归 depth+1 并包装结果；UNQUOTE/UNQUOTE-SPLICING 在 depth>0 时递归 depth-1 并包装为 (UNQUOTE ...)；修复 depth==1 时错误求值 UNQUOTE 参数的问题；所有 vsym 使用大写符号名与 reader 一致）
44. `unquote`/`unquote-splicing` 在 depth>0 时未递归处理 — 已修复（与 Bug #43 同修复，统一 depth>0 路径，移除了错误的 depth==1 特殊分支）
45. `loop` 的 `for x = expr` 子句在 expr 中引用其他循环变量时报 undefined
46. `load` 不支持 `:if-does-not-exist nil` 关键字参数 — 已确认已实现（builtinLoad 已处理 :if-does-not-exist 和 :if-exists 关键字参数）
47. `stringp`/`numberp` 谓词函数未实现 — 已修复（builtinStringP 和 builtinNumberP 已注册）
48. `loop` 不支持 `being each present-symbol/external-symbol of package` 子句 — 已确认已实现（loop 宏支持 being 子句和 present-symbol/external-symbol）
49. `random` 函数接受浮点数参数时总是返回 0（截断为整数导致 rand.Intn(0/1)）— 已修复（builtinRandom 对 VNum 区分浮点和整数路径，浮点参数使用 rand.Float64）
50. `macroexpand` 不展开 quasiquote 形式（返回原始形式不变）— 已确认已实现（builtinMacroexpand 处理 quasiquote 分支返回 `(quote expanded)`)
51. `loop` 不支持解构模式如 `(for (a b) in list)` — 已确认已实现（loop 宏已支持 destr-specs/destructuring-bind 包装解构模式）
52. `functionp` 谓词函数未实现 — 已确认已实现（builtinFunctionP 正确判断 VPrim/VFunc/VGeneric）
53. `defun` 接受 `(setf name)` 作为函数名 — 已修复
54. `ignore-errors` 错误时未返回 `(values nil condition)` — 已修复（成功时返回 multiVal(result, vnil())，错误时返回 multiVal(vnil(), condition)）
55. `nth-value` 无法从 VMultiVal 正确提取第 n 个值 — 已修复（nth-value 特殊形式使用 multiValList 获取多值列表并正确索引）
56. `delete-if`/`nsubstitute-if` 谓词函数调用方式错误（eval 而非 callFnOnSeq）— 已修复（使用 callFn 正确调用谓词函数，支持 #'evenp 等函数引用）
57. `delete-duplicates` 使用指针相等而非值相等判断重复 — 已修复（使用 toString 进行值比较；添加 VArray 和 VStr 委托给 remove-duplicates）
58. `*random-state*` 未初始化 — 已修复
59. `coerce` 不支持 `'vector` 和 `'array` 结果类型 — 已修复（coerce 的 vector/array 分支支持列表和字符串输入）
60. `typep` 不处理复合 `vector` 类型说明符如 `(vector *)` 且不识别字符串为 vector/array — 已修复（字符串是 CL 中的 vector 和 array 子类型）
61. `logand`/`logior`/`logxor` 对非整数参数静默转为0而非报type-error — 已修复（已有 `isIntegerValue` 检查和 `signalTypeError` 返回）
62. `(setf (values ...) ...)` 未实现 — 已修复（setf 展开器处理 values 子形式）
63. `char-name` 对 C1 控制字符（128-159）返回 nil — 已修复（返回 "C128"、"C129" 等实现定义名称，C0 未命名控制字符也返回 "C0"、"C1" 等）
64. `type-of` 返回 `"unknown"` 对于 `pathname`、`random-state`、`array`、`integer`（大整数）类型 — 已修复（typeStr 返回正确类型名称；typepCheckRec 符号分支添加 PATHNAME、RANDOM-STATE、PACKAGE、READTABLE、METHOD、RESTART、GENERIC、INSTANCE、HASH-TABLE、CHARACTER、STREAM、CLASS、MACRO、BOOLEAN、SEQUENCE、ATOM、RATIONAL、COMPLEX 类型检查；复合类型说明符分支添加相同类型检查并修复缩进）
65. `typep`/`subtypep` 类型比较大小写不敏感问题（符号名大写后比较失败）— 已修复（subtypepChecks 使用 strings.ToUpper 标准化类型名后比较，simpleSubtype 也使用大写比较）
66. `destructuring-bind` 的 Go 实现中 lambda-list 关键字大小写不匹配（`&rest` vs `&REST`）— 已确认已实现（bindPatternRec 使用 strings.ToUpper 比较）
67. `set!`/`setq` 不更新 globalEnv 中已定义的全局变量 — 已确认已实现（eval 中的 SETQ 分支检查 globalEnv.Get 并更新）
68. `isNil()` 不识别 VSym "NIL"（导致 length/butlast 等函数对 nil 报错）— 已修复（isNil 函数添加 VSym 且 strings.EqualFold(v.str, "nil") 检查）
69. `find-all-symbols` 函数未实现 — 已确认已实现（builtinFindAllSymbols 已存在于代码中）
70. `coerce` 类型说明符大小写敏感（`'STRING` vs `'string`）— 已确认已实现（builtinCoerce 使用 strings.ToLower 标准化类型名）
71. 关键字参数大小写不匹配（reader 大写化后 Go 侧用小写匹配）— 已确认已实现（seqParseKeys 使用大写符号名比较 ":KEY", ":TEST" 等）
72. `checked-compile` 宏引用 bug（`eval` 未正确展开变量）— 不适用（sbcl 特有）
73. `destructuring-bind` 点对模式匹配 nil 值时 Go nil 指针崩溃 — 已修复（添加 nil/VNil 元素检查）
74. `format ~e` 科学记数法指数有前导零（`(format nil "~e" 42.0)` 返回 "4.2E+01" 而非 "4.2E+1"）— 已修复（重写 ~E 分支：使用 strconv.FormatFloat 正确精度，去除指数前导零，确保尾数有小数点）
75. `string-capitalize` 不支持 `:start`/`:end` 关键字参数 — 已修复
76. `nstring-upcase/downcase/capitalize` 不支持 VArray 和 fill-pointer — 已修复（VArray 分支修改 arr.elements 中的 VChar 元素，尊重 fillPtr 限制）
77. `(setf fill-pointer)` 未实现 — 已修复（fill-pointer-setf 已注册在 builtin table）
78. `butlast` 对 dotted list 处理错误 — 已修复（重写 n>0 路径：收集 elems/cdrs，使用 cdrs[len-1] 作为 dottedTail；移除覆盖 Go 内置的 Lisp 定义）
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

92. `evalQuasiquote` 对 `UNQUOTE`/`UNQUOTE-SPLICING`/`QUASIQUOTE` 符号名大小写敏感匹配（reader 产生小写符号，Go 代码用大写字符串比较导致逗号/splice 在 backquote 内不被识别）— 已修复（使用 strings.EqualFold 进行大小写不敏感比较）

93. `expt` 对复数基底的整数指数幂返回错误结果（`(expt #c(0 1) 2)` 返回 0 而非 -1）— 已修复（新增 VComplex 分支，使用二进制幂法计算复数整数幂）

94. `arrayToString` 对 nil 数组元素崩溃（未初始化数组的 nil 元素导致 elem.typ 访问 Go nil 指针）— 已修复（添加 nil/VNil 元素检查）

95. `parseParams` 对 lambda 点对语法 `(lambda (a . rest) ...)` 中 rest 参数未捕获（`((lambda (a . rest) rest) 1 2 3)` 报 "undefined: REST"）— 已修复（在 parseParams 循环体开头增加 VSym 检测，当 v 变为 VSym 时作为 &rest 参数返回）

96. `seqParseKeys` 完全不支持 `:test-not` 关键字参数（导致 `member`、`find`、`position`、`count`、`remove`、`substitute` 等函数无法使用 `:test-not`，且 `:test` 测试函数参数顺序与 CL 规范相反（`(element item)` 应为 `(item element)`））— 已修复（`seqParseKeys` 增加 `testNotFn` 返回值，更新所有 18 个调用方添加额外的 `_` 忽略该值，`builtinMember` 和 `testItemMatchFull` 正确实现 `:test-not` 语义和 CL 规范参数顺序）

97. `assoc` 函数为 Lisp 简易实现（仅使用 `equal?` 比较），不支持 `:test`、`:test-not`、`:key` 关键字参数（`assoc-if` 的 Go 版本已实现但不完整）— 已修复（添加 `builtinAssoc` Go 实现，支持完整的 `:test`、`:test-not`、`:key` 参数，移除 Lisp 定义）

98. `destructuring-bind` 不支持 `&key` 的 supplied-p 变量（如 `(destructuring-bind (&key (x (funcall x) x-supplied)) () ...)` 报 "undefined: X-SUPPLIED"）— 已修复（bindPatternRec 中 &key 参数解析已正确处理 (var default supplied-p) 三元素形式的 supplied-p 符号绑定和默认值求值）

99. `destructuring-bind` 的 `&key` 默认值不生效（如 `(destructuring-bind (&key (a 99)) () a)` 返回 `nil` 而非 `99`）— 已修复（bindPatternRec 中 &key 分支在 keyValMap 未找到匹配时正确 eval 默认值并设置 supplied-p 为 false）

100. `character` 函数未实现（ANSI CL 标准函数，接受字符设计符返回字符）— 已修复（添加 `builtinCharacter`）

101. `constantp` 函数未实现（ANSI CL 标准函数，检查形式是否为常量）— 已修复（添加 `builtinConstantp` 和 `isConstant` 辅助函数）

102. `coerce` 的 `character` 类型不支持符号设计符（如 `(coerce 'a 'character)` 应返回 `#\A`）— 已修复（添加 VSym 单字符符号支持）

103. `coerce` 的 `character` 类型对多字符字符串应报错而非静默取首字符 — 已修复（添加长度>1的错误检查）

104. `coerce` 不支持 `simple-vector` 结果类型 — 已修复（添加到 `vector` case 的同组处理）

105. `incf`/`decf` 带 delta 表达式时先读取 place 后求值 delta（如 `(let ((x 1)) (flet ((d () (setf x (* 2 x)))) (incf x (d))) x)` 返回 3 而非 4）— 已修复（修改宏展开用 gensym 先绑定 delta 表达式，`(let* ((g delta-expr)) (setf place (+ place g)))`）

106. `handler-case` 条件类型匹配大小写不敏感（`findClass` 和 `classHasAncestorRec` 使用严格字符串比较，但 reader 全大写化符号名，导致 `(handler-case (error "test") (error (c) ...))` 无法匹配）— 已修复（`findClass` 添加 `strings.ToUpper` 回退，`classHasAncestorRec` 改用 `strings.EqualFold` 比较）

107. `coerce` 的 `list` 类型不支持复数（如 `(coerce #c(3 4) 'list)` 应返回 `(3 4)`）— 已修复（添加 VComplex 分支）

108. `pi`、`most-positive-fixnum`、`most-negative-fixnum` 等 CL 标准常量未定义 — 已修复（在 Go 初始化代码中添加 `pi`、`most-positive-fixnum`、`most-negative-fixnum` 常量）

## 新发现并修复的 Bug（第二轮测试）

109. `listp` 函数未实现 — 已修复（添加 `builtinListP` Go 函数，注册到 builtin table）

110. `coerce` 到 `float`/`single-float`/`double-float` 返回整数而非浮点数（`(coerce 1 'float)` 返回 `1` 而非 `1.0`，`vnum` 存储后 `toString` 按整数打印）— 已修复（添加 `isFloat` 标志到 Value 结构体，`vfloat()` 创建标记为浮点的 VNum 值，coerce float 分支改用 `vfloat`，toString 对 `isFloat` 值强制打印小数点）

111. 词法分析器不区分整数和浮点数字面量（`1` 和 `1.0` 在内部表示完全相同）— 已修复（Tok 结构体添加 `isFlt` 标志，lexNum 解析到小数点或指数标记时设置 `isFlt=true`，reader 对 `isFlt` token 使用 `vfloat`）

112. 算术运算（`+`, `-`, `*`, `/`）结果丢失浮点标记（`(+ 1 2.0)` 返回 `3` 而非 `3.0`）— 已修复（在 builtinAdd/builtinSub/builtinMul/builtinDiv 的最终 float 返回路径添加 `isFloat` 检查，添加 `numOrFloat` 辅助函数自动传播浮点标记）

113. `floatp` 无法识别 `coerce` 产生的浮点数（`1.0 == math.Trunc(1.0)` 导致判断为整数）— 已修复（使用 `isFloat` 标志而非值是否等于其整数部分来判断）

114. `typep`/`subtypep` 的 INTEGER/FLOAT 类型区分不尊重浮点标记（`(typep 1.0 'integer)` 返回 `#t`）— 已修复（在 conditions.go 的 `typepCheckRec` 和符号类型分支中，INTEGER 检查添加 `!val.isFloat` 条件，FLOAT 检查改为 `val.isFloat`）

115. `type-of` 对所有 VNum 返回 "number" 而非区分整数和浮点数 — 已修复（typeStr 对 `isFloat` VNum 返回 "single-float"，对非浮点整数返回 "integer"）

116. `reverse`/`nreverse` 对点状列表（dotted list）丢失尾部（`(nreverse '(1 2 3 . 4))` 返回 `(3 2 1)` 而非 `(3 2 1 . 4)`）— 已修复（改用 `for i := 0; i < len(elems); i++ { tail = cons(elems[i], tail) }` 方式重建反转列表，保留原始尾部值）

## 新发现并修复的 Bug（第三轮 SBCL 测试）

117. `char-name` 对 `(code-char 127)` 返回 "Del" 而非 "Rubout"（ANSI CL 要求 code-char(127) 与 #\rubout 为同一字符，名为 "Rubout"）— 已修复（builtinCharName 添加 code 127 优先返回 "Rubout"）

118. 复数浮点显示丢失 `.0` 后缀（`(coerce 1.0 '(complex float))` 打印 `#c(1 0)` 而非 `#c(1.0 0.0)`，虚部 0.0 显示为 0）— 已修复（在 TComplex 解析器中提前计算 hasFloat 标志，对简化后的 VNum 结果也设置 isFloat=true，使 toString 正确打印小数点）

119. `coerce` 到 `(complex rational)` 类型不产生复数（`(coerce 1/2 '(complex rational))` 返回 `1/2` 而非 `#c(1/2 0)`）— 已修复（coerce 中 compound complex 类型使用 vcomplexAlways）

120. `subtypep` 返回单值而非双值（CL 要求 `(values subtypedefinite-p)`，microlisp 返回裸 `#t`）— 已确认已修复（代码已返回 multiVal 双值）

121. `format ~s` 打印符号大写（`foo` 打印为 `FOO`，CL 默认应保持读取时的大小写）— 非Bug（CL reader 默认大写化符号名，~s 打印大写是正确行为）

122. `map` 不支持 `'list` 结果类型（`(map 'list #'1+ '(1 2 3))` 报错 "unsupported result-type: LIST"）— 已修复（builtinMap switch 使用大写符号名匹配 "LIST"/"CONS"/"VECTOR"/"STRING"）

## 新发现并修复的 Bug（第三轮 SBCL 测试 — backq/hash/random/setf/list）

123. `macroexpand` 对 backquote 形式返回求值结果而非代码形式（`(macroexpand '`#(() a #(#() nil x) #()))` 返回 `#(...)` 而非 `'(quote #(...))`）— 已修复（builtinMacroexpand 的 quasiquote 分支返回 `(list (vsym "quote") expanded)` 而非直接返回 expanded）

124. `sxhash` 对列表返回相同哈希值（`(sxhash '(1 2 3))` 和 `(sxhash '(3 2 1))` 返回相同值，违反哈希质量不变性）— 已修复（sxhashSeen VPair 分支改用 `h = golden ^ car_hash; h *= 31; h += cdr_hash; h ^= h >> 33` 混合公式，确保元素位置影响哈希）

## 新发现但未修复的 Bug（第三轮 SBCL 测试）

125. `read-from-string` 返回 `(value . position)` cons 对，但 `eval` 双重求值时 `((quasiquote ...))` 被当作函数调用导致 "not a procedure: pair" 错误（双重反引号 `(eval (eval (read-from-string expr)))` 测试失败）— 已修复（与 Bug #126 同根因：eval 对 car 为 quasiquote 的列表形式增加特殊处理）

126. `eval` 对 `(quasiquote ...)` 列表求值时，当 `quasiquote` 符号出现在嵌套列表首元素时被当作函数调用（`(eval '((quasiquote (unquote (*RR* *SS*)))))` 报错 "not a procedure: pair" 而非识别为 backquote 形式）— 已修复（在 eval 的 VPair 分支中，检测 car 为 `(quasiquote ...)` 或 `(backquote ...)` 时，调用 evalQuasiquote 展开并返回结果，而非尝试函数调用）

## 第五轮测试发现的 Bug（advanced_numerics/advanced_clos/advanced_packages/list-pure）

131. `typeStr` 返回小写类型名称（如 "rational", "complex", "integer"）而非 ANSI CL 要求的大写名称（如 "RATIONAL", "COMPLEX", "INTEGER"）— 已修复（将 typeStr 中所有返回值改为大写）

132. `#C(n 0)` 复数字面量在虚部为 0 时未简化为实数（ANSI CL 要求 `#C(5 0)` => `5`）— 已修复（TComplex 解析器分支从 `vcomplexAlways` 改为 `vcomplex`）

133. `coerce` 到 `'complex` 类型时对零虚部值返回简化后的实数而非复数（`(coerce 0.5 'complex)` 返回 `0.5` 而非 `#c(0.5 0.0)`）— 已修复（`coerce` 中 plain `'complex` 默认分支改用 `vcomplexAlways`）

## 新发现但未修复的 Bug

127. `random` 对某些值报错 "limit must be >= 1" — 已修复（重构 builtinRandom：VBigInt 使用 big.Int.Rand 避免浮点精度溢出；VRat 截断为整数；VNum 区分浮点和整数路径；添加 vbigInt 辅助函数）

128. `handler-case` 无法捕获 Go 层返回的错误 — 已修复（在 handler-case 评估 valForm 后，若存在 Go 错误，将其转换为 simple-error 条件）

129. `#*` 位向量字面量语法未实现 — 已修复（lexer 添加 `#*` 解析，生成 TVector token，parser 返回 bitVec，sxhashSeen 添加 VArray 分支）

130. `#.` (sharp-dot) 读者宏未实现 — 已修复（lexer 添加 TSharpDot token 类型，Parser.readExpr 中读取下一个表达式并立即 eval）

### 未修复的 Bug（子代理声称修复但实际丢失）

118. 复数浮点显示丢失 `.0` 后缀（`#c(1.0 0.0)` 被 vcomplex 简化为 1 而非 #c(1.0 0.0)，formatComplexPart 代码存在但被 reader 简化绕过）— 已修复（在 TComplex 解析器中提前计算 hasFloat 标志，对简化后的 VNum 结果也设置 isFloat=true）

137. `make-condition` 不评估 `:initform` — 已修复（重写 builtinMakeCondition，遍历 CPL 的 :initform 值并 eval，支持 :initarg 到 slot 名映射，条件类定义改为带 :initform/:initarg 的完整规格）

138. `princ-to-string` 对条件实例返回 `"#<instance ...>"` 而非格式化消息 — 已修复（在 princToString 的 VInstance 分支检测 condition 类祖先，读取 format-control/format-arguments 槽并用 formatMessage 格式化输出）

139. `with-condition-restarts` 宏未实现 — 已修复（添加 defmacro 定义，使用 list/quote 构建 unwind-protect 形式，添加 %associate-restarts-with-condition 和 %dissociate-restarts-with-condition Go 存根函数）

140. `type-error-datum`/`type-error-expected-type` 条件访问器未实现 — 已修复（添加 type-error-datum、type-error-expected-type、stream-error-stream、file-error-pathname、arithmetic-error-operation、arithmetic-error-operands、package-error-package 七个 defun 访问器函数）

## 新修复的 Bug（本次会话）

141. `format ~f` 对整数不添加 `.0`（`(format nil "~f" 42)` 返回 "42" 而非 "42.0"）— 已修复（在 ~F 格式分支中，检查输出是否包含小数点，不包含则追加 ".0"）

142. `format ~c` 对字符打印 `#\a` 而非 `a`（`(format nil "~c" #\a)` 返回 "#\\a"）— 已修复（~C 分支对 VChar 类型直接 WriteRune 输出字符本身，~@C 使用 prin1 转义形式，~:C 拼写非打印字符名称）

143. `format ~%` 和 `~&` 不接受重复计数参数（`(format nil "~3%")` 只产生 1 个换行）— 已修复（从 params 读取重复计数，默认值为 1）

144. `prog`/`prog*` 宏未实现 — 已修复（添加 defmacro 定义，展开为 `(block nil (let/let* bindings (tagbody body)))`，变量无初始化时默认为 nil）

145. `set-dispatch-macro-character` 注册的 dispatch 函数未被 parser 调用 — 已修复（添加 TSharpMacro token 类型，lexer 检查 currentReadtable.dispFns 并匹配用户注册的字符，parser 的 invokeDispatchMacro 调用 (function stream sub-char)）

## 新修复的 Bug（本次会话 — 第四轮）

146. `ignore-errors` 成功时未返回 `(values result nil)` — 已修复（成功路径返回 multiVal(result, vnil())，错误路径返回 multiVal(vnil(), condition)）

147. `nstring-upcase/downcase/capitalize` 对 VArray 输入返回字符串而非修改后的数组 — 已修复（VArray 分支返回 v 而非 vstr(vecToString(v))）

148. `delete-duplicates`/`delete-if`/`delete-if-not` 不支持向量输入 — 已修复（对 VArray/VStr 输入委托给 remove-duplicates/remove-if/remove-if-not）

149. `gethash` 未返回 `(values value present-p)` 双值 — 已修复（找到时返回 multiVal(value, vbool(true))，未找到时返回 multiVal(default, vnil())）

150. `array-has-fill-pointer-p`/`adjustable-array-p`/`array-displacement` 未实现 — 已修复（添加 builtinArrayHasFillPointerP、builtinAdjustableArrayP、builtinArrayDisplacement）

## 新修复的 Bug（本次会话 — 第五轮）

151. CLOS method combinations 未实现 — 已修复（VGeneric 添加 methodCombo 字段，实现 standard/progn/and/or/list/append/nconc/min/max/+ 方法组合分发，set-method-combination 内置函数，defgeneric 支持 :method-combination 选项，isTypeSpecializerMatch/methodSpecificity 使用大写类型名）

152. `with-standard-io-syntax` 宏未实现 — 已修复（Go 初始化添加 *print-radix*/*print-readably*/*print-gensym*/*print-array*/*read-base*/*read-default-float-format*/*read-eval*/*read-suppress* 变量，defmacro 用 unwind-protect save/restore 所有 IO 变量）

153. `with-hash-table-iterator` 宏未实现 — 已修复（实现为 let+macrolet 包装，内部绑定 hash-table-keys 和 index 计数器，迭代器宏返回 (values key value present-p) 三值）

## 新修复的 Bug（本次会话 — 第六轮）

154. `make-array` 对向量 `:initial-contents` 返回 `#(NIL NIL NIL)`（`arrayFillRecursive` 仅遍历 VPair 列表，不处理 VArray 向量）— 已修复（添加 VArray 分支遍历 arr.elements，支持单维和多维情况）

155. CLOS 方法特化优先级不正确（`integer` 方法和 `number` 方法同时匹配整数参数时按定义顺序选择而非类型特化性）— 已修复（`methodSpecificity` 为所有内置类型添加分层分数：INTEGER=130 < RATIONAL=140 < REAL=150 < NUMBER=200，更具体的类型分数更低；添加 `typeIsSubtype` 全局函数用于默认分支的子类型判断）

156. `pairlis` 不接受可选第三参数 `alist`（`(pairlis '(a b) '(1 2) '((c . 3)))` 应返回 `((A . 1) (B . 2) (C . 3))`）— 已修复（添加 `len(args) >= 3` 时取 `alist` 参数，nconc 结果到 alist）

157. `nsubstitute-if` 不支持字符串输入（对字符串调用报错或返回原串）— 已修复（添加 VStr 分支，遍历 rune 切片，谓词匹配时替换字符，支持 :start/:end/:count/:from-end 关键字参数）

158. `format ~e` 科学记数法指数有前导零（`(format nil "~e" 42.0)` 返回 "4.2E+01" 而非 "4.2E+1"）— 已修复（重写 ~E 分支：去除指数前导零，确保尾数有小数点，支持 e 参数控制指数最小位数）

159. 缺少 ANSI CL 标准谓词函数：`bit-vector-p`、`simple-vector-p`、`simple-bit-vector-p`、`simple-string-p`、`packagep`、`compiled-function-p` — 已修复（添加对应 Go 实现，注册到 builtin table）

160. 缺少 ANSI CL 标准函数：`char-int`、`special-operator-p`、`parse-namestring`、`gentemp` — 已修复（添加 Go 实现，special-operator-p 使用硬编码的特殊操作符集合，gentemp 使用 internSymbol 驻留符号）

161. 缺少 ANSI CL 数组函数：`array-dimension`、`array-in-bounds-p`、`array-row-major-index`、`bit`、`sbit` — 已修复（添加 Go 实现，bit/sbit 检查输入为位向量）

## 新修复的 Bug（本次会话）

162. `defmacro` 的 `&optional`/`&key` 参数默认值不生效 — 已修复（重写 parseMacroParams 返回 optDefaults 和 keySpecs，更新 defmacro/define-macro/macrolet 存储这些字段，重写 expandMacro 在参数缺失时评估默认值表达式，支持关键字参数匹配）

163. `destructuring-bind` 中 `&rest` 在 `&optional` 之后不被处理（`&optional` 内部循环吞噬了 `&rest` 关键字）— 已修复（在 &optional 循环中添加 lambda-list 关键字检测，遇到 &rest/&key/&aux/&allow-other-keys 时 break；将循环后的 return nil 替换为 continue 让外层循环处理后续关键字）

164. `float` 函数不设置 isFloat 标志（`(float 1)` 返回 `1` 而非 `1.0`）— 已修复（改用 vfloat 创建返回值）

165. `format ~g` 精度计算错误（使用固定小数位数而非有效位数）— 已修复（根据指数值计算固定格式的小数位数：decPlaces = digits - 1 - exp）

166. `upgraded-array-element-type` 标准函数缺失 — 已修复（添加 builtinUpgradedArrayElementType 内置函数）

167. `makeSimpleCondition` 槽名大小写不匹配（`instSlots` 使用小写键但 `instanceSlotWithInheritance` 用大写查找，导致 `princToString` 对 error 条件返回 `"#<instance ...>"` 而非格式化消息）— 已修复（统一使用大写键名 "MESSAGE"、"FORMAT-CONTROL"、"FORMAT-ARGUMENTS"）

168. `hash-table-size` 返回条目计数而非桶数量（ANSI CL 要求 `hash-table-size` 返回容量/bucket数，`hash-table-count` 返回条目数）— 已修复（改用 `len(ht.hashTab.table)` 返回桶数量）

169. `hash-table-rehash-threshold` 标准函数缺失 — 已修复（HashTable 结构体添加 rehashThreshold 字段，添加 builtinHashTableRehashThreshold，make-hash-table 支持 :rehash-threshold 关键字参数）

170. 缺少 ANSI CL 标准函数：`subst-if-not`、`nsubst-if-not` — 已修复（添加 builtinSubstIfNot 和 builtinNsubstIfNot，谓词返回 false 时替换）

171. 缺少 ANSI CL 标准流函数：`open-stream-p`、`stream-element-type`、`read-char-no-hang`、`read-preserving-whitespace`、`set-syntax-from-char`、`file-string-length` — 已修复（添加 Go 实现到 streams.go）

172. 缺少 ANSI CL 环境函数：`lisp-implementation-type`、`lisp-implementation-version`、`machine-type`、`machine-version`、`machine-instance`、`software-type`、`software-version`、`short-site-name`、`long-site-name` — 已修复（添加 Go 实现）

173. 缺少 ANSI CL 时间函数：`decode-universal-time`、`encode-universal-time` — 已修复（decode 返回9个多值：秒分时日月年星期夏令时时区，encode 接受6+参数返回通用时间）

174. `invoke-restart` 未传递 restart ID（`restartInvoke` 结构体的 id 字段始终为0，导致 `restart-case` 无法精确匹配重启入口）— 已修复（builtinInvokeRestart 在 panic 时传递 `r.id`）

175. 缺少 ANSI CL 标准流变量：`*standard-input*`、`*error-output*`、`*query-io*`、`*terminal-io*` 未绑定到全局环境 — 已修复（initStdStreams 中添加 Set 调用绑定这些变量）

176. `restart-case` 中 `conditionError` 传播时过早截断 restartStack，导致外层 handler-case 无法通过 invoke-restart 激活重启 — 注：这是一个深层架构问题，需要重大重构 restart-case 的 defer 机制

## 新修复的 Bug（本次会话 — 第七轮）

177. `copy-structure` 标准函数缺失 — 已修复（添加 builtinCopyStructure Go 函数，对 VInstance 浅拷贝 instSlots map）

178. `user-homedir-pathname` 标准函数缺失 — 已修复（添加 builtinUserHomedirPathname，使用 os.UserHomeDir 获取主目录路径并构建 pathname）

179. `ldb-test` 标准函数缺失 — 已修复（添加 builtinLdbTest，返回指定位域是否非零）

180. `mask-field` 标准函数缺失 — 已修复（添加 builtinMaskField，返回指定位域的值，语义等同于 ldb）

181. `deposit-field` 标准函数缺失 — 已修复（添加 builtinDepositField，将新字节值插入指定位域，语义等同于 dpb 但参数顺序不同）

182. `asin`/`acos` 标准函数缺失 — 已修复（添加 builtinAsin 和 builtinAcos，使用 math.Asin/math.Acos 实现，对超出 [-1,1] 范围的参数返回复数）

183. `rationalize` 标准函数缺失 — 已修复（添加 builtinRationalize，使用 big.Rat.SetFloat64 将浮点数转为有理数）

184. `internal-time-units-per-second` 标准函数/常量缺失 — 已修复（在 initLib 中定义为返回 1000000 的函数）

185. `get-decoded-time` 标准函数缺失 — 已修复（在 initLib 中定义为调用 decode-universal-time 的便捷函数）

186. `with-open-stream` 标准宏缺失 — 已修复（在 initLib 中定义为 unwind-protect + close 包装）

187. `with-input-from-string` 标准宏缺失 — 已修复（在 initLib 中定义为 make-string-input-stream + unwind-protect 包装）

188. `with-output-to-string` 标准宏缺失 — 已修复（在 initLib 中定义为 make-string-output-stream + get-output-stream-string 包装）

189. `pprint`/`pp` 标准函数缺失 — 已修复（在 initLib 中定义为绑定 *print-pretty* 为 true 后调用 princ + terpri）

190. `pprint-newline`/`pprint-dispatch`/`pprint-tab`/`pprint-logical-block`/`pprint-pop`/`pprint-exit-if-list-exhausted` 标准函数缺失 — 已修复（添加为 stub 实现，pprint-logical-block 为宏）

191. `with-package-iterator` 标准宏缺失 — 已修复（添加为简化 stub 宏）

192. `make-echo-stream` 标准函数缺失 — 已修复（添加 builtinMakeEchoStream，LispStream 添加 isEcho 字段，readChar 中 echo 模式将读取字符写入输出流）

193. `echo-stream-p`/`string-stream-p`/`echo-stream-input-stream`/`echo-stream-output-stream` 标准函数缺失 — 已修复（添加 Go 实现到 streams.go）

194. `unuse-package` 标准函数缺失 — 已修复（添加 builtinUnusePackage，从包的 used 列表中移除指定包）

195. `room` 标准函数缺失 — 已修复（添加 builtinRoom，使用 runtime.ReadMemStats 报告内存使用情况）

196. `make-load-form`/`make-load-form-saving-slots` 标准函数缺失 — 已修复（添加 Go 基础实现）

197. ANSI CL 浮点常量缺失（`most-positive-single-float`、`most-positive-double-float`、`least-positive-single-float` 等） — 已修复（在 initGlobalEnv 中添加 26 个浮点常量）

198. ANSI CL `boole-xxx` 常量缺失（`boole-clr`、`boole-and`、`boole-ior` 等 16 个） — 已修复（在 initGlobalEnv 中添加，值 0-15 对应 builtinBoole 中的 switch 分支）

## 新修复的 Bug（本次会话 — 第八轮）

199. `reduce` 的 `:from-end` 关键字参数不生效（`seqParseKeys` 返回了 `fromEnd` 值但 `builtinSeqReduce` 用 `_` 忽略了它，总是从左到右归约）— 已修复（`builtinSeqReduce` 根据 `fromEnd` 分支：从左到右时 `(fn accumulator element)`，从右到左时 `(fn element accumulator)`，并反转遍历方向）

200. `maphash` 不对第一个参数进行函数类型检查（对空哈希表静默返回 nil，对非空表报通用 "not a function" 而非 type-error）— 已修复（`builtinMaphash` 在循环前验证第一个参数是 VPrim/VFunc/VGeneric，否则报错 "first argument must be a function"）

201. `read-from-string` 不支持 `:start`/`:end` 关键字参数且返回的位置不正确（`builtinReadFromString` 完全忽略关键字参数，`parseExpr` 不返回位置信息）— 已修复（添加 `parseExprWithPos` 返回 `(value position)`，Lexer 添加 `prevEndPos` 字段在 `next()` 开头记录位置，`builtinReadFromString` 解析 `:start`/`:end`/`:eof-error-p`/`:eof-value` 关键字参数并返回 `(values obj original-pos)`）

202. `merge` 对 `'string` 结果类型返回字符列表而非字符串（Lisp 包装器将结果类型传给 `seq-merge` 但被忽略，Go 函数总是返回 list）— 已修复（Lisp 包装器用 `coerce` 将 `seq-merge` 的结果转换为请求的类型）
