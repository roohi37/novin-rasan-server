// Copyright 2024 Teamgram Authors
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
	"github.com/teamgram/proto/mtproto"
)

// HelpGetPassportConfig
// help.getPassportConfig#c661ad08 hash:int = help.PassportConfig;
func (c *PassportCore) HelpGetPassportConfig(in *mtproto.TLHelpGetPassportConfig) (*mtproto.Help_PassportConfig, error) {
	// نوین رسان: upstream استاب بود → پاسخ خالیِ معتبر (بدون خطای Enterprise).
	_ = in
	// passport در این سرور فعال نیست → پیکربندی خالی معتبر
	return mtproto.MakeTLHelpPassportConfig(&mtproto.Help_PassportConfig{
		Hash:           0,
		CountriesLangs: mtproto.MakeTLDataJSON(&mtproto.DataJSON{Data: "{\"countries_langs\":{}}"}).To_DataJSON(),
	}).To_Help_PassportConfig(), nil
}
