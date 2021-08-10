package download

import (
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/sirupsen/logrus"
)

type DownloadInterface interface {
	Add(path string, url string) error
}

func New(dw DownloadInterface, path PathInterface) *dm {
	return &dm{
		dw:   dw,
		path: path,
	}
}

// dm 管理下载
type dm struct {
	dw   DownloadInterface
	path PathInterface
}

func (d *dm) Created(collection *data.Collection) {
	return
}

func (d dm) Added(detail *data.Source, collection *data.Collection) {
	if err := d.dw.Add(d.path.Path(collection), detail.GetDownloadUrl()); err != nil {
		logrus.Error("添加下载任务失败:", err)
	}
}
