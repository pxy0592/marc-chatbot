# Theme State Model

## Theme

```text
light | dark
```

## ThemePreference

- Storage key: `marc-chatbot-theme`
- Valid values: `light`, `dark`
- Invalid/missing value: use system preference, then light fallback

## State flow

```text
load page
  -> valid saved preference? use it
  -> otherwise system prefers dark? dark
  -> otherwise light
  -> apply data-theme to document root
  -> user toggles
  -> update React state + document root + localStorage
```
