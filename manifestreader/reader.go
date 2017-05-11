package manifestreader

import (
	"io/ioutil"
	"encoding/json"
)

type ManifestType struct {
	Author string
	Files []FileType
	ManifestType string
	ManifestVersion int
	Minecraft MinecraftType
	Name string
	Overrides string
	ProjectID int
	Version string
}

type FileType struct {
	FileID int
	ProjectID int
	Required bool
}

type ModloadersType struct {
	Id string
	Primary bool
}

type MinecraftType struct {
	Modloaders []ModloadersType
	Version string
}

var manifestObject ManifestType

func Read(manifest string) error {
	file, err := ioutil.ReadFile(manifest)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &manifestObject)
	if err != nil {
		return err
	}
	return nil
}

func OverridesDir() string {
	return manifestObject.Overrides
}

func Modlist() []FileType {
	return manifestObject.Files
}

func Modloaders() []ModloadersType {
	return manifestObject.Minecraft.Modloaders
}

func ModpackVersion() string {
	return manifestObject.Version
}

func ModpackName() string {
	return manifestObject.Name
}

func MinecraftVersion() string {
	return manifestObject.Minecraft.Version
}
