// Copyright 2022 Teamgram Authors
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: teamgramio (teamgram.io@gmail.com)
//

package core

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/teamgram/proto/mtproto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// NovinRasanVersionFile — مسیر version.json که اسکریپت publish APK آن را می‌نویسد.
// در صورت نیاز با متغیر محیطی NR_VERSION_JSON قابل override است.
const NovinRasanVersionFile = "/var/www/dana-files/novin-rasan/version.json"

type novinVersionInfo struct {
	VersionCode int    `json:"version_code"`
	VersionName string `json:"version_name"`
	ApkUrl      string `json:"apk_url"`
	Size        int64  `json:"size"`
}

func (c *ConfigurationCore) readVersionInfo() *novinVersionInfo {
	path := os.Getenv("NR_VERSION_JSON")
	if path == "" {
		path = NovinRasanVersionFile
	}

	b, err := os.ReadFile(path)
	if err != nil {
		c.Logger.Errorf("help.getAppUpdate: read %s failed: %v", path, err)
		return nil
	}

	var v novinVersionInfo
	if err := json.Unmarshal(b, &v); err != nil {
		c.Logger.Errorf("help.getAppUpdate: parse %s failed: %v", path, err)
		return nil
	}
	if v.VersionName == "" || v.ApkUrl == "" {
		c.Logger.Errorf("help.getAppUpdate: %s missing version_name/apk_url", path)
		return nil
	}
	return &v
}

// HelpGetAppUpdate
// help.getAppUpdate#522d5a7d source:string = help.AppUpdate;
//
// نوین رسان: هندلر upstream خالی بود (NoAppUpdate). اینجا نسخه و لینک APK
// از version.json خوانده می‌شود. اگر نسخهٔ کلاینت قدیمی‌تر باشد، appUpdate
// برگردانده می‌شود (can_not_skip=false تا اختیاری باشد)، وگرنه noAppUpdate.
func (c *ConfigurationCore) HelpGetAppUpdate(in *mtproto.TLHelpGetAppUpdate) (*mtproto.Help_AppUpdate, error) {
	source := ""
	if in != nil {
		source = in.Source
	}

	v := c.readVersionInfo()
	if v == nil {
		// فایل نسخه در دسترس نیست → به کلاینت بگو آپدیتی نیست (رفتار امن)
		return mtproto.MakeTLHelpNoAppUpdate(nil).To_Help_AppUpdate(), nil
	}

	// نسخهٔ کلاینت از source می‌آید (مثل "NovinRasan Android 0.8.48 (1234)")
	// مقایسه بر اساس version_code اگر پیدا شود، وگرنه version_name.
	clientCode := extractVersionCode(source)
	if clientCode > 0 && clientCode >= int32(v.VersionCode) {
		return mtproto.MakeTLHelpNoAppUpdate(nil).To_Help_AppUpdate(), nil
	}
	if clientCode == 0 {
		clientName := extractVersionName(source)
		if clientName != "" && clientName == v.VersionName {
			return mtproto.MakeTLHelpNoAppUpdate(nil).To_Help_AppUpdate(), nil
		}
	}

	c.Logger.Infof("help.getAppUpdate: client source=%q -> offering v%s (code %d)", source, v.VersionName, v.VersionCode)

	rValue := mtproto.MakeTLHelpAppUpdate(&mtproto.Help_AppUpdate{
		CanNotSkip: false,
		Id:         int32(v.VersionCode),
		Version:    v.VersionName,
		Text:       "نسخهٔ جدید نوین رسان در دسترس است. برای دریافت آخرین امکانات و رفع اشکالات، به‌روزرسانی کنید.",
		Entities:   nil,
		Document:   nil,
		Url:        wrapperspb.String(v.ApkUrl),
		Sticker:    nil,
	}).To_Help_AppUpdate()

	return rValue, nil
}

// extractVersionCode — استخراج عدد نسخه از متن source (اولین عدد ۳+ رقمی داخل پرانتز یا انتها).
func extractVersionCode(source string) int32 {
	// الگوی رایج: "App 0.8.48 (1234)" → 1234
	if i := strings.LastIndex(source, "("); i >= 0 {
		if j := strings.Index(source[i:], ")"); j > 1 {
			if n, err := strconv.Atoi(strings.TrimSpace(source[i+1 : i+j])); err == nil {
				return int32(n)
			}
		}
	}
	return 0
}

// extractVersionName — استخراج نام نسخه (مثل 0.8.48) از source.
func extractVersionName(source string) string {
	fields := strings.Fields(source)
	for _, f := range fields {
		parts := strings.Split(f, ".")
		if len(parts) >= 2 {
			allNum := true
			for _, p := range parts {
				if _, err := strconv.Atoi(p); err != nil {
					allNum = false
					break
				}
			}
			if allNum {
				return f
			}
		}
	}
	return ""
}
