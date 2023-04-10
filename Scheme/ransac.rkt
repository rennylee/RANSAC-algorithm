#lang racket
(define (readXYZ fileIn)
    (let ((sL (map (lambda s (string-split (car s)))
        (cdr (file->lines fileIn)))))
    (map (lambda (L)
        (map (lambda (s)
            (if (eqv? (string->number s) #f)
                s
                (string->number s))) L)) sL)))

(define (random-points Ps)
  (list (list-ref Ps (random (length Ps)))
        (list-ref Ps (random (length Ps)))
        (list-ref Ps (random (length Ps)))))


(define (plane P1 P2 P3)
  (let ((a1 (- (car P2) (car P1)))
        (b1 (- (cadr P2) (cadr P1)))
        (c1 (- (caddr P2) (caddr P1)))
        (a2 (- (car P3) (car P1)))
        (b2 (- (cadr P3) (cadr P1)))
        (c2 (- (caddr P3) (caddr P1))))
    (let ((a (- (* b1 c2) (* b2 c1)))
          (b (- (* a2 c1) (* a1 c2)))
          (c (- (* a1 b2) (* b1 a2)))
          (d (+ (* a1 (+ (car P1) (car P2) (car P3)))
                (* b1 (+ (cadr P1) (cadr P2) (cadr P3)))
                (* c1 (+ (caddr P1) (caddr P2) (caddr P3))))))
      (list a b c d))))

(define (support plane points)
  (let ((count 0))
    (let ((a (car plane))
          (b (cadr plane))
          (c (caddr plane))
          (d (cadddr plane)))
      (do ((i 0 (+ i 1)))
          ((= i (length points)) (cons count plane))
        (let ((point (list-ref points i))
              (x (car points))
              (y (cadr points))
              (z (caddr points)))
          (when (= (+ (* a x) (* b y) (* c z)) d)
            (set! count (+ count 1))))))))

            
(define (dominantPlane Ps k)
  (let loop ((bestSupport 0)
             (bestPlane '())
             (iterations k))
    (if (= iterations 0)
        bestPlane
        (let* ((points (random-points Ps))
               (plane (plane (car points) (cadr points) (caddr points)))
               (supportCount (support plane Ps)))
          (if (> supportCount bestSupport)
              (loop supportCount plane (- iterations 1))
              (loop bestSupport bestPlane (- iterations 1)))))))

(define (ransacNumberOfIteration confidence percentage)
  (let ((p (expt (- 1 confidence) 3))
        (k (/ (log (- 1 (expt percentage 3))) (log p))))
    (ceiling k)))


(define (planeRANSAC filename confidence percentage eps)
  (let ((Ps (readXYZ filename))
    (k (ransacNumberOfIteration confidence percentage)))
  (dominantPlane Ps K)))






         