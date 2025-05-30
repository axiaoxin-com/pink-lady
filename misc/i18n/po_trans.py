#! /usr/bin/env python3
#
# ENV Install:
# virtualenv venv --python=python3.8
# source ./venv/bin/activate
# ./venv/bin/pip-3.8 install -r requirements.txt

import re
import sys

import markdown2
import polib
from pygtrans import Translate

CN_SITENAME = "pink-lady"
EN_SITENAME = "pink-lady"

transcli = Translate()

def is_markdown_content(text: str) -> bool:
    """Check if the text contains markdown syntax."""
    markdown_patterns = [
        r'^#\s',           # 标题
        r'^>\s',           # 引用
        r'^\s*[-*]\s',     # 列表
        r'```',            # 代码块
        r'\[.*?\]\(.*?\)', # 链接
        r'\*\*.*?\*\*',    # 粗体
        r'_.*?_',          # 斜体
    ]
    return any(re.search(pattern, text, re.MULTILINE) for pattern in markdown_patterns)

def convert_to_html(text: str) -> str:
    """Convert markdown to HTML if the text contains markdown syntax."""
    if is_markdown_content(text):
        return markdown2.markdown(
            text,
            extras=[
                'fenced-code-blocks',
                'tables',
                'header-ids',
                'target-blank-links',
            ]
        )
    return text

def chunks(lst, n):
    """Yield successive n-sized chunks from lst."""
    for i in range(0, len(lst), n):
        yield lst[i : i + n]

def main():
    target = "en"
    isinit = False
    if sys.argv[1]:
        target = sys.argv[1]
    if sys.argv[2] == "init":
        isinit = True

    # load po file
    pofile = polib.pofile("../../statics/i18n/" + target + "/LC_MESSAGES/messages.po", wrapwidth=0)
    print("Auto trans po file by pygtrans:", target)
    msgids = dict()
    for entry in pofile:
        if entry.msgid == "%s":
            continue

        # 特殊中文品牌词的指定英文名称替换，避免翻译为非特定的英文品牌词
        _m = entry.msgid.replace(CN_SITENAME, EN_SITENAME)

        # 如果是markdown内容，转换为HTML
        _m = convert_to_html(_m)

        if entry.msgid and not entry.msgstr:
            msgids[_m] = entry.msgid
        if isinit and entry.msgid == entry.msgstr:
            msgids[_m] = entry.msgid
        if entry.fuzzy and not entry.obsolete:
            msgids[_m] = entry.msgid

    print("msgids count:", len(msgids))
    if len(msgids) == 0:
        return

    msgs = dict()
    for rmids in chunks(list(msgids.keys()), 4000):
        trans = transcli.translate(rmids, target=target)
        for idx, rmid in enumerate(rmids):
            if trans[idx].translatedText:
                mid = msgids[rmid]
                msgs[mid] = trans[idx].translatedText
    print("msgs count:", len(msgs))

    for entry in pofile:
        if (
            (entry.msgid and not entry.msgstr)
            or entry.fuzzy
            or (isinit and entry.msgid == entry.msgstr)
        ):
            if entry.msgid == "%s":
                continue
            msgstr = msgs.get(entry.msgid)
            if not msgstr:
                continue

            entry.msgstr = msgstr

            # 英文品牌词的再次确认覆盖
            if entry.msgid == CN_SITENAME:
                entry.msgstr = EN_SITENAME

            print("trans:", entry.msgid, "to:", entry.msgstr)
            entry.flags = [f for f in entry.flags if f != "fuzzy"]

    pofile.save("../../statics/i18n/" + target + "/LC_MESSAGES/messages.po")

if __name__ == "__main__":
    main()
