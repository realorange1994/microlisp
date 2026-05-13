;; ============================================================
;; MicroLisp Adapted SBCL Tests: seq.pure.lisp
;; ============================================================

;; --- Remove with start/end/from-end ---
(let* ((orig '(1 2 3 2 6 1 2 4 1 3 2 7))
       (x (copy-seq orig))
       (y (remove 3 x :from-end t :start 1 :end 5))
       (z (remove 2 x :from-end t :start 1 :end 5)))
  (assert (equalp orig x) "remove: orig unchanged")
  (assert (equalp y '(1 2 2 6 1 2 4 1 3 2 7)) "remove from-end start end")
  (assert (equalp z '(1 3 6 1 2 4 1 3 2 7)) "remove from-end start end z"))

;; --- Substitute with start/end/from-end ---
(let* ((orig '(a a a a a a a a a a))
       (y (substitute 'x 'a orig :start 2 :end 7)))
  (assert (equal y '(a a x x x x x a a a)) "substitute with start end"))

;; --- Substitute-if with count and from-end ---
(let* ((orig '(a a a a a a a a a a))
       (y (substitute-if 'x (lambda (x) (eq x 'a)) orig :start 2 :end 7 :count 3)))
  (assert (equal y '(a a x x x a a a a a)) "substitute-if count"))

;; --- Position basic ---
(assert (= (position 3 '(1 2 3 4 5)) 2) "position basic")
(assert (null (position 3 '(1 2 3 4 5) :start 3)) "position not found")
(assert (= (position 3 '(1 2 3 4 5) :from-end t) 2) "position from-end")

;; --- Count ---
(assert (= (count 3 '(1 2 3 4 3 5)) 2) "count basic")

;; --- Find ---
(assert (= (find 3 '(1 2 3 4 5)) 3) "find basic")

;; --- Position with key ---
(assert (= (position 2 '(1 2 3 4 5) :key #'1+) 0) "position with key")

;; --- Concatenate ---
(assert (equal (concatenate 'list '(1 2) '(3 4)) '(1 2 3 4)) "concatenate list")

;; --- Map ---
(assert (equal (map 'list #'1+ '(1 2 3)) '(2 3 4)) "map list")

;; --- Reduce ---
(assert (= (reduce #'+ '(1 2 3 4 5)) 15) "reduce +")
(assert (= (reduce #'+ '(1)) 1) "reduce single")

;; --- Reverse ---
(assert (equal (reverse '(1 2 3)) '(3 2 1)) "reverse list")

;; --- Subseq ---
(assert (equal (subseq '(1 2 3 4 5) 2) '(3 4 5)) "subseq start")
(assert (equal (subseq '(1 2 3 4 5) 1 3) '(2 3)) "subseq start end")

;; --- Sort ---
(let ((x (list 3 1 4 1 5 9 2 6)))
  (let ((y (sort (copy-seq x) #'<)))
    (assert (equal y '(1 1 2 3 4 5 6 9)) "sort <")))

;; --- Merge ---
(assert (equal (merge 'list '(1 3 5) '(2 4 6) #'<) '(1 2 3 4 5 6)) "merge lists")