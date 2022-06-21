package lib

type Map[K comparable, V comparable] map[K]V

func (mp *Map[K, V]) GetOrInsert(k K, alt V) V {
	var zero V
	if (*mp)[k] == zero {
		(*mp)[k] = alt
	}
	return (*mp)[k]
}
