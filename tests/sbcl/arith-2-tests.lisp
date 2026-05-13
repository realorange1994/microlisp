;; ============================================================
;; MicroLisp Adapted SBCL Tests: arith-2.pure.lisp
;; ============================================================

;; --- Min/Max type error tests ---
(assert (null (ignore-errors (min '(1 2 3)))) "min type-error on list")
(assert (= (min -1) -1) "min single arg")
(assert (null (ignore-errors (min 1 #(1 2 3)))) "min type-error on vector")
(assert (= (min 10 11) 10) "min two args")
(assert (= (min 1.5 2) 1.5) "min float with integer")
(assert (= (min 5.0 -3) -3) "min mixed float/int")
(assert (null (ignore-errors (max #c(4 3)))) "max type-error on complex")
(assert (= (max 0) 0) "max single arg")
(assert (= (max -1 10.0) 10.0) "max mixed")
(assert (null (ignore-errors (max 3 "foo"))) "max type-error on string")
(assert (= (max -3 0) 0) "max two args")

;; --- Arithmetic tests ---
(assert (= (+ 3.0) 3.0) "+ single float")
(assert (= (+ 1 2) 3) "+ two args")
(assert (= (+ 3.0 4.0) 7.0) "+ floats")
(assert (= (* 3.0) 3.0) "* single float")
(assert (= (* 1 2) 2) "* two args")
(assert (= (* 3.0 4.0) 12.0) "* floats")

;; --- Logbitwise type error tests ---
(assert (null (ignore-errors (logand #(1)))) "logand type-error on vector")
(assert (= (logand 1) 1) "logand single arg")
(assert (null (ignore-errors (logior 3.0))) "logior type-error on float")
(assert (= (logior 4) 4) "logior single arg")
(assert (null (ignore-errors (logxor #c(2 3)))) "logxor type-error on complex")
(assert (= (logxor -6) -6) "logxor single arg")

;; --- ASH tests ---
(dotimes (i 41)
  (assert (= (ash (1- (ash 1 32)) (- i))
             (if (< i 32)
                 (1- (ash 1 (- 32 i)))
                 0))
          (format nil "ash right shift ~d" i)))

;; --- GCD ---
(assert (= (gcd most-negative-fixnum most-negative-fixnum) (- most-negative-fixnum)) "gcd most-negative")
(assert (= (gcd most-negative-fixnum 48) 16) "gcd negative fixnum with 48")

;; --- Log operations ---
(assert (logtest -3 (lognot 5)) "logtest lognot")
(assert (not (logtest -3 (lognot -3))) "logtest lognot 2")

;; --- Case with or-patterns ---
(assert (= (case 0 ((0 -3) 1) (t 2)) 1) "case or-pattern 0")
(assert (= (case -3 ((0 -3) 1) (t 2)) 1) "case or-pattern -3")
(assert (= (case 3 ((0 -3) 1) (t 2)) 2) "case or-pattern default")
(assert (= (case 1 ((0 -3) 1) (t 2)) 2) "case or-pattern other")
(assert (= (case -1 ((-1 0) 0) (t 1)) 0) "case or-pattern -1")
(assert (= (case 0 ((-1 0) 0) (t 1)) 0) "case or-pattern 0")
(assert (= (case 1 ((-1 0) 0) (t 1)) 1) "case or-pattern default")

;; --- Log operations with negatives ---
(assert (= (logand 7702 -1) 7702) "logand with -1")
(assert (= (logorc2 3 -1) 3) "logorc2 with all ones")
(assert (= (logandc1 -1 0) 0) "logandc1 simplified")

;; --- DPB tests ---
(assert (= (dpb 90 (byte 63 8) 81) 23121) "dpb large byte")
(assert (= (dpb 1 (byte 32 32) 1) 4294967297) "dpb 32-bit field")

;; --- Mask-field tests ---
(assert (= (mask-field (byte 78 0) 35) 35) "mask-field large byte")

;; --- Rem tests ---
(assert (= (rem -2 2) 0) "rem negative even")
(assert (= (rem -3 2) -1) "rem negative odd")
(assert (= (rem 2 2) 0) "rem positive even")
(assert (= (rem 3 2) 1) "rem positive odd")
(assert (= (rem 3 4) 3) "rem smaller divisor")
(assert (= (rem -3 4) -3) "rem negative smaller divisor")

;; --- Logxor/Logior/Logand patterns ---
(assert (= (logxor 5 4) 1) "logxor basic")
(assert (= (logxor -5 -5) 0) "logxor same negative")
(assert (= (logxor 0 0) 0) "logxor zeros")
(assert (= (logior 0 1) 1) "logior basic")
(assert (= (logior -1 0) -1) "logior -1")
(assert (= (logior 5 3) 7) "logior overlapping bits")
(assert (= (logand 7 3) 3) "logand basic")
(assert (= (logand -1 5) 5) "logand -1")
(assert (= (logand 5 0) 0) "logand 0")

;; --- Expt 0 0 ---
(assert (= (expt 0 0) 1) "expt 0 0")

;; --- Isqrt tests ---
(assert (= (isqrt 0) 0) "isqrt 0")
(assert (= (isqrt 1) 1) "isqrt 1")
(assert (= (isqrt 2) 1) "isqrt 2")
(assert (= (isqrt 3) 1) "isqrt 3")
(assert (= (isqrt 4) 2) "isqrt 4")
(assert (= (isqrt 9) 3) "isqrt 9")
(assert (= (isqrt 16) 4) "isqrt 16")
(assert (= (isqrt 100) 10) "isqrt 100")

;; --- ASH edge cases ---
(assert (= (ash 10 -2) 2) "ash right by 2")
(assert (= (ash -7514499718243589878 -2) -1878624929560897470) "ash negative right")

;; --- Complex division ---
(assert (= (/ 0 #c(1.0 3.0)) #c(0.0 0.0)) "complex div zero")
(assert (= (/ 5 #c(0.0 1.0)) #c(0.0 -5.0)) "complex div by i")
(assert (= (/ -2.0 #c(0.0 1.0)) #c(0.0 2.0)) "complex div neg by i")

;; --- Numeric inequality type errors ---
(assert-error (= 'feep) type-error)
(assert-error (< #c(0s0 1s0)) type-error)
(assert-error (<= #c(0s0 1s0)) type-error)
(assert-error (> #c(0s0 1s0)) type-error)
(assert-error (>= #c(0s0 1s0)) type-error)
(assert-error (= 3 'feep) type-error)
(assert-error (< 3 'feep) type-error)
(assert-error (< 0 0 'feep) type-error)
(assert-error (= 0 0 'feep) type-error)

;; --- Min/Max complex type errors ---
(assert-error (min #c(1s0 -2s0)) type-error)
(assert-error (max #c(1s0 -2s0)) type-error)

;; --- GCD 0 x returns (abs x) ---
(dolist (x (list -10 (* 3 most-negative-fixnum)))
  (assert (= (gcd 0 x) (abs x)) "gcd 0 x = abs x"))

;; --- LCM non-negative ---
(assert (= (lcm 4 -10) 20) "lcm 4 -10 = 20")
(assert (= (lcm 0 0) 0) "lcm 0 0 = 0")

;; --- Bignum multiplication ---
(assert (= (* 966082078641 419216044685) 404997107848943140073085) "bignum multiplication")

;; --- ASH negative bignum ---
(assert (= (ash (1- most-negative-fixnum) (1- most-negative-fixnum)) -1) "ash negative bignum")

;; --- LOGCOUNT basic tests ---
(assert (= (logcount 0) 0) "logcount 0")
(assert (= (logcount 1) 1) "logcount 1")
(assert (= (logcount 7) 3) "logcount 7")
(assert (= (logcount 8) 1) "logcount 8")
(assert (= (logcount -1) 0) "logcount -1")
(assert (= (logcount -7) 2) "logcount -7")
(assert (= (logcount -8) 3) "logcount -8")

;; --- SIGNUM ---
(assert (= (signum 5) 1) "signum positive")
(assert (= (signum -5) -1) "signum negative")
(assert (= (signum 0) 0) "signum zero")

;; --- GCD positive for large numbers ---
(assert (plusp (gcd 20286123923750474264166990598656 680564733841876926926749214863536422912)) "gcd positive large")

;; --- Truncate tests ---
(assert (equal (multiple-value-list (truncate 3 4)) '(0 3)) "truncate smaller")
(assert (equal (multiple-value-list (truncate -3 4)) '(0 -3)) "truncate negative")
(assert (equal (multiple-value-list (truncate 4 4)) '(1 0)) "truncate equal")
(assert (equal (multiple-value-list (truncate -4 4)) '(-1 0)) "truncate negative equal")
