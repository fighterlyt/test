package main

import (
	"context"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	run("通过名称判断是否有经验", "workers.rego", "workers", "data.ages.skilled", map[string]interface{}{
		"name": "张三",
	})
	run("获取所有的hostname", "std.rego", "std", "data.example.hostnames[name]", nil)
	run(`获取单个存在的hostname"hydrogen"`, "std.rego", "std", `data.example.hostnames["hydrogen"]`, nil)
	run(`获取单个不存在的hostname"test"`, "std.rego", "std", `data.example.hostnames["test"]`, nil)

	run(`获取所有的instance`, "std.rego", "std", `data.example.instances`, nil)
	run(`获取所有的instance[x]`, "std.rego", "std", `data.example.instances[x]`, nil)

	run(`获取所有的instance[x]`, "match.rego", "match", `data.match.match`, map[string]interface{}{
		"upperOrderArea": "410000",
		"workerArea":     "411200",
	})

	run(`获取所有的instance[x]`, "match.rego", "match", `data.match.exactMatch`, map[string]interface{}{
		"upperOrderArea": "411200",
		"workerArea":     "411200",
	})
}
func run(desc, moduleFile, module, queryContent string, data map[string]interface{}) {

	current, _ := os.Getwd()
	moduleContent, err := ioutil.ReadFile(filepath.Join(current, "opa", "general", "data", moduleFile))

	if err != nil {
		log.Panicf("加载定义失败[%s]\n", err.Error())
	}

	ctx := context.Background()

	// 规则定义，key会出现在错误信息中，用于识别
	compiler, err := ast.CompileModules(map[string]string{
		module: string(moduleContent),
	})

	if err != nil {
		log.Panicf("编译定义失败[%s]\n", err.Error())
	}

	query, err := rego.New(
		rego.Query(queryContent),
		rego.Compiler(compiler),
	).PrepareForEval(ctx)

	if err != nil {
		log.Panicf("构建预查询失败[%s]\n", err.Error())
	}

	results, err := query.Eval(ctx, rego.EvalInput(data))

	rego.Trace()
	if err != nil {
		log.Printf("[%s]处理错误%s", desc, err.Error())
	} else if len(results) == 0 {
		log.Printf("[%s]处理结果不确定", desc)
	} else {
		for i := 0; i < len(results); i++ {
			for j := 0; j < len(results[i].Expressions); j++ {
				log.Printf("[%s]处理结果[%d][%d],[%v][%v]", desc, i, j, results[i].Expressions[j].Value, results[i].Expressions[j].Text)

			}
		}

	}
}
