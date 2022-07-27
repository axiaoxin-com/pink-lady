cd ./misc/i18n

# 提取html文件翻译词
find ../.. -name "*.html" | xargs perl -pe "s/_i18n .+? \"(.+?)\"/{{ gettext(\"\1\") }}/g" | xgettext --no-wrap --no-location --language=c --from-code=UTF-8 --output=html.pot -
find ../.. -name "*.html" | xargs perl -pe "s/[\"\']/,/g;s/_i18n .+? \`(.+?)\`/{{ gettext(\"\1\") }}/g" | xgettext --no-wrap --no-location --language=c --from-code=UTF-8 --output=html2.pot -
# 提取go文件翻译词
find ../.. -name "*.go" | xargs perl -pe "s/CtxI18n\(c, /gettext(/g"  | xgettext --no-wrap --no-location --language=c --from-code=UTF-8 --output=go.pot -
find ../.. -name "*.go" | xargs perl -pe "s/LangI18n\(.+?, /gettext(/g"  | xgettext --no-wrap --no-location --language=c --from-code=UTF-8 --output=go2.pot -
find ../.. -name "*.go" | xargs perl -pe "s/I18nString/gettext/g"  | xgettext --no-wrap --no-location --language=c --from-code=UTF-8 --output=go3.pot -
# 合并多个pot为一个pot
rm messages.pot
xgettext --no-wrap --no-location *.pot -o messages.pot
# 删除单个的pot文件
rm html.pot html2.pot go.pot go2.pot go3.pot

echo "====> 合并pot到po"
init="update"
if test ! -e ../../statics/i18n/en/LC_MESSAGES/messages.po; then
    ./po_init.sh
    init="init"
else
    msgmerge --no-wrap ../../statics/i18n/en/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/en/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/de/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/de/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/es/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/es/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/fr/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/fr/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/it/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/it/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/ja/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/ja/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/ko/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/ko/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/pt/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/pt/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/ru/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/ru/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/tr/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/tr/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/zh-Hant/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/zh-Hant/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/vi/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/vi/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/ar/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/ar/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/hi/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/hi/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/bn/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/bn/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/id/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/id/LC_MESSAGES/messages.po
    msgmerge --no-wrap ../../statics/i18n/th/LC_MESSAGES/messages.po messages.pot -o ../../statics/i18n/th/LC_MESSAGES/messages.po
fi


echo "====> 谷歌自动翻译"
echo "===>" $init

. ./venv/bin/activate

transed="false"

msgidsCount=`(python ./po_trans.py en $init|tail -1)`
echo "en $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py en
    transed="true"
fi

msgidsCount=`(python ./po_trans.py de $init|tail -1)`
echo "de $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py de
    transed="true"
fi

msgidsCount=`(python ./po_trans.py es $init|tail -1)`
echo "es $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py es
    transed="true"
fi

msgidsCount=`(python ./po_trans.py fr $init|tail -1)`
echo "fr $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py fr
    transed="true"
fi

msgidsCount=`(python ./po_trans.py it $init|tail -1)`
echo "it $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py it
    transed="true"
fi

msgidsCount=`(python ./po_trans.py ja $init|tail -1)`
echo "ja $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py ja
    transed="true"
fi

msgidsCount=`(python ./po_trans.py ko $init|tail -1)`
echo "ko $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py ko
    transed="true"
fi

msgidsCount=`(python ./po_trans.py pt $init|tail -1)`
echo "pt $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py pt
    transed="true"
fi

msgidsCount=`(python ./po_trans.py ru $init|tail -1)`
echo "ru $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py ru
    transed="true"
fi

msgidsCount=`(python ./po_trans.py tr $init|tail -1)`
echo "tr $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py tr
    transed="true"
fi

msgidsCount=`(python ./po_trans.py zh-Hant $init|tail -1)`
echo "zh-Hant $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py zh-Hant
    transed="true"
fi

msgidsCount=`(python ./po_trans.py vi $init|tail -1)`
echo "vi $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py vi
    transed="true"
fi

msgidsCount=`(python ./po_trans.py ar $init|tail -1)`
echo "ar $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py ar
    transed="true"
fi

msgidsCount=`(python ./po_trans.py hi $init|tail -1)`
echo "hi $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py hi
    transed="true"
fi

msgidsCount=`(python ./po_trans.py bn $init|tail -1)`
echo "bn $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py bn
    transed="true"
fi

msgidsCount=`(python ./po_trans.py id $init|tail -1)`
echo "id $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py id
    transed="true"
fi

msgidsCount=`(python ./po_trans.py th $init|tail -1)`
echo "th $msgidsCount"
if [ "$msgidsCount" != "msgids count: 0" ]; then
    python ./fix_msgstr.py th
    transed="true"
fi


echo "====> po内容清理"
if [ `uname` = 'Darwin' ]; then
    sed -i  '' '/^#~ .*/d' `grep "#~ " -rl ../../statics/i18n/*/LC_MESSAGES/messages.po`
else
    sed -i '/^#~ .*/d' `grep "#~ " -rl ../../statics/i18n/*/LC_MESSAGES/messages.po`
fi

if [ "$transed" = "true" ]; then
    exit 0
else
    exit 5
fi
