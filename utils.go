package main

// func apply[I interface{}, O interface{}](vs []I, f func(I) O) []O {
//    vsm := make([]O, len(vs))
//    for i, v := range vs {
//        vsm[i] = f(v)
//    }
//    return vsm
//}

func apply(vs []interface{}, f func(interface{}) interface{}) interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
