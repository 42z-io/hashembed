package hashembed

// Config holds the configuration for the [HashedFS].
type Config struct {
	Hasher            FileHasher  // mechanism used to hash the file
	Renamer           FileRenamer // mechanism used to rename the file
	AllowedExtensions []string    // a list of extensions that will have content hashes generated for
}

// Default configuration for [HashedFS].
var ConfigDefault = Config{
	Hasher:            Sha256Hasher{},
	Renamer:           ExtensionRenamer,
	AllowedExtensions: []string{"js", "json", "png", "bmp", "jpeg", "jpg", "css", "ico"},
}
