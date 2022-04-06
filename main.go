package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-vk-api/vk"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	ovk "golang.org/x/oauth2/vk"
)

var (
	conf    *oauth2.Config
	logFile *os.File
	logger  *logrus.Logger
	err     error
	tmpls   []*template.Template
	domain  = "192.168.2.5:4100"
)

func main() {
	logFile, err = os.OpenFile(`errors.log.json`, os.O_RDWR|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(errors.Wrap(err, `Не удалось открыть файл для логов с ошибками`))
	}

	logger = logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.JSONFormatter{})

	//tmpls = []*template.Template{}

	//f, err := os.Open(`/home/ivan/goprj/src/myProjects/oauth/templates/index.html`)
	//if err != nil {
	//	logger.Fatalln(err)
	//}
	//b, err := ioutil.ReadAll(f)
	//
	//indexTmpl, err := template.New(``).Parse(string(b))
	//if err != nil {
	//	logger.Fatalln(err)
	//}
	//
	//f, err = os.Open(`/home/ivan/goprj/src/myProjects/oauth/templates/me.html`)
	//if err != nil {
	//	logger.Fatalln(err)
	//}
	//b, err = ioutil.ReadAll(f)
	//
	//meTmpl, err := template.New(``).Parse(string(b))
	//if err != nil {
	//	logger.Fatalln(err)
	//}
	//
	//tmpls = append(tmpls, indexTmpl)
	//tmpls = append(tmpls, meTmpl)

	conf = &oauth2.Config{
		ClientID:     "8013084",
		ClientSecret: "cpKTaWZcYHZedKi6WdTk",
		//RedirectURL:  "http://192.168.2.5:4100/mainpage",
		RedirectURL: fmt.Sprintf("\"http://%s/testauf\"", domain),
		Scopes:      []string{"email"},
		Endpoint:    ovk.Endpoint,
	}

	e := echo.New()
	e.GET(`/authVK`, mpage)
	e.GET(`/testauf`, auth)
	fmt.Println("http://" + domain + "/authVK")
	logger.Fatalln(e.Start(domain))
}

func mpage(c echo.Context) (err error) {
	time.Sleep(10 * time.Second)
	url := conf.AuthCodeURL("izi", oauth2.AccessTypeOffline)
	log.Println(url)
	logger.Infoln(url)
	//b := bytes.Buffer{}
	//err = tmpls[0].Execute(&b, url)
	//if err != nil {
	//	logger.Fatalln(err)
	//	return c.NoContent(http.StatusInternalServerError)
	//}
	//return c.HTML(http.StatusOK, b.String())
	//return c.Redirect(301, url)
	return c.NoContent(501)
}

func auth(c echo.Context) (err error) {
	ctx := context.Background()
	qpar := c.QueryParams()
	code := qpar.Get(`code`)
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		logger.Fatalln(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	client, err := vk.NewClientWithOptions(vk.WithToken(tok.AccessToken))
	if err != nil {
		log.Fatal(err)
	}
	client.Lang = `ru`
	user := getCurrentUser(client)
	logger.Infoln(user)
	return c.NoContent(http.StatusOK)
}

func getCurrentUser(api *vk.Client) []map[string]interface{} {
	users := []map[string]interface{}{}

	err = api.CallMethod("users.get", vk.RequestParams{
		"fields": "photo_400_orig,city",
	}, &users)

	if err != nil {
		logger.Fatalln(err)
	}

	return users
}
