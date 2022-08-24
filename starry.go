package starry

import (
  "fmt"
  "io"
  "os"
)

type Instruction interface {
  Decode() string
}

type Nop     struct {}
type Push    struct { Value int }
type Pop     struct {}
type Dup     struct {}
type Swap    struct {}
type Rotate  struct {}
type Plus    struct {}
type Minus   struct {}
type Mul     struct {}
type Div     struct {}
type Mod     struct {}
type Label   struct { Value int }
type JumpNZ  struct { Value int }
type OutputN struct {}
type OutputC struct {}
type InputN  struct {}
type InputC  struct {}

func (i Nop)     Decode() string { return "nop" }
func (i Push)    Decode() string { return fmt.Sprintf("push %d", i.Value) }
func (i Pop)     Decode() string { return "pop" }
func (i Dup)     Decode() string { return "dup" }
func (i Swap)    Decode() string { return "swap" }
func (i Rotate)  Decode() string { return "rotate" }
func (i Plus)    Decode() string { return "+" }
func (i Minus)   Decode() string { return "-" }
func (i Mul)     Decode() string { return "*" }
func (i Div)     Decode() string { return "/" }
func (i Mod)     Decode() string { return "%" }
func (i Label)   Decode() string { return fmt.Sprintf("label %d", i.Value) }
func (i JumpNZ)  Decode() string { return fmt.Sprintf("jumpnz %d", i.Value) }
func (i OutputN) Decode() string { return "output number" }
func (i OutputC) Decode() string { return "output character" }
func (i InputN)  Decode() string { return "input number" }
func (i InputC)  Decode() string { return "input character" }

type VMError struct {
  m string
  p int
}

func (e VMError) Error() string {
  return fmt.Sprintf("pc=%d %s", e.p, e.m)
}

type VM struct {
  Stack  []int
  PC     int
  Inst   []Instruction
  Label  map[int]int
  Stdin  io.Reader
  Stdout io.Writer
}

func NewVM(inst []Instruction, label map[int]int) *VM {
  return &VM{
    Stack:  []int{},
    PC:     0,
    Inst:   inst,
    Label:  label,
    Stdin:  os.Stdin,
    Stdout: os.Stdout,
  }
}

func (v *VM) SetStdin(r io.Reader) {
  v.Stdin = r
}

func (v *VM) SetStdout(w io.Writer) {
  v.Stdout = w
}

func (v *VM) Run() (int, error) {

  for len(v.Inst) > v.PC {
    err := v.step()
    if err != nil {
      return 1, err
    }
  }

  return 0, nil
}

func (v *VM) RunWithDebug() (int, error) {

  for len(v.Inst) > v.PC {
    err := v.step()
    if err != nil {
      return 1, err
    }
    fmt.Fprintf(os.Stderr, "[Debug] Stack=%v PC=%v\n", v.Stack, v.PC)
  }

  return 0, nil
}

func (v *VM) push(value int) {
  v.Stack = append(v.Stack, value)
}

func (v *VM) pop() int {
  if len(v.Stack) < 1 {
    panic(VMError{m:"insufficient stack size", p:v.PC})
  }
  top    := v.Stack[len(v.Stack)-1]
  v.Stack = v.Stack[:len(v.Stack)-1]
  return top
}

func (v *VM) dup() {
  top := v.pop()
  v.push(top)
  v.push(top)
}

func (v *VM) swap() {
  x := v.pop()
  y := v.pop()
  v.push(x)
  v.push(y)
}

func (v *VM) rotate() {
  x := v.pop()
  y := v.pop()
  z := v.pop()
  v.push(x)
  v.push(z)
  v.push(y)
}

func (v *VM) plus() {
  x := v.pop()
  y := v.pop()
  v.push(y + x)
}

func (v *VM) minus() {
  x := v.pop()
  y := v.pop()
  v.push(y - x)
}

func (v *VM) mul() {
  x := v.pop()
  y := v.pop()
  v.push(y * x)
}

func (v *VM) div() {
  x := v.pop()
  y := v.pop()
  if x == 0 {
    panic(VMError{m:"division by zero", p:v.PC})
  }
  v.push(y / x)
}

func (v *VM) mod() {
  x := v.pop()
  y := v.pop()
  if x == 0 {
    panic(VMError{m:"division by zero", p:v.PC})
  }
  v.push(y % x)
}

