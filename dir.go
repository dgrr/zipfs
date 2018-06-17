package main

import (
	"context"
	"os"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

var (
	bytePool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 512)
		},
	}
)

type Dir struct{}

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
	ff, ok := cacheFiles[name]
	if ok {
		return &File{
			f: ff,
		}, nil
	}
	return nil, fuse.ENOENT
}

func (d *Dir) ReadDirAll(_ context.Context) ([]fuse.Dirent, error) {
	var res []fuse.Dirent
	for _, f := range cacheFiles {
		var de fuse.Dirent
		if f.Mode().IsDir() {
			de.Type = fuse.DT_Dir
		}
		de.Name = f.Name
		res = append(res, de)
	}
	return res, nil
}
