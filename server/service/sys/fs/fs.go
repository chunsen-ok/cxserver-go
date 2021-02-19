package fs

import (
	"cxfw/service/internal/router"
	"cxfw/types"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const (
	resDir = "./.res"
)

func Init(r gin.IRouter) {
	os.Mkdir(resDir, os.ModeDir)

	g := r.Group("/fs")
	g.POST("/", router.Route(add))
	g.DELETE("/", router.Route(del))
	g.GET("/", get)
	g.GET("/all/", router.Route(getAll))
}

func add(c *gin.Context) (int, interface{}, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	files := form.File["uploads"]

	for _, file := range files {
		resname := filepath.Base(file.Filename)
		resPath := filepath.Join(resDir, resname)
		c.SaveUploadedFile(file, resPath)
	}

	return http.StatusOK, nil, nil
}

// [DELETE] /api/sys/fs/?name=<resname>
func del(c *gin.Context) (int, interface{}, error) {
	dres := filepath.Base(c.Query("name"))
	if len(dres) != 0 {
		filepath.Walk(resDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && path != resDir {
				return filepath.SkipDir
			}

			resname := filepath.Base(path)
			if dres == resname {
				resPath := filepath.Join(resDir, resname)
				os.Remove(resPath)
				return errors.New("stop walk")
			}

			return nil
		})
	}

	return http.StatusOK, nil, nil
}

// [GET] /api/sys/fs/?name=<resname>
func get(c *gin.Context) {
	resname := filepath.Base(c.Query("name"))
	if len(resname) != 0 {
		respath := filepath.Join(resDir, resname)
		c.FileAttachment(respath, resname)
		return
	}
	c.JSON(http.StatusOK, types.Response{Err: nil, Body: nil})
}

func getAll(c *gin.Context) (int, interface{}, error) {
	resLst := make([]string, 0)
	filepath.Walk(resDir, func(path string, info os.FileInfo, err error) error {
		fmt.Println("--> res:", path)
		if info.IsDir() && path != resDir {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			resname := filepath.Base(path)
			resLst = append(resLst, resname)
		}

		return nil
	})

	return http.StatusOK, resLst, nil
}
