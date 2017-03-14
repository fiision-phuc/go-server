package server

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/phuc0302/go-server/util"
)

var (
	// Cfg references to public config's instance.
	Cfg *Config

	// router references to router's instance.
	router *Router

	// redirectPaths references to HTTP status redirect instructions.
	redirectPaths map[int]string
)

// HandleGroupFunc defines type alias for group func callback handler.
type HandleGroupFunc func()

// HandleContextFunc defines type alias for request context func callback handler.
//
// @param
// - context {RequestContext} (a RequestContext's instance, will be created by router)
type HandleContextFunc func(*RequestContext)

// Adapter defines type alias for HandleContextFunc func decorator.
//
// @param
// - func {HandleContextFunc} (an implementation of HandleContextFunc)
//
// @return
// - func {HandleContextFunc} (a wrapper func to handle before & after events before an actual func)
//
//
// Thank Mat Ryer for instruction on how to implement 'Adapter Pattern'.
// @link: https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.7g2v827ux
type Adapter func(HandleContextFunc) HandleContextFunc

// Adapt generates decorator for HandleContextFunc func.
//
// @param
// - f {HandleContextFunc} (an implementation of HandleContextFunc)
// - adapters {Adapter} (a list of adapter func that user wish to be executed before an actual func)
//
// @return
// - func {HandleContextFunc} (a wrapper func to handle before & after events before an actual func in revert order of adapters)
//
//
// Thank Mat Ryer for instruction on how to implement 'Adapter Pattern'.
// @link: https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.7g2v827ux
func Adapt(f HandleContextFunc, adapters ...Adapter) HandleContextFunc {
	for i := len(adapters) - 1; i >= 0; i-- {
		adapter := adapters[i]
		f = adapter(f)
	}
	return f
}

// Initialize will init server either in sandbox mode or production mode.
//
// @param
// - sandboxMode {bool} (instruction in which config file should be loaded)
func Initialize(sandboxMode bool) {
	// Load config file
	if sandboxMode {
		Cfg = LoadConfig(Debug)
	} else {
		Cfg = LoadConfig(Release)
	}
	router = new(Router)
}

// Run will start HTTP server.
func Run() {
	address := generateAddress()
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    Cfg.ReadTimeout,
		WriteTimeout:   Cfg.WriteTimeout,
		MaxHeaderBytes: Cfg.HeaderSize,
		Handler:        serveHTTP(),
	}
	logrus.Infof("listening on %s", address)
	logrus.Fatal(server.ListenAndServe())
}

// RunTLS will start HTTPS server.
func RunTLS(certFile string, keyFile string) {
	if sslPath := util.GetEnv(util.SSLPath); len(sslPath) > 0 {
		certFile = fmt.Sprintf("%s/%s", sslPath, certFile)
		keyFile = fmt.Sprintf("%s/%s", sslPath, keyFile)
	}

	address := generateAddress()
	server := &http.Server{
		Addr:           address,
		ReadTimeout:    Cfg.ReadTimeout,
		WriteTimeout:   Cfg.WriteTimeout,
		MaxHeaderBytes: Cfg.HeaderSize,
		Handler:        serveHTTP(),
	}
	logrus.Infof("listening on %s\n", address)
	logrus.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}

// GroupRoute routes all URLs with same prefixURI.
//
// @param
// - prefixURI {string} (the prefix for url)
// - handler {HandleGroupFunc} (the callback func)
func GroupRoute(prefixURI string, handler HandleGroupFunc) {
	router.GroupRoute(prefixURI, handler)
}

// BindCopy routes copy request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindCopy(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Copy, patternURL, handler)
}

// BindDelete routes delete request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindDelete(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Delete, patternURL, handler)
}

// BindGet routes get request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindGet(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Get, patternURL, handler)
}

// BindHead routes head request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindHead(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Head, patternURL, handler)
}

// BindLink routes link request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindLink(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Link, patternURL, handler)
}

// BindOptions routes options request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindOptions(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Options, patternURL, handler)
}

// BindPatch routes patch request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindPatch(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Patch, patternURL, handler)
}

// BindPost routes post request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindPost(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Post, patternURL, handler)
}

// BindPurge routes purge request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindPurge(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Purge, patternURL, handler)
}

// BindPut routes put request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindPut(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Put, patternURL, handler)
}

// BindUnlink routes unlink request to registered handler.
//
// @param
// - patternURL {string} (the URL matching pattern)
// - handler {HandleContextFunc} (the callback func)
func BindUnlink(patternURL string, handler HandleContextFunc) {
	router.BindRoute(Unlink, patternURL, handler)
}

// generateAddress returns a string represent a domain and port that the server will listen on.
//
// @return
// - address {string} (the domain:port that server will listen on)
func generateAddress() (address string) {
	if port := util.GetEnv(util.Port); len(port) > 0 {
		address = fmt.Sprintf("%s:%s", Cfg.Host, port)
	} else {
		address = fmt.Sprintf("%s:%d", Cfg.Host, Cfg.Port)
	}
	return
}

// serveHTTP returns an implementation for http.Handler.
//
// @return
// - handler {http.Handler} (the http.Handler implementation)
func serveHTTP() http.Handler {
	methodsValidation := regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(Cfg.AllowMethods, "|")))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Recovery(w, r)
		method := strings.ToLower(r.Method)
		path := httprouter.CleanPath(r.URL.Path)

		/* Condition validation: validate request method */
		if !methodsValidation.MatchString(method) {
			panic(util.Status405())
		}

		// Find route to handle request
		if route, pathParams := router.MatchRoute(method, path); route != nil {
			context := CreateContext(w, r)
			if pathParams != nil {
				context.PathParams = pathParams
			}
			route.InvokeHandler(context)
		} else {
			if len(Cfg.StaticFolders) > 0 && method == Get {
				for prefix, folder := range Cfg.StaticFolders {

					if strings.HasPrefix(path, prefix) {
						path = strings.Replace(path, prefix, folder, 1)

						if file, err := os.Open(path); err == nil {
							defer file.Close()

							if info, _ := file.Stat(); !info.IsDir() {
								http.ServeContent(w, r, path, info.ModTime(), file)
								return
							}
						}
						panic(util.Status404())
					}
				}
			}
			panic(util.Status503())
		}
	})
}
