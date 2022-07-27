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

transcli = Translate()


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
        _m = (
            entry.msgid.replace("赛可心理测试", "PsycTest")
            .replace("PsycTestMBTI", "PsycTest MBTI")
            .replace("# ", "#")
            .replace("?wx_fmt=png&amp;wxfrom=5&amp;wx_lazy=1&amp;wx_co=1", "")
        )
        # 删除markdown链接描述
        _m = re.sub(r"\[(.*?)\]\s?\((.+?)\)", lambda m: "[](" + m.group(2) + ")", _m)

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
                msgs[mid] = (
                    trans[idx]
                    .translatedText.replace("&quot;", '"')
                    .replace("&#39;", "'")
                    .replace("&gt;", ">")
                )
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
            print("trans:", entry.msgid, "to:", entry.msgstr)
            entry.flags = [f for f in entry.flags if f != "fuzzy"]

    pofile.save("../../statics/i18n/" + target + "/LC_MESSAGES/messages.po")


if __name__ == "__main__":
    main()
