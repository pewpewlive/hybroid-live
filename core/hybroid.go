package core

import "fmt"

type LevelManifest struct {
	Name              string         `toml:"name" json:"name"`
	Descriptions      []string       `toml:"descriptions" json:"descriptions"`
	Information       string         `toml:"information" json:"information"`
	EntryPoint        string         `toml:"entry_point" json:"entry_point"`
	IsCasual          bool           `toml:"casual" json:"has_score_leaderboard"`
	MedalRequirements map[string]int `toml:"medal_requirements" json:"rank_thresholds_1p,omitempty"`
}

type ProjectConfig struct {
	Name            string `toml:"name"` // should be kebab-case
	OutputDirectory string `toml:"output_directory"`
}

type HybroidConfig struct {
	Level   LevelManifest `toml:"level"`
	Project ProjectConfig `toml:"project"`
	//Packages        []PackageConfig `toml:"packages"`
}

type FileInformation struct {
	DirectoryPath string // The directory the file is located at (relative)
	FileName      string // The name of the file (without an extension)
	FileExtension string // The extension of the file
}

func (fi *FileInformation) Path() string {
	if fi.DirectoryPath == "." {
		return fmt.Sprintf("%s%s", fi.FileName, fi.FileExtension)
	}

	return fmt.Sprintf("%s/%s%s", fi.DirectoryPath, fi.FileName, fi.FileExtension)
}

func (fi *FileInformation) NewPath(start string, end string) string {
	if fi.DirectoryPath == "." {
		return fmt.Sprintf("%s/%s%s", start, fi.FileName, end)
	}

	return fmt.Sprintf("%s/%s/%s%s", start, fi.DirectoryPath, fi.FileName, end)
}
