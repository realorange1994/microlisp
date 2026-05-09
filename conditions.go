package main

import (
	"fmt"
	"os"
	"strings"
)

// applyAndResolveTailCall applies fn to args and resolves any tailCall errors.
// This is needed because apply returns tailCall for VFunc, which defers (like warn's)
// would run before the tailCall body is actually evaluated.
func applyAndResolveTailCall(fn *Value, args *Value, env *Env) (*Value, error) {
	result, err := apply(fn, args, env)
	if err != nil {
		if tc, ok := err.(*tailCall); ok {
			// Resolve the tailCall by evaluating in the eval loop
			for {
				result, err = eval(tc.form, tc.env)
				if err != nil {
					if tc2, ok := err.(*tailCall); ok {
						tc = tc2
						continue
					}
					return nil, err
				}
				return result, nil
			}
		}
		return nil, err
	}
	return result, nil
}

// If matched, enters debugger (or calls *debugger-hook*) before handlers run.
func checkBreakOnSignals(cond *Value) {
	breakOn, err := globalEnv.Get("*break-on-signals*")
	if err != nil || isNil(breakOn) {
		return
	}
	if typepCheck(cond, breakOn, globalEnv) {
		if hook, e := globalEnv.Get("*debugger-hook*"); e == nil && hook != nil && hook.typ == VFunc {
			apply(hook, list(cond, vnil()), globalEnv)
			return
		}
		// Default debugger: print message and continue
		fmt.Fprintf(os.Stderr, "\n;; DEBUGGER BREAK: condition signaled: %s\n", toString(cond))
	}
}

// -------- Condition System Builtins --------

// conditionMatchesType checks if a condition instance matches a handler type symbol
// by traversing the CLOS class hierarchy.
func conditionMatchesType(cond *Value, typeSymbol string) bool {
	if cond.instClass == nil {
		return typeSymbol == "condition"
	}
	return classHasAncestor(cond.instClass, typeSymbol)
}

// classHasAncestor checks if cls (or its ancestors) has a class with the given name.
func classHasAncestor(cls *Value, name string) bool {
	seen := make(map[*Value]bool)
	return classHasAncestorRec(cls, name, seen)
}

func classHasAncestorRec(cls *Value, name string, seen map[*Value]bool) bool {
	if cls == nil {
		return name == "condition"
	}
	if seen[cls] {
		return false // cycle detected
	}
	seen[cls] = true
	if cls.str == name {
		return true
	}
	if cls.classParents != nil {
		for _, parent := range cls.classParents {
			if classHasAncestorRec(parent, name, seen) {
				return true
			}
		}
	}
	return false
}

// classMatchesCondition checks if a handler type symbol matches a condition.
func classMatchesCondition(typeSym string, cond *Value) bool {
	if cond == nil || cond.typ != VInstance || cond.instClass == nil {
		return typeSym == "condition"
	}
	return classHasAncestor(cond.instClass, typeSym)
}

// findClass looks up a class by name, checking the class registry first.
func findClass(name string) *Value {
	if cls, ok := classRegistry[name]; ok {
		return cls
	}
	return globalEnv.bindings[name]
}

// signalDivisionByZero creates a division-by-zero condition and checks handlers.
// If a handler catches it, panics with handledError. Otherwise returns a fallback error.
func signalDivisionByZero() error {
	cond := makeSimpleCondition("division-by-zero", "division by zero")
	checkBreakOnSignals(cond)
	if len(handlerStack) > 0 {
		for i := len(handlerStack) - 1; i >= 0; i-- {
			h := handlerStack[i]
			if conditionMatchesType(cond, h.typeSymbol) {
				fn := h.handlerFn
				if fn.typ == VPrim {
					result, err := fn.fn([]*Value{cond})
					if err != nil {
						panic(fmt.Errorf("handler-function panicked: %v", err))
					}
					panic(&handledError{condition: cond, result: result})
				} else if fn.typ == VFunc {
					result, err := apply(fn, cons(cond, vnil()), h.env)
					if err != nil {
						panic(fmt.Errorf("handler-function panicked: %v", err))
					}
					panic(&handledError{condition: cond, result: result})
				}
			}
		}
	}
	return fmt.Errorf("division by zero")
}

