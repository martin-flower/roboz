# intersection

## warning

* This algorithm does not scale beyond **200** or so commands. The memory footprint is relatively low, but the comparison between the lines, and specifically needing to track the line number of points with intersections, (particularly lines which overlap in their length) is very compute intensive.

## algorithm

1. convert all commands into lines (start point, end point)
1. calculate the total number of points in all the lines (including duplicates)
1. take the list of lines (in any order),
    * for each line, see which other lines it intersects with (just compare with the lines to the right)
    * if a point is an intersection between two lines, we record the line number of the second line
    * the second line is only counted once
1. the number of unique points is the total number of points in all the lines, minus the total number of points with at least one intersection

It is necessary to know the second line number in the comparison in order to avoid double-counting

Take for example four lines:

* line1 : (0,0 - 5,0)
* line2 : (0,0 - 5,0)
* line3 : (0,0 - 5,0)
* line4 : (0,0 - 5,0)

The total number of points in the lines is 24

Intersections count:
* line 1 : 0
* line 2 : 6 (compare line1/line2)
* line 3 : 6 (compare line1/line3 - do not count line2/line3 intersections)
* line 4 : 6 (compare line1/line4 - do not count line2/line4 or line3/line4 intersections)

Total intersections count **18**

Total unique points in commands: **24** - 18 = 6
