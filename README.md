# DocGo
GoDoc html 文档反向到 go 文件命令行工具.

无参数, 直接在当前目录生成子目录, 固定转换 http://mygodoc.oschina.mopaas.com 的中文翻译文档. 缺少合并英文注释功能. 已知转换过程中的问题:


转换过程中产生的未知问题

    http://mygodoc.oschina.mopaas.com/pkg/text_scanner.htm
    46:20: expected 'IDENT', found '|' (and 3 more errors)

html格式不标准

    http://mygodoc.oschina.mopaas.com/pkg/crypto_dsa.htm#GenerateKey
    http://mygodoc.oschina.mopaas.com/pkg/crypto_ecdsa.htm#PrivateKey

下列没翻译

    http://mygodoc.oschina.mopaas.com/pkg/debug_dwarf.htm
    http://mygodoc.oschina.mopaas.com/pkg/debug_elf.htm
    http://mygodoc.oschina.mopaas.com/pkg/debug_gosym.htm
    http://mygodoc.oschina.mopaas.com/pkg/debug_macho.htm
    http://mygodoc.oschina.mopaas.com/pkg/debug_pe.htm
    http://mygodoc.oschina.mopaas.com/pkg/debug_plan9obj.htm
    http://mygodoc.oschina.mopaas.com/pkg/go_ast.htm
    http://mygodoc.oschina.mopaas.com/pkg/go_build.htm
    http://mygodoc.oschina.mopaas.com/pkg/go_scanner.htm
    http://mygodoc.oschina.mopaas.com/pkg/go_token.htm
    http://mygodoc.oschina.mopaas.com/pkg/regexp_syntax.htm
    http://mygodoc.oschina.mopaas.com/pkg/syscall.htm
    http://mygodoc.oschina.mopaas.com/pkg/testing.htm
    http://mygodoc.oschina.mopaas.com/pkg/testing_iotest.htm
    http://mygodoc.oschina.mopaas.com/pkg/testing_quick.htm
    http://mygodoc.oschina.mopaas.com/pkg/text_template_parse.htm

