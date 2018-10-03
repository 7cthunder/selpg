# selpg

## selpg 程序逻辑

> selpg 是从文本输入选择页范围的使用程序。该输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入。
>
> selpg 首先处理所有的命令行参数。在扫描了所有的选项参数（也就是那些以连字符为前缀的参数）后，如果 selpg 发现还有一个参数，则它会接受该参数为输入文件的名称并尝试打开它以进行读取。如果没有其它参数，则 selpg 假定输入来自标准输入。



## 参数处理

这里我们实现的是 Golang 版本，所以使用一个 `pflag` 的包就能完成参数处理：

```go
/**
 * init args:
 *
 * "page length": 72(default)
 * "page type": false(default) -> lNumber/page, true -> \f as next page
 * "print dest": ""(default) -> no print dest
 */
func initArgs(args *SelpgArgs) {
	flag.IntVarP(&args.startPage, "sNumber", "s", -1, "start page")
	flag.IntVarP(&args.endPage, "eNumber", "e", -1, "end page")
	flag.IntVarP(&args.pageLen, "lNumber", "l", 72, "lines/page")
	flag.BoolVarP(&args.pageType, "formFeed", "f", false, "form-feed-delimited")
	flag.StringVarP(&args.printDest, "dest", "d", "", "print dest")
	flag.Parse()
}
```

### "-sNumber" 和 "-eNumber"强制选项

selpg 要求用户用两个命令行参数`-s`和`-e`来表示要抽取的页面范围的起始页和结束页；`selpg`程序会简单处理这两个数字的逻辑关系，例如是否大于0，结束页是否大于等于起始页。

### "-lNumber" 和 "-f"可选选项

selpg 可以处理两种输入文本：

*类型 1：*该类文本的页行数固定。这是缺省类型，因此不必给出选项进行说明。也就是说，如果既没有给出“-lNumber”也没有给出“-f”选项，则 selpg 会理解为页有固定的长度（每页 72 行）。

*类型 2：*该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C 中用“\f”表示）定界。

### “-dDestination"可选选项

[原程序实现](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html) 该选项允许用户使用“-dDestination”选项将选定的页直接发送至打印机。这里，“Destination”应该是 lp 命令“-d”选项（请参阅“man lp”）可接受的打印目的地名称。

**但是本人没有打印设备，所以改用 cat 命令来替代它，因为这里主要就是 pipe 处理**



## 使用 selpg

1. `$ selpg -s=1 -e=1 input_file`

   该命令将把“input_file”的第 1 页写至标准输出（也就是屏幕），因为这里没有重定向或管道。

2. `$ selpg -s=1 -e=1 < input_file`

   该命令与示例 1 所做的工作相同，但在本例中，selpg 读取标准输入，而标准输入已被 shell／内核重定向为来自“input_file”而不是显式命名的文件名参数。输入的第 1 页被写至屏幕。

3. `$ other_command | selpg -s=10 -e=20`

   “other_command”的标准输出被 shell／内核重定向至 selpg 的标准输入。将第 10 页到第 20 页写至 selpg 的标准输出（屏幕）。

4. `$ selpg -s=10 -e=20 input_file >output_file`

   selpg 将第 10 页到第 20 页写至标准输出；标准输出被 shell／内核重定向至“output_file”。

5. `$ selpg -s=10 -e=20 input_file 2>error_file`

   selpg 将第 10 页到第 20 页写至标准输出（屏幕）；所有的错误消息被 shell／内核重定向至“error_file”。请注意：在“2”和“>”之间不能有空格；这是 shell 语法的一部分（请参阅“man bash”或“man sh”）。

6. `$ selpg -s=10 -e=20 input_file >output_file 2>error_file`

   selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至“error_file”。当“input_file”很大时可使用这种调用；您不会想坐在那里等着 selpg 完成工作，并且您希望对输出和错误都进行保存。

7. `$ selpg -s=10 -e=20 input_file >output_file 2>/dev/null`

   selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至 /dev/null（空设备），这意味着错误消息被丢弃了。设备文件 /dev/null 废弃所有写至它的输出，当从该设备文件读取时，会立即返回 EOF。

8. `$ selpg -s=10 -e=20 input_file >/dev/null`

   selpg 将第 10 页到第 20 页写至标准输出，标准输出被丢弃；错误消息在屏幕出现。这可作为测试 selpg 的用途，此时您也许只想（对一些测试情况）检查错误消息，而不想看到正常输出。

9. `selpg -s=10 -e=20 input_file | other_command`

   selpg 的标准输出透明地被 shell／内核重定向，成为“other_command”的标准输入，第 10 页到第 20 页被写至该标准输入。“other_command”的示例可以是 lp，它使输出在系统缺省打印机上打印。“other_command”的示例也可以 wc，它会显示选定范围的页中包含的行数、字数和字符数。“other_command”可以是任何其它能从其标准输入读取的命令。错误消息仍在屏幕显示。

10. `$ selpg -s=10 -e=20 input_file 2>error_file | other_command`

    与上面的示例 9 相似，只有一点不同：错误消息被写至“error_file”。

