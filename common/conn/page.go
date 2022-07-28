package conn

type Page struct {
	PageNo     int
	PageSize   int
	TotalPage  int
	TotalCount int
	FirstPage  bool
	LastPage   bool
	List       interface{}
}

func new(count int, pageNo int, pageSize int, list interface{}) *Page {
	tp := count / pageSize
	if count % pageSize > 0 {
		tp = count / pageSize + 1
	}
	return &Page{PageNo: pageNo, PageSize: pageSize, TotalPage: tp, TotalCount: count, FirstPage: pageNo == 1, LastPage: pageNo == tp, List: list}
}


func PageUtil(count int64, pageNo int, pageSize int, list interface{}) *Page {
	return new(int(count), pageNo, pageSize, list)
}
