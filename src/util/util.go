package util

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

func SHA1(str string) string{
	//产生一个散列值得方式是sha1.New()
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(str))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来对现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)
	//SHA1 值经常以 16 进制输出，使用%x 来将散列结果格式化为 16 进制字符串。
	return fmt.Sprintf("%x", bs)
}

//JSON字符串转换为map对象
func JsonStr2Map(jsonstr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonstr), &result); err != nil {
		return nil, err
	}
	return result, nil
}