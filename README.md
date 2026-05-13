# MicroLisp

A Common Lisp-compatible interpreter written in Go.

## Overview

MicroLisp is a Lisp interpreter targeting ANSI CL compatibility, implemented as a standalone Go binary. It features a full reader/evaluator/print loop, CLOS object system, condition system, package system, and comprehensive sequence operations.

## Features

### Core Language
- **Reader**: Full readtable system with `set-macro-character`, `set-dispatch-macro-character`, `readtable-case` (upcase/downcase/preserve/invert), `|escaped symbols|`, `#p""` pathname literals, `#*` bit vectors, `#.` sharp-dot, `#\` character literals
- **Evaluator**: Lexical scoping with dynamic extent, special forms, `block`/`return-from`, `catch`/`throw`, `tagbody`/`go`
- **Quasiquote**: Full nested backquote/splice support with correct depth tracking
- **Multi-value**: `values`, `multiple-value-bind`, `multiple-value-call`, `nth-value`

### CLOS Object System
- `defclass` with `:accessor`, `:initarg`, `:initform`, `:conc-name`, `:print-function`
- `defmethod` with `:before`, `:after`, `:around` method qualifiers
- `defgeneric` with method combinations (`standard`, `progn`, `and`, `or`, `list`, `append`, `nconc`, `min`, `max`, `+`)
- EQL specializers, `find-method`, `remove-method`, `compute-applicable-methods`
- `class-of`, `find-class`, `(setf find-class)`, `(setf class-name)`, `ensure-generic-function`

### Condition System
- `define-condition` with inheritance, `:initform`, `:initarg`, `:reader`
- `handler-case`, `handler-bind`, `restart-case`, `restart-bind`
- `signal`, `error`, `warn`, `cerror`, `invoke-restart`, `find-restart`
- Type-error, stream-error, file-error, arithmetic-error conditions with accessors

### Package System
- `defpackage`, `in-package`, `export`, `unexport`, `import`, `unuse-package`
- `find-package`, `find-symbol`, `intern`, `make-package`, `delete-package`
- `list-all-packages`, `symbol-name`, `symbol-package`, `gentemp`
- `CL-USER` as default package with `COMMON-LISP` and `KEYWORD` packages

### Loop Iteration
- `for x =`, `for x in`, `for x on`, `for x across`, `for x from/to/downto/by`
- `being the hash-keys/of`, `being each present-symbol/external-symbol of`
- `with`, `when`, `unless`, `while`, `until`, `collect`, `sum`, `count`, `append`, `nconc`
- Destructuring in loop variables

### Sequence Operations
- `map`, `mapcar`, `mapc`, `mapcan`, `maplist`, `mapcon`, `map-into`, `mapinto`
- `reduce` with `:from-end`, `:initial-value`
- `subseq`, `replace`, `merge`, `fill`, `search`, `mismatch`, `count`, `count-if`
- `sort`, `stable-sort`
- `find`, `find-if`, `find-if-not`, `position`, `position-if`
- `remove`, `remove-if`, `remove-if-not`, `remove-duplicates`
- `delete`, `delete-if`, `delete-if-not`, `delete-duplicates`
- `substitute`, `substitute-if`, `substitute-if-not`
- `nsubstitute`, `nsubstitute-if`, `nsubstitute-if-not`
- `concatenate`, `make-sequence`, `copy-seq`, `reverse`, `nreverse`, `coerce`

### String & Character
- `string-upcase`, `string-downcase`, `string-capitalize` (and `nstring-*` variants)
- `string=`, `string/=`, `string<`, `string>`, `string<=`, `string>=` (with `:start`/`:end` keywords)
- `char-equal`, `char-not-equal`, `char<`, `char>`, `char<=`, `char>=` (multi-argument)
- `char-name`, `char-code`, `code-char`, `character`, `char-int`
- Character type hierarchy: `character > base-char > standard-char`

### Format System
- `~a`, `~s`, `~d`, `~b`, `~o`, `~x`, `~nR` (arbitrary radix 2-36)
- `~f`, `~e`, `~g` (float formatting with precision control)
- `~c` (character, `~@c` escaped, `~:c` spelled out)
- `~%`, `~&` (with repeat count), `~t` (tabulation with `colnum`/`colinc`)
- `~w` (write), `~?` (recursive formatting), `~@?` (variant)
- `~/name/` (user-defined format functions)
- `~{...~}`, `~(...~)`, `~^` (iteration, conditional, abort)

### Numeric System
- Arbitrary precision integers via `math/big` (`big.Int`)
- Rational numbers via `math/big` (`big.Rat`) — `(/ 3 6)` → `1/2`
- Complex numbers: `#c(real imag)`, arithmetic, `expt` with complex base
- Float functions: `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `sinh`, `cosh`, `tanh`, `asinh`, `acosh`, `atanh`
- `floor`, `ceiling`, `truncate`, `round` (and `ffloor`, `fceiling`, `ftruncate`, `fround`)
- `gcd`, `lcm`, `mod`, `rem`, `abs`, `signum`, `sqrt`, `expt`, `exp`, `log`
- Bit operations: `logand`, `logior`, `logxor`, `lognot`, `logandc1`, `logandc2`, `logorc1`, `logorc2`, `lognand`, `lognor`
- `ldb`, `ldb-test`, `dpb`, `mask-field`, `deposit-field`, `byte`, `byte-size`, `byte-position`
- Float introspection: `decode-float`, `integer-decode-float`, `scale-float`, `float-radix`, `float-digits`, `float-precision`
- Constants: `pi`, `most-positive/negative-fixnum`, `most-positive/negative-single/double-float`, `least-positive/negative-*`
- `boole-xxx` constants (0-15)

### Arrays & Vectors
- `make-array` with `:dimensions`, `:element-type`, `:initial-element`, `:initial-contents`, `:fill-pointer`, `:adjustable`
- `aref`, `adjust-array`, `array-dimension`, `array-dimensions`, `array-element-type`
- `array-in-bounds-p`, `array-row-major-index`, `array-total-size`, `array-has-fill-pointer-p`
- `adjustable-array-p`, `array-displacement`, `bit`, `sbit`
- `#(...)` vector literals, `#*` bit vector literals
- Fill pointer: `fill-pointer`, `(setf fill-pointer)`, `vector-push`, `vector-push-extend`, `vector-pop`

