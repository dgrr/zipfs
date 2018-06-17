package main

import (
	"archive/zip"
	"context"
	"io"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// File fuse files
type File struct {
	f *zip.File
	r io.ReadCloser
}

func (f *File) Attr(_ context.Context, attr *fuse.Attr) error {
	attr.Size = uint64(f.f.FileInfo().Size())
	attr.Mode = f.f.Mode()
	attr.Mtime = f.f.Modified
	attr.Ctime = attr.Mtime
	attr.Crtime = attr.Ctime
	attr.Uid = uid
	attr.Gid = gid
	return nil
}

func (f *File) Open(_ context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	r, err := f.f.Open()
	if err == nil {
		f.r = r
	}
	return f, err
}

func (f *File) Read(_ context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	buf := make([]byte, req.Size)
	n, err := io.ReadFull(f.r, buf)
	if err == io.ErrUnexpectedEOF || err == io.EOF {
		err = nil
	}
	resp.Data = buf[:n]
	return err
}

func (f *File) Release(_ context.Context, req *fuse.ReleaseRequest) (err error) {
	return f.r.Close()
}

func (f *File) Write(_ context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	return fuse.ENOSYS
}
