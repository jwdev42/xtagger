# xtagger
## usage
    xtagger [ OPTIONS ] COMMAND
## commands
#### nonterminals for all commands
    PATHS := PATH [ PATHS ]
    NAMES := name NAME [ and NAMES ]
    NAME is the identifier for a specific tag, must be a printable unicode string.
    PATH is a path to a file or directory.
    OPTIONS refer to command line options.
    SIZE_SPEC :~ ^[1-9][0-9]*(K|M|G|T)?$
### command print
    print { [ CONSTRAINT ] [ records ] [ by NAMES ] | untagged } for PATHS
#### tag-specific nonterminals
    CONSTRAINT := { valid | invalid }
##### valid
Valid prints the xtagger attribute for files that have at least one valid record.
##### invalid
Invalid prints the xtagger attribute for files that have no valid records. Files that have no record at all are not considered invalid.
##### untagged
Untagged prints files that have no records.
### command tag
    tag [ CONSTRAINT ] as NAME [ up to SIZE_SPEC ] for PATHS
Command **tag** tags a file or files in a directory.
#### tag-specific nonterminals
    CONSTRAINT := { untagged | invalid }
##### untagged 
If *untagged* is activated, only files that have no record yet will be tagged. If there is at least one record, valid or invalid, the file will be skipped.
##### invalid
If *invalid* is set, only files that don't have a valid record will be tagged. If a file already has a valid record, it will be skipped.
#### optional total size limit
If *up to SIZE_SPEC* is set after *NAME*, xtagger will only tag files as long as their total size sum is smaller than or equal to the limit set by *SIZE_SPEC*.
### command untag
    xbackup untag CONSTRAINT for PATHS
#### tag-specific nonterminals
    CONSTRAINT := { all | invalid | NAMES [ if invalid ] }
##### all
All removes all records.
##### invalid
Invalid removes all invalid records.
##### tag NAME
Tag NAME removes the record with the given name if it exists. The phrase *if invalid* can optionally be added after the name, then the record will only be removed if it is invalid.
### command invalidate
    invalidate { all | NAMES } for PATHS
Command **invalidate** marks records as invalid if the stored hash does not match the file hash anymore.
### command revalidate
    revalidate { all | NAMES } for PATHS
Command **revalidate** marks invalid records as valid again if the stored hash matches the file hash.
