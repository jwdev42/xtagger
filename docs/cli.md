# xbackup
## usage
    xtagger [ OPTIONS ] COMMAND
## commands
#### nonterminals for all commands
    PATHS := PATH [ PATHS ]
    NAME is the identifier for a specific tag, must be a printable string.
    PATH is a path to a file or directory.
    OPTIONS refer to command line options.
### command print
    print [ CONSTRAINT ] for PATHS
#### tag-specific nonterminals
    CONSTRAINT := { valid | invalid | untagged }
##### valid
Valid prints the xtagger attribute for files that have at least one valid record.
##### invalid
Invalid prints the xtagger attribute for files that have no valid records. Files that have no record at all are not considered invalid.
##### untagged
Untagged prints files that have no records.
### command tag
    tag [ CONSTRAINT ] as NAME in PATHS
Command **tag** tags a file or files in a directory.
#### tag-specific nonterminals
    CONSTRAINT := { untagged | invalid }
##### untagged 
If untagged is activated, only files that have no record yet will be tagged. If there is at least one record, valid or invalid, the file will be skipped.
##### invalid
If invalid is set, only files that don't have a valid record will be tagged. If a file already has a valid record, it will be skipped.
### command untag
    xbackup untag CONSTRAINT from PATHS
#### tag-specific nonterminals
    CONSTRAINT := { all | invalid | REMOVE_CONSTRAINT }
    REMOVE_CONSTRAINT := tag NAME [ if invalid ]
##### all
All removes all records.
##### invalid
Invalid removes all invalid records.
##### tag NAME
Tag NAME removes the record with the given name if it exists. The phrase *if invalid* can optionally be added after the name, then the record will only be removed if it is invalid.
