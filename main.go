package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
)

type Minutes struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

var minutesList = []Minutes{}

func main() {
    r := gin.Default()

    // 議事録一覧取得
    r.GET("/minutes", func(c *gin.Context) {
        c.JSON(http.StatusOK, minutesList)
    })

    // 議事録追加
    r.POST("/minutes", func(c *gin.Context) {
        var m Minutes
        if err := c.ShouldBindJSON(&m); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        m.ID = len(minutesList) + 1
        minutesList = append(minutesList, m)
        c.JSON(http.StatusOK, m)
    })

    // 議事録削除
    r.DELETE("/minutes/:id", func(c *gin.Context) {
        idParam := c.Param("id")
        var deleted bool
        for i, m := range minutesList {
            if idParam == "" {
                continue
            }
            if fmt.Sprintf("%d", m.ID) == idParam {
                minutesList = append(minutesList[:i], minutesList[i+1:]...)
                deleted = true
                break
            }
        }
        if deleted {
            c.JSON(http.StatusOK, gin.H{"message": "deleted"})
        } else {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        }
    })

    r.Run(":8080")
}