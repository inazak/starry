# starry

Starry is a stack-based (esoteric) programming language "like a twinkling starry sky", 
as described in "The Weird Programming Language in Ruby".

The only characters used are `` +*.,`' `` and spaces, the rest are ignored.
The number of spaces before a character determines the type of instruction.
See also [esolangs pages](https://esolangs.org/wiki/Starry).

Starryとは『Rubyで作る奇妙なプログラミング言語』で説明されている、
「きらめく星空のような」スタックベースのオリジナル言語。
使用する文字は `` +*.,`' `` とスペースのみで、他は無視される。
文字の前のスペースの数で命令の種類を決定する。
詳しくは [esolangs pages](https://esolangs.org/wiki/Starry) を参照。


## Installation

windows binary is [here](https://github.com/inazak/starry/releases) or


## How to use

```
Usage: 
  starry [OPTION] SOURCEFILE

  OPTION:
    -i or -inst  ... print instruction code.
    -d or -debug ... run with debug print.
```


## Examples

Print Hello World.
```
$ cat helloworld.txt
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
$
$ starry.exe helloworld.txt
Hello, world!
$
```

Print up to 10th Fibonacci number.
```
$ cat fibonacci.txt
      +     +               +       `
   + +   +     * +  .               +
 .   +   +      + * +       '
$
$ starry.exe fibonacci.txt
1
1
2
3
5
8
13
21
34
55
$
```

print instruction code.
```
$ starry.exe -i fibonacci.txt
[000] push 1
[001] push 0
[002] push 10
[003] label 7
[004] rotate
[005] dup
[006] rotate
[007] +
[008] dup
[009] output number
[010] push 10
[011] output character
[012] rotate
[013] rotate
[014] push 1
[015] -
[016] dup
[017] jumpnz 7
```

