package miyatama

type StringCollection interface {
	Add(s string) StringCollection
	Get(i int) string
	Size() int
}

type StringArray struct {
	items []string
}

func (array StringArray) Add(s string) StringArray {
	array.items = append(array.items, s)
	return array
}

func (array StringArray) Get(i int) string {
	return array.items[i]
}

func (array StringArray) Size() int {
	return len(array.items)
}
