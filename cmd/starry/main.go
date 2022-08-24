package main

import (
  "flag"
  "fmt"
  "os"
  "io/ioutil"
  "github.com/inazak/starry"
)

var usage=`
starry is a stack-based esoteric programming language. 

Usage: 
  starry [OPTION] SOURCEFILE

  OPTION:
    -i or -inst  ... print instruction code.
    -d or -debug ... run with debug print.
`

var optionsInst  bool
var optionsDebug bool

func main() {

  flag.BoolVar(&optionsInst,  "inst",  false, "print decoded instruction code.")
  flag.BoolVar(&optionsInst,  "i",     false, "print decoded instruction code.")
  flag.BoolVar(&optionsDebug, "debug", false, "run with debug print.")
  flag.BoolVar(&optionsDebug, "d",     false, "run with debug print.")
  flag.Parse()

  if len(flag.Args()) != 1 {
    fmt.Printf("%s", usage)
    os.Exit(1)
  }
  filename := flag.Args()[0]

  text, err := readTextFile(filename)
  if err != nil {
    fmt.Printf("File Open Error: %s\n", filename)
    os.Exit(1)
  }

  parser := starry.NewParser(text)
  insts, label := parser.Parse()

  if len(parser.Error) != 0 {
    for _, msg := range parser.Error {
      fmt.Printf("Parser Error: %s\n", msg)
    }
    os.Exit(1)
  }

  if optionsInst {
    for i, inst := range insts {
      fmt.Printf("[%03d] %s\n", i, inst.Decode())
    }
    os.Exit(0)
  }

  vm := starry.NewVM(insts, label)

  var result int

  if optionsDebug {
    result, err = vm.RunWithDebug()
  } else {
    result, err = vm.Run()
  }

  if result != 0 {
    fmt.Printf("VM Runtime Error: %v\n", err)
    os.Exit(1)
  }

  os.Exit(0)
}


func readTextFile(filepath string) (string, error) {

  file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
  if err != nil {
    return "", err
  }
  defer file.Close()

  bytes, err := ioutil.ReadAll(file)
  if err != nil {
    return "", err
  }

  return string(bytes), nil
}

