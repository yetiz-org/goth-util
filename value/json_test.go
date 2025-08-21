package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshal(t *testing.T) {
	// Test with simple string
	result := JsonMarshal("hello")
	assert.Equal(t, `"hello"`, result)

	// Test with integer
	result = JsonMarshal(42)
	assert.Equal(t, "42", result)

	// Test with boolean
	result = JsonMarshal(true)
	assert.Equal(t, "true", result)

	// Test with nil
	result = JsonMarshal(nil)
	assert.Equal(t, "null", result)

	// Test with slice
	slice := []int{1, 2, 3}
	result = JsonMarshal(slice)
	assert.Equal(t, "[1,2,3]", result)

	// Test with map
	m := map[string]int{"a": 1, "b": 2}
	result = JsonMarshal(m)
	// Note: map order is not guaranteed, so we check both possibilities
	assert.True(t, result == `{"a":1,"b":2}` || result == `{"b":2,"a":1}`)

	// Test with struct
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	person := Person{Name: "John", Age: 30}
	result = JsonMarshal(person)
	assert.Equal(t, `{"name":"John","age":30}`, result)

	// Test with empty struct
	result = JsonMarshal(struct{}{})
	assert.Equal(t, "{}", result)

	// Test with empty slice
	result = JsonMarshal([]int{})
	assert.Equal(t, "[]", result)

	// Test with empty map
	result = JsonMarshal(map[string]int{})
	assert.Equal(t, "{}", result)
}

func TestJsonMarshalComplexTypes(t *testing.T) {
	// Test with nested struct
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	person := Person{
		Name: "Alice",
		Age:  25,
		Address: Address{
			Street: "123 Main St",
			City:   "Boston",
		},
	}

	result := JsonMarshal(person)
	expected := `{"name":"Alice","age":25,"address":{"street":"123 Main St","city":"Boston"}}`
	assert.Equal(t, expected, result)

	// Test with pointer to struct
	result = JsonMarshal(&person)
	assert.Equal(t, expected, result)

	// Test with slice of structs
	people := []Person{person, {Name: "Bob", Age: 35}}
	result = JsonMarshal(people)
	assert.Contains(t, result, `"name":"Alice"`)
	assert.Contains(t, result, `"name":"Bob"`)
}

func TestJsonMarshalSpecialValues(t *testing.T) {
	// Test with float values
	result := JsonMarshal(3.14)
	assert.Equal(t, "3.14", result)

	result = JsonMarshal(0.0)
	assert.Equal(t, "0", result)

	// Test with negative numbers
	result = JsonMarshal(-42)
	assert.Equal(t, "-42", result)

	result = JsonMarshal(-3.14)
	assert.Equal(t, "-3.14", result)

	// Test with string with special characters
	result = JsonMarshal("hello\nworld")
	assert.Equal(t, `"hello\nworld"`, result)

	result = JsonMarshal("hello\"world")
	assert.Equal(t, `"hello\"world"`, result)

	// Test with unicode characters
	result = JsonMarshal("Hello 世界")
	assert.Equal(t, `"Hello 世界"`, result)
}

func TestJsonMarshalErrorCases(t *testing.T) {
	// Test with function (which cannot be marshaled to JSON)
	fn := func() {}
	result := JsonMarshal(fn)
	assert.Equal(t, "", result)

	// Test with channel (which cannot be marshaled to JSON)
	ch := make(chan int)
	result = JsonMarshal(ch)
	assert.Equal(t, "", result)

	// Test with complex number (which cannot be marshaled to JSON)
	complex := complex(1, 2)
	result = JsonMarshal(complex)
	assert.Equal(t, "", result)

	// Test with struct containing unexported fields with json tags that would cause error
	type BadStruct struct {
		fn func() `json:"function"`
	}
	bad := BadStruct{fn: func() {}}
	result = JsonMarshal(bad)
	// Note: unexported fields are ignored by JSON marshaling, so we get an empty object
	assert.Equal(t, "{}", result)
}

func TestJsonMarshalWithTags(t *testing.T) {
	// Test struct with json tags
	type Product struct {
		ID    int    `json:"id"`
		Name  string `json:"product_name"`
		Price float64 `json:"price"`
		Hidden string `json:"-"` // This field should be ignored
	}

	product := Product{
		ID:     1,
		Name:   "Laptop",
		Price:  999.99,
		Hidden: "secret",
	}

	result := JsonMarshal(product)
	expected := `{"id":1,"product_name":"Laptop","price":999.99}`
	assert.Equal(t, expected, result)
	assert.NotContains(t, result, "secret")
	assert.NotContains(t, result, "Hidden")
}

func TestJsonMarshalEdgeCases(t *testing.T) {
	// Test with empty string
	result := JsonMarshal("")
	assert.Equal(t, `""`, result)

	// Test with zero values
	result = JsonMarshal(0)
	assert.Equal(t, "0", result)

	result = JsonMarshal(false)
	assert.Equal(t, "false", result)

	// Test with pointer to nil
	var nilPtr *string
	result = JsonMarshal(nilPtr)
	assert.Equal(t, "null", result)

	// Test with interface{} containing different types
	var iface interface{} = "test"
	result = JsonMarshal(iface)
	assert.Equal(t, `"test"`, result)

	iface = 123
	result = JsonMarshal(iface)
	assert.Equal(t, "123", result)
}
