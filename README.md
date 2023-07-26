# Generate Student Completion Certificates

This program generates completion certificates for students based on a template image and a CSV file containing student names.

All this program does is input an image, and output many images with text written on them, pulled from a CSV. It outputs a new image for each line in the CSV. This is useful for bulk-creating custom certificates of completion to students with their name on it. 

## Installation

To install the program, run the following command:

```
go install github.com/bretfisher/gencert
```

## Usage

To generate a completion certificate, run the following command:

```
gencert -image template.png -names students.csv -output output-path
```

This command will generate a completion certificate using the `template.png` image as the template, the `students.csv` file as the data source, and save the output to `output-path/student-name.png`.

The following options are available:

```
  -center
        Center the text instead of aligning to the left
  -color string
        Font color name (default "black")
  -font string (required)
        Path to the TrueType (ttf) font file
  -height float
        Percentage of image height to start printing the name (default 0.5)
  -image string (required)
        Path to the input image file (png only)
  -names string (required)
        Path to the CSV file containing names, one name per row
  -output string
        Path to the output directory (default ".")
  -size float
        Font size in points (default 75)
  -width float
        Percentage of image width to start printing the name (default 0.5)
```

## Example

Suppose we have the following `students.csv` file:

```
Alice Jones
Bob Smith
Charlie Brown
```

And we've created a certificate template with something like Canva that has a blank area for student names.

We also need a TrueType (ttf) font file. We can download a font file from [Google Fonts](https://fonts.google.com/). For example, we can download the [Roboto](https://fonts.google.com/specimen/Roboto) font by clicking the "Select this font" button and then clicking the "Download" button in the bottom right corner of the page. This will download a zip file containing the font files. We can unzip the file and use the `Roboto-Regular.ttf` file as the font file.

We can generate completion certificates for each student using the following command:

```
gencert -image template.png -names students.csv -font Roboto-Regular.ttf -output output.png -size 100 -width 0.5 -height 0.5 -center
```

This will generate the following `output.png` image:

![Output Image](output.png)

## Required Inputs

The program requires the following inputs:

- A template image file in PNG format (`-image` option)
- A CSV file containing student names (`-names` option)
- A TTF font used to draw the student's name

## Expected Outputs

The program generates a completion certificate for each student in the CSV file, using the template image as the background and the student name as the text. The input and output files need to be PNG.

## Help Output

The program provides help output when the `-h` or `--help` option is used:

## LICENSE

See the [MIT LICENSE](./LICENSE)
