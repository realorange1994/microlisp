;; ============================================================
;; MicroLisp Adapted SBCL Tests: arith.pure.lisp
;; ============================================================

;; --- Fundamental arithmetic operations ---
(assert (= (+ 4 2) 6) "+ basic")
(assert (= (- 4 2) 2) "- basic")
(assert (= (* 4 2) 8) "* basic")
(assert (= (/ 4 2) 2) "/ basic")
(assert (= (expt 4 2) 16) "expt basic")

;; --- Coerce to complex float ---
(assert (= (coerce 1 '(complex float)) #c(1.0 0.0)) "coerce int to complex float")
(assert (= (coerce 1/2 '(complex float)) #c(0.5 0.0)) "coerce rational to complex float")
(assert (= (coerce #c(1 2) '(complex float)) #c(1.0 2.0)) "coerce complex to complex float")

;; --- Coerce within float bounds ---
(assert (= (coerce 1 '(single-float -1.0 2.0)) 1.0) "coerce within float bounds")