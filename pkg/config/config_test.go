package config

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	DescribeTable("LoadFromMap",
		func(field string, value string, expected interface{}) {
			type config struct {
				Foo string `config:"foo"`
				Bar int    `config:"bar"`
				Baz bool   `config:"baz"`
			}

			key := strings.ToLower(field)
			m := map[string]string{key: value}

			var c config
			err := LoadFromMap(m, &c)

			r := reflect.ValueOf(c)
			Expect(err).ToNot(HaveOccurred())
			Expect(r.FieldByName(field).Interface()).To(Equal(expected))
		},
		func(field string, value string, expected interface{}) string {
			key := strings.ToLower(field)
			return fmt.Sprintf("should load the field %s with the value %+v from the key %s", field, expected, key)
		},
		// string fields
		Entry(nil, "Foo", "foo", "foo"),
		Entry(nil, "Foo", "123", "123"),
		// int fields
		Entry(nil, "Bar", "0", 0),
		Entry(nil, "Bar", "123", 123),
		Entry(nil, "Bar", "-456", -456),
		Entry(nil, "Bar", "+789", 789),
		// bool fields
		Entry(nil, "Baz", "true", true),
		Entry(nil, "Baz", "false", false),
		Entry(nil, "Baz", "1", true),
		Entry(nil, "Baz", "0", false),
	)

	DescribeTable("LoadFromMapP",
		func(field string, prefix string, value string, expected interface{}) {
			type config struct {
				Foo string `config:"foo"`
				Bar int    `config:"bar"`
				Baz bool   `config:"baz"`
			}

			key := prefix + strings.ToLower(field)
			m := map[string]string{key: value}

			var c config
			err := LoadFromMapP(m, prefix, &c)

			r := reflect.ValueOf(c)
			Expect(err).ToNot(HaveOccurred())
			Expect(r.FieldByName(field).Interface()).To(Equal(expected))
		},
		func(field string, prefix string, value string, expected interface{}) string {
			key := prefix + strings.ToLower(field)
			return fmt.Sprintf("should load the field %s with the value %+v from the key %s", field, expected, key)
		},
		// string fields
		Entry(nil, "Foo", "prefix-", "foo", "foo"),
		// int fields
		Entry(nil, "Bar", "k8s.io/", "0", 0),
		// bool fields
		Entry(nil, "Baz", "my.prefix/", "true", true),
	)

	It("should ignore a non-required field that is missing", func() {
		type config struct {
			Foo string `config:"foo"`
		}

		m := map[string]string{"bar": "bar"}

		var c config
		err := LoadFromMap(m, &c)
		Expect(err).ToNot(HaveOccurred())
		Expect(c.Foo).To(Equal(""))
	})

	It("should throw an error if a required field is missing", func() {
		type config struct {
			Foo string `config:"foo,required"`
		}

		m := map[string]string{
			"bar": "bar",
		}

		var c config
		err := LoadFromMap(m, &c)
		Expect(err).To(HaveOccurred())
	})

	It("should ignore non-tagged or explicitly ignored struct fields", func() {
		type config struct {
			Foo string
			Bar int `config:"-"`
		}

		m := map[string]string{"foo": "foo", "bar": "1"}

		var c config
		err := LoadFromMap(m, &c)
		Expect(err).ToNot(HaveOccurred())
		Expect(c.Foo).To(Equal(""))
		Expect(c.Bar).To(Equal(0))
	})

})
