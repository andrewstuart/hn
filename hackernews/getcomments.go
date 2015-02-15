package hackernews

//Retreives comments for a given article
import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const HN_WIDTH_MULTIPLIER = 40

func (c *client) GetComments(a *Article) error {
	if a.Comments == nil {
		a.Comments = make([]*Comment, 0)
	}

	articleUrl := fmt.Sprintf("%s/item?id=%d", c.RootUrl, a.Id)

	req, err := http.NewRequest("GET", articleUrl, nil)

	if err != nil {
		return fmt.Errorf("Error creating http request: %v", err)
	}

	doc, err := c.doReq(req)

	if err != nil {
		return fmt.Errorf("Error doing http request: %v", err)
	}

	commentStack := make([]*Comment, 1, 10)

	doc.Find("span.comment").Each(func(i int, comment *goquery.Selection) {
		user := comment.Parent().Find("a").First().Text()
		text := html.UnescapeString(comment.Text())

		//Get around HN's unpredictable "reply" nesting
		if last5 := len(text) - 5; len(text) > 0 && last5 > 0 && text[last5:] == "reply" {
			text = text[:last5]
		}

		c := &Comment{
			User:     user,
			Text:     text,
			Comments: make([]*Comment, 0),
		}

		//Get id
		if idAttr, exists := comment.Prev().Find("a").Last().Attr("href"); exists {
			idSt := strings.Split(idAttr, "=")[1]

			if id, err := strconv.Atoi(idSt); err == nil {
				c.Id = id
			}
		}

		//Track the comment offset for nesting.
		if width, exists := comment.Parent().Prev().Prev().Find("img").Attr("width"); exists {
			//If for whatever reason this errors, we'll still add it to article (offset 0)
			offset, _ := strconv.Atoi(width)

			offset = offset / HN_WIDTH_MULTIPLIER

			stackHeight := len(commentStack) - 1 //Index of the last element in the stack

			//If comment is nested above parent (offset > stackHeight)
			if offset > stackHeight {
				//Grow stack
				commentStack = append(commentStack, c)
				commentStack[stackHeight].Comments = append(commentStack[stackHeight].Comments, c)
			} else {
				if offset < stackHeight {
					//Trim the stack all the way to offset
					commentStack = commentStack[:offset+1]
				}

				commentStack[offset] = c

				//Add the comment to its parents
				if offset == 0 {
					//If stack is empty, use article
					a.Comments = append(a.Comments, c)
				} else {
					//Otherwise, go to last in the stack
					commentStack[offset-1].Comments = append(commentStack[offset-1].Comments, c)
				}
			}
		}
	})

	return nil
}
