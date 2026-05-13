;; ============================================================
;; MicroLisp Adapted SBCL Tests: coerce.pure.lisp
;; ============================================================

;; --- Coerce number to float ---
(assert (= (coerce 3 'single-float) 3.0) "coerce int to single-float")
(assert (= (coerce 3 'float) 3.0) "coerce int to float")

;; --- Coerce integer to character ---
(assert (characterp (coerce 65 'character)) "coerce int to character")
(assert (char= (coerce 65 'character) #\A) "coerce 65 to #\A")

;; --- Coerce list to vector ---
(assert (equalp (coerce '(1 2 3) 'vector) #(1 2 3)) "coerce list to vector")

;; --- Coerce vector to list ---
(assert (equal (coerce #(1 2 3) 'list) '(1 2 3)) "coerce vector to list")

;; --- Coerce string to list ---
(assert (equal (coerce "abc" 'list) '(#\a #\b #\c)) "coerce string to list")

;; --- Coerce list to string ---
(assert (string= (coerce '(#\a #\b #\c) 'string) "abc") "coerce list to string")

;; --- Coerce string to vector ---
(assert (equalp (coerce "abc" 'vector) #(#\a #\b #\c)) "coerce string to vector")

;; --- Coerce vector to string ---
(assert (string= (coerce #(#\a #\b #\c) 'string) "abc") "coerce vector to string")

;; --- Coerce float to complex ---
(assert (= (coerce 3.0 'complex) #c(3.0 0.0)) "coerce float to complex")
(assert (= (coerce 3 'complex) #c(3 0)) "coerce int to complex")

;; --- Coerce identity ---
(assert (= (coerce 3 'integer) 3) "coerce int to integer")
(assert (char= (coerce #\a 'character) #\a) "coerce char to character")
(assert (string= (coerce "abc" 'string) "abc") "coerce string to string")

;; --- Coerce complex ---
(assert (= (coerce #c(3 0) 'complex) #c(3 0)) "coerce complex identity")

;; --- Coerce rational to float ---
(assert (= (coerce 1/2 'single-float) 0.5) "coerce rational to single-float")

;; --- Coerce complex to float ---
(assert (= (coerce #c(3.0 0.0) 'single-float) 3.0) "coerce complex real to float")