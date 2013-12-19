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
	"dbus/org/freedesktop/accounts"
	"dlib/dbus"
)

func (u *User) SetPassword(passwd, hint string) {
	u.UserInface.SetPassword(passwd, hint)
}

func (u *User) OnPropertiesChanged(propName string, old interface{}) {
	switch propName {
	case "AccountType":
		if v, ok := old.(int32); ok {
			if v != u.AccountType {
				u.UserInface.SetAccountType(u.AccountType)
			}
		}
	case "AutomaticLogin":
		if v, ok := old.(bool); ok {
			if v != u.AutomaticLogin {
				u.UserInface.SetAutomaticLogin(u.AutomaticLogin)
			}
		}
	case "IconFile":
		if v, ok := old.(string); ok {
			if v != u.IconFile {
				u.UserInface.SetIconFile(u.IconFile)
			}
		}
	case "Locked":
		if v, ok := old.(bool); ok {
			if v != u.Locked {
				u.UserInface.SetLocked(u.Locked)
			}
		}
	case "PasswordMode":
		if v, ok := old.(int32); ok {
			if v != u.PasswordMode {
				u.UserInface.SetPasswordMode(u.PasswordMode)
			}
		}
	case "UserName":
		if v, ok := old.(string); ok {
			if v != u.UserName {
				u.UserInface.SetUserName(u.UserName)
			}
		}
	}
}

func NewAccountUserManager(path string) *User {
	u := &User{}

	u.ObjectPath = path
	u.UserInface = accounts.GetUser(path)

	GetUserProperties(u)
	u.UserInface.ConnectChanged(func() {
		tmpUser := &User{}
		tmpUser.ObjectPath = u.ObjectPath
		tmpUser.UserInface = u.UserInface
		GetUserProperties(tmpUser)
		CompareUserManager(u, tmpUser)
	})

	_userMap[path] = u
	return u
}

func GetUserProperties(u *User) {
	userInface := u.UserInface
	if userInface == nil {
		return
	}
	u.AccountType = userInface.AccountType.Get()
	u.AutomaticLogin = userInface.AutomaticLogin.Get()
	u.IconFile = userInface.IconFile.Get()
	u.Locked = userInface.Locked.Get()
	u.LoginTime = userInface.LoginTime.Get()
	u.PasswordMode = userInface.PasswordMode.Get()
	u.UserName = userInface.UserName.Get()
}

func CompareUserManager(src, tmp *User) {
	if src == nil || tmp == nil {
		return
	}

	if src.AccountType != tmp.AccountType {
		src.AccountType = tmp.AccountType
		dbus.NotifyChange(src, "AccountType")
	}

	if src.AutomaticLogin != tmp.AutomaticLogin {
		src.AutomaticLogin = tmp.AutomaticLogin
		dbus.NotifyChange(src, "AutomaticLogin")
	}

	if src.IconFile != tmp.IconFile {
		src.IconFile = tmp.IconFile
		dbus.NotifyChange(src, "IconFile")
	}

	if src.Locked != tmp.Locked {
		src.Locked = tmp.Locked
		dbus.NotifyChange(src, "Locked")
	}

	if src.LoginTime != tmp.LoginTime {
		src.LoginTime = tmp.LoginTime
		dbus.NotifyChange(src, "LoginTime")
	}

	if src.PasswordMode != tmp.PasswordMode {
		src.PasswordMode = tmp.PasswordMode
		dbus.NotifyChange(src, "PasswordMode")
	}

	if src.UserName != tmp.UserName {
		src.UserName = tmp.UserName
		dbus.NotifyChange(src, "UserName")
	}
}

func DeleteUserManager(path string) {
	u := _userMap[path]
	if u != nil {
		return
	}

	accounts.DestroyUser(u.UserInface)
	dbus.UnInstallObject(u)
}
