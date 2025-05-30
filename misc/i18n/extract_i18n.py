#!/usr/bin/env python3

import os
import re
from pathlib import Path
from typing import List, Set
from babel.messages import Catalog, Message
from babel.messages.pofile import write_po
import markdown2

def find_files(base_dir: str, extensions: List[str]) -> List[str]:
    """Recursively find files with given extensions."""
    files = []
    for ext in extensions:
        files.extend([str(p) for p in Path(base_dir).rglob(f"*.{ext}")])
    return files

def extract_html_strings(file_path: str) -> Set[str]:
    """Extract strings from HTML files marked with _i18n."""
    strings = set()
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
        
        # Extract strings with double quotes: _i18n $lang "xxx"
        matches = re.findall(r'_i18n\s+\S+\s+"([^"]+)"', content)
        strings.update(matches)
        
        # Extract strings with backticks: _i18n $lang `xxx`
        matches = re.findall(r'_i18n\s+\S+\s+`([^`]+)`', content)
        strings.update(matches)
    
    return strings

def extract_go_strings(file_path: str) -> Set[str]:
    """Extract strings from Go files marked for translation."""
    strings = set()
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
        
        # Extract CtxI18n(c, "xxx")
        matches = re.findall(r'CtxI18n\([^,]+,\s*["`]([^"`]+)["`]', content)
        strings.update(matches)
        
        # Extract LangI18n("en", "xxx")
        matches = re.findall(r'LangI18n\([^,]+,\s*["`]([^"`]+)["`]', content)
        strings.update(matches)
        
        # Extract I18nString("xxx")
        matches = re.findall(r'I18nString\(["`]([^"`]+)["`]', content)
        strings.update(matches)
    
    return strings

def extract_markdown_strings(file_path: str) -> Set[str]:
    """Extract title, description and content from markdown files for translation."""
    strings = set()
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
        if not content: 
            return strings
        strings.add(content)

        lines = content.split('\n')
        
        # Extract title from first h1 heading
        title = "Untitled Flatpage"  # default title
        for line in lines:
            if line.startswith('# '):
                title = line[2:].strip()
                break
        strings.add(title)
        
        # Extract description from the first blockquote
        description = ""
        for line in lines:
            if line.startswith('> '):
                description = line[2:].strip()
                if description:  # 只有非空的描述才添加
                    strings.add(description)
                break   
    return strings

def write_pot_file(strings: Set[str], output_file: str):
    """Write strings to POT file using babel."""
    # Create a new catalog
    catalog = Catalog(
        project='pink-lady',
        version='1.0',
        copyright_holder='pink-lady',
        charset='utf-8'
    )
    
    # Add each string as a message to the catalog
    for string in sorted(strings):
        catalog.add(string, locations=[])
    
    # Write the catalog to POT file
    with open(output_file, 'wb') as f:
        write_po(f, catalog, width=0, no_location=True, omit_header=False)

def main():
    # Get base directory (two levels up from script location)
    base_dir = str(Path(__file__).parent.parent.parent)
    
    # Find all HTML and Go files
    html_files = find_files(base_dir, ['html'])
    go_files = find_files(base_dir, ['go'])
    
    # Find flatpages markdown files in default config directory
    flatpages_path = 'statics/flatpages'
    md_files = list(Path(base_dir).joinpath(flatpages_path).rglob('*.md'))
    
    # Extract strings
    strings = set()
    
    for html_file in html_files:
        strings.update(extract_html_strings(html_file))
    
    for go_file in go_files:
        strings.update(extract_go_strings(go_file))
        
    for md_file in md_files:
        strings.update(extract_markdown_strings(str(md_file)))
    
    # Generate messages.pot
    output_file = os.path.join(os.path.dirname(__file__), 'messages.pot')
    write_pot_file(strings, output_file)
    print(f"Generated {output_file} with {len(strings)} unique strings")

if __name__ == "__main__":
    main() 