### Hash Tables
- `make-hash-table` with `:test`, `:size`, `:rehash-threshold`
- `gethash`, `(setf gethash)`, `remhash`, `clrhash`, `hash-table-count`
- `hash-table-size`, `hash-table-rehash-threshold`, `hash-table-test`, `hash-table-p`
- `with-hash-table-iterator`, `maphash`

### Pathnames
- `#p""` literal syntax, `make-pathname`, `pathname`, `pathname-host`, `pathname-device`
- `pathname-name`, `pathname-type`, `pathname-version`, `pathname-directory`
- `merge-pathnames`, `translate-pathname`, `translate-logical-pathname`
- `logical-pathname-translations`, `(setf logical-pathname-translations)`
- `user-homedir-pathname`, `parse-namestring`

### Streams & I/O
- `open`, `close`, `read-char`, `unread-char`, `peek-char`, `read-line`
- `write-char`, `write-string`, `write-line`, `princ`, `prin1`, `print`
- `open-stream-p`, `stream-element-type`, `read-char-no-hang`, `file-string-length`
- `make-string-input-stream`, `make-string-output-stream`, `get-output-stream-string`
- `echo-stream`, `two-way-stream`, `concatenated-stream`, `broadcast-stream`
- `with-open-file`, `with-open-stream`, `with-input-from-string`, `with-output-to-string`
- `*standard-input*`, `*standard-output*`, `*error-output*`, `*query-io*`, `*terminal-io*`

### Defstruct
- `defstruct` with `:conc-name`, `:constructor`, `:print-function`, `:copier`
- `copy-structure`, automatically generated accessors and constructors

