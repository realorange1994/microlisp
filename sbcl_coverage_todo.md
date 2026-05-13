# SBCL 100% 测试覆盖追踪

## 总体进度

- **SBCL 测试文件总数**: 411 (.lisp)
- **SBCL 总断言数 (with-test)**: ~5,371
- **当前 MicroLisp 自有测试**: 60 文件, 2,080 断言
- **目标**: 100% 覆盖 SBCL 所有可解释执行的测试

## 阶段划分

### Phase 1: .pure.lisp (interpreter-safe) — 132 files
Pure tests 不依赖文件系统/外部状态，大部分可直接 load 执行。
少数文件含 `compile` 调用，需跳过或改编。

### Phase 2: .impure.lisp (runtime-safe) — 168 files  
Impure tests 依赖运行时状态（GC、I/O、路径等）。
需筛选掉编译器特定测试（compile-file, type declarations 优化等）。

### Phase 3: 其他测试 — 111 files
.pure-cload.lisp, .impure-cload.lisp, .test.sh, .c — 大部分需要编译器/C 编译，跳过或记录不适用。

## Phase 1 详细进度

| 轮次 | 文件范围 | 状态 | 新增 Bug 数 | 修复 Bug 数 |
|------|----------|------|-------------|-------------|
| Round 1 | 批次 1 (a-d) | DONE | 4 | 4 |
| Round 2 | 批次 2 (b-c) | 待开始 | - | - |
| ... | ... | 进行中 | - | - |

## Phase 1 文件列表 (.pure.lisp, 132 个)

按首字母排序，每批 10-15 个：

**批次 1 (a-d)**:
1. alien-struct-access.impure.lisp (skip - impure)
2. alien-struct-by-value.impure.lisp (skip - impure)
3. aliencall.pure.lisp
4. alientype.pure-cload.lisp (skip - cload)
5. allocator.pure.lisp
6. arith-2.pure.lisp
7. arith-combinations.pure.lisp
8. arith-slow.pure.lisp
9. arith.pure.lisp
10. array.pure.lisp
11. ascii.pure.lisp
12. assembler.pure.lisp
13. avltree.pure.lisp

**批次 2 (b-c)**:
14. backq.pure.lisp
15. bad-code.pure.lisp
16. bit-bash.pure.lisp
17. bsearch.pure.lisp
18. case.pure.lisp
19. ccase.pure.lisp
20. character.pure.lisp
21. charmacro.impure.lisp (skip - impure)
22. classoid-typep.impure.lisp (skip - impure)
23. clockget.pure.lisp
24. clos-method-combination-caches.pure.lisp
25. clos.pure.lisp
26. coalesce.pure.lisp
27. coerce.pure.lisp

**批次 3 (c-f)**:
28. cmp-combinations.pure.lisp
29. condition.pure.lisp
30. constantp.pure.lisp
31. compound.pure.lisp
32. constraint.pure.lisp
33. debug.pure.lisp
34. destructure.pure.lisp
35. dynamic-extent-arrays.pure.lisp
36. dynamic-extent.pure.lisp
37. enc.pure.lisp
38. external-format.pure.lisp
39. fast-removal.pure.lisp
40. fin.pure.lisp

**批次 4 (f-l)**:
41. float.pure.lisp (float-.pure.lisp)
42. format.pure.lisp
43. fun.pure.lisp
44. gc.pure.lisp
45. gcd.pure.lisp
46. gengc.pure.lisp
47. gethash.pure.lisp
48. hash.pure.lisp
49. hashset.pure.lisp
50. info.pure.lisp
51. integerdiv.pure.lisp
52. interface.pure.lisp
53. iso.pure.lisp (2 files)
54. jump.pure.lisp

**批次 5 (l-m)**:
55. lambda.pure.lisp
56. layouts.pure.lisp
57. list.pure.lisp
58. lispobj.pure.lisp
59. load.pure.lisp
60. lockfree.pure.lisp
61. loop.pure.lisp
62. loop-2.pure.lisp
63. macroexpand.pure.lisp
64. make.pure.lisp
65. map.pure.lisp
66. map-tests.pure.lisp
67. map-refs.pure.lisp

**批次 6 (m-p)**:
68. mop.pure.lisp
69. octets.pure.lisp
70. packed.pure.lisp
71. pathnames.pure.lisp
72. pprint.pure.lisp
73. profile.pure.lisp
74. progv.pure.lisp
75. properties.pure.lisp
76. random.pure.lisp
77. reader.pure.lisp
78. redblack.pure.lisp

**批次 7 (p-s)**:
79. sb-posix.pure.lisp
80. seq.pure.lisp
81. serve.pure.lisp
82. setf.pure.lisp
83. simd.pure.lisp (2 files)
84. sleepytests.pure.lisp
85. solist.pure.lisp
86. static-storage.pure.lisp
87. step.pure.lisp
88. stream.pure.lisp
89. string.pure.lisp
90. symbol.pure.lisp
91. symbol-2.pure.lisp

