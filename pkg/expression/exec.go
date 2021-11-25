package expression

import (
	"bytes"
	"fmt"
	"sync/atomic"

	"github.com/xushiwei/qlang"

	"github.com/pkg/errors"
)

// CreateExecer 创建表达式执行器
func CreateExecer(libs string) Execer {
	return &execExp{
		libs:       libs,
		predefined: predefined{data: make([]pairs, 0, 4)},
		imports:    make(map[string]string),
	}
}

type execExp struct {
	predefined
	libs string

	imports map[string]string
}

func (e *execExp) ScriptImport(model string) {
	e.imports[model] = ""
}
func (e *execExp) ScriptImportAlias(model, alias string) {
	e.imports[model] = alias
}
func (e *execExp) SetLibs(libs string) {
	e.libs = libs
}
func (e execExp) Exec(ctx ExpContext, exp string) (out *OutData, err error) {
	ql := qlangFromContext(ctx)
	ql.SetLibs(e.libs)
	resultKey, expdata := e.parse(ctx, exp)

	ok := make(chan struct{})

	go func() {
		defer close(ok)
		err = e.exec(ql, expdata)
		if err != nil {
			// 错误处理
			err = errors.Wrapf(err, "表达式( %s )执行失败:%v.", exp, err)
			return
		}
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
		if err != nil {
			err = errors.Wrapf(err, "执行失败:%v", err)
		}
	case <-ok:
		if err == nil {
			o := ql.Var(resultKey)
			out = &OutData{Result: o}
		}
	}

	return
}
func (e execExp) exec(ql *qlang.Qlang, exp []byte) (err error) {
	defer func() {
		if err != nil {
			return
		}
		if e := recover(); e != nil {
			err = errors.Errorf("执行表达式( %s )失败:%v", exp, e)
			err = errors.WithStack(err)
		}
	}()

	qlang.SetDumpCode("3")
	err = ql.SafeExec(exp, "")
	if err != nil {
		// 错误处理
		err = errors.Wrapf(err, "表达式( %s )执行失败:%v.", exp, err)
	}
	return
}
func (e execExp) parse(ctx ExpContext, exp string) (string, []byte) {

	ec := ctx.(*expContext)
	buff := bytes.NewBuffer(nil)

	key := creResultKey()
	parseExeTpl(buff, &tplOption{
		Import:    e.imports,
		ExecerVar: e.data,
		CtxVar:    ec.data,
		ResultKey: key,
		Exp:       exp,
	})

	return key, buff.Bytes()
}

func (e execExp) parsePredefined(key string, ps []pairs, buff *bytes.Buffer) {
	if len(ps) > 0 {

		buff.WriteString("// start " + key + " Predefined\n")
		for _, p := range ps {
			buff.WriteString(p.Key)
			buff.WriteString(" = ")
			buff.WriteString(p.Value)
			buff.WriteString("\n")
		}
		buff.WriteString("// end\n")
	}
}

var (
	counter = int64(0)
)

func creResultKey() string {
	c := atomic.AddInt64(&counter, 1)
	return fmt.Sprintf("__result%d__", c)
}
