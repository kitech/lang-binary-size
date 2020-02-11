package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var cnt = 100
	flag.IntVar(&cnt, "cnt", cnt, "gen funcs count")
	flag.Parse()

	var xgs []FuncGener
	xgs = append(xgs, &GoFuncGen{})
	xgs = append(xgs, &NimFuncGen{})
	xgs = append(xgs, &CppFuncGen{})
	xgs = append(xgs, &CFuncGen{})
	// xgs = append(xgs, &VFuncGen{})

	for _, xg := range xgs {
		xg.Init(cnt)
		xg.Generate()
		xg.Compile()
		xg.Getresult()
	}
}

type FuncGener interface {
	Init(cnt int)
	Generate()
	Compile()
	Getresult()
}

type basegen struct {
	cnt  int
	code string
}

const newline = "\n"
const fnpfx = "A_long_func_name_maybe_very_long_"

///
type GoFuncGen struct {
	basegen
}

func (fg *GoFuncGen) Init(cnt int) {
	fg.cnt = cnt
}

func (fg *GoFuncGen) Generate() {
	cnt := fg.cnt
	code := "package main" + newline
	code += "/*\n*/" + newline
	code += "import \"C\"" + newline
	code += "import \"unsafe\"" + newline
	code += "var keeper_val = 0" + newline
	code += "func keeper() {keeper_val++}" + newline

	for i := 0; i < cnt; i++ {
		code += fmt.Sprintf("//go:noinline") + newline
		code += fmt.Sprintf("func %s%d (a0 int, a1 string) {", fnpfx, i) + newline
		code += "var innerstr string" + newline
		code += "var lineno = 123" + newline
		code += "var fromidx = 456" + newline
		code += "var toperr unsafe.Pointer" + newline
		code += "a0 = fromidx" + newline
		code += "a0 = lineno" + newline
		code += "a1 = innerstr" + newline
		code += "toperr = unsafe.Pointer(&toperr)" + newline
		code += "keeper()" + newline
		code += "}" + newline
	}

	code += "func main() {" + newline
	for i := 0; i < cnt; i++ {
		code += fmt.Sprintf("%s%d (123, \"foo\")", fnpfx, i) + newline
	}
	code += "}" + newline

	fg.code = code
	os.Mkdir("tmp", 0755)
	err := ioutil.WriteFile("./tmp/binsize.go", []byte(code), 0644)
	if err != nil {
		log.Println(err)
	}
}
func (fg *GoFuncGen) Compile() {
	cmd := exec.Command("go", "build", "-o", "binsize-go", "./tmp/binsize.go")
	errcc, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err, string(errcc))
	}
	cmd = exec.Command("ls", "-lh", "binsize-go")
	btime := time.Now()
	outcc, err := cmd.CombinedOutput()
	log.Println(string(outcc), time.Since(btime))

	{
		// cmdstr := `go build -v -p 1 -gcflags "-N -l" -ldflags "-w -s" -o binsize-go tmp/binsize.go`
		cmd := exec.Command("go", "build",
			"-gcflags", "-N -l", "-ldflags", "-w -s",
			"-o", "binsize-go", "./tmp/binsize.go")
		errcc, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err, string(errcc))
		}
		cmd = exec.Command("ls", "-lh", "binsize-go")
		btime := time.Now()
		outcc, err := cmd.CombinedOutput()
		log.Println(string(outcc), time.Since(btime))
	}
}
func (fg *GoFuncGen) Getresult() {

}

///
type CFuncGen struct {
	basegen
}

func (fg *CFuncGen) Init(cnt int) {
	fg.cnt = cnt
}

