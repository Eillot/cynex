// Reactor 提供指定路径的处理器绑定
// 支持浏览器Get与Post方法
package http

import (
	"cynex/cache"
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var reactor *Reactor

type Reactor struct {
	w http.ResponseWriter // http.ResponseWriter
	r *http.Request       // *http.Request

	pathTree *ptree       // 路径树
	dynamic  *cache.Cache // 路由动态缓存

	extraForm map[string]string // 路径变量
}

func (re *Reactor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	re.r = r
	re.w = w
	key := r.RequestURI + ":" + strings.ToUpper(r.Method)
	var function reflect.Value
	if val, err := re.dynamic.Get(key); err == nil {
		function = val.(reflect.Value)
		function.Call(nil)
	} else {
		if function, err = re.getHandler(key); err == nil {
			function.Call(nil)
			re.dynamic.Set(key, function)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found."))
		}
	}

}

// Get 提供GET方式的HTTP访问
// 参数 url:将要注册处理的请求路径;comp:使用此组件中的方法处理请求;function:使用（指定组件中的）此方法处理请求;
func Get(url string, comp interface{}, function string) {
	v := reflect.ValueOf(comp)
	hf := v.MethodByName(function)
	reactor.register(url, hf, "GET")
}

// Post 提供POST方式的HTTP访问
// 参数 url:将要注册处理的请求路径;comp:使用此组件中的方法处理请求;function:使用（指定组件中的）此方法处理请求;
func Post(url string, comp interface{}, function string) {
	v := reflect.ValueOf(comp)
	hf := v.MethodByName(function)
	reactor.register(url, hf, "POST")
}

func init() {
	pathTree := &ptree{
		root: &node{
			name: "",
			sub:  *new([]*node),
			val:  reflect.Value{},
		},
	}
	reactor = &Reactor{
		pathTree:  pathTree,
		dynamic:   cache.NewCache(),
		extraForm: make(map[string]string),
	}
}

// 路径处理注册
func (re *Reactor) register(path string, function reflect.Value, method string) {
	key := path + ":" + method
	re.addHandler(key, function)
}

type node struct {
	name string
	val  reflect.Value
	sub  []*node
}

type ptree struct {
	root *node
}

func (re *Reactor) addHandler(path string, function reflect.Value) {
	splits := strings.Split(path, "/")
	cnode := re.pathTree.root
	for i, val := range splits {
		if strings.TrimSpace(val) == "" {
			continue
		}
		if i < len(splits)-2 {
			if n, err := re.checkSetExist(cnode.sub, val); err == nil {
				cnode = n
			} else {
				n := new(node)
				n.val = reflect.Value{}
				n.name = val
				cnode.sub = append(cnode.sub, n)
				cnode = n
			}
			continue
		} else {
			if n, err := re.checkSetExist(cnode.sub, val); err == nil {
				n.val = function
			} else {
				n := new(node)
				n.val = function
				n.name = val
				cnode.sub = append(cnode.sub, n)
			}

		}
	}
}

func (re *Reactor) checkSetExist(sub []*node, name string) (*node, error) {
	for _, n := range sub {
		if n.name == name {
			return n, nil
		}
	}
	return nil, errors.New("not exist")
}

func (re *Reactor) getHandler(path string) (reflect.Value, error) {
	rerr := errors.New("not exist")
	splits := strings.Split(path, "/")
	cnode := re.pathTree.root
	for i, val := range splits {
		if strings.TrimSpace(val) == "" {
			continue
		}
		if i < len(splits)-2 {
			if n, err := re.checkGetExist(cnode.sub, val); err == nil {
				cnode = n
				continue
			} else {
				return reflect.Value{}, rerr
			}
		} else {
			if n, err := re.checkGetExist(cnode.sub, val); err == nil {
				return n.val, nil
			} else {
				return reflect.Value{}, rerr
			}
		}
	}
	return reflect.Value{}, rerr
}

func (re *Reactor) checkGetExist(sub []*node, name string) (*node, error) {
	rerr := errors.New("not exist")
	for _, n := range sub {
		if strings.Index(n.name, "(") == 0 && strings.Index(n.name, ")") == len(n.name)-1 {
			// 正则匹配
			reg, err := regexp.Compile(n.name[1 : len(n.name)-1])
			if err != nil {
				panic(err)
			}
			if reg.MatchString(name) {
				return n, nil
			} else {
				return nil, rerr
			}
		}
		if strings.Index(n.name, "[") == 0 && strings.Index(n.name, "]") == len(n.name)-1 {
			// 常用匹配
			t := n.name[1 : len(n.name)-1]
			if t == "*" {
				return n, nil
			}
			if strings.Index(t, "*") == 0 && strings.LastIndex(t, "*") == len(t)-1 {
				if strings.Contains(name, t[1:len(t)-1]) {
					return n, nil
				}
				return nil, rerr
			}
			if strings.Index(t, "*") == 0 {
				m := t[1:]
				if i := strings.Index(name, m); i > -1 && i+len(m) == len(name) {
					return n, nil
				}
				return nil, rerr
			}
			if strings.Index(t, "*") == len(t)-1 {
				m := t[:len(t)-1]
				if strings.Index(name, m) == 0 {
					return n, nil
				} else {
					return nil, rerr
				}
			}

		}
		if strings.Index(n.name, "{") == 0 && strings.Index(n.name, "}") == len(n.name)-1 {
			// 路径变量匹配
			re.extraForm[n.name[1:len(n.name)-1]] = name
			return n, nil
		}
		// 静态匹配
		if n.name == name {
			return n, nil
		}
	}
	return nil, rerr
}