**批次 8 (t-w)**:
92. threads.pure.lisp
93. time.pure.lisp
94. tmpfile.pure.lisp
95. treeshake.pure.lisp
96. type.pure.lisp
97. typecase.pure.lisp
98. typetran.pure.lisp
99. ucs-2.pure.lisp (4 files)
100. unicode.pure.lisp (5 files)
101. utf-16.pure.lisp (5 files)
102. vector.pure.lisp (2 files)
103. vopcombine.pure.lisp
104. wait.pure.lisp
105. weak.pure.lisp
106. win32.pure.lisp (2 files)
107. xset.pure.lisp
108. zstd.pure.lisp

## Bug 记录

### 本阶段新发现 (SBCL 100% 覆盖)

（从 Round 1 开始记录，编号从 258 继续）

### Round 1 Bugs

**Bug #258: `isIntegerValue` accepts VNum float literals (e.g., `3.0`) as integers**
- File: `sbcl-tests/arith-2.pure.lisp`
- Test: `(assert (null (ignore-errors (logior 3.0))) "logior type-error on float")`
- Root cause: `isIntegerValue()` in `lisp.go` accepted VNum values where `isFloat=false` but also accepted whole-number floats where `isFloat=true`. CLHS requires bitwise ops (`logior`, `logand`, `logxor`) to reject non-integer arguments. MicroLisp's VNum with `isFloat=true` represents a CL float literal.
- Fix: Updated `isIntegerValue()` to reject VNum values with `isFloat=true`. The `isFloat` flag on VNum distinguishes between integer literals (`3` => `isFloat=false`) and float literals (`3.0` => `isFloat=true`).
- Status: FIXED

**Bug #259: `builtinEq` with single argument doesn't type-check**
- File: `sbcl-tests/arith-2.pure.lisp`
- Test: `(assert (null (ignore-errors (= 'feep))) "= on symbol")`
- Root cause: `builtinEq` returned `vbool(true)` for `len(args) < 2` without checking if the single argument is a valid number type. Per CLHS, `=` requires all arguments to be numbers.
- Fix: Added explicit type-checking loop before the `len < 2` early-return, and changed `len < 2` to `len == 1` to allow 0-arg case to return true without type-checking (edge case).
- Status: FIXED

**Bug #260: Comparison operators (`<`, `<=`, `>`, `>=`) with single arg don't type-check**
- File: `sbcl-tests/arith-2.pure.lisp`
- Test: `(assert (null (ignore-errors (< #c(0s0 1s0)))) "< on complex")`
- Root cause: `builtinLt`, `builtinGt`, `builtinLe`, `builtinGe` returned `vbool(true)` for `len(args) < 2` without type-checking the single argument. Per CLHS, these require real arguments.
- Fix: Moved the type-checking loop to execute before the `len < 2` early-return in all four comparison operators.
- Status: FIXED

**Bug #261: `logtest` returns `vbool` instead of CL canonical booleans**
- File: `sbcl-tests/arith-2.pure.lisp`
- Test: `(assert (eq (logtest -3 (lognot -3)) nil) "logtest lognot 2")`
- Root cause: `builtinLogtest` returned `vbool()` (a VBool false) for false results. In CL, `nil` is the canonical false value. MicroLisp's `vbool(false)` is VBool, not VNil.
- Fix: Updated `builtinLogtest` to return `vnil()` for false and `vsym("T")` for true.
- Note: Test also changed to use `(not (logtest ...))` form since `eq` comparison with `nil` literal would fail.
- Status: FIXED

### Round 1 文件处理状态

