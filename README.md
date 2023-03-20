# csv tools

The aims of this tool set is to process a csv file with sorting, filtering etc then export to another csv file.

## Features
1. Can read file and handle errors: not found, cannot open, not a csv file.
2. Can read content. Make sure it contains all required fields. Can handle other content and csv related errors.
3. Can sort content of a field. _Test the sorted results and any errors expected if any._
4. Can save to another file with processed data and handle errors. _Test successful and failed saving_
5. Processes includes reorder columns, order rows with one or multiple keys; extract a subset of data with
   fewer columns or rows or both.