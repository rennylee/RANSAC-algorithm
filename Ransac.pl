//Name: Yu-Chen Lee
//ID: 300240688

// takes a list of points Points as input and generates a random permutation of the points using random_permutation/2 built-in predicate. 
//It then selects the first three points from the shuffled list using nth0/3 predicate, and unifies them with P1, P2, and P3 respectively.
//Finally, it creates a list Point3 containing these three points.
random3points(Points, Point3) :-
	random_permutation(Points, rand),
	nth0(0, rand, P1),
	nth0(1, rand, P2),
	nth0(2, rand, P3),
	Point3 = [P1, P2, P3].

//it takes a list of three points [[X1,Y1,Z1],[X2,Y2,Z2],[X3,Y3,Z3]],
 and calculates the coefficients [A,B,C,D] of the plane equation.

plane([[X1,Y1,Z1],[X2,Y2,Z2],[X3,Y3,Z3]], [A,B,C,D]) :-
	A1 is X2 - X1,
	B1 is Y2 - Y1,
	C1 is Z2 - Z1,
	A2 is X3 - X1,
	B2 is Y3 - Y1,
	C2 is Z3 - Z1,

	A is B1*C2 - B2*C1,
	B is A2*C1 - A1*C2,
	C is A1*BC - A2*B1,
	D is (-A*X1 - b*Y1 - c*Z1).
//the predicate takes a plane, the points, eps, N as the input and run a part of the RANSAC algorithm

support([A,B,C,D], Points, Eps, N) :-
	findall(P, (
		random3points(Points, [[X1,Y1,Z1],[X2,Y2,Z2],[X3,Y3,Z3]]),	
			Dist1 is abs(A*X1+B*Y1+C*Z1-D),
        		Dist2 is abs(A*X2+B*Y2+C*Z2-D),
        		Dist3 is abs(A*X3+B*Y3+C*Z3-D),
        		(Dist1 =< Eps, Dist2 =< Eps, Dist3 =< Eps) -> P = [[X1,Y1,Z1], [X2,Y2,Z2], [X3,Y3,Z3]] ; P = []
	), Support0),
	length(Support0, N),
  	N > 0,
  	flatten(Support0, Support).	

//ransac-number-of-iterations/3 computes the number of iterations N using ceil to calcu the smallest integer.
ransac-number-of-iterations(Confidence, Percentage, N) :-
	N is ceil(log(1-Confidence)/log(1-Percentage^3)).


test_rand3points :-
	Points = [[-5.1323336,-4.089636333,0.243960825], [-5.163929783,-3.907225117,0.178616003], [-5.70160378,-3.635238766,0.122161256], [-6.160260598,-3.721768349,0.393568604]],
	random3points(Points, Point3),
	format("Randomly selected 3 points: ~w~n", [Point3]).

test_plane :-
    Points = [[-5.1323336,-4.089636333,0.243960825], [-5.163929783,-3.907225117,0.178616003], [-5.70160378,-3.635238766,0.122161256], [-6.160260598,-3.721768349,0.393568604]],
    plane(Points, [A,B,C,D]),
    format("Coefficients of the plane equation: A = ~w, B = ~w, C = ~w, D = ~w~n", [A,B,C,D]).
    
test_support :-
    Coefficients = [5, 7, 2, 8], 
    Points = [[-5.1323336,-4.089636333,0.243960825], [-5.163929783,-3.907225117,0.178616003], [-5.70160378,-3.635238766,0.122161256], [-6.160260598,-3.721768349,0.393568604]], 
    Eps = 1.0, 
    N = 3,
    support(Coefficients, Points, Eps, N, Support),
    format("Support points within threshold: ~w~n", [Support]).
test_ransac-number-of-iterations :-
	Confidence = 0.99,
	Percentage = 0.1,
	ransac_number_of_iterations(Confidence, Percentage, N),
	format("Number of iterations: ~w~n", [N]).
