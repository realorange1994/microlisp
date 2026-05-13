;; ============================================================
;; MicroLisp Adapted SBCL Tests: character.pure.lisp
;; ============================================================

;; --- Named characters ---
(assert (characterp (name-char "Newline")) "name-char Newline")
(assert (characterp (name-char "Space")) "name-char Space")
(assert (characterp (name-char "Tab")) "name-char Tab")
(assert (characterp (name-char "Page")) "name-char Page")
(assert (characterp (name-char "Rubout")) "name-char Rubout")
(assert (characterp (name-char "Return")) "name-char Return")
(assert (characterp (name-char "Backspace")) "name-char Backspace")

;; name-char / code-char consistency for standard characters
(assert (eql (name-char "Newline") (code-char 10)) "name-char Newline = code-char 10")
(assert (eql (name-char "Space") (code-char 32)) "name-char Space = code-char 32")
(assert (eql (name-char "Tab") (code-char 9)) "name-char Tab = code-char 9")
(assert (eql (name-char "Page") (code-char 12)) "name-char Page = code-char 12")
(assert (eql (name-char "Rubout") (code-char 127)) "name-char Rubout = code-char 127")
(assert (eql (name-char "Return") (code-char 13)) "name-char Return = code-char 13")
(assert (eql (name-char "Backspace") (code-char 8)) "name-char Backspace = code-char 8")

;; name-char for non-existent name
(assert (null (name-char 'foo)) "name-char foo")

;; --- Basic character predicates ---
(assert (characterp #\a) "characterp a")
(assert (characterp #\Space) "characterp Space")
(assert (characterp #\Newline) "characterp Newline")

;; --- Character comparison ---
(assert (char= #\a #\a) "char= a a")
(assert (char/= #\a #\b) "char/= a b")
(assert (char< #\a #\b) "char< a b")
(assert (char> #\b #\a) "char> b a")
(assert (char<= #\a #\a) "char<= a a")
(assert (char>= #\b #\a) "char>= b a")

;; --- Case predicates ---
(assert (upper-case-p #\A) "upper-case-p A")
(assert (lower-case-p #\a) "lower-case-p a")
(assert (not (upper-case-p #\a)) "not upper-case-p a")
(assert (not (lower-case-p #\A)) "not lower-case-p A")
(assert (both-case-p #\A) "both-case-p A")
(assert (both-case-p #\a) "both-case-p a")
(assert (not (both-case-p #\1)) "not both-case-p 1")

;; --- Alphanumeric predicate ---
(assert (alphanumericp #\a) "alphanumericp a")
(assert (alphanumericp #\A) "alphanumericp A")
(assert (alphanumericp #\1) "alphanumericp 1")
(assert (not (alphanumericp #\Space)) "not alphanumericp Space")

;; --- Character equality (case-insensitive) ---
(assert (char-equal #\a #\A) "char-equal a A")
(assert (char-equal #\A #\a) "char-equal A a")
(assert (not (char-equal #\a #\b)) "not char-equal a b")

;; --- Graphic character predicate ---
;; Per CLHS, graphic chars are printed chars that are NOT #\Space or #\Newline
(assert (graphic-char-p #\a) "graphic-char-p a")
(assert (not (graphic-char-p #\Space)) "not graphic-char-p Space")
(assert (not (graphic-char-p #\Newline)) "not graphic-char-p Newline")

;; --- Alpha character predicate ---
(assert (alpha-char-p #\a) "alpha-char-p a")
(assert (alpha-char-p #\A) "alpha-char-p A")
(assert (not (alpha-char-p #\1)) "not alpha-char-p 1")

;; --- Digit character predicate ---
(assert (digit-char-p #\5) "digit-char-p 5")
(assert (not (digit-char-p #\a)) "not digit-char-p a")

;; --- code-char / char-code roundtrip ---
(assert (= (char-code (code-char 65)) 65) "char-code code-char 65")
(assert (= (char-code #\A) 65) "char-code A")
(assert (= (char-code #\a) 97) "char-code a")
(assert (= (char-code #\0) 48) "char-code 0")

;; --- char-upcase / char-downcase ---
(assert (char= (char-upcase #\a) #\A) "char-upcase a")
(assert (char= (char-upcase #\A) #\A) "char-upcase A unchanged")
(assert (char= (char-downcase #\A) #\a) "char-downcase A")
(assert (char= (char-downcase #\a) #\a) "char-downcase a unchanged")

;; --- Case-insensitive comparisons (exhaustive for 0-127) ---
(dotimes (i 128)
  (let* ((char (code-char i))
         (down (char-downcase char))
         (up (char-upcase char)))
    (assert (char-equal char char) (format nil "char-equal self ~d" i))
    (when (char/= char down)
      (assert (char-equal char down) (format nil "char-equal down ~d" i)))
    (when (char/= char up)
      (assert (char-equal char up) (format nil "char-equal up ~d" i)))))

;; --- Standard-char predicates across 0-127 ---
(dotimes (i 128)
  (let ((char (code-char i)))
    (when (typep char 'standard-char)
      (if (find char "abcdefghijklmnopqrstuvwxyz")
          (assert (lower-case-p char) (format nil "lower ~c" char))
          (assert (not (lower-case-p char)) (format nil "not lower ~c" char)))
      (if (find char "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
          (assert (upper-case-p char) (format nil "upper ~c" char))
          (assert (not (upper-case-p char)) (format nil "not upper ~c" char)))
      (if (find char "0123456789")
          (assert (digit-char-p char) (format nil "digit ~c" char))
          (assert (not (digit-char-p char)) (format nil "not digit ~c" char))))))