### Macros & Symbols
- `defmacro`, `macrolet`, `symbol-macrolet`, `define-symbol-macro`
- `macroexpand`, `macroexpand-1`, `define-compiler-macro`, `compiler-macro-function`
- `get-macro-character`, `set-macro-character`, `make-dispatch-macro-character`
- `read`, `read-preserving-whitespace`, `read-from-string`, `read-delimited-list`
- `gensym`, `gentemp`, `gensym-counter`, `symbolp`, `keywordp`, `boundp`, `makunbound`

### Control Flow
- `if`, `when`, `unless`, `cond`, `case`, `ecase`, `typecase`, `etypecase`
- `loop`, `do`, `do*`, `dotimes`, `dolist`, `prog`, `prog*`
- `block`, `return-from`, `return`, `catch`, `throw`, `tagbody`, `go`

### Setf Extensions
- `(setf car)`, `(setf cdr)`, `(setf aref)`, `(setf gethash)`, `(setf fill-pointer)`
- `(setf values)`, `(setf symbol-value)`, `(setf symbol-function)`, `(setf macro-function)`
- `(setf class-name)`, `(setf find-class)`, `(setf logical-pathname-translations)`
- `defsetf` with `&environment` support
- `shiftf`, `rotatef`, `incf`, `decf`

### Environment & Introspection
- `functionp`, `compiled-function-p`, `special-operator-p`, `constantp`
- `variable-information`, `function-information`, `declaration-information`
- `apropos`, `apropos-list`, `describe`, `documentation`
- `lisp-implementation-type/version`, `machine-type/version/instance`
- `software-type/version`, `short/long-site-name`, `room`
- `*posix-argv*`, `internal-time-units-per-second`

### Time
- `decode-universal-time`, `encode-universal-time`, `get-decoded-time`
- `get-universal-time`, `get-internal-real-time`, `get-internal-run-time`

### FFI
- `ffi` special form for calling Go functions via reflection

## Building

```bash
go build -o microlisp .
```

## Usage

```bash
# Interactive REPL
./microlisp

# Run a Lisp file
./microlisp < file.lisp

# Or use load from REPL
(load "file.lisp")
```

## Project Structure

| File | Description |
|------|-------------|
| `lisp.go` | Core interpreter: reader, evaluator, printer, builtins, macros (28k lines) |
| `conditions.go` | Condition system: define-condition, handler-case, restart-case (1.5k lines) |
| `streams.go` | Stream system: file/string/echo/two-way/broadcast streams (2.3k lines) |

## Test Suite

60 test files covering all subsystems:

- **Core**: `core.lisp`, `list*.lisp`, `closures*.lisp`, `macros.lisp`
- **Numbers**: `numbers-enhanced.lisp`, `advanced_numerics.lisp`, `numbers-edge-cases.lisp`
- **Strings/Chars**: `strings.lisp`, `characters*.lisp`, `character-*.lisp`
- **Format**: `format-tests.lisp`, `format-advanced.lisp`
- **CLOS**: `advanced_clos.lisp`
- **Conditions**: `conditions-*.lisp`, `panic-bugs.lisp`
- **Packages**: `advanced_packages.lisp`, `packages-advanced.lisp`
- **Sequences**: `sequences*.lisp`, `sequence-keyword-args.lisp`
- **Arrays/Hash**: `array-edge-cases.lisp`, `hash-table-*.lisp`, `hash_tables*.lisp`
- **Readtable**: `readtable-tests.lisp`
- **Types**: `type-tests.lisp`, `coerce*.lisp`
- **Loop**: `loop-iteration-edge-cases.lisp`
- **Specialized**: `ffi.lisp`, `tco.lisp`, `advanced_gc.lisp`, `advanced_tco.lisp`
- **ANSI**: `ansi_tests.lisp`, `sbcl_derived_tests.lisp`
- **Test framework**: `framework.lisp`

## Bug History

257 bugs have been identified and fixed through systematic testing against SBCL test suites. See `todo.md` for the complete changelog.

## Common Lisp Compatibility

MicroLisp aims for ANSI CL compatibility where practical for an interpreter. Key deviations:
- No compilation to native code (interpreted only)
- No multiprocessing/threading support
- FFI is Go-specific via `ffi` special form
- Some advanced reader features (like `#=` sharpequal) may not be fully implemented