func (v *VM) jumpnz(addr int) {
  if _, ok := v.Label[addr]; !ok {
    panic(VMError{m:"jump target is not found", p:v.PC})
  }
  top := v.pop()
  if top != 0 {
    v.PC = v.Label[addr]
  } else {
    v.PC += 1
  }
}

func (v *VM) step() (err error) {

  // capture VMError from panic
  defer func() {
    if rec := recover(); rec != nil {
      if _, ok := rec.(VMError); ok {
        err = rec.(VMError)
      } else {
        panic(rec)
      }
    }
  }()

  inst := v.Inst[v.PC]

  switch inst.(type) {
    case Nop:
      //do nothing
    case Push:
      push := inst.(Push)
      v.push(push.Value)
    case Pop:
      v.pop()
    case Dup:
      v.dup()
    case Swap:
      v.swap()
    case Rotate:
      v.rotate()
    case Plus:
      v.plus()
    case Minus:
      v.minus()
    case Mul:
      v.mul()
    case Div:
      v.div()
    case Mod:
      v.mod()
    case Label:
      //do nothing
    case JumpNZ:
      jumpnz := inst.(JumpNZ)
      v.jumpnz(jumpnz.Value)
      return nil // v.PC is updated
    case OutputN:
      top := v.pop()
      fmt.Fprintf(v.Stdout, "%d", top)
    case OutputC:
      top := v.pop()
      fmt.Fprintf(v.Stdout, "%c", top)
    case InputN:
      var i int
      fmt.Fscanf(v.Stdin, "%d", &i)
      v.push(i)
    case InputC:
      var i int
      fmt.Fscanf(v.Stdin, "%c", &i)
      v.push(i)
    default:
      panic(VMError{m:"unknown instruction", p:v.PC})
  }

  v.PC += 1
  return nil
}


type Parser struct {
  Src   string
  Line  int
  Error []string
}

func NewParser(src string) *Parser {
  return &Parser{
    Src:   src,
    Line:  0,
    Error: []string{},
  }
}

func (p *Parser) AddError(m string) {
  msg := fmt.Sprintf("line %d - %s", p.Line, m)
  p.Error = append(p.Error, msg)
}

func (p *Parser) Parse() (list []Instruction, label map[int]int) {
  list  = []Instruction{}
  label = make(map[int]int)

  spaces := 0

  for _, c := range p.Src {
    switch c {
      case ' ':
        spaces += 1
      case '+':
        if spaces == 0 {
          p.AddError("zero space on '+'")
        } else {
          inst := getInstForStack(spaces)
          list = append(list, inst)
        }
        spaces = 0
      case '*':
        inst := getInstForCalc(spaces)
        list = append(list, inst)
        spaces = 0
      case '.':
        inst := getInstForOutput(spaces)
        list = append(list, inst)
        spaces = 0
      case ',':
        inst := getInstForInput(spaces)
        list = append(list, inst)
        spaces = 0
      case '`':
        if _, ok := label[spaces] ; ok {
          p.AddError("duplicated label")
        } else {
          label[spaces] = len(list)
          inst := Label{ Value: spaces }
          list = append(list, inst)
        }
        spaces = 0
      case '\'':
        inst := JumpNZ{ Value: spaces }
        list = append(list, inst)
        spaces = 0
      case '\n':
        p.Line += 1
      default:
        //ignored
    }
  } // for range loop

  return list, label
}


func getInstForStack(spaces int) Instruction {
  var inst Instruction
  switch spaces {
    case 0:
      inst = Nop{} // but nothing reach this case
    case 1:
      inst = Dup{}
    case 2:
      inst = Swap{}
    case 3:
      inst = Rotate{}
    case 4:
      inst = Pop{}
    default:
      inst = Push{ Value: spaces - 5 }
  }
  return inst
}

func getInstForCalc(spaces int) Instruction {
  var inst Instruction
  switch spaces % 5 {
    case 0:
      inst = Plus{}
    case 1:
      inst = Minus{}
    case 2:
      inst = Mul{}
    case 3:
      inst = Div{}
    case 4:
      inst = Mod{}
  }
  return inst
}

func getInstForOutput(spaces int) Instruction {
  var inst Instruction
  switch spaces % 2 {
    case 0:
      inst = OutputN{}
    case 1:
      inst = OutputC{}
  }
  return inst
}

func getInstForInput(spaces int) Instruction {
  var inst Instruction
  switch spaces % 2 {
    case 0:
      inst = InputN{}
    case 1:
      inst = InputC{}
  }
  return inst
}

