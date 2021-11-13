package main

import (
    "context"
    "encoding/csv"
    "errors"
    "flag"
    "fmt"
    "os"
    "strings"

    "google.golang.org/api/indexing/v3"
    "google.golang.org/api/option"
)

func getUrlsFromFile(filename string) ([]string, error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }

    var urls []string
    csvReader := csv.NewReader(f)
    csvReader.Read() // skip header
    for {
        record, err := csvReader.Read()
        if err != nil {
            break 
        }

        url := strings.TrimSpace(record[0])
        if len(url) != 0 {
            urls = append(urls, url)
        }
    }
    return urls, nil
}

func getUrls(targets []string) ([]string, error) {
    var urls []string

    for _, target := range(targets) {
        if strings.HasPrefix(target, "https://") || strings.HasPrefix(target, "http://") {
            urls = append(urls, target)
        } else {
            urlsFromFile, err := getUrlsFromFile(target)
            if err != nil {
                return nil, err
            }
            urls = append(urls, urlsFromFile...)
        }
    }
    return urls, nil
}

func parseArgs() (string, []string, error) {
    flag.Parse()
    if flag.NArg() < 2 {
        return "", nil, errors.New("specify mode and url")
    }

    urls, err := getUrls(flag.Args()[1:flag.NArg()])
    if err != nil {
        return  "", nil, err
    }

    mode := flag.Arg(0)
    if strings.Compare(mode, "update") == 0 {
        return "URL_UPDATED", urls, nil
    } else if strings.Compare(mode, "delete") == 0 {
        return "URL_DELETED", urls, nil
    } else {
        return "", nil, errors.New("mode can be update or delete")
    }
}

func notifyUrl(mode, url string) error {
    urlToNotify := indexing.UrlNotification{ Type: mode, Url: url }
    fmt.Println(urlToNotify)

    ctx := context.Background()
    svc, err := indexing.NewService(ctx, option.WithCredentialsFile("./search-console-api.json"))
    ntf := indexing.NewUrlNotificationsService(svc)
    resp, err := ntf.Publish(&urlToNotify).Do()
    if err != nil {
        return err
    }

    fmt.Println(resp)
    return nil
}

func main() {
    mode, urls, err := parseArgs()
    if err != nil {
	fmt.Println(err)
	os.Exit(1)
    }
    fmt.Printf("mode=%s\n", mode)
    for _, url := range(urls) {
        fmt.Println(url)
    }

    for _, url := range(urls) {
        err := notifyUrl(mode, url)
        if err != nil {
            fmt.Println(err)
        }
    }

    os.Exit(0)
}


