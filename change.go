package vbody

type Change struct {
	r *Reader
}

func (T *Change) Update(key any, val any) {
	T.r.m.RLock()
	defer T.r.m.RUnlock()

	switch v := key.(type) {
	case string:
		if _, ok := T.r.M[v]; !ok {
			return
		}
		T.r.M[v] = val
	case int:
		if len(T.r.A) > v {
			T.r.A[v] = val
		}
	}
}

func (T *Change) Set(key any, val any) {
	T.r.m.RLock()
	defer T.r.m.RUnlock()

	switch v := key.(type) {
	case string:
		T.r.M[v] = val
	case int:
		if v >= len(T.r.A) {
			l := v - len(T.r.A) + 1
			a := make([]any, l)
			T.r.A = append(T.r.A, a...)
		}
		T.r.A[v] = val
	}
}

func (T *Change) Delete(key any) {
	T.r.m.RLock()
	defer T.r.m.RUnlock()

	switch v := key.(type) {
	case string:
		delete(T.r.M, v)
	case int:
		if len(T.r.A) > v {
			T.r.A = append(T.r.A[:v], T.r.A[v+1:]...)
		}
	}
}
