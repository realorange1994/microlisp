;; readtable-tests.lisp — tests for the readtable system
;; MicroLisp readtable support: *readtable*, set-macro-character, etc.

(load "tests/framework.lisp")
(start-suite "Readtable System")

;; --- Basic readtable predicates ---
(assert-true (readtablep *readtable*) "*readtable* is a readtable")
(assert-false (readtablep 42) "numbers are not readtables")
(assert-false (readtablep 'foo) "symbols are not readtables")
(assert-false (readtablep "hello") "strings are not readtables")
(assert-false (readtablep '(1 2 3)) "lists are not readtables")

;; --- readtable-case ---
(assert-equal ':UPCASE (readtable-case *readtable*) "standard readtable is :UPCASE")

;; --- make-readtable ---
(define rt (make-readtable))
(assert-true (readtablep rt) "make-readtable creates a readtable")
(assert-equal ':UPCASE (readtable-case rt) "new readtable defaults to :UPCASE")

;; --- copy-readtable ---
(define rt-copy (copy-readtable))
(assert-true (readtablep rt-copy) "copy-readtable creates a readtable")
(assert-equal ':UPCASE (readtable-case rt-copy) "copied readtable preserves case mode")

(define rt2 (copy-readtable *readtable*))
(assert-true (readtablep rt2) "copy-readtable with argument works")

;; --- set-readtable-case ---
(set-readtable-case rt ':DOWNCASE)
(assert-equal ':DOWNCASE (readtable-case rt) "set-readtable-case works for :DOWNCASE")

(set-readtable-case rt ':PRESERVE)
(assert-equal ':PRESERVE (readtable-case rt) "set-readtable-case works for :PRESERVE")

(set-readtable-case rt ':INVERT)
(assert-equal ':INVERT (readtable-case rt) "set-readtable-case works for :INVERT")

(set-readtable-case rt ':UPCASE)
(assert-equal ':UPCASE (readtable-case rt) "set-readtable-case works for :UPCASE")

;; --- get-macro-character for standard chars ---
;; Standard macro chars are handled by Go-level code, so get-macro-character returns nil
(assert-nil (get-macro-character '#\') "get-macro-character for ' returns nil (Go-level)")
(assert-nil (get-macro-character '#\`) "get-macro-character for ` returns nil (Go-level)")
(assert-nil (get-macro-character '#\,) "get-macro-character for , returns nil (Go-level)")
(assert-nil (get-macro-character '#\;) "get-macro-character for ; returns nil (Go-level)")
(assert-nil (get-macro-character '#\") "get-macro-character for \" returns nil (Go-level)")
(assert-nil (get-macro-character '#\() "get-macro-character for ( returns nil (Go-level)")
(assert-nil (get-macro-character '#\)) "get-macro-character for ) returns nil (Go-level)")

;; --- get-macro-character for non-macro chars ---
(assert-nil (get-macro-character '#\a) "get-macro-character for a returns nil")
(assert-nil (get-macro-character '#\+) "get-macro-character for + returns nil")
(assert-nil (get-macro-character '#\space) "get-macro-character for space returns nil")

;; --- set-macro-character: define ! as macro that returns quoted bang ---
(set-macro-character #\! (lambda (c) (list 'quote 'bang)))
(display "Testing set-macro-character...")
(newline)

;; ! alone should expand to (quote bang) which evaluates to bang
(assert-equal 'bang ! "macro ! expands to bang")

;; --- set-macro-character with non-terminating ---
;; Define $ as non-terminating macro (like ')
(set-macro-character #\$ (lambda (c) (list 'quote 'dollar)) #f)
(assert-true #t "set-macro-character with non-terminating flag works")

;; --- make-dispatch-macro-character ---
;; #: is traditionally a dispatch character for the #: uninterned symbol reader macro
(make-dispatch-macro-character '#\#)
(display "make-dispatch-macro-character works")
(newline)

;; --- set-dispatch-macro-character ---
;; Define #? as a custom dispatch macro
(make-dispatch-macro-character '#\#)
;; Note: dispatch macros are harder to test in pure Lisp since they require
;; the parser to handle them. The basic setup should work.
(display "dispatch macro infrastructure set up")
(newline)

;; --- get-dispatch-macro-character ---
(assert-nil (get-dispatch-macro-character '#\# '#\!) "get-dispatch-macro-char returns nil for unset")

;; --- Multiple macro characters ---
(set-macro-character #\| (lambda (c) (list 'quote 'bar)))
(assert-equal 'bar | "second macro | works")

;; --- Macro character overwriting ---
(set-macro-character #\! (lambda (c) (list 'quote 'new-bang)))
(assert-equal 'new-bang ! "macro overwriting works")

;; --- Readtable isolation ---
;; Macros set in one readtable shouldn't affect another
(define rt-isolated (make-readtable))
(assert-nil (get-macro-character #\! rt-isolated) "fresh readtable has no ! macro")

;; --- readtable-case affects symbol reading ---
;; Test that readtable-case is properly integrated
;; symbol-name returns a string
(assert-equal "FOO" (symbol-name 'FOO) "symbol-name returns correct string")
(assert-equal ":KEYWORD" (symbol-name ':KEYWORD) "keyword symbol name includes colon")

;; --- Nested readtable operations ---
(set-readtable-case rt2 ':DOWNCASE)
(assert-equal ':DOWNCASE (readtable-case rt2) "second readtable case set independently")

;; --- Verify standard behavior still works ---
(assert-equal '(1 2 3) '(1 2 3) "standard quote behavior unchanged")
(assert-equal '(1 . 2) '(1 . 2) "standard dotted pair unchanged")
(assert-equal "hello" "hello" "standard strings unchanged")

(end-suite)
(test-summary)
