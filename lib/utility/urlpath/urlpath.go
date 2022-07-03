package urlpath

type Routes struct {
	Path string
}

func New(url string, postfix ...string) *Routes {
	return &Routes{Path: append(url, postfix)}
}

func (r Routes) URL(postfix ...string) string {
	return append(r.Path, postfix)
}

func (r Routes) Group(postfix string, handler func(router Routes)) {
	if handler != nil {
		r.Path = r.Path + postfix
		handler(r)
	}
}

func (r Routes) NewPath(newPath string, postfix ...string) Routes {
	r.Path = append(r.Path+newPath, postfix)
	return r
}

func append(data string, listAppend []string) string {
	for _, str := range listAppend {
		if str == "/" {
			continue
		}
		data += str
	}
	return data
}
