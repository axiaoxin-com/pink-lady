# gettext

```
Original C Sources ───> Preparation ───> Marked C Sources ───╮
                                                             │
              ╭─────────<─── GNU gettext Library             │
╭─── make <───┤                                              │
│             ╰─────────<────────────────────┬───────────────╯
│                                            │
│   ╭─────<─── PACKAGE.pot <─── xgettext <───╯   ╭───<─── PO Compendium
│   │                                            │              ↑
│   │                                            ╰───╮          │
│   ╰───╮                                            ├───> PO editor ───╮
│       ├────> msgmerge ──────> LANG.po ────>────────╯                  │
│   ╭───╯                                                               │
│   │                                                                   │
│   ╰─────────────<───────────────╮                                     │
│                                 ├─── New LANG.po <────────────────────╯
│   ╭─── LANG.gmo <─── msgfmt <───╯
│   │
│   ╰───> install ───> /.../LANG/PACKAGE.mo ───╮
│                                              ├───> "Hello world!"
╰───────> install ───> /.../bin/PROGRAM ───────╯
```

![png](https://www.atjiang.com/assets/images/gnu-gettext.png)

- xgettext工具从源文件中提取待翻译语句，生成pot模板文件
- 再用msginit从.pot文件生成特定语言版本的.po文件  `msginit --input=messages.pot --local=zh_CN.po`
- 使用msgfmt工具生成mo文件 `msgfmt zh_CN.po -o zh_CN.mo`
- 当源文件更新后，生成新的.pot文件，使用msgmerge工具将新的.pot文件和已翻译的po文件合并成新的.po文件 `msgmerge messages.pot zh_CN.po -o zh_CN2.po`
