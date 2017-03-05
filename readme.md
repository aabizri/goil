# goilscrap: Go IsepLive Scrapper

A scrapper for iseplive.fr, a student website

```
go get github.com/aabizri/goilscrap
```

## Example:

Start a session
```golang
import "net/http"
session,err := goilscrap.Login(USERNAME,PASSWORD, &http.Client{})
```

Scrap the website
```golang
studentList,err := session.GetStudentList()
```        

Export it in csv format through whatever io.writer you prefer, here with `os.Stdout`
```golang
err = goilscrap.ExportToCSV(studentList,os.Stdout)
```


