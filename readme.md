# goilscrap: Go IsepLive Scrapper

A scrapper for iseplive.fr, a student website

```
go get github.com/aabizri/goilscrap
```

## Example:

Start a session
```golang
session,err := Login(USERNAME,PASSWORD, &http.Client{})
```

Scrap the website
```golang
studentList,err := sess.GetStudentList()
```        

Export it in csv format through whatever io.writer you prefer, here with `os.Stdout`
```golang
err = ExportToCSV(studentList,os.Stdout)
```


