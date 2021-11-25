query = fn(query,args...) {
	return sqlctx.Query(query,args...)
}
one = fn(query, args...) {
	return sqlctx.One(query,args...)
}
querySliceStr = fn(query,key,args...) {
	return SliceStr(sqlctx.Query(query,args...),key)
}
export query, one, querySliceStr