package ezXMLTreeMaker

type XMLTree struct {
	name    string
	value   string
	subTree []XMLTree
	attr    map[string]string
}

func NewXMLTree(name, value string) XMLTree {
	return XMLTree{name: name, value: value, subTree: []XMLTree{}, attr: make(map[string]string)}
}
func (xmlTree *XMLTree) SetValue(v string) {
	xmlTree.value = v
}
func (xmlTree *XMLTree) SetAttr(k string, v string) {
	xmlTree.attr[k] = v
}
func (xmlTree *XMLTree) AddNode(n XMLTree) {
	xmlTree.subTree = append(xmlTree.subTree, n)
}
func (xmlTree *XMLTree) StrValue() string {
	rs := ""
	rs += "<" + xmlTree.name
	for key, value := range xmlTree.attr {
		rs += " " + key + `="` + value + `"`
	}
	rs += ">"
	if len(xmlTree.subTree) > 0 {
		for _, subTree := range xmlTree.subTree {
			rs += "\r\n" + subTree.StrValue()
		}
		rs += "\r\n"
	} else {
		rs += xmlTree.value
	}
	rs += "</" + xmlTree.name + ">"
	return rs
}
