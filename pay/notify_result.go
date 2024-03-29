package pay

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/arden/wechat/util"
	"github.com/spf13/cast"
	"reflect"
	"sort"
	"strings"
)

// doc: https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_7&index=8

// NotifyResult 下单回调
type NotifyResult struct {
	ReturnCode *string `xml:"return_code"`
	ReturnMsg  *string `xml:"return_msg"`

	AppID              *string `xml:"appid" json:"appid"`
	MchID              *string `xml:"mch_id"`
	DeviceInfo         *string `xml:"device_info"`
	NonceStr           *string `xml:"nonce_str"`
	Sign               *string `xml:"sign"`
	SignType           *string `xml:"sign_type"`
	ResultCode         *string `xml:"result_code"`
	ErrCode            *string `xml:"err_code"`
	ErrCodeDes         *string `xml:"err_code_des"`
	OpenID             *string `xml:"openid"`
	IsSubscribe        *string `xml:"is_subscribe"`
	TradeType          *string `xml:"trade_type"`
	BankType           *string `xml:"bank_type"`
	TotalFee           *int    `xml:"total_fee"`
	SettlementTotalFee *int    `xml:"settlement_total_fee"`
	FeeType            *string `xml:"fee_type"`
	CashFee            *string `xml:"cash_fee"`
	CashFeeType        *string `xml:"cash_fee_type"`
	CouponFee          *int    `xml:"coupon_fee"`
	CouponCount        *int    `xml:"coupon_count"`

	// coupon_type_$n 这里只声明 3 个，如果有更多的可以自己组合
	CouponType0 *string `xml:"coupon_type_0"`
	CouponType1 *string `xml:"coupon_type_1"`
	CouponType2 *string `xml:"coupon_type_2"`
	CouponID0   *string `xml:"coupon_id_0"`
	CouponID1   *string `xml:"coupon_id_1"`
	CouponID2   *string `xml:"coupon_id_2"`
	CouponFeed0 *string `xml:"coupon_fee_0"`
	CouponFeed1 *string `xml:"coupon_fee_1"`
	CouponFeed2 *string `xml:"coupon_fee_2"`

	TransactionID *string `xml:"transaction_id"`
	OutTradeNo    *string `xml:"out_trade_no"`
	Attach        *string `xml:"attach"`
	TimeEnd       *string `xml:"time_end"`
}

// NotifyResp 消息通知返回
type NotifyResp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// VerifySign 验签
func (pcf *Pay) VerifySign(notifyRes NotifyResult) bool {
	// STEP1, 转换 struct 为 map，并对 map keys 做排序
	resMap := structs.Map(notifyRes)

	sortedKeys := make([]string, 0, len(resMap))
	for k := range resMap {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	// STEP2, 对key=value的键值对用&连接起来，略过空值 & sign
	var signStrings string
	for _, k := range sortedKeys {
		value := fmt.Sprintf("%v", cast.ToString(resMap[k]))
		if value != "" && strings.ToLower(k) != "sign" {
			signStrings = signStrings + getTagKeyName(k, &notifyRes) + "=" + value + "&"
		}
	}

	// STEP3, 在键值对的最后加上key=API_KEY
	signStrings = signStrings + "key=" + pcf.PayKey

	// STEP4, 进行MD5签名并且将所有字符转为大写.
	sign := util.MD5Sum(signStrings)
	if sign != *notifyRes.Sign {
		return false
	}
	return true
}

func getTagKeyName(key string, notifyRes *NotifyResult) string {
	s := reflect.TypeOf(notifyRes).Elem()
	f, _ := s.FieldByName(key)
	name := f.Tag.Get("xml")
	return name
}
