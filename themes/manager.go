/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
        "dlib/glib-2.0"
        "io/ioutil"
        "os"
        "path"
)

type Manager struct {
        ThemeList    []string
        CurrentTheme string  `access:"readwrite"`
        GtkThemeList []string
        //GtkBasePath     string
        IconThemeList []string
        //IconBasePath    string
        CursorThemeList []string
        //CursorBasePath  string
        FontThemeList  []string
        BackgroundList []string
        SoundThemeList []string
        pathNameMap    map[string]PathInfo
}

func (op *Manager) GetPathViaName(name string) (string, bool) {
        if obj := op.getThemeObject(name); obj != nil {
                return obj.path, true
        }

        return "", false
}

func (op *Manager) SetGtkTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        setGtkThemeViaXSettings(name)
        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(name, obj.IconTheme, obj.GtkCursorTheme,
                        obj.GtkFontName, obj.BackgroundFile, obj.SoundThemeName)
                op.updateCurrentTheme(v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetIconTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        setIconThemeViaXSettings(name)
        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, name,
                        obj.GtkCursorTheme, obj.GtkFontName, obj.BackgroundFile, obj.SoundThemeName)
                op.updateCurrentTheme(v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetGtkCursorTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        setGtkCursorThemeViaXSettings(name)
        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, obj.IconTheme,
                        name, obj.GtkFontName, obj.BackgroundFile, obj.SoundThemeName)
                op.updateCurrentTheme(v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetGtkFontTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        setGtkFontThemeViaXSettings(name)
        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, obj.IconTheme,
                        obj.GtkCursorTheme, name, obj.BackgroundFile, obj.SoundThemeName)
                op.updateCurrentTheme(v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetBackgroundFile(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, obj.IconTheme,
                        obj.GtkCursorTheme, obj.GtkFontName, name, obj.SoundThemeName)
                op.updateCurrentTheme(v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) SetSoundTheme(name string) (string, bool) {
        if len(name) <= 0 {
                return op.CurrentTheme, false
        }

        if obj := op.getThemeObject(op.CurrentTheme); obj != nil {
                v := op.setTheme(obj.GtkTheme, obj.IconTheme,
                        obj.GtkCursorTheme, obj.GtkFontName, obj.BackgroundFile, name)
                op.updateCurrentTheme(v)
                return v, true
        }

        return op.CurrentTheme, false
}

func (op *Manager) setTheme(gtk, icon, cursor, gtkFont, bg, sound string) string {
        for _, path := range op.ThemeList {
                name, ok := isThemeExist(gtk, icon, cursor, gtkFont, bg, sound, path)
                if !ok {
                        continue
                } else {
                        return name
                }
        }

        createTheme("Custom", gtk, icon, cursor, gtkFont, bg, sound)
        updateThemeObj(op.pathNameMap)

        return "Custom"
}

func getThemeList() []PathInfo {
        return getThemeThumbList()
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getGtkThemeList() []PathInfo {
        valid := getValidGtkThemes()
        thumb := getGtkThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isElementExist(v, thumb) {
                        list = append(list, v)
                }
        }

        return list
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getIconThemeList() []PathInfo {
        valid := getValidIconThemes()
        thumb := getIconThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isElementExist(v, thumb) {
                        list = append(list, v)
                }
        }

        return list
}

/*
   Return only contains thumbnail theme
   First, get all local themes
   Then, determine whether contains thumbnail
*/
func getCursorThemeList() []PathInfo {
        valid := getValidCursorThemes()
        thumb := getCursorThumbList()

        list := []PathInfo{}
        for _, v := range valid {
                if isElementExist(v, thumb) {
                        list = append(list, v)
                }
        }

        return list
}

// Has not yet been determined
func getFontNameList() []string {
        return []string{}
}

/*
   Unlimited
   Return all bg.
*/
func getBackgroundList() []string {
        return []string{}
}

/*
   Return all sound theme names.
*/
func getSoundThemeList() []string {
        list := []string{}
        files, err := ioutil.ReadDir(SOUND_THEME_PATH)
        if err != nil {
                logObject.Error("%v", err)
                return list
        }

        for _, f := range files {
                if f.IsDir() {
                        // check if index.theme file exists
                        tmpf, err := os.Open(path.Join(SOUND_THEME_PATH, f.Name(), SOUND_THEME_MAIN_FILE))
                        if err != nil {
                                continue
                        }
                        defer tmpf.Close()
                        tmpi, err := tmpf.Stat()
                        if err != nil {
                                continue
                        }
                        if !tmpi.IsDir() {
                                list = append(list, f.Name())
                        }
                }
        }

        return list
}

func isThemeExist(gtk, icon, cursor, gtkFont, bg, sound, path string) (string, bool) {
        obj, ok := themeObjMap[path]
        if !ok {
                return "", false
        }

        if gtk != obj.GtkTheme || icon != obj.IconTheme ||
                cursor != obj.GtkCursorTheme || gtkFont != obj.GtkFontName ||
                obj.BackgroundFile != bg || obj.SoundThemeName != sound {
                return "", false
        }

        return obj.Name, true
}

func createTheme(name, gtk, icon, cursor, gtkFont, bg, sound string) bool {
        homeDir := getHomeDir()
        path := homeDir + THUMB_LOCAL_THEME_PATH + "/" + name
        logObject.Info("Theme Dir:%s\n", path)
        if ok, _ := objUtil.IsFileExist(path); !ok {
                logObject.Info("Create Theme Dir: %s\n", path)
                err := os.MkdirAll(path, 0755)
                if err != nil {
                        logObject.Info("Mkdir '%s' failed: %v\n", path, err)
                        return false
                }
        }

        filename := path + "/" + "theme.ini"
        logObject.Info("Theme Config File:%s\n", filename)
        if ok, _ := objUtil.IsFileExist(filename); !ok {
                logObject.Info("Create Theme Config File: %s\n", filename)
                f, err := os.Create(filename)
                if err != nil {
                        logObject.Info("Create '%s' failed: %v\n",
                                filename, err)
                        return false
                }
                f.Close()
        }

        mutex.Lock()
        defer mutex.Unlock()
        keyFile := glib.NewKeyFile()
        defer keyFile.Free()
        ok, err := keyFile.LoadFromFile(filename, glib.KeyFileFlagsKeepComments)
        if !ok {
                logObject.Warning("LoadKeyFile '%s' failed\n", filename)
                return false
        }

        keyFile.SetString(THEME_GROUP_THEME, THEME_KEY_NAME, name)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_GTK, gtk)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_ICONS, icon)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_CURSOR, cursor)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_GTK_FONT, gtkFont)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_BG, bg)
        keyFile.SetString(THEME_GROUP_COMPONENT, THEME_KEY_SOUND, sound)

        _, contents, err1 := keyFile.ToData()
        if err1 != nil {
                logObject.Warning("KeyFile '%s' ToData failed: %s\n",
                        filename, err)
                return false
        }

        writeDatasToKeyFile(contents, filename)

        return true
}

