package lib

type Map[K comparable, V comparable] map[K]V

func (mp *Map[K, V]) GetOrInsert(k K, alt V) V {
	var zero V
	if (*mp)[k] == zero {
		(*mp)[k] = alt
	}
	return (*mp)[k]
}

func NewMapFromSlice[K comparable, V any](keys []K, values []V) map[K]V {
	mp := make(map[K]V, len(keys))
	for i := range keys {
		mp[keys[i]] = values[i]
	}
	return mp
}

func MapSlice[T, U any](a []T, f func(a T) U) []U {
	b := make([]U, len(a))
	for i, elem := range a {
		b[i] = f(elem)
	}
	return b
}

type Pair[T, U any] struct {
	First  T
	Second U
}

func MakeSliceFromMap[K comparable, V any](mp map[K]V) []Pair[K, V] {
	pairs := make([]Pair[K, V], len(mp))
	for k, v := range mp {
		pairs = append(pairs, Pair[K, V]{k, v})
	}
	return pairs
}
