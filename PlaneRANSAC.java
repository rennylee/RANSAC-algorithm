/*

    Name: Yu-Chen Lee
    ID: 300240688


*/
import java.io.IOException;
import java.lang.Math;
import java.util.ArrayList;
import java.util.List;

public class PlaneRANSAC {

    private double eps;

    private PointCloud pcs;

    Plane3D dom_plane;

    public PlaneRANSAC (PointCloud pc) {
        this.pcs = pc;
    }

    public void setEps(double eps) {
        this.eps = eps;
    }

    public double getEps() {
        return eps;
    }

    public int getNumberOfIterations (double confidence, double percentageOfPointsOnPlane) {
        return (int) ((Math.log(1-confidence))/(Math.log(1-(Math.pow(percentageOfPointsOnPlane, 3)))));
    }

    public void run(int numberOfIterations, String filename) throws IOException {

        pcs = new PointCloud(filename);
        for (int t=1; t <4; t++ ) {
            int best_sup = 0;

            List<Point3D> container = new ArrayList<>();
            for (int i = 0; i < numberOfIterations; i++) {
                int curr_sup = 0;

                Point3D pt1 = pcs.getPt();
                Point3D pt2 = pcs.getPt();
                Point3D pt3 = pcs.getPt();

                Plane3D curr_plane = new Plane3D(pt1, pt2, pt3);

                for (int j = 0; j < pcs.size(); j++) {
                    if (curr_plane.getDistance(pcs.getPC().get(j)) < eps) {
                        curr_sup++;
                    }
                }
                if (curr_sup > best_sup) {
                    dom_plane = curr_plane;
                    best_sup = curr_sup;
                }
            }
            for (int j = 0; j < pcs.size(); j++) {
                if (dom_plane.getDistance(pcs.getPC().get(j)) < eps) {
                    container.add(pcs.getPC().get(j));
                    pcs.getPC().remove(j);
                }
            }
            PointCloud holder = new PointCloud();
            holder.setPc(container);
            holder.save("newPointCloud" + "_p" + t+ ".xyz");

        }
        pcs.save("newPointCloud_p0.xyz");
    }


    public static void main(String[] args) throws IOException {

        PlaneRANSAC planeRANSAC = new PlaneRANSAC(new PointCloud("PointCloud1.xyz"));
        double confidence = 0.99;
        double percentageOfPointsOnPlane = 0.3;

        int numIterations = planeRANSAC.getNumberOfIterations(confidence, percentageOfPointsOnPlane);

        String filename = "PointCloud1.xyz";
        planeRANSAC.setEps(0.15);
        planeRANSAC.run(numIterations, filename);

    }

}