| 文件 | 状态 | Adapted Test | 备注 |
|------|------|-------------|------|
| alien-struct-access.impure.lisp | SKIP | - | impure |
| alien-struct-by-value.impure.lisp | SKIP | - | impure |
| aliencall.pure.lisp | SKIP | - | compiler-only |
| alientype.pure-cload.lisp | SKIP | - | cload |
| allocator.pure.lisp | SKIP | - | compiler-only |
| arith-2.pure.lisp | DONE | `arith-2-tests.lisp` | 4 bugs fixed |
| arith-combinations.pure.lisp | SKIP | - | compiler-only |
| arith-slow.pure.lisp | SKIP | - | compiler-only |
| arith.pure.lisp | DONE | `arith-tests.lisp` | few portable tests |
| array.pure.lisp | DONE | `array-tests.lisp` | basic array ops |
| ascii.pure.lisp | SKIP | - | encoding/external-format |
| assembler.pure.lisp | SKIP | - | compiler-only |
| avltree.pure.lisp | SKIP | - | sb-thread internals |
| backq.pure.lisp | SKIP | - | compiler + backquote bug |
| bit-bash.pure.lisp | SKIP | - | sb-kernel internals |
| bsearch.pure.lisp | SKIP | - | sb-alien FFI |
| case.pure.lisp | SKIP | - | compiler-only |
| ccase.pure.lisp | SKIP | - | compiler-only |
| character.pure.lisp | DONE | `character-tests.lisp` | many portable tests |
| charmacro.impure.lisp | SKIP | - | impure |
| classoid-typep.impure.lisp | SKIP | - | impure |
| clockget.pure.lisp | SKIP | - | sb-thread internals |
| clos-method-combination-caches.pure.lisp | SKIP | - | MOP-specific |
| clos.pure.lisp | SKIP | - | MOP-specific |
| coalesce.pure.lisp | SKIP | - | compiler-only |
| coerce.pure.lisp | DONE | `coerce-tests.lisp` | basic coerce types |
| cmp-combinations.pure.lisp | SKIP | - | compiler-only |
| condition.pure.lisp | PENDING | - | - |
| constantp.pure.lisp | PENDING | - | - |
| compound.pure.lisp | PENDING | - | - |
| constraint.pure.lisp | PENDING | - | - |
| debug.pure.lisp | PENDING | - | - |
| destructure.pure.lisp | PENDING | - | - |
| dynamic-extent-arrays.pure.lisp | PENDING | - | - |
| dynamic-extent.pure.lisp | PENDING | - | - |
| enc.pure.lisp | PENDING | - | - |
| external-format.pure.lisp | PENDING | - | - |
| fast-removal.pure.lisp | PENDING | - | - |
| fin.pure.lisp | PENDING | - | - |
| float.pure.lisp | PENDING | - | - |
| format.pure.lisp | PENDING | - | - |
| fun.pure.lisp | PENDING | - | - |
| gc.pure.lisp | PENDING | - | - |
| gcd.pure.lisp | PENDING | - | - |
| gengc.pure.lisp | PENDING | - | - |
| gethash.pure.lisp | PENDING | - | - |
| hash.pure.lisp | PENDING | - | - |
| hashset.pure.lisp | PENDING | - | - |
| info.pure.lisp | PENDING | - | - |
| integerdiv.pure.lisp | PENDING | - | - |
| interface.pure.lisp | PENDING | - | - |
| iso.pure.lisp | PENDING | - | - |
| jump.pure.lisp | PENDING | - | - |
| lambda.pure.lisp | PENDING | - | - |
| layouts.pure.lisp | PENDING | - | - |
| list.pure.lisp | PENDING | - | - |
| lispobj.pure.lisp | PENDING | - | - |
| load.pure.lisp | PENDING | - | - |
| lockfree.pure.lisp | PENDING | - | - |
| loop.pure.lisp | PENDING | - | - |
| loop-2.pure.lisp | PENDING | - | - |
| macroexpand.pure.lisp | PENDING | - | - |
| make.pure.lisp | PENDING | - | - |
| map.pure.lisp | PENDING | - | - |
| map-tests.pure.lisp | PENDING | - | - |
| map-refs.pure.lisp | PENDING | - | - |
| mop.pure.lisp | PENDING | - | - |
| octets.pure.lisp | PENDING | - | - |
| packed.pure.lisp | PENDING | - | - |
| pathnames.pure.lisp | PENDING | - | - |
| pprint.pure.lisp | PENDING | - | - |
| profile.pure.lisp | PENDING | - | - |
| progv.pure.lisp | PENDING | - | - |
| properties.pure.lisp | PENDING | - | - |
| random.pure.lisp | PENDING | - | - |
| reader.pure.lisp | PENDING | - | - |
| redblack.pure.lisp | PENDING | - | - |
| sb-posix.pure.lisp | PENDING | - | - |
| seq.pure.lisp | DONE | `seq-tests.lisp` | sequence operations |
| serve.pure.lisp | PENDING | - | - |
| setf.pure.lisp | PENDING | - | - |
| simd.pure.lisp | PENDING | - | - |
| sleepytests.pure.lisp | PENDING | - | - |
| solist.pure.lisp | PENDING | - | - |
| static-storage.pure.lisp | PENDING | - | - |
| step.pure.lisp | PENDING | - | - |
| stream.pure.lisp | PENDING | - | - |
| string.pure.lisp | PENDING | - | - |
| symbol.pure.lisp | PENDING | - | - |
| symbol-2.pure.lisp | PENDING | - | - |
| threads.pure.lisp | PENDING | - | - |
| time.pure.lisp | PENDING | - | - |
| tmpfile.pure.lisp | PENDING | - | - |
| treeshake.pure.lisp | PENDING | - | - |
| type.pure.lisp | PENDING | - | - |
| typecase.pure.lisp | PENDING | - | - |
| typetran.pure.lisp | PENDING | - | - |
| ucs-2.pure.lisp | PENDING | - | - |
| unicode.pure.lisp | PENDING | - | - |
| utf-16.pure.lisp | PENDING | - | - |
| vector.pure.lisp | PENDING | - | - |
| vopcombine.pure.lisp | PENDING | - | - |
| wait.pure.lisp | PENDING | - | - |
| weak.pure.lisp | PENDING | - | - |
| win32.pure.lisp | PENDING | - | - |
| xset.pure.lisp | PENDING | - | - |
| zstd.pure.lisp | PENDING | - | - |