// makeSimpleCondition is a helper to create a condition instance of the given class.
func makeSimpleCondition(className, msg string) *Value {
	cond := gcv()
	cond.typ = VInstance
	cond.instClass = findClass(className)
	if cond.instClass == nil {
		cond.instClass = findClass("condition")
	}
	cond.instSlots = map[string]*Value{
		"message":          vstr(msg),
		"format-control":   vstr(msg),
		"format-arguments": vnil(),
	}
	return cond
}

// formatMessage applies ~a/~A substitutions to a format string.
func formatMessage(format string, args []*Value) string {
	for _, a := range args {
		format = strings.Replace(format, "~a", toString(primaryValue(a)), 1)
		format = strings.Replace(format, "~A", toString(primaryValue(a)), 1)
	}
	return format
}

func builtinError(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("error: need at least 1 argument")
	}
	datum := args[0]
	var msg string
	if datum.typ == VStr {
		msg = datum.str
	} else if datum.typ == VInstance {
		// If datum is already a condition, use it directly
		cond := datum
		if len(args) > 1 {
			msg = formatMessage(toString(primaryValue(args[1])), args[2:])
			cond.instSlots["message"] = vstr(msg)
		}
		checkBreakOnSignals(cond)
		if len(handlerStack) > 0 {
			for i := len(handlerStack) - 1; i >= 0; i-- {
				h := handlerStack[i]
				if conditionMatchesType(cond, h.typeSymbol) {
					fn := h.handlerFn
					if fn.typ == VPrim {
						result, err := fn.fn([]*Value{cond})
						if err != nil {
							panic(fmt.Errorf("handler-function panicked: %v", err))
						}
						panic(&handledError{condition: cond, result: result})
					} else if fn.typ == VFunc {
						result, err := apply(fn, cons(cond, vnil()), h.env)
						if err != nil {
							panic(fmt.Errorf("handler-function panicked: %v", err))
						}
						panic(&handledError{condition: cond, result: result})
					}
				}
			}
			panic(&handledError{condition: cond, result: nil})
		}
		if slotMsg, ok := cond.instSlots["message"]; ok {
			return nil, fmt.Errorf("error: %s", toString(slotMsg))
		}
		return nil, fmt.Errorf("error: %s", toString(cond))
	}
	// For non-VInstance datum: create simple-error condition
	if datum.typ != VStr {
		msg = toString(datum)
	}
	if len(args) > 1 {
		msg = formatMessage(msg, args[1:])
	}
	cond := makeSimpleCondition("simple-error", msg)

	checkBreakOnSignals(cond)

	// Walk handler stack
	if len(handlerStack) > 0 {
		for i := len(handlerStack) - 1; i >= 0; i-- {
			h := handlerStack[i]
			if conditionMatchesType(cond, h.typeSymbol) {
				fn := h.handlerFn
				if fn.typ == VPrim {
					result, err := fn.fn([]*Value{cond})
					if err != nil {
						panic(fmt.Errorf("handler-function panicked: %v", err))
					}
					panic(&handledError{condition: cond, result: result})
				} else if fn.typ == VFunc {
					result, err := apply(fn, cons(cond, vnil()), h.env)
					if err != nil {
						panic(fmt.Errorf("handler-function panicked: %v", err))
					}
					panic(&handledError{condition: cond, result: result})
				}
			}
		}
		panic(&handledError{condition: cond, result: nil})
	}
	return nil, fmt.Errorf("error: %s", msg)
}

