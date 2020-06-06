package main

import (
	"context"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"log"
)

const (
	// 规则定义
	module = `
		package access
		# 默认都不许访问
		default allow = false

		allow {
			# 识别为admin,可以访问
			input.identity = "admin"
		}

		allow {
			# 或者是GET方法
			input.method = "GET"
		}
	`
)

func main() {
	ctx := context.Background()

	// 规则定义，key会出现在错误信息中，用于识别
	compiler, err := ast.CompileModules(map[string]string{
		"access": module,
	})

	if err != nil {
		panic(err.Error())
	}

	// 进行查询
	query, err := rego.New(
		rego.Query("data.access.allow"),
		rego.Compiler(compiler),
	).PrepareForEval(ctx)

	if err != nil {
		panic("准备查询错误" + err.Error())
	}
	datas := map[string]map[string]interface{}{
		"任意GET": {
			"method": "GET",
		},
		"admin任意操作": {
			"identity": "admin",
		},
		"非adminGET": {
			"method":   "GET",
			"identity": "bob",
		},
		"非adminPost": {
			"method":   "POST",
			"identity": "bob",
		},
		"非法数据": {
			"method": 1,
		},
	}
	for key, value := range datas {
		// Run evaluation.
		results, err := query.Eval(ctx, rego.EvalInput(value))

		if err != nil {
			log.Printf("处理[%s]错误%s", key, err.Error())
		} else if len(results) == 0 {
			log.Printf("处理[%s]结果不确定", key)
		} else {
			log.Printf("处理[%s]结果[%v]", key, results[0].Expressions[0].Value.(bool))

		}
	}

}
