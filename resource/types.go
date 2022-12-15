package resource

const (
	// Company represents the company name string file info entry in the resources table
	Company = "CompanyName"
	// FileDescription represents the file description entry in the resources table
	FileDescription = "FileDescription"
	// FileVersion represents the file version entry in the resources table
	FileVersion = "FileVersion"
	// OriginalFilename the name of the original executable in the resources table
	OriginalFilename = "OriginalFilename"
	// LegalCopyright represents the copyright notice in the resources directory table
	LegalCopyright = "LegalCopyright"
	// ProductName is the product name entry in the resources table
	ProductName = "ProductName"
	// ProductVersion is the product version entry in the resources table
	ProductVersion = "ProductVersion"
)

// PE contains various headers that identifies the format and characteristics of the executable files.
type VersionResources map[string]string