// builtinCError implements (cerror continue-format-control datum &rest arguments)
// CL spec: signals a simple-error condition, establishes a continue restart.
// If no handler catches it, prints the error and returns nil (implicit continue).
func builtinCError(args []*Value) (*Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("cerror: need at least 2 arguments (continue-format-control datum)")
	}
	contMsg := toString(primaryValue(args[0]))
	datum := args[1]
	var errMsg string
	if datum.typ == VStr {
		errMsg = datum.str
		if len(args) > 2 {
			errMsg = formatMessage(errMsg, args[2:])
		}
	} else {
		errMsg = toString(primaryValue(datum))
	}

	cond := makeSimpleCondition("simple-error", errMsg)

	// Establish continue restart (CL spec: cerror establishes a continue restart)
	continueEntry := restartEntry{
		name: "continue",
		handlerFn: &Value{typ: VPrim, fn: func(_ []*Value) (*Value, error) {
			return vnil(), nil
		}},
		env: globalEnv,
	}
	restartStack = append(restartStack, continueEntry)
	defer func() {
		restartStack = restartStack[:len(restartStack)-1]
	}()

	checkBreakOnSignals(cond)

	// Walk handler stack
	if len(handlerStack) > 0 {
		for i := len(handlerStack) - 1; i >= 0; i-- {
			h := handlerStack[i]
			if conditionMatchesType(cond, h.typeSymbol) {
				fn := h.handlerFn
				if fn.typ == VPrim {
					result, err := fn.fn([]*Value{cond})
					if err != nil {
						panic(fmt.Errorf("handler-function panicked: %v", err))
					}
					panic(&handledError{condition: cond, result: result})
				} else if fn.typ == VFunc {
					result, err := apply(fn, cons(cond, vnil()), h.env)
					if err != nil {
						panic(fmt.Errorf("handler-function panicked: %v", err))
					}
					panic(&handledError{condition: cond, result: result})
				}
			}
		}
		panic(&handledError{condition: cond, result: nil})
	}

	// No handler matched — print error with continue message, return nil
	fmt.Fprintf(os.Stderr, "Error: %s\nContinue: %s\n", errMsg, contMsg)
	return vnil(), nil
}

// builtinWarn implements (warn format-string &rest args)
// CL spec: signals a simple-warning condition, establishes muffle-warning restart.
// If no handler handles it (or handler invokes muffle-warning), suppresses printing.
// Unlike error, warn does NOT transfer control via panic — handlers are called
// and if they return normally, warn continues to check the muffled flag.
func builtinWarn(args []*Value) (*Value, error) {
	msg := "Warning"
	if len(args) >= 1 {
		msg = toString(primaryValue(args[0]))
	}
	if len(args) > 1 {
		msg = formatMessage(msg, args[1:])
	}

	cond := makeSimpleCondition("simple-warning", msg)

	// Check if muffle-warning restart already exists on the stack
	// (e.g., from an outer restart-case). If so, don't establish our own —
	// we want invoke-restart to find the outer one so restart-case can
	// evaluate its body and return the result.
	hasOuterMuffle := false
	for i := len(restartStack) - 1; i >= 0; i-- {
		if restartStack[i].name == "muffle-warning" {
			hasOuterMuffle = true
			break
		}
	}

	// Establish muffle-warning restart (CL spec: warn establishes this restart).
	// Only establish if no outer one exists, to avoid shadowing.
	muffled := false
	savedLen := len(restartStack)
	if !hasOuterMuffle {
		restartStack = append(restartStack, restartEntry{
			name: "muffle-warning",
			handlerFn: &Value{typ: VPrim, fn: func(_ []*Value) (*Value, error) {
				muffled = true
				return vnil(), nil
			}},
			env: globalEnv,
		})
		defer func() {
			restartStack = restartStack[:savedLen]
		}()
	}

	checkBreakOnSignals(cond)

	// Walk handler stack (like signal: do NOT panic — just call handlers).
	for i := len(handlerStack) - 1; i >= 0; i-- {
		h := handlerStack[i]
		if conditionMatchesType(cond, h.typeSymbol) {
			fn := h.handlerFn
			if fn.typ == VPrim {
				fn.fn([]*Value{cond})
			} else if fn.typ == VFunc {
				applyAndResolveTailCall(fn, cons(cond, vnil()), h.env)
			}
			break
		}
	}

	// After signaling, try to invoke muffle-warning. If an outer restart-case
	// established it, invoke-restart finds its nil entry and panics with
	// restartInvoke — the panic propagates to restart-case's defer, which
	// evaluates the body and returns the result. If no outer restart-case,
	// invoke-restart finds our VPrim entry, sets muffled=true, and returns.
	if !muffled {
		_, _ = builtinInvokeRestart([]*Value{vsym("muffle-warning")})
	}

	// If not muffled, print warning
	if !muffled {
		fmt.Fprintf(os.Stderr, "WARNING: %s\n", msg)
	}
	return vnil(), nil
}

