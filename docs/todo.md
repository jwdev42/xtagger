# TODO
## Commands
### print
Print a file's attribute if there is at least one record.
#### Flags
##### --name
Only print a file's attribute if there is a record with the specified name.
### print_valid
Only print a file's attribute if there is at least one valid record.
#### Flags
##### --name
Only print a file's attribute if there is a valid record with the specified name.
### print_invalid
Only print a file's attribute if there is no valid record.
#### Flags
##### --name
Only print a file's attribute if there is an invalid record with the specified name.
### tag
Tag files if they don't have any record for the given name yet. If there already exists a matching record within a file, it will be skipped.
#### Flags
##### --name
The name for the new record(s)
