;; ============================================================
;; MicroLisp Adapted SBCL Tests: array.pure.lisp
;; ============================================================

;; --- Basic array creation and access ---
(assert (= (length (make-array 10)) 10) "make-array size")
(assert (= (array-rank (make-array '(3 4))) 2) "array-rank 2d")
(assert (= (array-dimension (make-array '(3 4 5)) 0) 3) "array-dimension 0")
(assert (= (array-dimension (make-array '(3 4 5)) 1) 4) "array-dimension 1")
(assert (= (array-dimension (make-array '(3 4 5)) 2) 5) "array-dimension 2")
(assert (= (array-total-size (make-array '(3 4))) 12) "array-total-size")

;; --- Vector fill pointer ---
(let ((v (make-array 5 :fill-pointer 3 :initial-contents '(1 2 3 4 5))))
  (assert (= (length v) 3) "fill-pointer length")
  (assert (= (aref v 2) 3) "fill-pointer aref")
  (vector-push 99 v)
  (assert (= (length v) 4) "after vector-push length")
  (assert (= (aref v 3) 99) "after vector-push aref"))

;; --- Array type checking ---
(assert (typep (make-array 5) 'array) "typep array")
(assert (typep (make-array 5) 'vector) "typep vector")
(assert (typep (make-array '(2 3)) 'array) "typep 2d-array")

;; --- Bit Vector ---
(let ((bv (make-array 8 :element-type 'bit :initial-contents '(0 0 0 0 0 0 0 0))))
  (assert (= (bit bv 0) 0) "bit 0")
  (assert (= (bit bv 3) 0) "bit 3"))

;; --- Array element type ---
(assert (equal (array-element-type (make-array 5 :element-type 'character)) 'character) "array-element-type character")

;; --- Adjustable arrays ---
(let ((a (make-array 5 :initial-contents '(1 2 3 4 5))))
  (let ((b (adjust-array a 7 :fill-pointer t)))
    (assert (= (length b) 7) "adjust-array extended length")
    (assert (= (aref b 4) 5) "adjust-array preserved elements")))

;; --- Adjustable array with displaced ---
(let ((a (make-array 10 :initial-contents '(0 1 2 3 4 5 6 7 8 9))))
  (let ((b (adjust-array a 5)))
    (assert (= (length b) 5) "adjust-array shrink length")
    (assert (= (aref b 4) 4) "adjust-array shrink element")))