func builtinSignal(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("signal: need at least 1 argument")
	}
	datum := args[0]
	var cond *Value
	if datum.typ == VInstance {
		cond = datum
	} else if datum.typ == VSym {
		cond = gcv()
		cond.typ = VInstance
		cond.instClass = findClass(datum.str)
		if cond.instClass == nil {
			cond.instClass = findClass("condition")
		}
		cond.instSlots = map[string]*Value{}
		if len(args) > 1 {
			cond.instSlots["message"] = vstr(toString(primaryValue(args[1])))
		}
	} else {
		cond = datum
	}

	checkBreakOnSignals(cond)

	for i := len(handlerStack) - 1; i >= 0; i-- {
		h := handlerStack[i]
		if conditionMatchesType(cond, h.typeSymbol) {
			fn := h.handlerFn
			if fn.typ == VPrim {
				result, err := fn.fn([]*Value{cond})
				if err != nil {
					return nil, err
				}
				return result, nil
			} else if fn.typ == VFunc {
				result, err := apply(fn, cons(cond, vnil()), h.env)
				if err != nil {
					return nil, err
				}
				return result, nil
			}
		}
	}
	return vnil(), nil
}

func builtinInvokeRestart(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("invoke-restart: need restart name")
	}
	pv := primaryValue(args[0])
	if pv == nil || pv.typ != VSym {
		return nil, fmt.Errorf("invoke-restart: restart name must be a symbol")
	}
	name := pv.str
	restArgs := vnil()
	if len(args) > 1 {
		restArgs = list(args[1:]...)
	}
	for i := len(restartStack) - 1; i >= 0; i-- {
		r := restartStack[i]
		if r.name == name {
			if r.handlerFn != nil {
				return apply(r.handlerFn, restArgs, r.env)
			}
			panic(&restartInvoke{name: name, args: restArgs})
		}
	}
	return nil, fmt.Errorf("invoke-restart: no restart named %s", name)
}

func builtinComputeRestarts(args []*Value) (*Value, error) {
	var result *Value = vnil()
	for i := len(restartStack) - 1; i >= 0; i-- {
		r := restartStack[i]
		result = cons(vsym(r.name), result)
	}
	return result, nil
}

func builtinFindRestart(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("find-restart: need restart name")
	}
	pv := primaryValue(args[0])
	if pv == nil || pv.typ != VSym {
		return nil, fmt.Errorf("find-restart: restart name must be a symbol")
	}
	name := pv.str
	for i := len(restartStack) - 1; i >= 0; i-- {
		if restartStack[i].name == name {
			return vsym(name), nil
		}
	}
	return vnil(), nil
}

func builtinMakeCondition(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("make-condition: need type")
	}
	typeVal := primaryValue(args[0])
	typeName := ""
	if typeVal.typ == VSym {
		typeName = typeVal.str
	} else if typeVal.typ == VInstance && typeVal.instClass != nil {
		typeName = typeVal.instClass.str
	}
	cond := gcv()
	cond.typ = VInstance
	cond.instClass = findClass(typeName)
	if cond.instClass == nil {
		cond.instClass = findClass("condition")
	}
	cond.instSlots = map[string]*Value{}
	for i := 1; i+1 < len(args); i += 2 {
		key := primaryValue(args[i])
		if key.typ == VSym && len(key.str) > 0 && key.str[0] == ':' {
			cond.instSlots[key.str[1:]] = primaryValue(args[i+1])
		} else if key.typ == VSym {
			cond.instSlots[key.str] = primaryValue(args[i+1])
		}
	}
	return cond, nil
}

// -------- typep --------

func builtinTypep(args []*Value) (*Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("typep: need value and type-specifier")
	}
	return vbool(typepCheck(args[0], args[1], globalEnv)), nil
}

// -------- Symbol property lists --------

