#!/usr/bin/env python3
"""
PO文件翻译脚本
用于自动翻译未翻译的PO文件内容
"""

import re
import sys
from typing import Dict, List, Generator, Optional
from pathlib import Path

import polib
from pygtrans import Translate
import markdown2

# 配置常量
SITE_NAME_MAPPING = {
    "raw": "pink-lady",
    "translated": "pink-lady"
}

MARKDOWN_PATTERNS = [
    r'^#\s',           # 标题
    r'^>\s',           # 引用
    r'^\s*[-*]\s',     # 列表
    r'```',            # 代码块
    r'\[.*?\]\(.*?\)', # 链接
    r'\*\*.*?\*\*',    # 粗体
    r'_.*?_',          # 斜体
]

MARKDOWN_EXTRAS = [
    'fenced-code-blocks',    # 代码块支持
    'tables',                # 表格支持
    'header-ids',            # 标题ID
    'target-blank-links',    # 链接新窗口打开
    'code-friendly',         # 代码友好
    'preserve_whitespace',   # 保留空白
    'footnotes',             # 脚注
    'smarty-pants',          # 智能标点
    'wiki-tables',           # 更多表格格式支持
    'xml',                   # XML 支持
    'def_list',              # 定义列表
    'markdown-in-html',      # HTML中的markdown支持
    'strike',                # 删除线支持
    'task_list',             # GFM任务列表
    'mermaid',               # mermaid图表支持
]

