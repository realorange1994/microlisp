;;;; Adapted from sbcl-tests/condition.pure.lisp
;;;; Portable condition system tests

;;; Test 1: make-condition with file-error is printable
(let ((c (make-condition 'file-error :pathname "foo")))
  (assert (stringp (prin1-to-string c)))
  (assert (stringp (princ-to-string c)))
  (format t "PASS: file-error printable~%"))

;;; Test 2: handler-bind with handler-case basic signal handling
(assert (eq :handled
            (handler-bind ((simple-error (lambda (c)
                                           (declare (ignore c))
                                           (return-from test :handled)))
                           (error (lambda (c)
                                    (declare (ignore c))
                                    (return-from test :error-handled))))
              (signal "test error")
              :not-handled)))

;;; Test 3: handler-bind signal 'condition
(let ((called nil))
  (handler-bind ((condition (lambda (c)
                              (declare (ignore c))
                              (setf called t))))
    (signal 'condition))
  (assert called))

;;; Test 4: handler-bind with function reference
(let ((called nil))
  (flet ((my-handler (c)
           (declare (ignore c))
           (setf called t)))
    (handler-bind ((condition #'my-handler))
      (signal 'condition)))
  (assert called))

;;; Test 5: handler-case with declare
(assert (null (handler-case (error "foo")
                (error ()
                  (declare (optimize speed))))))

;;; Test 6: simple-condition without args
(let ((sc (make-condition 'simple-condition)))
  (assert (not (simple-condition-format-control sc)))
  (assert (not (simple-condition-format-arguments sc)))
  (assert (stringp (prin1-to-string sc))))

;;; Test 7: make-condition with non-condition class should signal type-error
(assert (typep (nth-value 1
                          (ignore-errors
                           (make-condition 'standard-class)))
               'type-error)))

;;; Test 8: make-condition with undefined condition type
(let ((result (ignore-errors (make-condition 'no-such-condition-xyz))))
  (assert (not (car result)))
  (let ((err (cdr result)))
    (assert (typep err 'type-error))
    (assert (stringp (princ-to-string err)))))

;;; Test 9: handler-bind empty bindings
(assert (typep (nth-value 1
                          (ignore-errors (handler-bind () (error "Foo"))))
               'error)))

;;; Test 10: restart-case basic functionality
(assert (eq :ok
            (restart-case
                (handler-bind ((error (lambda (c)
                                        (declare (ignore c))
                                        (invoke-restart 'my-restart :ok))))
                  (error "test"))
              (my-restart (x) x))))

;;; Test 11: restart-case multiple values
(assert (equal '(1 2)
               (multiple-value-list
                (restart-case
                    (handler-bind ((error (lambda (c)
                                            (declare (ignore c))
                                            (invoke-restart 'foo 1 2))))
                      (error "test"))
                  (foo (a b) (values a b)))))))

;;; Test 12: restart-case with no restart args
(assert (eq 'bar
            (restart-case
                (handler-bind ((error (lambda (c)
                                        (declare (ignore c))
                                        (invoke-restart 'foo))))
                  (error "test"))
              (foo () 'bar))))

;;; Test 13: cerror with continue restart
(let ((x 0))
  (handler-bind ((simple-error (lambda (c)
                                 (declare (ignore c))
                                 (incf x)
                                 (continue))))
    (cerror "Continue" "test error ~A" 42))
  (assert (= x 1)))

;;; Test 14: type-error-datum and type-error-expected-type
(let ((result (ignore-errors (funcall (lambda (x) (check-type x fixnum) x) t))))
  (assert (not (car result)))
  (let ((err (cdr result)))
    (assert (typep err 'type-error))
    (assert (eq (type-error-datum err) t))))

;;; Test 15: error with format control and arguments
(let ((result (ignore-errors (error "Value ~A is bad" 42))))
  (assert (not (car result)))
  (assert (typep (cdr result) 'simple-error)))

;;; Test 16: restart-case evaluation order
(let (result)
  (with-condition-restarts (progn
                             (push 1 result)
                             (make-condition 'error))
      (progn (push 2 result) nil)
    (push 3 result))
  (assert (equal result '(3 2 1))))

;;; Test 17: simple-error format control
(let ((c (make-condition 'simple-error :format-control "Hello ~A" :format-arguments '(world))))
  (assert (string= "Hello ~A" (simple-condition-format-control c)))
  (assert (equal '(world) (simple-condition-format-arguments c))))

;;; Test 18: signal warning with muffle-warning
(handler-case
    (handler-bind ((warning #'muffle-warning))
      (signal 'warning))
  (:no-error (&rest args)
    (assert (null args))))

;;; Test 19: malformed handler-bind bindings
(flet ((test (binding)
         (assert (eq :ok
                     (handler-case
                         (macroexpand `(handler-bind (,binding)))
                       (simple-error (e)
                         :ok))))))
  (test 1)
  (test '())
  (test '(error))
  (test '(error #'print :foo)))

;;; Test 20: restart-case malformed clauses
(flet ((test (clause)
         (assert (eq :ok
                     (handler-case
                         (macroexpand `(restart-case (error "foo") ,clause))
                       (simple-error (e)
                         :ok))))))
  (test :report)
  (test ())
  (test (foo)))

(format t "All condition tests passed.~%")