func builtinCopySymbol(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("copy-symbol: need a symbol")
	}
	sym := args[0]
	if sym.typ != VSym {
		return nil, fmt.Errorf("copy-symbol: expected a symbol")
	}
	copy := gcv()
	copy.typ = VSym
	copy.str = sym.str
	// Copy plist if second arg is non-nil (CL: copy-symbol returns a new uninterned symbol)
	if len(args) >= 2 && !isNil(args[1]) {
		copy.plist = sym.plist
	}
	return copy, nil
}

func builtinGet(args []*Value) (*Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("get: need symbol and indicator")
	}
	sym := args[0]
	indicator := args[1]
	if sym.typ != VSym {
		return nil, fmt.Errorf("get: expected a symbol")
	}
	// Walk plist: (indicator1 value1 indicator2 value2 ...)
	plist := sym.plist
	seen := make(map[*Value]bool)
	for !isNil(plist) && plist.typ == VPair && !isNil(plist.cdr) && plist.cdr.typ == VPair {
		if seen[plist] {
			break
		} // circular plist
		seen[plist] = true
		if eqVal(plist.car, indicator) {
			return plist.cdr.car, nil
		}
		plist = plist.cdr.cdr
	}
	return vnil(), nil
}

func builtinPutprop(args []*Value) (*Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("putprop: need symbol, value, indicator")
	}
	sym := args[0]
	value := args[1]
	indicator := args[2]
	if sym.typ != VSym {
		return nil, fmt.Errorf("putprop: expected a symbol")
	}
	// Walk plist: if indicator exists, update value; else append
	plist := sym.plist
	ppSeen := make(map[*Value]bool)
	for !isNil(plist) && plist.typ == VPair && !isNil(plist.cdr) && plist.cdr.typ == VPair {
		if ppSeen[plist] {
			break
		}
		ppSeen[plist] = true
		if eqVal(plist.car, indicator) {
			plist.cdr.car = value
			return value, nil
		}
		plist = plist.cdr.cdr
	}
	// Not found, append (indicator value) to plist
	newEntry := cons(indicator, cons(value, vnil()))
	if sym.plist == nil || isNil(sym.plist) {
		sym.plist = newEntry
	} else {
		// Append to end of plist (iterative, with cycle detection)
		appendToList := sym.plist
		appendToSeen := make(map[*Value]bool)
		for !isNil(appendToList.cdr) && appendToList.cdr.typ == VPair {
			if appendToSeen[appendToList] {
				break
			}
			appendToSeen[appendToList] = true
			appendToList = appendToList.cdr
		}
		appendToList.cdr = newEntry
	}
	return value, nil
}

func builtinRemprop(args []*Value) (*Value, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("remprop: need symbol and indicator")
	}
	sym := args[0]
	indicator := args[1]
	if sym.typ != VSym {
		return nil, fmt.Errorf("remprop: expected a symbol")
	}
	plist := sym.plist
	var prev *Value = nil
	remSeen := make(map[*Value]bool)
	for !isNil(plist) && plist.typ == VPair && !isNil(plist.cdr) && plist.cdr.typ == VPair {
		if remSeen[plist] {
			break
		}
		remSeen[plist] = true
		if eqVal(plist.car, indicator) {
			// Remove this indicator+value pair
			if prev == nil {
				sym.plist = plist.cdr.cdr
			} else {
				prev.cdr = plist.cdr.cdr
			}
			return vsym("t"), nil
		}
		prev = plist
		plist = plist.cdr.cdr
	}
	return vnil(), nil
}

// -------- get-setf for (setf (get s 'foo) val) --------
func builtinGetSetf(args []*Value) (*Value, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("get-setf: need newval, symbol, indicator")
	}
	newVal := args[0]
	sym := args[1]
	indicator := args[2]
	_, err := builtinPutprop([]*Value{sym, newVal, indicator})
	if err != nil {
		return nil, err
	}
	return newVal, nil
}

func typepCheck(val *Value, typeSpec *Value, env *Env) bool {
	return typepCheckRec(val, typeSpec, env, make(map[*Value]bool))
}

