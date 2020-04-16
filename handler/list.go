package handler

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/labstack/echo"
)

var paginator = &pagination.Paginator{}

// NewSlice ...
func NewSlice(start, count, step int) []int {
	s := make([]int, count)
	for i := range s {
		s[i] = start
		start += step
	}
	return s
}

// FindListTest 테스트
func (h *Handler) FindListTest(c echo.Context) error {
	usernames := []string{"Larry Ellison", "Carlos Slim Helu", "Mark Zuckerberg", "Amancio Ortega ", "Jeff Bezos", " Warren Buffett ", "Bill Gates"}
	postsPerPage := 5
	paginator = pagination.NewPaginator(c.Request(), postsPerPage, len(usernames))
	fmt.Println(paginator.Offset())

	idrange := NewSlice(paginator.Offset(), postsPerPage, 1)

	//create a new page list that shows up on html
	myusernames := []string{}
	for _, num := range idrange {
		//Prevent index out of range errors
		if num <= len(usernames)-1 {
			myuser := usernames[num]
			myusernames = append(myusernames, myuser)
		}
	}

	return c.JSON(http.StatusOK, List(c, *paginator, myusernames))
}

// List ...
func List(c echo.Context, paginator pagination.Paginator, data interface{}) map[string]interface{} {
	host := c.Scheme() + "://" + c.Request().Host
	var previous *string
	if paginator.HasPrev() {
		prev := host + paginator.PageLinkPrev()
		previous = &prev
	}

	var next *string
	if paginator.HasNext() {
		n := host + paginator.PageLinkNext()
		next = &n
	}

	return map[string]interface{}{
		"next":     next,
		"previous": previous,
		"results":  data,
	}
}
