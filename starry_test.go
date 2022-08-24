package starry

import (
  "testing"
  "bytes"
)

func TestParser(t *testing.T) {

  ps := []struct {
    Text string
    Inst string
  }{
    { Text: " +",
      Inst: "dup", },
    { Text: "  +",
      Inst: "swap", },
    { Text: "   +",
      Inst: "rotate", },
    { Text: "    +",
      Inst: "pop", },
    { Text: "     +",
      Inst: "push 0", },
    { Text: "      +",
      Inst: "push 1", },
    { Text: "*",
      Inst: "+", },
    { Text: " *",
      Inst: "-", },
    { Text: "  *",
      Inst: "*", },
    { Text: "   *",
      Inst: "/", },
    { Text: "    *",
      Inst: "%", },
    { Text: "     *",
      Inst: "+", },
    { Text: ".",
      Inst: "output number", },
    { Text: " .",
      Inst: "output character", },
    { Text: "  .",
      Inst: "output number", },
    { Text: ",",
      Inst: "input number", },
    { Text: " ,",
      Inst: "input character", },
    { Text: "  ,",
      Inst: "input number", },
    { Text: "`",
      Inst: "label 0", },
    { Text: " `",
      Inst: "label 1", },
    { Text: "'",
      Inst: "jumpnz 0", },
    { Text: " '",
      Inst: "jumpnz 1", },
  }

  for i, p := range ps {
    parser  := NewParser(p.Text)
    list, _ := parser.Parse()

    if len(parser.Error) != 0 {
      t.Errorf("No.%d Parser Error: %s", i, parser.Error)
    }

    if list[0].Decode() != p.Inst {
      t.Errorf("No.%d expect=%v got=%v", i, p.Inst, list[0].Decode())
    }
  }
}

func TestStackOp(t *testing.T) {
  vm  := NewVM([]Instruction{
                 Push{ Value: 88 },
                 Push{ Value: 2 },
                 Div{}, // 88/2=44
                 Push{ Value: 33 },
                 Minus{}, // 44-33=11
                 Push{ Value: 6 },
                 Mod{}, // 11%6=5
                 Dup{}, // 5 | 5
                 Push{ Value: 3 },
                 Mul{}, // 5 | 5*3=15
                 Swap{}, // 15 | 5
                 Push{ Value: 8 }, // 15 | 5 | 8
                 Rotate{}, // 8 | 15 | 5
                 Mul{}, // 8 | 15*5=75
                 Plus{}, // 8+75=83
                 OutputC{}, // 'S'
               }, nil)


  out := &bytes.Buffer{}
  vm.SetStdout(out)

  result, _ := vm.Run()

  if result != 0 {
    t.Fatalf("VM return code expect=%d, got=%d", 0, result)
  }

  expect := "S"
  if out.String() != expect {
    t.Fatalf("Stdout expect=`%s` got=`%s`", expect, out.String())
  }
}

func TestJumpOp(t *testing.T) {

  label := make(map[int]int)
  label[9] = 1
  vm  := NewVM([]Instruction{
                 Push{ Value: 4 },
                 Label{ Value: 9 },
                 Push{ Value: 46 },
                 OutputC{}, // print '.'
                 Push{ Value: 1 },
                 Minus{},
                 Dup{}, // because jumpnz op use pop
                 JumpNZ{ Value: 9 },
                 Push{ Value: 33 },
                 OutputC{}, // print '!'
               },
               label)

  out := &bytes.Buffer{}
  vm.SetStdout(out)

  result, _ := vm.Run()

  if result != 0 {
    t.Fatalf("VM return code expect=%d, got=%d", 0, result)
  }

  expect := "....!"
  if out.String() != expect {
    t.Fatalf("Stdout expect=`%s` got=`%s`", expect, out.String())
  }
}

func TestNumberInputOp(t *testing.T) {
  vm  := NewVM([]Instruction{
                 InputN{},
                 OutputN{},
               }, nil)

  testdata := "999"
  in  := bytes.NewBufferString(testdata)
  out := &bytes.Buffer{}
  vm.SetStdin(in)
  vm.SetStdout(out)

  result, _ := vm.Run()

  if result != 0 {
    t.Fatalf("VM return code expect=%d, got=%d", 0, result)
  }

  if out.String() != testdata {
    t.Fatalf("Stdout expect=`%s` got=`%s`", testdata, out.String())
  }
}

func TestCharInputOp(t *testing.T) {
  vm  := NewVM([]Instruction{
                 InputC{},
                 OutputC{},
               }, nil)

  testdata := "Z"
  in  := bytes.NewBufferString(testdata)
  out := &bytes.Buffer{}
  vm.SetStdin(in)
  vm.SetStdout(out)

  result, _ := vm.Run()

  if result != 0 {
    t.Fatalf("VM return code expect=%d, got=%d", 0, result)
  }

  if out.String() != testdata {
    t.Fatalf("Stdout expect=`%s` got=`%s`", testdata, out.String())
  }
}

func TestHelloWorld(t *testing.T) {

  text :=`
            +               +  *       +    
 * + .        +              +  *       +   
  *     * + .            +     * + . + .    
    +     * + .              +            + 
 *         +     * * + .                 + *
 + .              + +  *           +     *  
   * + .             + * + .        +     * 
+ .           + * + .             + * + .   
           +            +  *         +     *
 * + .
`

  parser := NewParser(text)
  inst, label := parser.Parse()

  if len(parser.Error) != 0 {
    t.Fatalf("Parser Error: %v", parser.Error)
  }

  buf := &bytes.Buffer{}
  vm  := NewVM(inst, label)
  vm.SetStdout(buf)

  result, _ := vm.Run()

  if result != 0 {
    t.Fatalf("VM return code expect=%d, got=%d", 0, result)
  }

  expect := "Hello, world!"
  if buf.String() != expect {
    t.Fatalf("Stdout expect=`%s` got=`%s`", expect, buf.String())
  }
}