func typepCheckRec(val *Value, typeSpec *Value, env *Env, seen map[*Value]bool) bool {
	if seen[typeSpec] {
		return false
	}
	if isNil(typeSpec) || typeSpec.typ != VSym {
		if typeSpec.typ == VPair && typeSpec.car != nil && typeSpec.car.typ == VSym {
			seen[typeSpec] = true
			switch strings.ToUpper(typeSpec.car.str) {
			case "AND":
				// (and type1 type2 ...) - all must match
				body := typeSpec.cdr
				for !isNil(body) {
					if !typepCheckRec(val, body.car, env, seen) {
						return false
					}
					body = body.cdr
				}
				return true
			case "OR":
				// (or type1 type2 ...) - any must match
				body := typeSpec.cdr
				for !isNil(body) {
					if typepCheckRec(val, body.car, env, seen) {
						return true
					}
					body = body.cdr
				}
				return false
			case "NOT":
				// (not type) - must NOT match
				if typeSpec.cdr == nil || isNil(typeSpec.cdr) || typeSpec.cdr.typ != VPair {
					return false
				}
				return !typepCheckRec(val, typeSpec.cdr.car, env, seen)
			case "SATISFIES":
				// (satisfies fn) - predicate must return true
				body := typeSpec.cdr
				if !isNil(body) && body.typ == VPair && body.car != nil && body.car.typ == VSym {
					fnName := body.car.str
					fn, err := env.Get(fnName)
					if err == nil {
						if fn.typ == VPrim {
							result, err := fn.fn([]*Value{val})
							return err == nil && !isNil(result)
						} else if fn.typ == VFunc {
							result, err := apply(fn, cons(val, vnil()), env)
							return err == nil && !isNil(result)
						}
					}
				}
				return false
			case "EQL":
				// (eql value) - must be eql to value
				body := typeSpec.cdr
				if !isNil(body) {
					return eqlCheck(val, body.car)
				}
				return false
			case "MEMBER":
				// (member v1 v2 ...) - must be eql to one of the values
				body := typeSpec.cdr
				for !isNil(body) {
					if eqlCheck(val, body.car) {
						return true
					}
					body = body.cdr
				}
				return false
			case "ARRAY":
				// (array element-type) - check if array
				if val.typ != VArray {
					return false
				}
				return true
			case "VECTOR":
				// (vector) - check if 1D array
				// (vector element-type) - check if 1D array
				if val.typ != VArray {
					return false
				}
				return len(val.array.dims) == 1
			case "CONS":
				// (cons) - check if it's a pair
				// (cons car-type) - check car matches car-type
				// (cons car-type cdr-type) - check car and cdr match their types
				if val.typ != VPair {
					return false
				}
				// Check cdr (second element if present)
				if typeSpec.cdr != nil && !isNil(typeSpec.cdr) {
					if typeSpec.cdr.typ == VPair && typeSpec.cdr.car != nil {
						carType := typeSpec.cdr.car
						cdrRemaining := typeSpec.cdr.cdr
						// Check car against car-type
						if !typepCheckRec(val.car, carType, env, seen) {
							return false
						}
						// If there's a cdr-type, check val.cdr against it
						if cdrRemaining != nil && !isNil(cdrRemaining) {
							if cdrRemaining.typ == VPair && cdrRemaining.car != nil {
								if !typepCheckRec(val.cdr, cdrRemaining.car, env, seen) {
									return false
								}
							}
						}
					}
				}
				return true
			case "INTEGER":
				return (val.typ == VNum && val.num == float64(int64(val.num))) || val.typ == VRat || val.typ == VBigInt
			case "FLOAT":
				return val.typ == VNum && val.num != float64(int64(val.num))
			case "NUMBER":
				return val.typ == VNum || val.typ == VRat || val.typ == VComplex || val.typ == VBigInt
			case "REAL":
				return val.typ == VNum || val.typ == VRat || val.typ == VBigInt
			case "STRING":
				return val.typ == VStr
			case "SYMBOL":
				return val.typ == VSym || val.typ == VNil
			case "LIST":
				return val.typ == VPair || val.typ == VNil
			case "FUNCTION":
				return val.typ == VPrim || val.typ == VFunc || val.typ == VGeneric
			default:
				// Try as a class name
				if cls := findClass(typeSpec.car.str); cls != nil && cls.typ == VClass {
					if val.typ == VInstance && val.instClass != nil {
						return classHasAncestor(val.instClass, cls.str)
					}
				}
				return false
			}
		}
		return false
	}
	// Symbol type specifier
	typeName := strings.ToUpper(typeSpec.str)
	if typeName == "T" {
		return true
	}
	if typeName == "NULL" {
		return isNil(val)
	}
	if typeName == "NIL" {
		return false // nil is the empty type - never matches any value
	}
	if typeName == "INTEGER" {
		return (val.typ == VNum && val.num == float64(int64(val.num))) || val.typ == VRat || val.typ == VBigInt
	}
	if typeName == "FLOAT" {
		return val.typ == VNum && val.num != float64(int64(val.num))
	}
	if typeName == "NUMBER" {
		return val.typ == VNum || val.typ == VRat || val.typ == VComplex || val.typ == VBigInt
	}
	if typeName == "REAL" {
		return val.typ == VNum || val.typ == VRat || val.typ == VBigInt
	}
	if typeName == "RATIONAL" {
		return val.typ == VRat || val.typ == VNum || val.typ == VBigInt
	}
	if typeName == "COMPLEX" {
		return val.typ == VComplex
	}
	if typeName == "STRING" {
		return val.typ == VStr
	}
	if typeName == "SYMBOL" {
		return val.typ == VSym || val.typ == VNil
	}
	if typeName == "LIST" {
		return val.typ == VPair || val.typ == VNil
	}
	if typeName == "CONS" || typeName == "PAIR" {
		return val.typ == VPair
	}
	if typeName == "FUNCTION" {
		return val.typ == VPrim || val.typ == VFunc || val.typ == VGeneric
	}
	if typeName == "HASH-TABLE" {
		return val.typ == VVHash
	}
	if typeName == "CHARACTER" {
		return val.typ == VChar
	}
	if typeName == "BASE-CHAR" {
		return val.typ == VChar
	}
	if typeName == "STANDARD-CHAR" {
		if val.typ != VChar {
			return false
		}
		ch := val.ch
		// Standard chars: graphic chars in range 32-126 plus #\Newline
		return (ch >= 32 && ch <= 126) || ch == '\n'
	}
	if typeName == "EXTENDED-CHAR" {
		// In non-Unicode implementations, there are no extended chars
		return false
	}
	if typeName == "STREAM" {
		return val.typ == VStream
	}
	if typeName == "ARRAY" {
		return val.typ == VArray
	}
	if typeName == "VECTOR" {
		return val.typ == VArray && len(val.array.dims) == 1
	}
	if typeName == "PROCEDURE" {
		return val.typ == VPrim || val.typ == VFunc
	}
	if typeName == "MACRO" {
		return val.typ == VMacro
	}
	if typeName == "CLASS" {
		return val.typ == VClass
	}
	if typeName == "BOOLEAN" {
		return val.typ == VBool || isNil(val)
	}
	if typeName == "SEQUENCE" {
		return val.typ == VStr || val.typ == VPair || val.typ == VNil || val.typ == VArray
	}
	if typeName == "ATOM" {
		return val.typ != VPair
	}
	// Try as a class name
	if val.typ == VInstance && val.instClass != nil {
		return classHasAncestor(val.instClass, typeName)
	}
	return false
}