在以上涉及标准输出或标准错误重定向的任一示例中，用“>>”替代“>”将把输出或错误数据附加在目标文件后面，而不是覆盖目标文件（当目标文件存在时）或创建目标文件（当目标文件不存在时）。

11. `$ selpg -s=10 -e=20 -l=66 input_file`

    该命令将页长设置为 66 行，这样 selpg 就可以把输入当作被定界为该长度的页那样处理。第 10 页到第 20 页被写至 selpg 的标准输出（屏幕）。

12. `$ selpg -s=10 -e=20 -f input_file`

    假定页由换页符定界。第 10 页到第 20 页被写至 selpg 的标准输出（屏幕）。

13. `$ selpg -s=10 -e=20 -d=lp1 input_file`

    第 10 页到第 20 页由管道输送至命令“lp -dlp1”，该命令将使输出在打印机 lp1 上打印。 **我们这里是无论什么打印设备，都是使用`cat -n`命令在标准输出中打印内容**


## 测试

使用测试文件`input_file_line` 和 `input_file_f`来代表两种文本，其中 `input_file_f` 随机生成 `'\f'` 。

1. `$ selpg -s=1 -e=1 input_file_line`

   屏幕输出为：

   ```bash
   Line   0
   Line   1
   Line   2
   Line   3
   Line   4
   Line   5
   Line   6
   Line   7
   Line   8
   Line   9
   Line  10
   ...
   Line  69
   Line  70
   Line  71
   ```

2. `$ selpg -s=1 -e=1 < input_file_line`

   屏幕输出为：

   ```bash
   Line   0
   Line   1
   Line   2
   Line   3
   Line   4
   Line   5
   Line   6
   Line   7
   Line   8
   Line   9
   Line  10
   ...
   Line  69
   Line  70
   Line  71
   ```

3. `$ cat -n input_file_line | selpg -s=1 -e=2`

   屏幕输出为：

   ```bash
        1  Line   0
        2  Line   1
        3  Line   2
        4  Line   3
        5  Line   4
        6  Line   5
        7  Line   6
        8  Line   7
        9  Line   8
       10  Line   9
   	...
      142  Line 141
      143  Line 142
      144  Line 143
   ```

4. `$ selpg -s=1 -e=2 input_file_line >output_file`

   查看`output_file`：

   ```bash
        1  Line   0
        2  Line   1
        3  Line   2
        4  Line   3
        5  Line   4
        6  Line   5
        7  Line   6
        8  Line   7
        9  Line   8
       10  Line   9
       ...
      142  Line 141
      143  Line 142
      144  Line 143
   ```

5. `$ selpg -s=1 -e=2 input_file 2>error_file`

   查看 `error_file`，因为我们的输入文件没有 `input_file`:

   ```bash
   selpg: open input_file: no such file or directory
   ```

6. `$ selpg -s=1 -e=2 input_file_line >output_file 2>error_file`

   查看`output_file`：

   ```bash
        1  Line   0
        2  Line   1
        3  Line   2
        4  Line   3
        5  Line   4
        6  Line   5
        7  Line   6
        8  Line   7
        9  Line   8
       10  Line   9
       ...
      142  Line 141
      143  Line 142
      144  Line 143
   ```

   查看 `error_file`:

   ```bash
   
   ```

7. `$ selpg -s=1 -e=2 input_file_line | cat -n`

   屏幕输出为：

   ```bash
        1  Line   0
        2  Line   1
        3  Line   2
        4  Line   3
        5  Line   4
        6  Line   5
        7  Line   6
        8  Line   7
        9  Line   8
       10  Line   9
       ...
      142  Line 141
      143  Line 142
      144  Line 143
   ```

8. `selpg -s=1 -e=2 -l=66 input_file_line`

   屏幕输出为：

   ```bash
   Line   0
   Line   1
   Line   2
   Line   3
   Line   4
   Line   5
   Line   6
   Line   7
   Line   8
   Line   9
   Line  10
   ...
   Line  69
   Line  70
   Line  131
   ```

9. `selpg -s=1 -e=2 -f input_file_f`

   屏幕输出为：

   ```bash
   Line   0
   Line   1
   Line   2
   Line   3
   Line   4
   Line   5
   Line   6
   Line   7
   Line   8
   Line   9
   ```

   `input_file_f` 前23行，`^L`在vim中表示的就是`\f`：

   ```bash
     1 Line   0
     2 Line   1
     3 ^LLine   2
     4 Line   3
     5 Line   4
     6 Line   5
     7 Line   6
     8 Line   7
     9 Line   8
    10 Line   9
    11 ^LLine  10
    12 Line  11
    13 Line  12
    14 ^LLine  13
    15 ^LLine  14
    16 Line  15
    17 Line  16
    18 Line  17
    19 Line  18
    20 Line  19
    21 Line  20
    22 Line  21
    23 Line  22
   
   ```

10. `selpg -s=1 -e=2 -d=lp1 input_file_line`

    查看屏幕输出：

    ```bash
         1  Line   0
         2  Line   1
         3  Line   2
         4  Line   3
         5  Line   4
         6  Line   5
         7  Line   6
         8  Line   7
         9  Line   8
        10  Line   9
        ...
       142  Line 141
       143  Line 142
       144  Line 143
    ```