func (fg *CFuncGen) Generate() {
	cnt := fg.cnt
	code := "//" + newline

	code += "int keeper_val = 0;" + newline
	code += "void keeper() {keeper_val++;}" + newline

	for i := 0; i < cnt; i++ {
		// code += fmt.Sprintf("//go:noinline") + newline
		code += fmt.Sprintf("void %s%d (int a0, const char* a1) {", fnpfx, i) + newline
		code += "  char* innerstr;" + newline
		code += "  int lineno;" + newline
		code += "  int fromidx;" + newline
		code += "  a0 = lineno;" + newline
		code += "  a0 = fromidx;" + newline
		// code += "  a1 = innerstr" + newline
		code += "keeper();" + newline
		code += "}" + newline
	}

	code += "int main(int argc, char**argv) {" + newline
	for i := 0; i < cnt; i++ {
		code += fmt.Sprintf("%s%d (123, \"foo\");", fnpfx, i) + newline
	}
	code += "return 0;" + newline
	code += "}" + newline

	fg.code = code
	os.Mkdir("tmp", 0755)
	err := ioutil.WriteFile("./tmp/binsize.c", []byte(code), 0644)
	if err != nil {
		log.Println(err)
	}
}
func (fg *CFuncGen) Compile() {
	cmd := exec.Command("gcc", "-g", "-O2", "-o", "binsize-c", "./tmp/binsize.c")
	errcc, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err, string(errcc))
	}
	cmd = exec.Command("ls", "-lh", "binsize-c")
	btime := time.Now()
	outcc, err := cmd.CombinedOutput()
	log.Println(string(outcc), time.Since(btime))

	{
		cmd := exec.Command("gcc", "-Os", "-o", "binsize-c", "./tmp/binsize.c")
		errcc, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err, string(errcc))
		}
		cmd = exec.Command("ls", "-lh", "binsize-c")
		btime := time.Now()
		outcc, err := cmd.CombinedOutput()
		log.Println(string(outcc), time.Since(btime))

	}
}
func (fg *CFuncGen) Getresult() {

}

///
type CppFuncGen struct {
	basegen
}

func (fg *CppFuncGen) Init(cnt int) {
	fg.cnt = cnt
}

func (fg *CppFuncGen) Generate() {
	cnt := fg.cnt
	code := "//" + newline

	code += "int keeper_val = 0;" + newline
	code += "void keeper() {keeper_val++;}" + newline

	for i := 0; i < cnt; i++ {
		// code += fmt.Sprintf("//go:noinline") + newline
		code += fmt.Sprintf("void %s%d (int a0, const char* a1) {", fnpfx, i) + newline
		code += "keeper();" + newline
		code += "}" + newline
	}

	code += "int main(int argc, char**argv) {" + newline
	for i := 0; i < cnt; i++ {
		code += fmt.Sprintf("%s%d (123, \"foo\");", fnpfx, i) + newline
	}
	code += "return 0;" + newline
	code += "}" + newline

	fg.code = code
	os.Mkdir("tmp", 0755)
	err := ioutil.WriteFile("./tmp/binsize.cpp", []byte(code), 0644)
	if err != nil {
		log.Println(err)
	}
}
func (fg *CppFuncGen) Compile() {
	cmd := exec.Command("g++", "-g", "-O2", "-o", "binsize-cpp", "./tmp/binsize.cpp")
	errcc, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err, string(errcc))
	}
	cmd = exec.Command("ls", "-lh", "binsize-cpp")
	btime := time.Now()
	outcc, err := cmd.CombinedOutput()
	log.Println(string(outcc), time.Since(btime))

	{
		cmd := exec.Command("g++", "-Os", "-o", "binsize-cpp", "./tmp/binsize.cpp")
		errcc, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err, string(errcc))
		}
		cmd = exec.Command("ls", "-lh", "binsize-cpp")
		btime := time.Now()
		outcc, err := cmd.CombinedOutput()
		log.Println(string(outcc), time.Since(btime))

	}
}
func (fg *CppFuncGen) Getresult() {

}

///
type NimFuncGen struct {
	basegen
}

func (fg *NimFuncGen) Init(cnt int) {
	fg.cnt = cnt
}

