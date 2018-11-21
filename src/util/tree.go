package util

import (
	"fmt"
)

// Slice2Tree 将切片数据转换为树形数据
func Slice2Tree(sliceDatas []map[string]interface{}, idField, pidField string) []map[string]interface{} {
	var r []map[string]interface{}
	index := make(map[string]interface{})

	for _, val := range sliceDatas {
		id := fmt.Sprint(val[idField])
		index[id] = val
	}

	for _, val := range sliceDatas {
		pid := fmt.Sprint(val[pidField])
		if _, ok := index[pid]; !ok || pid == "" {
			r = append(r, val)
		} else {
			pval := index[pid].(map[string]interface{})
			if _, ok := pval["children"]; !ok {
				var n []map[string]interface{}
				n = append(n, val)
				pval["children"] = &n
			} else {
				nodes := pval["children"].(*[]map[string]interface{})
				*nodes = append(*nodes, val)
			}
		}
	}
	return r
}

// ConvertToViewTree 转换树形数据为视图树数据
func ConvertToViewTree(treeDatas []map[string]interface{}, labelField, valueField, keyField string) []map[string]interface{} {
	for _, node := range treeDatas {
		node["title"] = node[labelField]
		node["value"] = node[valueField]
		node["key"] = node[keyField]
		child, ok := node["children"]
		if ok {
			node["children"] = ConvertToViewTree(*child.(*[]map[string]interface{}), labelField, valueField, keyField)
		}
	}
	return treeDatas
}
