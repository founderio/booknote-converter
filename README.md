# Booknote Converter

This tool converts my personal notation of transactions to a csv format that can be read by [Homebank](http://homebank.free.fr/en/index.php).  
Written for personal usage. The motivation behind this is that I write down all payments made with cash, i.e. anything I don't have a receipt for, to get a more complete overview of my finances. In the past I have transferred the notes I take manually, but since the notation is easy to understand I have automated this step.

The notation is as simple as this:

```
2017:
20.12.
3,80 Bakery
10 Pizza

30.12.
2,40 Gas Station

2018:
2.1.
5,0 Post Office

```

- Empty lines will be ignored.
- The parser is relatively lenient with the number formats, but the date format DD.MM is a hard requirement. The year HAS to be in the format 'YYYY:'
- The trailing dot for dates as shown above is optional, you can even add ':'
- Commas and dots are equal in function - you use what you feel comfortable with.  
- The year number can also be passed in as a parameter if you only have transactions in the same year.

## Building from source

The supplied makefile simply runs `go build` with a custom output file. On Unix systems, simply run `make` from the source directory.

## Using the tool

Run `booknote-converter -h` to get this overview:

```
Usage of ./booknote-converter:
  -input string
        Input file for noted transactions (default "input.txt")
  -output string
        Output file for generated csv (default "output.csv")
  -year string
        The initial year to use, leave blank if the year is specified as 'yyyy:' in the file
```