func builtinAbort(args []*Value) (*Value, error) {
	return builtinInvokeRestart(append([]*Value{vsym("abort")}, args...))
}

func builtinContinue(args []*Value) (*Value, error) {
	return builtinInvokeRestart(append([]*Value{vsym("continue")}, args...))
}

func builtinMuffleWarning(args []*Value) (*Value, error) {
	return builtinInvokeRestart(append([]*Value{vsym("muffle-warning")}, args...))
}

func builtinStoreValue(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("store-value: need a value")
	}
	return builtinInvokeRestart(append([]*Value{vsym("store-value")}, args...))
}

func builtinUseValue(args []*Value) (*Value, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("use-value: need a value")
	}
	return builtinInvokeRestart(append([]*Value{vsym("use-value")}, args...))
}

func eqlCheck(a, b *Value) bool {
	a = primaryValue(a)
	b = primaryValue(b)
	if a.typ != b.typ {
		return false
	}
	switch a.typ {
	case VNum:
		return a.num == b.num
	case VStr:
		return a.str == b.str
	case VSym:
		return a == b
	case VChar:
		return a.ch == b.ch
	case VBool:
		return a == b
	case VNil:
		return true
	default:
		return a == b
	}
}

// -------- Debugging Builtins --------

// builtinBreak implements (break format-string &rest args)
// CL spec: signals a break condition (class `break`, subclass of serious-condition),
// establishes a continue restart, and enters the debugger if not handled.
func builtinBreak(args []*Value) (*Value, error) {
	msg := "BREAK"
	if len(args) >= 1 {
		if args[0].typ == VStr {
			msg = args[0].str
			if len(args) > 1 {
				msg = formatMessage(msg, args[1:])
			}
		} else {
			msg = toString(primaryValue(args[0]))
		}
	}

	// Create break condition (class `break`, subclass of serious-condition)
	cond := gcv()
	cond.typ = VInstance
	cond.instClass = findClass("break")
	if cond.instClass == nil {
		cond.instClass = findClass("serious-condition")
	}
	cond.instSlots = map[string]*Value{
		"message":          vstr(msg),
		"format-control":   vstr(msg),
		"format-arguments": vnil(),
	}

	// Establish continue restart (CL spec: break establishes continue restart)
	continueEntry := restartEntry{
		name: "continue",
		handlerFn: &Value{typ: VPrim, fn: func(_ []*Value) (*Value, error) {
			return vnil(), nil
		}},
		env: globalEnv,
	}
	restartStack = append(restartStack, continueEntry)
	defer func() {
		restartStack = restartStack[:len(restartStack)-1]
	}()

	// Walk handler stack
	if len(handlerStack) > 0 {
		for i := len(handlerStack) - 1; i >= 0; i-- {
			h := handlerStack[i]
			if conditionMatchesType(cond, h.typeSymbol) {
				fn := h.handlerFn
				if fn.typ == VPrim {
					return fn.fn([]*Value{cond})
				} else if fn.typ == VFunc {
					return apply(fn, cons(cond, vnil()), h.env)
				}
			}
		}
	}

	// Check *debugger-hook*
	if hook, err := globalEnv.Get("*debugger-hook*"); err == nil && hook != nil && hook.typ == VFunc {
		result, _ := apply(hook, list(cond, vnil()), globalEnv)
		return result, nil
	}

	// Default debugger behavior: print message and return nil
	fmt.Fprintf(os.Stderr, "\n;; BREAK: %s\n", msg)
	return vnil(), nil
}

