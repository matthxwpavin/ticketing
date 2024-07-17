package pmux

import (
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
	"github.com/matthxwpavin/ticketing/fmts"
	"github.com/matthxwpavin/ticketing/logging/sugar"

	"go.uber.org/zap"
)

type Router struct {
	r    *mux.Router
	opts *RouterOptions
}

func NewRouter(opts ...RouterOptionsFunc) *Router {
	var o RouterOptions
	for _, opt := range opts {
		opt(&o)
	}
	return &Router{r: mux.NewRouter(), opts: &o}
}

func (s *Router) HandleFunc(path string, handler HandleFunc) *Route {
	r := s.r.HandleFunc(path, handler.HTTPHandleFunc(s.opts))
	return NewRoute(r, s.opts)
}

func (s *Router) PathPrefix(tpl string) *Route {
	return NewRoute(s.r.PathPrefix(tpl), s.opts)
}

func (s *Router) Router() *mux.Router {
	return s.r
}

type RouterOptions struct {
	Validator *validator.Validate
	Trans     ut.Translator
	Logger    *zap.SugaredLogger
}

type RouterOptionsFunc func(opts *RouterOptions)

func WithValidatorAndTranslation(opts *RouterOptions) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validator := validator.New(validator.WithRequiredStructEnabled())

	en_translations.RegisterDefaultTranslations(validator, trans)
	opts.Validator = validator
	opts.Trans = trans
}

func WithLogger(opts *RouterOptions) {
	logger, err := sugar.New()
	if err != nil {
		fmts.Panicf("failed to new logger: %v", err)
	}
	opts.Logger = logger
}

type HandleFunc func(*ResponseWriter, *Request)

func (f HandleFunc) HTTPHandleFunc(opts *RouterOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if opts.Logger != nil {
			ctx := sugar.WithContext(r.Context(), opts.Logger)
			r = r.WithContext(ctx)
		}
		f(
			&ResponseWriter{
				ResponseWriter: w,
				trans:          opts.Trans,
			},
			&Request{
				r:         r,
				validator: opts.Validator,
			},
		)
	}
}

type Route struct {
	r    *mux.Route
	opts *RouterOptions
}

func NewRoute(r *mux.Route, opts *RouterOptions) *Route {
	return &Route{r: r, opts: opts}
}

func (s *Route) Subrouter() *Router {
	return &Router{r: s.r.Subrouter(), opts: s.opts}
}

func (s *Route) Handler(handler http.Handler) *Route {
	s.r.Handler(handler)
	return s
}

func (s *Route) HandlerFunc(handler HandleFunc) *Route {
	s.r.HandlerFunc(handler.HTTPHandleFunc(s.opts))
	return s
}
