package main

import (
	"github.com/antchfx/xquery/xml"
	"fmt"
	"strings"
)
var (
	data=`<JBXX>
        <XYRBH>201702010001</XYRBH>
        <!-- 嫌疑人编号 -->
        <XM>张三</XM>
        <!-- 姓名 -->
        <CYM></CYM>
        <!-- 曾用名 -->
        <CH>小三</CH>
        <!-- 绰号 -->
        <ZJLX_DM>证件类型代码</ZJLX_DM>
        <!-- 证件类型代码，代码类型， 9910 -->
        <ZJLX_MC>居民身份证</ZJLX_MC>
        <!--证件类型名称 -->
        <ZJHM>513031198510124512</ZJHM>
        <!-- 证件号码 -->
        <XB_DM>9909180000001</XB_DM>
        <!-- 性别代码，代码类型，9909 -->
        <XB_MC>男性</XB_MC>
        <!-- 性别名称 -->
        <MZ_DM>9912180100001</MZ_DM>
        <!-- 民族代码，代码类型，9912 -->
        <MZ_MC>汉族</MZ_MC>
        <!-- 民族名称 -->
        <CSRQ>1985-10-12</CSRQ>
        <!-- 出生日期 -->
        <ZASNL>32</ZASNL>
        <!-- 作案时年龄 -->
        <GJ_DM>9911180200001</GJ_DM>
        <!-- 国籍代码，代码类型，9911 -->
        <GJ_MC>中国</GJ_MC>
        <!-- 国籍名称 -->
        <HJSZD_DM>9913000370102</HJSZD_DM>
        <!-- 户籍所在地代码，代码类型，9913 -->
        <HJSZD_MC>山东省济南市历下区</HJSZD_MC>
        <!-- 户籍所在地名称 -->
        <ZSD_DM>9913000370102</ZSD_DM>
        <!-- 住所地代码，代码类型，9913 -->
        <ZSD_MC>山东省济南市历下区</ZSD_MC>
        <!-- 住所地名称 -->
        <ZSDXXDZ>山东省济南市历下区</ZSDXXDZ>
        <!-- 住所地详细地址 -->
        <GZDW></GZDW>
        <!-- 工作单位 -->
        <GZDWSZD_DM></GZDWSZD_DM>
        <!-- 工作单位所在地代码，代码类型，9913 -->
        <GZDWSZD_MC></GZDWSZD_MC>
        <!-- 工作单位所在地名称 -->
        <ZJ_DM>9914180800001</ZJ_DM>
        <!-- 职级代码，代码类型，9914 -->
        <ZJ_MC>无职级</ZJ_MC>
        <!-- 职级名称 -->
        <ZW>无职务</ZW>
        <!-- 职务 -->
        <ZY_DM>9947199700074</ZY_DM>
        <!-- 职业代码，代码类型，9947 -->
        <ZY_MC>无业</ZY_MC>
        <!-- 职业名称 -->
        <SF_DM>9916192601000</SF_DM>
        <!-- 身份代码，代码类型，9916 -->
        <SF_MC>农民</SF_MC>
        <!-- 身份名称 -->
        <QTGZSF_DM>9970010000990</QTGZSF_DM>
        <!-- 其他关注身份代码，代码类型，9970 -->
        <QTGZSF_MC>非上述关注身份</QTGZSF_MC>
        <!-- 其他关注身份名称 -->
        <SFDWLD>N</SFDWLD>
        <!-- 是否单位领导，是是Y,否为N -->
        <SFDRSZ>N</SFDRSZ>
        <!-- 是否实职，是是Y,否为N -->
        <SJYZK_DM>9915180600013</SJYZK_DM>
        <!-- 受教育情况代码，代码类型，9915 -->
        <SJYZK_MC>受教育状况不详</SJYZK_MC>
        <!-- 受教育情况名称 -->
        <ZZMM_DM>9917180500001</ZZMM_DM>
        <!-- 政治面貌代码，代码类型，9917 -->
        <ZZMM_MC>群众</ZZMM_MC>
        <!-- 政治面貌名称 -->
        <RDDB_DM>9918181500001</RDDB_DM>
        <!-- 人大代码，代码类型，9918 -->
        <RDDB_MC>无</RDDB_MC>
        <!-- 人大代表名称 -->
        <ZXWY_DM>9919181400001</ZXWY_DM>
        <!-- 政协委员代码，代码类型，9919 -->
        <ZXWY_MC>无</ZXWY_MC>
        <!-- 政协委员名称 -->
        <SFNCJCZZRY>Y</SFNCJCZZRY>
        <!-- 是否农村基层组织人员 -->
        <SFCJZZRY>Y</SFCJZZRY>
        <!-- 是否村级组织人员 -->
        <SFCJR>N</SFCJR>
        <!--是否残疾人 -->
        <FDDLR></FDDLR>
        <!-- 是否法定代理人 -->
        <WCNFZXYRDJHRQK_DM></WCNFZXYRDJHRQK_DM>
        <!-- 未成年人监护人情况代码，代码类型，9920 -->
        <WCNFZXYRDJHRQK_MC></WCNFZXYRDJHRQK_MC>
        <!-- 未成年人监护人情况名称 -->
      </JBXX>
`
)
func main(){

	if node,err:=xmlquery.Parse(strings.NewReader(data));err!=nil{
		panic(err.Error())
	}else{
		visit(node)
	}
}
/*
type Node struct {
    Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

    Type         NodeType
    Data         string
    Prefix       string
    NamespaceURI string
    Attr         []xml.Attr
    // contains filtered or unexported fields
}
*/

//visit 宽度优先
func visit(node *xmlquery.Node) {
	fmt.Println("visit")
	next:=make([]*xmlquery.Node,0,10)
	fmt.Printf("data=%s\t prefix=%s\t url=%s\n",node.Data,node.Prefix,node.NamespaceURI)
	
	for ;node!=nil;node=node.NextSibling{
		if node.FirstChild!=nil{
			next=append(next,node.FirstChild)
		}
	}
	for _,n:=range next{
		visit(n)
	}
}