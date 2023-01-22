package globals

var Secret = []byte("CHANGE ME")

const Userkey = "user"

var AvailableVendors = []string{"postgres", "mssql", "sqlite3", "mysql", "oracle"}
var CreatebleVendors = []string{"postgres", "mysql", "sqlite3"}
var DownloadableVendors = []string{"postgres:latest", "mysql:latest"}
