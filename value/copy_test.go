package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCast(t *testing.T) {
	// Test successful string cast
	result := Cast[string]("hello")
	assert.Equal(t, "hello", result)

	// Test successful int cast
	result2 := Cast[int](42)
	assert.Equal(t, 42, result2)

	// Test failed cast - returns zero value
	result3 := Cast[string](42)
	assert.Equal(t, "", result3)

	// Test with nil
	result4 := Cast[string](nil)
	assert.Equal(t, "", result4)

	// Test with custom struct
	type Person struct {
		Name string
		Age  int
	}
	person := Person{Name: "John", Age: 30}
	result5 := Cast[Person](person)
	assert.Equal(t, person, result5)

	// Test failed cast to struct
	result6 := Cast[Person]("not a person")
	assert.Equal(t, Person{}, result6)

	// Test with pointer types
	str := "test"
	result7 := Cast[*string](&str)
	assert.Equal(t, &str, result7)

	// Test slice cast
	slice := []int{1, 2, 3}
	result8 := Cast[[]int](slice)
	assert.Equal(t, slice, result8)
}

func TestCopyPrimitiveTypes(t *testing.T) {
	// Test int copying
	var fromInt, toInt int = 42, 0
	Copy(fromInt, &toInt)
	assert.Equal(t, 42, toInt)

	// Test string copying
	var fromStr, toStr string = "hello", ""
	Copy(fromStr, &toStr)
	assert.Equal(t, "hello", toStr)

	// Test bool copying
	var fromBool, toBool bool = true, false
	Copy(fromBool, &toBool)
	assert.Equal(t, true, toBool)

	// Test float copying
	var fromFloat, toFloat float64 = 3.14, 0.0
	Copy(fromFloat, &toFloat)
	assert.Equal(t, 3.14, toFloat)

	// Test with pointers
	fromIntPtr := 100
	var toIntPtr int
	Copy(&fromIntPtr, &toIntPtr)
	assert.Equal(t, 100, toIntPtr)
}

func TestCopyStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		City string
	}

	from := Person{Name: "John", Age: 30, City: "New York"}
	var to Person

	Copy(from, &to)
	assert.Equal(t, "John", to.Name)
	assert.Equal(t, 30, to.Age)
	assert.Equal(t, "New York", to.City)

	// Test with pointer source
	Copy(&from, &to)
	assert.Equal(t, from, to)
}

func TestCopyStructsPartialMatch(t *testing.T) {
	type PersonFrom struct {
		Name    string
		Age     int
		Email   string
		Unknown int
	}

	type PersonTo struct {
		Name string
		Age  int
		City string // This field won't be copied
	}

	from := PersonFrom{Name: "Alice", Age: 25, Email: "alice@test.com", Unknown: 99}
	var to PersonTo

	Copy(from, &to)
	assert.Equal(t, "Alice", to.Name)
	assert.Equal(t, 25, to.Age)
	assert.Equal(t, "", to.City) // Should remain empty
}

func TestCopyMaps(t *testing.T) {
	from := map[string]int{"a": 1, "b": 2}
	var to map[string]int

	Copy(from, &to)
	assert.Equal(t, from, to)

	// Test with pointer source
	Copy(&from, &to)
	assert.Equal(t, from, to)
}

func TestCopyIncompatibleTypes(t *testing.T) {
	// Test copying between incompatible types
	var fromInt int = 42
	var toString string

	Copy(fromInt, &toString)
	assert.Equal(t, "", toString) // Should remain unchanged

	// Test struct to primitive
	type Person struct {
		Name string
	}
	person := Person{Name: "John"}
	var toInt int

	Copy(person, &toInt)
	assert.Equal(t, 0, toInt) // Should remain unchanged
}

func TestCopyNonSettableDestination(t *testing.T) {
	// Test copying to non-settable destination (value instead of pointer)
	var from, to int = 42, 0
	Copy(from, to) // Note: not &to
	assert.Equal(t, 0, to) // Should remain unchanged
}

func TestCopyNilValues(t *testing.T) {
	// Test with nil source
	var to int = 42
	Copy(nil, &to)
	assert.Equal(t, 42, to) // Should remain unchanged

	// Test with nil destination
	from := 100
	Copy(from, nil) // Should not panic
}

func TestCopyComplexStructs(t *testing.T) {
	type Address struct {
		Street string
		City   string
	}

	type PersonFrom struct {
		Name    string
		Age     int
		Address Address
	}

	type PersonTo struct {
		Name    string
		Age     int
		Address Address
		Phone   string // Extra field
	}

	from := PersonFrom{
		Name: "Bob",
		Age:  35,
		Address: Address{
			Street: "123 Main St",
			City:   "Boston",
		},
	}

	var to PersonTo
	Copy(from, &to)

	assert.Equal(t, "Bob", to.Name)
	assert.Equal(t, 35, to.Age)
	assert.Equal(t, Address{Street: "123 Main St", City: "Boston"}, to.Address)
	assert.Equal(t, "", to.Phone) // Should remain empty
}
