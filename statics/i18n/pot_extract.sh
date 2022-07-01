# 提取html文件翻译词
find ../.. -name "*.html" | xargs perl -pe "s/{{\s*_text [\"\`](.+?)[\"\`]\s*}}/{{ gettext(\"\1\") }}/g" | xgettext --no-wrap --no-location --language=c --from-code=UTF-8 --output=html.pot -
# 提取go文件翻译词
find ../.. -name "*.go" | xargs perl -pe "s/gettext.Gettext/gettext/g"  | xgettext --no-wrap --no-location --language=c --from-code=UTF-8 --output=go.pot -
# 提取js文件翻译词
find ../.. -name "custom.js" | xargs perl -pe "s/String\([\"\'](.+?)[\"\']\)/gettext(\"\1\")/g" | xgettext --no-wrap --no-location --language=c --from-code=UTF-8 --output=js.pot -
# 合并多个pot为一个pot
rm messages.pot
xgettext --no-wrap --no-location *.pot -o messages.pot
# 删除单个的pot文件
rm html.pot go.pot js.pot
