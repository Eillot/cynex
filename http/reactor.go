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
	ResponseWriter http.ResponseWriter // http.ResponseWriter
	Request        *http.Request       // *http.Request

	pathTree  *ptree       // 路径树
	pathCache *cache.Cache // 路由动态缓存

	extraForm map[string]string // 路径变量
}

func (re *Reactor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	re.Request = r
	re.ResponseWriter = w
	re.Request.ParseForm()
	key := strings.ToUpper(r.Method) + ":" + r.RequestURI
	var function reflect.Value
	if val, err := re.pathCache.Get(key); err == nil {
		function = val.(reflect.Value)
		function.Call(nil)
	} else {
		if function, err = re.getMatchHandler(key); err == nil {
			function.Call(nil)
			iLeft := strings.Index(key, "{")
			iRight := strings.Index(key, "}")
			if iLeft < 0 && iRight < 0 {
				// 不包含变量路径时，存储缓存
				// 包含变量路径时，需将路径参数值保存至Request，故不可使用缓存
				go re.pathCache.Set(key, function)
			}
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
			name: "ROOT",
			sub:  *new([]*node),
			val:  reflect.Value{},
		},
	}
	reactor = &Reactor{
		pathTree:  pathTree,
		pathCache: cache.NewCache(),
		extraForm: make(map[string]string),
	}
}

// 路径处理注册
func (re *Reactor) register(path string, function reflect.Value, method string) {
	key := method + ":" + path
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

// 保存路由配置
func (re *Reactor) addHandler(path string, function reflect.Value) {
	splits := strings.Split(path, "/")
	cnode := re.pathTree.root
	for i, val := range splits {
		if strings.TrimSpace(val) == "" {
			continue
		}
		if i < len(splits)-1 {
			if n, err := re.inSet(cnode.sub, val); err == nil {
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
			if n, err := re.inSet(cnode.sub, val); err == nil {
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

// 是否已经设定路由
func (re *Reactor) inSet(sub []*node, name string) (*node, error) {
	for _, n := range sub {
		if n.name == name {
			return n, nil
		}
	}
	return nil, errors.New("not exist")
}

// 获取匹配当前请求路径的处理方法
func (re *Reactor) getMatchHandler(path string) (reflect.Value, error) {
	rErr := errors.New("not exist")
	splits := strings.Split(path, "/")
	curNode := re.pathTree.root
	for i, val := range splits {
		if strings.TrimSpace(val) == "" {
			continue
		}
		if i < len(splits)-1 {
			if n, err := re.exists(curNode.sub, val); err == nil {
				curNode = n
				continue
			}
			return reflect.Value{}, rErr
		} else {
			if n, err := re.exists(curNode.sub, val); err == nil {
				return n.val, nil
			}
			return reflect.Value{}, rErr
		}
	}
	return reflect.Value{}, rErr
}

// 是否存在已经保存的路由设置
// 支持全路径匹配、正则匹配、常用匹配、变量路径
func (re *Reactor) exists(sub []*node, name string) (*node, error) {
	rErr := errors.New("not exist")
	for _, n := range sub {
		if strings.Index(n.name, "(") == 0 && strings.Index(n.name, ")") == len(n.name)-1 {
			// 正则匹配
			reg, err := regexp.Compile(n.name[1 : len(n.name)-1])
			if err != nil {
				panic(err)
			}
			if reg.MatchString(name) {
				return n, nil
			}
			return nil, rErr
		}
		if strings.Index(n.name, "[") == 0 && strings.Index(n.name, "]") == len(n.name)-1 {
			// 常用匹配，使用*号匹配常用字符串
			t := n.name[1 : len(n.name)-1]
			if t == "*" {
				return n, nil
			}
			if strings.Index(t, "*") == 0 && strings.LastIndex(t, "*") == len(t)-1 {
				if strings.Contains(name, t[1:len(t)-1]) {
					return n, nil
				}
				return nil, rErr
			}
			if strings.Index(t, "*") == 0 {
				m := t[1:]
				if i := strings.Index(name, m); i > -1 && i+len(m) == len(name) {
					return n, nil
				}
				return nil, rErr
			}
			if strings.Index(t, "*") == len(t)-1 {
				m := t[:len(t)-1]
				if strings.Index(name, m) == 0 {
					return n, nil
				}
				return nil, rErr
			}

		}
		if strings.Index(n.name, "{") == 0 && strings.Index(n.name, "}") == len(n.name)-1 {
			// 变量路径匹配，含有路径变量，将路径中的值作为匹配变量的值存储Request
			formName := n.name[1 : len(n.name)-1]
			formValue := re.Request.Form[formName]
			formValue = append(formValue, name)
			re.Request.Form[formName] = formValue
			return n, nil
		}
		// 全路径匹配
		if n.name == name {
			return n, nil
		}
	}
	return nil, rErr
}