// goErrorToCondition converts a Go error to a Lisp condition object.
// For file-related errors (containing "load: open"), creates a file-error condition.
// Otherwise creates a simple-error condition.
func goErrorToCondition(err error) *Value {
	msg := err.Error()
	cond := gcv()
	cond.typ = VInstance

	// Try to create a file-error for file-related errors
	if strings.Contains(msg, "load: open") || strings.Contains(msg, "open ") {
		cond.instClass = findClass("file-error")
		if cond.instClass != nil {
			// Try to extract filename from error message
			// Format is typically: "load: open <filename>: <error>"
			parts := strings.SplitN(msg, "open ", 2)
			if len(parts) >= 2 {
				fnameParts := strings.SplitN(parts[1], ":", 2)
				fname := strings.TrimSpace(fnameParts[0])
				fname = strings.Trim(fname, "\"")
				cond.instSlots = map[string]*Value{
					"file-pathname":  vstr(fname),
					"message":        vstr(msg),
					"format-control": vstr(msg),
					"format-arguments": vnil(),
				}
				return cond
			}
		}
	}

	// Default: simple-error condition
	cond.instClass = findClass("simple-error")
	if cond.instClass == nil {
		cond.instClass = findClass("error")
		if cond.instClass == nil {
			cond.instClass = findClass("condition")
		}
	}
	cond.instSlots = map[string]*Value{
		"message":          vstr(msg),
		"format-control":   vstr(msg),
		"format-arguments": vnil(),
	}
	return cond
}
