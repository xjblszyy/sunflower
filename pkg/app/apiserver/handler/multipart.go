package handler

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"path"
	"strings"

	"github.com/h2non/filetype"
	"go.uber.org/zap"

	"sunflower/config"
)

func multipartFile(ctx context.Context, req io.ReadCloser, contentType string) ([]multipart.File, error) {
	defer req.Close()

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		zap.L().Error("parse media type failed", zap.Error(err))
		return nil, MakeInternalServerError(ctx, err.Error())
	}
	mr := multipart.NewReader(req, params["boundary"])

	fileMaxSize := config.C.File.MaxSize
	maxMemory := fileMaxSize + int64(1<<20)
	form, err := mr.ReadForm(maxMemory)
	if err != nil {
		return nil, MakeInternalServerError(ctx, err.Error())
	}
	fileHandles, ok := form.File["file"]
	if !ok {
		return nil, MakeBadRequestError(ctx, "没有上传任何文件")
	}
	if len(fileHandles) > config.C.File.MaxCount {
		return nil, MakeBadRequestError(ctx, "超过最大文件上传数量")
	}

	res := make([]multipart.File, 0, len(fileHandles))
	for i := 0; i < len(fileHandles); i++ {
		fileHandle := fileHandles[i]
		if fileHandle.Size > fileMaxSize {
			ferr := newFileSizeErr(fileMaxSize)
			if ferr != nil {
				return nil, MakeBadRequestError(ctx, ferr.Error())
			}
		}

		fileName := fileHandle.Filename
		file, err := fileHandle.Open()
		if err != nil {
			zap.L().Error("open file failed", zap.Error(err))
			return nil, MakeInternalServerError(ctx, err.Error())
		}

		if !CheckFileType(fileName, file) {
			return nil, MakeBadRequestError(ctx, "只能上传xlsx类型文件")
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			zap.L().Error("file seek failed", zap.Error(err))
			return nil, MakeInternalServerError(ctx, err.Error())
		}
		res = append(res, file)
	}

	return res, nil
}

// FileSizeErr 文件size 错误
type FileSizeErr struct {
	Message string
}

// newFileSizeErr 构建FileSizeErr
// @param fileMaxSize 允许的最大文件大小 单位：byte
func newFileSizeErr(fileMaxSize int64) error {
	var errStr string
	switch {
	case fileMaxSize >= 1024*1024:
		errStr = fmt.Sprintf("文件大小不能超过 %d MB", fileMaxSize>>20)
	case fileMaxSize >= 1024:
		errStr = fmt.Sprintf("文件大小不能超过 %d KB", fileMaxSize>>10)
	default:
		errStr = fmt.Sprintf("文件大小不能超过 %d B", fileMaxSize)
	}
	return FileSizeErr{
		Message: errStr,
	}
}

func (fe FileSizeErr) Error() string {
	return fe.Message
}

func CheckFileType(fileName string, file multipart.File) bool {
	var validExt = map[string]bool{
		"xlsx": true,
	}

	// 先通过拓展名去判断是否正确
	fileExtension := path.Ext(fileName)
	ext := strings.ToLower(strings.Replace(fileExtension, ".", "", -1))
	_, isValid := validExt[ext]
	if !isValid {
		zap.L().Error("file type is not valid", zap.String("Extension", fileExtension))
		return false
	}

	// 判断只能上传csv类型
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		zap.L().Error("read file failed", zap.Error(err))
		return false
	}

	kind, err := filetype.Match(fileData)
	if err != nil {
		zap.L().Error("file type.Match fail", zap.Error(err))
		return false
	}

	if kind == filetype.Unknown {
		zap.L().Error("file type.Unknown")
		return false
	}

	_, isValid = validExt[kind.Extension]
	if !isValid {
		zap.L().Error("file type is not valid", zap.String("Extension", kind.Extension))
		return false
	}
	return true
}
