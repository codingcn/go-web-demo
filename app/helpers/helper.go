package helpers

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 获取md5明文字串
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	x := h.Sum(nil)
	y := fmt.Sprintf("%x", x)
	return y
}
func DaySub(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours() / 24)
}
func httpGet() {
	resp, err := http.Get("http://www.test.com/demo/accept.php?id=1")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
func httpPost() {
	resp, err := http.Post("http://www.test.com/demo/accept.php",
		"application/x-www-form-urlencoded",
		strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
func httpPostForm() {
	resp, err := http.PostForm("http://www.test.com/demo/accept.php",
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}
func httpDo() {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://www.test.com/demo/accept.php", strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

// 取随机数
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// 判断值是否在数组中，字符串操作
func InArrayStr(value string, arr []string) bool {
	if len(arr) <= 0 {
		return false
	}

	for _, v := range arr {
		if value == v {
			return true
		}
	}

	return false
}

// 判断值是否在数组中，整数操作
func InArrayInt(value int, arr []int) bool {
	if len(arr) <= 0 {
		return false
	}

	for _, v := range arr {
		if value == v {
			return true
		}
	}

	return false
}
