package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net/http"
    "github.com/gin-gonic/gin"
    // "fmt"
)

type Minutes struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

// var minutesList = []Minutes{}

func main() {
    dsn := "admin:gijirokugo@tcp(gijiroku-database.ctwwukcgkoxu.ap-southeast-1.rds.amazonaws.com:3306)/gijiroku-database?parseTime=true"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    r := gin.Default()


    r.GET("/minutes", func(c *gin.Context) {
    rows, err := db.Query("SELECT id, title, content FROM minutes")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()
    var list []Minutes
    for rows.Next() {
        var m Minutes
        if err := rows.Scan(&m.ID, &m.Title, &m.Content); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        list = append(list, m)
    }
    c.JSON(http.StatusOK, list)
})

    // 議事録追加
    r.POST("/minutes", func(c *gin.Context) {
    var m Minutes
    if err := c.ShouldBindJSON(&m); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    res, err := db.Exec("INSERT INTO minutes (title, content) VALUES (?, ?)", m.Title, m.Content)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    id, _ := res.LastInsertId()
    m.ID = int(id)
    c.JSON(http.StatusOK, m)
})

    // 議事録削除
    r.DELETE("/minutes/:id", func(c *gin.Context) {
    idParam := c.Param("id")
    res, err := db.Exec("DELETE FROM minutes WHERE id = ?", idParam)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    count, _ := res.RowsAffected()
    if count > 0 {
        c.JSON(http.StatusOK, gin.H{"message": "deleted"})
    } else {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
    }
})

    r.Run(":8080")
}