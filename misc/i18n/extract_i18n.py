#!/usr/bin/env python3

import os
import re
import importlib
import importlib.util
from pathlib import Path
from typing import List, Set, Dict, Any, Callable

import tomli
from babel.messages import Catalog, Message
from babel.messages.pofile import write_po


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
        title = ""
        for line in lines:
            if line.startswith('# '):
                title = line[2:].strip()
                if title:
                    strings.add(title)
                break

        # Extract description from the first blockquote
        description = ""
        for line in lines:
            if line.startswith('> '):
                description = line[2:].strip()
                if description:  # 只有非空的描述才添加
                    strings.add(description)
                break
    return strings

def extract_config_nav_names(config_path: str) -> Set[str]:
    """Extract nav_name values from flatpages configuration in TOML file."""
    strings = set()
    try:
        with open(config_path, 'rb') as f:
            config = tomli.load(f)
            if 'flatpages' in config and 'dirs' in config['flatpages']:
                for dir_config in config['flatpages']['dirs']:
                    if 'nav_name' in dir_config:
                        strings.add(dir_config['nav_name'])
    except Exception as e:
        print(f"Error reading config file: {e}")
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

def load_plugins(plugins_dir: str) -> Dict[str, Callable[[], Set[str]]]:
    """Load custom string extraction plugins from the plugins directory."""
    plugins = {}
    plugins_path = Path(plugins_dir)
    
    if not plugins_path.exists():
        return plugins
        
    for plugin_file in plugins_path.glob('*.py'):
        if plugin_file.name == '__init__.py':
            continue
            
        try:
            # Load the plugin module
            spec = importlib.util.spec_from_file_location(plugin_file.stem, plugin_file)
            if spec is None or spec.loader is None:
                continue
                
            module = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(module)
            
            # Check if the module has the required extract_strings function
            if hasattr(module, 'extract_strings'):
                plugins[plugin_file.stem] = module.extract_strings
                print(f"Loaded plugin: {plugin_file.stem}")
        except Exception as e:
            print(f"Error loading plugin {plugin_file.name}: {e}")
            
    return plugins

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

    # Extract nav_name from config file
    config_filename = 'config.default.toml'
    config_path = os.path.join(base_dir, config_filename)
    strings.update(extract_config_nav_names(config_path))

    # Load and run custom plugins
    plugins_dir = os.path.join(os.path.dirname(__file__), 'plugins')
    plugins = load_plugins(plugins_dir)
    for plugin_name, extract_func in plugins.items():
        try:
            plugin_strings = extract_func()
            strings.update(plugin_strings)
            print(f"Plugin {plugin_name} extracted {len(plugin_strings)} strings")
        except Exception as e:
            print(f"Error running plugin {plugin_name}: {e}")

    # Generate messages.pot
    output_file = os.path.join(os.path.dirname(__file__), 'messages.pot')
    write_pot_file(strings, output_file)
    print(f"Generated {output_file} with {len(strings)} unique strings")

if __name__ == "__main__":
    main()
