package main

import (
	"os"
	"path"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

var (
	bytePool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 512)
		},
	}
)

type Dir struct {
	Name string
}

func (d *Dir) Attr(_ context.Context, attr *fuse.Attr) error {
	attr.Inode = 0
	attr.Uid = uid
	attr.Gid = gid
	attr.Mode = os.ModeDir | 0755
	return nil
}

func (d *Dir) Create(_ context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	return nil, nil, fuse.ENOSYS
}

func (d *Dir) Mkdir(_ context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	return nil, fuse.ENOSYS
}

func (d *Dir) Remove(_ context.Context, req *fuse.RemoveRequest) error {
	return fuse.ENOSYS
}

func (d *Dir) Rename(_ context.Context, req *fuse.RenameRequest, _ fs.Node) error {
	return fuse.ENOSYS
}

func (d *Dir) Lookup(_ context.Context, name string) (fs.Node, error) {
	target := path.Join(d.Name, name)
	info, err := rootFs.Stat(target)
	if err != nil {
		return nil, fuse.ENOENT
	}
	if info.IsDir() {
		return &Dir{
			Name: target,
		}, nil
	}
	ff, ok := cacheFiles[target]
	if ok {
		return &File{
			f: ff,
		}, nil
	}
	return nil, fuse.ENOENT
}

func (d *Dir) ReadDirAll(_ context.Context) ([]fuse.Dirent, error) {
	var res []fuse.Dirent
	file, err := rootFs.Open(d.Name)
	if err != nil {
		panic(err)
	}
	files, err := file.Readdir(0)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		var de fuse.Dirent
		if f.IsDir() {
			de.Type = fuse.DT_Dir
		} else {
			de.Type = fuse.DT_File
		}
		de.Name = f.Name()
		res = append(res, de)
	}
	return res, nil
}
