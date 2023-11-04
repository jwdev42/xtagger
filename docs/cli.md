# xbackup
## commands
#### nonterminals for all commands
PATHS := PATH [ PATHS ]

NAME is the identifier for the tag, must be a printable string.

PATH is a path to a file or directory.

OPTIONS refer to command line options.
### command print
    print [ OPTIONS ] [ CONSTRAINT ] for PATHS
#### tag-specific nonterminals
    CONSTRAINT := { valid | invalid }
##### valid
Valid only prints files that have at least one valid record.
##### invalid
Invalid only prints files that have no valid records. Files that have no record at all are not considered invalid.
### command tag
    tag [ OPTIONS ] [ CONSTRAINT ] as NAME in PATHS
Command **tag** tags a file or files in a directory.
#### tag-specific nonterminals
    CONSTRAINT := { untagged | invalid }
##### untagged 
If untagged is activated, only files that have no record yet will be tagged. If there is at least one record, valid or invalid, the file will be skipped.
##### invalid
If invalid is set, only files that don't have a valid record will be tagged. If a file already has a valid record, it will be skipped.
### command remove
    xbackup remove [-name name] [-non-interactive=bool] dir...
