package utils

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/loveRyujin/go-mall/common/enum"
	"regexp"
	"time"
)

func CopyProperties(dst, src interface{}) error {
	if err := copier.CopyWithOption(dst, src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (dst interface{}, err error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, errors.New("src type is not time.Time")
					}
					return s.Format(enum.TimeFormatHyphenedYMDHIS), nil
				},
			},
			{
				SrcType: copier.String,
				DstType: time.Time{},
				Fn: func(src interface{}) (dst interface{}, err error) {
					s, ok := src.(string)
					if !ok {
						return nil, errors.New("src type is not time format string")
					}
					pattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`
					matched, _ := regexp.MatchString(pattern, s)
					if matched {
						return time.Parse(enum.TimeFormatHyphenedYMDHIS, s)
					}
					return nil, errors.New("src type is not time format string")
				},
			},
		},
	}); err != nil {
		return err
	}
	return nil
}
