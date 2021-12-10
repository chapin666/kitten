package expression_test

import (
	"database/sql"
	"github.com/chapin666/kitten/pkg/expression/ext"
	"log"
	"os"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/chapin666/kitten/pkg/expression"
)

func TestMain(m *testing.M) {
	expression.GlobalImport("test", map[string]interface{}{
		"testAdd": func(a, b int) int {
			return a + b
		},
	})
	os.Exit(m.Run())
}

func TestResultBool(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", args{"1==1"}, true, false},
		{"2", args{"1+1"}, true, false},
		{"3", args{"1-1"}, false, false},
		{"4+", args{"(1==2)"}, false, false},
		{"5", args{`"b"`}, true, false},
		{"6", args{`"true"`}, true, false},
		{"7", args{`"on"`}, true, false},
		{"8", args{`"off"`}, false, false},
		{"9", args{`"false"`}, false, false},
		{"10", args{`[1,2]`}, true, false},
		{"11", args{`[]`}, false, false},
		{"12", args{`]`}, false, true},
		{"float1", args{`0.1`}, true, false},
		{"float2", args{`1.1`}, true, false},
		{"byte1", args{`byte(1)`}, true, false},
		{"byte2", args{`byte(0)`}, false, false},
		{"nil", args{`nil`}, false, false},
		{"un", args{`a`}, false, true},
		{"un", args{`1==a`}, false, false},
		{"un", args{`a==b`}, true, false},
		{"global1", args{`global.test_1==1`}, true, false},
		{"global2", args{`global.test_a=="a"`}, true, false},
		{"fun1_1", args{`fun1(global.test_1)`}, true, false},
		{"fun1_2", args{`fun1(2)`}, false, false},
	}
	exp := createDataTypeExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execBool(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultInt(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int1", args{"1+1"}, 2, false},
		{"int 2", args{"1*1"}, 1, false},
		{"int 3", args{"4/2"}, 2, false},
		{"int 4", args{"5/2"}, 2, false},
		{"int 5", args{"5"}, 5, false},
		{"var1", args{"a"}, 0, true},
		{"str1", args{`"a"`}, 0, true},
		{"str2", args{`"10"`}, 10, false},
		{"slice1", args{`[1,2]`}, 2, false},
		{"slice2", args{`[]`}, 0, false},
		{"map1", args{`{}`}, 0, false},
		{"map2", args{`{"a":1}`}, 1, false},
		{"map3", args{`{"a":1,"b":"a"}`}, 2, false},
		{"map4", args{`{"a":1,b:"a"}`}, 0, true},
		{"bool1", args{`true`}, 1, false},
		{"bool2", args{`1==2`}, 0, false},
		{"bool3", args{`2>1`}, 1, false},
		{"bool4", args{`2<1`}, 0, false},
		{"float1", args{`1.1`}, 1, false},
		{"float2", args{`0.9`}, 0, false},
		{"float3", args{`-0.9`}, 0, false},
		{"float4", args{`-1.9`}, -1, false},
		{"byte1", args{`byte(1)`}, 1, false},
		{"byte2", args{`byte(2)`}, 2, false},
		{"byte3", args{`byte(0)`}, 0, false},
		{"testAdd", args{`test.testAdd(1,2)`}, 3, false},
		{"testAdd", args{`test.testAdd(3,5)`}, 8, false},
		{"testAdd", args{`test.testAdd(10,20)`}, 30, false},
		{"ctx_10", args{`test.testAdd(ctx_10,20)`}, 30, false},

		{"nil", args{`nil`}, 0, false},
	}
	exp := createDataTypeExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execInt(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultFloat(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int1", args{"1+1"}, 2, false},
		{"int 2", args{"1*1"}, 1, false},
		{"int 3", args{"4/2"}, 2, false},
		{"int 4", args{"5/2"}, 2, false},
		{"int 5", args{"5"}, 5, false},
		{"var1", args{"a"}, 0, true},
		{"str1", args{`"a"`}, 0, true},
		{"str2", args{`"10"`}, 10, false},
		{"slice1", args{`[1,2]`}, 0, true},
		{"slice2", args{`[]`}, 0, true},
		{"map1", args{`{}`}, 0, true},
		{"map2", args{`{"a":1}`}, 0, true},
		{"map3", args{`{"a":1,"b":"a"}`}, 0, true},
		{"map4", args{`{"a":1,b:"a"}`}, 0, true},
		{"bool1", args{`true`}, 1, false},
		{"bool2", args{`1==2`}, 0, false},
		{"bool3", args{`2>1`}, 1, false},
		{"bool4", args{`2<1`}, 0, false},
		{"float1", args{`1.1`}, 1.1, false},
		{"float2", args{`0.9`}, 0.9, false},
		{"float3", args{`-0.9`}, -0.9, false},
		{"float4", args{`5.0/2.0`}, 2.5, false},
		{"float5", args{`-1.9`}, -1.9, false},
		{"float6", args{`float(5)/float(2)`}, 2.5, false},
		{"byte1", args{`byte(1)`}, 1, false},
		{"byte2", args{`byte(2)`}, 2, false},
		{"byte3", args{`byte(0)`}, 0, false},

		{"nil", args{`nil`}, 0, false},
	}
	exp := createDataTypeExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execFloat(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultString(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"int1", args{"1+1"}, "2", false},
		{"int 2", args{"1*1"}, "1", false},
		{"int 3", args{"4/2"}, "2", false},
		{"int 4", args{"5/2"}, "2", false},
		{"int 5", args{"5"}, "5", false},
		{"var1", args{"a"}, "", true},
		{"str1", args{`"a"`}, "a", false},
		{"str2", args{`"10"`}, "10", false},
		{"slice1", args{`[1,2]`}, "[1 2]", false},
		{"slice2", args{`[]`}, "[]", false},
		{"map1", args{`{}`}, "map[]", false},
		{"map2", args{`{"a":1}`}, `map[a:1]`, false},
		{"map3", args{`{"a":1,"b":"a"}`}, `map[a:1 b:a]`, false},
		{"map4", args{`{"a":1,b:"a"}`}, "", true},
		{"bool1", args{`true`}, "true", false},
		{"bool2", args{`1==2`}, "false", false},
		{"bool3", args{`2>1`}, "true", false},
		{"bool4", args{`2<1`}, "false", false},
		{"float1", args{`1.1`}, "1.1", false},
		{"float2", args{`0.9`}, "0.9", false},
		{"float3", args{`-0.9`}, "-0.9", false},
		{"float4", args{`5.0/2.0`}, "2.5", false},
		{"float5", args{`-1.9`}, "-1.9", false},
		{"float6", args{`float(5)/float(2)`}, "2.5", false},
		{"byte1", args{`byte(1)`}, "1", false},
		{"byte2", args{`byte(2)`}, "2", false},
		{"byte3", args{`byte(0)`}, "0", false},
		{"nil", args{`nil`}, "", false},
		{"ctx_a", args{"ctx_a"}, "a", false},
	}
	exp := createDataTypeExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execStr(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultSliceString(t *testing.T) {
	type args struct {
		scriptCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", args{`["1","2","3"]`}, []string{"1", "2", "3"}, false},
		{"2", args{`["1"]`}, []string{"1"}, false},
		{"3", args{`[""]`}, []string{""}, false},
		{"4", args{`[]`}, nil, true},
		{"5", args{`nil`}, nil, false},
		{"6", args{`a`}, nil, true},
	}
	exp := createDataTypeExpression()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exp.execSliceStr(tt.args.scriptCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("execSliceStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("execSliceStr() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestResultSQL(t *testing.T) {
	exp := createSQLExpression()
	out, err := exp.execSql(`sql.querySliceStr("select * from f_flow where id < ?", "id", 10)`)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%#v", out)
}

type dataTypeExpression struct {
	exe expression.Execer
}

func createDataTypeExpression() *dataTypeExpression {
	exp := expression.CreateExecer("")
	exp.PredefinedJson("global", map[string]interface{}{
		"test_1": 1,
		"test_a": "a",
	})
	exp.PredefinedVar("fun1", `fn(a) {
		return 1==a
	}`)
	exp.PredefinedVar("fun2", `fn(a) {
		return 1==a
	}`)
	return &dataTypeExpression{
		exe: exp,
	}
}

func (t *dataTypeExpression) exec(exp string) (*expression.OutData, error) {
	expCtx := expression.CreateExpContext()

	expCtx.AddVar("ctx_10", 10)
	expCtx.AddVar("ctx_a", "a")

	return t.exe.Exec(expCtx, exp)
}

func (t *dataTypeExpression) execBool(exp string) (bool, error) {
	out, err := t.exec(exp)
	if err != nil {
		return false, err
	}
	return out.Bool()
}

func (t *dataTypeExpression) execInt(exp string) (int, error) {
	out, err := t.exec(exp)
	if err != nil {
		return 0, err
	}
	return out.Int()
}

func (t *dataTypeExpression) execFloat(exp string) (float64, error) {
	out, err := t.exec(exp)
	if err != nil {
		return 0, err
	}

	return out.Float()
}

func (t *dataTypeExpression) execStr(exp string) (string, error) {
	out, err := t.exec(exp)
	if err != nil {
		return "", err
	}

	return out.String()
}

func (t *dataTypeExpression) execSliceStr(exp string) ([]string, error) {
	out, err := t.exec(exp)
	if err != nil {
		return nil, err
	}
	return out.SliceStr()
}


type sqlExpression struct {
	exe expression.Execer
}

func createSQLExpression() *sqlExpression {
	pwd, _ := os.Getwd()
	exp := expression.CreateExecer(pwd)
	exp.ScriptImportAlias("ext", "sql")
	return &sqlExpression{
		exe: exp,
	}
}

func (t *sqlExpression) execSql(sqlStr string) (interface{}, error) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/flows")
	if err != nil {
		log.Fatalf("open mysql faild: %s", err.Error())
	}
	ext.Reg(db)

	expCtx := expression.CreateExpContext()
	//expCtx.AddVar("ctx", context.Background())
	//expCtx.AddVar("db", db)
	out, err := t.exe.Exec(expCtx, sqlStr)
	if err != nil {
		return nil, err
	}
	return out, nil
}
