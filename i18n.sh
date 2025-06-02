cd ./misc/i18n

# 提取翻译词生成messages.pot
if ! python ./extract_i18n.py; then
    echo "Error: Failed to extract i18n messages"
    exit 1
fi

# 定义所有支持的语言列表
LANGUAGES=(
    "en"
    "ja"
    "ko"
    "es"
    "pt"
    "fr"
    "de"
    "it"
    "ru"
    "tr"
    "zh-Hant"
    "vi"
    "ar"
    "hi"
    "bn"
    "id"
    "th"
)

echo "====> 合并pot到po"
if test ! -e ../../statics/i18n/en/LC_MESSAGES/messages.po; then
    # 首次创建po文件
    for lang in "${LANGUAGES[@]}"; do
        mkdir -p "../../statics/i18n/${lang}/LC_MESSAGES"
        if ! msginit --no-translator --no-wrap --input=messages.pot --local="${lang}" -o "../../statics/i18n/${lang}/LC_MESSAGES/messages.po"; then
            echo "Error: Failed to initialize po file for $lang"
            exit 1
        fi
    done
else
    # 更新已存在的po文件
    for lang in "${LANGUAGES[@]}"; do
        if ! msgmerge --no-wrap "../../statics/i18n/${lang}/LC_MESSAGES/messages.po" messages.pot -o "../../statics/i18n/${lang}/LC_MESSAGES/messages.po"; then
            echo "Error: Failed to merge pot file for $lang"
            exit 1
        fi
    done
fi


echo "====> 谷歌自动翻译"
failed_langs=()
success_langs=()

# 使用LANGUAGES数组遍历所有语言
for lang in "${LANGUAGES[@]}"; do
    echo "Processing $lang..."
    if ! python ./po_trans.py "$lang"; then
        failed_langs+=("$lang")
        echo "Error: Translation failed for $lang"
    else
        success_langs+=("$lang")
        echo "Successfully processed $lang"
    fi
done

# 报告翻译结果
echo "====> 翻译结果报告"
echo "成功翻译的语言: ${success_langs[*]}"
if [ ${#failed_langs[@]} -gt 0 ]; then
    echo "翻译失败的语言: ${failed_langs[*]}"
    echo "警告: 部分语言翻译失败，请检查上述失败的语言"
fi

echo "====> po内容清理"
if [ `uname` = 'Darwin' ]; then
    sed -i  '' '/^#~ .*/d' `grep "#~ " -rl ../../statics/i18n/*/LC_MESSAGES/messages.po`
else
    sed -i '/^#~ .*/d' `grep "#~ " -rl ../../statics/i18n/*/LC_MESSAGES/messages.po`
fi

# 如果有失败的语言，返回非零状态码，但不中断执行
if [ ${#failed_langs[@]} -gt 0 ]; then
    exit 5
fi

exit 0