func writeDatasToKeyFile(contents, filename string) {
        if len(filename) <= 0 {
                return
        }

        f, err := os.Create(filename + "~")
        if err != nil {
                logObject.Warning("OpenFile '%s' failed: %v\n",
                        filename+"~", err)
                return
        }
        defer f.Close()

        if _, err = f.WriteString(contents); err != nil {
                logObject.Warning("Write in '%s' failed: %v\n",
                        filename+"~", err)
                return
        }
        f.Sync()
        os.Rename(filename+"~", filename)
}

func newManager() *Manager {
        m := &Manager{}

        m.pathNameMap = make(map[string]PathInfo)
        m.setPropName("ThemeList")
        m.setPropName("CurrentTheme")
        m.setPropName("GtkThemeList")
        m.setPropName("IconThemeList")
        m.setPropName("CursorThemeList")
        m.setPropName("SoundThemeList")

        m.listenSettingsChanged()
        homeDir := getHomeDir()
        m.listenThemeDir(THEMES_PATH)
        m.listenThemeDir(homeDir + THEMES_LOCAL_PATH)
        m.listenThemeDir(ICONS_PATH)
        m.listenThemeDir(homeDir + ICONS_LOCAL_PATH)
        m.listenThemeDir(THUMB_BASE_PATH)
        m.listenThemeDir(homeDir + THUMB_LOCAL_BASE_PATH)

        return m
}
