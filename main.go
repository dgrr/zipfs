package main

import (
	"archive/zip"
	"log"
	"os"
	"path"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/afero"
)

var (
	uid        = uint32(os.Getuid())
	gid        = uint32(os.Getgid())
	zipfile    string
	rootFs     = afero.NewMemMapFs()
	cacheFiles = make(map[string]*zip.File)
)

func main() {
	if len(os.Args) < 3 {
		log.Printf("%s <zip file> <mount point>\n", os.Args[0])
		os.Exit(0)
	}

	zipfile = os.Args[1]
	file, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	ctx := daemon.Context{
		LogFileName: path.Join(os.TempDir(), "zipfs"),
	}
	dmn, err := ctx.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if dmn != nil { // parent
		return
	}

	created := false
	if _, err := os.Stat(os.Args[2]); err != nil {
		created = true
		os.Mkdir(os.Args[2], 0755)
	}
	for _, file := range file.File {
		name := path.Join("/", file.Name)
		dir := path.Dir(name)
		rootFs.MkdirAll(dir, 0777)
		f, err := rootFs.Create(name)
		if err == nil {
			f.Close()
		}
		cacheFiles[name] = file
	}
	c, err := fuse.Mount(
		os.Args[2],
		fuse.FSName("zipfs"),
		fuse.Subtype("zipfs"),
		fuse.LocalVolume(),
		fuse.AllowOther(),
		fuse.VolumeName("zip-volume"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	err = fs.Serve(c, &FS{})
	if err != nil {
		log.Fatalln(err)
	}
	if created {
		os.Remove(os.Args[2])
	}
}

var _ fs.FS = (*FS)(nil)

type FS struct{}

func (root *FS) Root() (fs.Node, error) {
	dir := &Dir{"/"}
	return dir, nil
}
