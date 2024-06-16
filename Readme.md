# cenc

テキストファイルの文字コードを変更するソフトウェアです  
対応している文字コードは２つです  

utf-8  
shift-jis  

# 使い方

```
cenc.exe [file code] [out code] [file path]
file code
  ファイルの文字コード [auto | utf-8 | shift-jis]

out code
  変換先の文字コード [utf-8 | shift-jis]

file path
  変換するテキストファイルのパス
```

```
example
  cenc.exe shift-jis utf-8 test.text
```