func (fg *NimFuncGen) Generate() {
	cnt := fg.cnt
	code := "#//" + newline

	code += "var keeper_val = 0" + newline + newline
	code += "proc keeper()=\n  keeper_val+=1" + newline + newline

	for i := 0; i < cnt; i++ {
		// code += fmt.Sprintf("//go:noinline") + newline
		code += fmt.Sprintf("proc %s%d (a0: int, a1: string)=", fnpfx, i) + newline
		code += "  var innerstr: string" + newline
		code += "  var lineno: int" + newline
		code += "  var fromidx: int" + newline
		// code += "  a0 = lineno" + newline
		// code += "  a0 = fromidx" + newline
		// code += "  a1 = innerstr" + newline
		code += "  keeper()" + newline + newline
		// code += "}" + newline
	}

	code += "proc main()=" + newline
	for i := 0; i < cnt; i++ {
		code += fmt.Sprintf("  %s%d(123, \"foo\")", fnpfx, i) + newline
	}
	// code += "return 0;" + newline
	// code += "}" + newline

	code += "main()" + newline

	fg.code = code
	os.Mkdir("tmp", 0755)
	err := ioutil.WriteFile("./tmp/binsize.nim", []byte(code), 0644)
	if err != nil {
		log.Println(err)
	}
}
func (fg *NimFuncGen) Compile() {
	cmd := exec.Command("nim", "c", "-o:binsize-nim", "./tmp/binsize.nim")
	errcc, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err, string(errcc))
	}
	cmd = exec.Command("ls", "-lh", "binsize-nim")
	btime := time.Now()
	outcc, err := cmd.CombinedOutput()
	log.Println(string(outcc), time.Since(btime))

	{
		cmd := exec.Command("nim", "c", "--opt:size",
			"-o:binsize-nim", "./tmp/binsize.nim")
		errcc, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err, string(errcc))
		}
		cmd = exec.Command("ls", "-lh", "binsize-nim")
		btime := time.Now()
		outcc, err := cmd.CombinedOutput()
		log.Println(string(outcc), time.Since(btime))

	}
}
func (fg *NimFuncGen) Getresult() {

}

///
type VFuncGen struct {
	basegen
}

func (fg *VFuncGen) Init(cnt int) {
	fg.cnt = cnt
}

func (fg *VFuncGen) Generate() {
	cnt := fg.cnt
	code := "module main" + newline
	code += "/*\n*/" + newline
	// code += "import \"C\"" + newline
	// code += "var keeper_val = 0" + newline
	code += "fn keeper() {/*keeper_val++*/ }" + newline

	fnpfx2 := strings.ToLower(fnpfx)
	for i := 0; i < cnt; i++ {
		// code += fmt.Sprintf("//go:noinline") + newline
		code += fmt.Sprintf("fn %s%d (a0 int, a1 string) {", fnpfx2, i) + newline
		code += "keeper()" + newline
		code += "}" + newline
	}

	code += "fn main() {" + newline
	for i := 0; i < cnt; i++ {
		code += fmt.Sprintf("%s%d (123, \"foo\")", fnpfx2, i) + newline
	}
	code += "}" + newline

	fg.code = code
	os.Mkdir("tmp", 0755)
	err := ioutil.WriteFile("./tmp/binsize.v", []byte(code), 0644)
	if err != nil {
		log.Println(err)
	}
}
func (fg *VFuncGen) Compile() {
	cmd := exec.Command("v", "build", "-o", "binsize-v", "./tmp/binsize.v")
	errcc, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err, string(errcc))
	}
	cmd = exec.Command("ls", "-lh", "binsize-v")
	btime := time.Now()
	outcc, err := cmd.CombinedOutput()
	log.Println(string(outcc), time.Since(btime))

	if false {
		// cmdstr := `go build -v -p 1 -gcflags "-N -l" -ldflags "-w -s" -o binsize-go tmp/binsize.go`
		cmd := exec.Command("go", "build",
			"-gcflags", "-N -l", "-ldflags", "-w -s",
			"-o", "binsize-go", "./tmp/binsize.go")
		errcc, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err, string(errcc))
		}
		cmd = exec.Command("ls", "-lh", "binsize-go")
		btime := time.Now()
		outcc, err := cmd.CombinedOutput()
		log.Println(string(outcc), time.Since(btime))
	}
}
func (fg *VFuncGen) Getresult() {

}
