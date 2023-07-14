# csv tools

The aims of this tool is to process a csv file by sorting, filtering, etc. The results can be saved to a csv file.
It requires the first line of a csv file is the names for its fields.

## Features of Processor
1. Processor primarily works on csv with titles.
1. Processor can sort content by a field or multiple fields. Processor can sort string and integer fields.
1. Processor can save to another file with processed data and handle errors.
1. Processor can extract a subset of data with fewer fields.
1. Processor can remove duplicates or rows satisfy one or multiple row conditions.
1. Processor can split rows into a slice of new Processors which each has certain columns are identical. Sorting is a pre-request.
1. Extract, Convert, Split, Write, Filter, Unique create views of original data, so change on views will change the original data.
1. Clone creates a new Processor which has the same content but independent to the source.
1. An [`Operation`](row_ops.go#L12) contains an [`Action`](row_ops.go#L9) and a [`IsFunc`](processor.go#L216) which works on a row (a slice of string).
1. For checking rows, [`IsFunc`](processor.go#L216) should works on the whole elements of a row. To qualify all rows, simply return `true`. Because 
   all elements are available to `IsFunc`, it means we can check numbers of element in various ways.
1. Rows can be changed an [`Operation`](row_ops.go#L12).
1. [`IsFunc`](processor.go#L216) and [`Action`](row_ops.go#L9) in [`Operation`](row_ops.go#L12) works on the same slice of string.
1. `Processor.Replace` can have multiple `Operation` which means rows can be changed in sequence.
1. `Processor.Derive` add one more column by deriving new content based on two marked columns (by their positions).
