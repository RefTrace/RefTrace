package nf

import "reft-go/parser"

func mkFromChannel() parser.IClassNode {
	path := parser.MakeWithoutCaching("Channel(List(Path))")
	arg := parser.NewParameter(parser.STRING_TYPE, "Path")
	mn := parser.NewMethodNode("fromPath", 0, path, []*parser.Parameter{arg}, nil, nil)
	cn := parser.MakeWithoutCaching("Channel")
	cn.AddMethod(mn)
	return cn
}
