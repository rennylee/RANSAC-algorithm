/*

    Name: Yu-Chen Lee
    ID: 300240688


*/
import java.lang.Math;

public class Plane3D {

    private double a;
    private double b;
    private double c;
    private double d;


    public Plane3D(Point3D p1, Point3D p2, Point3D p3) {

        double a1 = p2.getX() - p1.getX();
        double b1 = p2.getY() - p1.getY();
        double c1 = p2.getZ() - p1.getZ();

        double a2 = p3.getX() - p1.getX();
        double b2 = p3.getY() - p2.getY();
        double c2 = p3.getZ() - p1.getZ();

        this.a = b1 * c2 - b2 * c1;
        this.b  = a2 * c1 - a1 * c2;
        this.c= a1 * b2 - b1 * a2;
        this.d = (- a * p1.getX() - b * p1.getY() - c * p1.getZ());
        //System.out.println(A + "" + B + "" + +C + "" + D);
    }

    public Plane3D(double a, double b, double c, double d) {
        this.a=a;
        this.b=b;
        this.c=c;
        this.d=d;
    }


    public double getDistance(Point3D pt) {


        //System.out.println(a + "" + B + "" + +C + "" + D);
        d = Math.abs((a* pt.getX() + b* pt.getY() + c* pt.getZ() + d));

        double vector = Math.sqrt(a*a +a*a + a*a);

        return d/vector;
    }

}
