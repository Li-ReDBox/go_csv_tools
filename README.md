# csv tools

The aims of this tool is to process a csv file by sorting, filtering etc.. The results can be saved to a csv file.
It requires the first line of a csv file is the names for its fields.

## Features of Processor
1. Processor can sort content by a field or multiple fields. Processor can sort string and integer fields.
1. Processor can save to another file with processed data and handle errors.
1. Processor can extract a subset of data with fewer fields.
1. Processor can remove duplicates or rows satisfy one or multiple row conditions.
1. Processor can split rows into a slice of new Processors which each has certain columns are identical. Sorting is a pre-request.
1. Extract, Convert, Split, Write, Filter, Unique create views of original data, so change on views will change the original data.
1. Clone creates a new Processor which has the same content but independent to the source.


## TODO: new features
1. Do some content mainipulations based on some conditions: e.g. Sell (S), add - in front of Quantity in place without addinga new column.
1. Do some content mainipulations based on some conditions: e.g. Buy (B), add - in front of Net Proceeds ($) in place without addinga new column.
1. number operations: this needs to convert string columns to numerical columns, save the results to a new column or replace a column.
1. The previous feature can be applied to multiple columns as long as the conditions are the same: e.g. 