# goilscrap: Go IsepLive Scrapper

A scrapper for iseplive.fr, a student website

```
go get github.com/aabizri/goil
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
err = goil.ExportToCSV(studentList,os.Stdout)
```

Post something
```golang
post := NewPost("Hello World", Divers, true, false)
err := session.PublishPost(post)
```

Delete a post
```golang
err := session.DeletePost(1750)
```

