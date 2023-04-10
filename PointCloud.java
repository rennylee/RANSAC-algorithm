/*

    Name: Yu-Chen Lee
    ID: 300240688


*/
import java.io.*;
import java.util.*;

public class PointCloud {

    private List<Point3D> pts;
    private Random random;

    PointCloud(String filename) throws FileNotFoundException{
        this.pts = new ArrayList<>();
        this.random = new Random();
        File file = new File(filename);
        Scanner scanner = new Scanner(file);
        scanner.nextLine();
        while (scanner.hasNextLine()) {
            String[] container = scanner.nextLine().split("\\s+");

            Point3D point3D = new Point3D(Double.parseDouble(container[0]), Double.parseDouble(container[1]), Double.parseDouble(container[2]));
            pts.add(point3D);
        }
        scanner.close();
    }

    PointCloud() {
        pts = new ArrayList<>();
    }

    public List<Point3D> getPC() {
        return pts;
    }

    public void setPc(List<Point3D> pc){
        this.pts = pc;
    }

    public void addPoint(Point3D pt) {
        pts.add(pt);
    }

    public Point3D getPt() {
        return pts.get(random.nextInt(pts.size()));
    }

    public int size() {
        return pts.size();
    }

    public void save (String filename) throws IOException {
        BufferedWriter bufferedWriter = new BufferedWriter( new FileWriter(filename, false));
        bufferedWriter.write("x, y, z");
        for(Point3D p: pts) {
            bufferedWriter.append(p.getX() + "," + p.getY() + "," + p.getZ());
            bufferedWriter.newLine();
        }
        bufferedWriter.close();
    }

    public final Iterator<Point3D> iterator() {
        return new PointCloudIterator();
    }


    private class PointCloudIterator implements Iterator<Point3D> {

        private int index;

        @Override
        public boolean hasNext() {
            return index < pts.size();
        }

        @Override
        public Point3D next() {
            if (!hasNext()) {
                throw new NoSuchElementException();
            }
            return pts.get(index++);
        }

        @Override
        public void remove() {
            pts.remove(--index);
        }

    }


}
