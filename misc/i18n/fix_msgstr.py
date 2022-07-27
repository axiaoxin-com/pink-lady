#! /usr/bin/env python3
#
# ENV Install:
# virtualenv venv --python=python3.8
# source ./venv/bin/activate
# ./venv/bin/pip-3.8 install -r requirements.txt

import re
import sys

import polib
from pygtrans import Translate


def main():
    target = "en"
    if sys.argv[1]:
        target = sys.argv[1]

    pofile = polib.pofile("../../statics/i18n/" + target + "/LC_MESSAGES/messages.po")
    for entry in pofile:
        # 标题换行修复
        entry.msgstr = re.sub(r"(?<!\n\n)(#+)\s+", r"\n\n\g<1>", entry.msgstr)
        entry.msgstr = entry.msgstr.replace("#\n\n#", "##")
        # 链接/图片修复
        entry.msgstr = fix_links(entry.msgstr)
        # 删除多余的\n
        entry.msgstr = re.sub("\n\n\n+", "\n\n", entry.msgstr)

    pofile.save("../../statics/i18n/" + target + "/LC_MESSAGES/messages.po")


def fix_links(text):
    # 匹配markdown链接
    pattern = r"\[(.*?)\]\s?\((.+?)\)"
    # 删除markdown url中的空格
    text = re.sub(
        pattern, lambda m: "[" + m.group(1) + "]" + "(" + m.group(2).replace(" ", "") + ")", text
    )
    # 图片! []() fix
    pattern = r"\!\s?\[(.*?)\]\s?\((.+?)\)"
    text = re.sub(
        pattern, lambda m: "![" + m.group(1) + "]" + "(" + m.group(2).replace(" ", "") + ")", text
    )
    # 图片换行
    text = re.sub(r"(?<!\n\n)\!\[.*?\]\(.+?\)", r"\n\n\g<0>\n\n", text)
    return text


if __name__ == "__main__":
    main()
