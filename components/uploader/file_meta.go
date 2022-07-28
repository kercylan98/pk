package uploader

import (
	"github.com/google/uuid"
	"path"
	"time"
)

// HTTP上传的文件的文件元信息。
type FileMeta struct {
	id           string     // 文件ID(uuid)
	name         string     // 文件名称
	storageName  string     // 在文件系统中的名称
	size         int64      // 文件大小
	ext          string     // 文件扩展名
	uploadAt     time.Time  // 文件上传时间
	hash         string     // 文件的hash
	isFastUpload bool       // 是否为快速上传返回的文件元信息
	savePath	 string		// 完整的存储路径
}

// 获取存储路径。
func (slf *FileMeta) GetSavePath() string {
	return slf.savePath
}

// 检查这个文件是否是快速上传的。
func (slf *FileMeta) IsFastUpload() bool {
	return slf.isFastUpload
}

// 获取文件id。
func (slf *FileMeta) GetId() string {
	return slf.id
}

// 获取文件的真实名称。
func (slf *FileMeta) GetName() string {
	return slf.name
}

// 获取该文件在文件系统中的存储的文件名。
func (slf *FileMeta) GetStorageName() string {
	return slf.storageName
}

// 获取文件占用的空间大小。
func (slf *FileMeta) GetSize() int64 {
	return slf.size
}

// 获取文件的拓展名。
func (slf *FileMeta) GetExt() string {
	return slf.ext
}

// 获取文件的上传时间。
func (slf *FileMeta) GetUploadAt() time.Time {
	return slf.uploadAt
}

// 获取文件的hash值。
func (slf *FileMeta) GetHash() string {
	return slf.hash
}

// 构建一个完备的文件元信息。
func NewFileMeta(name string, size int64, savePath string) *FileMeta {
	uid := uuid.New().String()
	ext := path.Ext(name)
	this := &FileMeta{
		id: uid,
		name: name,
		storageName: uid + ext,
		size: size,
		ext: ext,
		uploadAt: time.Now(),
		savePath: savePath,
	}
	return this
}

