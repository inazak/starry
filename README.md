# starry

Starry is a stack-based (esoteric) programming language.
See also [esolangs pages](https://esolangs.org/wiki/Starry).

Starryとは『Rubyで作る奇妙なプログラミング言語』で説明されている、
「きらめく星空のような」スタックベースのオリジナル言語。
使用する文字は `+*.,``'` とスペースのみで、他は無視される。
文字の前のスペースの数で命令の種類を決定する。
詳しくは [esolangs pages](https://esolangs.org/wiki/Starry) を参照。


## Command line

```
Usage: 
  starry [OPTION] SOURCEFILE

  OPTION:
    -i or -inst  ... print decoded instruction code.
```


## Examples

Hello World
```
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
```

Result
```
$ starry.exe helloworld.txt
Hello, world!
```


## Requirements

golang.

