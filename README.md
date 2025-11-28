# Windows rule
```
windowrulev2 = float, class:^(algo)$
windowrulev2 = pin, class:^(algo)$
windowrulev2 = center, class:^(algo)$
windowrulev2 = size 400 360, class:^(algo)$
```

# algo.conf (Example)
```
Browser=google-chrome
Terminal=alacritty
Theme=rose-pine
```

# theme.conf (Example)
```
selected_bg=#403d52
selected_text=#e0def4
indicator=#c4a7e7
text_color=#908caa
input_text=#e0def4
input_placeholder=#808080
input_prompt=#c4a7e7
```

# custom file (example)
```
Name=Monitor|Exec=hx ~/.config/hypr/monitor.conf|Terminal=true
Name=Windows|Exec=hx ~/.config/hypr/windows.conf|Terminal=true
Name=Input|Exec=hx ~/.config/hypr/input.conf|Terminal=true
Name=Look and Feel|Exec=hx ~/.config/hypr/look-and-feel.conf|Terminal=true
Name=Keybindings|Exec=hx ~/.config/hypr/keybindings.conf|Terminal=true
Name=Hyprpaper|Exec=hx ~/.config/hypr/hyprpaper.conf|Terminal=true
Name=Default Programs|Exec=hx ~/.config/hypr/default-programs.conf|Terminal=true
Name=Alacritty|Exec=hx ~/.config/alacritty/alacritty.toml|Terminal=true
NName=Alias|Exec=hx ~/.config/bash-config/aliases|Terminal=true
```
