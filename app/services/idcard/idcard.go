package idcard

import (
	"errors"
	"strconv"
	"time"
)

var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var valid_value = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
var valid_province = []string{
	"11", // 北京市
	"12", // 天津市
	"13", // 河北省
	"14", // 山西省
	"15", // 内蒙古自治区
	"21", // 辽宁省
	"22", // 吉林省
	"23", // 黑龙江省
	"31", // 上海市
	"32", // 江苏省
	"33", // 浙江省
	"34", // 安徽省
	"35", // 福建省
	"36", // 山西省
	"37", // 山东省
	"41", // 河南省
	"42", // 湖北省
	"43", // 湖南省
	"44", // 广东省
	"45", // 广西壮族自治区
	"46", // 海南省
	"50", // 重庆市
	"51", // 四川省
	"52", // 贵州省
	"53", // 云南省
	"54", // 西藏自治区
	"61", // 陕西省
	"62", // 甘肃省
	"63", // 青海省
	"64", // 宁夏回族自治区
	"65", // 新疆维吾尔自治区
	"71", // 台湾省
	"81", // 香港特别行政区
	"91", // 澳门特别行政区
}

// Check citizen number 18 valid.
func IsValidCitizenNo18(citizenNo18 *[]byte) bool {
	nLen := len(*citizenNo18)
	if nLen != 18 {
		return false
	}

	nSum := 0
	for i := 0; i < nLen-1; i++ {
		n, _ := strconv.Atoi(string((*citizenNo18)[i]))
		nSum += n * weight[i]
	}
	mod := nSum % 11
	if valid_value[mod] == (*citizenNo18)[17] {
		return true
	}

	return false
}

// Convert citizen 15 to 18.
func Citizen15To18(citizenNo15 []byte) []byte {
	nLen := len(citizenNo15)
	if nLen != 15 {
		return nil
	}

	citizenNo18 := make([]byte, 0)
	citizenNo18 = append(citizenNo18, citizenNo15[:6]...)
	citizenNo18 = append(citizenNo18, '1', '9')
	citizenNo18 = append(citizenNo18, citizenNo15[6:]...)

	sum := 0
	for i, v := range citizenNo18 {
		n, _ := strconv.Atoi(string(v))
		sum += n * weight[i]
	}
	mod := sum % 11
	citizenNo18 = append(citizenNo18, valid_value[mod])

	return citizenNo18
}

func IsLeapYear(nYear int) bool {
	if nYear <= 0 {
		return false
	}

	if (nYear%4 == 0 && nYear%100 != 0) || nYear%400 == 0 {
		return true
	}

	return false
}

// Check birthday's year month day valid.
func CheckBirthdayValid(nYear, nMonth, nDay int) bool {
	if nYear < 1900 || nMonth <= 0 || nMonth > 12 || nDay <= 0 || nDay > 31 {
		return false
	}

	curYear, curMonth, curDay := time.Now().Date()
	if nYear == curYear {
		if nMonth > int(curMonth) {
			return false
		} else if nMonth == int(curMonth) && nDay > curDay {
			return false
		}
	}

	if 2 == nMonth {
		if IsLeapYear(nYear) && nDay > 29 {
			return false
		} else if nDay > 28 {
			return false
		}
	} else if 4 == nMonth || 6 == nMonth || 9 == nMonth || 11 == nMonth {
		if nDay > 30 {
			return false
		}
	}

	return true
}

// Check province error_code valid.
func CheckProvinceValid(citizenNo []byte) bool {
	provinceCode := make([]byte, 0)
	provinceCode = append(provinceCode, citizenNo[:2]...)
	provinceStr := string(provinceCode)

	for i, _ := range valid_province {
		if provinceStr == valid_province[i] {
			return true
		}
	}

	return false
}

// Check citizen number valid.
func IsValidCitizenNo(citizenNo *[]byte) bool {
	nLen := len(*citizenNo)
	if nLen != 15 && nLen != 18 {
		return false
	}

	for i, v := range *citizenNo {
		n, _ := strconv.Atoi(string(v))
		if n >= 0 && n <= 9 {
			continue
		}

		if v == 'X' && i == 16 {
			continue
		}

		return false
	}

	if !CheckProvinceValid(*citizenNo) {
		return false
	}

	if nLen == 15 {
		*citizenNo = Citizen15To18(*citizenNo)
		if citizenNo == nil {
			return false
		}
	} else if !IsValidCitizenNo18(citizenNo) {
		return false
	}

	nYear, _ := strconv.Atoi(string((*citizenNo)[6:10]))
	nMonth, _ := strconv.Atoi(string((*citizenNo)[10:12]))
	nDay, _ := strconv.Atoi(string((*citizenNo)[12:14]))
	if !CheckBirthdayValid(nYear, nMonth, nDay) {
		return false
	}

	return true
}

// Get information from citizen number. Birthday, gender, province mask.
func GetCitizenNoInfo(citizenNo []byte) (err error, birthday int64, isMale bool, addrMask int) {
	err = nil
	birthday = 0
	isMale = false
	addrMask = 0
	if !IsValidCitizenNo(&citizenNo) {
		err = errors.New("Invalid citizen number.")
		return
	}

	// Birthday information.
	nYear, _ := strconv.Atoi(string(citizenNo[6:10]))
	nMonth, _ := strconv.Atoi(string(citizenNo[10:12]))
	nDay, _ := strconv.Atoi(string(citizenNo[12:14]))
	birthday = time.Date(nYear, time.Month(nMonth), nDay, 0, 0, 0, 0, time.Local).Unix()

	// Gender information.
	genderMask, _ := strconv.Atoi(string(citizenNo[16]))
	if genderMask%2 == 0 {
		isMale = false
	} else {
		isMale = true
	}

	// Address error_code mask.
	addrMask, _ = strconv.Atoi(string(citizenNo[:2]))

	return
}
