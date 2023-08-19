package cli

const (
	ArgKeyInvalid   ArgKey = iota
	ArgKeyInput            //Input directories or files
	ArgKeyName             //Name of the xtag
	ArgKeyRecursive        //Recursive flag
	ArgKeyFollowSymlinks
	ArgKeyBackup
)

type ArgKey int