class TranslationManager:
    def __init__(self, target_lang: str):
        self.target_lang = target_lang
        self.translator = Translate()
        self.po_file_path = Path(f"../../statics/i18n/{target_lang}/LC_MESSAGES/messages.po")

    def is_markdown_content(self, text: str) -> bool:
        """检查文本是否包含markdown语法"""
        return any(re.search(pattern, text, re.MULTILINE) for pattern in MARKDOWN_PATTERNS)

    def convert_to_html(self, text: str) -> str:
        """将markdown转换为HTML"""
        if not self.is_markdown_content(text):
            return text

        html = markdown2.markdown(text, extras=MARKDOWN_EXTRAS)
        return html

    @staticmethod
    def chunk_list(items: List, chunk_size: int) -> Generator[List, None, None]:
        """将列表分割成指定大小的块"""
        for i in range(0, len(items), chunk_size):
            yield items[i:i + chunk_size]

    def get_untranslated_entries(self, po_file: polib.POFile) -> Dict[str, str]:
        """获取需要翻译的条目"""
        untranslated = {}

        for entry in po_file:
            if entry.msgid == "%s":
                continue

            processed_msgid = entry.msgid.replace(
                SITE_NAME_MAPPING["raw"],
                SITE_NAME_MAPPING["translated"]
            )
            processed_msgid = self.convert_to_html(processed_msgid)

            if (entry.msgid and not entry.msgstr) or \
               entry.msgid == entry.msgstr or \
               (entry.fuzzy and not entry.obsolete):
                untranslated[processed_msgid] = entry.msgid

        return untranslated

    def handle_code_blocks(self, original: str, translated_text: str) -> str:
        """处理包含代码块的翻译文本
        翻译后html代码被压缩为一行，代码块的缩进换行等格式丢失，需要尽最大可能的还原

        Args:
            original: 原始文本
            translated_text: 翻译后的文本

        Returns:
            处理后的翻译文本
        """
        # 匹配所有代码块，包括<pre><code>和<pre>的情况
        original_blocks = list(re.finditer(
            r'(<pre[^>]*>)(?:<code[^>]*>)?(.*?)(?:</code>)?(</pre>)',
            original,
            re.DOTALL
        ))

        # 获取翻译后文本中的所有代码块
        translated_blocks = list(re.finditer(
            r'(<pre[^>]*>)(?:<code[^>]*>)?(.*?)(?:</code>)?(</pre>)',
            translated_text,
            re.DOTALL
        ))

        # 确保代码块数量匹配
        if len(original_blocks) == len(translated_blocks):
            # 从后向前替换，避免位置偏移问题
            for i in range(len(original_blocks) - 1, -1, -1):
                original_block = original_blocks[i]
                translated_block = translated_blocks[i]

                # 获取原始代码块的各个部分
                original_pre_open = original_block.group(1)  # <pre> 开始标签
                original_content = original_block.group(2)    # 代码内容
                original_pre_close = original_block.group(3)  # </pre> 结束标签

                # 获取翻译后代码块的各个部分
                translated_pre_open = translated_block.group(1)
                translated_pre_close = translated_block.group(3)

                # 检查原始代码块是否包含 <code> 标签
                has_code_tag = "<code>" in original_block.group(0)
                # 代码块内出现连续\n时，写入msgstr出现多个\n，这时msgstr的html被当作markdown解析会导致html被切断，出现嵌套的pre，导致代码块错乱，用&nbsp;\n代替
                original_content = original_content.replace('\n', '&nbsp;\n')
                # 构建新的代码块
                if has_code_tag:
                    new_block = f"{translated_pre_open}<code>{original_content}</code>{translated_pre_close}"
                else:
                    new_block = f"{translated_pre_open}{original_content}{translated_pre_close}"

                # 替换翻译后文本中的代码块
                translated_text = (
                    translated_text[:translated_block.start()] +
                    new_block +
                    translated_text[translated_block.end():]
                )
        else:
            print(f"Warning: Code block count mismatch in translation. Original: {len(original_blocks)}, Translated: {len(translated_blocks)}")

        return translated_text

    def translate_entries(self, untranslated: Dict[str, str]) -> Dict[str, str]:
        """翻译条目"""
        translations = {}

        for chunk in self.chunk_list(list(untranslated.keys()), 4000):
            translated = self.translator.translate(chunk, target=self.target_lang)
            for original, translation in zip(chunk, translated):
                if translation.translatedText:
                    translated_text = translation.translatedText
                    # 如果原始文本包含pre标签，保留代码格式
                    if "<pre>" in original:
                        translated_text = self.handle_code_blocks(original, translated_text)
                        # print("original", original, "translated", translated_text)

                    translations[untranslated[original]] = translated_text
        return translations

    def update_po_file(self, po_file: polib.POFile, translations: Dict[str, str]) -> None:
        """更新PO文件中的翻译"""
        for entry in po_file:
            if entry.msgid == "%s":
                continue

            if (entry.msgid and not entry.msgstr) or \
               entry.fuzzy or \
               (entry.msgid == entry.msgstr):

                if entry.msgid in translations:
                    entry.msgstr = translations[entry.msgid]
                    if entry.msgid == SITE_NAME_MAPPING["raw"]:
                        entry.msgstr = SITE_NAME_MAPPING["translated"]
                    entry.flags = [f for f in entry.flags if f != "fuzzy"]

    def process(self) -> int:
        """处理翻译流程"""
        try:
            po_file = polib.pofile(str(self.po_file_path), wrapwidth=0)
            print(f"Processing translations for {self.target_lang}...")

            untranslated = self.get_untranslated_entries(po_file)
            if not untranslated:
                print(f"No new translations needed for {self.target_lang}")
                return 0

            translations = self.translate_entries(untranslated)
            self.update_po_file(po_file, translations)

            po_file.save(str(self.po_file_path))
            print(f"Successfully translated {len(translations)} messages for {self.target_lang}")
            return 0

        except Exception as e:
            print(f"Error processing {self.target_lang}: {str(e)}", file=sys.stderr)
            return 1

def main() -> int:
    """主函数"""
    target_lang = sys.argv[1] if len(sys.argv) > 1 else "en"
    manager = TranslationManager(target_lang)
    return manager.process()

if __name__ == "__main__":
    sys.exit(main())
