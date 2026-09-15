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
// نوین رسان: این هندلر upstream خالی بود (JSONObject بدون کلید) و کلاینت
// تنظیمات حیاتی را نمی‌گرفت. اینجا یک config واقعی و سبک برگردانده می‌شود.

package core

import (
	"github.com/teamgram/proto/mtproto"
)

// نوین رسان: مقادیر پیکربندی اپ (کلیدهای پرکاربرد در MessagesController/...).
// فقط کلیدهایی که اپ واقعاً می‌خواند و نبودشان رفتار ناقص می‌سازد.
func novinAppConfig() []*mtproto.JSONObjectValue {
	str := func(k, v string) *mtproto.JSONObjectValue {
		return mtproto.MakeTLJsonObjectValue(&mtproto.JSONObjectValue{
			Key:   k,
			Value: mtproto.MakeTLJsonString(&mtproto.JSONValue{Value_STRING: v}).To_JSONValue(),
		}).To_JSONObjectValue()
	}
	num := func(k string, v float64) *mtproto.JSONObjectValue {
		return mtproto.MakeTLJsonObjectValue(&mtproto.JSONObjectValue{
			Key:   k,
			Value: mtproto.MakeTLJsonNumber(&mtproto.JSONValue{Value_FLOAT64: v}).To_JSONValue(),
		}).To_JSONObjectValue()
	}
	bl := func(k string, v bool) *mtproto.JSONObjectValue {
		var bv *mtproto.Bool
		if v {
			bv = mtproto.MakeTLBoolTrue(nil).To_Bool()
		} else {
			bv = mtproto.MakeTLBoolFalse(nil).To_Bool()
		}
		return mtproto.MakeTLJsonObjectValue(&mtproto.JSONObjectValue{
			Key:   k,
			Value: mtproto.MakeTLJsonBool(&mtproto.JSONValue{Value_BOOL: bv}).To_JSONValue(),
		}).To_JSONObjectValue()
	}
	obj := func(k string, v []*mtproto.JSONObjectValue) *mtproto.JSONObjectValue {
		return mtproto.MakeTLJsonObjectValue(&mtproto.JSONObjectValue{
			Key: k,
			Value: mtproto.MakeTLJsonObject(&mtproto.JSONValue{
				Value_VECTORJSONOBJECTVALUE: v,
			}).To_JSONValue(),
		}).To_JSONObjectValue()
	}

	return []*mtproto.JSONObjectValue{
		// —— ورود/احراز هویت ——
		bl("qr_login_camera", true),
		bl("qr_login_code", true),
		num("authorization_autoconfirm_period", 604800),
		bl("android_check_reset_langpack", false),
		bl("android_collect_device_stats", false),

		// —— چت/گروه/کانال (حدود واقعی) ——
		num("channels_limit_default", 500),
		num("channels_limit_premium", 1000),
		num("channels_public_limit_default", 10),
		num("channels_public_limit_premium", 20),
		num("chatlist_invites_limit_default", 3),
		num("chatlist_invites_limit_premium", 100),
		num("chatlists_joined_limit_default", 2),
		num("chatlists_joined_limit_premium", 20),
		num("chat_pinned_count", 5),
		bl("autoarchive_setting_available", true),

		// —— پیام/مدیا ——
		num("caption_length_limit_default", 1024),
		num("caption_length_limit_premium", 4096),
		num("about_length_limit_default", 70),
		num("about_length_limit_premium", 140),
		num("message_length_limit", 4096),
		num("audio_bitrate", 64000),
		str("gif_search_username", "gif"),
		str("venue_search_username", "foursquare"),
		str("img_search_username", "bing"),

		// —— ایموجی/دایس (بازی‌های داخلی) ——
		bl("emojies_sounds", true),
		bl("emojies_animated_zoom", true),
		bl("emojies_send_dice", true),
		obj("emojies_send_dice_success", []*mtproto.JSONObjectValue{
			num("🎲", 6), num("🎯", 6), num("🏀", 5),
			num("⚽", 5), num("🎰", 64), num("🎳", 6),
		}),

		// —— سشن/امنیت ——
		bl("chat_read_mark", true),
		num("chat_read_mark_size_threshold", 50),
		num("chat_read_mark_expire_period", 604800),
		bl("call_requests_disabled", false),
		bl("background_connection", true),

		// —— پریمیوم: غیرفعال (این سرور پریمیوم ندارد) ——
		bl("premium_purchase_blocked", true),
		num("premium_gift_attach_menu_icon", 0),
		num("premium_gift_text_field_icon", 0),
		str("premium_manage_subscription_url", ""),
		str("premium_bot_username", ""),
	}
}

// HelpGetAppConfig61E3F854
// help.getAppConfig#61e3f854 hash:int = help.AppConfig;
func (c *ConfigurationCore) HelpGetAppConfig61E3F854(in *mtproto.TLHelpGetAppConfig61E3F854) (*mtproto.Help_AppConfig, error) {
	_ = in
	return mtproto.MakeTLHelpAppConfig(&mtproto.Help_AppConfig{
		Hash: 0,
		Config: mtproto.MakeTLJsonObject(&mtproto.JSONValue{
			Value_VECTORJSONOBJECTVALUE: novinAppConfig(),
		}).To_JSONValue(),
	}).To_Help_AppConfig(), nil
}
